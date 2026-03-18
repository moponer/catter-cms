package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// ThemeConfig 主题配置模型 - 用于持久化存储主题配置
type ThemeConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ThemeName string    `gorm:"uniqueIndex;size:100" json:"theme_name"`
	Config    string    `gorm:"type:text" json:"config"` // JSON格式的配置
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (ThemeConfig) TableName() string {
	return "theme_configs"
}

// GetThemeConfig 获取主题配置
func GetThemeConfig(themeName string) (map[string]interface{}, error) {
	var themeConfig ThemeConfig
	result := DB.Where("theme_name = ?", themeName).First(&themeConfig)

	if result.Error != nil {
		// 配置不存在，返回空配置
		return make(map[string]interface{}), nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(themeConfig.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal theme config: %v", err)
	}

	return config, nil
}

// SaveThemeConfig 保存主题配置
func SaveThemeConfig(themeName string, config map[string]interface{}) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal theme config: %v", err)
	}

	var themeConfig ThemeConfig
	result := DB.Where("theme_name = ?", themeName).First(&themeConfig)

	if result.Error != nil {
		// 创建新配置
		themeConfig = ThemeConfig{
			ThemeName: themeName,
			Config:    string(configJSON),
		}
		return DB.Create(&themeConfig).Error
	}

	// 更新现有配置
	themeConfig.Config = string(configJSON)
	return DB.Save(&themeConfig).Error
}

// DeleteThemeConfig 删除主题配置
func DeleteThemeConfig(themeName string) error {
	return DB.Where("theme_name = ?", themeName).Delete(&ThemeConfig{}).Error
}

// MergeThemeConfig 合并主题配置（将新配置合并到现有配置中）
func MergeThemeConfig(themeName string, newConfig map[string]interface{}) error {
	existingConfig, err := GetThemeConfig(themeName)
	if err != nil {
		return err
	}

	// 合并配置
	for key, value := range newConfig {
		existingConfig[key] = value
	}

	return SaveThemeConfig(themeName, existingConfig)
}
