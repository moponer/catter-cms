// Package theme Repeater字段处理
// 支持动态列表数据的添加、删除、排序和验证
package theme

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// RepeaterField Repeater字段处理器
type RepeaterField struct {
	Config *RepeaterConfig
}

// NewRepeaterField 创建新的Repeater字段处理器
func NewRepeaterField(config *RepeaterConfig) *RepeaterField {
	if config == nil {
		config = &RepeaterConfig{
			MinItems: 0,
			MaxItems: 0, // 0表示无限制
			Fields:   []ModuleField{},
		}
	}
	return &RepeaterField{Config: config}
}

// Validate 验证Repeater数据
func (rf *RepeaterField) Validate(items []interface{}) error {
	// 检查最小项目数
	if rf.Config.MinItems > 0 && len(items) < rf.Config.MinItems {
		return fmt.Errorf("至少需要 %d 个项目，当前有 %d 个", rf.Config.MinItems, len(items))
	}

	// 检查最大项目数
	if rf.Config.MaxItems > 0 && len(items) > rf.Config.MaxItems {
		return fmt.Errorf("最多允许 %d 个项目，当前有 %d 个", rf.Config.MaxItems, len(items))
	}

	// 验证每个项目
	for i, item := range items {
		if err := rf.validateItem(i, item); err != nil {
			return fmt.Errorf("项目 %d 验证失败: %w", i+1, err)
		}
	}

	return nil
}

// validateItem 验证单个项目
func (rf *RepeaterField) validateItem(index int, item interface{}) error {
	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("项目必须是对象类型")
	}

	for _, field := range rf.Config.Fields {
		value, exists := itemMap[field.Name]

		// 检查必填字段
		if field.Required {
			if !exists || rf.isEmpty(value) {
				return fmt.Errorf("字段 '%s' 是必填的", field.Label)
			}
		}

		// 验证字段类型
		if exists && !rf.isEmpty(value) {
			if err := rf.validateFieldType(field, value); err != nil {
				return fmt.Errorf("字段 '%s' %w", field.Label, err)
			}
		}
	}

	return nil
}

// validateFieldType 验证字段类型
func (rf *RepeaterField) validateFieldType(field ModuleField, value interface{}) error {
	switch field.Type {
	case "text", "textarea", "code":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("必须是字符串类型")
		}
	case "number":
		switch v := value.(type) {
		case float64:
			if field.Min != 0 && v < field.Min {
				return fmt.Errorf("不能小于 %v", field.Min)
			}
			if field.Max != 0 && v > field.Max {
				return fmt.Errorf("不能大于 %v", field.Max)
			}
		case int, int64:
			// OK
		default:
			return fmt.Errorf("必须是数字类型")
		}
	case "switch":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("必须是布尔类型")
		}
	case "select":
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("必须是字符串类型")
		}
		// 验证选项是否有效
		valid := false
		for _, opt := range field.Options {
			if opt.Value == strValue {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("无效的选择值")
		}
	}

	return nil
}

// isEmpty 检查值是否为空
func (rf *RepeaterField) isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}

// ProcessRepeaterFields 处理数据中的Repeater字段
// 将JSON字符串解析为对象数组
func ProcessRepeaterFields(data map[string]interface{}, fields []ModuleField) error {
	for _, field := range fields {
		if field.Type != "repeater" || field.Repeater == nil {
			continue
		}

		value, exists := data[field.Name]
		if !exists {
			continue
		}

		// 如果已经是数组，跳过
		if _, ok := value.([]interface{}); ok {
			continue
		}

		// 尝试解析JSON字符串
		if strValue, ok := value.(string); ok {
			var items []interface{}
			if err := json.Unmarshal([]byte(strValue), &items); err != nil {
				return fmt.Errorf("failed to parse repeater field %s: %w", field.Name, err)
			}
			data[field.Name] = items
		}
	}

	return nil
}

// ConvertRepeaterToJSON 将Repeater字段转换为JSON字符串
// 用于存储到数据库
func ConvertRepeaterToJSON(data map[string]interface{}, fields []ModuleField) error {
	for _, field := range fields {
		if field.Type != "repeater" || field.Repeater == nil {
			continue
		}

		value, exists := data[field.Name]
		if !exists {
			continue
		}

		// 如果已经是字符串，跳过
		if _, ok := value.(string); ok {
			continue
		}

		// 转换为JSON字符串
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to convert repeater field %s to JSON: %w", field.Name, err)
		}
		data[field.Name] = string(jsonBytes)
	}

	return nil
}

// DeepCopyMap 深度拷贝map
func DeepCopyMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}

	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			result[k] = DeepCopyMap(val)
		case []interface{}:
			result[k] = deepCopySlice(val)
		default:
			result[k] = val
		}
	}
	return result
}

// deepCopySlice 深度拷贝切片
func deepCopySlice(s []interface{}) []interface{} {
	if s == nil {
		return nil
	}

	result := make([]interface{}, len(s))
	for i, v := range s {
		switch val := v.(type) {
		case map[string]interface{}:
			result[i] = DeepCopyMap(val)
		case []interface{}:
			result[i] = deepCopySlice(val)
		default:
			result[i] = val
		}
	}
	return result
}
