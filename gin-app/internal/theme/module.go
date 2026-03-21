// Package theme 提供主题系统的核心功能
// 包括模块管理、Repeater字段支持和条件显示
package theme

import (
	"fmt"
	"sync"
)

// ModuleManager 管理主题模块的注册和配置
type ModuleManager struct {
	modules map[string]*ModuleDefinition
	mu      sync.RWMutex
}

// ModuleDefinition 定义一个可配置的模块
type ModuleDefinition struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Category    string                 `json:"category" yaml:"category"`
	Icon        string                 `json:"icon" yaml:"icon"`
	Enabled     bool                   `json:"enabled" yaml:"enabled"`
	Order       int                    `json:"order" yaml:"order"`
	Config      map[string]interface{} `json:"config" yaml:"config"`
	Fields      []ModuleField          `json:"fields" yaml:"fields"`
}

// ModuleField 定义模块配置字段
type ModuleField struct {
	Name         string      `json:"name" yaml:"name"`
	Label        string      `json:"label" yaml:"label"`
	Type         string      `json:"type" yaml:"type"`
	Description  string      `json:"description" yaml:"description"`
	Default      interface{} `json:"default" yaml:"default"`
	Required     bool        `json:"required" yaml:"required"`
	ShowIf       string      `json:"show_if" yaml:"show_if"`
	Placeholder  string      `json:"placeholder" yaml:"placeholder"`
	Options      []Option    `json:"options,omitempty" yaml:"options,omitempty"`
	Repeater     *RepeaterConfig `json:"repeater,omitempty" yaml:"repeater,omitempty"`
}

// Option 定义下拉选项
type Option struct {
	Label string      `json:"label" yaml:"label"`
	Value interface{} `json:"value" yaml:"value"`
}

// RepeaterConfig Repeater字段配置
type RepeaterConfig struct {
	MinItems    int           `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems    int           `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	AddLabel    string        `json:"addLabel,omitempty" yaml:"addLabel,omitempty"`
	ItemLabel   string        `json:"itemLabel,omitempty" yaml:"itemLabel,omitempty"`
	Collapsible bool          `json:"collapsible,omitempty" yaml:"collapsible,omitempty"`
	Fields      []ModuleField `json:"fields" yaml:"fields"`
}

// PageTemplate 定义页面模板
type PageTemplate struct {
	Name        string `json:"name" yaml:"name"`
	File        string `json:"file" yaml:"file"`
	Description string `json:"description" yaml:"description"`
	Thumbnail   string `json:"thumbnail" yaml:"thumbnail"`
	Category    string `json:"category" yaml:"category"`
}

// NewModuleManager 创建新的模块管理器
func NewModuleManager() *ModuleManager {
	return &ModuleManager{
		modules: make(map[string]*ModuleDefinition),
	}
}

// Register 注册一个模块
func (mm *ModuleManager) Register(module *ModuleDefinition) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if module.ID == "" {
		return fmt.Errorf("module ID cannot be empty")
	}

	if _, exists := mm.modules[module.ID]; exists {
		return fmt.Errorf("module %s already registered", module.ID)
	}

	mm.modules[module.ID] = module
	return nil
}

// Unregister 注销一个模块
func (mm *ModuleManager) Unregister(moduleID string) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	delete(mm.modules, moduleID)
}

// Get 获取模块定义
func (mm *ModuleManager) Get(moduleID string) (*ModuleDefinition, bool) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	m, ok := mm.modules[moduleID]
	return m, ok
}

// GetAll 获取所有模块
func (mm *ModuleManager) GetAll() []*ModuleDefinition {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	modules := make([]*ModuleDefinition, 0, len(mm.modules))
	for _, m := range mm.modules {
		modules = append(modules, m)
	}
	return modules
}

// GetEnabled 获取启用的模块（按Order排序）
func (mm *ModuleManager) GetEnabled() []*ModuleDefinition {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	modules := make([]*ModuleDefinition, 0)
	for _, m := range mm.modules {
		if m.Enabled {
			modules = append(modules, m)
		}
	}

	// 按Order排序
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Order < modules[j].Order
	})

	return modules
}

// ValidateConfig 验证模块配置
func (m *ModuleDefinition) ValidateConfig(config map[string]interface{}) error {
	for _, field := range m.Fields {
		value, exists := config[field.Name]

		// 检查必填字段
		if field.Required && (!exists || value == nil || value == "") {
			return fmt.Errorf("field %s is required", field.Name)
		}

		// 验证 Repeater 字段
		if field.Type == "repeater" && field.Repeater != nil {
			if items, ok := value.([]interface{}); ok {
				if field.Repeater.MinItems > 0 && len(items) < field.Repeater.MinItems {
					return fmt.Errorf("field %s requires at least %d items", field.Name, field.Repeater.MinItems)
				}
				if field.Repeater.MaxItems > 0 && len(items) > field.Repeater.MaxItems {
					return fmt.Errorf("field %s allows at most %d items", field.Name, field.Repeater.MaxItems)
				}
			}
		}
	}

	return nil
}
