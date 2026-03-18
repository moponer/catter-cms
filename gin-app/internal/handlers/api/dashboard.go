package api

import (
	"github.com/gin-gonic/gin"
	"catter-cms-gin/internal/models"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalPosts       int64 `json:"total_posts"`
	PublishedPosts   int64 `json:"published_posts"`
	DraftPosts       int64 `json:"draft_posts"`
	TotalUsers       int64 `json:"total_users"`
	TotalCategories  int64 `json:"total_categories"`
	TotalTags        int64 `json:"total_tags"`
	TotalMedia       int64 `json:"total_media"`
	TotalComments    int64 `json:"total_comments"`
	TotalPages       int64 `json:"total_pages"`
}

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// 统计文章
	models.DB.Model(&models.Post{}).Count(&stats.TotalPosts)
	models.DB.Model(&models.Post{}).Where("status = ?", "published").Count(&stats.PublishedPosts)
	models.DB.Model(&models.Post{}).Where("status = ?", "draft").Count(&stats.DraftPosts)

	// 统计用户
	models.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	// 统计分类
	models.DB.Model(&models.Category{}).Count(&stats.TotalCategories)

	// 统计标签
	models.DB.Model(&models.Tag{}).Count(&stats.TotalTags)

	// 统计媒体
	models.DB.Model(&models.Media{}).Count(&stats.TotalMedia)

	// 统计评论
	models.DB.Model(&models.Comment{}).Count(&stats.TotalComments)

	// 统计页面
	models.DB.Model(&models.Page{}).Count(&stats.TotalPages)

	SuccessResponse(c, stats)
}

// GetRecentPosts 获取最近文章
func GetRecentPosts(c *gin.Context) {
	var posts []models.Post
	models.DB.Order("created_at DESC").Limit(10).Find(&posts)

	// 转换为前端需要的格式
	var result []gin.H
	for _, post := range posts {
		result = append(result, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"status":     post.Status,
			"created_at": post.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	SuccessResponse(c, result)
}
