// Package theme 提供模板函数扩展
// 支持条件渲染、模块引用等功能
package theme

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// TemplateFuncs 返回主题相关的模板函数
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// 条件显示
		"if":    tmplIf,
		"and":   tmplAnd,
		"or":    tmplOr,
		"not":   tmplNot,
		"eq":    tmplEq,
		"ne":    tmplNe,
		"gt":    tmplGt,
		"lt":    tmplLt,
		"gte":   tmplGte,
		"lte":   tmplLte,

		// 数据处理
		"get":       tmplGet,
		"has":       tmplHas,
		"len":       tmplLen,
		"default":   tmplDefault,
		"coalesce":  tmplCoalesce,

		// 字符串处理
		"upper":     strings.ToUpper,
		"lower":     strings.ToLower,
		"title":     strings.Title,
		"trim":      strings.TrimSpace,
		"replace":   strings.ReplaceAll,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,

		// 模块渲染
		"renderModule": tmplRenderModule,
		"hasModule":    tmplHasModule,

		// 调试
		"dump": tmplDump,
		"type": tmplType,
	}
}

// tmplIf 返回条件的真假值
func tmplIf(condition interface{}) bool {
	if condition == nil {
		return false
	}
	switch v := condition.(type) {
	case bool:
		return v
	case string:
		return v != "" && v != "false" && v != "0"
	case int, int8, int16, int32, int64:
		return v != 0
	case float32, float64:
		return v != 0
	case []interface{}:
		return len(v) > 0
	default:
		return true
	}
}

// tmplAnd 逻辑与
func tmplAnd(args ...interface{}) bool {
	for _, arg := range args {
		if !tmplIf(arg) {
			return false
		}
	}
	return len(args) > 0
}

// tmplOr 逻辑或
func tmplOr(args ...interface{}) bool {
	for _, arg := range args {
		if tmplIf(arg) {
			return true
		}
	}
	return false
}

// tmplNot 逻辑非
func tmplNot(arg interface{}) bool {
	return !tmplIf(arg)
}

// tmplEq 等于
func tmplEq(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// tmplNe 不等于
func tmplNe(a, b interface{}) bool {
	return !tmplEq(a, b)
}

// tmplGt 大于
func tmplGt(a, b interface{}) bool {
	return compareNumbers(a, b) > 0
}

// tmplLt 小于
func tmplLt(a, b interface{}) bool {
	return compareNumbers(a, b) < 0
}

// tmplGte 大于等于
func tmplGte(a, b interface{}) bool {
	return compareNumbers(a, b) >= 0
}

// tmplLte 小于等于
func tmplLte(a, b interface{}) bool {
	return compareNumbers(a, b) <= 0
}

// compareNumbers 比较两个数字，返回 -1, 0, 1
func compareNumbers(a, b interface{}) int {
	aFloat := toFloat64(a)
	bFloat := toFloat64(b)

	if aFloat < bFloat {
		return -1
	}
	if aFloat > bFloat {
		return 1
	}
	return 0
}

// toFloat64 转换为float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

// tmplGet 从map中获取值，支持点号路径
func tmplGet(data interface{}, path string) interface{} {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	keys := strings.Split(path, ".")
	current := dataMap

	for i, key := range keys {
		if i == len(keys)-1 {
			return current[key]
		}

		next, ok := current[key].(map[string]interface{})
		if !ok {
			return nil
		}
		current = next
	}

	return nil
}

// tmplHas 检查map中是否存在路径
func tmplHas(data interface{}, path string) bool {
	return tmplGet(data, path) != nil
}

// tmplLen 获取长度
func tmplLen(v interface{}) int {
	switch val := v.(type) {
	case string:
		return len(val)
	case []interface{}:
		return len(val)
	case map[string]interface{}:
		return len(val)
	default:
		return 0
	}
}

// tmplDefault 如果值为空则返回默认值
func tmplDefault(value, defaultValue interface{}) interface{} {
	if value == nil || value == "" {
		return defaultValue
	}
	return value
}

// tmplCoalesce 返回第一个非空值
func tmplCoalesce(values ...interface{}) interface{} {
	for _, v := range values {
		if v != nil && v != "" {
			return v
		}
	}
	return nil
}

// tmplRenderModule 渲染模块
func tmplRenderModule(moduleName string, data interface{}) template.HTML {
	// 这里只是一个占位符，实际渲染逻辑需要结合模板引擎
	// 返回安全HTML
	return template.HTML(fmt.Sprintf("<!-- module: %s -->", moduleName))
}

// tmplHasModule 检查模块是否存在
func tmplHasModule(moduleName string) bool {
	manager := GetModuleManager()
	_, exists := manager.Get(moduleName)
	return exists
}

// tmplDump 调试输出
func tmplDump(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}

// tmplType 返回值的类型
func tmplType(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return reflect.TypeOf(v).String()
}

// EvaluateCondition 便捷函数：评估条件表达式
func EvaluateCondition(condition string, data map[string]interface{}) (bool, error) {
	evaluator := NewConditionEvaluator(data)
	return evaluator.Evaluate(condition)
}

// ShouldShowField 检查字段是否应该显示（基于show_if条件）
func ShouldShowField(field ModuleField, data map[string]interface{}) bool {
	if field.ShowIf == "" {
		return true
	}

	result, err := EvaluateCondition(field.ShowIf, data)
	if err != nil {
		// 如果条件解析失败，默认显示
		return true
	}

	return result
}
