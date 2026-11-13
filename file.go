// https://pan.baidu.com/union/doc/9l4gmjcmj
package bdpan

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// GetFileListAll 递归获取文件列表
// https://pan.baidu.com/union/doc/Zksg0sb73
func GetFileListAll(accessToken string, req *GetFileListAllReq) (*GetFileListAllRes, error) {
	_, r, err := GetClient().
		MultimediafileApi.Xpanfilelistall(context.Background()).
		AccessToken(accessToken).
		Path(req.Path).
		Recursion(int32(req.Recursion)).
		Web(strconv.Itoa(req.Web)).
		Start(int32(req.Start)).
		Limit(int32(req.Limit)).
		Order(req.Order).
		Desc(int32(req.Desc)).
		Execute()

	return ToInterface[GetFileListAllRes](r, err)
}

// 获取文件列表
// https://pan.baidu.com/union/doc/nksg0sat9
func GetFileList(accessToken string, req *GetFileListReq) (*GetFileListRes, error) {
	_, r, _ := GetClient().
		FileinfoApi.Xpanfilelist(context.Background()).
		AccessToken(accessToken).
		Dir(req.Dir).
		Web(fmt.Sprintf("%d", req.Web)).
		Start(fmt.Sprintf("%d", req.Start)).
		Order(req.Order).
		Desc(req.Desc).
		Limit(req.Limit).
		Execute()
	return ToInterface[GetFileListRes](r)
}

// 获取文件详情
// https://pan.baidu.com/union/doc/Fksg0sbcm
func GetFileInfo(accessToken string, req *GetFileInfoReq) (*GetFileInfoRes, error) {
	batchReq := NewBatchGetFileListReq(req.FSID)
	batchReq.Dlink = req.Dlink
	res, err := BatchGetFileInfo(accessToken, batchReq)
	if err != nil {
		return nil, err
	}
	if len(res.List) == 0 {
		return nil, fmt.Errorf("%d not found", req.FSID)
	}
	return &GetFileInfoRes{
		FileInfo: *res.List[0],
	}, nil
}

// 批量获取文件详情
func BatchGetFileInfo(accessToken string, req *BatchGetFileInfoReq) (*BatchGetFileInfoRes, error) {
	api := GetClient().
		MultimediafileApi.Xpanmultimediafilemetas(context.Background()).
		AccessToken(accessToken).
		Dlink(req.GetDlink()).
		Fsids(req.GetFSIDString())
	if req.Blocks > 0 {
		api = api.Blocks(int32(req.Blocks))
	}
	_, r, err := api.Execute()
	if err != nil {
		return nil, ToAPIError(r, err)
	}
	return ToInterface[BatchGetFileInfoRes](r)
}

// 搜索文件
func SearchFile(accessToken string, req *SearchFileReq) (*SearchFileRes, error) {
	_, r, _ := GetClient().
		FileinfoApi.Xpanfilesearch(context.Background()).
		AccessToken(accessToken).
		Key(req.Key).
		Recursion(strconv.Itoa(req.Recursion)).
		Dir(req.Dir).
		Execute()
	return ToInterface[SearchFileRes](r)
}

// 删除文件列表
func DeleteFiles(accessToken string, paths ...string) (*ManageFileRes, error) {
	filelist := make([]*FileManager, 0)
	for _, p := range paths {
		filelist = append(filelist, NewFileManager(p, "", ""))
	}

	req := NewManageFileReq(OperaDelete, filelist)
	return ManageFile(accessToken, req)
}

// 移动文件列表
func MoveFiles(accessToken string, dir string, paths ...string) (*ManageFileRes, error) {
	filelist := make([]*FileManager, 0)
	for _, p := range paths {
		name := filepath.Base(p)
		filelist = append(filelist, NewFileManager(p, dir, name))
	}
	req := NewManageFileReq(OperaMove, filelist)
	return ManageFile(accessToken, req)
}

// 重命名
// RenameFiles(token, NewFileManager("/apps/bdpan/test.go", "", "new.go"))
func RenameFiles(accessToken string, files ...*FileManager) (*ManageFileRes, error) {
	req := NewManageFileReq(OperaRename, files)
	return ManageFile(accessToken, req)
}

// 管理文件
// https://pan.baidu.com/union/doc/mksg0s9l4
func ManageFile(accessToken string, req *ManageFileReq) (*ManageFileRes, error) {
	var r *http.Response
	switch req.Opera {
	case OperaDelete:
		r, _ = GetClient().
			FilemanagerApi.
			Filemanagerdelete(context.Background()).
			AccessToken(accessToken).
			Async(int32(req.Async)).
			Ondup(string(req.Ondup)).
			Filelist(req.GetFilelistString()).
			Execute()
	case OperaMove:
		r, _ = GetClient().
			FilemanagerApi.
			Filemanagermove(context.Background()).
			AccessToken(accessToken).
			Async(int32(req.Async)).
			Ondup(string(req.Ondup)).
			Filelist(req.GetFilelistString()).
			Execute()
	case OperaRename:
		r, _ = GetClient().
			FilemanagerApi.
			Filemanagerrename(context.Background()).
			AccessToken(accessToken).
			Async(int32(req.Async)).
			Ondup(string(req.Ondup)).
			Filelist(req.GetFilelistString()).
			Execute()
	}
	return ToInterface[ManageFileRes](r)
}

// 预上传
// https://pan.baidu.com/union/doc/3ksg0s9r7
// 预上传是通知网盘云端新建一个上传任务，网盘云端返回唯一ID uploadid 来标识此上传任务。
func PreCreateFile(accessToken string, req *PreCreateFileReq) (*PreCreateFileRes, error) {
	return PreCreateFileWithContext(context.Background(), accessToken, req)
}

// PreCreateFileWithContext 是 PreCreateFile 的 context 版本，ctx 取消时请求会被中断。
func PreCreateFileWithContext(ctx context.Context, accessToken string, req *PreCreateFileReq) (*PreCreateFileRes, error) {
	_, r, err := GetClient().
		FileuploadApi.Xpanfileprecreate(ctx).
		AccessToken(accessToken).
		Path(req.Path).
		Isdir(req.Isdir).
		Size(req.Size).
		Autoinit(req.Autoinit).
		Rtype(req.Rtype).
		BlockList(req.GetBlockListString()).
		ContentMD5(req.ContentMD5).
		Uploadid(req.Uploadid).
		SliceMD5(req.SliceMD5).
		LocalCTime(req.LocalCTime).
		LocalMTime(req.LocalMTime).
		Execute()
	if err != nil {
		return nil, ToAPIError(r, err)
	}
	return ToInterface[PreCreateFileRes](r)
}

// 分片上传
// https://pan.baidu.com/union/doc/nksg0s9vi
// 分片上传，这里是实际的文件内容传送部分。一般多为大于4MB的文件，需将文件以4MB为单位切分，对切分后得到的n个分片一一调用该接口进行传送。
func UploadFilePart(accessToken string, req *UploadFilePartReq) (*UploadFilePartRes, error) {
	return UploadFilePartWithContext(context.Background(), accessToken, req)
}

// UploadFilePartWithContext 是 UploadFilePart 的 context 版本，ctx 取消时请求会被中断。
func UploadFilePartWithContext(ctx context.Context, accessToken string, req *UploadFilePartReq) (*UploadFilePartRes, error) {
	_, r, err := GetClient().
		FileuploadApi.Pcssuperfile2(ctx).
		AccessToken(accessToken).
		Path(req.Path).
		Uploadid(req.Uploadid).
		Partseq(fmt.Sprintf("%d", req.Partseq)).
		Type_(req.Type).
		File(req.File).
		Execute()
	if err != nil {
		return nil, ToAPIError(r, err)
	}
	return ToInterface[UploadFilePartRes](r)
}

// UploadFilePartWithReader 上传分片，内容从 reader 流式读取，不整片读入内存。
//
// 与 UploadFilePart 的区别是调用方可以直接传 io.SectionReader 等按需读取的
// 视图，不需要先落临时文件；reader 的生命周期由调用方管理。
func UploadFilePartWithReader(
	ctx context.Context,
	accessToken string,
	req *UploadFilePartReq,
	reader io.Reader,
) (*UploadFilePartRes, error) {
	part := GetClient().
		FileuploadApi.Pcssuperfile2(ctx).
		AccessToken(accessToken).
		Path(req.Path).
		Uploadid(req.Uploadid).
		Partseq(fmt.Sprintf("%d", req.Partseq)).
		Type_(req.Type).
		FileReader(reader)
	if size, ok := readerLength(reader); ok {
		part = part.FileSize(size)
	}
	_, r, err := part.Execute()
	if err != nil {
		return nil, ToAPIError(r, err)
	}
	return ToInterface[UploadFilePartRes](r)
}

// readerLength 推断 reader 的总字节数，未知时返回 false。
func readerLength(reader io.Reader) (int64, bool) {
	if sizer, ok := reader.(interface{ Size() int64 }); ok {
		return sizer.Size(), true
	}
	if file, ok := reader.(*os.File); ok {
		if st, err := file.Stat(); err == nil {
			return st.Size(), true
		}
	}
	return 0, false
}

// 创建文件
// https://pan.baidu.com/union/doc/rksg0sa17
// 将多个文件分片合并成一个文件，生成文件基本信息，完成文件的上传最后一步。
func CreateFile(accessToken string, req *CreateFileReq) (*CreateFileRes, error) {
	return CreateFileWithContext(context.Background(), accessToken, req)
}

// CreateFileWithContext 是 CreateFile 的 context 版本，ctx 取消时请求会被中断。
func CreateFileWithContext(ctx context.Context, accessToken string, req *CreateFileReq) (*CreateFileRes, error) {
	_, r, err := GetClient().
		FileuploadApi.Xpanfilecreate(ctx).
		AccessToken(accessToken).
		Path(req.Path).
		Isdir(req.Isdir).
		Size(req.Size).
		Rtype(req.Rtype).
		Uploadid(req.Uploadid).
		BlockList(req.GetBlockListString()).
		Execute()
	if err != nil {
		return nil, ToAPIError(r, err)
	}
	return ToInterface[CreateFileRes](r)
}
