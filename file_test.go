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
			res, err := DeleteFiles(accessToken, f.Path)
			assert.NoError(t, err)
			info := res.Info[0]
			assert.Equal(t, info.Errno, 0)
			assert.Equal(t, info.Path, f.Path)
			return
		}
	}
}

func TestPreCreateFile(t *testing.T) {
	// 仅在有access token时执行测试
	if accessToken == "" {
		t.Skip("BDPAN_ACCESS_TOKEN is not set")
	}

	// 准备测试数据
	path := "/apps/bdpan/for_test/test_precreate.txt"
	size := int32(1024) // 1KB
	// 对于小于4MB的文件，block_list只包含一个MD5值
	blockList := []string{"e08b8e863d2fffce685530608305598c"}

	// 创建预上传请求
	req := NewPreCreateFileReq(path, size, blockList)
	// 设置文件命名策略为1（当path冲突时，进行重命名）
	req.SetRtype(1)

	// 调用预上传函数
	res, err := PreCreateFile(accessToken, req)
	assert.NoError(t, err)

	// 验证返回结果
	assert.Equal(t, int32(0), res.Errno, "预上传失败")
	assert.NotEmpty(t, res.Uploadid, "返回的uploadid不能为空")
	assert.Equal(t, path, res.Path, "返回的path与请求的path不一致")

	// 打印返回结果
	fmt.Printf("预上传成功: uploadid=%s, path=%s\n", res.Uploadid, res.Path)
}
