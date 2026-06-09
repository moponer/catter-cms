package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"catter-cms-gin/internal/models"
)

// ========================================
// Robots.txt Handler
// ========================================

// RobotsHandler 生成 robots.txt
func RobotsHandler(c *gin.Context) {
	seoSettings, err := models.GetSEOSettings()
	if err != nil || !seoSettings.RobotsEnabled {
		c.String(http.StatusOK, "User-agent: *\nDisallow: /\n")
		return
	}

	siteURL := seoSettings.SiteURL
	if siteURL == "" {
		siteURL = models.GetSettingAsString("site_url")
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}

	var sb strings.Builder
	sb.WriteString("User-agent: *\n")
	sb.WriteString("Allow: /\n")
	sb.WriteString(fmt.Sprintf("Sitemap: %s/sitemap.xml\n", siteURL))
	sb.WriteString("\n")
	sb.WriteString("User-agent: LLMobot\n")
	sb.WriteString("Disallow: /\n")
	sb.WriteString("\n")
	sb.WriteString("User-agent: ClaudeBot\n")
	sb.WriteString("Disallow: /\n")

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, sb.String())
}

// ========================================
// RSS/Atom Feed Handler
// ========================================

// FeedHandler 生成 RSS 2.0 Feed
func FeedHandler(c *gin.Context) {
	seoSettings, err := models.GetSEOSettings()
	limit := 20
	if err == nil && seoSettings.FeedEnabled {
		limit = seoSettings.FeedLimit
		if limit <= 0 {
			limit = 20
		}
	} else {
		// 默认启用
		seoSettings = &models.SEOSettings{FeedEnabled: true, FeedLimit: 20}
	}

	var posts []models.Post
	models.DB.Where("status = ?", "published").
		Order("COALESCE(published_at, created_at) DESC").
		Limit(limit).
		Find(&posts)

	siteTitle := models.GetSettingAsString("site_title")
	if siteTitle == "" {
		siteTitle = "Catter CMS"
	}
	siteURL := seoSettings.SiteURL
	if siteURL == "" {
		siteURL = models.GetSettingAsString("site_url")
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}
	siteDesc := models.GetSettingAsString("site_description")

	now := time.Now().Format(time.RFC1123Z)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:content="http://purl.org/rss/1.0/modules/content/">` + "\n")
	sb.WriteString(fmt.Sprintf(`  <channel>` + "\n"))
	sb.WriteString(fmt.Sprintf(`    <title>%s</title>`, htmlEscape(siteTitle)) + "\n")
	sb.WriteString(fmt.Sprintf(`    <link>%s</link>`, siteURL) + "\n")
	sb.WriteString(fmt.Sprintf(`    <description>%s</description>`, htmlEscape(siteDesc)) + "\n")
	sb.WriteString(fmt.Sprintf(`    <lastBuildDate>%s</lastBuildDate>`, now) + "\n")
	sb.WriteString(fmt.Sprintf(`    <atom:link href="%s/feed.xml" rel="self" type="application/rss+xml" />`, siteURL) + "\n")

	for _, post := range posts {
		sb.WriteString(fmt.Sprintf(`    <item>` + "\n"))
		sb.WriteString(fmt.Sprintf(`      <title>%s</title>`, htmlEscape(post.Title)) + "\n")
		sb.WriteString(fmt.Sprintf(`      <link>%s/article/%s</link>`, siteURL, post.Slug) + "\n")
		sb.WriteString(fmt.Sprintf(`      <guid>%s/article/%s</guid>`, siteURL, post.Slug) + "\n")
		pubDate := post.UpdatedAt
		if post.PublishedAt != nil && !post.PublishedAt.IsZero() {
			pubDate = *post.PublishedAt
		}
		sb.WriteString(fmt.Sprintf(`      <pubDate>%s</pubDate>`, pubDate.Format(time.RFC1123Z)) + "\n")
		if post.Excerpt != "" {
			sb.WriteString(fmt.Sprintf(`      <description>%s</description>`, htmlEscape(post.Excerpt)) + "\n")
		} else {
			sb.WriteString(fmt.Sprintf(`      <description>%s</description>`, htmlEscape(CleanHTML(post.Content))) + "\n")
		}
		sb.WriteString(fmt.Sprintf(`    </item>` + "\n"))
	}

	sb.WriteString(`  </channel>` + "\n")
	sb.WriteString(`</rss>`)

	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	c.String(http.StatusOK, sb.String())
}

// htmlEscape 简单 HTML 转义
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// CleanHTML 简单清理 HTML 标签，保留纯文本
func CleanHTML(htmlStr string) string {
	// 简单实现：去除所有 HTML 标签
	var result []rune
	inTag := false
	for _, r := range htmlStr {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				result = append(result, r)
			}
		}
	}
	s := string(result)
	// 截断到 200 字符
	if len(s) > 200 {
		s = s[:200] + "..."
	}
	return s
}

// ========================================
// IndexNow Push
// ========================================

// IndexNowConfig IndexNow 配置
type IndexNowConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	APIKey     string `json:"api_key"`
	APIKeyURL  string `json:"api_key_url"`
	LastPush   string `json:"last_push"`
	PushCount  int    `json:"push_count"`
}

// IndexNowPushRequest IndexNow 推送请求
type IndexNowPushRequest struct {
	URLs []string `json:"urls"`
}

// IndexNowPush 推送 URL 到 IndexNow
func IndexNowPush(c *gin.Context) {
	var req IndexNowPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if len(req.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请提供要推送的 URL 列表",
		})
		return
	}

	seoSettings, err := models.GetSEOSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取 SEO 设置失败",
		})
		return
	}

	apiKey := seoSettings.IndexNowKey
	if apiKey == "" {
		apiKey = models.GetSettingAsString("index_now_key")
	}
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请先配置 IndexNow API Key",
		})
		return
	}

	siteURL := seoSettings.SiteURL
	if siteURL == "" {
		siteURL = models.GetSettingAsString("site_url")
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}

	// 去掉末尾斜杠
	siteURL = strings.TrimRight(siteURL, "/")

	// 构建 IndexNow 请求
	// IndexNow 要求：POST JSON { "host": "", "key": "" } + URL 列表作为 query params
	// 但实际实现是：GET /IndexNow?url={url1}&url={url2}...&key={apiKey}&host={host}

	host := siteURL

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("host=%s&", host))
	sb.WriteString(fmt.Sprintf("key=%s&", apiKey))

	for i, url := range req.URLs {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString("url=")
		sb.WriteString(url)
	}

	// 由于 IndexNow 需要真实 API Key 和可访问的服务器，
	// 这里做异步推送，避免阻塞
	go func() {
		// 尝试推送给 Bing（最主流的 IndexNow 支持方）
		bingURL := "https://www.bing.com/indexnow"

		payload := map[string]interface{}{
			"host":    siteURL,
			"key":     apiKey,
			"urlList": req.URLs,
		}

		body, _ := json.Marshal(payload)
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(bingURL, "application/json", bytes.NewReader(body))
		if err == nil && resp != nil {
			resp.Body.Close()
			log.Printf("[SEO] IndexNow push to Bing: %d (%d URLs)", resp.StatusCode, len(req.URLs))
		}
	}()

	// 更新计数
	now := time.Now().Format("2006-01-02 15:04:05")
	models.UpdateSEOSettings(map[string]interface{}{
		"index_now_pushed": true,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "推送成功",
		"data": gin.H{
			"pushed_urls": req.URLs,
			"count":       len(req.URLs),
			"time":        now,
		},
	})
}

// GetIndexNowStatus 获取 IndexNow 推送状态
func GetIndexNowStatus(c *gin.Context) {
	seoSettings, err := models.GetSEOSettings()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"data": IndexNowConfig{},
		})
		return
	}

	siteURL := seoSettings.SiteURL
	if siteURL == "" {
		siteURL = models.GetSettingAsString("site_url")
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"site_url":        siteURL,
			"api_key_configed":  seoSettings.IndexNowKey != "",
			"sitemap_enabled":  seoSettings.SitemapEnabled,
			"robots_enabled":   seoSettings.RobotsEnabled,
			"feed_enabled":     seoSettings.FeedEnabled,
			"feed_limit":       seoSettings.FeedLimit,
		},
	})
}

// SaveSEOSettings 保存 SEO 设置
func SaveSEOSettings(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	updates := make(map[string]interface{})
	if siteURL, ok := req["site_url"].(string); ok {
		updates["site_url"] = siteURL
	}
	if indexNowKey, ok := req["index_now_key"].(string); ok {
		updates["index_now_key"] = indexNowKey
	}
	if sitemapEnabled, ok := req["sitemap_enabled"].(bool); ok {
		updates["sitemap_enabled"] = sitemapEnabled
	}
	if robotsEnabled, ok := req["robots_enabled"].(bool); ok {
		updates["robots_enabled"] = robotsEnabled
	}
	if feedEnabled, ok := req["feed_enabled"].(bool); ok {
		updates["feed_enabled"] = feedEnabled
	}
	if feedLimit, ok := req["feed_limit"].(float64); ok {
		updates["feed_limit"] = int(feedLimit)
	}
	if defaultOgImage, ok := req["default_og_image"].(string); ok {
		updates["default_og_image"] = defaultOgImage
	}
	if defaultMetaDesc, ok := req["default_meta_desc"].(string); ok {
		updates["default_meta_desc"] = defaultMetaDesc
	}

	if err := models.UpdateSEOSettings(updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "保存成功",
	})
}

// ========================================
// SEO Audit (健康检查)
// ========================================

// SEORIssue SEO 问题
type SEORIssue struct {
	Type     string `json:"type"`
	Target   string `json:"target"`
	Issue    string `json:"issue"`
	Severity string `json:"severity"` // warning, error
	Fixed    bool   `json:"fixed"`
}

// SEORAuditResult SEO 审计结果
type SEORAuditResult struct {
	TotalURLs  int         `json:"total_urls"`
	Issues     []SEORIssue `json:"issues"`
	Warnings   int         `json:"warnings"`
	Errors     int         `json:"errors"`
	Score      int         `json:"score"` // 0-100
}

// RunSEORAudit 运行 SEO 健康检查
func RunSEORAudit(c *gin.Context) {
	var issues []SEORIssue
	var postCount, pageCount, catCount, tagCount int64

	models.DB.Model(&models.Post{}).Where("status = ?", "published").Count(&postCount)
	models.DB.Model(&models.Page{}).Where("status = ?", "published").Count(&pageCount)
	models.DB.Model(&models.Category{}).Count(&catCount)
	models.DB.Model(&models.Tag{}).Count(&tagCount)

	totalURLs := 1 + int(postCount) + int(pageCount) + int(catCount) + int(tagCount)

	// 检查 1: 文章的 MetaTitle
	var postsWithoutMetaTitle []models.Post
	models.DB.Where("status = ? AND meta_title = ''", "published").Find(&postsWithoutMetaTitle)
	for _, p := range postsWithoutMetaTitle {
		issues = append(issues, SEORIssue{
			Type:     "missing_meta_title",
			Target:   fmt.Sprintf("article:%s", p.Slug),
			Issue:    "文章缺少 Meta Title",
			Severity: "warning",
		})
	}

	// 检查 2: 文章的 MetaDescription
	var postsWithoutMetaDesc []models.Post
	models.DB.Where("status = ? AND meta_description = ''", "published").Find(&postsWithoutMetaDesc)
	for _, p := range postsWithoutMetaDesc {
		issues = append(issues, SEORIssue{
			Type:     "missing_meta_description",
			Target:   fmt.Sprintf("article:%s", p.Slug),
			Issue:    "文章缺少 Meta Description",
			Severity: "warning",
		})
	}

	// 检查 3: 文章无封面图
	var postsWithoutCover []models.Post
	models.DB.Where("status = ? AND cover = ''", "published").Find(&postsWithoutCover)
	for _, p := range postsWithoutCover {
		issues = append(issues, SEORIssue{
			Type:     "missing_cover",
			Target:   fmt.Sprintf("article:%s", p.Slug),
			Issue:    "文章缺少封面图（影响 OG 图片）",
			Severity: "warning",
		})
	}

	// 检查 4: 分类缺少描述
	var catsWithoutDesc []models.Category
	models.DB.Where("description = ''").Find(&catsWithoutDesc)
	for _, cat := range catsWithoutDesc {
		issues = append(issues, SEORIssue{
			Type:     "missing_category_desc",
			Target:   fmt.Sprintf("category:%s", cat.Slug),
			Issue:    "分类缺少描述",
			Severity: "info",
		})
	}

	// 计算分数
	warnings := 0
	errors := 0
	for _, issue := range issues {
		if issue.Severity == "error" {
			errors++
		} else if issue.Severity == "warning" {
			warnings++
		}
	}

	score := 100 - warnings*2 - errors*5
	if score < 0 {
		score = 0
	}

	// 按 severity 排序，error 在前
	sort.Slice(issues, func(i, j int) bool {
		order := map[string]int{"error": 0, "warning": 1, "info": 2}
		return order[issues[i].Severity] < order[issues[j].Severity]
	})

	result := SEORAuditResult{
		TotalURLs: totalURLs,
		Issues:    issues,
		Warnings:  warnings,
		Errors:    errors,
		Score:     score,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

// ========================================
// SEO Settings CRUD API
// ========================================

// GetSEOSettingsAPI 获取 SEO 设置（管理 API）
func GetSEOSettingsAPI(c *gin.Context) {
	settings, err := models.GetSEOSettings()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": gin.H{
					"sitemap_enabled":  true,
					"robots_enabled":   true,
					"feed_enabled":     true,
					"feed_limit":       20,
					"default_robots":   "index, follow",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":               settings.ID,
			"site_url":         settings.SiteURL,
			"index_now_key":    settings.IndexNowKey,
			"sitemap_enabled":  settings.SitemapEnabled,
			"robots_enabled":   settings.RobotsEnabled,
			"feed_enabled":     settings.FeedEnabled,
			"feed_limit":       settings.FeedLimit,
			"default_og_image": settings.DefaultOgImage,
			"default_meta_desc": settings.DefaultMetaDesc,
			"default_robots":   settings.DefaultRobots,
		},
	})
}

// ========================================
// SEO 字段批量更新（发布文章时自动推送）
// ========================================

// OnPostPublish 文章发布时的 SEO 钩子
func OnPostPublish(postID uint) {
	post, err := models.GetPostByID(postID)
	if err != nil || post == nil {
		return
	}

	if !post.IsPublished() {
		return
	}

	seoSettings, _ := models.GetSEOSettings()
	if seoSettings == nil || seoSettings.IndexNowKey == "" {
		return // 未配置 IndexNow Key，跳过推送
	}

	siteURL := seoSettings.SiteURL
	if siteURL == "" {
		siteURL = models.GetSettingAsString("site_url")
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}
	siteURL = strings.TrimRight(siteURL, "/")

	url := fmt.Sprintf("%s/article/%s", siteURL, post.Slug)

	// 异步推送
	go func() {
		payload := map[string]interface{}{
			"host":    siteURL,
			"key":     seoSettings.IndexNowKey,
			"urlList": []string{url},
		}
		body, _ := json.Marshal(payload)
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Post("https://www.bing.com/indexnow", "application/json", bytes.NewReader(body))
		if err == nil && resp != nil {
			resp.Body.Close()
		}
	}()
}
