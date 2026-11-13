package bdpan

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	sdk "github.com/wxnacy/go-bdpan/openapi"
)

// setTestServer 把全局 client 指向测试服务器，测试结束后恢复。
func setTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: server.URL}}
	// 预上传/创建/分片接口在 OperationServers 里写死了线上地址，测试需要一并替换
	for name := range cfg.OperationServers {
		cfg.OperationServers[name] = sdk.ServerConfigurations{{URL: server.URL}}
	}
	clientMu.Lock()
	previous := client
	client = sdk.NewAPIClient(cfg)
	clientMu.Unlock()
	t.Cleanup(func() {
		clientMu.Lock()
		client = previous
		clientMu.Unlock()
		server.Close()
	})
	return server
}

// TestPreCreateFileIncludesResumeFields 验证预上传请求会带上断点续传相关字段。
func TestPreCreateFileIncludesResumeFields(t *testing.T) {
	var form url.Values
	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		form = r.PostForm
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = io.WriteString(w, `{"errno":0,"uploadid":"resume-id","block_list":[1]}`)
	})

	req := NewPreCreateFileReq("/apps/demo/a.bin", 4*1024*1024, []string{"abcd"})
	req.SetUploadid("prev-upload-id").
		SetSliceMD5("slice-md5-value").
		SetLocalCTime("1700000000").
		SetLocalMTime("1700000001")

	res, err := PreCreateFile("token", req)
	require.NoError(t, err)
	require.Equal(t, "resume-id", res.Uploadid)

	require.Equal(t, "prev-upload-id", form.Get("uploadid"))
	require.Equal(t, "slice-md5-value", form.Get("slice-md5"))
	require.Equal(t, "1700000000", form.Get("local_ctime"))
	require.Equal(t, "1700000001", form.Get("local_mtime"))
	require.Equal(t, `["abcd"]`, form.Get("block_list"))
}

// TestPreCreateFileOmitsEmptyResumeFields 验证空值不会以空字段形式发送。
func TestPreCreateFileOmitsEmptyResumeFields(t *testing.T) {
	var form url.Values
	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		form = r.PostForm
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = io.WriteString(w, `{"errno":0,"uploadid":"new-id","block_list":[1]}`)
	})

	req := NewPreCreateFileReq("/apps/demo/a.bin", 4*1024*1024, []string{"abcd"})
	_, err := PreCreateFile("token", req)
	require.NoError(t, err)

	for _, key := range []string{"uploadid", "slice-md5", "local_ctime", "local_mtime"} {
		require.False(t, form.Has(key), "空字段 %s 不应发送", key)
	}
}
