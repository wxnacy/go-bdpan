package bdpan

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	accessToken string
	testDir     string = "/apps/bdpan/for_test"
)

func init() {
	accessToken = os.Getenv("BDPAN_ACCESS_TOKEN")
}

func TestGetFileList(t *testing.T) {
	req := NewGetFileListReq()
	res, err := GetFileList(accessToken, req)
	assert.NoError(t, err)
	fileList := res.List
	assert.GreaterOrEqual(t, len(fileList), 1)
}

func TestGetFileInfo(t *testing.T) {
	listRes, err := GetFileList(accessToken, NewGetFileListReq())
	reqInfo := listRes.List[0]
	fsid := reqInfo.FSID

	req := NewGetFileInfoReq(fsid)
	res, err := GetFileInfo(accessToken, req)
	assert.NoError(t, err)
	assert.Equal(t, res.Path, reqInfo.Path)

	fsid = 241766811636032
	req = NewGetFileInfoReq(fsid)
	res, err = GetFileInfo(accessToken, req)
	assert.Errorf(t, err, "%d not found", fsid)
}

func TestBatchGetFileInfo(t *testing.T) {
	listRes, err := GetFileList(accessToken, NewGetFileListReq())
	reqInfo := listRes.List[0]
	fsid := reqInfo.FSID
	req := NewBatchGetFileListReq(fsid)
	for _, f := range listRes.List {
		req.AppendFSID(f.FSID)
	}
	res, err := BatchGetFileInfo(accessToken, req)
	assert.NoError(t, err)
	assert.Equal(t, len(res.List), len(listRes.List))
}

// func TestSearchFile(t *testing.T) {
// req := NewSearchFileReq("/1视频/1电影")
// req.Recursion = 0
// res, err := SearchFile(accessToken, req)
// assert.NoError(t, err)

// files := res.List
// assert.Equal(t, len(files), 1)
// }

func TestDeleteFile(t *testing.T) {
	listRes, err := GetFileList(accessToken, NewGetFileListReq().SetDir(testDir).SetLimit(100))
	assert.NoError(t, err)
	for _, f := range listRes.List {
		fmt.Println(f.Path)
		if !f.IsDir() {
			fmt.Printf("Delete path: %s\n", f.Path)
			res, err := DeleteFile(accessToken, f.Path)
			assert.NoError(t, err)
			info := res.Info[0]
			assert.Equal(t, info.Errno, 0)
			assert.Equal(t, info.Path, f.Path)
			return
		}
	}
}
