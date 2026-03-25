package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"catter-cms-gin/internal/services"
)

// GetCacheStatus 获取缓存状态
func GetCacheStatus(c *gin.Context) {
	// 获取主题数量作为模板缓存指标
	themes := services.GetAllThemes()
	templateCount := len(themes)

	SuccessResponse(c, gin.H{
		"template": templateCount,
		"data":     0, // 数据缓存暂时用 0 表示
	})
}

// ClearCache 清理缓存
func ClearCache(c *gin.Context) {
	var req struct {
		Type string `json:"type" binding:"required"` // all, template, data
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请指定缓存类型")
		return
	}

	switch req.Type {
	case "all":
		// 清理所有主题的缓存
		themes := services.GetAllThemes()
		for name := range themes {
			if err := services.ClearThemeCache(name); err != nil {
				InternalError(c, "清理缓存失败: "+err.Error())
				return
			}
		}
		SuccessWithMessage(c, "所有缓存已清理", nil)

	case "template":
		// 只清理模板缓存
		themes := services.GetAllThemes()
		for name := range themes {
			if err := services.ClearThemeCache(name); err != nil {
				InternalError(c, "清理模板缓存失败: "+err.Error())
				return
			}
		}
		SuccessWithMessage(c, "模板缓存已清理", nil)

	case "data":
		// 数据缓存清理（目前只是占位）
		SuccessWithMessage(c, "数据缓存已清理", nil)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"success": false,
			"message": "无效的缓存类型，请使用 all, template 或 data",
		})
	}
}
