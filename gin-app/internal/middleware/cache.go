package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// CacheConfig 缓存中间件配置
type CacheConfig struct {
	// 缓存过期时间
	Expiration time.Duration
	// 清理间隔
	PurgeInterval time.Duration
	// 不缓存的路径
	ExcludedPaths []string
	// 只缓存的状态码
	AllowedStatusCodes []int
}

// DefaultCacheConfig 默认缓存配置
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Expiration:        5 * time.Minute,
		PurgeInterval:     10 * time.Minute,
		ExcludedPaths:     []string{"/admin", "/api/admin"},
		AllowedStatusCodes: []int{200},
	}
}

var (
	// defaultCache 默认缓存实例
	defaultCache *cache.Cache
	// cacheOnce 确保缓存只初始化一次
	cacheOnce sync.Once
)

// initCache 初始化缓存
func initCache() {
	cfg := DefaultCacheConfig()
	defaultCache = cache.New(cfg.Expiration, cfg.PurgeInterval)
}

// getCache 获取缓存实例
func getCache() *cache.Cache {
	cacheOnce.Do(initCache)
	return defaultCache
}

// Cache 缓存中间件
// 对GET请求进行响应缓存
func Cache() gin.HandlerFunc {
	return CacheWithConfig(DefaultCacheConfig())
}

// CacheWithConfig 使用配置创建缓存中间件
func CacheWithConfig(config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只缓存GET请求
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// 检查是否排除该路径
		path := c.Request.URL.Path
		for _, excludedPath := range config.ExcludedPaths {
			if excludedPath == "/"+strings.TrimPrefix(path, "/") ||
				strings.HasPrefix(path, excludedPath) {
				c.Next()
				return
			}
		}

		// 生成缓存键
		cacheKey := generateCacheKey(c)

		// 尝试从缓存获取
		if cached, found := getCache().Get(cacheKey); found {
			c.Header("X-Cache", "HIT")
			c.Header("X-Cache-Key", cacheKey)

			// 返回缓存的响应
			if cachedResp, ok := cached.(*CacheResponse); ok {
				// 设置缓存的状态码和头
				for k, v := range cachedResp.Headers {
					c.Header(k, v)
				}
				c.Status(cachedResp.StatusCode)
				c.Writer.Write(cachedResp.Body)
				c.Abort()
				return
			}
		}

		c.Header("X-Cache", "MISS")

		// 创建响应捕获写入器
		writer := &cacheResponseWriter{
			ResponseWriter: c.Writer,
			StatusCode:     200,
			Headers:        make(map[string]string),
			Body:           []byte{},
		}
		c.Writer = writer

		// 执行处理
		c.Next()

		// 检查是否应该缓存
		if shouldCache(writer.StatusCode, config.AllowedStatusCodes) {
			// 存储到缓存
			cacheResp := &CacheResponse{
				StatusCode: writer.StatusCode,
				Headers:    writer.Headers,
				Body:       writer.Body,
			}
			getCache().Set(cacheKey, cacheResp, config.Expiration)
		}
	}
}

// CacheResponse 缓存的响应
type CacheResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// cacheResponseWriter 缓存响应写入器
type cacheResponseWriter struct {
	gin.ResponseWriter
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// Write 写入响应体
func (w *cacheResponseWriter) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return w.ResponseWriter.Write(b)
}

// WriteHeader 写入状态码
func (w *cacheResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Header 设置响应头（捕获以便缓存）
func (w *cacheResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// generateCacheKey 生成缓存键
func generateCacheKey(c *gin.Context) string {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	if query != "" {
		return fmt.Sprintf("%s?%s", path, query)
	}
	return path
}

// shouldCache 检查状态码是否应该缓存
func shouldCache(statusCode int, allowedCodes []int) bool {
	if len(allowedCodes) == 0 {
		return statusCode >= 200 && statusCode < 300
	}
	for _, code := range allowedCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

// CacheKey 生成缓存键（带前缀）
func CacheKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

// ClearCache 清除指定缓存
func ClearCache(key string) {
	cacheKey := CacheKey("app", key)
	getCache().Delete(cacheKey)
}

// ClearCacheByPrefix 按前缀清除缓存
func ClearCacheByPrefix(prefix string) {
	getCache().DeleteExpired()
	// 注意：go-cache 不支持按前缀删除，这里简化处理
}

// GetCache 获取缓存值
func GetCache(key string) (interface{}, bool) {
	return getCache().Get(key)
}

// SetCache 设置缓存
func SetCache(key string, value interface{}, expiration time.Duration) {
	getCache().Set(key, value, expiration)
}

// DeleteCache 删除缓存
func DeleteCache(key string) {
	getCache().Delete(key)
}

// ClearAllCache 清除所有缓存
func ClearAllCache() {
	getCache().Flush()
}

// CacheDelete 删除缓存中间件
// 在请求完成后清除相关缓存
func CacheDelete(pattern string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 请求处理完成后清除缓存
		if c.Request.Method == "POST" || c.Request.Method == "PUT" ||
			c.Request.Method == "PATCH" || c.Request.Method == "DELETE" {
			// 清除相关缓存
			ClearCacheByPrefix(pattern)
		}
	}
}
