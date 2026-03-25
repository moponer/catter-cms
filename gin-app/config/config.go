package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
	Security SecurityConfig `yaml:"security"`
	Env      string         `yaml:"env"` // 环境: development, production
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port      string `yaml:"port"`
	Host      string `yaml:"host"`
	Mode      string `yaml:"mode"`
	ThemePath string `yaml:"theme_path"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Database string `yaml:"database"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"` // 小时
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize     int64  `yaml:"max_size"`
	Path        string `yaml:"path"`
	AllowedExt  string `yaml:"allowed_ext"`
	StorageType string `yaml:"storage_type"` // local, oss, s3
	OSSConfig   OSSConfig `yaml:"oss_config,omitempty"`
}

// OSSConfig OSS配置
type OSSConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Bucket          string `yaml:"bucket"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	AllowOrigins []string `yaml:"allow_origins"`
	CSP          string   `yaml:"csp"`
	EnableCSP    bool     `yaml:"enable_csp"`
}

var cfg *Config

// Load 加载配置
func Load(path string) (*Config, error) {
	// 获取环境变量中的配置路径
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		path = envPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Config file not found at %s, using default config", path)
		cfg = defaultConfig()
		return cfg, nil
	}

	loadedCfg := &Config{}
	if err := yaml.Unmarshal(data, loadedCfg); err != nil {
		log.Printf("Failed to parse config file: %v, using default config", err)
		cfg = defaultConfig()
		return cfg, nil
	}

	// 从环境变量覆盖（支持开发/生产切换）
	applyEnvOverrides(loadedCfg)

	// 验证配置
	if err := validateConfig(loadedCfg); err != nil {
		log.Printf("Config validation failed: %v", err)
		return nil, err
	}

	cfg = loadedCfg
	return cfg, nil
}

// Get 获取配置
func Get() *Config {
	if cfg == nil {
		cfg = defaultConfig()
	}
	return cfg
}

// applyEnvOverrides 从环境变量应用配置覆盖
func applyEnvOverrides(c *Config) {
	// 端口
	if port := os.Getenv("PORT"); port != "" {
		c.Server.Port = port
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		c.Server.Port = port
	}

	// 环境模式
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		c.Server.Mode = mode
	}
	if env := os.Getenv("APP_ENV"); env != "" {
		c.Env = env
	}

	// JWT 密钥（生产环境必须从环境变量读取）
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		c.JWT.Secret = secret
	}

	// 数据库路径
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		c.Database.Database = dbPath
	}

	// 上传路径
	if uploadPath := os.Getenv("UPLOAD_PATH"); uploadPath != "" {
		c.Upload.Path = uploadPath
	}
}

// validateConfig 验证配置
func validateConfig(c *Config) error {
	// 检查是否为生产环境且使用默认 JWT 密钥
	if c.Env == "production" || c.Server.Mode == "release" {
		if c.JWT.Secret == "your-secret-key-change-in-production" ||
		   c.JWT.Secret == "" ||
		   len(c.JWT.Secret) < 32 {
			return fmt.Errorf("production mode requires a strong JWT secret (at least 32 characters)")
		}
	}

	// 检查端口是否有效
	if c.Server.Port == "" {
		c.Server.Port = "3000"
	}

	// 检查环境标识
	if c.Env == "" {
		if c.Server.Mode == "release" {
			c.Env = "production"
		} else {
			c.Env = "development"
		}
	}

	return nil
}

// IsProduction 检查是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Env == "production" || c.Server.Mode == "release"
}

// IsDevelopment 检查是否为开发环境
func (c *Config) IsDevelopment() bool {
	return !c.IsProduction()
}

// GetAdminDistPath 获取管理后台静态文件目录
func (c *Config) GetAdminDistPath() string {
	// 支持从环境变量或配置指定路径
	if path := os.Getenv("ADMIN_DIST_PATH"); path != "" {
		return path
	}

	// 尝试多个常见路径
	paths := []string{
		"../admin/dist",
		"./admin/dist",
		filepath.Join(GetAppRoot(), "admin/dist"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// 默认返回第一个路径（实际使用时会被检测）
	return paths[0]
}

// GetThemePath 获取主题目录
func (c *Config) GetThemePath() string {
	if c.Server.ThemePath != "" {
		return c.Server.ThemePath
	}

	paths := []string{
		"./templates/themes",
		"../templates/themes",
		filepath.Join(GetAppRoot(), "templates/themes"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return "./templates/themes"
}

// GetUploadsPath 获取上传文件目录
func (c *Config) GetUploadsPath() string {
	if c.Upload.Path != "" {
		return c.Upload.Path
	}
	return "./uploads"
}

// GetAppRoot 获取应用根目录
func GetAppRoot() string {
	// 获取当前执行文件的目录
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

// GenerateSecureSecret 生成安全的随机密钥
func GenerateSecureSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetFullPort 返回完整地址格式 (如 "0.0.0.0:3000")
func (c *Config) GetFullPort() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// GetPublicURL 返回公开访问的完整 URL
func (c *Config) GetPublicURL() string {
	return fmt.Sprintf("http://%s:%s", c.Server.Host, c.Server.Port)
}

// LoadFromEnv 从环境变量加载最小配置（用于容器化部署）
func LoadFromEnv() *Config {
	c := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "3000"),
			Host: getEnvOrDefault("SERVER_HOST", "0.0.0.0"),
			Mode: getEnvOrDefault("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Driver:   getEnvOrDefault("DB_DRIVER", "sqlite"),
			Database: getEnvOrDefault("DB_DATABASE", "catter.db"),
		},
		JWT: JWTConfig{
			Secret: getEnvOrDefault("JWT_SECRET", ""),
			Expire: 24,
		},
		Upload: UploadConfig{
			MaxSize:    10 * 1024 * 1024,
			Path:       getEnvOrDefault("UPLOAD_PATH", "uploads"),
			AllowedExt: ".jpg,.jpeg,.png,.gif,.webp",
		},
		Env: getEnvOrDefault("APP_ENV", "production"),
	}

	// 如果 JWT_SECRET 太短或为空，生成一个
	if len(c.JWT.Secret) < 32 {
		secret, err := GenerateSecureSecret(32)
		if err != nil {
			log.Printf("Failed to generate secure secret: %v, using fallback", err)
			c.JWT.Secret = "insecure-default-change-me"
		} else {
			c.JWT.Secret = secret
		}
	}

	cfg = c
	return c
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// MergeConfig 合并配置（用于覆盖）
func MergeConfig(base, override *Config) *Config {
	if override == nil {
		return base
	}
	// 这里可以实现深度合并逻辑
	result := *base
	// ... 简化实现，根据需要扩展
	return &result
}

// String 返回配置的字符串表示（用于日志，不包含敏感信息）
func (c *Config) String() string {
	return fmt.Sprintf("Config{Env:%s, Port:%s, Host:%s, DB:%s, Mode:%s}",
		c.Env, c.Server.Port, c.Server.Host, c.Database.Database, c.Server.Mode)
}

// defaultConfig 默认配置
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "3000",
			Host: "0.0.0.0",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Driver:   "sqlite",
			Database: "catter.db",
		},
		JWT: JWTConfig{
			Secret: "dev-only-secret-do-not-use-in-production",
			Expire: 24,
		},
		Upload: UploadConfig{
			MaxSize:    10 * 1024 * 1024, // 10MB
			Path:       "uploads",
			AllowedExt: ".jpg,.jpeg,.png,.gif,.webp",
		},
		Security: SecurityConfig{
			AllowOrigins: []string{"*"},
			EnableCSP:    false,
		},
		Env: "development",
	}
}

// InitConfig 初始化配置（支持多种加载方式）
func InitConfig(configPaths ...string) (*Config, error) {
	// 1. 首先尝试环境变量指定的路径
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		if cfg, err := Load(envPath); err == nil {
			log.Printf("Config loaded from environment CONFIG_PATH: %s", envPath)
			return cfg, nil
		}
	}

	// 2. 尝试指定的配置路径
	for _, path := range configPaths {
		if cfg, err := Load(path); err == nil {
			log.Printf("Config loaded from: %s", path)
			return cfg, nil
		}
	}

	// 3. 尝试常见的配置路径
	commonPaths := []string{
		"./config.yaml",
		"../config.yaml",
		"../../config.yaml",
		filepath.Join(GetAppRoot(), "config.yaml"),
		"/etc/catter-cms/config.yaml",
	}

	for _, path := range commonPaths {
		if cfg, err := Load(path); err == nil {
			log.Printf("Config loaded from common path: %s", path)
			return cfg, nil
		}
	}

	// 4. 如果是生产环境，尝试从环境变量加载最小配置
	if os.Getenv("APP_ENV") == "production" || os.Getenv("GIN_MODE") == "release" {
		log.Println("No config file found in production mode, loading from environment variables")
		return LoadFromEnv(), nil
	}

	// 5. 使用默认配置
	log.Println("Using default configuration")
	return defaultConfig(), nil
}

// Save 保存配置到文件
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// UpdateFromCLI 从命令行参数更新配置
func (c *Config) UpdateFromCLI(port, host, env, dbPath string) {
	if port != "" {
		c.Server.Port = port
	}
	if host != "" {
		c.Server.Host = host
	}
	if env != "" {
		c.Env = env
		if env == "production" {
			c.Server.Mode = "release"
		} else {
			c.Server.Mode = "debug"
		}
	}
	if dbPath != "" {
		c.Database.Database = dbPath
	}
}

// ValidateSecurity 检查安全配置
func (c *Config) ValidateSecurity() []string {
	var warnings []string

	// 检查 CORS 配置
	if len(c.Security.AllowOrigins) == 0 {
		warnings = append(warnings, "No CORS origins configured, using defaults")
	}

	// 检查 JWT 密钥强度
	if len(c.JWT.Secret) < 32 {
		warnings = append(warnings, "JWT secret is too short (< 32 characters)")
	}

	if strings.Contains(c.JWT.Secret, "secret") || strings.Contains(c.JWT.Secret, "password") {
		warnings = append(warnings, "JWT secret appears to contain common words, consider using a stronger secret")
	}

	// 检查是否启用了 CSP
	if !c.Security.EnableCSP && c.IsProduction() {
		warnings = append(warnings, "Content Security Policy (CSP) is disabled in production")
	}

	return warnings
}
