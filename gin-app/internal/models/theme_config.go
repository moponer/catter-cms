package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
// 注意：现在统一返回扁平格式（dot-notation key），与 services/theme.go 保持一致
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

	// 直接返回扁平格式，与 services/theme.go 的 FlattenMap 输出保持一致
	return config, nil
}

// flattenToNested 将扁平结构（点号分隔的键）转换为嵌套结构
func flattenToNested(flat map[string]interface{}) map[string]interface{} {
	if flat == nil {
		return make(map[string]interface{})
	}

	nested := make(map[string]interface{})
	for key, value := range flat {
		// 检查键名是否包含点号
		if strings.Contains(key, ".") {
			// 按照点号分割键名
			parts := strings.Split(key, ".")
			current := nested

			// 遍历所有部分，除了最后一个
			for i := 0; i < len(parts)-1; i++ {
				part := parts[i]
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				if nextMap, ok := current[part].(map[string]interface{}); ok {
					current = nextMap
				} else {
					// 如果已存在但不是map，覆盖它
					newMap := make(map[string]interface{})
					current[part] = newMap
					current = newMap
				}
			}

			// 设置最后一个部分的值
			current[parts[len(parts)-1]] = value
		} else {
			// 没有点号，直接设置
			nested[key] = value
		}
	}
	return nested
}

// DeepMerge 将 src 深度合并到 dst，src 中的值会覆盖 dst 中的同名值
// 对于 map 类型，会递归合并而不是整体替换
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = make(map[string]interface{})
	}
	if src == nil {
		return dst
	}

	for key, srcVal := range src {
		if dstVal, exists := dst[key]; exists {
			// 如果两边都是 map，递归合并
			srcMap, srcOk := srcVal.(map[string]interface{})
			dstMap, dstOk := dstVal.(map[string]interface{})
			if srcOk && dstOk {
				dst[key] = DeepMerge(dstMap, srcMap)
				continue
			}
		}
		// 否则直接覆盖
		dst[key] = srcVal
	}
	return dst
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
		if err := DB.Create(&themeConfig).Error; err != nil {
			log.Printf("[ERROR] Failed to create theme config for %s: %v", themeName, err)
			return err
		}
		log.Printf("[INFO] Created theme config for %s", themeName)
		return nil
	}

	// 更新现有配置
	themeConfig.Config = string(configJSON)
	if err = DB.Save(&themeConfig).Error; err != nil {
		log.Printf("[ERROR] Failed to update theme config for %s: %v", themeName, err)
		return err
	}
	log.Printf("[INFO] Updated theme config for %s", themeName)
	return nil
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
