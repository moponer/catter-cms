package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gin-app/internal/theme"
	"gin-app/pkg/utils"
)

// ThemeConfigHandler 主题配置处理器
type ThemeConfigHandler struct {
	themeManager *theme.Manager
}

// NewThemeConfigHandler 创建配置处理器
func NewThemeConfigHandler(tm *theme.Manager) *ThemeConfigHandler {
	return &ThemeConfigHandler{
		themeManager: tm,
	}
}

// GetCurrentThemeConfig 获取当前主题配置
// @Summary 获取当前主题配置
// @Tags 主题配置
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/themes/current/config [get]
func (h *ThemeConfigHandler) GetCurrentThemeConfig(c *gin.Context) {
	config, err := h.themeManager.GetThemeConfig("")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "获取配置失败: "+err.Error())
		return
	}

	utils.Success(c, config)
}

// GetThemeConfig 获取指定主题配置
// @Summary 获取指定主题配置
// @Tags 主题配置
// @Produce json
// @Param id path string true "主题ID"
// @Success 200 {object} utils.Response
// @Router /api/themes/{id}/config [get]
func (h *ThemeConfigHandler) GetThemeConfig(c *gin.Context) {
	themeID := c.Param("id")
	if themeID == "" {
		utils.Error(c, http.StatusBadRequest, "主题ID不能为空")
		return
	}

	config, err := h.themeManager.GetThemeConfig(themeID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "获取配置失败: "+err.Error())
		return
	}

	utils.Success(c, config)
}

// SaveCurrentThemeConfig 保存当前主题配置
// @Summary 保存当前主题配置
// @Tags 主题配置
// @Accept json
// @Produce json
// @Param config body map[string]interface{} true "配置数据"
// @Success 200 {object} utils.Response
// @Router /api/themes/current/config [put]
func (h *ThemeConfigHandler) SaveCurrentThemeConfig(c *gin.Context) {
	var configData map[string]interface{}
	if err := c.ShouldBindJSON(&configData); err != nil {
		utils.Error(c, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}

	if err := h.themeManager.SaveThemeConfig("", configData); err != nil {
		utils.Error(c, http.StatusInternalServerError, "保存配置失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "配置保存成功"})
}

// SaveThemeConfig 保存指定主题配置
// @Summary 保存指定主题配置
// @Tags 主题配置
// @Accept json
// @Produce json
// @Param id path string true "主题ID"
// @Param config body map[string]interface{} true "配置数据"
// @Success 200 {object} utils.Response
// @Router /api/themes/{id}/config [put]
func (h *ThemeConfigHandler) SaveThemeConfig(c *gin.Context) {
	themeID := c.Param("id")
	if themeID == "" {
		utils.Error(c, http.StatusBadRequest, "主题ID不能为空")
		return
	}

	var configData map[string]interface{}
	if err := c.ShouldBindJSON(&configData); err != nil {
		utils.Error(c, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}

	if err := h.themeManager.SaveThemeConfig(themeID, configData); err != nil {
		utils.Error(c, http.StatusInternalServerError, "保存配置失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "配置保存成功"})
}

// ResetCurrentThemeConfig 重置当前主题配置
// @Summary 重置当前主题配置
// @Tags 主题配置
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/themes/current/config/reset [post]
func (h *ThemeConfigHandler) ResetCurrentThemeConfig(c *gin.Context) {
	if err := h.themeManager.ResetThemeConfig(""); err != nil {
		utils.Error(c, http.StatusInternalServerError, "重置配置失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "配置已重置"})
}

// ResetThemeConfig 重置指定主题配置
// @Summary 重置指定主题配置
// @Tags 主题配置
// @Produce json
// @Param id path string true "主题ID"
// @Success 200 {object} utils.Response
// @Router /api/themes/{id}/config/reset [post]
func (h *ThemeConfigHandler) ResetThemeConfig(c *gin.Context) {
	themeID := c.Param("id")
	if themeID == "" {
		utils.Error(c, http.StatusBadRequest, "主题ID不能为空")
		return
	}

	if err := h.themeManager.ResetThemeConfig(themeID); err != nil {
		utils.Error(c, http.StatusInternalServerError, "重置配置失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "配置已重置"})
}

// GetCurrentThemeSchema 获取当前主题设置定义
// @Summary 获取当前主题设置定义
// @Tags 主题配置
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/themes/current/schema [get]
func (h *ThemeConfigHandler) GetCurrentThemeSchema(c *gin.Context) {
	schema, err := h.themeManager.GetThemeSchema("")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "获取设置定义失败: "+err.Error())
		return
	}

	utils.Success(c, schema)
}

// GetThemeSchema 获取指定主题设置定义
// @Summary 获取指定主题设置定义
// @Tags 主题配置
// @Produce json
// @Param id path string true "主题ID"
// @Success 200 {object} utils.Response
// @Router /api/themes/{id}/schema [get]
func (h *ThemeConfigHandler) GetThemeSchema(c *gin.Context) {
	themeID := c.Param("id")
	if themeID == "" {
		utils.Error(c, http.StatusBadRequest, "主题ID不能为空")
		return
	}

	schema, err := h.themeManager.GetThemeSchema(themeID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "获取设置定义失败: "+err.Error())
		return
	}

	utils.Success(c, schema)
}
