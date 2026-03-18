package api

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
		Name:        theme.Name,
		Slug:        theme.Name,
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
		"slug":        theme.Name,
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

	// TODO: 实现从文件系统删除主题的逻辑
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
