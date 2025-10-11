package bdpan

import (
	"encoding/json"
	"fmt"
	"os"
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

func NewGetFileListAllReq(path string) *GetFileListAllReq {
	return &GetFileListAllReq{
		Path:      path,
		Recursion: 1, // Recursive
		Limit:     1000,
		Start:     0,
	}
}

type GetFileListAllReq struct {
	Path      string
	Recursion int
	Web       int
	Start     uint64
	Limit     uint32
	Order     string
	Desc      int
}

type GetFileListAllRes struct {
	List      []*FileInfo `json:"list"`
	RequestID string      `json:"request_id"`
	HasMore   int         `json:"has_more"`
	Cursor    uint64      `json:"cursor"`
	ErrorRes
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

func (r *BatchGetFileInfoReq) AppendFSID(fsid ...uint64) *BatchGetFileInfoReq {
	r.FSIds = append(r.FSIds, fsid...)
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

type (
	Opera string
	Ondup string
	Async int32
)

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
	ErrorRes
	Path   string `json:"path,omitempty"`
	ToPath string `json:"to_path,omitempty"`
}

type ManageFileRes struct {
	ErrorRes
	Infos     []FileManagerInfo `json:"info,omitempty"`
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
			ext := filepath.Ext(f.Path)
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

// PreCreateFileReq 预上传请求参数
func NewPreCreateFileReq(path string, size int32, blockList []string) *PreCreateFileReq {
	return &PreCreateFileReq{
		Path:      path,
		Size:      size,
		Isdir:     0,
		Autoinit:  1,
		BlockList: blockList,
	}
}

// PreCreateFileReq 预上传请求参数
// https://pan.baidu.com/union/doc/3ksg0s9r7
// 预上传是通知网盘云端新建一个上传任务，网盘云端返回唯一ID uploadid 来标识此上传任务。
type PreCreateFileReq struct {
	// 上传后使用的文件绝对路径，需要urlencode
	Path string
	// 文件和目录两种情况：上传文件时，表示文件的大小，单位B；上传目录时，表示目录的大小，目录的话大小默认为0
	Size int32
	// 是否为目录，0 文件，1 目录
	Isdir int32
	// 固定值1
	Autoinit int32
	// 文件各分片MD5数组的json串
	// 如果上传的文件小于4MB，其md5值（32位小写）即为block_list字符串数组的唯一元素
	// 如果上传的文件大于4MB，需要将上传的文件按照4MB大小在本地切分成分片，不足4MB的分片自动成为最后一个分片，所有分片的md5值（32位小写）组成的字符串数组即为block_list
	BlockList []string
	// 文件命名策略
	// 1 表示当path冲突时，进行重命名
	// 2 表示当path冲突且block_list不同时，进行重命名
	// 3 当云端存在同名文件时，对该文件进行覆盖
	Rtype int32
	// 上传ID
	Uploadid string
	// 文件MD5，32位小写
	ContentMD5 string
	// 文件校验段的MD5，32位小写，校验段对应文件前256KB
	SliceMD5 string
	// 客户端创建时间， 默认为当前时间戳
	LocalCTime string
	// 客户端修改时间，默认为当前时间戳
	LocalMTime string
}

func (r *PreCreateFileReq) SetRtype(rtype int32) *PreCreateFileReq {
	r.Rtype = rtype
	return r
}

func (r *PreCreateFileReq) SetUploadid(uploadid string) *PreCreateFileReq {
	r.Uploadid = uploadid
	return r
}

func (r *PreCreateFileReq) SetContentMD5(contentMD5 string) *PreCreateFileReq {
	r.ContentMD5 = contentMD5
	return r
}

func (r *PreCreateFileReq) SetSliceMD5(sliceMD5 string) *PreCreateFileReq {
	r.SliceMD5 = sliceMD5
	return r
}

func (r *PreCreateFileReq) SetLocalCTime(localCTime string) *PreCreateFileReq {
	r.LocalCTime = localCTime
	return r
}

func (r *PreCreateFileReq) SetLocalMTime(localMTime string) *PreCreateFileReq {
	r.LocalMTime = localMTime
	return r
}

// GetBlockListString 返回block_list的json字符串
func (r *PreCreateFileReq) GetBlockListString() string {
	bytesData, _ := json.Marshal(r.BlockList)
	return string(bytesData)
}

// PreCreateFileRes 预上传响应参数
type PreCreateFileRes struct {
	ErrorRes
	// 文件的绝对路径
	Path string `json:"path,omitempty"`
	// 上传唯一ID标识此上传任务
	Uploadid string `json:"uploadid,omitempty"`
	// 返回类型，系统内部状态字段
	ReturnType int32 `json:"return_type,omitempty"`
	// 需要上传的分片序号列表，索引从0开始
	BlockList []int `json:"block_list,omitempty"`
	// 请求ID
	RequestId int64 `json:"request_id,omitempty"`
}

// UploadFilePartReq 分片上传请求参数
func NewUploadFilePartReq(path string, uploadid string, partseq int) *UploadFilePartReq {
	return &UploadFilePartReq{
		Path:     path,
		Uploadid: uploadid,
		Partseq:  partseq,
		Type:     "tmpfile",
	}
}

// UploadFilePartReq 分片上传请求参数
// https://pan.baidu.com/union/doc/nksg0s9vi
// 分片上传，这里是实际的文件内容传送部分。一般多为大于4MB的文件，需将文件以4MB为单位切分，对切分后得到的n个分片一一调用该接口进行传送。
type UploadFilePartReq struct {
	// 上传后使用的文件绝对路径，需要urlencode，需要与上一个阶段预上传precreate接口中的path保持一致
	Path string
	// 上一个阶段预上传precreate接口下发的uploadid
	Uploadid string
	// 文件分片的位置序号，从0开始，参考上一个阶段预上传precreate接口返回的block_list
	Partseq int
	// 固定值 tmpfile
	Type string
	// 要进行传送的本地文件分片
	File *os.File
}

// UploadFilePartRes 分片上传响应参数
type UploadFilePartRes struct {
	ErrorRes
	// 文件切片云端md5
	Md5 string `json:"md5,omitempty"`
}

// CreateFileReq 创建文件请求参数
func NewCreateFileReq(path string, size int32, isdir int32, uploadid string, blockList []string) *CreateFileReq {
	return &CreateFileReq{
		Path:      path,
		Size:      size,
		Isdir:     isdir,
		Uploadid:  uploadid,
		BlockList: blockList,
	}
}

// CreateFileReq 创建文件请求参数
// https://pan.baidu.com/union/doc/rksg0sa17
// 将多个文件分片合并成一个文件，生成文件基本信息，完成文件的上传最后一步。
type CreateFileReq struct {
	// 上传后使用的文件绝对路径，需要urlencode，需要与预上传precreate接口中的path保持一致
	Path string
	// 文件或目录的大小，必须要和文件真实大小保持一致，需要与预上传precreate接口中的size保持一致
	Size int32
	// 是否目录，0 文件、1 目录，需要与预上传precreate接口中的isdir保持一致
	Isdir int32
	// 预上传precreate接口下发的uploadid
	Uploadid string
	// 文件各分片md5数组的json串需要与预上传precreate接口中的block_list保持一致，同时对应分片上传superfile2接口返回的md5，且要按照序号顺序排列
	BlockList []string
	// 文件命名策略，默认0
	// 0 为不重命名，返回冲突
	// 1 为只要path冲突即重命名
	// 2 为path冲突且block_list不同才重命名
	// 3 为覆盖，需要与预上传precreate接口中的rtype保持一致
	Rtype int32
	// 客户端创建时间(精确到秒)，默认为当前时间戳
	LocalCTime int64
	// 客户端修改时间(精确到秒)，默认为当前时间戳
	LocalMTime int64
}

func (r *CreateFileReq) SetRtype(rtype int32) *CreateFileReq {
	r.Rtype = rtype
	return r
}

func (r *CreateFileReq) SetLocalCTime(localCTime int64) *CreateFileReq {
	r.LocalCTime = localCTime
	return r
}

func (r *CreateFileReq) SetLocalMTime(localMTime int64) *CreateFileReq {
	r.LocalMTime = localMTime
	return r
}

// GetBlockListString 返回block_list的json字符串
func (r *CreateFileReq) GetBlockListString() string {
	bytesData, _ := json.Marshal(r.BlockList)
	return string(bytesData)
}

// CreateFileRes 创建文件响应参数
type CreateFileRes struct {
	ErrorRes
	// 文件在云端的唯一标识ID
	FSId uint64 `json:"fs_id,omitempty"`
	// 文件的MD5，只有提交文件时才返回，提交目录时没有该值
	Md5 string `json:"md5,omitempty"`
	// 文件名
	ServerFilename string `json:"server_filename,omitempty"`
	// 分类类型, 1 视频 2 音频 3 图片 4 文档 5 应用 6 其他 7 种子
	Category int32 `json:"category,omitempty"`
	// 上传后使用的文件绝对路径
	Path string `json:"path,omitempty"`
	// 文件大小，单位B
	Size uint64 `json:"size,omitempty"`
	// 文件创建时间
	Ctime uint64 `json:"ctime,omitempty"`
	// 文件修改时间
	Mtime uint64 `json:"mtime,omitempty"`
	// 是否目录，0 文件、1 目录
	Isdir int32 `json:"isdir,omitempty"`
	// 请求ID
	RequestId int64 `json:"request_id,omitempty"`
}
