package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the unified error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// CustomRecovery is a secure recovery middleware that:
// - Captures all panics and logs full stack traces server-side
// - Returns a generic error to the client (never exposes stack traces)
// - Works identically in both Debug and Release modes
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log full details server-side for debugging
				log.Printf("[PANIC RECOVER] %v\n%s", err, debug.Stack())

				// Return generic error to client - NEVER expose stack traces
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal Server Error",
					Message: "An unexpected error occurred. Please try again later.",
					Code:    http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ErrorHandler handles application errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("[ERROR] %v", e.Err)
			}
		}
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: "The requested resource does not exist: " + c.Request.URL.Path,
			Code:    http.StatusNotFound,
		})
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			Error:   "Method Not Allowed",
			Message: fmt.Sprintf("Method %s is not allowed for %s", c.Request.Method, c.Request.URL.Path),
			Code:    http.StatusMethodNotAllowed,
		})
	}
}

// JSONValidator validates request body JSON format
func JSONValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if contentType == "" {
			c.Next()
			return
		}

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
				if c.Request.Body != nil {
					var js interface{}
					if err := c.ShouldBindJSON(&js); err != nil {
						c.JSON(http.StatusBadRequest, ErrorResponse{
							Error:   "Invalid JSON",
							Message: "Request body must be valid JSON",
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

// RequestIDHandler generates a unique ID for each request
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

// generateRequestID generates a cryptographically secure request ID
func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback: should never happen with crypto/rand
		return fmt.Sprintf("fallback-%d", 0)
	}
	return hex.EncodeToString(b)
}
