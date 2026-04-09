package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"catter-cms-gin/config"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// Init 初始化 Redis 连接
func Init(cfg config.RedisConfig) error {
	if !cfg.Enabled {
		log.Println("Redis is disabled, using no-op cache")
		return nil
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		client = nil
		return fmt.Errorf("Redis connection failed (%s): %w", addr, err)
	}

	log.Printf("Redis connected: %s (DB %d)", addr, cfg.DB)
	return nil
}

// IsEnabled 检查 Redis 是否启用且连接正常
func IsEnabled() bool {
	return client != nil
}

// Reconnect 重新连接 Redis（配置变更后调用）
func Reconnect(cfg config.RedisConfig) error {
	// 关闭旧连接
	if client != nil {
		_ = client.Close()
		client = nil
	}

	return Init(cfg)
}

// Close 关闭 Redis 连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// prefix 包装 key 前缀
func prefix(key string) string {
	cfg := config.Get()
	return cfg.Redis.KeyPrefix + key
}

// Get 获取缓存值（反序列化为 JSON 目标）
func Get(key string, dest interface{}) error {
	if !IsEnabled() {
		return fmt.Errorf("cache: redis not enabled")
	}

	val, err := client.Get(ctx, prefix(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache: key not found")
		}
		return fmt.Errorf("cache: get failed: %w", err)
	}

	if dest != nil {
		if err := json.Unmarshal([]byte(val), dest); err != nil {
			return fmt.Errorf("cache: json unmarshal failed: %w", err)
		}
	}

	return nil
}

// GetString 获取字符串缓存值
func GetString(key string) (string, error) {
	if !IsEnabled() {
		return "", fmt.Errorf("cache: redis not enabled")
	}

	val, err := client.Get(ctx, prefix(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("cache: key not found")
		}
		return "", fmt.Errorf("cache: get failed: %w", err)
	}

	return val, nil
}

// Set 设置缓存值（自动 JSON 序列化）
func Set(key string, value interface{}, ttl time.Duration) error {
	if !IsEnabled() {
		return nil // 未启用时静默跳过
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: json marshal failed: %w", err)
	}

	if err := client.Set(ctx, prefix(key), data, ttl).Err(); err != nil {
		return fmt.Errorf("cache: set failed: %w", err)
	}

	return nil
}

// SetString 设置字符串缓存值
func SetString(key string, value string, ttl time.Duration) error {
	if !IsEnabled() {
		return nil
	}

	if err := client.Set(ctx, prefix(key), value, ttl).Err(); err != nil {
		return fmt.Errorf("cache: set failed: %w", err)
	}

	return nil
}

// Del 删除缓存
func Del(keys ...string) error {
	if !IsEnabled() {
		return nil
	}

	prefixed := make([]string, len(keys))
	for i, k := range keys {
		prefixed[i] = prefix(k)
	}

	if err := client.Del(ctx, prefixed...).Err(); err != nil {
		return fmt.Errorf("cache: del failed: %w", err)
	}

	return nil
}

// InvalidatePattern 按模式批量删除缓存（如 "post:*"）
func InvalidatePattern(pattern string) error {
	if !IsEnabled() {
		return nil
	}

	p := prefix(pattern)
	iter := client.Scan(ctx, 0, p, 100).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if len(keys) > 0 {
		if err := client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("cache: invalidate pattern failed: %w", err)
		}
		log.Printf("Cache invalidated: pattern=%s, keys=%d", p, len(keys))
	}

	return nil
}

// FlushDB 清空当前数据库的所有缓存
func FlushDB() error {
	if !IsEnabled() {
		return nil
	}

	return client.FlushDB(ctx).Err()
}

// Stats 返回缓存统计信息
func Stats() map[string]interface{} {
	if !IsEnabled() {
		return map[string]interface{}{
			"enabled":  false,
			"status":   "disabled",
			"keys":     0,
			"hit_rate": "N/A",
		}
	}

	dbSize, _ := client.DBSize(ctx).Result()

	return map[string]interface{}{
		"enabled": true,
		"status":  "connected",
		"keys":    dbSize,
	}
}
