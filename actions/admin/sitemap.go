package admin

import (
	"log"
	"net/http"

	"catter-cms/services"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// SitemapHandler serves the sitemap.xml
func SitemapHandler(c buffalo.Context) error {
	sitemapService := services.GetSitemapService()

	// Get base URL from settings if available
	if baseURL, ok := c.Session().Get("base_url").(string); ok {
		if err := sitemapService.SetBaseURL(baseURL); err != nil {
			log.Printf("Invalid base URL in session: %v", err)
		}
	}

	// Generate sitemap XML
	xmlContent, err := sitemapService.Generate()
	if err != nil {
		log.Printf("Failed to generate sitemap: %v", err)
		return c.Error(http.StatusInternalServerError, err)
	}

	// Set appropriate headers
	c.Response().Header().Set("Content-Type", "application/xml")
	c.Response().Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	return c.Render(http.StatusOK, render.String(xmlContent))
}

// RebuildSitemapHandler regenerates the sitemap
func RebuildSitemapHandler(c buffalo.Context) error {
	sitemapService := services.GetSitemapService()

	_, err := sitemapService.Rebuild()
	if err != nil {
		log.Printf("Failed to rebuild sitemap: %v", err)
		return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
			"success": false,
			"message": "重建sitemap失败: " + err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"success": true,
		"message": "Sitemap重建成功",
	}))
}

// SitemapInfoHandler returns sitemap statistics
func SitemapInfoHandler(c buffalo.Context) error {
	sitemapService := services.GetSitemapService()

	count, err := sitemapService.GetURLCount()
	if err != nil {
		return c.Render(http.StatusInternalServerError, render.JSON(map[string]interface{}{
			"success": false,
			"message": "获取sitemap统计信息失败: " + err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"url_count": count,
		},
	}))
}
