package bdpan

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsRetryable 覆盖网络错误、HTTP 状态码与业务 errno 的可重试判定。
func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "context canceled", err: context.Canceled, want: false},
		{name: "context deadline", err: context.DeadlineExceeded, want: true},
		{name: "5xx", err: &ErrorRes{StatusCode: 503}, want: true},
		{name: "429", err: &ErrorRes{StatusCode: 429}, want: true},
		{name: "400 参数错误", err: &ErrorRes{StatusCode: 400, Errno: 31023}, want: false},
		{name: "限流 errno", err: &ErrorRes{Errno: 9013}, want: true},
		{name: "鉴权 errno", err: &ErrorRes{Errno: 31045}, want: false},
		{name: "配额超限 errno", err: &ErrorRes{Errno: 31365}, want: false},
		{name: "旧 ApiErr 限流", err: &ApiErr{Errno: 4}, want: true},
		{name: "旧 ApiErr 参数", err: &ApiErr{Errno: 2}, want: false},
		{name: "包装后的错误", err: fmt.Errorf("upload: %w", &ErrorRes{StatusCode: 500}), want: true},
		{name: "errno=0 但带 error 字段", err: &ErrorRes{ErrorMsg: "临时失败"}, want: true},
		{name: "errno=0 的 502", err: &ErrorRes{StatusCode: 502}, want: true},
		{name: "errno=0 的 404", err: &ErrorRes{StatusCode: 404, ErrorMsg: "不存在"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsRetryable(tt.err))
		})
	}
}

// TestErrorResMessage 验证错误文案优先使用服务端说明，不再把 errno=0 显示成「请求成功」。
func TestErrorResMessage(t *testing.T) {
	tests := []struct {
		name string
		err  ErrorRes
		want string
	}{
		{name: "error_description 优先", err: ErrorRes{Errno: 0, ErrorDesc: "描述错误", ShowMsg: "提示"}, want: "描述错误"},
		{name: "show_msg", err: ErrorRes{Errno: 0, ShowMsg: "请稍后重试"}, want: "请稍后重试"},
		{name: "error 字段", err: ErrorRes{Errno: 0, ErrorMsg: "分片上传失败"}, want: "分片上传失败"},
		{name: "errmsg", err: ErrorRes{Errno: 31023, Errmsg: "参数错误"}, want: "参数错误"},
		{name: "HTTP 状态兜底", err: ErrorRes{Errno: 0, StatusCode: http.StatusBadGateway}, want: "HTTP 502 Bad Gateway"},
		{name: "errno 映射", err: ErrorRes{Errno: 31023}, want: "参数错误"},
		{name: "errno=0 且无说明", err: ErrorRes{}, want: "未知错误"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.err.Error())
		})
	}
}

// TestToAPIError 验证非 2xx 响应会被转换为带状态码与 show_msg 的 *ErrorRes。
func TestToAPIError(t *testing.T) {
	t.Run("503 可重试", func(t *testing.T) {
		setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, `{"errno":-1,"errmsg":"服务不可用","show_msg":"请稍后重试"}`)
		})

		_, err := PreCreateFile("token", NewPreCreateFileReq("/apps/demo/a.bin", 4, []string{"abcd"}))
		require.Error(t, err)

		var apiErr *ErrorRes
		require.True(t, errors.As(err, &apiErr), "应返回 *ErrorRes，实际: %T %v", err, err)
		require.Equal(t, int32(-1), apiErr.Errno)
		require.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
		require.Equal(t, "请稍后重试", apiErr.ShowMsg)
		require.True(t, IsRetryable(err))
	})

	t.Run("400 不可重试", func(t *testing.T) {
		setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"errno":31023,"errmsg":"参数错误"}`)
		})

		_, err := CreateFile("token", NewCreateFileReq("/apps/demo/a.bin", 4, 0, "uploadid", []string{"abcd"}))
		require.Error(t, err)

		var apiErr *ErrorRes
		require.True(t, errors.As(err, &apiErr))
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		require.Equal(t, int32(31023), apiErr.Errno)
		require.False(t, IsRetryable(err))
	})
}
