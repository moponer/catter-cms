package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig 日志中间件配置
type LoggerConfig struct {
	// 输出文件路径，空为标准输出
	Output string
	// 自定义格式化函数
	Formatter func(param gin.LogFormatterParams) string
}

// DefaultFormatter 默认日志格式化函数
func DefaultFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] %s | %3d | %13v | %15s | %-7s %s\n",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.Method,
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
	)
}

// Logger 日志中间件
// 记录每个HTTP请求的详细信息
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithConfig 使用配置创建日志中间件
func LoggerWithConfig(cfg LoggerConfig) gin.HandlerFunc {
	formatter := cfg.Formatter
	if formatter == nil {
		formatter = DefaultFormatter
	}

	var output *os.File
	if cfg.Output != "" {
		var err error
		output, err = os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			output = os.Stdout
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		stop := time.Now()
		latency := stop.Sub(start)

		// 获取请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// 格式化日志
		param := gin.LogFormatterParams{
			TimeStamp:   stop,
			Latency:     latency,
			ClientIP:    clientIP,
			Method:      method,
			StatusCode:  statusCode,
			Path:        path,
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		// 添加查询参数
		if query != "" {
			param.Path = path + "?" + query
		}

		// 写入日志
		fmt.Fprint(output, formatter(param))
	}
}

// RequestLogger 请求日志中间件（简化版）
// 只记录基本请求信息
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		fmt.Printf("[%s] %s %s %d %v %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			latency,
			clientIP,
		)
	}
}

// AccessLogWriter 自定义响应写入器，用于捕获响应状态
type AccessLogWriter struct {
	gin.ResponseWriter
	StatusCode int
	Body       []byte
}

// Write 写入响应
func (w *AccessLogWriter) Write(b []byte) (int, error) {
	w.Body = b
	return w.ResponseWriter.Write(b)
}

// WriteHeader 写入响应头
func (w *AccessLogWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// AccessLog 访问日志中间件
// 记录更详细的访问信息，包括响应状态和内容长度
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 包装响应写入器
		writer := &AccessLogWriter{
			ResponseWriter: c.Writer,
			StatusCode:     200,
		}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		contentLength := c.Writer.Size()

		fmt.Printf("[ACCESS] %s | %3d | %12v | %15s | %-7s %s | %s | %d bytes\n",
			time.Now().Format("2006-01-02 15:04:05"),
			writer.StatusCode,
			latency,
			clientIP,
			method,
			path,
			userAgent,
			contentLength,
		)
	}
}

// StructuredLogger 结构化日志中间件
// 输出JSON格式的日志，便于日志收集和分析
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()
		contentLength := c.Writer.Size()

		if query != "" {
			path = path + "?" + query
		}

		// 输出JSON格式日志
		fmt.Printf(`{"time":"%s","client_ip":"%s","method":"%s","path":"%s","status":%d,"latency":%d,"user_agent":"%s","content_length":%d}`+"\n",
			time.Now().Format(time.RFC3339),
			clientIP,
			method,
			path,
			statusCode,
			latency.Milliseconds(),
			userAgent,
			contentLength,
		)
	}
}
