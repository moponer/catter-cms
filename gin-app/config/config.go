package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
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
	MaxSize    int64  `yaml:"max_size"`
	Path       string `yaml:"path"`
	AllowedExt string `yaml:"allowed_ext"`
}

var cfg *Config

// Load 加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		log.Println("Config file not found, using default config")
		cfg = defaultConfig()
		return cfg, nil
	}

	loadedCfg := &Config{}
	if err := yaml.Unmarshal(data, loadedCfg); err != nil {
		log.Printf("Failed to parse config file: %v, using default config", err)
		cfg = defaultConfig()
		return cfg, nil
	}

	// 从环境变量覆盖
	if port := os.Getenv("PORT"); port != "" {
		loadedCfg.Server.Port = port
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

// defaultConfig 默认配置
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "3001",
			Host: "0.0.0.0",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Driver:   "sqlite",
			Database: "catter.db",
		},
		JWT: JWTConfig{
			Secret: "your-secret-key-change-in-production",
			Expire: 24,
		},
		Upload: UploadConfig{
			MaxSize:    10 * 1024 * 1024, // 10MB
			Path:       "uploads",
			AllowedExt: ".jpg,.jpeg,.png,.gif,.webp",
		},
	}
}
