package bdpan

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewGetFileListAllReqUsesStableOrder 固定 listall 排序，避免默认 time 排序在
// 大量同秒文件上翻页重复/漏项（bdpan-cli DBG-001）。
func TestNewGetFileListAllReqUsesStableOrder(t *testing.T) {
	req := NewGetFileListAllReq("/apps/demo")
	require.Equal(t, "name", req.Order)
	require.Equal(t, 0, req.Desc)
}

// TestGetFileListAllSendsOrder 验证请求确实把 order=name 发给了服务端。
func TestGetFileListAllSendsOrder(t *testing.T) {
	var query url.Values
	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errno":0,"has_more":0,"cursor":0,"list":[]}`))
	})

	_, err := GetFileListAll("token", NewGetFileListAllReq("/apps/demo"))
	require.NoError(t, err)
	require.Equal(t, "name", query.Get("order"))
}
