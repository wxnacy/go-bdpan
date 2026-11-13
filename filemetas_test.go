package bdpan

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBatchGetFileInfoBlocks 验证 blocks=1 会作为查询参数发送，并解析 block_list。
func TestBatchGetFileInfoBlocks(t *testing.T) {
	var query url.Values
	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = io.WriteString(w, `{"errno":0,"list":[{"fs_id":1,"md5":"abc","block_list":["b0","b1"]}]}`)
	})

	req := NewBatchGetFileListReq(1)
	req.Blocks = 1
	res, err := BatchGetFileInfo("token", req)
	require.NoError(t, err)
	require.Equal(t, "1", query.Get("blocks"))
	require.Len(t, res.List, 1)
	require.Equal(t, []string{"b0", "b1"}, res.List[0].BlockList)
}

// TestBatchGetFileInfoNoBlocks 验证未显式开启时不发送 blocks 参数，保持旧行为。
func TestBatchGetFileInfoNoBlocks(t *testing.T) {
	var query url.Values
	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = io.WriteString(w, `{"errno":0,"list":[{"fs_id":1,"md5":"abc"}]}`)
	})

	_, err := BatchGetFileInfo("token", NewBatchGetFileListReq(1))
	require.NoError(t, err)
	require.False(t, query.Has("blocks"))
}
