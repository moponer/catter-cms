package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/models"

	"github.com/yuin/goldmark"
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

// FlattenMap 将嵌套 map 转换为扁平 map（key 用点号连接）
// 例如: {appearance: {colors: {primary: "#000"}}} -> {"appearance.colors.primary": "#000"}
func FlattenMap(nested map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	flattenRecursive(nested, "", result)
	return result
}

// flattenToNestedForTemplate 将扁平 map（点号分隔的 key）转换为嵌套结构
// 例如: {"modules.hero.enabled": true} -> {"modules": {"hero": {"enabled": true}}}
// 与 FlattenMap 相反，用于将存储的扁平格式转换为模板可访问的嵌套格式
func flattenToNestedForTemplate(flat map[string]interface{}) map[string]interface{} {
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

func flattenRecursive(m map[string]interface{}, prefix string, result map[string]interface{}) {
	for key, value := range m {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}
		if childMap, ok := value.(map[string]interface{}); ok {
			flattenRecursive(childMap, newKey, result)
		} else {
			result[newKey] = value
		}
	}
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

// TemplateSet 模板集 - 存储组合后的所有模板
// 使用组合模板集可以让 {{template "name"}} 调用正常工作
type TemplateSet struct {
	Root     *template.Template                    // 根模板（包含所有已解析的模板）
	Template map[string]*template.Template        // 按名称索引的模板
}

// ThemeService 主题服务
type ThemeService struct {
	activeTheme   *Theme
	themes        map[string]*Theme
	themesPath    string
	templateMu    sync.RWMutex
	templateCache map[string]*TemplateSet // 模板缓存: "themeName" -> *TemplateSet
}

// 全局主题服务
var themeService *ThemeService

// InitThemeService 初始化主题服务
func InitThemeService() error {
	var themesPath string

	// 1. 优先使用配置文件中的 theme_path
	if cfg := config.Get(); cfg != nil && cfg.Server.ThemePath != "" {
		if _, err := os.Stat(cfg.Server.ThemePath); err == nil {
			themesPath = cfg.Server.ThemePath
			log.Printf("[INFO] Using theme_path from config: %s", themesPath)
		} else {
			log.Printf("[WARN] Config theme_path '%s' not found, trying auto-detection", cfg.Server.ThemePath)
		}
	}

	// 2. 配置不存在或无效时，自动检测路径
	if themesPath == "" {
		themesPaths := []string{
			"./templates/themes",         // 根目录运行（方案一：独立部署）
			"../templates/themes",        // gin-app 目录运行（方案二）
			"./gin-app/templates/themes", // 根目录但模板在 gin-app 下
		}

		for _, path := range themesPaths {
			if _, err := os.Stat(path); err == nil {
				themesPath = path
				log.Printf("[INFO] Auto-detected themes at: %s", themesPath)
				break
			}
		}
	}

	if themesPath == "" {
		log.Printf("[WARN] Themes directory not found, themes will not be available")
		themesPath = "./templates/themes" // fallback
	}

	themeService = &ThemeService{
		themes:        make(map[string]*Theme),
		themesPath:    themesPath,
		templateCache: make(map[string]*TemplateSet),
	}

	// 加载主题
	if err := themeService.LoadThemes(themeService.themesPath); err != nil {
		log.Printf("[WARN] Failed to load themes: %v", err)
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

		theme, err := s.LoadTheme(entry.Name(), themePath)
		if err != nil {
			log.Printf("[WARN] Failed to load theme %s: %v", entry.Name(), err)
			continue
		}

		s.themes[entry.Name()] = theme
		log.Printf("[INFO] Loaded theme: %s (v%s)", theme.Name, theme.Version)
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

	// 自动检测主题缩略图
	if theme.Thumbnail == "" {
		// 常见的截图文件名
		screenshotFiles := []string{"screenshot.png", "screenshot.jpg", "thumbnail.png", "thumbnail.jpg", "preview.png", "preview.jpg"}
		for _, screenshot := range screenshotFiles {
			screenshotPath := filepath.Join(path, screenshot)
			if _, err := os.Stat(screenshotPath); err == nil {
				// 缩略图路径使用 /themes 前缀（与 main.go 中的 r.Static("/themes", ...) 对应）
				theme.Thumbnail = "/themes/" + name + "/" + screenshot
				break
			}
		}
	}

	// 关键修复：先将 YAML 的嵌套 config 扁平化，再与数据库的扁平格式合并
	// 这样 DeepMerge(flat, flat) 就能正确处理所有 key
	if theme.Config != nil {
		theme.Config = FlattenMap(theme.Config)
	}

	// 从数据库加载主题配置（数据库中已经是扁平格式）
	savedConfig, err := models.GetThemeConfig(name)
	if err != nil {
		log.Printf("[WARN] Failed to load theme config for %s: %v", name, err)
	} else if len(savedConfig) > 0 {
		// 扁平合并：数据库值覆盖 YAML 默认值
		if theme.Config == nil {
			theme.Config = make(map[string]interface{})
		}
		theme.Config = models.DeepMerge(theme.Config, savedConfig)
	}

	return theme, nil
}

// SetActiveTheme 设置活动主题
func (s *ThemeService) SetActiveTheme(name string) error {
	log.Printf("[DEBUG] SetActiveTheme called with: %s", name)
	log.Printf("[DEBUG] Currently loaded themes: %v", getLoadedThemeNames())

	theme, ok := s.themes[name]
	if !ok {
		// 主题不在内存中，尝试重新加载所有主题
		log.Printf("[WARN] Theme %s not found in memory, reloading themes...", name)
		if err := s.LoadThemes(s.themesPath); err != nil {
			return fmt.Errorf("failed to reload themes: %v", err)
		}
		theme, ok = s.themes[name]
		if !ok {
			log.Printf("[ERROR] Theme %s still not found after reload. Available themes: %v", name, getLoadedThemeNames())
			return fmt.Errorf("theme not found: %s (even after reload)", name)
		}
	}

	log.Printf("[DEBUG] Theme %s found, path: %s", theme.Name, theme.Path)

	// 清理旧主题的模板缓存，确保使用新主题
	s.templateMu.Lock()
	delete(s.templateCache, name)
	s.templateMu.Unlock()

	s.activeTheme = theme

	// 持久化到数据库
	if err := models.SetActiveTheme(name); err != nil {
		log.Printf("[WARN] Failed to save active theme to database: %v", err)
	}

	// 同步预加载所有模板（页面 + 模块），确保立即可用
	s.preloadAllTemplates(theme)

	return nil
}

// getLoadedThemeNames 辅助函数，用于日志
func getLoadedThemeNames() []string {
	if themeService == nil {
		return []string{}
	}
	names := make([]string, 0, len(themeService.themes))
	for name := range themeService.themes {
		names = append(names, name)
	}
	return names
}

// getMapKeys 获取 map 的所有键
func getMapKeys(m map[string]interface{}) []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// preloadAllTemplates 预加载主题的所有模板（页面 + 模块）
// 使用组合解析方式，使 {{template "模块名"}} 调用正常工作
func (s *ThemeService) preloadAllTemplates(theme *Theme) {
	if theme == nil {
		return
	}

	// 收集所有模板文件
	var templateFiles []string
	themePath := theme.Path

	// 获取当前工作目录
	cwd, _ := os.Getwd()
	log.Printf("[DEBUG] CWD: %s", cwd)
	log.Printf("[DEBUG] Preloading templates for theme: %s, path: %s", theme.Name, themePath)

	// 添加 modules 目录下的模板（如果存在）
	modulesPath := filepath.Join(themePath, "modules")
	if info, err := os.Stat(modulesPath); err == nil && info.IsDir() {
		moduleFiles, _ := filepath.Glob(filepath.Join(modulesPath, "*.html"))
		templateFiles = append(templateFiles, moduleFiles...)
		log.Printf("[DEBUG] Found %d module templates in: %s", len(moduleFiles), modulesPath)
	}

	// 添加 pages 目录下的模板
	pagesPath := filepath.Join(themePath, "pages")
	if info, err := os.Stat(pagesPath); err == nil && info.IsDir() {
		pageFiles, _ := filepath.Glob(filepath.Join(pagesPath, "*.html"))
		templateFiles = append(templateFiles, pageFiles...)
		log.Printf("[DEBUG] Found %d page templates in: %s", len(pageFiles), pagesPath)
	}

	if len(templateFiles) == 0 {
		log.Printf("[ERROR] No templates found for theme: %s at path: %s", theme.Name, themePath)
		return
	}

	log.Printf("[DEBUG] Total template files to parse: %d", len(templateFiles))

	// 使用 ParseFiles 解析所有模板（它们之间可以互相引用）
	funcs := getTemplateFuncsForTheme(theme)
	tmpl := template.New("").Funcs(funcs)

	parsedTmpl, err := tmpl.ParseFiles(templateFiles...)
	if err != nil {
		log.Printf("[ERROR] Failed to parse templates for %s: %v", theme.Name, err)
		return
	}

	// 构建模板索引
	templateMap := make(map[string]*template.Template)
	for _, tmpl := range parsedTmpl.Templates() {
		// 从完整路径中提取模板名称（去掉目录和扩展名）
		// 例如: "pages/home.html" -> "home", "modules/header.html" -> "header"
		fullName := tmpl.Name()
		baseName := filepath.Base(fullName)
		templateName := strings.TrimSuffix(baseName, ".html")
		templateMap[templateName] = tmpl
		// 同时存储完整路径作为键（兼容旧代码）
		templateMap[fullName] = tmpl
	}

	// 存储到缓存
	set := &TemplateSet{
		Root:     parsedTmpl,
		Template: templateMap,
	}

	s.templateMu.Lock()
	s.templateCache[theme.Name] = set
	s.templateMu.Unlock()

	log.Printf("[INFO] Preloaded %d templates for theme: %s", len(templateFiles), theme.Name)
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

// getTemplateCacheKey 生成模板缓存 key
func getTemplateCacheKey(themeName, templateName string) string {
	return themeName + "/" + templateName
}

// clearTemplateCache 清除指定主题的模板缓存
func clearTemplateCache(themeName string) {
	if themeService == nil {
		return
	}
	themeService.templateMu.Lock()
	defer themeService.templateMu.Unlock()

	// 直接删除该主题的 TemplateSet
	delete(themeService.templateCache, themeName)
	log.Printf("[INFO] Cleared template cache for theme: %s", themeName)
}

// RenderTemplate 渲染模板 - 使用 Go 标准模板引擎（带组合缓存）
func RenderTemplate(templateName string, data map[string]interface{}) (template.HTML, error) {
	theme := GetActiveTheme()
	if theme == nil {
		return "", fmt.Errorf("no active theme")
	}

	log.Printf("[DEBUG] RenderTemplate called: theme=%s, template=%s", theme.Name, templateName)
	log.Printf("[DEBUG] Theme path: %s", theme.Path)
	log.Printf("[DEBUG] Theme config keys: %v", getMapKeys(theme.Config))

	// 验证模板文件存在
	templatePath := filepath.Join(theme.Path, "pages", templateName+".html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	// 准备模板数据
	templateData := make(map[string]interface{})
	for key, value := range data {
		templateData[key] = value
	}

	// 添加主题配置（同时提供 config 和 themeConfig 两种访问方式）
	nestedConfig := flattenToNestedForTemplate(theme.Config)
	templateData["config"] = nestedConfig
	if nestedConfig == nil {
		templateData["config"] = make(map[string]interface{})
	}
	templateData["themeConfig"] = templateData["config"]

	// 添加站点信息
	if site, ok := data["site"].(map[string]interface{}); ok {
		if configMap, ok := templateData["config"].(map[string]interface{}); ok {
			configMap["siteTitle"] = site["title"]
			configMap["siteDescription"] = site["description"]
			configMap["siteYear"] = site["year"]
		}
		if title, ok := site["title"].(string); ok {
			templateData["siteName"] = title
		}
	}

	// 添加菜单数据
	headerMenu, _ := models.GetMenuByLocation("header")
	templateData["headerMenu"] = headerMenu
	footerMenu, _ := models.GetMenuByLocation("footer")
	templateData["footerMenu"] = footerMenu
	sidebarMenu, _ := models.GetMenuByLocation("sidebar")
	templateData["sidebarMenu"] = sidebarMenu

	// 尝试从缓存获取模板集
	themeService.templateMu.RLock()
	set, found := themeService.templateCache[theme.Name]
	themeService.templateMu.RUnlock()

	if !found {
		// 缓存未命中，重新加载所有模板
		themeService.preloadAllTemplates(theme)
		themeService.templateMu.RLock()
		set, found = themeService.templateCache[theme.Name]
		themeService.templateMu.RUnlock()
	}

	if set == nil || set.Template == nil {
		return "", fmt.Errorf("failed to load template set for theme: %s", theme.Name)
	}

	// 从模板集中获取指定模板
	_, ok := set.Template[templateName]
	if !ok {
		return "", fmt.Errorf("template %s not found in theme %s", templateName, theme.Name)
	}

	// 执行模板（使用 Root 的 ExecuteTemplate 确保正确渲染具名模板）
	var buf bytes.Buffer
	if err := set.Root.ExecuteTemplate(&buf, templateName+".html", templateData); err != nil {
		log.Printf("[ERROR] ExecuteTemplate %q failed: %v", templateName, err)
		return "", fmt.Errorf("template execute error: %v", err)
	}

	return template.HTML(buf.String()), nil
}

// getTemplateFuncsForTheme 获取模板自定义函数（接受主题参数）
func getTemplateFuncsForTheme(theme *Theme) template.FuncMap {
	if theme == nil {
		theme = GetActiveTheme()
	}

	return template.FuncMap{
		// getConfig函数 - 用于获取嵌套配置值（支持点分隔路径）
		"getConfig": func(path string, defaultValue interface{}) interface{} {
			if theme == nil || theme.Config == nil {
				return defaultValue
			}

			// 分割路径
			keys := strings.Split(path, ".")
			current := theme.Config

			for i, key := range keys {
				if i == len(keys)-1 {
					// 最后一层，返回值
					if val, ok := current[key]; ok {
						return val
					}
					return defaultValue
				}
				// 中间层，继续深入
				if next, ok := current[key].(map[string]interface{}); ok {
					current = next
				} else {
					return defaultValue
				}
			}

			return defaultValue
		},

		// len函数 - 用于获取长度
		"len": func(v interface{}) int {
			if v == nil {
				return 0
			}
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
				return rv.Len()
			default:
				return 0
			}
		},

		// asset_path函数 - 用于生成资源路径
		"asset_path": func(path string) string {
			if theme == nil {
				return "/" + path
			}
			return "/themes/" + theme.Name + "/" + path
		},

		// current_theme函数 - 获取当前主题名称
		"current_theme": func() string {
			if theme == nil {
				return ""
			}
			return theme.Name
		},

		// safe函数 - 标记HTML为安全（不转义）
		// 注意：仅在确认内容安全时使用，否则可能导致XSS
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},

		// default函数 - 提供默认值
		"default": func(value interface{}, defaultValue interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},

		// lower函数 - 转小写
		"lower": func(s string) string {
			return strings.ToLower(s)
		},

		// upper函数 - 转大写
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},

		// trim函数 - 去除空格
		"trim": func(s string) string {
			return strings.TrimSpace(s)
		},

		// gt函数 - 大于比较（支持 int 和 int64）
		"gt": func(a, b interface{}) bool {
			toInt64 := func(v interface{}) int64 {
				switch n := v.(type) {
				case int:
					return int64(n)
				case int64:
					return n
				case float64:
					return int64(n)
				default:
					return 0
				}
			}
			return toInt64(a) > toInt64(b)
		},

		// lt函数 - 小于比较（支持 int 和 int64）
		"lt": func(a, b interface{}) bool {
			toInt64 := func(v interface{}) int64 {
				switch n := v.(type) {
				case int:
					return int64(n)
				case int64:
					return n
				case float64:
					return int64(n)
				default:
					return 0
				}
			}
			return toInt64(a) < toInt64(b)
		},

		// eq函数 - 等于比较
		"eq": func(a, b interface{}) bool {
			return a == b
		},

		// add函数 - 加法（支持 int 和 int64）
		"add": func(a, b interface{}) interface{} {
			switch va := a.(type) {
			case int:
				if vb, ok := b.(int); ok {
					return va + vb
				}
				if vb, ok := b.(int64); ok {
					return int64(va) + vb
				}
			case int64:
				if vb, ok := b.(int64); ok {
					return va + vb
				}
				if vb, ok := b.(int); ok {
					return va + int64(vb)
				}
			case float64:
				if vb, ok := b.(float64); ok {
					return va + vb
				}
			}
			return 0
		},

		// sub函数 - 减法（支持 int 和 int64）
		"sub": func(a, b interface{}) interface{} {
			switch va := a.(type) {
			case int:
				if vb, ok := b.(int); ok {
					return va - vb
				}
				if vb, ok := b.(int64); ok {
					return int64(va) - vb
				}
			case int64:
				if vb, ok := b.(int64); ok {
					return va - vb
				}
				if vb, ok := b.(int); ok {
					return va - int64(vb)
				}
			}
			return 0
		},

		// iterate函数 - 生成数字序列（用于分页）
		"iterate": func(n int64) []int {
			result := make([]int, int(n))
			for i := 0; i < int(n); i++ {
				result[i] = i
			}
			return result
		},

		// markdownHTML函数 - 将 Markdown 转换为 HTML
		// 如果内容已经是 HTML（以 < 开头），则直接输出
		"markdownHTML": func(content string) template.HTML {
			// 检测是否已经是 HTML 内容
			trimmed := strings.TrimSpace(content)
			if strings.HasPrefix(trimmed, "<") {
				// 已经是 HTML，直接返回
				return template.HTML(content)
			}
			// 是 Markdown，转换为 HTML
			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(content), &buf); err != nil {
				return template.HTML(content)
			}
			return template.HTML(buf.String())
		},

		// json_decode函数 - 将 JSON 字符串解析为切片或 map
		// 用于从配置中读取列表类型的数据，如 features、services 等
		"json_decode": func(jsonStr string) interface{} {
			if jsonStr == "" {
				return []interface{}{}
			}
			var result interface{}
			if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
				log.Printf("[WARN] json_decode error: %v, input: %s", err, jsonStr)
				return []interface{}{}
			}
			return result
		},

		// ensureSlice函数 - 确保返回切片类型
		// 如果传入的已经是切片，直接返回；如果是字符串，解析后返回
		"ensureSlice": func(v interface{}) []interface{} {
			if v == nil {
				return []interface{}{}
			}
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Slice, reflect.Array:
				result := make([]interface{}, rv.Len())
				for i := 0; i < rv.Len(); i++ {
					result[i] = rv.Index(i).Interface()
				}
				return result
			case reflect.String:
				// 如果是字符串，尝试解析为JSON
				var decoded interface{}
				if err := json.Unmarshal([]byte(rv.String()), &decoded); err == nil {
					if s, ok := decoded.([]interface{}); ok {
						return s
					}
				}
				return []interface{}{}
			}
			return []interface{}{}
		},

		// substr函数 - 字符串截取
		// Go模板没有内置的substr，需要自己实现
		"substr": func(s string, start, length int) string {
			if start < 0 {
				start = 0
			}
			if start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},

		// mod函数 - 取模运算
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},

		// first函数 - 返回切片的前 n 个元素
		// Go模板没有内置的first，需要自己实现
		// 注意：返回原切片的类型，而不是 []interface{}
		"first": func(n int, slice interface{}) interface{} {
			if slice == nil {
				return nil
			}
			rv := reflect.ValueOf(slice)
			switch rv.Kind() {
			case reflect.Slice, reflect.Array:
				if n < 0 {
					n = 0
				}
				if n > rv.Len() {
					n = rv.Len()
				}
				// 创建相同类型的切片
				newSlice := reflect.MakeSlice(rv.Type(), n, n)
				for i := 0; i < n; i++ {
					newSlice.Index(i).Set(rv.Index(i))
				}
				return newSlice.Interface()
			default:
				return nil
			}
		},

		// int函数 - 将值转换为 int 类型
		// 用于模板中将数字字符串或 float64 转换为 int
		"int": func(v interface{}) int {
			switch val := v.(type) {
			case int:
				return val
			case int64:
				return int(val)
			case float64:
				return int(val)
			case float32:
				return int(val)
			case string:
				// 尝试解析字符串为数字
				if i, err := strconv.Atoi(val); err == nil {
					return i
				}
				return 0
			default:
				return 0
			}
		},

		// div函数 - 整数除法
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},

		// dateFormat函数 - 日期格式化
		// 格式字符串使用 Go 的时间格式化语法（参考：2006-01-02 15:04:05）
		"dateFormat": func(format string, date interface{}) string {
			if date == nil {
				return ""
			}
			var t time.Time
			switch v := date.(type) {
			case time.Time:
				t = v
			case *time.Time:
				if v != nil {
					t = *v
				}
			case string:
				// 尝试解析字符串日期
				if parsed, err := time.Parse("2006-01-02", v); err == nil {
					t = parsed
				} else if parsed, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
					t = parsed
				}
			}
			if t.IsZero() {
				return ""
			}
			return t.Format(format)
		},

		// in函数 - 判断元素是否在切片/数组中
		"in": func(item interface{}, slice interface{}) bool {
			if slice == nil {
				return false
			}
			rv := reflect.ValueOf(slice)
			switch rv.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < rv.Len(); i++ {
					if reflect.DeepEqual(item, rv.Index(i).Interface()) {
						return true
					}
				}
			case reflect.Map:
				for _, key := range rv.MapKeys() {
					if reflect.DeepEqual(item, key.Interface()) {
						return true
					}
				}
			case reflect.String:
				// 如果slice是字符串，item也必须是字符串
				if s, ok := item.(string); ok {
					return strings.Contains(rv.String(), s)
				}
			}
			return false
		},

		// where函数 - 过滤切片中的元素（模拟数据库查询）
		"where": func(slice interface{}, field string, value interface{}) []interface{} {
			if slice == nil {
				return []interface{}{}
			}
			rv := reflect.ValueOf(slice)
			if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
				return []interface{}{}
			}

			result := make([]interface{}, 0)
			for i := 0; i < rv.Len(); i++ {
				item := rv.Index(i).Interface()
				itemVal := reflect.ValueOf(item)

				// 尝试获取字段值
				var fieldVal interface{}

				// 直接字段访问
				if f := itemVal.FieldByName(field); f.IsValid() {
					fieldVal = f.Interface()
				} else {
					// 尝试通过map访问
					if itemMap, ok := item.(map[string]interface{}); ok {
						fieldVal = itemMap[field]
					}
				}

				if fieldVal != nil && reflect.DeepEqual(fieldVal, value) {
					result = append(result, item)
				}
			}
			return result
		},
	}
}

// getTemplateFuncs 获取模板自定义函数
func getTemplateFuncs() template.FuncMap {
	return getTemplateFuncsForTheme(nil)
}

// convertPlushToGoTemplate 将Plush语法转换为Go模板语法
// 完整支持：变量定义、条件判断、循环、对象访问等
func convertPlushToGoTemplate(content string) string {
	result := content

	// 1. 移除let变量定义行（Plush特有，Go模板中用.访问数据）
	// <%= let heroEnabled = true %> -> 整行移除
	letPattern := regexp.MustCompile(`<%=\s*let\s+\w+\s*=\s*[^%]+%>`)
	result = letPattern.ReplaceAllString(result, "")

	// 2. 处理复杂的配置访问表达式
	// 将 <%= (config["appearance"] && config["appearance"]["colors"] && config["appearance"]["colors"]["primary"]) || "#3B82F6" %>
	// 简化为 {{ index .config "appearance.colors.primary" | default "#3B82F6" }}
	// 但由于Go模板的限制，我们先用更简单的方式处理

	// 3. 处理简单的config访问
	// <%= config.siteTitle %> -> {{ .config.siteTitle }}
	result = regexp.MustCompile(`<%=\s*config\.(\w+)\s*%>`).ReplaceAllString(result, "{{ .config.$1 }}")

	// 4. 处理config["key"]访问
	// <%= config["siteTitle"] %> -> {{ index .config "siteTitle" }}
	result = regexp.MustCompile(`<%=\s*config\["([^"]+)"\]\s*%>`).ReplaceAllString(result, `{{ index .config "$1" }}`)

	// 5. 处理嵌套的config访问（简化处理）
	// <%= (config["appearance"] && config["appearance"]["colors"] && config["appearance"]["colors"]["primary"]) || "#3B82F6" %>
	// 这是一个复杂的表达式，需要用自定义函数处理
	// 暂时使用正则替换，提取默认值
	configPattern := regexp.MustCompile(`<%=\s*\(config\["[^"]+"\].*\|\|\s*"([^"]+)"\)\s*%>`)
	if configPattern.MatchString(result) {
		// 提取默认值并简化为直接输出默认值（后续可通过模板函数改进）
		matches := configPattern.FindStringSubmatch(result)
		if len(matches) > 1 {
			result = configPattern.ReplaceAllString(result, matches[1])
		}
	}

	// 6. 替换 <%= ... %> 为 {{ ... }}
	result = strings.ReplaceAll(result, "<%=", "{{")
	result = strings.ReplaceAll(result, "%>", "}}")

	// 7. 替换 <% ... %> 为 {{ ... }} (for控制语句)
	result = strings.ReplaceAll(result, "<%", "{{")

	// 8. 处理for循环语法
	// plush: <%= for (item) in items { %> 或 <%= for (item) in items { %>
	// go: {{ range .items }}
	// 需要提取变量名和集合名
	forPattern := regexp.MustCompile(`{{\s*for\s+\((\w+)\)\s+in\s+(\w+)\s*{?\s*}}`)
	result = forPattern.ReplaceAllString(result, "{{ range .{{- $2 := . }}$2 }}\n{{- range $i, $1 := . }}")

	// 更简单的for循环处理
	result = regexp.MustCompile(`{{\s*for\s+\((\w+)\)\s+in\s+(\w+)\s*}}`).ReplaceAllString(result, "{{ range .{{- $2 := . }}$2 }}")

	// 修正：直接使用Go模板的range语法
	result = regexp.MustCompile(`{{\s*for\s+\((\w+)\)\s+in\s+(\w+)\s*{?\s*}}`).ReplaceAllString(result, "{{ range .{{- $2 }} }}")

	// 9. 处理if条件
	// plush: <%= if (condition) { %> ... <% } %>
	// go: {{ if .condition }} ... {{ end }}
	result = regexp.MustCompile(`{{\s*if\s+\(([^)]+)\)\s*{?\s*}}`).ReplaceAllString(result, "{{ if $1 }}")

	// 10. 处理else if
	result = regexp.MustCompile(`{{\s*}\s*else\s*if\s+\(([^)]+)\)\s*{?\s*}}`).ReplaceAllString(result, "{{ else if $1 }}")

	// 11. 处理else
	result = regexp.MustCompile(`{{\s*}\s*else\s*{?\s*}}`).ReplaceAllString(result, "{{ else }}")

	// 12. 处理end/}
	result = regexp.MustCompile(`{{\s*}\s*}}`).ReplaceAllString(result, "{{ end }}")

	// 13. 清理空的模板控制语句
	result = regexp.MustCompile(`{{\s*}}`).ReplaceAllString(result, "")

	// 14. 处理布尔比较
	// == false -> not
	result = regexp.MustCompile(`(\w+)\s*==\s*false`).ReplaceAllString(result, "not .$1")

	// == true -> 直接访问
	result = regexp.MustCompile(`(\w+)\s*==\s*true`).ReplaceAllString(result, ".$1")

	// 15. 清理let语句中的变量赋值（已废弃）
	// 例如：{{ showSidebar = layout != "one-column" }} 这种表达式Go模板无法处理
	// 需要在后端预处理或使用自定义函数

	return result
}

// RenderTemplateString 渲染模板字符串
func RenderTemplateString(content string, data map[string]interface{}) (template.HTML, error) {
	// 使用标准模板引擎
	tmpl, err := template.New("dynamic").Funcs(getTemplateFuncs()).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
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

	// 从默认配置开始
	config := make(map[string]interface{})
	if theme.Config != nil {
		for key, value := range theme.Config {
			config[key] = value
		}
	}

	// 从数据库加载配置并深合并
	savedConfig, err := models.GetThemeConfig(themeName)
	if err == nil && savedConfig != nil {
		config = models.DeepMerge(config, savedConfig)
	}

	return config, nil
}

// UpdateThemeConfig 更新主题配置（保存到数据库并更新内存）
func UpdateThemeConfig(themeName string, config map[string]interface{}) error {
	// 保存到数据库
	if err := models.SaveThemeConfig(themeName, config); err != nil {
		return err
	}

	// 从数据库重新加载配置（会自动转换为嵌套结构）
	reloadedConfig, err := models.GetThemeConfig(themeName)
	if err != nil {
		return fmt.Errorf("failed to reload theme config: %v", err)
	}

	// 更新内存中的配置
	theme := GetTheme(themeName)
	if theme != nil {
		theme.Config = reloadedConfig
		// 清除该主题的模板缓存，下次渲染会重新加载
		clearTemplateCache(themeName)
	} else {
		log.Printf("[WARN] Theme %s not found in memory", themeName)
	}

	return nil
}

// ClearThemeCache 清理主题所有缓存（模板缓存 + 重新加载主题配置）
func ClearThemeCache(themeName string) error {
	if themeService == nil {
		return fmt.Errorf("theme service not initialized")
	}

	// 1. 清理模板缓存
	themeService.templateMu.Lock()
	delete(themeService.templateCache, themeName)
	themeService.templateMu.Unlock()
	log.Printf("[INFO] Cleared template cache for theme: %s", themeName)

	// 2. 重新加载主题配置
	themePath := filepath.Join(themeService.themesPath, themeName)
	theme, err := themeService.LoadTheme(themeName, themePath)
	if err != nil {
		return fmt.Errorf("failed to reload theme: %v", err)
	}

	themeService.themes[themeName] = theme

	// 3. 如果是当前激活的主题，更新 activeTheme 指针并预加载模板
	if themeService.activeTheme != nil && themeService.activeTheme.Name == themeName {
		themeService.activeTheme = theme
		themeService.preloadAllTemplates(theme)
		log.Printf("[INFO] Reloaded active theme: %s", themeName)
	}

	log.Printf("[INFO] Theme cache cleared and reloaded: %s", themeName)
	return nil
}

// DeleteThemeConfig 删除主题配置
func DeleteThemeConfig(themeName string) error {
	// 删除数据库中的配置
	if err := models.DeleteThemeConfig(themeName); err != nil {
		return err
	}

	// 清除内存中的配置
	theme := GetTheme(themeName)
	if theme != nil {
		theme.Config = make(map[string]interface{})
		clearTemplateCache(themeName)
	}

	log.Printf("[INFO] Deleted theme config for %s", themeName)
	return nil
}
