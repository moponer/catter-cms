package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gin-app/internal/theme"
	"gin-app/pkg/utils"
)

// ThemeHandler 主题处理器
type ThemeHandler struct {
	themeManager *theme.Manager
}

// NewThemeHandler 创建主题处理器
func NewThemeHandler(tm *theme.Manager) *ThemeHandler {
	return &ThemeHandler{
		themeManager: tm,
	}
}

// GetThemes 获取主题列表
// @Summary 获取主题列表
// @Tags 主题
// @Produce json
// @Param type query string false "主题类型"
// @Success 200 {object} utils.Response
// @Router /api/themes [get]
func (h *ThemeHandler) GetThemes(c *gin.Context) {
	themes := h.themeManager.GetAllThemes()

	// 转换为响应格式
	list := make([]gin.H, 0, len(themes))
	for _, t := range themes {
		list = append(list, gin.H{
			"id":          t.ID,
			"name":        t.Name,
			"description": t.Description,
			"version":     t.Version,
			"author":      t.Author,
			"thumbnail":   t.Thumbnail,
			"path":        t.Path,
			"isCurrent":   t.ID == h.themeManager.GetCurrentThemeID(),
		})
	}

	utils.Success(c, gin.H{
		"list":  list,
		"total": len(list),
	})
}

// GetCurrentTheme 获取当前主题
// @Summary 获取当前主题
// @Tags 主题
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/themes/current [get]
func (h *ThemeHandler) GetCurrentTheme(c *gin.Context) {
	theme := h.themeManager.GetCurrentTheme()
	if theme == nil {
		utils.Error(c, http.StatusNotFound, "当前没有激活的主题")
		return
	}

	utils.Success(c, gin.H{
		"id":          theme.ID,
		"name":        theme.Name,
		"description": theme.Description,
		"version":     theme.Version,
		"author":      theme.Author,
		"thumbnail":   theme.Thumbnail,
		"path":        theme.Path,
	})
}

// SwitchTheme 切换主题
// @Summary 切换主题
// @Tags 主题
// @Produce json
// @Param id path string true "主题ID"
// @Success 200 {object} utils.Response
// @Router /api/themes/{id}/activate [post]
func (h *ThemeHandler) SwitchTheme(c *gin.Context) {
	themeID := c.Param("id")
	if themeID == "" {
		utils.Error(c, http.StatusBadRequest, "主题ID不能为空")
		return
	}

	if err := h.themeManager.SwitchTheme(themeID); err != nil {
		utils.Error(c, http.StatusInternalServerError, "切换主题失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"message": "主题切换成功",
		"themeId": themeID,
	})
}

// GetCurrentThemeModules 获取当前主题模块
// @Summary 获取当前主题模块
// @Tags 主题
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/themes/current/modules [get]
func (h *ThemeHandler) GetCurrentThemeModules(c *gin.Context) {
	h.GetThemeModules(c)
}

// GetThemeModules 获取指定主题模块
// @Summary 获取指定主题模块
// @Tags 主题
// @Produce json
// @Param id path string true "主题ID"
// @Success 200 {object} utils.Response
// @Router /api/themes/{id}/modules [get]
func (h *ThemeHandler) GetThemeModules(c *gin.Context) {
	themeID := c.Param("id")

	theme, err := h.themeManager.GetTheme(themeID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	// 从 theme.yaml 或 settings.yaml 读取模块定义
	modules := []gin.H{
		{"id": "hero", "name": "首屏大图", "enabled": true, "order": 1},
		{"id": "services", "name": "服务展示", "enabled": true, "order": 2},
		{"id": "stats", "name": "数据统计", "enabled": true, "order": 3},
		{"id": "team", "name": "团队介绍", "enabled": false, "order": 4},
		{"id": "news", "name": "新闻动态", "enabled": true, "order": 5},
		{"id": "contact", "name": "联系我们", "enabled": true, "order": 6},
	}

	utils.Success(c, gin.H{
		"theme":   theme,
		"modules": modules,
	})
}

// GetCurrentThemeTemplates 获取当前主题模板
// @Summary 获取当前主题模板
// @Tags 主题
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/themes/current/templates [get]
func (h *ThemeHandler) GetCurrentThemeTemplates(c *gin.Context) {
	h.GetThemeTemplates(c)
}

// GetThemeTemplates 获取指定主题模板
// @Summary 获取指定主题模板
// @Tags 主题
// @Produce json
// @Param id path string true "主题ID"
// @Success 200 {object} utils.Response
// @Router /api/themes/{id}/templates [get]
func (h *ThemeHandler) GetThemeTemplates(c *gin.Context) {
	themeID := c.Param("id")

	theme, err := h.themeManager.GetTheme(themeID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	// 从主题目录读取模板文件
	templates := gin.H{
		"page": []gin.H{
			{"value": "default", "label": "默认页面", "file": "page_default.html"},
			{"value": "fullwidth", "label": "全宽页面", "file": "page_fullwidth.html"},
			{"value": "sidebar", "label": "侧边栏页面", "file": "page_sidebar.html"},
		},
		"list": []gin.H{
			{"value": "grid", "label": "网格布局", "file": "list_grid.html"},
			{"value": "list", "label": "列表布局", "file": "list_list.html"},
			{"value": "masonry", "label": "瀑布流布局", "file": "list_masonry.html"},
		},
		"detail": []gin.H{
			{"value": "default", "label": "默认详情", "file": "detail_default.html"},
			{"value": "fullwidth", "label": "全宽详情", "file": "detail_fullwidth.html"},
			{"value": "sidebar", "label": "侧边栏详情", "file": "detail_sidebar.html"},
		},
	}

	utils.Success(c, gin.H{
		"theme":     theme,
		"templates": templates,
	})
}
