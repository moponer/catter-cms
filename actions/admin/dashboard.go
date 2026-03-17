package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalPosts      int `json:"total_posts"`
	PublishedPosts  int `json:"published_posts"`
	DraftPosts      int `json:"draft_posts"`
	TotalUsers      int `json:"total_users"`
	TotalCategories int `json:"total_categories"`
	TotalTags       int `json:"total_tags"`
	TotalMedia      int `json:"total_media"`
}

// DashboardStatsAPI 获取仪表盘统计数据
func DashboardStatsAPI(c buffalo.Context) error {
	stats := DashboardStats{}

	// 统计文章数
	posts := &models.Posts{}
	if err := models.DB.All(posts); err == nil {
		stats.TotalPosts = len(*posts)
		for _, p := range *posts {
			if p.Status == "published" {
				stats.PublishedPosts++
			} else if p.Status == "draft" {
				stats.DraftPosts++
			}
		}
	}

	// 统计用户数
	users := &models.Users{}
	if err := models.DB.All(users); err == nil {
		stats.TotalUsers = len(*users)
	}

	// 统计分类数
	categories := &models.Categories{}
	if err := models.DB.All(categories); err == nil {
		stats.TotalCategories = len(*categories)
	}

	// 统计标签数
	tags := &models.Tags{}
	if err := models.DB.All(tags); err == nil {
		stats.TotalTags = len(*tags)
	}

	// 统计媒体数
	media := &models.MediaItems{}
	if err := models.DB.All(media); err == nil {
		stats.TotalMedia = len(*media)
	}

	return api.Success(c, stats)
}

// DashboardRecentPostsAPI 获取最近文章
func DashboardRecentPostsAPI(c buffalo.Context) error {
	posts := &models.Posts{}
	query := models.DB.Order("created_at DESC").Limit(5)

	if err := query.All(posts); err != nil {
		return api.InternalError(c, "获取文章失败")
	}

	return api.Success(c, posts)
}

// DashboardHandler 显示仪表盘页面 (用于服务端渲染，如果需要)
func DashboardHandler(c buffalo.Context) error {
	return c.Render(200, render.HTML("<html><body><h1>Dashboard</h1><p>Admin SPA is served at /admin</p></body></html>"))
}
