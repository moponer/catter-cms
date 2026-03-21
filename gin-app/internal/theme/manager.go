package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Manager 主题管理器
type Manager struct {
	mu            sync.RWMutex
	themes        map[string]*Theme
	currentTheme  string
	templatesDir  string
	configDir     string
}

// Theme 主题信息
type Theme struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Version     string                 `json:"version" yaml:"version"`
	Author      string                 `json:"author" yaml:"author"`
	Thumbnail   string                 `json:"thumbnail" yaml:"thumbnail"`
	Path        string                 `json:"path"`
	Settings    map[string]interface{} `json:"settings"`
}

// NewManager 创建主题管理器
func NewManager(templatesDir, configDir string) *Manager {
	return &Manager{
		themes:       make(map[string]*Theme),
		templatesDir: templatesDir,
		configDir:    configDir,
	}
}

// SetTemplateDir 设置模板目录
func (m *Manager) SetTemplateDir(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.templatesDir = dir
}

// SetConfigDir 设置配置目录
func (m *Manager) SetConfigDir(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configDir = dir
}

// LoadThemes 加载所有主题
func (m *Manager) LoadThemes() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries, err := os.ReadDir(m.templatesDir)
	if err != nil {
		return fmt.Errorf("读取主题目录失败: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		themePath := filepath.Join(m.templatesDir, entry.Name())
		theme, err := m.loadThemeInfo(entry.Name(), themePath)
		if err != nil {
			continue // 跳过无法加载的主题
		}

		m.themes[theme.ID] = theme
	}

	// 加载当前主题配置
	m.loadCurrentTheme()

	return nil
}

// loadThemeInfo 加载主题信息
func (m *Manager) loadThemeInfo(id, path string) (*Theme, error) {
	// 读取 theme.yaml
	themeFile := filepath.Join(path, "theme.yaml")
	data, err := os.ReadFile(themeFile)
	if err != nil {
		return nil, err
	}

	theme := &Theme{
		ID:   id,
		Path: path,
	}

	// 使用简单的 YAML 解析（实际项目中使用 yaml 库）
	// 这里简化处理，直接解析为 map
	// TODO: 使用 yaml.Unmarshal

	// 加载配置
	configFile := filepath.Join(m.configDir, id+".json")
	if configData, err := os.ReadFile(configFile); err == nil {
		json.Unmarshal(configData, &theme.Settings)
	}

	return theme, nil
}

// loadCurrentTheme 加载当前主题
func (m *Manager) loadCurrentTheme() {
	// 从文件或数据库读取当前主题
	currentFile := filepath.Join(m.configDir, "current.txt")
	if data, err := os.ReadFile(currentFile); err == nil {
		m.currentTheme = string(data)
	}

	// 如果没有设置或设置的主题不存在，使用第一个
	if m.currentTheme == "" || m.themes[m.currentTheme] == nil {
		for id := range m.themes {
			m.currentTheme = id
			break
		}
	}
}

// GetTheme 获取主题
func (m *Manager) GetTheme(id string) (*Theme, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if id == "" {
		id = m.currentTheme
	}

	theme, ok := m.themes[id]
	if !ok {
		return nil, fmt.Errorf("主题不存在: %s", id)
	}

	return theme, nil
}

// GetThemeConfig 获取主题配置
func (m *Manager) GetThemeConfig(themeID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return nil, fmt.Errorf("主题不存在: %s", themeID)
	}

	// 如果设置为空，返回空 map
	if theme.Settings == nil {
		return make(map[string]interface{}), nil
	}

	return theme.Settings, nil
}

// SaveThemeConfig 保存主题配置
func (m *Manager) SaveThemeConfig(themeID string, config map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return fmt.Errorf("主题不存在: %s", themeID)
	}

	// 更新内存中的配置
	theme.Settings = config

	// 保存到文件
	configFile := filepath.Join(m.configDir, themeID+".json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}

// GetThemeSchema 获取主题设置定义
func (m *Manager) GetThemeSchema(themeID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return nil, fmt.Errorf("主题不存在: %s", themeID)
	}

	// 读取 settings.yaml
	settingsFile := filepath.Join(theme.Path, "settings.yaml")
	if _, err := os.Stat(settingsFile); err != nil {
		// 如果没有 settings.yaml，返回基础结构
		return map[string]interface{}{
			"fields": []map[string]interface{}{},
		}, nil
	}

	// TODO: 解析 YAML 文件
	// 这里简化处理，实际应该使用 yaml 库解析
	return map[string]interface{}{
		"fields": []map[string]interface{}{},
	}, nil
}

// ResetThemeConfig 重置主题配置
func (m *Manager) ResetThemeConfig(themeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return fmt.Errorf("主题不存在: %s", themeID)
	}

	// 清空配置
	theme.Settings = make(map[string]interface{})

	// 删除配置文件
	configFile := filepath.Join(m.configDir, themeID+".json")
	os.Remove(configFile)

	return nil
}

// GetCurrentThemeID 获取当前主题ID
func (m *Manager) GetCurrentThemeID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentTheme
}

// GetCurrentTheme 获取当前主题
func (m *Manager) GetCurrentTheme() *Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.themes[m.currentTheme]
}

// GetAllThemes 获取所有主题
func (m *Manager) GetAllThemes() []*Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*Theme, 0, len(m.themes))
	for _, theme := range m.themes {
		list = append(list, theme)
	}
	return list
}

// SwitchTheme 切换主题
func (m *Manager) SwitchTheme(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.themes[id]; !ok {
		return fmt.Errorf("主题不存在: %s", id)
	}

	m.currentTheme = id

	// 保存到文件
	currentFile := filepath.Join(m.configDir, "current.txt")
	if err := os.WriteFile(currentFile, []byte(id), 0644); err != nil {
		return fmt.Errorf("保存当前主题失败: %w", err)
	}

	return nil
}

// GetThemeConfig 获取主题配置
func (m *Manager) GetThemeConfig(themeID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return nil, fmt.Errorf("主题不存在: %s", themeID)
	}

	// 如果设置为空，返回空 map
	if theme.Settings == nil {
		return make(map[string]interface{}), nil
	}

	return theme.Settings, nil
}

// SaveThemeConfig 保存主题配置
func (m *Manager) SaveThemeConfig(themeID string, config map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return fmt.Errorf("主题不存在: %s", themeID)
	}

	// 更新内存中的配置
	theme.Settings = config

	// 保存到文件
	configFile := filepath.Join(m.configDir, themeID+".json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}

// GetThemeSchema 获取主题设置定义
func (m *Manager) GetThemeSchema(themeID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return nil, fmt.Errorf("主题不存在: %s", themeID)
	}

	// 读取 settings.yaml
	settingsFile := filepath.Join(theme.Path, "settings.yaml")
	if _, err := os.Stat(settingsFile); err != nil {
		// 如果没有 settings.yaml，返回基础结构
		return map[string]interface{}{
			"fields": []map[string]interface{}{},
		}, nil
	}

	// TODO: 解析 YAML 文件
	// 这里简化处理，实际应该使用 yaml 库解析
	return map[string]interface{}{
		"fields": []map[string]interface{}{},
	}, nil
}

// ResetThemeConfig 重置主题配置
func (m *Manager) ResetThemeConfig(themeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if themeID == "" {
		themeID = m.currentTheme
	}

	theme, ok := m.themes[themeID]
	if !ok {
		return fmt.Errorf("主题不存在: %s", themeID)
	}

	// 清空配置
	theme.Settings = make(map[string]interface{})

	// 删除配置文件
	configFile := filepath.Join(m.configDir, themeID+".json")
	os.Remove(configFile)

	return nil
}

// GetCurrentThemeID 获取当前主题ID
func (m *Manager) GetCurrentThemeID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentTheme
}
