package services

import (
	"sync"
	"time"
)

// ThemeCache 主题配置缓存
type ThemeCache struct {
	cache     sync.Map
	duration  time.Duration
	enabled   bool
	mu        sync.RWMutex
}

// CacheItem 缓存项
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// 全局主题缓存实例
var themeCache = &ThemeCache{
	duration: 5 * time.Minute,
	enabled:  true,
}

// Get 从缓存获取配置
func (c *ThemeCache) Get(key string) (interface{}, bool) {
	if !c.enabled {
		return nil, false
	}

	value, ok := c.cache.Load(key)
	if !ok {
		return nil, false
	}

	item, ok := value.(CacheItem)
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(item.ExpiresAt) {
		c.cache.Delete(key)
		return nil, false
	}

	return item.Data, true
}

// Set 设置缓存
func (c *ThemeCache) Set(key string, data interface{}) {
	if !c.enabled {
		return
	}

	item := CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(c.duration),
	}

	c.cache.Store(key, item)
}

// Delete 删除缓存
func (c *ThemeCache) Delete(key string) {
	c.cache.Delete(key)
}

// Clear 清空所有缓存
func (c *ThemeCache) Clear() {
	c.cache.Range(func(key, value interface{}) bool {
		c.cache.Delete(key)
		return true
	})
}

// SetDuration 设置缓存时长
func (c *ThemeCache) SetDuration(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.duration = duration
}

// Enable 启用缓存
func (c *ThemeCache) Enable() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enabled = true
}

// Disable 禁用缓存
func (c *ThemeCache) Disable() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enabled = false
}

// GetGlobalThemeCache 获取全局主题缓存实例
func GetGlobalThemeCache() *ThemeCache {
	return themeCache
}
