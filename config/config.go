package config

import (
	"github.com/gobuffalo/envy"
)

// Config 应用配置
type Config struct {
	ServerPort    string
	DatabaseDSN   string
	SessionSecret string
	CacheEnabled  bool
	CacheTTL      int
	Environment   string
	LogLevel      string
}

var (
	// AppConfig 应用配置实例
	AppConfig Config
)

// LoadConfig 加载配置
func LoadConfig() error {
	// 使用 envy.Get 获取环境变量
	AppConfig = Config{
		ServerPort:    envy.Get("SERVER_PORT", "3000"),
		DatabaseDSN:   envy.Get("DATABASE_DSN", "catter.db"),
		SessionSecret: envy.Get("SESSION_SECRET", "your-secret-key-please-change-in-production"),
		Environment:   envy.Get("ENVIRONMENT", "development"),
		LogLevel:      envy.Get("LOG_LEVEL", "debug"),
		CacheTTL:      300,
		CacheEnabled:  true,
	}
	return nil
}

// GetServerPort 获取服务器端口
func GetServerPort() string {
	return AppConfig.ServerPort
}

// GetDatabaseDSN 获取数据库连接字符串
func GetDatabaseDSN() string {
	return AppConfig.DatabaseDSN
}

// GetSessionSecret 获取会话密钥
func GetSessionSecret() string {
	return AppConfig.SessionSecret
}

// IsCacheEnabled 是否启用缓存
func IsCacheEnabled() bool {
	return AppConfig.CacheEnabled
}

// GetCacheTTL 获取缓存过期时间
func GetCacheTTL() int {
	return AppConfig.CacheTTL
}

// GetEnvironment 获取环境
func GetEnvironment() string {
	return AppConfig.Environment
}

// GetLogLevel 获取日志级别
func GetLogLevel() string {
	return AppConfig.LogLevel
}

// IsProduction 是否为生产环境
func IsProduction() bool {
	return AppConfig.Environment == "production"
}

// IsDevelopment 是否为开发环境
func IsDevelopment() bool {
	return AppConfig.Environment == "development"
}
