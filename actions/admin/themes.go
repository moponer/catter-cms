package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"catter-cms/services"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// ThemesListAPI 获取主题列表API
func ThemesListAPI(c buffalo.Context) error {
	themes := &models.Themes{}
	query := models.DB.Order("created_at DESC")

	if err := query.All(themes); err != nil {
		return api.InternalError(c, "获取主题列表失败")
	}

	// 为每个主题添加是否有模板文件的信息
	result := make([]map[string]interface{}, len(*themes))
	configService := services.NewThemeConfigService("templates/themes")
	
	for i, theme := range *themes {
		// 总是尝试从theme.yml读取最新的缩略图配置
		thumbnailPath := theme.Thumbnail
		if config, err := configService.LoadThemeConfig(theme.Slug); err == nil && config.Thumbnail != "" {
			thumbnailPath = "/themes/" + theme.Slug + "/" + config.Thumbnail
		}
		
		result[i] = map[string]interface{}{
			"id":          theme.ID,
			"name":        theme.Name,
			"slug":        theme.Slug,
			"version":     theme.Version,
			"description": theme.Description,
			"author":      theme.Author,
			"thumbnail":   thumbnailPath,
			"is_active":   theme.IsActive,
			"created_at":  theme.CreatedAt,
			"updated_at":  theme.UpdatedAt,
			"has_templates": services.ThemeHasTemplates(theme.Slug),
		}
	}

	return api.Success(c, result)
}

// ThemeDetailAPI 获取主题详情API
func ThemeDetailAPI(c buffalo.Context) error {
	themeID := c.Param("id")
	theme := &models.Theme{}

	if err := models.DB.Find(theme, themeID); err != nil {
		return api.NotFound(c, "主题不存在")
	}

	return api.Success(c, theme)
}

// ThemesActivateAPI 激活主题API
func ThemesActivateAPI(c buffalo.Context) error {
	themeID := c.Param("id")
	theme := &models.Theme{}

	if err := models.DB.Find(theme, themeID); err != nil {
		return api.NotFound(c, "主题不存在")
	}

	// 检查主题是否有模板文件
	if !services.ThemeHasTemplates(theme.Slug) {
		return api.BadRequest(c, "该主题缺少模板文件，无法激活。请确保主题目录包含 pages/home.plush.html 等模板文件。")
	}

	// 先取消所有主题的激活状态
	allThemes := &models.Themes{}
	if err := models.DB.All(allThemes); err == nil {
		for _, t := range *allThemes {
			if t.ID != theme.ID {
				t.IsActive = false
				models.DB.Update(&t)
			}
		}
	}

	// 激活当前主题
	theme.IsActive = true
	if err := models.DB.Update(theme); err != nil {
		return api.InternalError(c, "激活主题失败")
	}

	// 更新系统设置
	setting := &models.Setting{}
	if err := models.DB.Where("key = ?", "active_theme").First(setting); err != nil {
		// 创建新设置
		setting = &models.Setting{
			Key:      "active_theme",
			Value:    theme.Slug,
			Type:     "string",
			Category: "theme",
		}
		models.DB.Create(setting)
	} else {
		setting.Value = theme.Slug
		models.DB.Update(setting)
	}

	// 清除主题缓存，确保前台立即使用新主题
	services.ClearThemeCache()

	return api.SuccessWithMessage(c, "主题激活成功", theme)
}

// ThemesGetConfigAPI 获取主题配置API
func ThemesGetConfigAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")

	// 获取主题信息
	theme := &models.Theme{}
	if err := models.DB.Where("slug = ?", themeSlug).First(theme); err != nil {
		return api.NotFound(c, "主题不存在")
	}

	// 加载主题配置定义
	themeConfigService := services.NewThemeConfigService("templates/themes")
	config, err := themeConfigService.LoadThemeConfig(themeSlug)
	if err != nil {
		// 如果配置文件不存在，返回基本配置字段
		return api.Success(c, map[string]interface{}{
			"theme": theme,
			"config": map[string]interface{}{
				"name":    theme.Name,
				"slug":    theme.Slug,
				"version": "1.0.0",
				"fields":  getDefaultConfigFields(),
				"values":  getDefaultConfigValues(),
			},
		})
	}

	// 转换为前端期望的格式
	fields := convertConfigToFields(config)
	values := extractConfigValues(config)

	return api.Success(c, map[string]interface{}{
		"theme":  theme,
		"config": map[string]interface{}{
			"name":    config.Name,
			"version": config.Version,
			"fields":  fields,
			"values":  values,
		},
	})
}

// getDefaultConfigFields 返回默认配置字段
func getDefaultConfigFields() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "primary_color", "label": "主色调", "type": "color"},
		{"name": "secondary_color", "label": "次要颜色", "type": "color"},
		{"name": "background_color", "label": "背景颜色", "type": "color"},
		{"name": "text_color", "label": "文字颜色", "type": "color"},
		{"name": "seo_title", "label": "SEO 标题", "type": "text"},
		{"name": "seo_description", "label": "SEO 描述", "type": "textarea"},
		{"name": "seo_keywords", "label": "SEO 关键词", "type": "text"},
		{"name": "custom_css", "label": "自定义 CSS", "type": "textarea"},
		{"name": "custom_js", "label": "自定义 JavaScript", "type": "textarea"},
	}
}

// getDefaultConfigValues 返回默认配置值
func getDefaultConfigValues() map[string]interface{} {
	return map[string]interface{}{
		"primary_color":    "#3B82F6",
		"secondary_color":  "#764ba2",
		"background_color": "#f5f7fa",
		"text_color":       "#1f2937",
		"seo_title":        "",
		"seo_description":  "",
		"seo_keywords":     "",
		"custom_css":       "",
		"custom_js":        "",
	}
}

// convertConfigToFields 将主题配置转换为前端字段格式
func convertConfigToFields(config *services.ThemeConfig) []map[string]interface{} {
	fields := []map[string]interface{}{
		{"name": "primary_color", "label": "主色调", "type": "color"},
		{"name": "secondary_color", "label": "次要颜色", "type": "color"},
		{"name": "background_color", "label": "背景颜色", "type": "color"},
		{"name": "text_color", "label": "文字颜色", "type": "color"},
		{"name": "border_color", "label": "边框颜色", "type": "color"},
		{"name": "seo_title", "label": "SEO 标题", "type": "text", "placeholder": "网站标题"},
		{"name": "seo_description", "label": "SEO 描述", "type": "textarea", "placeholder": "网站描述"},
		{"name": "seo_keywords", "label": "SEO 关键词", "type": "text", "placeholder": "关键词，用逗号分隔"},
		{"name": "custom_css", "label": "自定义 CSS", "type": "textarea", "placeholder": "输入自定义CSS样式"},
		{"name": "custom_js", "label": "自定义 JavaScript", "type": "textarea", "placeholder": "输入自定义JavaScript代码"},
	}
	return fields
}

// extractConfigValues 从主题配置中提取值
func extractConfigValues(config *services.ThemeConfig) map[string]interface{} {
	return map[string]interface{}{
		"primary_color":     config.Config.Appearance.Colors.Primary,
		"secondary_color":   config.Config.Appearance.Colors.Secondary,
		"background_color":  config.Config.Appearance.Colors.Background,
		"text_color":        config.Config.Appearance.Colors.Text,
		"border_color":      config.Config.Appearance.Colors.Border,
		"seo_title":         config.Config.SEO.Title,
		"seo_description":   config.Config.SEO.Description,
		"seo_keywords":      config.Config.SEO.Keywords,
		"custom_css":        config.Config.Custom.CSS,
		"custom_js":         config.Config.Custom.JS,
	}
}

// ThemesSaveConfigAPI 保存主题配置API
func ThemesSaveConfigAPI(c buffalo.Context) error {
	themeSlug := c.Param("slug")
	configData := map[string]interface{}{}

	if err := c.Bind(&configData); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	themeConfigService := services.NewThemeConfigService("templates/themes")
	if err := themeConfigService.UpdateThemeConfig(themeSlug, configData); err != nil {
		return api.InternalError(c, "保存主题配置失败")
	}

	return api.SuccessWithMessage(c, "主题配置保存成功", nil)
}

// ThemesIndexHandler 主题列表页面 (已弃用，SPA模式)
func ThemesIndexHandler(c buffalo.Context) error {
	return c.Render(200, render.HTML("<html><body><p>Themes SPA</p></body></html>"))
}

// ThemesNewHandler 新建主题页面 (已弃用)
func ThemesNewHandler(c buffalo.Context) error {
	return c.Render(200, render.HTML("<html><body><p>Themes SPA</p></body></html>"))
}

// ThemesCreateHandler 创建主题 (已弃用，使用API)
func ThemesCreateHandler(c buffalo.Context) error {
	theme := &models.Theme{}
	if err := c.Bind(theme); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	if err := models.DB.Create(theme); err != nil {
		return api.InternalError(c, "创建主题失败")
	}

	return api.SuccessWithMessage(c, "主题创建成功", theme)
}

// ThemesActivateHandler 激活主题 (兼容旧API)
func ThemesActivateHandler(c buffalo.Context) error {
	return ThemesActivateAPI(c)
}

// ThemesDeactivateHandler 禁用主题
func ThemesDeactivateHandler(c buffalo.Context) error {
	themeID := c.Param("id")
	theme := &models.Theme{}

	if err := models.DB.Find(theme, themeID); err != nil {
		return api.NotFound(c, "主题不存在")
	}

	theme.IsActive = false
	if err := models.DB.Update(theme); err != nil {
		return api.InternalError(c, "禁用主题失败")
	}

	return api.SuccessWithMessage(c, "主题禁用成功", nil)
}

// ThemesDeleteHandler 删除主题
func ThemesDeleteHandler(c buffalo.Context) error {
	themeID := c.Param("id")
	theme := &models.Theme{}

	if err := models.DB.Find(theme, themeID); err != nil {
		return api.NotFound(c, "主题不存在")
	}

	if err := models.DB.Destroy(theme); err != nil {
		return api.InternalError(c, "删除主题失败")
	}

	return api.SuccessWithMessage(c, "主题删除成功", nil)
}

// ThemesConfigHandler 获取主题配置 (兼容旧API)
func ThemesConfigHandler(c buffalo.Context) error {
	return ThemesGetConfigAPI(c)
}

// ThemesSaveConfigHandler 保存主题配置 (兼容旧API)
func ThemesSaveConfigHandler(c buffalo.Context) error {
	return ThemesSaveConfigAPI(c)
}

// unused - avoid compile error
var _ = strconv.Itoa
