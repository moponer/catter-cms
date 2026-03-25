package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 安全头中间件
// 添加各种安全相关的 HTTP 头
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止 MIME 类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 防止点击劫持
		c.Header("X-Frame-Options", "SAMEORIGIN")

		// XSS 保护
		c.Header("X-XSS-Protection", "1; mode=block")

		// 禁用 IE 的兼容模式
		c.Header("X-UA-Compatible", "IE=edge")

		// 强制使用 HTTPS（如果配置了）
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 引用来源策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 禁用缓存（对于 API，敏感页面）
		// 对于静态资源，应该在路由层面单独处理缓存
		if c.Request.URL.Path == "/api/health" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Header("Pragma", "no-cache")
		}

		c.Next()
	}
}

// CacheHeaders 缓存头中间件
// 用于静态资源和公开页面
func CacheHeaders(maxAge string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age="+maxAge)
		c.Next()
	}
}

// NoCacheHeaders 禁用缓存头
func NoCacheHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}
