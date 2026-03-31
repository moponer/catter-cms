package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
	"catter-cms-gin/internal/services"
)

// PageRequest 页面请求结构
type PageRequest struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Content         string `json:"content"`
	Excerpt         string `json:"excerpt"`
	Status          string `json:"status"`
	Template        string `json:"template"`
	TemplateData    string `json:"template_data"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	SortOrder       int    `json:"sort_order"`
}

// ListPages 获取页面列表
func ListPages(c *gin.Context) {
	var pages []models.Page
	query := models.DB

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Model(&models.Page{}).Count(&total)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&pages).Error; err != nil {
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

// GetPage 获取单个页面
func GetPage(c *gin.Context) {
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

// CreatePage 创建页面
func CreatePage(c *gin.Context) {
	var req PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if req.Title == "" {
		BadRequest(c, "页面标题不能为空")
		return
	}

	// 生成 slug
	if req.Slug == "" {
		req.Slug = services.GenerateSlug(req.Title)
	}

	// 检查 slug 是否已存在
	var count int64
	models.DB.Model(&models.Page{}).Where("slug = ?", req.Slug).Count(&count)
	if count > 0 {
		BadRequest(c, "页面别名已存在")
		return
	}

	page := models.Page{
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Status:          req.Status,
		Template:        req.Template,
		TemplateData:    req.TemplateData,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		SortOrder:       req.SortOrder,
	}

	if page.Status == "" {
		page.Status = "draft"
	}

	if err := models.DB.Create(&page).Error; err != nil {
		InternalError(c, "创建页面失败")
		return
	}

	SuccessWithMessage(c, "页面创建成功", page)
}

// UpdatePage 更新页面
func UpdatePage(c *gin.Context) {
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

	var req PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 检查 slug 是否已存在（排除自己）
	if req.Slug != "" && req.Slug != page.Slug {
		var count int64
		models.DB.Model(&models.Page{}).Where("slug = ? AND id != ?", req.Slug, id).Count(&count)
		if count > 0 {
			BadRequest(c, "页面别名已存在")
			return
		}
	}

	updates := map[string]interface{}{
		"title":             req.Title,
		"slug":              req.Slug,
		"content":          req.Content,
		"excerpt":           req.Excerpt,
		"status":            req.Status,
		"template":          req.Template,
		"template_data":     req.TemplateData,
		"meta_title":        req.MetaTitle,
		"meta_description":  req.MetaDescription,
		"sort_order":       req.SortOrder,
	}

	if err := models.DB.Model(&page).Updates(updates).Error; err != nil {
		InternalError(c, "更新页面失败")
		return
	}

	SuccessWithMessage(c, "页面更新成功", page)
}

// DeletePage 删除页面
func DeletePage(c *gin.Context) {
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

	if err := models.DB.Delete(&page).Error; err != nil {
		InternalError(c, "删除页面失败")
		return
	}

	SuccessWithMessage(c, "页面删除成功", nil)
}
