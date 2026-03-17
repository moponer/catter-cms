package middleware

import (
	"fmt"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/patrickmn/go-cache"
)

var (
	// defaultCache 默认缓存
	defaultCache *cache.Cache
)

// init 初始化缓存
func init() {
	defaultCache = cache.New(5*time.Minute, 10*time.Minute)
}

// CacheMiddleware 缓存中间件
func CacheMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		cacheKey := c.Request().URL.String()
		
		// 检查缓存
		if _, found := defaultCache.Get(cacheKey); found {
			c.Response().Header().Set("X-Cache", "HIT")
		}
		
		// 执行下一个处理器
		err := next(c)
		
		// 缓存结果（仅缓存成功的GET请求）
		if err == nil && c.Request().Method == "GET" {
			c.Response().Header().Set("X-Cache", "MISS")
		}
		
		return err
	}
}

// CacheKey 生成缓存键
func CacheKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

// ClearCache 清除缓存
func ClearCache(key string) {
	cacheKey := CacheKey("app", key)
	defaultCache.Delete(cacheKey)
}

// ClearCacheByPrefix 按前缀清除缓存
func ClearCacheByPrefix(prefix string) {
	items := defaultCache.Items()
	for k := range items {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			defaultCache.Delete(k)
		}
	}
}
