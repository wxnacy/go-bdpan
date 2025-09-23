package bdpan

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/wxnacy/go-tools"
)

func NewGetFileListReq() *GetFileListReq {
	return &GetFileListReq{
		Dir:   "/",
		Web:   1,
		Start: 0,
		Limit: 1000,
		Order: "name",
	}
}

type GetFileListReq struct {
	Dir   string
	Web   int
	page  int
	Start int
	Limit int32
	// 排序字段：默认为name；
	// time表示先按文件类型排序，后按修改时间排序；
	// name表示先按文件类型排序，后按文件名称排序；
	// size表示先按文件类型排序，后按文件大小排序。
	Order string
	// 默认为升序，设置为1实现降序 （注：排序的对象是当前目录下所有文件，不是当前分页下的文件）
	Desc int32
}

func (r *GetFileListReq) SetLimit(limit int32) *GetFileListReq {
	r.Limit = limit
	return r
}

func (r *GetFileListReq) SetPage(page int) *GetFileListReq {
	r.Start = (page - 1) * int(r.Limit)
	return r
}

func (r *GetFileListReq) SetDir(dir string) *GetFileListReq {
	r.Dir = dir
	return r
}

type GetFileListRes struct {
	GuidInfo string      `json:"guid_info"`
	Errmsg   string      `json:"errmsg"`
	List     []*FileInfo `json:"list"`
}

func NewGetFileInfoReq(fsid uint64) *GetFileInfoReq {
	return &GetFileInfoReq{
		FSID: fsid,
		// 是否需要下载地址，0为否，1为是，默认为0。获取到dlink后，参考下载文档进行下载操作
		// 下载文档: https://pan.baidu.com/union/doc/pkuo3snyp
		Dlink: 1,
	}
}

type GetFileInfoReq struct {
	Dlink int
	FSID  uint64
}

type GetFileInfoRes struct {
	FileInfo
}

func NewBatchGetFileListReq(fsid uint64) *BatchGetFileInfoReq {
	return &BatchGetFileInfoReq{
		FSIds: []uint64{fsid},
		// 是否需要下载地址，0为否，1为是，默认为0。获取到dlink后，参考下载文档进行下载操作
		// 下载文档: https://pan.baidu.com/union/doc/pkuo3snyp
		Dlink: 1,
	}
}

type BatchGetFileInfoReq struct {
	Dlink int
	FSIds []uint64
}

func (r *BatchGetFileInfoReq) AppendFSID(fsid uint64) *BatchGetFileInfoReq {
	r.FSIds = append(r.FSIds, fsid)
	return r
}

func (r *BatchGetFileInfoReq) GetFSIDString() string {
	bytesData, _ := json.Marshal(r.FSIds)
	return string(bytesData)
}

func (r *BatchGetFileInfoReq) GetDlink() string {
	return fmt.Sprintf("%d", r.Dlink)
}

type BatchGetFileInfoRes struct {
	List []*FileInfo `json:"list"`
}

func NewSearchFileReq(path string) *SearchFileReq {
	dir, name := filepath.Split(path)
	return &SearchFileReq{
		Dir: dir,
		Key: name,
	}
}

type SearchFileReq struct {
	Dir       string
	Key       string
	Recursion int
}

type SearchFileRes struct {
	HasMore int         `json:"has_more"`
	List    []*FileInfo `json:"list"`
}

type Opera string
type Ondup string
type Async int32

const (
	OndupFail      Ondup = "fail"
	OndupNewCopy         = "newcopy"
	OndupOverwrite       = "overwrite"
	OndupSkip            = "skip"

	AsyncSync Async = iota
	AsyncSelfAdaptation
	AsyncAsync

	OperaMove   Opera = "move"
	OperaCopy         = "copy"
	OperaDelete       = "delete"
	OperaRename       = "rename"
)

func NewFileManager(path, dest, newname string) *FileManager {
	return &FileManager{
		Path:    path,
		Dest:    dest,
		Newname: newname,
	}
}

type FileManager struct {
	Path    string `json:"path,omitempty"`
	Newname string `json:"newname,omitempty"`
	Dest    string `json:"dest,omitempty"`
	Ondup   string `json:"ondup,omitempty"`
}

func NewManageFileReq(
	opera Opera,
	filelist []*FileManager,
) *ManageFileReq {
	return &ManageFileReq{
		Async:    AsyncSelfAdaptation,
		Ondup:    OndupFail,
		Filelist: filelist,
		Opera:    opera,
	}
}

// https://pan.baidu.com/union/doc/mksg0s9l4?from=open-sdk-go
type ManageFileReq struct {
	Filelist []*FileManager
	// 0 同步，1 自适应，2 异步
	Async Async
	// 全局ondup,遇到重复文件的处理策略,
	// fail(默认，直接返回失败)、newcopy(重命名文件)、overwrite、skip
	Ondup Ondup
	// 文件操作参数，可实现文件复制、移动、重命名、删除，依次对应的参数值为：copy、move、rename、delete
	Opera Opera
}

func (r ManageFileReq) GetFilelistString() string {
	bytes, _ := json.Marshal(r.Filelist)
	return string(bytes)
}

type FileManagerInfo struct {
	Errno int    `json:"errno,omitempty"`
	Path  string `json:"path,omitempty"`
}

type ManageFileRes struct {
	Info      []FileManagerInfo `json:"info,omitempty"`
	RequestId int64             `json:"request_id,omitempty"`
	Taskid    int64             `json:"taskid,omitempty"`
}

type FileInfo struct {
	FSID           uint64            `json:"fs_id"`
	Path           string            `json:"path"`
	Size           int               `json:"size"`
	FileType       int               `json:"isdir"`
	Filename       string            `json:"filename"`
	ServerFilename string            `json:"server_filename"`
	Category       int               `json:"category"`
	Dlink          string            `json:"dlink"`
	MD5            string            `json:"md5"`
	Thumbs         map[string]string `json:"thumbs"`
	ServerCTime    int64             `json:"server_ctime"`
	ServerMTime    int64             `json:"server_mtime"`
	LocalCTime     int64             `json:"local_ctime"`
	LocalMTime     int64             `json:"local_mtime"`
}

func (f FileInfo) GetFilename() string {
	if f.ServerFilename != "" {
		return f.ServerFilename
	}
	return f.Filename
}

func (f FileInfo) GetSize() string {
	return tools.FormatSize(int64(f.Size))
}

func (f FileInfo) formatTime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

func (f FileInfo) GetServerCTime() string {
	return f.formatTime(f.ServerCTime)
}

func (f FileInfo) GetServerMTime() string {
	return f.formatTime(f.ServerMTime)
}

func (f FileInfo) GetLocalCTime() string {
	return f.formatTime(f.LocalCTime)
}

func (f FileInfo) GetLocalMTime() string {
	return f.formatTime(f.LocalMTime)
}

func (f FileInfo) IsDir() bool {
	if f.FileType == 1 {
		return true
	} else {
		return false
	}
}

func (f FileInfo) GetFileTypeIcon() string {
	if f.IsDir() {
		return ""
	}
	icon, ok := GetIconByPath(f.GetFilename())
	if !ok {
		icon = GetDefaultFileIcon()
	}
	return icon.Icon
}

func (f FileInfo) GetFileTypeEmoji() string {
	if f.IsDir() {
		// 🗂️
		return "\U0001f5c2"
	} else {
		switch f.Category {
		case 1:
			// 📹
			return "\U0001f4f9"
		case 2:
			// 🎵
			return "\U0001f3b5"
		case 3:
			// 🖼️
			return "\U0001f5bc"
		case 4:
			// 📄
			return "\U0001f4c4"
		case 5:
			// 🚀
			return "\U0001f680"
		case 6:
			var ext = filepath.Ext(f.Path)
			switch ext {
			case ".zip":
				return ""
			}
			// 其他 🤷
			return "\U0001f937"
		case 7:
			// 种子 🤷
			return "\U0001f937"
		}
		// 🤷
		return "\U0001f937"
	}
}

func (f FileInfo) GetFileType() string {
	if f.IsDir() {
		return "文件夹"
	} else {
		switch f.GetCategory() {
		case "其他":
			return "文件"
		default:
			return f.GetCategory()
		}
	}
}

func (f FileInfo) GetCategory() string {
	// 文件类型，1 视频、2 音频、3 图片、4 文档、5 应用、6 其他、7 种子
	switch f.Category {
	case 1:
		return "视频"
	case 2:
		return "音频"
	case 3:
		return "图片"
	case 4:
		return "文档"
	case 5:
		return "应用"
	case 6:
		return "其他"
	case 7:
		return "种子"
	}
	return "未知"
}
