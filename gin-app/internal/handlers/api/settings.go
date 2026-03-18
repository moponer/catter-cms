package api

import (
	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
	"catter-cms-gin/internal/services"
)

// GetSettings 获取设置
func GetSettings(c *gin.Context) {
	settings, err := models.GetAllSettings()
	if err != nil {
		InternalError(c, "获取设置失败")
		return
	}
	SuccessResponse(c, settings)
}

// SaveSettings 保存设置
func SaveSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 更新数据库中的设置
	if err := models.SetSettings(settings); err != nil {
		InternalError(c, "保存设置失败")
		return
	}

	// 如果设置了激活主题，同时更新主题服务
	if activeTheme, ok := settings["active_theme"].(string); ok && activeTheme != "" {
		services.SetActiveTheme(activeTheme)
	}

	// 返回更新后的所有设置
	allSettings, _ := models.GetAllSettings()
	SuccessWithMessage(c, "设置保存成功", allSettings)
}

// GetSetting 获取单个设置
func GetSetting(c *gin.Context) {
	key := c.Param("key")

	value, err := models.GetSetting(key)
	if err != nil {
		NotFound(c, "设置不存在")
		return
	}

	SuccessResponse(c, gin.H{
		"key":   key,
		"value": value,
	})
}

// UpdateSetting 更新单个设置
func UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var req struct {
		Value interface{} `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if err := models.SetSetting(key, req.Value); err != nil {
		InternalError(c, "更新设置失败")
		return
	}

	// 如果是激活主题设置，同时更新主题服务
	if key == "active_theme" {
		if themeName, ok := req.Value.(string); ok {
			services.SetActiveTheme(themeName)
		}
	}

	SuccessWithMessage(c, "设置更新成功", gin.H{
		"key":   key,
		"value": req.Value,
	})
}

// GetActiveThemeSetting 获取当前激活的主题设置
func GetActiveThemeSetting(c *gin.Context) {
	activeTheme := models.GetActiveTheme()
	SuccessResponse(c, gin.H{
		"active_theme": activeTheme,
	})
}

// SetActiveThemeSetting 设置激活的主题
func SetActiveThemeSetting(c *gin.Context) {
	var req struct {
		Theme string `json:"theme" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请提供主题名称")
		return
	}

	// 验证主题是否存在
	theme := services.GetTheme(req.Theme)
	if theme == nil {
		BadRequest(c, "主题不存在")
		return
	}

	// 保存到数据库
	if err := models.SetActiveTheme(req.Theme); err != nil {
		InternalError(c, "保存主题设置失败")
		return
	}

	// 更新主题服务
	if err := services.SetActiveTheme(req.Theme); err != nil {
		InternalError(c, "激活主题失败")
		return
	}

	SuccessWithMessage(c, "主题已激活", gin.H{
		"active_theme": req.Theme,
	})
}
