package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"catter-cms-gin/internal/models"
)

// ========================================
// Robots.txt Handler
// ========================================

// pushClient 推送用 HTTP 客户端（包级单例，复用连接池）
var pushClient = &http.Client{Timeout: 10 * time.Second}

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
	BaiduVerify string `json:"baidu_verify"`
	GoogleVerify string `json:"google_verify"`
	HasBaiduVerify bool `json:"has_baidu_verify"`
	HasGoogleVerify bool `json:"has_google_verify"`
}

// IndexNowPushRequest IndexNow 推送请求
type IndexNowPushRequest struct {
	URLs []string `json:"urls"`
}

// pushResult 单次推送结果
type pushResult struct {
	name    string
	success bool
	message string
}

// IndexNowPush 推送 URL 到搜索引擎（IndexNow + 百度/360/搜狗主动推送）
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

	siteURL := seoSettings.SiteURL
	if siteURL == "" {
		siteURL = models.GetSettingAsString("site_url")
	}
	if siteURL == "" {
		siteURL = "http://localhost:3000"
	}
	siteURL = strings.TrimRight(siteURL, "/")

	urls := req.URLs

	// 异步推送，不阻塞响应
	go func() {
		results := pushToAllEngines(siteURL, urls, seoSettings)

		// 汇总日志
		successCount := 0
		for _, r := range results {
			if r.success {
				successCount++
			}
			log.Printf("[SEO] Push %s: %s", r.name, r.message)
		}
		log.Printf("[SEO] Push complete: %d/%d engines succeeded, %d URLs", successCount, len(results), len(urls))

		// 更新计数
		now := time.Now().Format("2006-01-02 15:04:05")
		models.UpdateSEOSettings(map[string]interface{}{
			"index_now_pushed": true,
			"pushed_at":        now,
		})
	}()

	now := time.Now().Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "推送请求已提交，正在后台处理",
		"data": gin.H{
			"pushed_urls": urls,
			"count":       len(urls),
			"time":        now,
		},
	})
}

// pushToAllEngines 并发推送到所有搜索引擎
func pushToAllEngines(siteURL string, urls []string, seoSettings *models.SEOSettings) []pushResult {
	var results []pushResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 1. IndexNow (Bing) — 需要 API Key
	if seoSettings.IndexNowKey != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			payload := map[string]interface{}{
				"host":    siteURL,
				"key":     seoSettings.IndexNowKey,
				"urlList": urls,
			}
			body, _ := json.Marshal(payload)
			resp, err := pushClient.Post("https://www.bing.com/indexnow", "application/json", bytes.NewReader(body))
			msg := ""
			ok := false
			if err != nil {
				msg = fmt.Sprintf("错误: %v", err)
			} else if resp != nil {
				resp.Body.Close()
				msg = fmt.Sprintf("HTTP %d", resp.StatusCode)
				ok = (resp.StatusCode >= 200 && resp.StatusCode < 300)
			} else {
				msg = "无响应"
			}
			mu.Lock()
			results = append(results, pushResult{name: "Bing (IndexNow)", success: ok, message: msg})
			mu.Unlock()
		}()
	}

	// 2. 百度站长主动推送（每批最多 1000 URL）
	if seoSettings.BaiduPushToken != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 百度 API 要求 POST 纯文本，每行一个 URL
			baiduURL := fmt.Sprintf("https://data.zz.baidu.com/urls?site=%s&token=%s", siteURL, seoSettings.BaiduPushToken)

			// 分批推送（百度限制每批 1000 条）
			batchSize := 1000
			for i := 0; i < len(urls); i += batchSize {
				end := i + batchSize
				if end > len(urls) {
					end = len(urls)
				}
				batch := urls[i:end]
				bodyStr := strings.Join(batch, "\n")

				resp, err := pushClient.Post(baiduURL, "text/plain", bytes.NewReader([]byte(bodyStr)))
				msg := ""
				ok := false
				if err != nil {
					msg = fmt.Sprintf("错误: %v", err)
			} else if resp != nil {
				resp.Body.Close()
					msg = fmt.Sprintf("HTTP %d (批次 %d-%d)", resp.StatusCode, i+1, end)
					ok = (resp.StatusCode >= 200 && resp.StatusCode < 300)
				} else {
					msg = "无响应"
				}
				mu.Lock()
				results = append(results, pushResult{name: "百度", success: ok, message: msg})
				mu.Unlock()
			}
		}()
	}

	// 3. 360 站长推送
	if seoSettings.ThreeSixZeroPushToken != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 360 推送 API：POST 纯文本 URL 列表
			threeSixZeroURL := fmt.Sprintf("https://data.zhanzhang.360.cn/?token=%s", seoSettings.ThreeSixZeroPushToken)
			bodyStr := strings.Join(urls, "\n")

			resp, err := pushClient.Post(threeSixZeroURL, "text/plain", bytes.NewReader([]byte(bodyStr)))
			msg := ""
			ok := false
			if err != nil {
				msg = fmt.Sprintf("错误: %v", err)
			} else if resp != nil {
				resp.Body.Close()
				msg = fmt.Sprintf("HTTP %d", resp.StatusCode)
				ok = (resp.StatusCode >= 200 && resp.StatusCode < 300)
			} else {
				msg = "无响应"
			}
			mu.Lock()
			results = append(results, pushResult{name: "360", success: ok, message: msg})
			mu.Unlock()
		}()
	}

	// 4. 搜狗站长推送
	if seoSettings.SogouPushToken != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 搜狗主动推送 API
			sogouURL := "https://www.sogou.com/web/FrequentSubmit"

			// 搜狗需要 POST form 数据
			formData := fmt.Sprintf("site=%s&token=%s&urls=%s", siteURL, seoSettings.SogouPushToken, strings.Join(urls, "\n"))
			resp, err := pushClient.Post(sogouURL, "application/x-www-form-urlencoded", bytes.NewReader([]byte(formData)))
			msg := ""
			ok := false
			if err != nil {
				msg = fmt.Sprintf("错误: %v", err)
			} else if resp != nil {
				resp.Body.Close()
				msg = fmt.Sprintf("HTTP %d", resp.StatusCode)
				ok = (resp.StatusCode >= 200 && resp.StatusCode < 300)
			} else {
				msg = "无响应"
			}
			mu.Lock()
			results = append(results, pushResult{name: "搜狗", success: ok, message: msg})
			mu.Unlock()
		}()
	}

	wg.Wait()
	return results
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

	// 计算推送次数（基于 IndexNowPushed 状态）
	pushCount := 0
	if seoSettings.IndexNowPushed {
		pushCount = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"site_url":              siteURL,
			"api_key_configed":      seoSettings.IndexNowKey != "",
			"sitemap_enabled":       seoSettings.SitemapEnabled,
			"robots_enabled":        seoSettings.RobotsEnabled,
			"feed_enabled":          seoSettings.FeedEnabled,
			"feed_limit":            seoSettings.FeedLimit,
			"baidu_push_configed":   seoSettings.BaiduPushToken != "",
			"three_six_zero_push_configed": seoSettings.ThreeSixZeroPushToken != "",
			"sogou_push_configed":   seoSettings.SogouPushToken != "",
			"push_count":            pushCount,
			"last_push":             "", // 可以后续扩展从数据库读取
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
	if defaultRobots, ok := req["default_robots"].(string); ok {
		updates["default_robots"] = defaultRobots
	}
	if baiduToken, ok := req["baidu_push_token"].(string); ok {
		updates["baidu_push_token"] = baiduToken
	}
	if threeSixZeroToken, ok := req["three_six_zero_push_token"].(string); ok {
		updates["three_six_zero_push_token"] = threeSixZeroToken
	}
	if sogouToken, ok := req["sogou_push_token"].(string); ok {
		updates["sogou_push_token"] = sogouToken
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
	models.DB.Select("id, slug").Where("status = ? AND meta_title = ''", "published").Find(&postsWithoutMetaTitle)
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
	models.DB.Select("id, slug").Where("status = ? AND meta_description = ''", "published").Find(&postsWithoutMetaDesc)
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
	models.DB.Select("id, slug").Where("status = ? AND cover = ''", "published").Find(&postsWithoutCover)
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
	models.DB.Select("id, slug").Where("description = ''").Find(&catsWithoutDesc)
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

	// 百分比算法：基于总 URL 数和问题数的比例
	totalItems := totalURLs + warnings + errors*2
	score := 100
	if totalItems > 0 {
		score = int(float64(totalURLs) / float64(totalItems) * 100)
	}
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
			"id":                    settings.ID,
			"site_url":              settings.SiteURL,
			"index_now_key":         settings.IndexNowKey,
			"sitemap_enabled":       settings.SitemapEnabled,
			"robots_enabled":        settings.RobotsEnabled,
			"feed_enabled":          settings.FeedEnabled,
			"feed_limit":            settings.FeedLimit,
			"default_og_image":      settings.DefaultOgImage,
			"default_meta_desc":     settings.DefaultMetaDesc,
			"default_robots":        settings.DefaultRobots,
			"baidu_push_token":      settings.BaiduPushToken,
			"three_six_zero_push_token": settings.ThreeSixZeroPushToken,
			"sogou_push_token":      settings.SogouPushToken,
		},
	})
}

// ========================================
// SEO 字段批量更新（发布文章时自动推送）
// ========================================

// OnPostPublish 文章发布时的 SEO 钩子（自动推送到所有配置的搜索引擎）
func OnPostPublish(postID uint) {
	post, err := models.GetPostByID(postID)
	if err != nil || post == nil {
		return
	}

	if !post.IsPublished() {
		return
	}

	seoSettings, _ := models.GetSEOSettings()
	if seoSettings == nil {
		return
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
	urls := []string{url}

	// 检查是否有任意引擎可推
	hasAnyEngine := seoSettings.IndexNowKey != "" || seoSettings.BaiduPushToken != "" ||
		seoSettings.ThreeSixZeroPushToken != "" || seoSettings.SogouPushToken != ""
	if !hasAnyEngine {
		return
	}

	// 异步推送到所有已配置的引擎
	go pushToAllEngines(siteURL, urls, seoSettings)
}
