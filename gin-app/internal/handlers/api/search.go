package api

import (
	"strings"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// GlobalSearch 全局搜索
func GlobalSearch(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		keyword = c.Query("keyword")
	}
	if keyword == "" {
		BadRequest(c, "请提供搜索关键词")
		return
	}

	keyword = strings.TrimSpace(keyword)
	searchPattern := "%" + keyword + "%"

	// 搜索文章
	var posts []models.Post
	models.DB.Where("title LIKE ? OR content LIKE ? OR excerpt LIKE ?", 
		searchPattern, searchPattern, searchPattern).
		Where("status = ?", "published").
		Limit(20).
		Find(&posts)

	// 搜索页面
	var pages []models.Page
	models.DB.Where("title LIKE ? OR content LIKE ?", 
		searchPattern, searchPattern).
		Where("status = ?", "published").
		Limit(20).
		Find(&pages)

	SuccessResponse(c, gin.H{
		"keyword": keyword,
		"posts":   posts,
		"pages":   pages,
		"total":   len(posts) + len(pages),
	})
}
