package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GzipConfig Gzip中间件配置
type GzipConfig struct {
	// 压缩级别 (1-9)
	Level int
	// 不压缩的最小内容长度
	MinLength int
	// 不压缩的内容类型
	ExcludedExtensions []string
	// 不压缩的路径前缀
	ExcludedPaths []string
}

// DefaultGzipConfig 默认Gzip配置
func DefaultGzipConfig() GzipConfig {
	return GzipConfig{
		Level:           gzip.DefaultCompression,
		MinLength:       1024, // 最小1KB
		ExcludedExtensions: []string{
			".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg",
			".mp3", ".mp4", ".webm", ".avi",
			".zip", ".rar", ".7z", ".tar", ".gz",
			".pdf", ".doc", ".docx", ".xls", ".xlsx",
		},
		ExcludedPaths: []string{
			"/api/health",
			"/api/ping",
		},
	}
}

// Gzip Gzip压缩中间件
// 对响应内容进行gzip压缩，节省带宽
func Gzip() gin.HandlerFunc {
	return GzipWithConfig(DefaultGzipConfig())
}

// GzipWithConfig 使用配置创建Gzip中间件
func GzipWithConfig(config GzipConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否支持gzip
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next()
			return
		}

		// 检查是否排除该路径
		path := c.Request.URL.Path
		for _, excludedPath := range config.ExcludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				c.Next()
				return
			}
		}

		// 检查扩展名是否排除
		for _, ext := range config.ExcludedExtensions {
			if strings.HasSuffix(path, ext) {
				c.Next()
				return
			}
		}

		// 创建gzip写入器
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		// 设置响应头
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// 包装响应写入器以跟踪实际写入的字节数
		ginWriter := &gzipResponseWriter{
			ResponseWriter: c.Writer,
			Writer:         gz,
		}
		c.Writer = ginWriter

		c.Next()

		// 如果内容太短，不使用压缩
		if ginWriter.written < config.MinLength {
			// 恢复原始写入器
			c.Writer = c.Writer.(*gzipResponseWriter).ResponseWriter
			c.Header("Content-Encoding", "")
		}
	}
}

// gzipResponseWriter 包装响应写入器
type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer   *gzip.Writer
	written int
}

// Write 写入数据
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	w.written += n
	return n, err
}

// WriteHeader 写入状态码
func (w *gzipResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

// GzipDecompress Gzip解压中间件
// 用于处理客户端发送的gzip压缩请求
func GzipDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncoding := c.GetHeader("Content-Encoding")

		if strings.Contains(contentEncoding, "gzip") {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid gzip content",
				})
				c.Abort()
				return
			}
			defer reader.Close()

			// 替换body
			body, err := io.ReadAll(reader)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to decompress gzip content",
				})
				c.Abort()
				return
			}

			// 创建新的Reader以支持多次读取
			c.Request.Body = io.NopCloser(strings.NewReader(string(body)))
		}

		c.Next()
	}
}

// CompressLevel 自定义压缩级别
// level: 1-9 (gzip.NoCompression 到 gzip.BestCompression)
func CompressLevel(level int) gin.HandlerFunc {
	config := DefaultGzipConfig()
	config.Level = level
	return GzipWithConfig(config)
}

// BestCompression 最佳压缩
func BestCompression() gin.HandlerFunc {
	config := DefaultGzipConfig()
	config.Level = gzip.BestCompression
	return GzipWithConfig(config)
}

// FastCompression 快速压缩
func FastCompression() gin.HandlerFunc {
	config := DefaultGzipConfig()
	config.Level = 1 // 最低压缩级别，最快
	return GzipWithConfig(config)
}
