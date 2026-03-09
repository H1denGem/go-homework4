package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 自定义错误类型
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// 创建错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrBadRequest       = NewAppError(http.StatusBadRequest, "请求参数错误", nil)
	ErrInternalServer   = NewAppError(http.StatusInternalServerError, "服务器内部错误", nil)
	ErrInvalidToken     = NewAppError(http.StatusUnauthorized, "无效的令牌", nil)
	ErrMissingToken     = NewAppError(http.StatusUnauthorized, "缺少认证令牌", nil)
	ErrUserExists       = NewAppError(http.StatusBadRequest, "用户已存在", nil)
	ErrUserNotFound     = NewAppError(http.StatusNotFound, "用户不存在", nil)
	ErrInvalidPassword  = NewAppError(http.StatusUnauthorized, "密码错误", nil)
	ErrPostNotFound     = NewAppError(http.StatusNotFound, "文章不存在", nil)
	ErrCommentNotFound  = NewAppError(http.StatusNotFound, "评论不存在", nil)
	ErrPermissionDenied = NewAppError(http.StatusForbidden, "无权操作此资源", nil)
)

// 处理错误并返回合适的 HTTP 状态码
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.Code, appErr.Message)
		return
	}

	// 未知错误，返回 500
	Error(c, http.StatusInternalServerError, "服务器内部错误")
}

// 检查错误是否为特定类型
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// 包装错误
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return NewAppError(http.StatusInternalServerError, message, err)
}
