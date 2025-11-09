package errors

// ErrorCode represents an application error code
type ErrorCode int

// Success code
const (
	Success ErrorCode = 0
)

// Client error codes (1000-1999)
const (
	// Common client errors
	BadRequest    ErrorCode = 1000 // 请求参数错误
	Unauthorized  ErrorCode = 1001 // 未授权
	Forbidden     ErrorCode = 1003 // 禁止访问
	NotFound      ErrorCode = 1004 // 资源不存在
	Validation    ErrorCode = 1005 // 数据验证失败
	Conflict      ErrorCode = 1009 // 资源冲突
	TooManyReq    ErrorCode = 1029 // 请求过于频繁
	InvalidToken  ErrorCode = 1030 // 无效的令牌
	ExpiredToken  ErrorCode = 1031 // 令牌已过期
	InvalidFormat ErrorCode = 1032 // 数据格式错误

	// User related errors
	UserNotFound     ErrorCode = 1100 // 用户不存在
	InvalidPassword  ErrorCode = 1101 // 密码错误
	UserExists       ErrorCode = 1102 // 用户已存在
	InvalidUserState ErrorCode = 1103 // 用户状态异常

	// Auth related errors
	InvalidCredentials ErrorCode = 1200 // 无效的凭证
	SessionExpired     ErrorCode = 1201 // 会话已过期
	InvalidSession     ErrorCode = 1202 // 无效的会话
)

// Server error codes (5000-5999)
const (
	// Common server errors
	InternalServer ErrorCode = 5000 // 服务器内部错误
	Database       ErrorCode = 5001 // 数据库错误
	Cache          ErrorCode = 5002 // 缓存错误
	NotImplemented ErrorCode = 5003 // 功能未实现
	ThirdParty     ErrorCode = 5004 // 第三方服务错误
	Config         ErrorCode = 5005 // 配置错误
	Network        ErrorCode = 5006 // 网络错误
	Timeout        ErrorCode = 5007 // 超时错误
)

// Standard error messages
var standardMessages = map[ErrorCode]string{
	Success:       "success",
	BadRequest:    "请求参数错误",
	Unauthorized:  "未授权",
	Forbidden:     "禁止访问",
	NotFound:      "资源不存在",
	Validation:    "数据验证失败",
	Conflict:      "资源冲突",
	TooManyReq:    "请求过于频繁",
	InvalidToken:  "无效的令牌",
	ExpiredToken:  "令牌已过期",
	InvalidFormat: "数据格式错误",

	UserNotFound:     "用户不存在",
	InvalidPassword:  "密码错误",
	UserExists:       "用户已存在",
	InvalidUserState: "用户状态异常",

	InvalidCredentials: "无效的凭证",
	SessionExpired:     "会话已过期",
	InvalidSession:     "无效的会话",

	InternalServer: "服务器内部错误",
	Database:       "数据库错误",
	Cache:          "缓存错误",
	NotImplemented: "功能未实现",
	ThirdParty:     "第三方服务错误",
	Config:         "配置错误",
	Network:        "网络错误",
	Timeout:        "超时错误",
}

// GetStandardMessage returns the standard message for an error code
func GetStandardMessage(code ErrorCode) string {
	if msg, ok := standardMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// IsClientError checks if the error code is a client error
func IsClientError(code ErrorCode) bool {
	return code >= 1000 && code < 2000
}

// IsServerError checks if the error code is a server error
func IsServerError(code ErrorCode) bool {
	return code >= 5000 && code < 6000
}

// IsSuccess checks if the error code indicates success
func IsSuccess(code ErrorCode) bool {
	return code == Success
}
