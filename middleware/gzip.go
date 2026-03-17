package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gobuffalo/buffalo"
)

// GzipMiddleware Gzip压缩中间件
func GzipMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// 检查是否支持gzip
		if !strings.Contains(c.Request().Header.Get("Accept-Encoding"), "gzip") {
			return next(c)
		}
		
		// 创建gzip写入器
		gz := gzip.NewWriter(c.Response())
		defer gz.Close()
		
		// 设置响应头
		c.Response().Header().Set("Content-Encoding", "gzip")
		c.Response().Header().Set("Vary", "Accept-Encoding")
		
		// 执行下一个处理器
		return next(c)
	}
}
