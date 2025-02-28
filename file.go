// https://pan.baidu.com/union/doc/9l4gmjcmj
package bdpan

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// 获取文件列表
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
func DeleteFiles(accessToken string, paths []string) (*ManageFileRes, error) {

	filelist := make([]*FileManager, 0)
	for _, p := range paths {
		filelist = append(filelist, NewFileManager(p, "", ""))
	}

	req := NewManageFileReq(OperaDelete, filelist)
	return ManageFile(accessToken, req)
}

func DeleteFile(accessToken string, path string) (*ManageFileRes, error) {
	paths := []string{path}
	return DeleteFiles(accessToken, paths)
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
	}
	return ToInterface[ManageFileRes](r)
}
