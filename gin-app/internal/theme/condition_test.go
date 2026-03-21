package theme

import (
	"testing"
)

func TestConditionEvaluator_Evaluate(t *testing.T) {
	data := map[string]interface{}{
		"enabled": true,
		"disabled": false,
		"type": "advanced",
		"count": 10,
		"name": "test",
		"nested": map[string]interface{}{
			"value": 5,
			"enabled": true,
		},
	}

	tests := []struct {
		name      string
		condition string
		want      bool
		wantErr   bool
	}{
		// 基本存在性检查
		{"field exists", "enabled", true, false},
		{"field not exists", "nonexistent", false, false},
		{"nested field exists", "nested.enabled", true, false},
		{"nested value", "nested.value", true, false},

		// 布尔值比较
		{"equals true", "enabled == true", true, false},
		{"equals false", "disabled == true", false, false},
		{"not equals", "enabled != false", true, false},

		// 字符串比较
		{"string equals", "type == 'advanced'", true, false},
		{"string not equals", "type == 'basic'", false, false},
		{"string with double quotes", `name == "test"`, true, false},

		// 数字比较
		{"number equals", "count == 10", true, false},
		{"number greater", "count > 5", true, false},
		{"number less", "count < 20", true, false},
		{"number greater equals", "count >= 10", true, false},
		{"number less equals", "count <= 10", true, false},
		{"nested number compare", "nested.value > 3", true, false},

		// 逻辑与
		{"and both true", "enabled && type == 'advanced'", true, false},
		{"and one false", "enabled && type == 'basic'", false, false},
		{"and both false", "disabled && type == 'basic'", false, false},
		{"and multiple", "enabled && count > 5 && type == 'advanced'", true, false},

		// 逻辑或
		{"or both true", "enabled || type == 'advanced'", true, false},
		{"or one true", "disabled || type == 'advanced'", true, false},
		{"or both false", "disabled || type == 'basic'", false, false},

		// 混合逻辑
		{"and or mixed", "enabled && (type == 'basic' || type == 'advanced')", true, false},
		{"complex condition", "(enabled && count > 5) || (disabled && type == 'basic')", true, false},

		// 空条件
		{"empty condition", "", true, false},
		{"only whitespace", "   ", true, false},

		// 否定
		{"not true", "!enabled", false, false},
		{"not false", "!disabled", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewConditionEvaluator(data)
			got, err := evaluator.Evaluate(tt.condition)

			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v for condition: %s", got, tt.want, tt.condition)
			}
		})
	}
}

func TestShouldShowField(t *testing.T) {
	data := map[string]interface{}{
		"enabled": true,
		"type": "advanced",
	}

	tests := []struct {
		name     string
		field    ModuleField
		expected bool
	}{
		{
			name:     "no show_if condition",
			field:    ModuleField{Name: "test", ShowIf: ""},
			expected: true,
		},
		{
			name:     "condition true",
			field:    ModuleField{Name: "test", ShowIf: "enabled"},
			expected: true,
		},
		{
			name:     "condition false",
			field:    ModuleField{Name: "test", ShowIf: "!enabled"},
			expected: false,
		},
		{
			name:     "complex condition true",
			field:    ModuleField{Name: "test", ShowIf: "enabled && type == 'advanced'"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldShowField(tt.field, data)
			if result != tt.expected {
				t.Errorf("ShouldShowField() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestModuleManager(t *testing.T) {
	mm := NewModuleManager()

	// 测试注册模块
	module := &ModuleDefinition{
		ID:          "test-module",
		Name:        "Test Module",
		Description: "A test module",
		Enabled:     true,
		Order:       1,
	}

	err := mm.Register(module)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试重复注册
	err = mm.Register(module)
	if err == nil {
		t.Error("Register() should return error for duplicate module")
	}

	// 测试获取模块
	retrieved, exists := mm.Get("test-module")
	if !exists {
		t.Error("Get() should return existing module")
	}
	if retrieved.Name != "Test Module" {
		t.Errorf("Get() returned wrong module, got %s", retrieved.Name)
	}

	// 测试获取所有模块
	allModules := mm.GetAll()
	if len(allModules) != 1 {
		t.Errorf("GetAll() returned %d modules, want 1", len(allModules))
	}

	// 测试注销模块
	mm.Unregister("test-module")
	_, exists = mm.Get("test-module")
	if exists {
		t.Error("Unregister() should remove module")
	}
}

func TestRepeaterField_Validate(t *testing.T) {
	config := &RepeaterConfig{
		MinItems: 1,
		MaxItems: 3,
		Fields: []ModuleField{
			{Name: "title", Type: "text", Required: true},
			{Name: "count", Type: "number"},
		},
	}

	rf := NewRepeaterField(config)

	tests := []struct {
		name    string
		items   []interface{}
		wantErr bool
	}{
		{
			name: "valid items",
			items: []interface{}{
				map[string]interface{}{"title": "Item 1", "count": 10},
				map[string]interface{}{"title": "Item 2"},
			},
			wantErr: false,
		},
		{
			name:    "empty items - below min",
			items:   []interface{}{},
			wantErr: true,
		},
		{
			name: "too many items",
			items: []interface{}{
				map[string]interface{}{"title": "1"},
				map[string]interface{}{"title": "2"},
				map[string]interface{}{"title": "3"},
				map[string]interface{}{"title": "4"},
			},
			wantErr: true,
		},
		{
			name: "missing required field",
			items: []interface{}{
				map[string]interface{}{"count": 10},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rf.Validate(tt.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeepCopyMap(t *testing.T) {
	original := map[string]interface{}{
		"name": "test",
		"nested": map[string]interface{}{
			"value": 42,
			"items": []interface{}{1, 2, 3},
		},
		"list": []interface{}{
			map[string]interface{}{"id": 1},
			map[string]interface{}{"id": 2},
		},
	}

	copy := DeepCopyMap(original)

	// 修改副本不应影响原
	copy["name"] = "modified"
	copy["nested"].(map[string]interface{})["value"] = 999
	copy["list"].([]interface{})[0].(map[string]interface{})["id"] = 999

	// 验证原数据未被修改
	if original["name"] != "test" {
		t.Error("DeepCopyMap: original name was modified")
	}
	if original["nested"].(map[string]interface{})["value"] != 42 {
		t.Error("DeepCopyMap: original nested.value was modified")
	}
	if original["list"].([]interface{})[0].(map[string]interface{})["id"] != 1 {
		t.Error("DeepCopyMap: original list[0].id was modified")
	}
}
