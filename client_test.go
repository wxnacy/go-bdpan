package bdpan

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestSetHTTPTimeout 验证设置超时后慢请求会按预期超时。
func TestSetHTTPTimeout(t *testing.T) {
	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = io.WriteString(w, `{"errno":0}`)
	})
	t.Cleanup(func() { SetHTTPTimeout(0) })

	SetHTTPTimeout(20 * time.Millisecond)

	start := time.Now()
	_, err := PreCreateFile("token", NewPreCreateFileReq("/apps/demo/a.bin", 4, []string{"abcd"}))
	require.Error(t, err)
	require.Less(t, time.Since(start), 150*time.Millisecond, "请求应被超时中断")

	var timeoutErr interface{ Timeout() bool }
	require.True(t, errors.As(err, &timeoutErr) && timeoutErr.Timeout(), "应返回超时错误，实际: %v", err)
}
