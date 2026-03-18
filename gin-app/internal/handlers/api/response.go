package api

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, 400, message)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, 401, message)
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, 403, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, 404, message)
}

// InternalError 500 错误
func InternalError(c *gin.Context, message string) {
	ErrorResponse(c, 500, message)
}
