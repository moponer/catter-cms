package api

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/cache"
)

// redisMu 保护动态切换 Redis 时的并发安全
var redisMu sync.RWMutex

// GetRedisConfig 获取 Redis 配置（密码脱敏）
func GetRedisConfig(c *gin.Context) {
	cfg := config.Get()
	r := cfg.Redis

	result := gin.H{
		"enabled":    r.Enabled,
		"host":       r.Host,
		"port":       r.Port,
		"password":   maskSecret(r.Password),
		"db":         r.DB,
		"key_prefix": r.KeyPrefix,
	}

	// 附加当前状态
	stats := cache.Stats()
	result["status"] = stats["status"]
	if keys, ok := stats["keys"]; ok {
		result["keys"] = keys
	}

	SuccessResponse(c, result)
}

// redisConfigRequest 前端提交的 Redis 配置
type redisConfigRequest struct {
	Enabled   bool   `json:"enabled"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Password  string `json:"password"`
	DB        int    `json:"db"`
	KeyPrefix string `json:"key_prefix"`
}

// SaveRedisConfig 保存 Redis 配置并动态重连
func SaveRedisConfig(c *gin.Context) {
	var req redisConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 合并密码：如果提交的是脱敏值（含 *），保留原值
	cfg := config.Get()
	if strings.Contains(req.Password, "*") {
		req.Password = cfg.Redis.Password
	}

	// 验证配置
	if req.Enabled {
		if req.Host == "" {
			BadRequest(c, "请填写 Redis 地址")
			return
		}
		if req.Port <= 0 || req.Port > 65535 {
			BadRequest(c, "Redis 端口无效（1-65535）")
			return
		}
	}

	// 构建 RedisConfig
	redisCfg := config.RedisConfig{
		Enabled:   req.Enabled,
		Host:      req.Host,
		Port:      req.Port,
		Password:  req.Password,
		DB:        req.DB,
		KeyPrefix: req.KeyPrefix,
	}
	if redisCfg.KeyPrefix == "" {
		redisCfg.KeyPrefix = "catter:"
	}

	// 如果启用，先测试连接
	if req.Enabled {
		if err := cache.Reconnect(redisCfg); err != nil {
			BadRequest(c, "Redis 连接失败: "+err.Error())
			return
		}
	} else {
		// 禁用则关闭连接
		cache.Close()
	}

	// 写入 config.yaml
	if err := writeRedisToConfigYAML(&req); err != nil {
		InternalError(c, "写入配置文件失败: "+err.Error())
		return
	}

	// 更新内存中的 config
	redisMu.Lock()
	cfg.Redis = redisCfg
	redisMu.Unlock()

	status := "disabled"
	if req.Enabled {
		status = "connected"
	}
	log.Printf("Redis config saved, status: %s", status)

	// 返回脱敏后的配置
	GetRedisConfig(c)
}

// TestRedisConnection 测试 Redis 连接
func TestRedisConnection(c *gin.Context) {
	var req redisConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 合并密码
	cfg := config.Get()
	if strings.Contains(req.Password, "*") {
		req.Password = cfg.Redis.Password
	}

	// 验证必填
	if req.Host == "" {
		BadRequest(c, "请填写 Redis 地址")
		return
	}
	if req.Port <= 0 || req.Port > 65535 {
		BadRequest(c, "Redis 端口无效")
		return
	}

	redisCfg := config.RedisConfig{
		Enabled:   true,
		Host:      req.Host,
		Port:      req.Port,
		Password:  req.Password,
		DB:        req.DB,
		KeyPrefix: req.KeyPrefix,
	}
	if redisCfg.KeyPrefix == "" {
		redisCfg.KeyPrefix = "catter:"
	}

	// 临时连接测试（不影响当前连接）
	if err := cache.Reconnect(redisCfg); err != nil {
		BadRequest(c, "连接失败: "+err.Error())
		return
	}

	// 获取统计信息
	stats := cache.Stats()

	SuccessWithMessage(c, "Redis 连接测试成功", gin.H{
		"status": stats["status"],
		"keys":   stats["keys"],
	})
}

// writeRedisToConfigYAML 将 Redis 配置写入 config.yaml（只修改 redis 段）
func writeRedisToConfigYAML(req *redisConfigRequest) error {
	configPath := findConfigPath()
	if configPath == "" {
		configPath = "config.yaml"
	}

	// 读取现有 YAML
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 反序列化为 map 保留所有字段
	var cfgMap map[string]interface{}
	if err := yaml.Unmarshal(data, &cfgMap); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 构建 redis 段
	prefix := req.KeyPrefix
	if prefix == "" {
		prefix = "catter:"
	}

	redisMap := map[string]interface{}{
		"enabled":    req.Enabled,
		"host":       req.Host,
		"port":       req.Port,
		"password":   req.Password,
		"db":         req.DB,
		"key_prefix": prefix,
	}

	// 如果密码为空，不写入（YAML 不输出空字符串字段）
	if req.Password == "" {
		delete(redisMap, "password")
	}
	if req.DB == 0 {
		redisMap["db"] = 0
	}
	if req.Port == 0 {
		redisMap["port"] = 6379
	}

	// 更新 redis 段
	cfgMap["redis"] = redisMap

	// 重新序列化
	out, err := yaml.Marshal(cfgMap)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写回文件
	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
