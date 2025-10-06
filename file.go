// https://pan.baidu.com/union/doc/9l4gmjcmj
package bdpan

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
)

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
	_, r, _ := GetClient().
		MultimediafileApi.Xpanmultimediafilemetas(
		context.Background()).AccessToken(accessToken).
		Dlink(req.GetDlink()).
		Fsids(req.GetFSIDString()).
		Execute()
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
	}
	return ToInterface[ManageFileRes](r)
}

// 预上传
// https://pan.baidu.com/union/doc/3ksg0s9r7
// 预上传是通知网盘云端新建一个上传任务，网盘云端返回唯一ID uploadid 来标识此上传任务。
func PreCreateFile(accessToken string, req *PreCreateFileReq) (*PreCreateFileRes, error) {
	_, r, _ := GetClient().
		FileuploadApi.Xpanfileprecreate(context.Background()).
		AccessToken(accessToken).
		Path(req.Path).
		Isdir(req.Isdir).
		Size(req.Size).
		Autoinit(req.Autoinit).
		Rtype(req.Rtype).
		BlockList(req.GetBlockListString()).
		ContentMD5(req.ContentMD5).
		Execute()
	return ToInterface[PreCreateFileRes](r)
}

// 分片上传
// https://pan.baidu.com/union/doc/nksg0s9vi
// 分片上传，这里是实际的文件内容传送部分。一般多为大于4MB的文件，需将文件以4MB为单位切分，对切分后得到的n个分片一一调用该接口进行传送。
func UploadFilePart(accessToken string, req *UploadFilePartReq) (*UploadFilePartRes, error) {
	_, r, _ := GetClient().
		FileuploadApi.Pcssuperfile2(context.Background()).
		AccessToken(accessToken).
		Path(req.Path).
		Uploadid(req.Uploadid).
		Partseq(fmt.Sprintf("%d", req.Partseq)).
		Type_(req.Type).
		File(req.File).
		Execute()
	return ToInterface[UploadFilePartRes](r)
}

// 创建文件
// https://pan.baidu.com/union/doc/rksg0sa17
// 将多个文件分片合并成一个文件，生成文件基本信息，完成文件的上传最后一步。
func CreateFile(accessToken string, req *CreateFileReq) (*CreateFileRes, error) {
	_, r, _ := GetClient().
		FileuploadApi.Xpanfilecreate(context.Background()).
		AccessToken(accessToken).
		Path(req.Path).
		Isdir(req.Isdir).
		Size(req.Size).
		Rtype(req.Rtype).
		Uploadid(req.Uploadid).
		BlockList(req.GetBlockListString()).
		Execute()
	return ToInterface[CreateFileRes](r)
}
