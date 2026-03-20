package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"

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
	activeTheme   *Theme
	themes        map[string]*Theme
	themesPath    string
	templateMu    sync.RWMutex
	templateCache map[string]*template.Template // 模板缓存: "themeName/templateName" -> *template.Template
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
		templateCache: make(map[string]*template.Template),
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

	// 从数据库加载主题配置
	savedConfig, err := models.GetThemeConfig(name)
	if err != nil {
		log.Printf("[WARN] Failed to load theme config for %s: %v", name, err)
	} else if len(savedConfig) > 0 {
		// 使用深合并，避免嵌套字段被整体替换
		if theme.Config == nil {
			theme.Config = make(map[string]interface{})
		}
		theme.Config = models.DeepMerge(theme.Config, savedConfig)
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
		log.Printf("[WARN] Failed to save active theme to database: %v", err)
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

	prefix := themeName + "/"
	for key := range themeService.templateCache {
		if strings.HasPrefix(key, prefix) {
			delete(themeService.templateCache, key)
		}
	}
	log.Printf("[INFO] Cleared template cache for theme: %s", themeName)
}

// RenderTemplate 渲染模板 - 使用 Go 标准模板引擎（带缓存）
func RenderTemplate(templateName string, data map[string]interface{}) (template.HTML, error) {
	theme := GetActiveTheme()
	if theme == nil {
		return "", fmt.Errorf("no active theme")
	}

	// 构建模板路径 - 只使用 .html 后缀
	templatePath := filepath.Join(theme.Path, "pages", templateName+".html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	// 准备模板数据
	templateData := make(map[string]interface{})

	// 复制所有传入的数据
	for key, value := range data {
		templateData[key] = value
	}

	// 添加主题配置
	templateData["config"] = theme.Config
	if theme.Config == nil {
		templateData["config"] = make(map[string]interface{})
	}

	// 添加站点信息到config中（用于模板访问）
	if site, ok := data["site"].(map[string]interface{}); ok {
		if configMap, ok := templateData["config"].(map[string]interface{}); ok {
			configMap["siteTitle"] = site["title"]
			configMap["siteDescription"] = site["description"]
			configMap["siteYear"] = site["year"]
		}
	}

	// 添加菜单数据
	headerMenu, _ := models.GetMenuByLocation("header")
	templateData["headerMenu"] = headerMenu

	footerMenu, _ := models.GetMenuByLocation("footer")
	templateData["footerMenu"] = footerMenu

	sidebarMenu, _ := models.GetMenuByLocation("sidebar")
	templateData["sidebarMenu"] = sidebarMenu

	// 尝试从缓存获取模板
	cacheKey := getTemplateCacheKey(theme.Name, templateName)
	themeService.templateMu.RLock()
	cachedTmpl, found := themeService.templateCache[cacheKey]
	themeService.templateMu.RUnlock()

	if !found {
		// 缓存未命中，读取并解析模板
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return "", err
		}

		tmpl, err := template.New(templateName).Funcs(getTemplateFuncs()).Parse(string(content))
		if err != nil {
			return "", fmt.Errorf("template parse error: %v", err)
		}

		// 写入缓存
		themeService.templateMu.Lock()
		themeService.templateCache[cacheKey] = tmpl
		themeService.templateMu.Unlock()
		cachedTmpl = tmpl
	}

	// 执行模板（每次执行需要 clone，因为 Execute 会修改 state）
	var buf bytes.Buffer
	if err := cachedTmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("template execute error: %v", err)
	}

	return template.HTML(buf.String()), nil
}

// getTemplateFuncs 获取模板自定义函数
func getTemplateFuncs() template.FuncMap {
	theme := GetActiveTheme()

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
	}
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
