package bdpan

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorRes struct {
	// 错误码
	Errno int32 `json:"errno,omitempty"`
}

func (e ErrorRes) IsError() bool {
	return e.Errno != 0
}

func (e ErrorRes) Error() string {
	return GetApiErrByErrno(e.Errno).Error()
}

type ApiError struct {
	ErrorRes
	StatusCode       int    `json:"status_code"`
	ErrMsg           string `json:"errmsg"`
	Erro             string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
}

func ToInterface[T any](r *http.Response) (*T, error) {
	var i T
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	GetLogger().Debugf("%s code %d\n", r.Request.URL.Path, r.StatusCode)
	GetLogger().Debugf("%s response %s\n", r.Request.URL.Path, string(bodyBytes))
	if r.StatusCode == 200 {
		if err := json.Unmarshal(bodyBytes, &i); err != nil {
			return nil, err
		}
	} else {
		var apiErr ApiError
		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			return nil, err
		}
		return nil, &apiErr

	}
	return &i, nil
}
