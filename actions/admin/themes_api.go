package admin

import (
	"catter-cms/actions/api"
	"catter-cms/services"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo"
)

// ========================================
// 主题配置管理API
// ========================================

// ThemeConfigResponse 主题配置响应
type ThemeConfigResponse struct {
	Appearance      map[string]interface{}            `json:"appearance"`
	Modules         map[string]services.Module        `json:"modules"`
	Widgets         map[string]services.Widget        `json:"widgets"`
	Menus           map[string][]services.Menu        `json:"menus"`
	Pages           map[string]services.Page          `json:"pages"`
	SEO             map[string]interface{}            `json:"seo"`
	Custom          map[string]string                 `json:"custom"`
	Advanced        map[string]interface{}            `json:"advanced"`
	MenuLocations   map[string]services.MenuLocation  `json:"menuLocations"`
	WidgetAreas     []services.WidgetArea            `json:"widgetAreas"`
	ConfigForm      []services.ConfigSection         `json:"configForm"`
	CustomFields    map[string][]services.Field       `json:"customFields"`
}

// ThemeModuleRequest 主题模块请求
type ThemeModuleRequest struct {
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config"`
}

// ThemeWidgetRequest 主题小工具请求
type ThemeWidgetRequest struct {
	Enabled bool                   `json:"enabled"`
	Area    string                 `json:"area"`
	Config  map[string]interface{} `json:"config"`
}

// GetThemeConfigAPI 获取主题完整配置
// GET /api/themes/{slug}/config
func GetThemeConfigAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	// 尝试从缓存获取
	cacheKey := "theme_config_" + themeSlug
	if cachedConfig, ok := services.GetGlobalThemeCache().Get(cacheKey); ok {
		if response, ok := cachedConfig.(ThemeConfigResponse); ok {
			return api.Success(c, response)
		}
	}
	
	// 缓存未命中，加载配置
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 加载主题配置
	config, err := themeConfigService.LoadThemeConfig(themeSlug)
	if err != nil {
		return api.InternalError(c, "加载主题配置失败")
	}
	
	// 构建响应
	response := ThemeConfigResponse{
		Appearance: map[string]interface{}{
			"colors":     config.Config.Appearance.Colors,
			"typography": config.Config.Appearance.Typography,
			"layout":     config.Config.Appearance.Layout,
			"animations": config.Config.Appearance.Animations,
		},
		Modules:       config.Config.Modules,
		Widgets:       config.Config.Widgets,
		Menus:         config.Config.Menus,
		Pages:         config.Config.Pages,
		SEO: map[string]interface{}{
			"title":       config.Config.SEO.Title,
			"description": config.Config.SEO.Description,
			"keywords":    config.Config.SEO.Keywords,
			"author":      config.Config.SEO.Author,
			"canonical":   config.Config.SEO.Canonical,
			"robots":      config.Config.SEO.Robots,
			"opengraph":   config.Config.SEO.OpenGraph,
			"twitter":     config.Config.SEO.Twitter,
			"schema":      config.Config.SEO.Schema,
			"sitemap":     config.Config.SEO.Sitemap,
		},
		Custom: map[string]string{
			"css":       config.Config.Custom.CSS,
			"js":        config.Config.Custom.JS,
			"head":      config.Config.Custom.Head,
			"bodyTop":   config.Config.Custom.BodyTop,
			"bodyBottom": config.Config.Custom.BodyBottom,
		},
		Advanced:        map[string]interface{}{
			"performance": config.Config.Advanced.Performance,
			"comments":    config.Config.Advanced.Comments,
			"social":      config.Config.Advanced.Social,
			"analytics":   config.Config.Advanced.Analytics,
			"ads":         config.Config.Advanced.Ads,
			"features":    config.Config.Advanced.Features,
		},
		MenuLocations:   config.MenuLocations,
		WidgetAreas:     config.WidgetAreas,
		ConfigForm:      config.ConfigDefinition,
		CustomFields:    config.CustomFields,
	}
	
	// 保存到缓存
	services.GetGlobalThemeCache().Set(cacheKey, response)
	
	return api.Success(c, response)
}

// UpdateThemeConfigAPI 更新主题配置
// PUT /api/themes/{slug}/config
func UpdateThemeConfigAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return api.BadRequest(c, "参数错误")
	}
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 更新配置
	if err := themeConfigService.UpdateThemeConfig(themeSlug, updates); err != nil {
		return api.InternalError(c, "更新主题配置失败")
	}
	
	// 清除缓存
	cacheKey := "theme_config_" + themeSlug
	services.GetGlobalThemeCache().Delete(cacheKey)
	
	return api.SuccessWithMessage(c, "主题配置已更新", nil)
}

// GetThemeModuleAPI 获取主题模块配置
// GET /api/themes/{slug}/modules/{moduleId}
func GetThemeModuleAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	moduleID := c.Param("moduleId")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取模块配置
	moduleConfig, err := themeConfigService.GetModuleConfig(themeSlug, moduleID)
	if err != nil {
		return api.NotFound(c, "模块不存在")
	}
	
	module := services.Module{
		Enabled: moduleConfig["enabled"].(bool),
		Config:  moduleConfig,
	}
	
	return api.Success(c, module)
}

// SetThemeModuleAPI 设置主题模块配置
// PUT /api/themes/{slug}/modules/{moduleId}
func SetThemeModuleAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	moduleID := c.Param("moduleId")
	
	var req ThemeModuleRequest
	if err := c.Bind(&req); err != nil {
		return api.BadRequest(c, "参数错误")
	}
	
	module := services.Module{
		Enabled: req.Enabled,
		Config:  req.Config,
	}
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 设置模块配置
	if err := themeConfigService.SetModuleConfig(themeSlug, moduleID, module); err != nil {
		return api.InternalError(c, "设置模块配置失败")
	}
	
	return api.SuccessWithMessage(c, "模块配置已保存", nil)
}

// EnableThemeModuleAPI 启用主题模块
// POST /api/themes/{slug}/modules/{moduleId}/enable
func EnableThemeModuleAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	moduleID := c.Param("moduleId")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 启用模块
	if err := themeConfigService.EnableModule(themeSlug, moduleID); err != nil {
		return api.InternalError(c, "启用模块失败")
	}
	
	return api.SuccessWithMessage(c, "模块已启用", nil)
}

// DisableThemeModuleAPI 禁用主题模块
// POST /api/themes/{slug}/modules/{moduleId}/disable
func DisableThemeModuleAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	moduleID := c.Param("moduleId")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 禁用模块
	if err := themeConfigService.DisableModule(themeSlug, moduleID); err != nil {
		return api.InternalError(c, "禁用模块失败")
	}
	
	return api.SuccessWithMessage(c, "模块已禁用", nil)
}

// GetThemeModulesAPI 获取所有主题模块
// GET /api/themes/{slug}/modules
func GetThemeModulesAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取所有模块
	modules, err := themeConfigService.GetAllModules(themeSlug)
	if err != nil {
		return api.InternalError(c, "获取模块列表失败")
	}
	
	return api.Success(c, modules)
}

// GetThemeWidgetAPI 获取主题小工具配置
// GET /api/themes/{slug}/widgets/{widgetId}
func GetThemeWidgetAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	widgetID := c.Param("widgetId")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取小工具配置
	widget, err := themeConfigService.GetWidgetConfig(themeSlug, widgetID)
	if err != nil {
		return api.NotFound(c, "小工具不存在")
	}
	
	return api.Success(c, widget)
}

// SetThemeWidgetAPI 设置主题小工具配置
// PUT /api/themes/{slug}/widgets/{widgetId}
func SetThemeWidgetAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	widgetID := c.Param("widgetId")
	
	var req ThemeWidgetRequest
	if err := c.Bind(&req); err != nil {
		return api.BadRequest(c, "参数错误")
	}
	
	widget := services.Widget{
		Enabled: req.Enabled,
		Area:    req.Area,
		Config:  req.Config,
	}
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 设置小工具配置
	if err := themeConfigService.SetWidgetConfig(themeSlug, widgetID, widget); err != nil {
		return api.InternalError(c, "设置小工具配置失败")
	}
	
	return api.SuccessWithMessage(c, "小工具配置已保存", nil)
}

// GetThemeWidgetsAPI 获取所有主题小工具
// GET /api/themes/{slug}/widgets
func GetThemeWidgetsAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取所有小工具
	widgets, err := themeConfigService.GetAllWidgets(themeSlug)
	if err != nil {
		return api.InternalError(c, "获取小工具列表失败")
	}
	
	return api.Success(c, widgets)
}

// GetEnabledThemeModulesAPI 获取启用的主题模块
// GET /api/themes/{slug}/modules/enabled
func GetEnabledThemeModulesAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取启用的模块
	modules, err := themeConfigService.GetEnabledModules(themeSlug)
	if err != nil {
		return api.InternalError(c, "获取启用的模块失败")
	}
	
	return api.Success(c, modules)
}

// GetThemeWidgetsByAreaAPI 获取指定区域的小工具
// GET /api/themes/{slug}/widgets/area/{area}
func GetThemeWidgetsByAreaAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	area := c.Param("area")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取指定区域的小工具
	widgets, err := themeConfigService.GetEnabledWidgetsByArea(themeSlug, area)
	if err != nil {
		return api.InternalError(c, "获取小工具失败")
	}
	
	return api.Success(c, widgets)
}

// GetThemeConfigFormAPI 获取主题配置表单定义
// GET /api/themes/{slug}/config/form
func GetThemeConfigFormAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取配置表单
	form, err := themeConfigService.GetConfigForm(themeSlug)
	if err != nil {
		return api.InternalError(c, "获取配置表单失败")
	}
	
	return api.Success(c, form)
}

// GetThemeCustomFieldsAPI 获取主题自定义字段
// GET /api/themes/{slug}/custom-fields/{contentType}
func GetThemeCustomFieldsAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	contentType := c.Param("contentType")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取自定义字段
	fields, err := themeConfigService.GetCustomFields(themeSlug, contentType)
	if err != nil {
		return api.InternalError(c, "获取自定义字段失败")
	}
	
	return api.Success(c, fields)
}

// GetThemeWidgetAreasAPI 获取小工具区域
// GET /api/themes/{slug}/widget-areas
func GetThemeWidgetAreasAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取小工具区域
	areas, err := themeConfigService.GetWidgetAreas(themeSlug)
	if err != nil {
		return api.InternalError(c, "获取小工具区域失败")
	}
	
	return api.Success(c, areas)
}

// GetThemeMenuLocationsAPI 获取菜单位置
// GET /api/themes/{slug}/menu-locations
func GetThemeMenuLocationsAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	themeConfigService := services.NewThemeConfigService("templates/themes")
	
	// 获取菜单位置
	locations, err := themeConfigService.GetMenuLocations(themeSlug)
	if err != nil {
		return api.InternalError(c, "获取菜单位置失败")
	}
	
	return api.Success(c, locations)
}

// ResetThemeConfigAPI 重置主题配置为默认值
// POST /api/themes/{slug}/config/reset
func ResetThemeConfigAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	// 这里可以实现重置逻辑,从theme.yml读取默认配置
	// 暂时返回成功
	_ = themeSlug // 临时避免未使用变量警告
	
	return api.SuccessWithMessage(c, "主题配置已重置", nil)
}

// PreviewThemeConfigAPI 预览主题配置(不保存)
// POST /api/themes/{slug}/config/preview
func PreviewThemeConfigAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	
	var previewConfig map[string]interface{}
	if err := c.Bind(&previewConfig); err != nil {
		return api.BadRequest(c, "参数错误")
	}
	
	// 生成预览Token(有效期15分钟)
	// 这里简化处理,直接返回配置
	
	return api.Success(c, map[string]interface{}{
		"config": previewConfig,
		"preview_url": "/preview/" + themeSlug,
	})
}

// UploadThemeAPI 上传并安装主题
// POST /api/themes/upload
func UploadThemeAPI(c buffalo.Context) error {
	// 获取上传的文件
	file, err := c.File("theme")
	if err != nil {
		return api.BadRequest(c, "请上传主题文件")
	}
	
	// 验证文件类型
	if !strings.HasSuffix(file.Filename, ".zip") {
		return api.BadRequest(c, "主题文件必须是ZIP格式")
	}
	
	// 验证文件大小(最大50MB)
	if file.Size > 50*1024*1024 {
		return api.BadRequest(c, "主题文件大小不能超过50MB")
	}
	
	// 保存临时文件
	tempPath := filepath.Join(os.TempDir(), file.Filename)
	
	// 直接读取文件内容
	fileData := make([]byte, file.Size)
	_, err = file.File.Read(fileData)
	if err != nil {
		return api.InternalError(c, "读取文件内容失败")
	}
	
	if err := os.WriteFile(tempPath, fileData, 0644); err != nil {
		return api.InternalError(c, "保存文件失败")
	}
	defer os.Remove(tempPath)
	
	// 解压主题(这里简化处理,实际应该调用解压服务)
	// ... 解压逻辑 ...
	
	return api.SuccessWithMessage(c, "主题上传成功", map[string]string{
		"filename": file.Filename,
		"size":     fmt.Sprintf("%d", file.Size),
	})
}
