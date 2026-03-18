package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// SitemapHandler 生成站点地图
func SitemapHandler(c *gin.Context) {
	baseURL := "http://localhost:3001"

	// 获取已发布的文章
	var posts []models.Post
	models.DB.Where("status = ?", "published").Find(&posts)

	// 获取已发布的页面
	var pages []models.Page
	models.DB.Where("status = ?", "published").Find(&pages)

	// 获取分类
	var categories []models.Category
	models.DB.Find(&categories)

	// 生成 XML
	xml := `<?xml version="1.0" encoding="UTF-8"?>`
	xml += "\n"
	xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	xml += "\n"

	// 首页
	xml += `  <url>`
	xml += "\n"
	xml += fmt.Sprintf(`    <loc>%s/</loc>`, baseURL)
	xml += "\n"
	xml += `    <changefreq>daily</changefreq>`
	xml += "\n"
	xml += `    <priority>1.0</priority>`
	xml += "\n"
	xml += `  </url>`
	xml += "\n"

	// 文章
	for _, post := range posts {
		xml += `  <url>`
		xml += "\n"
		xml += fmt.Sprintf(`    <loc>%s/blog/%s</loc>`, baseURL, post.Slug)
		xml += "\n"
		if !post.CreatedAt.IsZero() {
			xml += fmt.Sprintf(`    <lastmod>%s</lastmod>`, post.CreatedAt.Format("2006-01-02"))
		}
		xml += "\n"
		xml += `    <changefreq>weekly</changefreq>`
		xml += "\n"
		xml += `    <priority>0.8</priority>`
		xml += "\n"
		xml += `  </url>`
		xml += "\n"
	}

	// 页面
	for _, page := range pages {
		xml += `  <url>`
		xml += "\n"
		xml += fmt.Sprintf(`    <loc>%s/page/%s</loc>`, baseURL, page.Slug)
		xml += "\n"
		if !page.CreatedAt.IsZero() {
			xml += fmt.Sprintf(`    <lastmod>%s</lastmod>`, page.CreatedAt.Format("2006-01-02"))
		}
		xml += "\n"
		xml += `    <changefreq>monthly</changefreq>`
		xml += "\n"
		xml += `    <priority>0.6</priority>`
		xml += "\n"
		xml += `  </url>`
		xml += "\n"
	}

	// 分类
	for _, category := range categories {
		xml += `  <url>`
		xml += "\n"
		xml += fmt.Sprintf(`    <loc>%s/category/%s</loc>`, baseURL, category.Slug)
		xml += "\n"
		xml += `    <changefreq>weekly</changefreq>`
		xml += "\n"
		xml += `    <priority>0.5</priority>`
		xml += "\n"
		xml += `  </url>`
		xml += "\n"
	}

	xml += `</urlset>`

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, xml)
}

// SitemapIndexHandler 站点地图索引
func SitemapIndexHandler(c *gin.Context) {
	baseURL := "http://localhost:3001"

	xml := `<?xml version="1.0" encoding="UTF-8"?>`
	xml += "\n"
	xml += `<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	xml += "\n"
	xml += `  <sitemap>`
	xml += "\n"
	xml += fmt.Sprintf(`    <loc>%s/sitemap.xml</loc>`, baseURL)
	xml += "\n"
	xml += fmt.Sprintf(`    <lastmod>%s</lastmod>`, time.Now().Format("2006-01-02"))
	xml += "\n"
	xml += `  </sitemap>`
	xml += "\n"
	xml += `</sitemapindex>`

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, xml)
}

// SitemapInfo Sitemap信息结构
type SitemapInfo struct {
	URLCount      int    `json:"url_count"`
	PostCount     int    `json:"post_count"`
	PageCount     int    `json:"page_count"`
	CategoryCount int    `json:"category_count"`
	LastUpdated   string `json:"last_updated"`
}

// GetSitemapInfo 获取Sitemap信息
func GetSitemapInfo(c *gin.Context) {
	var postCount, pageCount, categoryCount int64

	models.DB.Model(&models.Post{}).Where("status = ?", "published").Count(&postCount)
	models.DB.Model(&models.Page{}).Where("status = ?", "published").Count(&pageCount)
	models.DB.Model(&models.Category{}).Count(&categoryCount)

	// 总URL数 = 首页 + 文章 + 页面 + 分类
	totalURL := 1 + int(postCount) + int(pageCount) + int(categoryCount)

	info := SitemapInfo{
		URLCount:      totalURL,
		PostCount:     int(postCount),
		PageCount:     int(pageCount),
		CategoryCount: int(categoryCount),
		LastUpdated:   time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// RebuildSitemap 重新生成Sitemap
func RebuildSitemap(c *gin.Context) {
	// 这里可以实现重新生成sitemap.xml文件的逻辑
	// 目前sitemap是动态生成的，这里返回成功即可
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sitemap重新生成成功",
		"data": gin.H{
			"last_updated": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}
