package bdpan

import (
	"encoding/json"
	"fmt"
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

// Error 返回给用户看的错误文案。
//
// 服务端的说明字段优先级高于 errno 映射：errno=0 的异常响应（例如网关错误、
// 只有 error 字段的分片失败）不能再显示成「请求成功」，否则会误导排查。
func (e ErrorRes) Error() string {
	switch {
	case e.ErrorDesc != "":
		return e.ErrorDesc
	case e.ShowMsg != "":
		return e.ShowMsg
	case e.ErrorMsg != "":
		return e.ErrorMsg
	case e.Errmsg != "":
		return e.Errmsg
	case e.StatusCode > 0:
		return fmt.Sprintf("HTTP %d %s", e.StatusCode, http.StatusText(e.StatusCode))
	case e.Errno != 0:
		return GetApiErrByErrno(e.Errno).Error()
	default:
		return "未知错误"
	}
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
	var body []byte
	if errWithBody, ok := err.(interface{ Body() []byte }); ok {
		body = errWithBody.Body()
		if len(body) > 0 {
			_ = json.Unmarshal(body, apiErr)
		}
	}
	if apiErr.Errmsg == "" {
		apiErr.Errmsg = resp.Status
	}

	// 错误分支此前不落日志，errno=0 的异常响应很难定位；把状态、响应体
	// （含 request_id）记录下来，日志文件固定 DEBUG 级别，便于事后排查。
	path := ""
	if resp.Request != nil && resp.Request.URL != nil {
		path = resp.Request.URL.Path
	}
	GetLogger().Debugf("%s error status=%d body=%s\n", path, resp.StatusCode, truncateLogBody(body))
	return apiErr
}

// truncateLogBody 限制错误响应体的日志长度，避免网关返回整页 HTML 时刷屏。
func truncateLogBody(body []byte) string {
	const maxLen = 4096
	if len(body) <= maxLen {
		return string(body)
	}
	return string(body[:maxLen]) + "...(truncated)"
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
