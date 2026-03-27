package api

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/services"
)

// ThemeInfo 主题信息（用于 API 响应）
type ThemeInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Thumbnail   string `json:"thumbnail"`
	IsActive    bool   `json:"is_active"`
}

// ListThemes 获取主题列表
func ListThemes(c *gin.Context) {
	themes := services.GetAllThemes()
	activeTheme := services.GetActiveTheme()

	var themeList []ThemeInfo
	id := uint(1)
	for slug, theme := range themes {
		info := ThemeInfo{
			ID:          id,
			Name:        theme.Name,
			Slug:        slug,
			Version:     theme.Version,
			Description: theme.Description,
			Author:      theme.Author,
			Thumbnail:   theme.Thumbnail,
			IsActive:    activeTheme != nil && activeTheme.Name == theme.Name,
		}
		themeList = append(themeList, info)
		id++
	}

	SuccessResponse(c, gin.H{
		"list":  themeList,
		"total": len(themeList),
	})
}

// GetTheme 获取单个主题
func GetTheme(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	info := ThemeInfo{
		ID:          0,
		Name:        theme.Name,
		Slug:        slug,
		Version:     theme.Version,
		Description: theme.Description,
		Author:      theme.Author,
		Thumbnail:   theme.Thumbnail,
	}

	SuccessResponse(c, info)
}

// GetThemeConfig 获取主题配置
func GetThemeConfig(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 获取完整配置（包含数据库中的配置）
	config, err := services.GetThemeFullConfig(slug)
	if err != nil {
		// 如果获取失败，返回内存中的配置
		config = theme.Config
	}

	SuccessResponse(c, config)
}

// SaveThemeConfig 保存主题配置
func SaveThemeConfig(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 解析配置
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 保存配置到数据库
	if err := services.UpdateThemeConfig(slug, config); err != nil {
		InternalError(c, "保存配置失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "配置保存成功", config)
}

// ThemeConfigFull 主题完整配置结构
type ThemeConfigFull struct {
	Name        string                   `json:"name"`
	Slug        string                   `json:"slug"`
	Version     string                   `json:"version"`
	Description string                   `json:"description"`
	Author      string                   `json:"author"`
	ConfigForm  []services.ConfigSection `json:"configForm"`
}

// GetThemeConfigFull 获取主题完整配置（包含配置表单定义和已保存的配置值）
func GetThemeConfigFull(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 直接使用已解析的 ConfigDefinition
	configForm := theme.ConfigDefinition
	if configForm == nil {
		configForm = []services.ConfigSection{}
	}

	// 获取已保存的配置值
	savedConfig, err := services.GetThemeFullConfig(slug)
	if err != nil {
		savedConfig = make(map[string]interface{})
	}

	// 构建响应结果，包含配置表单定义和已保存的配置值
	result := map[string]interface{}{
		"name":        theme.Name,
		"slug":        slug,
		"version":     theme.Version,
		"description": theme.Description,
		"author":      theme.Author,
		"configForm":  configForm,
	}

	// 将已保存的配置值合并到结果中
	for key, value := range savedConfig {
		result[key] = value
	}

	SuccessResponse(c, result)
}

// SaveThemeConfigFull 保存主题完整配置
func SaveThemeConfigFull(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 解析配置
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 保存配置到数据库
	if err := services.UpdateThemeConfig(slug, config); err != nil {
		InternalError(c, "保存配置失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "配置保存成功", config)
}

// ResetThemeConfig 重置主题配置
func ResetThemeConfig(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 重置为空配置
	theme.Config = make(map[string]interface{})

	// 同时清空数据库中的配置
	if err := services.UpdateThemeConfig(slug, make(map[string]interface{})); err != nil {
		InternalError(c, "重置配置失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "配置已重置", nil)
}

// ActivateTheme 激活主题
func ActivateTheme(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 设置活动主题（会同时保存到数据库）
	if err := services.SetActiveTheme(slug); err != nil {
		InternalError(c, "激活主题失败")
		return
	}

	SuccessWithMessage(c, "主题已激活", map[string]string{
		"active_theme": slug,
	})
}

// DeleteTheme 删除主题
func DeleteTheme(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 检查主题是否激活
	activeTheme := services.GetActiveTheme()
	if activeTheme != nil && activeTheme.Name == slug {
		BadRequest(c, "不能删除当前激活的主题")
		return
	}

	// 获取主题目录路径
	cfg := config.Get()
	themesPath := cfg.Server.ThemePath
	if themesPath == "" {
		themesPath = "../templates/themes"
	}
	themeDir := filepath.Join(themesPath, slug)

	// 删除主题目录
	if _, err := os.Stat(themeDir); err == nil {
		// 目录存在，尝试删除
		if err := os.RemoveAll(themeDir); err != nil {
			InternalError(c, "删除主题文件失败: "+err.Error())
			return
		}
	}

	// 删除数据库中的主题配置
	if err := services.DeleteThemeConfig(slug); err != nil {
		// 配置删除失败但文件已删除，记录警告但不返回错误
		log.Printf("警告：删除主题 %s 配置失败: %v", slug, err)
	}

	// 重新加载主题服务
	if err := services.InitThemeService(); err != nil {
		InternalError(c, "重新加载主题服务失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "主题删除成功", nil)
}

// UploadTheme 上传主题
func UploadTheme(c *gin.Context) {
	// 获取上传配置
	cfg := config.Get()
	themesPath := cfg.Server.ThemePath
	if themesPath == "" {
		themesPath = "../templates/themes"
	}

	// 确保主题目录存在
	if err := os.MkdirAll(themesPath, 0755); err != nil {
		InternalError(c, "创建主题目录失败")
		return
	}

	// 获取文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		BadRequest(c, "请选择要上传的主题文件")
		return
	}
	defer file.Close()

	// 检查文件扩展名
	if filepath.Ext(header.Filename) != ".zip" {
		BadRequest(c, "请上传 ZIP 格式的主题文件")
		return
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "theme-*.zip")
	if err != nil {
		InternalError(c, "创建临时文件失败")
		return
	}
	defer os.Remove(tempFile.Name())

	// 复制文件内容
	if _, err := io.Copy(tempFile, file); err != nil {
		InternalError(c, "保存文件失败")
		return
	}
	tempFile.Close()

	// 解压主题
	themeName, err := extractTheme(tempFile.Name(), themesPath)
	if err != nil {
		BadRequest(c, fmt.Sprintf("解压主题失败: %v", err))
		return
	}

	// 重新加载主题服务
	if err := services.InitThemeService(); err != nil {
		InternalError(c, fmt.Sprintf("重新加载主题失败: %v", err))
		return
	}

	SuccessWithMessage(c, "主题上传成功", map[string]string{
		"theme_name": themeName,
	})
}

// extractTheme 解压主题
func extractTheme(zipPath, destPath string) (string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var themeName string

	for _, file := range reader.File {
		// 获取主题名称（从第一级目录）
		if themeName == "" && file.FileInfo().IsDir() {
			themeName = file.Name
			if len(themeName) > 0 && themeName[len(themeName)-1] == '/' {
				themeName = themeName[:len(themeName)-1]
			}
		}

		filePath := filepath.Join(destPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, 0755)
			continue
		}

		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return "", err
		}

		// 创建文件
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return "", err
		}
		defer outFile.Close()

		// 复制内容
		rc, err := file.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		if _, err := io.Copy(outFile, rc); err != nil {
			return "", err
		}
	}

	if themeName == "" {
		return "", fmt.Errorf("无效的主题压缩包")
	}

	return themeName, nil
}

// GetThemeModules 获取主题模块列表
func GetThemeModules(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 获取模块配置
	modulesConfig := make(map[string]interface{})
	if config, err := services.GetThemeFullConfig(slug); err == nil {
		if m, ok := config["modules"].(map[string]interface{}); ok {
			modulesConfig = m
		}
	}

	// 从 theme.yaml 或 settings.yaml 读取模块定义
	modules := []map[string]interface{}{
		{"id": "hero", "name": "首屏大图", "description": "首页顶部的大图展示区域", "icon": "Monitor"},
		{"id": "services", "name": "服务展示", "description": "展示服务项目的模块", "icon": "Service"},
		{"id": "stats", "name": "数据统计", "description": "展示关键数据的统计模块", "icon": "TrendCharts"},
		{"id": "team", "name": "团队介绍", "description": "展示团队成员的模块", "icon": "UserFilled"},
		{"id": "news", "name": "新闻动态", "description": "展示最新新闻的模块", "icon": "Document"},
		{"id": "contact", "name": "联系我们", "description": "联系信息和表单模块", "icon": "Phone"},
	}

	// 合并配置（enabled, order）
	for i, module := range modules {
		moduleID := module["id"].(string)
		if config, ok := modulesConfig[moduleID].(map[string]interface{}); ok {
			if enabled, ok := config["enabled"].(bool); ok {
				module["enabled"] = enabled
			} else {
				module["enabled"] = true
			}
			if order, ok := config["order"].(int); ok {
				module["order"] = order
			} else {
				module["order"] = i + 1
			}
		} else {
			module["enabled"] = true
			module["order"] = i + 1
		}
	}

	SuccessResponse(c, gin.H{
		"theme":   theme.Name,
		"modules": modules,
	})
}

// UpdateModuleConfig 更新模块配置
func UpdateModuleConfig(c *gin.Context) {
	slug := c.Param("slug")
	moduleID := c.Param("moduleId")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	var configData map[string]interface{}
	if err := c.ShouldBindJSON(&configData); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 获取当前配置
	config, err := services.GetThemeFullConfig(slug)
	if err != nil {
		config = make(map[string]interface{})
	}

	// 更新模块配置
	if config["modules"] == nil {
		config["modules"] = make(map[string]interface{})
	}
	modules := config["modules"].(map[string]interface{})
	modules[moduleID] = configData

	// 保存配置
	if err := services.UpdateThemeConfig(slug, config); err != nil {
		InternalError(c, "保存配置失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "模块配置更新成功", gin.H{
		"moduleId": moduleID,
	})
}

// UpdateModuleOrder 更新模块排序
func UpdateModuleOrder(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	var reqData struct {
		ModuleIds []string `json:"moduleIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 获取当前配置
	config, err := services.GetThemeFullConfig(slug)
	if err != nil {
		config = make(map[string]interface{})
	}

	// 更新模块排序
	if config["modules"] == nil {
		config["modules"] = make(map[string]interface{})
	}
	modules := config["modules"].(map[string]interface{})

	// 更新每个模块的 order
	for i, moduleID := range reqData.ModuleIds {
		if moduleConfig, exists := modules[moduleID]; exists {
			if moduleMap, ok := moduleConfig.(map[string]interface{}); ok {
				moduleMap["order"] = i + 1
				modules[moduleID] = moduleMap
			}
		}
	}

	// 保存配置
	if err := services.UpdateThemeConfig(slug, config); err != nil {
		InternalError(c, "保存配置失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "模块排序更新成功", nil)
}

// GetThemeSettingsSchema 获取主题设置定义（schema）
func GetThemeSettingsSchema(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 从主题的 ConfigDefinition 获取设置定义
	schema := theme.ConfigDefinition
	if schema == nil {
		schema = []services.ConfigSection{}
	}

	SuccessResponse(c, gin.H{
		"schema": schema,
	})
}

// ExportThemeConfig 导出主题配置
func ExportThemeConfig(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 获取完整配置
	config, err := services.GetThemeFullConfig(slug)
	if err != nil {
		config = make(map[string]interface{})
	}

	// 添加主题元数据
	exportData := map[string]interface{}{
		"theme":     theme.Name,
		"slug":      slug,
		"version":   theme.Version,
		" exported": time.Now().Format("2006-01-02 15:04:05"),
		"config":    config,
	}

	SuccessResponse(c, exportData)
}

// ClearThemeCache 清理主题缓存（模板缓存 + 重新加载主题）
func ClearThemeCache(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 清理缓存并重新加载
	if err := services.ClearThemeCache(slug); err != nil {
		InternalError(c, "清理缓存失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "缓存清理成功，主题已重新加载", nil)
}

// ImportThemeConfig 导入主题配置
func ImportThemeConfig(c *gin.Context) {
	slug := c.Param("slug")

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	// 解析导入数据
	var importData struct {
		Config map[string]interface{} `json:"config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&importData); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 验证配置是否为有效 map
	if importData.Config == nil {
		BadRequest(c, "无效的配置文件")
		return
	}

	// 保存配置到数据库
	if err := services.UpdateThemeConfig(slug, importData.Config); err != nil {
		InternalError(c, "导入配置失败: "+err.Error())
		return
	}

	SuccessWithMessage(c, "配置导入成功", nil)
}

// GetPageTemplates 获取主题的页面模板定义
func GetPageTemplates(c *gin.Context) {
	slug := c.Param("slug")

	// 如果是 "current" 或空字符串，获取当前激活主题
	if slug == "current" || slug == "" {
		activeTheme := services.GetActiveTheme()
		if activeTheme == nil {
			NotFound(c, "没有激活的主题")
			return
		}
		SuccessResponse(c, gin.H{
			"templates":    activeTheme.PageTemplates,
			"templateMap": buildTemplateMap(activeTheme.PageTemplates),
		})
		return
	}

	theme := services.GetTheme(slug)
	if theme == nil {
		NotFound(c, "主题不存在")
		return
	}

	SuccessResponse(c, gin.H{
		"templates":    theme.PageTemplates,
		"templateMap": buildTemplateMap(theme.PageTemplates),
	})
}

// buildTemplateMap 将页面模板数组转为 map（key 为 type）
func buildTemplateMap(templates []services.PageTemplate) map[string]services.PageTemplate {
	result := make(map[string]services.PageTemplate)
	for _, t := range templates {
		result[t.Type] = t
	}
	return result
}
