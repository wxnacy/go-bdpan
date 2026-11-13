package bdpan

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUploadAPIsContextCanceled 验证三个上传接口都能被 context 取消。
func TestUploadAPIsContextCanceled(t *testing.T) {
	setTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = io.WriteString(w, `{"errno":0}`)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	partPath := filepath.Join(t.TempDir(), "part.bin")
	require.NoError(t, os.WriteFile(partPath, []byte("part"), 0o644))
	partFile, err := os.Open(partPath)
	require.NoError(t, err)
	defer partFile.Close()

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "PreCreateFile",
			call: func() error {
				_, err := PreCreateFileWithContext(ctx, "token", NewPreCreateFileReq("/apps/a.bin", 4, []string{"abcd"}))
				return err
			},
		},
		{
			name: "UploadFilePart",
			call: func() error {
				req := NewUploadFilePartReq("/apps/a.bin", "uploadid", 0)
				req.File = partFile
				_, err := UploadFilePartWithContext(ctx, "token", req)
				return err
			},
		},
		{
			name: "CreateFile",
			call: func() error {
				_, err := CreateFileWithContext(ctx, "token", NewCreateFileReq("/apps/a.bin", 4, 0, "uploadid", []string{"abcd"}))
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			require.Error(t, err)
			require.True(t, errors.Is(err, context.Canceled), "应返回 context.Canceled，实际: %v", err)
		})
	}
}
