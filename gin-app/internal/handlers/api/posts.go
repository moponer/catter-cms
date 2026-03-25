package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// PostRequest 文章请求结构
type PostRequest struct {
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Content         string     `json:"content"`
	Excerpt         string     `json:"excerpt"`
	Cover           string     `json:"cover"`
	Status          string     `json:"status"`
	MetaTitle       string     `json:"meta_title"`
	MetaDescription string     `json:"meta_description"`
	CategoryID      *uint      `json:"category_id"`
	TagIDs          []uint     `json:"tag_ids"`
	TagNames        []string   `json:"tag_names"` // 支持标签名称，自动创建
	PublishedAt     *string    `json:"published_at"`
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

	// 关键词搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 日期范围筛选
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate+" 00:00:00")
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
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

	// 检查 slug 唯一性
	var existingPost models.Post
	if err := models.DB.Where("slug = ?", req.Slug).First(&existingPost).Error; err == nil {
		BadRequest(c, "文章别名已存在，请使用其他名称")
		return
	}

	post := models.Post{
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Cover:           req.Cover,
		Status:          req.Status,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		CategoryID:      req.CategoryID,
		AuthorID:        currentUser.ID,
	}

	if post.Status == "" {
		post.Status = "draft"
	}

	if err := models.DB.Create(&post).Error; err != nil {
		InternalError(c, "创建文章失败")
		return
	}

	// 关联标签（支持 tag_names 和 tag_ids）
	var postTags []models.Tag
	if len(req.TagNames) > 0 {
		// 使用标签名称，自动创建不存在的标签
		for _, name := range req.TagNames {
			var tag models.Tag
			result := models.DB.Where("name = ?", name).First(&tag)
			if result.Error != nil {
				// 标签不存在，创建新标签
				tag = models.Tag{
					Name: name,
					Slug: generateSlug(name),
				}
				models.DB.Create(&tag)
			}
			postTags = append(postTags, tag)
		}
	} else if len(req.TagIDs) > 0 {
		// 使用标签 ID
		models.DB.Find(&postTags, req.TagIDs)
	}
	if len(postTags) > 0 {
		models.DB.Model(&post).Association("Tags").Append(&postTags)
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

	// 检查 slug 唯一性（如果修改了 slug）
	if req.Slug != "" && req.Slug != post.Slug {
		var existingPost models.Post
		if err := models.DB.Where("slug = ? AND id != ?", req.Slug, id).First(&existingPost).Error; err == nil {
			BadRequest(c, "文章别名已存在，请使用其他名称")
			return
		}
	}

	// 更新字段
	updates := map[string]interface{}{
		"title":            req.Title,
		"slug":             req.Slug,
		"content":          req.Content,
		"excerpt":          req.Excerpt,
		"cover":            req.Cover,
		"status":           req.Status,
		"meta_title":       req.MetaTitle,
		"meta_description": req.MetaDescription,
		"category_id":      req.CategoryID,
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

	// 更新标签关联（支持 tag_names 和 tag_ids）
	if req.TagNames != nil || req.TagIDs != nil {
		var postTags []models.Tag
		if len(req.TagNames) > 0 {
			// 使用标签名称，自动创建不存在的标签
			for _, name := range req.TagNames {
				var tag models.Tag
				result := models.DB.Where("name = ?", name).First(&tag)
				if result.Error != nil {
					// 标签不存在，创建新标签
					tag = models.Tag{
						Name: name,
						Slug: generateSlug(name),
					}
					models.DB.Create(&tag)
				}
				postTags = append(postTags, tag)
			}
		} else if len(req.TagIDs) > 0 {
			models.DB.Find(&postTags, req.TagIDs)
		}
		models.DB.Model(&post).Association("Tags").Replace(&postTags)
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

// BatchDeletePosts 批量删除文章
func BatchDeletePosts(c *gin.Context) {
	var req struct {
		IDs string `json:"ids"` // 逗号分隔的ID列表
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if req.IDs == "" {
		BadRequest(c, "请选择要删除的文章")
		return
	}

	// 解析 ID 列表
	var ids []uint
	idStrs := strings.Split(req.IDs, ",")
	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err == nil {
			ids = append(ids, uint(id))
		}
	}

	if len(ids) == 0 {
		BadRequest(c, "无效的文章ID")
		return
	}

	// 批量删除
	if err := models.DB.Where("id IN ?", ids).Delete(&models.Post{}).Error; err != nil {
		InternalError(c, "批量删除失败")
		return
	}

	SuccessWithMessage(c, "批量删除成功", gin.H{"deleted_count": len(ids)})
}

// BatchUpdateStatus 批量更新文章状态
func BatchUpdateStatus(c *gin.Context) {
	var req struct {
		IDs    string `json:"ids"`    // 逗号分隔的ID列表
		Status string `json:"status"` // 新状态
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if req.IDs == "" {
		BadRequest(c, "请选择要更新的文章")
		return
	}

	if req.Status == "" {
		BadRequest(c, "请指定新状态")
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{"draft": true, "published": true, "archived": true}
	if !validStatuses[req.Status] {
		BadRequest(c, "无效的状态值")
		return
	}

	// 解析 ID 列表
	var ids []uint
	idStrs := strings.Split(req.IDs, ",")
	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err == nil {
			ids = append(ids, uint(id))
		}
	}

	if len(ids) == 0 {
		BadRequest(c, "无效的文章ID")
		return
	}

	// 更新状态
	updates := map[string]interface{}{"status": req.Status}
	if req.Status == "published" {
		now := time.Now()
		updates["published_at"] = &now
	}

	if err := models.DB.Model(&models.Post{}).Where("id IN ?", ids).Updates(updates).Error; err != nil {
		InternalError(c, "批量更新失败")
		return
	}

	SuccessWithMessage(c, "批量更新成功", gin.H{"updated_count": len(ids)})
}

// generateSlug 生成 URL 友好的 slug
func generateSlug(title string) string {
	slug := ""
	hasAlphaNum := false

	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			slug += string(r)
			hasAlphaNum = true
		} else if r == ' ' || r == '-' {
			slug += "-"
		}
	}

	// 如果标题全是中文或其他非拉丁字符，使用时间戳
	if !hasAlphaNum || slug == "" {
		slug = "post-" + time.Now().Format("20060102150405")
	}

	return slug
}
