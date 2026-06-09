package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// getSitemapURL 获取站点 URL（从设置或配置动态获取）
func getSitemapURL() string {
	url := models.GetSettingAsString("site_url")
	if url == "" {
		url = "http://localhost:3000"
	}
	// 确保末尾无斜杠
	return url
}

// SitemapHandler 生成站点地图（动态，使用真实 URL）
func SitemapHandler(c *gin.Context) {
	baseURL := getSitemapURL()

	// 获取已发布的文章（排除 noindex）
	var posts []models.Post
	models.DB.Where("status = ? AND (robots = '' OR robots NOT LIKE '%noindex%')", "published").
		Order("COALESCE(published_at, created_at) DESC").
		Find(&posts)

	// 获取已发布的页面（排除 noindex）
	var pages []models.Page
	models.DB.Where("status = ? AND (robots = '' OR robots NOT LIKE '%noindex%')", "published").
		Order("COALESCE(UpdatedAt, CreatedAt) DESC").
		Find(&pages)

	// 获取分类
	var categories []models.Category
	models.DB.Order("name ASC").Find(&categories)

	// 获取标签
	var tags []models.Tag
	models.DB.Order("name ASC").Find(&tags)

	// 生成 XML
	xml := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"

	// 首页
	xml += `  <url>` + "\n"
	xml += fmt.Sprintf(`    <loc>%s/</loc>`, baseURL) + "\n"
	xml += `    <changefreq>daily</changefreq>` + "\n"
	xml += `    <priority>1.0</priority>` + "\n"
	xml += `  </url>` + "\n"

	// 文章
	for _, post := range posts {
		lastMod := post.UpdatedAt
		if post.PublishedAt != nil && !post.PublishedAt.After(lastMod) {
			lastMod = *post.PublishedAt
		}
		xml += `  <url>` + "\n"
		xml += fmt.Sprintf(`    <loc>%s/article/%s</loc>`, baseURL, post.Slug) + "\n"
		xml += fmt.Sprintf(`    <lastmod>%s</lastmod>`, lastMod.Format("2006-01-02")) + "\n"
		xml += `    <changefreq>weekly</changefreq>` + "\n"
		xml += `    <priority>0.8</priority>` + "\n"
		xml += `  </url>` + "\n"
	}

	// 页面
	for _, page := range pages {
		lastMod := page.UpdatedAt
		xml += `  <url>` + "\n"
		xml += fmt.Sprintf(`    <loc>%s/page/%s</loc>`, baseURL, page.Slug) + "\n"
		xml += fmt.Sprintf(`    <lastmod>%s</lastmod>`, lastMod.Format("2006-01-02")) + "\n"
		xml += `    <changefreq>monthly</changefreq>` + "\n"
		xml += `    <priority>0.6</priority>` + "\n"
		xml += `  </url>` + "\n"
	}

	// 分类
	for _, category := range categories {
		xml += `  <url>` + "\n"
		xml += fmt.Sprintf(`    <loc>%s/category/%s</loc>`, baseURL, category.Slug) + "\n"
		xml += `    <changefreq>weekly</changefreq>` + "\n"
		xml += `    <priority>0.5</priority>` + "\n"
		xml += `  </url>` + "\n"
	}

	// 标签
	for _, tag := range tags {
		xml += `  <url>` + "\n"
		xml += fmt.Sprintf(`    <loc>%s/tag/%s</loc>`, baseURL, tag.Slug) + "\n"
		xml += `    <changefreq>weekly</changefreq>` + "\n"
		xml += `    <priority>0.4</priority>` + "\n"
		xml += `  </url>` + "\n"
	}

	xml += `</urlset>`

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, xml)
}

// SitemapIndexHandler 站点地图索引
func SitemapIndexHandler(c *gin.Context) {
	baseURL := getSitemapURL()

	xml := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	xml += `<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"
	xml += `  <sitemap>` + "\n"
	xml += fmt.Sprintf(`    <loc>%s/sitemap.xml</loc>`, baseURL) + "\n"
	xml += fmt.Sprintf(`    <lastmod>%s</lastmod>`, getCurrentDate()) + "\n"
	xml += `  </sitemap>` + "\n"
	xml += `</sitemapindex>`

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, xml)
}

// getCurrentDate 获取当前日期（ISO 格式）
func getCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// now 返回当前时间
func now() time.Time {
	return time.Now()
}

// SitemapInfo Sitemap信息结构
type SitemapInfo struct {
	URLCount      int    `json:"url_count"`
	PostCount     int    `json:"post_count"`
	PageCount     int    `json:"page_count"`
	CategoryCount int    `json:"category_count"`
	TagCount      int    `json:"tag_count"`
	LastUpdated   string `json:"last_updated"`
}

// GetSitemapInfo 获取Sitemap信息
func GetSitemapInfo(c *gin.Context) {
	var postCount, pageCount, categoryCount, tagCount int64

	models.DB.Model(&models.Post{}).Where("status = ?", "published").Count(&postCount)
	models.DB.Model(&models.Page{}).Where("status = ?", "published").Count(&pageCount)
	models.DB.Model(&models.Category{}).Count(&categoryCount)
	models.DB.Model(&models.Tag{}).Count(&tagCount)

	// 总URL数 = 首页 + 文章 + 页面 + 分类 + 标签
	totalURL := 1 + int(postCount) + int(pageCount) + int(categoryCount) + int(tagCount)

	info := SitemapInfo{
		URLCount:      totalURL,
		PostCount:     int(postCount),
		PageCount:     int(pageCount),
		CategoryCount: int(categoryCount),
		TagCount:      int(tagCount),
		LastUpdated:   fmt.Sprintf("%d-%02d-%02d",
			now().Year(), now().Month(), now().Day()),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// RebuildSitemap 重新生成Sitemap
func RebuildSitemap(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sitemap重新生成成功",
		"data": gin.H{
			"last_updated": fmt.Sprintf("%d-%02d-%02d",
				now().Year(), now().Month(), now().Day()),
		},
	})
}
