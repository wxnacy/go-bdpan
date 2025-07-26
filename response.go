package bdpan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiError struct {
	StatusCode       int    `json:"status_code"`
	Errno            int32  `json:"errno,omitempty"`
	ErrMsg           string `json:"errmsg"`
	Erro             string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
}

func (e *ApiError) Error() string {
	// return fmt.Sprintf("%d: %s(%s)", e.StatusCode, e.Erro, e.ErrorDescription)
	return e.String()
}

func (e *ApiError) String() string {
	if e.ErrorCode > 0 {
		return fmt.Sprintf("%d[%s]", e.ErrorCode, e.ErrorMsg)
	} else if e.Erro != "" {
		return fmt.Sprintf("%s[%s]", e.Erro, e.ErrorDescription)
	} else {
		switch e.Errno {
		case -9:
			return ErrPathNotFound.Error()
		case -6:
			return ErrAccessFail.Error()
		case 2:
			return ErrParamError.Error()
		case 6:
			return ErrUserNoUse.Error()
		case 12:
			return ErrPathExists.Error()
		case 111:
			return ErrAccessTokenFail.Error()
		case 31034, 9013, 9019:
			return ErrApiFrequent.Error()
		default:
			return fmt.Sprintf("未知错误: %d", e.Errno)
		}
	}
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
