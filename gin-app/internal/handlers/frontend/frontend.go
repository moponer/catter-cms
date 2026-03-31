package frontend

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
	"catter-cms-gin/internal/services"
)

// parseJSON 解析 JSON 字符串为 map
func parseJSON(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// getSiteData 获取站点基础数据
func getSiteData() map[string]interface{} {
	siteTitle := models.GetSettingAsString("site_title")
	if siteTitle == "" {
		siteTitle = "Catter CMS"
	}
	return map[string]interface{}{
		"title":       siteTitle,
		"description": models.GetSettingAsString("site_description"),
		"url":         models.GetSettingAsString("site_url"),
		"logo":        models.GetSettingAsString("logo"),
		"favicon":     models.GetSettingAsString("favicon"),
		"year":        time.Now().Year(),
	}
}

// Render404 渲染 404 页面：优先使用当前主题的 404 模板，回退到系统 404.html
func Render404(c *gin.Context) {
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "页面未找到",
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusNotFound)
	// 尝试用当前主题渲染 404 模板
	if content, err := services.RenderTemplate("404", data); err == nil {
		c.String(http.StatusNotFound, string(content))
		return
	}
	// 回退到系统静态 404.html
	c.File("./templates/404.html")
}

// HomeHandler 首页
func HomeHandler(c *gin.Context) {
	var posts []models.Post
	models.DB.Where("status = ?", "published").
		Preload("Category").
		Preload("Tags").
		Order("COALESCE(published_at, created_at) DESC").
		Limit(10).
		Find(&posts)

	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("COALESCE(published_at, created_at) DESC").
		Limit(5).
		Find(&recentPosts)

	var categories []models.Category
	models.DB.Find(&categories)

	var tags []models.Tag
	models.DB.Find(&tags)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   "首页",
		"posts":       posts,
		"recentPosts": recentPosts,
		"categories":  categories,
		"tags":        tags,
	}

	theme := services.GetActiveTheme()
	if theme == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "没有激活的主题",
		})
		return
	}

	content, err := services.RenderTemplate("home", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "主题渲染失败",
			"error":   err.Error(),
			"theme":   theme.Name,
			"path":    theme.Path,
		})
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(content))
}

// ArticleShowHandler 文章详情页
func ArticleShowHandler(c *gin.Context) {
	slug := c.Param("slug")

	var post models.Post
	if err := models.DB.Where("slug = ? AND status = ?", slug, "published").
		Preload("Category").
		Preload("Tags").
		First(&post).Error; err != nil {
		Render404(c)
		return
	}

	models.DB.Model(&post).Update("views", post.Views+1)

	var categories []models.Category
	models.DB.Find(&categories)
	var tags []models.Tag
	models.DB.Find(&tags)

	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentPosts)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   post.Title,
		"post":        post,
		"categories":  categories,
		"tags":        tags,
		"recentPosts": recentPosts,
	}

	if content, err := services.RenderTemplate("article_show", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// 默认页面 - XSS 安全修复：使用 html.EscapeString 转义
	defaultHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Catter CMS</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
               background: #f5f5f5; min-height: 100vh; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; margin-bottom: 20px; }
        .meta { color: #666; font-size: 0.9em; margin-bottom: 20px; }
        .content { line-height: 1.8; color: #555; }
        .back-link { display: inline-block; margin-top: 20px; color: #667eea; text-decoration: none; }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <div class="meta">Views: %d | Published: %s</div>
        <div class="content">%s</div>
        <a href="/" class="back-link">Back to Home</a>
    </div>
</body>
</html>`,
		html.EscapeString(post.Title),
		html.EscapeString(post.Title),
		post.Views,
		post.PublishedAt.Format("2006-01-02"),
		html.EscapeString(post.Content))

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

// ArticleListHandler 文章列表页
func ArticleListHandler(c *gin.Context) {
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("pageSize"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	var posts []models.Post
	var total int64
	models.DB.Where("status = ?", "published").
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts)
	models.DB.Model(&models.Post{}).Where("status = ?", "published").Count(&total)

	var categories []models.Category
	var tags []models.Tag
	models.DB.Find(&categories)
	models.DB.Find(&tags)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  "文章",
		"posts":      posts,
		"categories": categories,
		"tags":       tags,
		"pagination": map[string]interface{}{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	content, err := services.RenderTemplate("article_list", data)
	if err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	log.Printf("[ERROR] article_list template render failed: %v", err)
	c.String(http.StatusOK, "<h1>文章列表</h1><p>请创建 article_list.html 模板</p>")
}

// PageShowHandler 页面详情页
func PageShowHandler(c *gin.Context) {
	slug := c.Param("slug")

	var page models.Page
	if err := models.DB.Where("slug = ? AND status = ?", slug, "published").First(&page).Error; err != nil {
		Render404(c)
		return
	}

	models.DB.Model(&page).Update("views", page.Views+1)

	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentPosts)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   page.Title,
		"page":        page,
		"recentPosts": recentPosts,
	}

	templateName := "page"
	if page.Template != "" && page.Template != "default" {
		templateName = "page-" + page.Template
	}

	// 解析页面模板专属数据
	if page.TemplateData != "" {
		data["templateData"] = page.TemplateData // 原始 JSON 字符串
		// 解析为 map 方便模板直接访问
		if varData, err := parseJSON(page.TemplateData); err == nil {
			data["td"] = varData
		}
	}

	if content, err := services.RenderTemplate(templateName, data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	if templateName != "page" {
		if content, err := services.RenderTemplate("page", data); err == nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, string(content))
			return
		}
	}

	defaultHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Catter CMS</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
               background: #f5f5f5; min-height: 100vh; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; margin-bottom: 20px; }
        .content { line-height: 1.8; color: #555; }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <div class="content">%s</div>
    </div>
</body>
</html>`,
		html.EscapeString(page.Title),
		html.EscapeString(page.Title),
		html.EscapeString(page.Content))

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

// CategoryHandler 分类页
func CategoryHandler(c *gin.Context) {
	slug := c.Param("slug")

	var category models.Category
	if err := models.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		Render404(c)
		return
	}

	var posts []models.Post
	models.DB.Where("category_id = ? AND status = ?", category.ID, "published").
		Preload("Category").
		Order("published_at DESC").
		Find(&posts)

	var categories []models.Category
	models.DB.Find(&categories)
	var tags []models.Tag
	models.DB.Find(&tags)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  category.Name,
		"category":   category,
		"posts":      posts,
		"categories": categories,
		"tags":       tags,
	}

	if content, err := services.RenderTemplate("category", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// 默认页面 - XSS 安全修复
	postsHTML := ""
	for _, post := range posts {
		postsHTML += fmt.Sprintf(`<div class="post-item"><h3><a href="/article/%s">%s</a></h3></div>`,
			html.EscapeString(post.Slug),
			html.EscapeString(post.Title))
	}
	if postsHTML == "" {
		postsHTML = "<p>No posts in this category.</p>"
	}

	defaultHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Catter CMS</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
               background: #f5f5f5; min-height: 100vh; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; margin-bottom: 20px; }
        .post-item { padding: 10px 0; border-bottom: 1px solid #eee; }
        .post-item a { color: #667eea; text-decoration: none; }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <p>%s</p>
        <div class="posts">%s</div>
    </div>
</body>
</html>`,
		html.EscapeString(category.Name),
		html.EscapeString(category.Name),
		html.EscapeString(category.Description),
		postsHTML)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

// CategoriesListHandler 分类列表页
func CategoriesListHandler(c *gin.Context) {
	var categories []models.Category
	models.DB.Find(&categories)

	type CategoryWithCount struct {
		models.Category
		PostCount int `json:"post_count"`
	}
	var categoriesWithCount []CategoryWithCount
	for _, cat := range categories {
		var count int64
		models.DB.Model(&models.Post{}).Where("category_id = ? AND status = ?", cat.ID, "published").Count(&count)
		categoriesWithCount = append(categoriesWithCount, CategoryWithCount{
			Category:  cat,
			PostCount: int(count),
		})
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  "分类",
		"categories": categoriesWithCount,
	}

	if content, err := services.RenderTemplate("categories", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>分类列表页</h1><p>请创建 categories.plush.html 模板</p>")
}

// TagsListHandler 标签列表页
func TagsListHandler(c *gin.Context) {
	var tags []models.Tag
	models.DB.Find(&tags)

	type TagWithCount struct {
		models.Tag
		PostCount int `json:"post_count"`
	}
	var tagsWithCount []TagWithCount
	for _, tag := range tags {
		var count int64
		models.DB.Table("post_tags").Where("tag_id = ?", tag.ID).Count(&count)
		tagsWithCount = append(tagsWithCount, TagWithCount{
			Tag:       tag,
			PostCount: int(count),
		})
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "标签",
		"tags":      tagsWithCount,
	}

	if content, err := services.RenderTemplate("tags", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>标签列表页</h1><p>请创建 tags.plush.html 模板</p>")
}

// TagHandler 标签详情页
func TagHandler(c *gin.Context) {
	slug := c.Param("slug")

	var tag models.Tag
	if err := models.DB.Where("slug = ?", slug).First(&tag).Error; err != nil {
		Render404(c)
		return
	}

	var posts []models.Post
	models.DB.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = ?", tag.ID, "published").
		Preload("Category").
		Order("posts.published_at DESC").
		Find(&posts)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": tag.Name,
		"tag":       tag,
		"posts":     posts,
	}

	if content, err := services.RenderTemplate("tag", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// XSS 安全修复
	c.String(http.StatusOK, fmt.Sprintf("<h1>标签: %s</h1><p>请创建 tag.plush.html 模板</p>",
		html.EscapeString(tag.Name)))
}

// ArchivesHandler 文章归档页
func ArchivesHandler(c *gin.Context) {
	var posts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Find(&posts)

	archives := make(map[string][]models.Post)
	for _, post := range posts {
		if post.PublishedAt != nil {
			key := post.PublishedAt.Format("2006年01月")
			archives[key] = append(archives[key], post)
		}
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "文章归档",
		"archives":  archives,
		"posts":     posts,
	}

	if content, err := services.RenderTemplate("archives", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>文章归档</h1><p>请创建 archives.plush.html 模板</p>")
}

// SearchHandler 搜索页
func SearchHandler(c *gin.Context) {
	query := c.Query("q")

	var posts []models.Post
	if query != "" {
		models.DB.Where("status = ? AND (title LIKE ? OR content LIKE ?)",
			"published",
			"%"+query+"%",
			"%"+query+"%").
			Preload("Category").
			Order("published_at DESC").
			Find(&posts)
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "搜索",
		"query":     query,
		"posts":     posts,
		"count":     len(posts),
	}

	if content, err := services.RenderTemplate("search", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// XSS 安全修复
	c.String(http.StatusOK, fmt.Sprintf("<h1>搜索: %s</h1><p>请创建 search.plush.html 模板</p>",
		html.EscapeString(query)))
}