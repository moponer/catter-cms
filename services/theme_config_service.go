package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigCacheEntry 配置缓存条目
type ConfigCacheEntry struct {
	Config    *ThemeConfig
	ExpiresAt time.Time
}

// 全局配置缓存
var (
	configCache     sync.Map
	cacheDuration   = 5 * time.Minute
	cacheEnabled    = true
)

// ThemeConfig 主题配置结构
type ThemeConfig struct {
	Name            string                 `yaml:"name"`
	Version         string                 `yaml:"version"`
	APIVersion      string                 `yaml:"apiVersion"`
	Description     string                 `yaml:"description"`
	Author          string                 `yaml:"author"`
	AuthorEmail     string                 `yaml:"authorEmail"`
	AuthorUrl       string                 `yaml:"authorUrl"`
	License         string                 `yaml:"license"`
	Tags            []string               `yaml:"tags"`
	Screenshots     []string               `yaml:"screenshots"`
	Thumbnail       string                 `yaml:"thumbnail"`
	Requires        map[string]string      `yaml:"requires"`
	Capabilities    Capabilities           `yaml:"capabilities"`
	Config          Config                 `yaml:"config"`
	ConfigDefinition []ConfigSection       `yaml:"configDefinition"`
	MenuLocations   map[string]MenuLocation `yaml:"menuLocations"`
	WidgetAreas     []WidgetArea           `yaml:"widgetAreas"`
	CustomFields    map[string][]Field     `yaml:"customFields"`
	Hooks           []Hook                 `yaml:"hooks"`
	Settings        map[string]interface{} `yaml:"settings"`
}

// Capabilities 主题能力声明
type Capabilities struct {
	Layouts       []string `yaml:"layouts"`
	Modules       []string `yaml:"modules"`
	Widgets       []string `yaml:"widgets"`
	PostFormats   []string `yaml:"postFormats"`
	CustomFields  bool     `yaml:"customFields"`
	CustomTemplates bool   `yaml:"customTemplates"`
	ColorSchemes  bool     `yaml:"colorSchemes"`
	DarkMode      bool     `yaml:"darkMode"`
}

// Config 主题配置
type Config struct {
	Appearance      Appearance            `yaml:"appearance"`
	Modules         map[string]Module     `yaml:"modules"`
	Widgets         map[string]Widget     `yaml:"widgets"`
	Menus           map[string][]Menu     `yaml:"menus"`
	Pages           map[string]Page       `yaml:"pages"`
	Taxonomies      map[string]Taxonomy   `yaml:"taxonomies"`
	SEO             SEO                   `yaml:"seo"`
	Custom          Custom                `yaml:"custom"`
	Advanced        Advanced              `yaml:"advanced"`
}

// Appearance 外观配置
type Appearance struct {
	Colors       Colors       `yaml:"colors"`
	Typography   Typography   `yaml:"typography"`
	Layout       Layout       `yaml:"layout"`
	Animations   Animations   `yaml:"animations"`
}

// Colors 颜色配置
type Colors struct {
	Primary       string `yaml:"primary"`
	Secondary     string `yaml:"secondary"`
	Accent        string `yaml:"accent"`
	Background    string `yaml:"background"`
	Text          string `yaml:"text"`
	TextMuted     string `yaml:"textMuted"`
	Border        string `yaml:"border"`
	Success       string `yaml:"success"`
	Warning       string `yaml:"warning"`
	Error         string `yaml:"error"`
	Info          string `yaml:"info"`
}

// Typography 字体配置
type Typography struct {
	BaseFont   string     `yaml:"baseFont"`
	HeadingFont string    `yaml:"headingFont"`
	CodeFont   string     `yaml:"codeFont"`
	FontSize   FontSize   `yaml:"fontSize"`
	LineHeight LineHeight `yaml:"lineHeight"`
}

// FontSize 字号配置
type FontSize struct {
	Base  string  `yaml:"base"`
	Scale float64 `yaml:"scale"`
}

// LineHeight 行高配置
type LineHeight struct {
	Base    string `yaml:"base"`
	Heading string `yaml:"heading"`
}

// Layout 布局配置
type Layout struct {
	Width        string `yaml:"width"`
	Sidebar      string `yaml:"sidebar"`
	Container    string `yaml:"container"`
	Gap          string `yaml:"gap"`
	SectionGap   string `yaml:"sectionGap"`
	BorderRadius string `yaml:"borderRadius"`
}

// Animations 动画配置
type Animations struct {
	Enabled    string `yaml:"enabled"`
	Duration   string `yaml:"duration"`
	Easing     string `yaml:"easing"`
	FadeIn     bool   `yaml:"fadeIn"`
	SlideUp    bool   `yaml:"slideUp"`
	Stagger    string `yaml:"stagger"`
}

// Module 模块配置
type Module struct {
	Enabled bool                   `yaml:"enabled"`
	Config  map[string]interface{} `yaml:"config"`
}

// Widget 小工具配置
type Widget struct {
	Enabled bool                   `yaml:"enabled"`
	Area    string                 `yaml:"area"`
	Config  map[string]interface{} `yaml:"config"`
}

// Menu 菜单项
type Menu struct {
	Label    string `yaml:"label"`
	URL      string `yaml:"url"`
	Icon     string `yaml:"icon"`
	Target   string `yaml:"target"`
	Enabled  bool   `yaml:"enabled"`
	Children []Menu `yaml:"children,omitempty"`
}

// Page 页面配置
type Page struct {
	Title    string                 `yaml:"title"`
	Template string                 `yaml:"template"`
	Layout   string                 `yaml:"layout"`
	Enabled  bool                   `yaml:"enabled"`
	Config   map[string]interface{} `yaml:"config"`
}

// Taxonomy 分类法配置
type Taxonomy struct {
	Enabled        bool   `yaml:"enabled"`
	Template       string `yaml:"template"`
	ShowDescription bool   `yaml:"showDescription"`
	ShowPostCount  bool   `yaml:"showPostCount"`
	PerPage        int    `yaml:"perPage"`
}

// SEO SEO配置
type SEO struct {
	Title       string      `yaml:"title"`
	Description string      `yaml:"description"`
	Keywords    interface{} `yaml:"keywords"` // 支持字符串或数组
	Author      string      `yaml:"author"`
	Canonical   string      `yaml:"canonical"`
	Robots      string      `yaml:"robots"`
	OpenGraph   OpenGraph   `yaml:"opengraph"`
	Twitter     Twitter     `yaml:"twitter"`
	Schema      Schema      `yaml:"schema"`
	Sitemap     Sitemap     `yaml:"sitemap"`
}

// OpenGraph Open Graph配置
type OpenGraph struct {
	Enabled   string `yaml:"enabled"`
	Type      string `yaml:"type"`
	Image     string `yaml:"image"`
	ImageAlt  string `yaml:"imageAlt"`
}

// Twitter Twitter Card配置
type Twitter struct {
	Enabled  string `yaml:"enabled"`
	Card     string `yaml:"card"`
	Creator  string `yaml:"creator"`
	Site     string `yaml:"site"`
	Image    string `yaml:"image"`
}

// Schema 结构化数据配置
type Schema struct {
	Enabled string `yaml:"enabled"`
	Type    string `yaml:"type"`
	Logo    string `yaml:"logo"`
	Publisher string `yaml:"publisher"`
}

// Sitemap 站点地图配置
type Sitemap struct {
	Enabled           string `yaml:"enabled"`
	Priority          string `yaml:"priority"`
	Changefreq        string `yaml:"changefreq"`
	IncludeCategories bool   `yaml:"includeCategories"`
	IncludeTags       bool   `yaml:"includeTags"`
	IncludeAuthors    bool   `yaml:"includeAuthors"`
}

// Custom 自定义代码
type Custom struct {
	CSS       string `yaml:"css"`
	JS        string `yaml:"js"`
	Head      string `yaml:"head"`
	BodyTop   string `yaml:"bodyTop"`
	BodyBottom string `yaml:"bodyBottom"`
}

// Advanced 高级配置
type Advanced struct {
	Performance Performance `yaml:"performance"`
	Comments    Comments    `yaml:"comments"`
	Social      Social      `yaml:"social"`
	Analytics   Analytics   `yaml:"analytics"`
	Ads         Ads         `yaml:"ads"`
	Features    Features    `yaml:"features"`
}

// Performance 性能优化
type Performance struct {
	LazyLoadImages   bool     `yaml:"lazyLoadImages"`
	LazyLoadIframes  bool     `yaml:"lazyLoadIframes"`
	PreconnectUrls   []string `yaml:"preconnectUrls"`
	PreloadCritical  bool     `yaml:"preloadCritical"`
}

// Comments 评论系统
type Comments struct {
	Enabled  string                 `yaml:"enabled"`
	Provider string                 `yaml:"provider"`
	Disqus   map[string]interface{} `yaml:"disqus"`
	Valine   map[string]interface{} `yaml:"valine"`
}

// Social 社交媒体
type Social struct {
	Twitter string `yaml:"twitter"`
	Github  string `yaml:"github"`
	Weibo   string `yaml:"weibo"`
	Email   string `yaml:"email"`
	Rss     bool   `yaml:"rss"`
}

// Analytics 统计分析
type Analytics struct {
	Enabled  string                 `yaml:"enabled"`
	Provider string                 `yaml:"provider"`
	GoogleId string                 `yaml:"googleId"`
}

// Ads 广告配置
type Ads struct {
	Enabled  bool   `yaml:"enabled"`
	AutoAds  bool   `yaml:"autoAds"`
}

// Features 其他功能
type Features struct {
	ReadingProgress bool `yaml:"readingProgress"`
	BackToTop       bool `yaml:"backToTop"`
	Toc             bool `yaml:"toc"`
	WordCount       bool `yaml:"wordCount"`
	ReadingTime     bool `yaml:"readingTime"`
	PostViews       bool `yaml:"postViews"`
}

// ConfigSection 配置区块
type ConfigSection struct {
	Section     string  `yaml:"section" json:"section"`
	Icon        string  `yaml:"icon" json:"icon"`
	Description string  `yaml:"description" json:"description"`
	Fields      []Field `yaml:"fields" json:"fields"`
}

// Field 配置字段
type Field struct {
	Name        string      `yaml:"name" json:"name"`
	Label       string      `yaml:"label" json:"label"`
	Type        string      `yaml:"type" json:"type"`
	Default     interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Placeholder string      `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	Required    bool        `yaml:"required,omitempty" json:"required,omitempty"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	ShowIf      string      `yaml:"show_if,omitempty" json:"showIf,omitempty"`
	Options     []Option    `yaml:"options,omitempty" json:"options,omitempty"`
	Language    string      `yaml:"language,omitempty" json:"language,omitempty"`
	Min         int         `yaml:"min,omitempty" json:"min,omitempty"`
	Max         int         `yaml:"max,omitempty" json:"max,omitempty"`
	Step        float64     `yaml:"step,omitempty" json:"step,omitempty"`
}

// Option 选项
type Option struct {
	Label string      `yaml:"label" json:"label"`
	Value interface{} `yaml:"value" json:"value"`
}

// MenuLocation 菜单位置
type MenuLocation struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	MaxDepth    int    `yaml:"maxDepth"`
}

// WidgetArea 小工具区域
type WidgetArea struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	MaxWidgets  int    `yaml:"maxWidgets"`
	Columns     int    `yaml:"columns,omitempty"`
}

// Hook 主题钩子
type Hook struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ThemeConfigService 主题配置服务
type ThemeConfigService struct {
	themesPath string
}

// NewThemeConfigService 创建主题配置服务
func NewThemeConfigService(themesPath string) *ThemeConfigService {
	return &ThemeConfigService{
		themesPath: themesPath,
	}
}

// GetConfigPath 获取配置文件路径
func (s *ThemeConfigService) GetConfigPath(themeSlug string) string {
	return filepath.Join(s.themesPath, themeSlug, "theme.yml")
}

// LoadThemeConfig 加载主题配置
func (s *ThemeConfigService) LoadThemeConfig(themeSlug string) (*ThemeConfig, error) {
	configPath := s.GetConfigPath(themeSlug)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme config file: %w", err)
	}

	var config ThemeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse theme config: %w", err)
	}

	return &config, nil
}

// SaveThemeConfig 保存主题配置
func (s *ThemeConfigService) SaveThemeConfig(themeSlug string, config *ThemeConfig) error {
	configPath := s.GetConfigPath(themeSlug)

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal theme config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme config file: %w", err)
	}

	return nil
}

// UpdateThemeConfig 更新主题配置(支持嵌套路径)
func (s *ThemeConfigService) UpdateThemeConfig(themeSlug string, updates map[string]interface{}) error {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return err
	}

	// 应用更新
	for key, value := range updates {
		if err := s.setNestedValue(&config.Config, key, value); err != nil {
			return fmt.Errorf("failed to set value for %s: %w", key, err)
		}
	}

	return s.SaveThemeConfig(themeSlug, config)
}

// setNestedValue 设置嵌套值(支持点分隔的路径)
func (s *ThemeConfigService) setNestedValue(config *Config, path string, value interface{}) error {
	parts := strings.Split(path, ".")

	switch parts[0] {
	case "appearance":
		return s.setAppearanceValue(&config.Appearance, parts[1:], value)
	case "modules":
		return s.setModuleValue(config.Modules, parts[1:], value)
	case "widgets":
		return s.setWidgetValue(config.Widgets, parts[1:], value)
	case "menus":
		// 菜单配置暂时不支持通过此方法更新
		return fmt.Errorf("menu updates not supported via this method")
	case "seo":
		return s.setSEOValue(&config.SEO, parts[1:], value)
	case "custom":
		return s.setCustomValue(&config.Custom, parts[1:], value)
	default:
		return fmt.Errorf("unknown config section: %s", parts[0])
	}
}

// setAppearanceValue 设置外观配置值
func (s *ThemeConfigService) setAppearanceValue(appearance *Appearance, parts []string, value interface{}) error {
	if len(parts) == 0 {
		return fmt.Errorf("invalid appearance path")
	}

	switch parts[0] {
	case "colors":
		return s.setColorValue(&appearance.Colors, parts[1:], value)
	case "typography":
		return s.setTypographyValue(&appearance.Typography, parts[1:], value)
	case "layout":
		return s.setLayoutValue(&appearance.Layout, parts[1:], value)
	case "animations":
		return s.setAnimationsValue(&appearance.Animations, parts[1:], value)
	default:
		return fmt.Errorf("unknown appearance section: %s", parts[0])
	}
}

// setColorValue 设置颜色值
func (s *ThemeConfigService) setColorValue(colors *Colors, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid color path")
	}

	switch parts[0] {
	case "primary":
		colors.Primary = value.(string)
	case "secondary":
		colors.Secondary = value.(string)
	case "accent":
		colors.Accent = value.(string)
	case "background":
		colors.Background = value.(string)
	case "text":
		colors.Text = value.(string)
	case "textMuted":
		colors.TextMuted = value.(string)
	case "border":
		colors.Border = value.(string)
	case "success":
		colors.Success = value.(string)
	case "warning":
		colors.Warning = value.(string)
	case "error":
		colors.Error = value.(string)
	case "info":
		colors.Info = value.(string)
	default:
		return fmt.Errorf("unknown color: %s", parts[0])
	}

	return nil
}

// setTypographyValue 设置字体值
func (s *ThemeConfigService) setTypographyValue(typography *Typography, parts []string, value interface{}) error {
	if len(parts) == 0 {
		return fmt.Errorf("invalid typography path")
	}

	switch parts[0] {
	case "baseFont":
		typography.BaseFont = value.(string)
	case "headingFont":
		typography.HeadingFont = value.(string)
	case "codeFont":
		typography.CodeFont = value.(string)
	default:
		return fmt.Errorf("unknown typography field: %s", parts[0])
	}

	return nil
}

// setLayoutValue 设置布局值
func (s *ThemeConfigService) setLayoutValue(layout *Layout, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid layout path")
	}

	switch parts[0] {
	case "width":
		layout.Width = value.(string)
	case "sidebar":
		layout.Sidebar = value.(string)
	case "container":
		layout.Container = value.(string)
	case "gap":
		layout.Gap = value.(string)
	case "sectionGap":
		layout.SectionGap = value.(string)
	case "borderRadius":
		layout.BorderRadius = value.(string)
	default:
		return fmt.Errorf("unknown layout field: %s", parts[0])
	}

	return nil
}

// setAnimationsValue 设置动画值
func (s *ThemeConfigService) setAnimationsValue(animations *Animations, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid animations path")
	}

	switch parts[0] {
	case "enabled":
		animations.Enabled = fmt.Sprintf("%v", value)
	case "duration":
		animations.Duration = value.(string)
	case "easing":
		animations.Easing = value.(string)
	case "fadeIn":
		animations.FadeIn = value.(bool)
	case "slideUp":
		animations.SlideUp = value.(bool)
	case "stagger":
		animations.Stagger = value.(string)
	default:
		return fmt.Errorf("unknown animations field: %s", parts[0])
	}

	return nil
}

// setModuleValue 设置模块值
func (s *ThemeConfigService) setModuleValue(modules map[string]Module, parts []string, value interface{}) error {
	if len(parts) < 2 {
		return fmt.Errorf("invalid module path")
	}

	moduleID := parts[0]
	module, exists := modules[moduleID]
	if !exists {
		// 创建新模块
		module = Module{
			Enabled: false,
			Config:  make(map[string]interface{}),
		}
	}

	if parts[1] == "enabled" {
		module.Enabled = value.(bool)
		modules[moduleID] = module
		return nil
	}

	if parts[1] == "config" {
		if len(parts) < 3 {
			return fmt.Errorf("invalid module config path")
		}
		// 设置嵌套配置值
		config := module.Config
		for i := 2; i < len(parts)-1; i++ {
			key := parts[i]
			if _, exists := config[key]; !exists {
				config[key] = make(map[string]interface{})
			}
			config = config[key].(map[string]interface{})
		}
		config[parts[len(parts)-1]] = value
		modules[moduleID] = module
	}

	return nil
}

// setWidgetValue 设置小工具值
func (s *ThemeConfigService) setWidgetValue(widgets map[string]Widget, parts []string, value interface{}) error {
	if len(parts) < 2 {
		return fmt.Errorf("invalid widget path")
	}

	widgetID := parts[0]
	widget, exists := widgets[widgetID]
	if !exists {
		// 创建新小工具
		widget = Widget{
			Enabled: false,
			Area:    "sidebar",
			Config:  make(map[string]interface{}),
		}
	}

	switch parts[1] {
	case "enabled":
		widget.Enabled = value.(bool)
	case "area":
		widget.Area = value.(string)
	case "config":
		if len(parts) < 3 {
			return fmt.Errorf("invalid widget config path")
		}
		// 设置嵌套配置值
		config := widget.Config
		for i := 2; i < len(parts)-1; i++ {
			key := parts[i]
			if _, exists := config[key]; !exists {
				config[key] = make(map[string]interface{})
			}
			config = config[key].(map[string]interface{})
		}
		config[parts[len(parts)-1]] = value
	default:
		return fmt.Errorf("unknown widget field: %s", parts[1])
	}

	widgets[widgetID] = widget
	return nil
}

// setSEOValue 设置SEO值
func (s *ThemeConfigService) setSEOValue(seo *SEO, parts []string, value interface{}) error {
	if len(parts) == 0 {
		return fmt.Errorf("invalid seo path")
	}

	switch parts[0] {
	case "title":
		seo.Title = value.(string)
	case "description":
		seo.Description = value.(string)
	case "keywords":
		// 支持字符串或数组类型的 keywords
		switch v := value.(type) {
		case string:
			seo.Keywords = v
		case []interface{}:
			// 将 []interface{} 转换为 []string
			keywords := make([]string, 0, len(v))
			for _, item := range v {
				if str, ok := item.(string); ok {
					keywords = append(keywords, str)
				}
			}
			seo.Keywords = keywords
		case []string:
			seo.Keywords = v
		default:
			seo.Keywords = fmt.Sprintf("%v", v)
		}
	case "author":
		seo.Author = value.(string)
	case "canonical":
		seo.Canonical = value.(string)
	case "robots":
		seo.Robots = value.(string)
	default:
		if len(parts) > 1 {
			switch parts[0] {
			case "opengraph":
				return s.setOpenGraphValue(&seo.OpenGraph, parts[1:], value)
			case "twitter":
				return s.setTwitterValue(&seo.Twitter, parts[1:], value)
			case "schema":
				return s.setSchemaValue(&seo.Schema, parts[1:], value)
			case "sitemap":
				return s.setSitemapValue(&seo.Sitemap, parts[1:], value)
			}
		}
		return fmt.Errorf("unknown seo field: %s", parts[0])
	}

	return nil
}

// setOpenGraphValue 设置Open Graph值
func (s *ThemeConfigService) setOpenGraphValue(og *OpenGraph, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid opengraph path")
	}

	switch parts[0] {
	case "enabled":
		og.Enabled = fmt.Sprintf("%v", value)
	case "type":
		og.Type = value.(string)
	case "image":
		og.Image = value.(string)
	case "imageAlt":
		og.ImageAlt = value.(string)
	default:
		return fmt.Errorf("unknown opengraph field: %s", parts[0])
	}

	return nil
}

// setTwitterValue 设置Twitter值
func (s *ThemeConfigService) setTwitterValue(twitter *Twitter, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid twitter path")
	}

	switch parts[0] {
	case "enabled":
		twitter.Enabled = fmt.Sprintf("%v", value)
	case "card":
		twitter.Card = value.(string)
	case "creator":
		twitter.Creator = value.(string)
	case "site":
		twitter.Site = value.(string)
	case "image":
		twitter.Image = value.(string)
	default:
		return fmt.Errorf("unknown twitter field: %s", parts[0])
	}

	return nil
}

// setSchemaValue 设置Schema值
func (s *ThemeConfigService) setSchemaValue(schema *Schema, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid schema path")
	}

	switch parts[0] {
	case "enabled":
		schema.Enabled = fmt.Sprintf("%v", value)
	case "type":
		schema.Type = value.(string)
	case "logo":
		schema.Logo = value.(string)
	case "publisher":
		schema.Publisher = value.(string)
	default:
		return fmt.Errorf("unknown schema field: %s", parts[0])
	}

	return nil
}

// setSitemapValue 设置Sitemap值
func (s *ThemeConfigService) setSitemapValue(sitemap *Sitemap, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid sitemap path")
	}

	switch parts[0] {
	case "enabled":
		sitemap.Enabled = fmt.Sprintf("%v", value)
	case "priority":
		sitemap.Priority = value.(string)
	case "changefreq":
		sitemap.Changefreq = value.(string)
	case "includeCategories":
		sitemap.IncludeCategories = value.(bool)
	case "includeTags":
		sitemap.IncludeTags = value.(bool)
	case "includeAuthors":
		sitemap.IncludeAuthors = value.(bool)
	default:
		return fmt.Errorf("unknown sitemap field: %s", parts[0])
	}

	return nil
}

// setCustomValue 设置自定义代码值
func (s *ThemeConfigService) setCustomValue(custom *Custom, parts []string, value interface{}) error {
	if len(parts) != 1 {
		return fmt.Errorf("invalid custom path")
	}

	switch parts[0] {
	case "css":
		custom.CSS = value.(string)
	case "js":
		custom.JS = value.(string)
	case "head":
		custom.Head = value.(string)
	case "bodyTop":
		custom.BodyTop = value.(string)
	case "bodyBottom":
		custom.BodyBottom = value.(string)
	default:
		return fmt.Errorf("unknown custom field: %s", parts[0])
	}

	return nil
}

// GetModuleConfig 获取模块配置
func (s *ThemeConfigService) GetModuleConfig(themeSlug string, moduleID string) (map[string]interface{}, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return nil, err
	}

	module, exists := config.Config.Modules[moduleID]
	if !exists {
		return nil, fmt.Errorf("module %s not found", moduleID)
	}

	return module.Config, nil
}

// SetModuleConfig 设置模块配置
func (s *ThemeConfigService) SetModuleConfig(themeSlug string, moduleID string, moduleConfig Module) error {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return err
	}

	config.Config.Modules[moduleID] = moduleConfig
	return s.SaveThemeConfig(themeSlug, config)
}

// GetWidgetConfig 获取小工具配置
func (s *ThemeConfigService) GetWidgetConfig(themeSlug string, widgetID string) (Widget, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return Widget{}, err
	}

	widget, exists := config.Config.Widgets[widgetID]
	if !exists {
		return Widget{}, fmt.Errorf("widget %s not found", widgetID)
	}

	return widget, nil
}

// SetWidgetConfig 设置小工具配置
func (s *ThemeConfigService) SetWidgetConfig(themeSlug string, widgetID string, widget Widget) error {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return err
	}

	config.Config.Widgets[widgetID] = widget
	return s.SaveThemeConfig(themeSlug, config)
}

// GetAllModules 获取所有模块
func (s *ThemeConfigService) GetAllModules(themeSlug string) (map[string]Module, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return nil, err
	}

	return config.Config.Modules, nil
}

// GetAllWidgets 获取所有小工具
func (s *ThemeConfigService) GetAllWidgets(themeSlug string) (map[string]Widget, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return nil, err
	}

	return config.Config.Widgets, nil
}

// GetEnabledModules 获取启用的模块
func (s *ThemeConfigService) GetEnabledModules(themeSlug string) ([]string, error) {
	modules, err := s.GetAllModules(themeSlug)
	if err != nil {
		return nil, err
	}

	var enabled []string
	for id, module := range modules {
		if module.Enabled {
			enabled = append(enabled, id)
		}
	}

	return enabled, nil
}

// GetEnabledWidgetsByArea 获取指定区域启用的小工具
func (s *ThemeConfigService) GetEnabledWidgetsByArea(themeSlug string, area string) ([]string, error) {
	widgets, err := s.GetAllWidgets(themeSlug)
	if err != nil {
		return nil, err
	}

	var enabled []string
	for id, widget := range widgets {
		if widget.Enabled && widget.Area == area {
			enabled = append(enabled, id)
		}
	}

	return enabled, nil
}

// EnableModule 启用模块
func (s *ThemeConfigService) EnableModule(themeSlug string, moduleID string) error {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return err
	}

	module, exists := config.Config.Modules[moduleID]
	if !exists {
		// 创建新模块
		module = Module{
			Enabled: true,
			Config:  make(map[string]interface{}),
		}
	} else {
		module.Enabled = true
	}

	config.Config.Modules[moduleID] = module
	return s.SaveThemeConfig(themeSlug, config)
}

// DisableModule 禁用模块
func (s *ThemeConfigService) DisableModule(themeSlug string, moduleID string) error {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return err
	}

	module, exists := config.Config.Modules[moduleID]
	if !exists {
		return fmt.Errorf("module %s not found", moduleID)
	}

	module.Enabled = false
	config.Config.Modules[moduleID] = module
	return s.SaveThemeConfig(themeSlug, config)
}

// ReorderModules 重新排序模块
func (s *ThemeConfigService) ReorderModules(themeSlug string, order []string) error {
	// 注意: 当前配置结构不支持直接排序,需要在渲染时处理
	// 这里返回成功,实际排序逻辑在前端或渲染时处理
	return nil
}

// GetWidgetAreas 获取小工具区域
func (s *ThemeConfigService) GetWidgetAreas(themeSlug string) ([]WidgetArea, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return nil, err
	}

	return config.WidgetAreas, nil
}

// GetMenuLocations 获取菜单位置
func (s *ThemeConfigService) GetMenuLocations(themeSlug string) (map[string]MenuLocation, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return nil, err
	}

	return config.MenuLocations, nil
}

// GetConfigForm 获取配置表单定义
func (s *ThemeConfigService) GetConfigForm(themeSlug string) ([]ConfigSection, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return nil, err
	}

	return config.ConfigDefinition, nil
}

// GetCustomFields 获取自定义字段定义
func (s *ThemeConfigService) GetCustomFields(themeSlug string, contentType string) ([]Field, error) {
	config, err := s.LoadThemeConfig(themeSlug)
	if err != nil {
		return nil, err
	}

	fields, exists := config.CustomFields[contentType]
	if !exists {
		return []Field{}, nil
	}

	return fields, nil
}

// ListAllThemes 列出所有主题
func (s *ThemeConfigService) ListAllThemes() ([]string, error) {
	entries, err := os.ReadDir(s.themesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	var themes []string
	for _, entry := range entries {
		if entry.IsDir() {
			configPath := filepath.Join(s.themesPath, entry.Name(), "theme.yml")
			if _, err := os.Stat(configPath); err == nil {
				themes = append(themes, entry.Name())
			}
		}
	}

	return themes, nil
}

// ThemeExists 检查主题是否存在
func (s *ThemeConfigService) ThemeExists(themeSlug string) bool {
	configPath := s.GetConfigPath(themeSlug)
	_, err := os.Stat(configPath)
	return err == nil
}
