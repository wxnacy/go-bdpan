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
	// ShowMsg 是服务端返回给用户看的提示信息
	ShowMsg string `json:"show_msg,omitempty"`
	// StatusCode 是 HTTP 状态码，仅本地填充，不参与 JSON 编解码
	StatusCode int `json:"-"`
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

// ToAPIError 把 SDK 返回的非 2xx 响应转换为带 HTTP 状态码与 show_msg 的结构化错误。
//
// SDK 的 Execute 在非 2xx 时返回 GenericOpenAPIError，会丢掉状态码；这里从响应体和
// 状态行补齐，调用方可用 errors.As 取出 *ErrorRes 判断 errno 与是否可重试。
func ToAPIError(resp *http.Response, err error) error {
	if err == nil {
		return nil
	}
	if resp == nil {
		return err
	}
	if resp.StatusCode < 300 {
		return err
	}

	apiErr := &ErrorRes{StatusCode: resp.StatusCode}
	if errWithBody, ok := err.(interface{ Body() []byte }); ok {
		if body := errWithBody.Body(); len(body) > 0 {
			_ = json.Unmarshal(body, apiErr)
		}
	}
	if apiErr.Errmsg == "" {
		apiErr.Errmsg = resp.Status
	}
	return apiErr
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
		apiErr.StatusCode = r.StatusCode
		return nil, &apiErr

	}
	return &i, nil
}
