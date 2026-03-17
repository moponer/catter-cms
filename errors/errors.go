package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

const (
	// 通用错误码
	ErrCodeSuccess ErrorCode = 0
	ErrCodeBadRequest ErrorCode = 400
	ErrCodeUnauthorized ErrorCode = 401
	ErrCodeForbidden ErrorCode = 403
	ErrCodeNotFound ErrorCode = 404
	ErrCodeConflict ErrorCode = 409
	ErrCodeUnprocessableEntity ErrorCode = 422
	ErrCodeInternalError ErrorCode = 500
)

// AppError 应用错误结构
type AppError struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	StatusCode int         `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError 创建新的应用错误
func NewAppError(code ErrorCode, message string, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: getHTTPStatusFromCode(code),
	}
}

// NewBadRequestError 创建400错误
func NewBadRequestError(message string, details interface{}) *AppError {
	return NewAppError(ErrCodeBadRequest, message, details)
}

// NewUnauthorizedError 创建401错误
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, nil)
}

// NewForbiddenError 创建403错误
func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrCodeForbidden, message, nil)
}

// NewNotFoundError 创建404错误
func NewNotFoundError(resource string, id interface{}) *AppError {
	message := fmt.Sprintf("%s not found", resource)
	if id != nil {
		message = fmt.Sprintf("%s with id %v not found", resource, id)
	}
	return NewAppError(ErrCodeNotFound, message, nil)
}

// NewConflictError 创建409错误
func NewConflictError(message string) *AppError {
	return NewAppError(ErrCodeConflict, message, nil)
}

// NewUnprocessableEntityError 创建422错误
func NewUnprocessableEntityError(message string, details interface{}) *AppError {
	return NewAppError(ErrCodeUnprocessableEntity, message, details)
}

// NewInternalServerError 创建500错误
func NewInternalServerError(message string) *AppError {
	return NewAppError(ErrCodeInternalError, message, nil)
}

// WrapError 包装已有错误
func WrapError(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewInternalServerError(message)
}

// getHTTPStatusFromCode 根据错误码获取HTTP状态码
func getHTTPStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeSuccess:
		return http.StatusOK
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case ErrCodeInternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}
