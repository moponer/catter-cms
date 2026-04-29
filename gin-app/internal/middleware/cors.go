package middleware

import (
	"catter-cms-gin/config"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS中间件配置
type CORSConfig struct {
	// 允许的来源
	AllowOrigins []string
	// 允许的HTTP方法
	AllowMethods []string
	// 允许的HTTP头
	AllowHeaders []string
	// 允许暴露的HTTP头
	ExposeHeaders []string
	// 是否允许携带凭证
	AllowCredentials bool
	// 预检请求缓存时间
	MaxAge int
}

// DefaultCORSConfig 默认CORS配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders: []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge: 12,
	}
}

// CORS CORS middleware
// Reads allowed origins from config.Security.AllowOrigins
func CORS() gin.HandlerFunc {
	cfg := config.Get()
	allowOrigins := cfg.Security.AllowOrigins
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	corsConfig := CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3600,
	}
	return CORSWithConfig(corsConfig)
}

// CORSWithConfig 使用配置创建CORS中间件
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许该来源
		allowed := false
		for _, allowOrigin := range config.AllowOrigins {
			if allowOrigin == "*" || allowOrigin == origin {
				allowed = true
				break
			}
			// 支持通配符子域名
			if strings.HasPrefix(allowOrigin, "*.") {
				suffix := strings.TrimPrefix(allowOrigin, "*.")
				if strings.HasSuffix(origin, suffix) {
					allowed = true
					break
				}
			}
		}

		if allowed {
			// 设置允许的来源
			if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			}

			// 设置允许的方法
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))

			// 设置允许的请求头
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))

			// 设置允许暴露的响应头
			if len(config.ExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
			}

			// 设置是否允许携带凭证
			if config.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			// 设置预检请求缓存时间
			if config.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
			}
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			// 设置Vary头
			c.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

			// 直接返回204状态码
			c.Status(204)
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecureCORS 安全CORS中间件
// 仅允许特定来源
func SecureCORS(allowedOrigins []string) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowOrigins = allowedOrigins
	config.AllowCredentials = true
	return CORSWithConfig(config)
}

// API CORS API专用的CORS中间件
func API_CORS() gin.HandlerFunc {
	config := CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders: []string{"Content-Length", "Content-Type", "Authorization", "X-Total-Count", "X-Page", "X-Per-Page"},
		AllowCredentials: true,
		MaxAge: 3600,
	}
	return CORSWithConfig(config)
}

// AdminCORS 管理后台专用的CORS中间件
func AdminCORS() gin.HandlerFunc {
	config := CORSConfig{
		AllowOrigins: []string{}, // 需要根据实际情况配置
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders: []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge: 3600,
	}
	return CORSWithConfig(config)
}

// CorsHandler 简单CORS处理（基于gin内置）
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
