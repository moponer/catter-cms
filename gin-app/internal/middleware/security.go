package middleware

import (
	"fmt"

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

		// 引用来源策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (CSP)
		csp := getCSPHeader(c)
		if csp != "" {
			c.Header("Content-Security-Policy", csp)
		}

		// 禁用缓存（对于 API 健康检查）
		if c.Request.URL.Path == "/api/health" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Header("Pragma", "no-cache")
		}

		c.Next()
	}
}

// getCSPHeader 根据环境获取 CSP 头
func getCSPHeader(c *gin.Context) string {
	// 开发环境使用宽松的 CSP
	if gin.Mode() == gin.DebugMode {
		return "default-src 'self' 'unsafe-inline' 'unsafe-eval'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' ws: wss: http: https:;"
	}

	// 生产环境使用严格的 CSP
	return "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'self'; base-uri 'self'; form-action 'self'"
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

// ForceHTTPS 强制 HTTPS 中间件
// 用于生产环境强制使用安全连接
func ForceHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否已经是 HTTPS
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Next()
			return
		}

		// 获取原始主机和路径
		host := c.GetHeader("X-Forwarded-Host")
		if host == "" {
			host = c.Request.Host
		}
		path := c.Request.URL.RequestURI()

		// 重定向到 HTTPS
		httpsURL := "https://" + host + path
		c.Redirect(301, httpsURL)
		c.Abort()
	}
}

// HSTSHeaders HSTS 安全头
// Strict-Transport-Security
func HSTSHeaders(maxAge int64, includeSubDomains bool, preload bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		hsts := fmt.Sprintf("max-age=%d", maxAge)
		if includeSubDomains {
			hsts += "; includeSubDomains"
		}
		if preload {
			hsts += "; preload"
		}
		c.Header("Strict-Transport-Security", hsts)
		c.Next()
	}
}
