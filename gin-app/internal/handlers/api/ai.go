package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"

	"catter-cms-gin/internal/models"
)

// ============================================
// Key 管理 API（JWT 认证）
// ============================================

// CreateAIKeyRequest 创建 AI Key 请求
type CreateAIKeyRequest struct {
	Name        string   `json:"name" binding:"required"`
	Permissions []string `json:"permissions"` // ["read", "write", "admin"]
	Scopes      []string `json:"scopes"`      // ["posts", "pages", ...]
	ExpiresIn   *int     `json:"expires_in"`  // 有效天数，nil 表示永不过期
}

// CreateAIKey 创建 AI Key
func CreateAIKey(c *gin.Context) {
	var req CreateAIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 获取当前用户
	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}
	currentUser := user.(*models.User)

	// 验证权限
	validPerms := map[string]bool{"read": true, "write": true, "admin": true}
	perms := req.Permissions
	if len(perms) == 0 {
		perms = []string{"read"}
	}
	for _, p := range perms {
		if !validPerms[p] {
			BadRequest(c, "无效的权限值: "+p)
			return
		}
	}

	// 验证范围
	validScopes := map[string]bool{
		"posts": true, "pages": true, "categories": true, "tags": true,
		"media": true, "comments": true, "settings": true, "all": true,
	}
	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = []string{"posts", "pages", "categories", "tags", "media", "comments", "settings"}
	}
	for _, s := range scopes {
		if !validScopes[s] {
			BadRequest(c, "无效的资源范围: "+s)
			return
		}
	}

	// 序列化权限和范围
	permsJSON, _ := json.Marshal(perms)
	scopesJSON, _ := json.Marshal(scopes)

	// 计算过期时间
	var expiresAt time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expiresAt = time.Now().Add(time.Duration(*req.ExpiresIn) * 24 * time.Hour)
	} else {
		expiresAt = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}

	aiToken := models.AIToken{
		Name:        req.Name,
		UserID:      currentUser.ID,
		Permissions: string(permsJSON),
		Scopes:      string(scopesJSON),
		ExpiresAt:   expiresAt,
		IsActive:    true,
	}

	if err := models.DB.Create(&aiToken).Error; err != nil {
		InternalError(c, "创建 Key 失败")
		return
	}

	// 返回完整 token（仅此一次可见）
	SuccessWithMessage(c, "Key 创建成功", gin.H{
		"id":            aiToken.ID,
		"name":          aiToken.Name,
		"token":         aiToken.Token, // 完整 token，仅创建时返回
		"token_prefix":  aiToken.TokenPrefix,
		"permissions":   aiToken.GetPermissions(),
		"scopes":        aiToken.GetScopes(),
		"expires_at":    aiToken.ExpiresAt,
		"created_at":    aiToken.CreatedAt,
	})
}

// ListAIKeys 获取 AI Key 列表
func ListAIKeys(c *gin.Context) {
	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}
	currentUser := user.(*models.User)

	var tokens []models.AIToken
	if err := models.DB.Where("user_id = ?", currentUser.ID).
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		InternalError(c, "获取 Key 列表失败")
		return
	}

	SuccessResponse(c, tokens)
}

// DeleteAIKey 撤销（删除）AI Key
func DeleteAIKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的ID")
		return
	}

	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}
	currentUser := user.(*models.User)

	result := models.DB.Where("id = ? AND user_id = ?", id, currentUser.ID).Delete(&models.AIToken{})
	if result.Error != nil {
		InternalError(c, "删除失败")
		return
	}
	if result.RowsAffected == 0 {
		NotFound(c, "Key 不存在")
		return
	}

	SuccessWithMessage(c, "Key 已撤销", nil)
}

// GetAIKeyLogs 获取 AI Key 审计日志
func GetAIKeyLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的ID")
		return
	}

	// 验证权限：只能查看自己的 key 日志
	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}
	currentUser := user.(*models.User)

	// 验证 key 属于当前用户
	var token models.AIToken
	if err := models.DB.Where("id = ? AND user_id = ?", id, currentUser.ID).First(&token).Error; err != nil {
		NotFound(c, "Key 不存在")
		return
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var logs []models.AIAuditLog
	var total int64
	query := models.DB.Model(&models.AIAuditLog{}).Where("token_id = ?", id)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)

	SuccessResponse(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ============================================
// 公开端点（无需认证）
// ============================================

// isBrowserRequest 检测是否为浏览器直接访问
func isBrowserRequest(c *gin.Context) bool {
	accept := c.GetHeader("Accept")
	return strings.Contains(accept, "text/html")
}

// respondAIOnly 对浏览器访问返回友好提示
func respondAIOnly(c *gin.Context, contentType string) {
	html := "<!DOCTYPE html><html><head><meta charset='utf-8'><title>AI Only</title>" +
		"<style>body{font-family:system-ui,sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;background:#f5f5f5;color:#333}" +
		".box{text-align:center;padding:48px;background:#fff;border-radius:12px;box-shadow:0 2px 12px rgba(0,0,0,.08);max-width:480px}" +
		"h1{font-size:20px;margin-bottom:8px}p{color:#666;line-height:1.6;font-size:14px;margin:8px 0}" +
		"code{background:#f0f0f0;padding:2px 6px;border-radius:4px;font-size:13px}</style></head>" +
		"<body><div class='box'>" +
		"<h1>AI 专用接口</h1>" +
		"<p>此接口仅供 AI 客户端（如 Cursor、WorkBuddy 等）通过程序访问，不面向人类用户直接使用。</p>" +
		"<p>如需配置 AI 集成，请前往管理后台的 <strong>AI 集成</strong> 页面按引导操作。</p>" +
		"<p>返回内容类型：<code>" + contentType + "</code></p>" +
		"</div></body></html>"
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("X-Robots-Tag", "noindex, nofollow")
	c.String(200, html)
}

// GetAISkill 获取 AI Skill 文档（读取静态文件）
func GetAISkill(c *gin.Context) {
	if isBrowserRequest(c) {
		respondAIOnly(c, "Skill Document (Markdown)")
		return
	}
	// 尝试从 deploy/public/ai/skill.md 读取，否则从源码 public/ai/skill.md 读取
	paths := []string{
		"public/ai/skill.md",
		"../public/ai/skill.md",
	}

	var content []byte
	var err error
	for _, p := range paths {
		content, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}

	if err != nil || len(content) == 0 {
		c.String(500, "Skill document not found")
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	// 标记 AI 专用
	c.Header("X-Content-Type", "ai-skill-document")
	c.String(200, string(content))
}

// GetAISchema 获取 OpenAPI Schema
func GetAISchema(c *gin.Context) {
	if isBrowserRequest(c) {
		respondAIOnly(c, "OpenAPI 3.0 Schema (JSON)")
		return
	}

	cmsURL := getCMSURL(c)

	schema := gin.H{
		"openapi": "3.0.3",
		"info": gin.H{
			"title":       "Catter CMS API",
			"description": "Catter CMS AI 集成 API",
			"version":     "1.0.0",
		},
		"servers": []gin.H{
			{"url": cmsURL, "description": "Catter CMS"},
		},
		"components": gin.H{
			"securitySchemes": gin.H{
				"bearerAuth": gin.H{
					"type":         "http",
					"scheme":       "bearer",
					"description":  "Catter CMS Key (cat_ 前缀)",
					"bearerFormat": "cat_xxx",
				},
			},
			"schemas": gin.H{
				"Error": gin.H{
					"type": "object",
					"properties": gin.H{
						"code":    gin.H{"type": "integer", "example": 400},
						"message": gin.H{"type": "string", "example": "参数错误"},
					},
				},
				"PaginatedResponse": gin.H{
					"type": "object",
					"properties": gin.H{
						"list":      gin.H{"type": "array"},
						"total":     gin.H{"type": "integer"},
						"page":      gin.H{"type": "integer"},
						"page_size": gin.H{"type": "integer"},
					},
				},
				"Post": gin.H{
					"type": "object",
					"properties": gin.H{
						"id":              gin.H{"type": "integer"},
						"title":           gin.H{"type": "string"},
						"slug":            gin.H{"type": "string"},
						"content":         gin.H{"type": "string"},
						"excerpt":         gin.H{"type": "string"},
						"cover":           gin.H{"type": "string"},
						"status":          gin.H{"type": "string", "enum": []string{"draft", "published", "archived"}},
						"category_id":     gin.H{"type": "integer"},
						"views":           gin.H{"type": "integer"},
						"created_at":      gin.H{"type": "string", "format": "date-time"},
						"updated_at":      gin.H{"type": "string", "format": "date-time"},
					},
				},
				"Page": gin.H{
					"type": "object",
					"properties": gin.H{
						"id":         gin.H{"type": "integer"},
						"title":      gin.H{"type": "string"},
						"slug":       gin.H{"type": "string"},
						"content":    gin.H{"type": "string"},
						"status":     gin.H{"type": "string"},
						"template":   gin.H{"type": "string"},
						"sort_order": gin.H{"type": "integer"},
						"created_at": gin.H{"type": "string", "format": "date-time"},
						"updated_at": gin.H{"type": "string", "format": "date-time"},
					},
				},
				"Category": gin.H{
					"type": "object",
					"properties": gin.H{
						"id":        gin.H{"type": "integer"},
						"name":      gin.H{"type": "string"},
						"slug":      gin.H{"type": "string"},
						"parent_id": gin.H{"type": "integer"},
						"sort":      gin.H{"type": "integer"},
					},
				},
				"Tag": gin.H{
					"type": "object",
					"properties": gin.H{
						"id":   gin.H{"type": "integer"},
						"name": gin.H{"type": "string"},
						"slug": gin.H{"type": "string"},
					},
				},
			},
		},
		"security": []gin.H{
			{"bearerAuth": []string{}},
		},
		"paths": buildOpenAPIPaths(),
	}

	c.JSON(200, schema)
}

// GetThemeDevSkill 获取主题开发 Skill 文档（读取静态文件）
func GetThemeDevSkill(c *gin.Context) {
	if isBrowserRequest(c) {
		respondAIOnly(c, "Theme Development Skill Document (Markdown)")
		return
	}
	// 尝试从 deploy/public/ai/themes-skill.md 读取，否则从源码 public/ai/themes-skill.md 读取
	paths := []string{
		"public/ai/themes-skill.md",
		"../public/ai/themes-skill.md",
	}

	var content []byte
	var err error
	for _, p := range paths {
		content, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}

	if err != nil || len(content) == 0 {
		c.String(500, "Theme development skill document not found")
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("X-Content-Type", "ai-theme-skill-document")
	c.String(200, string(content))
}

// GetAISetupCommands 获取 AI 集成设置命令（JWT 认证）
func GetAISetupCommands(c *gin.Context) {
	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}
	currentUser := user.(*models.User)

	// 查找用户的活跃 key
	var activeKey models.AIToken
	if err := models.DB.Where("user_id = ? AND is_active = ?", currentUser.ID, true).
		Order("created_at DESC").
		First(&activeKey).Error; err != nil {
		// 没有活跃 key，返回提示
		SuccessResponse(c, gin.H{
			"cms_url":    getCMSURL(c),
			"has_key":    false,
			"message":    "请先创建一个 Catter CMS Key",
		})
		return
	}

	if activeKey.IsExpired() {
		SuccessResponse(c, gin.H{
			"cms_url":    getCMSURL(c),
			"has_key":    false,
			"message":    "当前 Key 已过期，请创建新的 Key",
		})
		return
	}

	cmsURL := getCMSURL(c)

	SuccessResponse(c, gin.H{
		"cms_url":     cmsURL,
		"has_key":     true,
		"key_name":    activeKey.Name,
		"key_prefix":  activeKey.TokenPrefix,
		"permissions": activeKey.GetPermissions(),
		"scopes":      activeKey.GetScopes(),
		"expires_at":  activeKey.ExpiresAt,
		"commands": gin.H{
			"windows": fmt.Sprintf("[Environment]::SetEnvironmentVariable(\"CATTER_CMS_URL\", \"%s\", \"User\"); [Environment]::SetEnvironmentVariable(\"CATTER_CMS_KEY\", \"%s\", \"User\")", cmsURL, activeKey.Token),
			"mac":     fmt.Sprintf("echo 'export CATTER_CMS_URL=\"%s\"' >> ~/.zshrc && echo 'export CATTER_CMS_KEY=\"%s\"' >> ~/.zshrc && source ~/.zshrc", cmsURL, activeKey.Token),
			"linux":   fmt.Sprintf("echo 'export CATTER_CMS_URL=\"%s\"' >> ~/.bashrc && echo 'export CATTER_CMS_KEY=\"%s\"' >> ~/.bashrc && source ~/.bashrc", cmsURL, activeKey.Token),
		},
		"skill_prompt": fmt.Sprintf(`Catter CMS 集成：当用户要求管理网站内容（文章、页面、分类、标签、媒体、评论、设置）时，请按以下步骤操作：

【环境配置】
- CMS 地址：[Environment]::GetEnvironmentVariable("CATTER_CMS_URL", "User")
- API Key：[Environment]::GetEnvironmentVariable("CATTER_CMS_KEY", "User")

【首次使用 - 安装 Skill】
1. 检查记忆中是否已有 "Catter CMS Skill 已安装"
2. 如未安装，调用 save_memory 保存以下信息：
   "Catter CMS Skill 已安装。API认证：请求头 Authorization: Bearer {CATTER_CMS_KEY}。操作原则：先读后写、删除确认、默认草稿。"
3. 告知用户安装完成

【执行操作】
1. 使用 run_shell_command 通过 PowerShell Socket 发起 API 请求
2. 请求地址：{CATTER_CMS_URL}/api/v1/{endpoint}
3. 请求头：Authorization: Bearer {CATTER_CMS_KEY}
4. 参考 Skill 文档：{CATTER_CMS_URL}/api/ai/skill（可选，获取完整 API 规范）

【操作原则】
- 先读后写：创建/更新前先 GET 查看现有数据
- 删除确认：删除前向用户展示内容并确认
- 默认草稿：创建文章 status 默认为 draft`, cmsURL),
	})
}

// ============================================
// AI v1 API 代理（使用 AI Key 认证）
// 将现有 CMS API 暴露给 AI 客户端
// ============================================

// AIListPosts AI 文章列表
func AIListPosts(c *gin.Context) {
	var posts []models.Post
	query := models.DB.Preload("Category").Preload("Author").Preload("Tags")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	if keyword := c.Query("keyword"); keyword != "" {
		escaped := escapeLikeWildcards(keyword)
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+escaped+"%", "%"+escaped+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Model(&models.Post{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&posts).Error; err != nil {
		InternalError(c, "获取文章列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// AIGetPost AI 获取文章详情
func AIGetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的文章ID")
		return
	}

	var post models.Post
	if err := models.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, id).Error; err != nil {
		NotFound(c, "文章不存在")
		return
	}

	SuccessResponse(c, post)
}

// AICreatePost AI 创建文章
func AICreatePost(c *gin.Context) {
	var req struct {
		Title      string   `json:"title"`
		Slug       string   `json:"slug"`
		Content    string   `json:"content"`
		Excerpt    string   `json:"excerpt"`
		Cover      string   `json:"cover"`
		Status     string   `json:"status"`
		CategoryID *uint    `json:"category_id"`
		TagNames   []string `json:"tag_names"`
		ContentType string  `json:"content_type"` // "html" 或 "markdown"，默认 "html"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}
	if req.Title == "" {
		BadRequest(c, "文章标题不能为空")
		return
	}

	// 获取 AI token 关联的用户
	aiToken, _ := c.Get("ai_token")
	token := aiToken.(*models.AIToken)

	// 处理内容格式转换
	content := req.Content
	if req.ContentType == "markdown" {
		content = markdownToHTML(req.Content)
	}

	post := models.Post{
		Title:      req.Title,
		Content:    content,
		Excerpt:    req.Excerpt,
		Cover:      req.Cover,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		AuthorID:   token.UserID,
	}
	if post.Status == "" {
		post.Status = "draft"
	}
	if post.Slug == "" {
		// 简单 slug 生成
		post.Slug = generateSimpleSlug(req.Title)
	}

	if err := models.DB.Create(&post).Error; err != nil {
		InternalError(c, "创建文章失败")
		return
	}

	// 同步标签（如果有）
	if len(req.TagNames) > 0 {
		syncSimpleTags(&post, req.TagNames)
	}

	models.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, post.ID)
	SuccessWithMessage(c, "文章创建成功", post)
}

// AIUpdatePost AI 更新文章
func AIUpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的文章ID")
		return
	}

	var post models.Post
	if err := models.DB.First(&post, id).Error; err != nil {
		NotFound(c, "文章不存在")
		return
	}

	var req struct {
		Title      string   `json:"title"`
		Slug       string   `json:"slug"`
		Content    string   `json:"content"`
		Excerpt    string   `json:"excerpt"`
		Cover      string   `json:"cover"`
		Status     string   `json:"status"`
		CategoryID *uint    `json:"category_id"`
		TagNames   []string `json:"tag_names"`
		ContentType string  `json:"content_type"` // "html" 或 "markdown"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 处理内容格式转换
	content := req.Content
	if req.ContentType == "markdown" {
		content = markdownToHTML(req.Content)
	}

	updates := map[string]interface{}{
		"title":       req.Title,
		"slug":        req.Slug,
		"content":     content,
		"excerpt":     req.Excerpt,
		"cover":       req.Cover,
		"status":      req.Status,
		"category_id": req.CategoryID,
	}

	// 如果状态为 published 且没有发布时间，设置发布时间
	if req.Status == "published" && post.PublishedAt == nil {
		now := time.Now()
		updates["published_at"] = &now
	}

	if err := models.DB.Model(&post).Updates(updates).Error; err != nil {
		InternalError(c, "更新文章失败")
		return
	}

	// 同步标签
	if req.TagNames != nil {
		syncSimpleTags(&post, req.TagNames)
	}

	models.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, post.ID)
	SuccessWithMessage(c, "文章更新成功", post)
}

// markdownToHTML 将 Markdown 转换为 HTML
func markdownToHTML(content string) string {
	if content == "" {
		return ""
	}
	// 检测是否已经是 HTML 内容
	trimmed := strings.TrimSpace(content)
	if strings.HasPrefix(trimmed, "<") {
		return content
	}
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(content), &buf); err != nil {
		return content
	}
	return buf.String()
}

// AIDeletePost AI 删除文章
func AIDeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的文章ID")
		return
	}

	var post models.Post
	if err := models.DB.First(&post, id).Error; err != nil {
		NotFound(c, "文章不存在")
		return
	}

	if err := models.DB.Delete(&post).Error; err != nil {
		InternalError(c, "删除文章失败")
		return
	}

	SuccessWithMessage(c, "文章删除成功", nil)
}

// AIListPages AI 页面列表
func AIListPages(c *gin.Context) {
	var pages []models.Page
	query := models.DB

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	query.Model(&models.Page{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&pages).Error; err != nil {
		InternalError(c, "获取页面列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":      pages,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// AIGetPage AI 获取页面详情
func AIGetPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的页面ID")
		return
	}

	var page models.Page
	if err := models.DB.First(&page, id).Error; err != nil {
		NotFound(c, "页面不存在")
		return
	}

	SuccessResponse(c, page)
}

// AIListCategories AI 分类列表
func AIListCategories(c *gin.Context) {
	var categories []models.Category
	if err := models.DB.Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		InternalError(c, "获取分类列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":  categories,
		"total": len(categories),
	})
}

// AIListTags AI 标签列表
func AIListTags(c *gin.Context) {
	var tags []models.Tag
	query := models.DB

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	var total int64
	query.Model(&models.Tag{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tags).Error; err != nil {
		InternalError(c, "获取标签列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":      tags,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// AIListMedia AI 媒体列表
func AIListMedia(c *gin.Context) {
	var media []models.Media
	query := models.DB

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.Media{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&media).Error; err != nil {
		InternalError(c, "获取媒体列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":      media,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// AIListComments AI 评论列表
func AIListComments(c *gin.Context) {
	var comments []models.Comment
	query := models.DB.Preload("Post").Preload("User")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.Comment{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		InternalError(c, "获取评论列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":      comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// AIGetSettings AI 获取系统设置
func AIGetSettings(c *gin.Context) {
	var settings []models.Setting
	if err := models.DB.Find(&settings).Error; err != nil {
		InternalError(c, "获取设置失败")
		return
	}

	result := make(map[string]interface{})
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	SuccessResponse(c, result)
}

// ============================================
// 辅助函数
// ============================================

// getCMSURL 获取 CMS 地址
func getCMSURL(c *gin.Context) string {
	// 优先从请求中获取
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost:3000"
	}
	return scheme + "://" + host
}

// generateSimpleSlug 简单 slug 生成（AI 端使用，不依赖 services 包）
func generateSimpleSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// 只保留字母、数字和连字符
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	res := result.String()
	if res == "" {
		res = fmt.Sprintf("post-%d", time.Now().Unix())
	}
	return res
}

// syncSimpleTags 简单标签同步（AI 端使用）
func syncSimpleTags(post *models.Post, tagNames []string) {
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var tag models.Tag
		err := models.DB.Where("name = ?", name).First(&tag).Error
		if err != nil {
			tag = models.Tag{
				Name: name,
				Slug: generateSimpleSlug(name),
			}
			models.DB.Create(&tag)
		}
		models.DB.Exec("INSERT OR IGNORE INTO post_tags (post_id, tag_id) VALUES (?, ?)", post.ID, tag.ID)
	}
}

// buildOpenAPIPaths 构建 OpenAPI paths
func buildOpenAPIPaths() gin.H {
	paths := gin.H{}

	// 文章
	paths["/api/v1/posts"] = gin.H{
		"get":  buildPathGet("获取文章列表", "PaginatedResponse", "文章列表"),
		"post": buildPathPost("创建文章", "Post"),
	}
	paths["/api/v1/posts/{id}"] = gin.H{
		"get": buildPathGet("获取文章详情", "Post", ""),
		"put": buildPathPut("更新文章", "Post"),
		"delete": buildPathDelete("删除文章"),
	}

	// 页面
	paths["/api/v1/pages"] = gin.H{
		"get":  buildPathGet("获取页面列表", "PaginatedResponse", "页面列表"),
		"post": buildPathPost("创建页面", "Page"),
	}
	paths["/api/v1/pages/{id}"] = gin.H{
		"get":    buildPathGet("获取页面详情", "Page", ""),
		"put":    buildPathPut("更新页面", "Page"),
		"delete": buildPathDelete("删除页面"),
	}

	// 分类
	paths["/api/v1/categories"] = gin.H{
		"get":  buildPathGet("获取分类列表", "#/components/schemas/Category", "分类列表"),
		"post": buildPathPost("创建分类", "#/components/schemas/Category"),
	}
	paths["/api/v1/categories/{id}"] = gin.H{
		"put":    buildPathPut("更新分类", "#/components/schemas/Category"),
		"delete": buildPathDelete("删除分类"),
	}

	// 标签
	paths["/api/v1/tags"] = gin.H{
		"get":  buildPathGet("获取标签列表", "#/components/schemas/Tag", "标签列表"),
		"post": buildPathPost("创建标签", "#/components/schemas/Tag"),
	}
	paths["/api/v1/tags/{id}"] = gin.H{
		"put":    buildPathPut("更新标签", "#/components/schemas/Tag"),
		"delete": buildPathDelete("删除标签"),
	}

	// 媒体
	paths["/api/v1/media"] = gin.H{
		"get": buildPathGet("获取媒体列表", "PaginatedResponse", "媒体列表"),
	}
	paths["/api/v1/media/{id}"] = gin.H{
		"delete": buildPathDelete("删除媒体"),
	}

	// 评论
	paths["/api/v1/comments"] = gin.H{
		"get": buildPathGet("获取评论列表", "PaginatedResponse", "评论列表"),
	}
	paths["/api/v1/comments/{id}"] = gin.H{
		"delete": buildPathDelete("删除评论"),
	}

	// 设置
	paths["/api/v1/settings"] = gin.H{
		"get": buildPathGet("获取系统设置", "object", ""),
		"put": buildPathPut("更新系统设置", "object"),
	}

	return paths
}

func buildPathGet(summary, schemaRef, description string) gin.H {
	p := gin.H{
		"summary":     summary,
		"description": description,
		"parameters": []gin.H{
			{"name": "page", "in": "query", "schema": gin.H{"type": "integer", "default": 1}},
			{"name": "page_size", "in": "query", "schema": gin.H{"type": "integer", "default": 10}},
			{"name": "keyword", "in": "query", "schema": gin.H{"type": "string"}, "description": "搜索关键词"},
		},
		"responses": gin.H{
			"200": gin.H{"description": "成功"},
		},
	}
	return p
}

func buildPathPost(summary, schemaRef string) gin.H {
	return gin.H{
		"summary":     summary,
		"description": "需要 write 或 admin 权限",
		"requestBody": gin.H{
			"required": true,
			"content": gin.H{
				"application/json": gin.H{
					"schema": gin.H{"$ref": "#/components/schemas/" + schemaRef},
				},
			},
		},
		"responses": gin.H{
			"200": gin.H{"description": "成功"},
		},
	}
}

func buildPathPut(summary, schemaRef string) gin.H {
	return gin.H{
		"summary":     summary,
		"description": "需要 write 或 admin 权限",
		"parameters": []gin.H{
			{"name": "id", "in": "path", "required": true, "schema": gin.H{"type": "integer"}},
		},
		"requestBody": gin.H{
			"required": true,
			"content": gin.H{
				"application/json": gin.H{
					"schema": gin.H{"$ref": "#/components/schemas/" + schemaRef},
				},
			},
		},
		"responses": gin.H{
			"200": gin.H{"description": "成功"},
		},
	}
}

func buildPathDelete(summary string) gin.H {
	return gin.H{
		"summary":     summary,
		"description": "需要 admin 权限",
		"parameters": []gin.H{
			{"name": "id", "in": "path", "required": true, "schema": gin.H{"type": "integer"}},
		},
		"responses": gin.H{
			"200": gin.H{"description": "成功"},
		},
	}
}
