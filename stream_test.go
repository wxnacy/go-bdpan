package bdpan

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUploadFilePartWithReaderStreams 验证分片内容从 Reader 流式发送，服务端收到的字节与 md5 一致。
func TestUploadFilePartWithReaderStreams(t *testing.T) {
	content := bytes.Repeat([]byte("stream-chunk-"), 8192)
	var received []byte
	var receivedMD5 string

	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// 百度网盘分片接口不接受 chunked，必须带 Content-Length
		require.Greater(t, r.ContentLength, int64(0), "流式上传应带 Content-Length")
		require.Empty(t, r.TransferEncoding, "不应使用 chunked")
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		received, _ = io.ReadAll(file)
		sum := md5.Sum(received)
		receivedMD5 = hex.EncodeToString(sum[:])

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = fmt.Fprintf(w, `{"md5":%q}`, receivedMD5)
	})

	req := NewUploadFilePartReq("/apps/demo/a.bin", "uploadid", 3)
	res, err := UploadFilePartWithReader(context.Background(), "token", req, bytes.NewReader(content))
	require.NoError(t, err)
	require.Equal(t, content, received)
	require.Equal(t, receivedMD5, res.Md5)
}

// TestUploadFilePartWithFileStillWorks 验证旧的 *os.File 入口保持兼容。
func TestUploadFilePartWithFileStillWorks(t *testing.T) {
	content := []byte("legacy-file-part")
	partPath := filepath.Join(t.TempDir(), "part.bin")
	require.NoError(t, os.WriteFile(partPath, content, 0o644))
	file, err := os.Open(partPath)
	require.NoError(t, err)
	defer file.Close()

	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		got, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer got.Close()
		received, _ := io.ReadAll(got)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		sum := md5.Sum(received)
		_, _ = fmt.Fprintf(w, `{"md5":%q}`, hex.EncodeToString(sum[:]))
	})

	req := NewUploadFilePartReq("/apps/demo/a.bin", "uploadid", 0)
	req.File = file
	res, err := UploadFilePart("token", req)
	require.NoError(t, err)
	sum := md5.Sum(content)
	require.Equal(t, hex.EncodeToString(sum[:]), res.Md5)
}
