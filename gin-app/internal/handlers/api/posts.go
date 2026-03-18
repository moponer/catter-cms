package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// PostRequest 文章请求结构
type PostRequest struct {
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	Status      string     `json:"status"`
	CategoryID  *uint      `json:"category_id"`
	TagIDs      []uint     `json:"tag_ids"`
	PublishedAt *string    `json:"published_at"`
}

// ListPosts 获取文章列表
func ListPosts(c *gin.Context) {
	var posts []models.Post
	query := models.DB.Preload("Category").Preload("Author").Preload("Tags")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 获取总数
	var total int64
	query.Model(&models.Post{}).Count(&total)

	// 查询数据
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

// GetPost 获取单个文章
func GetPost(c *gin.Context) {
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

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if req.Title == "" {
		BadRequest(c, "文章标题不能为空")
		return
	}

	// 获取当前用户
	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}
	currentUser := user.(*models.User)

	// 生成 slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}

	post := models.Post{
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		AuthorID:   currentUser.ID,
	}

	if post.Status == "" {
		post.Status = "draft"
	}

	if err := models.DB.Create(&post).Error; err != nil {
		InternalError(c, "创建文章失败")
		return
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		models.DB.Find(&tags, req.TagIDs)
		models.DB.Model(&post).Association("Tags").Append(&tags)
	}

	// 重新加载关联
	models.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, post.ID)

	SuccessWithMessage(c, "文章创建成功", post)
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
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

	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 更新字段
	updates := map[string]interface{}{
		"title":       req.Title,
		"slug":        req.Slug,
		"content":     req.Content,
		"excerpt":     req.Excerpt,
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

	// 更新标签关联
	if req.TagIDs != nil {
		var tags []models.Tag
		models.DB.Find(&tags, req.TagIDs)
		models.DB.Model(&post).Association("Tags").Replace(&tags)
	}

	// 重新加载关联
	models.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, post.ID)

	SuccessWithMessage(c, "文章更新成功", post)
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
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

// generateSlug 生成 URL 友好的 slug
func generateSlug(title string) string {
	slug := ""
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			slug += string(r)
		} else if r == ' ' || r == '-' {
			slug += "-"
		}
	}
	return slug
}
