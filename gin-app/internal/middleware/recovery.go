package middleware

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
}

// Recovery 恢复中间件
// 捕获所有panic并返回统一的错误响应
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				log.Printf("[PANIC RECOVER] %v\n%s", err, debug.Stack())

				// 构建错误响应
				errMsg := fmt.Sprintf("%v", err)
				response := ErrorResponse{
					Error:   "Internal Server Error",
					Message: errMsg,
					Code:    http.StatusInternalServerError,
				}

				// 如果是开发模式，显示详细信息
				if gin.Mode() == gin.DebugMode {
					response.Details = gin.H{
						"stack": string(debug.Stack()),
						"error": errMsg,
					}
				}

				// 发送错误响应
				c.JSON(http.StatusInternalServerError, response)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ErrorHandler 错误处理中间件
// 统一处理应用程序错误
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("[ERROR] %v", e.Err)
			}
		}
	}
}

// NotFoundHandler 404处理
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: "请求的资源不存在: " + c.Request.URL.Path,
			Code:    http.StatusNotFound,
		})
	}
}

// MethodNotAllowedHandler 405处理
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			Error:   "Method Not Allowed",
			Message: fmt.Sprintf("不允许使用 %s 方法访问 %s", c.Request.Method, c.Request.URL.Path),
			Code:    http.StatusMethodNotAllowed,
		})
	}
}

// JSONValidator JSON验证中间件
// 验证请求体的JSON格式
func JSONValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只检查带有Body的请求
		if c.Request.Body == nil {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if contentType == "" {
			c.Next()
			return
		}

		// 检查是否是JSON请求
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
				// 尝试解析JSON以验证格式
				if c.Request.Body != nil {
					var js interface{}
					if err := c.ShouldBindJSON(&js); err != nil {
						c.JSON(http.StatusBadRequest, ErrorResponse{
							Error:   "Invalid JSON",
							Message: "请求体必须是有效的JSON格式",
							Code:    http.StatusBadRequest,
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// RequestIDHandler 请求ID中间件
// 为每个请求生成唯一的ID，便于追踪
func RequestIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成唯一的请求ID
func generateRequestID() string {
	// 使用时间戳和随机数生成唯一ID
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), getRandomInt())
}

// getRandomInt 获取随机整数（简化实现）
func getRandomInt() int {
	// 简单实现，实际项目可使用 crypto/rand
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100000)
}
