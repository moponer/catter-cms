package api

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// Response 统一API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c buffalo.Context, data interface{}) error {
	return c.Render(http.StatusOK, render.JSON(Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}))
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c buffalo.Context, message string, data interface{}) error {
	return c.Render(http.StatusOK, render.JSON(Response{
		Code:    0,
		Message: message,
		Data:    data,
	}))
}

// Error 错误响应
func Error(c buffalo.Context, code int, message string) error {
	status := http.StatusBadRequest
	if code == 401 {
		status = http.StatusUnauthorized
	} else if code == 403 {
		status = http.StatusForbidden
	} else if code == 404 {
		status = http.StatusNotFound
	} else if code == 500 {
		status = http.StatusInternalServerError
	}

	return c.Render(status, render.JSON(Response{
		Code:    code,
		Message: message,
	}))
}

// BadRequest 400错误
func BadRequest(c buffalo.Context, message string) error {
	return Error(c, 400, message)
}

// Unauthorized 401错误
func Unauthorized(c buffalo.Context, message string) error {
	return Error(c, 401, message)
}

// Forbidden 403错误
func Forbidden(c buffalo.Context, message string) error {
	return Error(c, 403, message)
}

// NotFound 404错误
func NotFound(c buffalo.Context, message string) error {
	return Error(c, 404, message)
}

// InternalError 500错误
func InternalError(c buffalo.Context, message string) error {
	return Error(c, 500, message)
}

// PagedSuccess 分页成功响应
func PagedSuccess(c buffalo.Context, list interface{}, total, page, pageSize int) error {
	return Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// JSON 直接返回JSON响应
func JSON(c buffalo.Context, status int, data interface{}) error {
	return c.Render(status, render.JSON(data))
}
