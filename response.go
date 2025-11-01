package bdpan

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorRes struct {
	// 错误码
	Errno     int32  `json:"errno,omitempty"`
	Errmsg    string `json:"errmsg,omitempty"`
	ErrorMsg  string `json:"error,omitempty"`
	ErrorDesc string `json:"error_description,omitempty"`
}

func (e ErrorRes) IsError() bool {
	return e.Errno != 0 || e.ErrorMsg != ""
}

func (e ErrorRes) Error() string {
	if e.ErrorDesc != "" {
		return e.ErrorDesc
	}
	return GetApiErrByErrno(e.Errno).Error()
}

func ToInterface[T any](r *http.Response, args ...any) (*T, error) {
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
		var apiErr ErrorRes
		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			return nil, err
		}
		return nil, &apiErr

	}
	return &i, nil
}
