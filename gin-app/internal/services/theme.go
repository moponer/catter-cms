package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"catter-cms-gin/internal/models"

	"github.com/gobuffalo/plush/v4"
	"gopkg.in/yaml.v3"
)

// ConfigSection 配置分组结构
type ConfigSection struct {
	Section     string  `yaml:"section" json:"section"`
	Label       string  `yaml:"label" json:"label"`
	Description string  `yaml:"description" json:"description"`
	Icon        string  `yaml:"icon" json:"icon"`
	Fields      []Field `yaml:"fields" json:"fields"`
}

// Field 配置字段结构
type Field struct {
	Name        string      `yaml:"name" json:"name"`
	Label       string      `yaml:"label" json:"label"`
	Type        string      `yaml:"type" json:"type"`
	Description string      `yaml:"description" json:"description"`
	Placeholder string      `yaml:"placeholder" json:"placeholder"`
	Required    bool        `yaml:"required" json:"required"`
	Default     interface{} `yaml:"default" json:"default"`
	Options     []Option    `yaml:"options,omitempty" json:"options,omitempty"`
	Min         float64     `yaml:"min,omitempty" json:"min,omitempty"`
	Max         float64     `yaml:"max,omitempty" json:"max,omitempty"`
	Step        float64     `yaml:"step,omitempty" json:"step,omitempty"`
	Rows        int         `yaml:"rows,omitempty" json:"rows,omitempty"`
	ShowAlpha   bool        `yaml:"showAlpha,omitempty" json:"showAlpha,omitempty"`
	ShowIf      string      `yaml:"show_if,omitempty" json:"show_if,omitempty"`
}

// Option 选项结构
type Option struct {
	Label       string `yaml:"label" json:"label"`
	Value       string `yaml:"value" json:"value"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// Theme 主题结构
type Theme struct {
	Name             string                 `yaml:"name"`
	Version          string                 `yaml:"version"`
	Description      string                 `yaml:"description"`
	Author           string                 `yaml:"author"`
	Screenshots      []string               `yaml:"screenshots"`
	Thumbnail        string                 `yaml:"thumbnail"`
	Path             string                 // 主题路径
	Config           map[string]interface{} `yaml:"config"`
	ConfigDefinition []ConfigSection        `yaml:"configDefinition"`
}

// ThemeService 主题服务
type ThemeService struct {
	activeTheme *Theme
	themes      map[string]*Theme
	themesPath  string
}

// 全局主题服务
var themeService *ThemeService

// InitThemeService 初始化主题服务
func InitThemeService() error {
	themeService = &ThemeService{
		themes:     make(map[string]*Theme),
		themesPath: "./templates/themes",
	}

	// 加载主题
	if err := themeService.LoadThemes(themeService.themesPath); err != nil {
		fmt.Printf("Warning: failed to load themes: %v\n", err)
	}

	// 从数据库获取激活的主题
	activeThemeName := models.GetActiveTheme()
	if activeThemeName == "" {
		activeThemeName = "catter-blog"
	}

	// 设置激活主题
	if err := themeService.SetActiveTheme(activeThemeName); err != nil {
		// 如果设置失败，使用第一个可用主题
		for name := range themeService.themes {
			if err := themeService.SetActiveTheme(name); err != nil {
				return err
			}
			// 保存到数据库
			models.SetActiveTheme(name)
			break
		}
	}

	return nil
}

// LoadThemes 加载所有主题
func (s *ThemeService) LoadThemes(themesPath string) error {
	// 检查目录是否存在
	if _, err := os.Stat(themesPath); os.IsNotExist(err) {
		return fmt.Errorf("themes directory not found: %s", themesPath)
	}

	// 遍历主题目录
	entries, err := os.ReadDir(themesPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		themePath := filepath.Join(themesPath, entry.Name())
		themeYml := filepath.Join(themePath, "theme.yml")

		if _, err := os.Stat(themeYml); err != nil {
			continue // 跳过没有 theme.yml 的目录
		}

		// 加载主题
		theme, err := s.LoadTheme(entry.Name(), themePath)
		if err != nil {
			fmt.Printf("Warning: failed to load theme %s: %v\n", entry.Name(), err)
			continue
		}

		s.themes[entry.Name()] = theme
		fmt.Printf("Loaded theme: %s (v%s)\n", theme.Name, theme.Version)
	}

	return nil
}

// LoadTheme 加载单个主题
func (s *ThemeService) LoadTheme(name, path string) (*Theme, error) {
	themeYml := filepath.Join(path, "theme.yml")

	data, err := os.ReadFile(themeYml)
	if err != nil {
		return nil, err
	}

	// 先解析 YAML
	theme := &Theme{}
	if err := yaml.Unmarshal(data, theme); err != nil {
		return nil, err
	}

	// 设置名称和路径（YAML中可能不包含，或者需要覆盖）
	theme.Name = name
	theme.Path = path

	// 从数据库加载主题配置
	savedConfig, err := models.GetThemeConfig(name)
	if err != nil {
		fmt.Printf("Warning: failed to load theme config for %s: %v\n", name, err)
	} else if len(savedConfig) > 0 {
		// 合并保存的配置
		if theme.Config == nil {
			theme.Config = make(map[string]interface{})
		}
		for key, value := range savedConfig {
			theme.Config[key] = value
		}
	}

	return theme, nil
}

// SetActiveTheme 设置活动主题
func (s *ThemeService) SetActiveTheme(name string) error {
	theme, ok := s.themes[name]
	if !ok {
		return fmt.Errorf("theme not found: %s", name)
	}

	s.activeTheme = theme

	// 持久化到数据库
	if err := models.SetActiveTheme(name); err != nil {
		fmt.Printf("Warning: failed to save active theme to database: %v\n", err)
	}

	return nil
}

// SetActiveTheme 设置活动主题（包级别函数）
func SetActiveTheme(name string) error {
	if themeService == nil {
		return fmt.Errorf("theme service not initialized")
	}
	return themeService.SetActiveTheme(name)
}

// GetActiveTheme 获取活动主题
func GetActiveTheme() *Theme {
	if themeService == nil {
		return nil
	}
	return themeService.activeTheme
}

// GetTheme 获取主题
func GetTheme(name string) *Theme {
	if themeService == nil {
		return nil
	}
	return themeService.themes[name]
}

// GetAllThemes 获取所有主题
func GetAllThemes() map[string]*Theme {
	if themeService == nil {
		return nil
	}
	return themeService.themes
}

// ReloadTheme 重新加载指定主题（用于配置更新后）
func ReloadTheme(name string) error {
	if themeService == nil {
		return fmt.Errorf("theme service not initialized")
	}

	themePath := filepath.Join(themeService.themesPath, name)
	theme, err := themeService.LoadTheme(name, themePath)
	if err != nil {
		return err
	}

	themeService.themes[name] = theme

	// 如果是当前激活的主题，更新activeTheme指针
	if themeService.activeTheme != nil && themeService.activeTheme.Name == name {
		themeService.activeTheme = theme
	}

	return nil
}

// SaveThemeConfigToDB 保存主题配置到数据库
func SaveThemeConfigToDB(themeName string, config map[string]interface{}) error {
	return models.SaveThemeConfig(themeName, config)
}

// RenderTemplate 渲染模板 - 使用Plush引擎
func RenderTemplate(templateName string, data map[string]interface{}) (template.HTML, error) {
	theme := GetActiveTheme()
	if theme == nil {
		return "", fmt.Errorf("no active theme")
	}

	// 构建模板路径
	templatePath := filepath.Join(theme.Path, "pages", templateName+".plush.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// 尝试 .html 后缀
		templatePath = filepath.Join(theme.Path, "pages", templateName+".html")
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return "", fmt.Errorf("template not found: %s", templateName)
		}
	}

	// 读取模板内容
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	// 使用 Plush 引擎渲染
	ctx := plush.NewContext()

	// 添加所有数据到上下文
	for key, value := range data {
		ctx.Set(key, value)
	}

	// 添加主题配置
	if theme.Config != nil {
		ctx.Set("config", theme.Config)
	}

	// 添加辅助函数
	ctx.Set("asset_path", func(path string) string {
		return "/themes/" + theme.Name + "/" + path
	})

	ctx.Set("current_theme", func() string {
		return theme.Name
	})

	// 添加hasWidget辅助函数，用于检查小工具是否在配置中启用
	ctx.Set("hasWidget", func(widgetName string) bool {
		if theme.Config == nil {
			return false
		}
		// 检查首页小工具配置
		if widgets, ok := theme.Config["pages.home.sidebarWidgets"].([]interface{}); ok {
			for _, w := range widgets {
				if w == widgetName {
					return true
				}
			}
		}
		// 检查文章页小工具配置
		if widgets, ok := theme.Config["pages.post.sidebarWidgets"].([]interface{}); ok {
			for _, w := range widgets {
				if w == widgetName {
					return true
				}
			}
		}
		// 检查分类页小工具配置
		if widgets, ok := theme.Config["pages.category.sidebarWidgets"].([]interface{}); ok {
			for _, w := range widgets {
				if w == widgetName {
					return true
				}
			}
		}
		return false
	})

	// 执行模板渲染
	rendered, err := plush.Render(string(content), ctx)
	if err != nil {
		// Plush渲染失败，尝试使用标准模板引擎作为降级方案
		return renderWithStandardTemplate(templateName, string(content), data)
	}

	return template.HTML(rendered), nil
}

// renderWithStandardTemplate 使用标准模板引擎渲染（降级方案）
func renderWithStandardTemplate(templateName, content string, data map[string]interface{}) (template.HTML, error) {
	// 转换Plush语法到Go模板语法
	convertedContent := convertPlushToGoTemplate(content)

	tmpl, err := template.New(templateName).Parse(convertedContent)
	if err != nil {
		return "", fmt.Errorf("template parse error: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute error: %v", err)
	}

	return template.HTML(buf.String()), nil
}

// convertPlushToGoTemplate 将Plush语法转换为Go模板语法（简单转换）
func convertPlushToGoTemplate(content string) string {
	// 这是一个简单的转换，处理最常见的Plush语法
	result := content

	// 替换 <%= ... %> 为 {{ ... }}
	result = strings.ReplaceAll(result, "<%=", "{{")
	result = strings.ReplaceAll(result, "%>", "}}")

	// 替换 <% ... %> 为 {{ ... }}
	result = strings.ReplaceAll(result, "<%", "{{")

	// 替换 for 循环语法
	// plush: <%= for (item) in items { %>
	// go: {{ range .items }}
	result = strings.ReplaceAll(result, "for (", "range ")
	result = strings.ReplaceAll(result, ") in ", " .")

	// 替换 end
	result = strings.ReplaceAll(result, "<% } %>", "{{ end }}")

	return result
}

// RenderTemplateString 渲染模板字符串
func RenderTemplateString(content string, data map[string]interface{}) (template.HTML, error) {
	// 尝试使用Plush渲染
	ctx := plush.NewContext()
	for key, value := range data {
		ctx.Set(key, value)
	}

	rendered, err := plush.Render(content, ctx)
	if err != nil {
		// 降级到标准模板
		tmpl, err := template.New("dynamic").Parse(content)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	}

	return template.HTML(rendered), nil
}

// GetThemeConfig 获取主题配置
func GetThemeConfig(key string) interface{} {
	theme := GetActiveTheme()
	if theme == nil {
		return nil
	}

	// 简化实现：直接返回顶层配置
	return theme.Config[key]
}

// GetThemePath 获取主题路径
func GetThemePath() string {
	theme := GetActiveTheme()
	if theme == nil {
		return ""
	}
	return theme.Path
}

// GetThemeFullConfig 获取主题的完整配置（包含数据库中的配置）
func GetThemeFullConfig(themeName string) (map[string]interface{}, error) {
	theme := GetTheme(themeName)
	if theme == nil {
		return nil, fmt.Errorf("theme not found: %s", themeName)
	}

	// 合并主题默认配置和数据库配置
	config := make(map[string]interface{})

	// 先加载默认配置
	if theme.Config != nil {
		for key, value := range theme.Config {
			config[key] = value
		}
	}

	// 从数据库加载配置并覆盖
	savedConfig, err := models.GetThemeConfig(themeName)
	if err == nil && savedConfig != nil {
		for key, value := range savedConfig {
			config[key] = value
		}
	}

	return config, nil
}

// UpdateThemeConfig 更新主题配置（保存到数据库并更新内存）
func UpdateThemeConfig(themeName string, config map[string]interface{}) error {
	// 保存到数据库
	if err := models.SaveThemeConfig(themeName, config); err != nil {
		return err
	}

	// 更新内存中的配置
	theme := GetTheme(themeName)
	if theme != nil {
		if theme.Config == nil {
			theme.Config = make(map[string]interface{})
		}
		for key, value := range config {
			theme.Config[key] = value
		}
	}

	return nil
}
