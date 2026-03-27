package models

import (
	"encoding/json"
	"fmt"
)

// Setting 系统设置模型
type Setting struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"uniqueIndex;size:100" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

// TableName 设置表名
func (Setting) TableName() string {
	return "settings"
}

// DefaultSettings 默认设置
var DefaultSettings = map[string]interface{}{
	"site_title":         "Catter CMS",
	"site_description":   "一个现代化的内容管理系统",
	"site_url":           "http://localhost:3001",
	"logo":               "",
	"favicon":            "",
	"default_theme":      "catter-blog",
	"active_theme":       "catter-blog",
	"posts_per_page":     10,
	"comment_moderation": true,
	"comment_approval":   true,
}

// InitSettings 初始化默认设置
func InitSettings() error {
	// 第一步：确保所有默认设置都存在
	for key, value := range DefaultSettings {
		var setting Setting
		result := DB.Where("key = ?", key).First(&setting)
		if result.Error != nil {
			// 设置不存在，创建它
			valueStr, _ := json.Marshal(value)
			setting = Setting{
				Key:   key,
				Value: string(valueStr),
			}
			if err := DB.Create(&setting).Error; err != nil {
				return fmt.Errorf("failed to create setting %s: %v", key, err)
			}
		}
	}

	// 第二步：检查并修复损坏的数据
	if err := validateAndFixSettings(); err != nil {
		// 只记录错误，不阻止启动
		fmt.Printf("WARNING: settings validation found issues: %v\n", err)
	}

	return nil
}

// validateAndFixSettings 检查并修复损坏的设置值
func validateAndFixSettings() error {
	var settings []Setting
	if err := DB.Find(&settings).Error; err != nil {
		return err
	}

	for _, setting := range settings {
		// 检查值是否是有效的 JSON
		if !json.Valid([]byte(setting.Value)) {
			fmt.Printf("WARNING: corrupted setting detected: %s = %s, resetting to default\n",
				setting.Key, truncateString(setting.Value, 50))

			// 如果有默认值，用默认值覆盖
			if defaultValue, exists := DefaultSettings[setting.Key]; exists {
				valueStr, _ := json.Marshal(defaultValue)
				DB.Model(&setting).Update("value", string(valueStr))
			} else {
				// 没有默认值，设为空字符串的JSON
				DB.Model(&setting).Update("value", `""`)
			}
		}
	}

	return nil
}

// truncateString 截断字符串用于日志显示
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// GetSetting 获取单个设置值
func GetSetting(key string) (string, error) {
	var setting Setting
	if err := DB.Where("key = ?", key).First(&setting).Error; err != nil {
		// 如果数据库中不存在，返回默认值
		if defaultValue, exists := DefaultSettings[key]; exists {
			valueStr, _ := json.Marshal(defaultValue)
			return string(valueStr), nil
		}
		return "", err
	}
	return setting.Value, nil
}

// GetSettingAsString 获取设置值为字符串
func GetSettingAsString(key string) string {
	value, err := GetSetting(key)
	if err != nil {
		if defaultValue, exists := DefaultSettings[key]; exists {
			if str, ok := defaultValue.(string); ok {
				return str
			}
		}
		return ""
	}
	// 尝试解析JSON字符串
	var result string
	if err := json.Unmarshal([]byte(value), &result); err == nil {
		return result
	}
	return value
}

// GetSettingAsBool 获取设置值为布尔值
func GetSettingAsBool(key string) bool {
	value, err := GetSetting(key)
	if err != nil {
		if defaultValue, exists := DefaultSettings[key]; exists {
			if b, ok := defaultValue.(bool); ok {
				return b
			}
		}
		return false
	}
	var result bool
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return false
	}
	return result
}

// GetSettingAsInt 获取设置值为整数
func GetSettingAsInt(key string) int {
	value, err := GetSetting(key)
	if err != nil {
		if defaultValue, exists := DefaultSettings[key]; exists {
			if i, ok := defaultValue.(int); ok {
				return i
			}
			if f, ok := defaultValue.(float64); ok {
				return int(f)
			}
		}
		return 0
	}
	var result float64
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return 0
	}
	return int(result)
}

// SetSetting 设置单个设置值
func SetSetting(key string, value interface{}) error {
	valueStr, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var setting Setting
	result := DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		// 创建新设置
		setting = Setting{
			Key:   key,
			Value: string(valueStr),
		}
		return DB.Create(&setting).Error
	}

	// 更新现有设置
	setting.Value = string(valueStr)
	return DB.Save(&setting).Error
}

// SetSettings 批量设置
func SetSettings(settings map[string]interface{}) error {
	for key, value := range settings {
		if err := SetSetting(key, value); err != nil {
			return err
		}
	}
	return nil
}

// GetAllSettings 获取所有设置
func GetAllSettings() (map[string]interface{}, error) {
	settings := make(map[string]interface{})

	// 先加载默认值
	for key, value := range DefaultSettings {
		settings[key] = value
	}

	// 从数据库加载覆盖
	var dbSettings []Setting
	if err := DB.Find(&dbSettings).Error; err != nil {
		return settings, nil // 返回默认值
	}

	for _, setting := range dbSettings {
		var value interface{}
		if err := json.Unmarshal([]byte(setting.Value), &value); err == nil {
			settings[setting.Key] = value
		} else {
			settings[setting.Key] = setting.Value
		}
	}

	return settings, nil
}

// GetActiveTheme 获取当前激活的主题
func GetActiveTheme() string {
	return GetSettingAsString("active_theme")
}

// SetActiveTheme 设置激活的主题
func SetActiveTheme(themeName string) error {
	return SetSetting("active_theme", themeName)
}
