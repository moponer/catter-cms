// Package theme 条件显示引擎
// 支持复杂的条件表达式评估，用于控制字段和模块的显示
package theme

import (
	"fmt"
	"strconv"
	"strings"
)

// ConditionEvaluator 条件表达式评估器
type ConditionEvaluator struct {
	data map[string]interface{}
}

// NewConditionEvaluator 创建新的条件评估器
func NewConditionEvaluator(data map[string]interface{}) *ConditionEvaluator {
	return &ConditionEvaluator{
		data: data,
	}
}

// Evaluate 评估条件表达式
// 支持的语法:
//   - "field.enabled" (字段存在且为真)
//   - "field.value == 'string'" (等于)
//   - "field.value != 'string'" (不等于)
//   - "field.num > 10" (大于)
//   - "field.enabled && field.type == 'x'" (逻辑与)
//   - "field.enabled || field.other == 'y'" (逻辑或)
//   - "(a || b) && c" (括号分组)
func (ce *ConditionEvaluator) Evaluate(condition string) (bool, error) {
	if strings.TrimSpace(condition) == "" {
		return true, nil
	}

	// 解析并评估表达式
	result, err := ce.parseExpression(condition)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate condition '%s': %w", condition, err)
	}

	return result, nil
}

// parseExpression 解析表达式（处理 || ）
func (ce *ConditionEvaluator) parseExpression(expr string) (bool, error) {
	expr = strings.TrimSpace(expr)

	// 先处理括号
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		// 检查是否是外层括号
		depth := 0
		isOuter := true
		for i, ch := range expr {
			if i == 0 {
				continue
			}
			if ch == '(' {
				depth++
			} else if ch == ')' {
				if depth == 0 && i < len(expr)-1 {
					isOuter = false
					break
				}
				depth--
			}
		}
		if isOuter {
			return ce.parseExpression(expr[1 : len(expr)-1])
		}
	}

	// 处理 || (OR) - 优先级最低
	orParts := ce.splitByOperator(expr, "||")
	if len(orParts) > 1 {
		for _, part := range orParts {
			result, err := ce.parseExpression(strings.TrimSpace(part))
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil
	}

	// 处理 && (AND)
	andParts := ce.splitByOperator(expr, "&&")
	if len(andParts) > 1 {
		for _, part := range andParts {
			result, err := ce.parseExpression(strings.TrimSpace(part))
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil
	}

	// 单个条件
	return ce.evaluateComparison(expr)
}

// splitByOperator 按操作符分割表达式，考虑括号
func (ce *ConditionEvaluator) splitByOperator(expr, op string) []string {
	var parts []string
	depth := 0
	lastIndex := 0

	for i := 0; i < len(expr); i++ {
		ch := expr[i]

		if ch == '(' || ch == '[' || ch == '{' {
			depth++
		} else if ch == ')' || ch == ']' || ch == '}' {
			depth--
		}

		if depth == 0 && i+len(op) <= len(expr) && expr[i:i+len(op)] == op {
			// 确保不是更大操作符的一部分
			if (op == "||" || op == "&&") && i > 0 && (expr[i-1] == '|' || expr[i-1] == '&') {
				continue
			}
			parts = append(parts, expr[lastIndex:i])
			lastIndex = i + len(op)
			i += len(op) - 1
		}
	}

	parts = append(parts, expr[lastIndex:])
	return parts
}

// evaluateComparison 评估比较表达式
func (ce *ConditionEvaluator) evaluateComparison(expr string) (bool, error) {
	expr = strings.TrimSpace(expr)

	// 定义所有支持的操作符（按长度降序，避免匹配错误）
	operators := []string{
		"!=", "==", ">=", "<=", ">", "<",
	}

	var op string
	var parts []string

	for _, candidate := range operators {
		if idx := strings.Index(expr, candidate); idx >= 0 {
			op = candidate
			parts = []string{
				strings.TrimSpace(expr[:idx]),
				strings.TrimSpace(expr[idx+len(candidate):]),
			}
			break
		}
	}

	// 如果没有找到操作符，只是字段存在性检查
	if op == "" {
		fieldPath := expr
		value := ce.getNestedValue(fieldPath)
		return ce.isTruthy(value), nil
	}

	if len(parts) != 2 {
		return false, fmt.Errorf("invalid comparison expression: %s", expr)
	}

	leftPath := parts[0]
	rightValue := parts[1]

	// 获取左边字段值
	leftValue := ce.getNestedValue(leftPath)

	// 解析右边值
	right := ce.parseValue(rightValue)

	// 执行比较
	return ce.compareValues(leftValue, op, right), nil
}

// getNestedValue 从嵌套map中获取值，支持点号路径
func (ce *ConditionEvaluator) getNestedValue(path string) interface{} {
	keys := strings.Split(path, ".")
	current := ce.data

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

// parseValue 解析字面量值
func (ce *ConditionEvaluator) parseValue(valueStr string) interface{} {
	valueStr = strings.TrimSpace(valueStr)

	// 字符串字面量
	if (strings.HasPrefix(valueStr, "'") && strings.HasSuffix(valueStr, "'")) ||
		(strings.HasPrefix(valueStr, "\"") && strings.HasSuffix(valueStr, "\"")) {
		return valueStr[1 : len(valueStr)-1]
	}

	// 布尔值
	if valueStr == "true" {
		return true
	}
	if valueStr == "false" {
		return false
	}

	// null/nil
	if valueStr == "null" || valueStr == "nil" {
		return nil
	}

	// 数字
	if num, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return num
	}

	// 作为字段路径处理
	return ce.getNestedValue(valueStr)
}

// isTruthy 判断值是否为真
func (ce *ConditionEvaluator) isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != "" && val != "false" && val != "0"
	case int, int8, int16, int32, int64:
		return val != 0
	case float32, float64:
		return val != 0
	case []interface{}:
		return len(val) > 0
	default:
		return true
	}
}

// compareValues 比较两个值
func (ce *ConditionEvaluator) compareValues(left interface{}, op string, right interface{}) bool {
	// 处理nil
	if left == nil && right == nil {
		return op == "==" || op == ">=" || op == "<="
	}
	if left == nil || right == nil {
		return op == "!="
	}

	// 字符串比较
	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)

	if leftIsStr && rightIsStr {
		switch op {
		case "==":
			return leftStr == rightStr
		case "!=":
			return leftStr != rightStr
		case ">":
			return leftStr > rightStr
		case "<":
			return leftStr < rightStr
		case ">=":
			return leftStr >= rightStr
		case "<=":
			return leftStr <= rightStr
		}
	}

	// 数字比较
	leftNum := ce.toFloat64(left)
	rightNum := ce.toFloat64(right)

	switch op {
	case "==":
		return leftNum == rightNum
	case "!=":
		return leftNum != rightNum
	case ">":
		return leftNum > rightNum
	case "<":
		return leftNum < rightNum
	case ">=":
		return leftNum >= rightNum
	case "<=":
		return leftNum <= rightNum
	}

	return false
}

// toFloat64 将值转换为float64
func (ce *ConditionEvaluator) toFloat64(v interface{}) float64 {
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

// Global module manager instance
var (
	globalModuleManager     *ModuleManager
	globalModuleManagerOnce sync.Once
)

// GetModuleManager 获取全局模块管理器
func GetModuleManager() *ModuleManager {
	globalModuleManagerOnce.Do(func() {
		globalModuleManager = NewModuleManager()
	})
	return globalModuleManager
}
