package bdpan

import (
	"fmt"
)

// https://pan.baidu.com/union/doc/okumlx17r
type ApiErr struct {
	Errno  int32  `json:"errno,omitempty"`
	ErrMsg string `json:"errmsg"`
	Remark string
}

func (e *ApiErr) Error() string {
	return e.ErrMsg
}

// https://pan.baidu.com/union/doc/okumlx17r
var (
	// 公共错误码
	Success = ApiErr{
		Errno:  0,
		ErrMsg: "请求成功",
		Remark: "成功",
	}

	ErrUnknown = ApiErr{
		Errno:  1,
		ErrMsg: "未知错误",
		Remark: "服务器内部错误，请再次请求，如果持续出现此类错误，请提交工单反馈",
	}

	ApiErrParamError = ApiErr{
		Errno:  2,
		ErrMsg: "参数错误",
		Remark: "检查必选参数是否都已填写；检查参数位置，有的参数是在url里，有的是在body里；检查每个参数的值是否正确",
	}

	ErrUnsupportedMethod = ApiErr{
		Errno:  3,
		ErrMsg: "未知的方法",
		Remark: "调用的API不存在，请检查后重新尝试",
	}

	ErrRequestLimitReached = ApiErr{
		Errno:  4,
		ErrMsg: "接口请求过于频繁",
		Remark: "接口调用次数已达到设定的上限，请控制请求频率",
	}

	ErrUnauthorizedIP = ApiErr{
		Errno:  5,
		ErrMsg: "请求来自未经授权的IP地址",
		Remark: "请求IP不在白名单内",
	}

	ErrNoPermission = ApiErr{
		Errno:  6,
		ErrMsg: "无权限访问该用户数据",
		Remark: "应用无权限访问该用户数据，请为应用配置权限",
	}

	ErrInvalidParameter = ApiErr{
		Errno:  100,
		ErrMsg: "请求参数无效",
		Remark: "检查请求参数是否正确",
	}

	ErrInvalidApiKey = ApiErr{
		Errno:  101,
		ErrMsg: "API Key无效",
		Remark: "检查API Key是否正确",
	}

	ErrInvalidAccessToken = ApiErr{
		Errno:  110,
		ErrMsg: "无效的access token",
		Remark: "access token无效或已过期",
	}

	ErrAccessTokenExpired = ApiErr{
		Errno:  111,
		ErrMsg: "access token过期",
		Remark: "access token已过期，请重新获取",
	}

	// 百度网盘特有错误码
	ErrOperationFailed = ApiErr{
		Errno:  -1,
		ErrMsg: "操作失败",
		Remark: "检查网络连接是否正常，尝试重新登录或重启百度网盘客户端",
	}

	ErrFileNotFound = ApiErr{
		Errno:  -3,
		ErrMsg: "文件不存在",
		Remark: "文件不存在",
	}

	ErrAuthFailed = ApiErr{
		Errno:  -6,
		ErrMsg: "身份验证失败",
		Remark: "access_token 是否有效; 授权是否成功；参考接入授权FAQ；阅读文档《使用入门->接入授权》章节",
	}

	ErrPathError = ApiErr{
		Errno:  -7,
		ErrMsg: "文件或目录名错误或无权访问",
		Remark: "文件或目录名有误",
	}

	ApiErrPathExists = ApiErr{
		Errno:  -8,
		ErrMsg: "文件或目录已存在",
		Remark: "文件或目录已存在",
	}

	ApiErrPathNotFound = ApiErr{
		Errno:  -9,
		ErrMsg: "文件或目录不存在",
		Remark: "文件或目录不存在",
	}

	ErrSaveFileExists = ApiErr{
		Errno:  10,
		ErrMsg: "转存文件已经存在",
		Remark: "转存文件已经存在",
	}

	ErrUserNotExists = ApiErr{
		Errno:  11,
		ErrMsg: "用户不存在",
		Remark: "uid不存在",
	}

	ErrBatchSaveFailed = ApiErr{
		Errno:  12,
		ErrMsg: "批量转存出错",
		Remark: "参数错误，检查转存源和目的是不是同一个uid，正常不应该是一个 uid",
	}

	ErrAsyncTaskRunning = ApiErr{
		Errno:  111,
		ErrMsg: "有其他异步任务正在执行",
		Remark: "稍后，可重新请求",
	}

	ErrAppReviewing = ApiErr{
		Errno:  20011,
		ErrMsg: "应用审核中",
		Remark: "仅限前10个完成OAuth授权的用户测试应用",
	}

	ErrAccessLimited = ApiErr{
		Errno:  20012,
		ErrMsg: "访问超限",
		Remark: "调用次数已达上限，触发限流",
	}

	ErrPermissionDenied = ApiErr{
		Errno:  20013,
		ErrMsg: "权限不足",
		Remark: "当前应用无接口权限",
	}

	ErrInvalidParam = ApiErr{
		Errno:  31023,
		ErrMsg: "参数错误",
		Remark: "检查必选参数是否都已填写；检查参数位置，有的参数是在url里，有的是在body里；检查每个参数的值是否正确",
	}

	ErrNoAccessPermission = ApiErr{
		Errno:  31024,
		ErrMsg: "没有访问权限",
		Remark: "检查授权应用方式",
	}

	ErrApiFreqLimit9013 = ApiErr{
		Errno:  9013,
		ErrMsg: "命中接口频控",
		Remark: "接口请求过于频繁，注意控制",
	}

	ErrApiFreqLimit9019 = ApiErr{
		Errno:  9019,
		ErrMsg: "命中接口频控",
		Remark: "接口请求过于频繁，注意控制",
	}

	ErrApiFreqLimit = ApiErr{
		Errno:  31034,
		ErrMsg: "命中接口频控",
		Remark: "接口请求过于频繁，注意控制",
	}

	ErrAccessTokenInvalid = ApiErr{
		Errno:  31045,
		ErrMsg: "access_token验证未通过",
		Remark: "请检查access_token是否过期，用户授权时是否勾选网盘权限等",
	}

	ErrFileAlreadyExists = ApiErr{
		Errno:  31061,
		ErrMsg: "文件已存在",
		Remark: "文件已存在",
	}

	ErrInvalidFilename = ApiErr{
		Errno:  31062,
		ErrMsg: "文件名无效",
		Remark: "检查是否包含特殊字符",
	}

	ErrUploadPathError = ApiErr{
		Errno:  31064,
		ErrMsg: "上传路径错误",
		Remark: "上传文件的绝对路径格式：/apps/申请接入时填写的产品名称",
	}

	ErrFilenameNotFound = ApiErr{
		Errno:  31066,
		ErrMsg: "文件名不存在",
		Remark: "排查文件是否存储，路径是否传错",
	}

	ErrFileMissing = ApiErr{
		Errno:  31190,
		ErrMsg: "文件不存在",
		Remark: "block_list参数是否正确；一般是分片上传阶段有问题；检查分片上传阶段，分片传完了么；size大小对不对，跟实际文件是否一致，跟预上传接口的size是否一致",
	}

	ErrFirstChunkTooSmall = ApiErr{
		Errno:  31299,
		ErrMsg: "第一个分片的大小小于4MB",
		Remark: "要等于4MB",
	}

	ErrNotMediaFile = ApiErr{
		Errno:  31301,
		ErrMsg: "非音视频文件",
		Remark: "文件类型是否是音视频",
	}

	ErrInvalidVideoFormat = ApiErr{
		Errno:  31304,
		ErrMsg: "视频格式不支持播放",
		Remark: "视频格式不支持播放",
	}

	ErrAntiHotlink = ApiErr{
		Errno:  31326,
		ErrMsg: "命中防盗链",
		Remark: "查看自己请求是否合理，User-Agent请求头是否正常",
	}

	ErrVideoBitrateTooHigh = ApiErr{
		Errno:  31338,
		ErrMsg: "当前视频码率太高暂不支持流畅播放",
		Remark: "用户下载后播放",
	}

	ErrIllegalMediaFile = ApiErr{
		Errno:  31339,
		ErrMsg: "非法媒体文件",
		Remark: "检查视频内容",
	}

	ErrVideoTranscoding = ApiErr{
		Errno:  31341,
		ErrMsg: "视频正在转码",
		Remark: "可重新请求",
	}

	ErrVideoTranscodeFailed = ApiErr{
		Errno:  31346,
		ErrMsg: "视频转码失败",
		Remark: "排查该文件是否是个正常的视频",
	}

	ErrVideoTooLong = ApiErr{
		Errno:  31347,
		ErrMsg: "当前视频太长，暂不支持在线播放",
		Remark: "建议用户下载后播放",
	}

	ErrParamAbnormal = ApiErr{
		Errno:  31355,
		ErrMsg: "参数异常",
		Remark: "一般是 uploadid 参数传的有问题，确认uploadid参数传的是否与预上传precreate接口下发的uploadid一致",
	}

	ErrUrlExpired = ApiErr{
		Errno:  31360,
		ErrMsg: "url过期",
		Remark: "请重新获取",
	}

	ErrSignError = ApiErr{
		Errno:  31362,
		ErrMsg: "签名错误",
		Remark: "请检查链接地址是否完整",
	}

	ErrMissingChunks = ApiErr{
		Errno:  31363,
		ErrMsg: "分片缺失",
		Remark: "分片是否全部上传；每个上传的分片是否正确；size大小是否正确，跟实际文件是否一致，跟预上传接口的size是否一致",
	}

	ErrChunkSizeExceeded = ApiErr{
		Errno:  31364,
		ErrMsg: "超出分片大小限制",
		Remark: "建议以4MB作为上限",
	}

	ErrFileSizeExceeded = ApiErr{
		Errno:  31365,
		ErrMsg: "文件总大小超限",
		Remark: "授权用户为普通用户时，单个分片大小固定为4MB，单文件总大小上限为4GB；授权用户为普通会员时，单个分片大小",
	}
)

// 根据Errno获取对应的ApiErr实例
func GetApiErrByErrno(errno int32) *ApiErr {
	errMap := map[int32]*ApiErr{
		0:     &Success,
		1:     &ErrUnknown,
		2:     &ApiErrParamError,
		3:     &ErrUnsupportedMethod,
		4:     &ErrRequestLimitReached,
		5:     &ErrUnauthorizedIP,
		6:     &ErrNoPermission,
		10:    &ErrSaveFileExists,
		11:    &ErrUserNotExists,
		12:    &ErrBatchSaveFailed,
		100:   &ErrInvalidParameter,
		101:   &ErrInvalidApiKey,
		110:   &ErrInvalidAccessToken,
		111:   &ErrAsyncTaskRunning,
		-1:    &ErrOperationFailed,
		-3:    &ErrFileNotFound,
		-6:    &ErrAuthFailed,
		-7:    &ErrPathError,
		-8:    &ApiErrPathExists,
		-9:    &ApiErrPathNotFound,
		20011: &ErrAppReviewing,
		20012: &ErrAccessLimited,
		20013: &ErrPermissionDenied,
		31023: &ErrInvalidParam,
		31024: &ErrNoAccessPermission,
		31034: &ErrApiFreqLimit,
		31045: &ErrAccessTokenInvalid,
		31061: &ErrFileAlreadyExists,
		31062: &ErrInvalidFilename,
		31064: &ErrUploadPathError,
		31066: &ErrFilenameNotFound,
		31190: &ErrFileMissing,
		31299: &ErrFirstChunkTooSmall,
		31301: &ErrNotMediaFile,
		31304: &ErrInvalidVideoFormat,
		31326: &ErrAntiHotlink,
		31338: &ErrVideoBitrateTooHigh,
		31339: &ErrIllegalMediaFile,
		31341: &ErrVideoTranscoding,
		31346: &ErrVideoTranscodeFailed,
		31347: &ErrVideoTooLong,
		31355: &ErrParamAbnormal,
		31360: &ErrUrlExpired,
		31362: &ErrSignError,
		31363: &ErrMissingChunks,
		31364: &ErrChunkSizeExceeded,
		31365: &ErrFileSizeExceeded,
	}

	if err, ok := errMap[errno]; ok {
		return err
	}

	// 如果没有找到对应的错误码，返回未知错误
	unknownErr := ErrUnknown
	unknownErr.Remark = fmt.Sprintf("未知的错误码: %d", errno)
	return &unknownErr
}
