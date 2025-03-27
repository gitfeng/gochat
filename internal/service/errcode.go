package service

// 定义错误码常量
const (
	ErrCodeSuccess        = 0
	ErrCodeInternalServer = 500
	ErrCodeInvalidRequest = 400
	ErrCodeUnauthorized   = 401
	ErrCodeForbidden      = 403
	ErrCodeNotFound       = 404
	ErrCodeConflict       = 409
	// 可以根据业务需求添加更多自定义错误码
	ErrCodeUserExists      = 1001
	ErrCodeUserNotFound    = 1002
	ErrCodeInvalidPassword = 1003
	ErrCodeInvalidToken    = 1004
)

// 定义错误码对应的错误信息
var errorMessages = map[int]string{
	ErrCodeSuccess:         "成功",
	ErrCodeInternalServer:  "服务器内部错误",
	ErrCodeInvalidRequest:  "无效请求",
	ErrCodeUnauthorized:    "未授权",
	ErrCodeForbidden:       "禁止访问",
	ErrCodeNotFound:        "未找到资源",
	ErrCodeConflict:        "资源冲突",
	ErrCodeUserExists:      "用户已存在",
	ErrCodeUserNotFound:    "用户未找到",
	ErrCodeInvalidPassword: "密码无效",
	ErrCodeInvalidToken:    "无效的令牌",
}

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code int) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
