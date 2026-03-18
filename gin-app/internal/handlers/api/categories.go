package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// CategoryRequest 分类请求结构
type CategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
}

// ListCategories 获取分类列表
func ListCategories(c *gin.Context) {
	var categories []models.Category

	// 查询所有分类（不分页）
	if err := models.DB.Order("sort_order ASC, created_at DESC").Find(&categories).Error; err != nil {
		InternalError(c, "获取分类列表失败")
		return
	}

	// 返回标准格式
	SuccessResponse(c, categories)
}

// GetCategoryTree 获取分类树
func GetCategoryTree(c *gin.Context) {
	var categories []models.Category

	// 查询所有分类
	if err := models.DB.Order("sort_order ASC, created_at DESC").Find(&categories).Error; err != nil {
		InternalError(c, "获取分类列表失败")
		return
	}

	// 构建树形结构
	tree := buildCategoryTree(categories, nil)

	SuccessResponse(c, tree)
}

// buildCategoryTree 构建分类树
func buildCategoryTree(categories []models.Category, parentID *uint) []gin.H {
	var result []gin.H
	for _, cat := range categories {
		if (parentID == nil && cat.ParentID == nil) || (parentID != nil && cat.ParentID != nil && *cat.ParentID == *parentID) {
			children := buildCategoryTree(categories, &cat.ID)
			item := gin.H{
				"id":          cat.ID,
				"name":        cat.Name,
				"slug":        cat.Slug,
				"description": cat.Description,
				"parent_id":   cat.ParentID,
				"sort_order":  cat.SortOrder,
				"created_at":  cat.CreatedAt.Format("2006-01-02 15:04:05"),
				"updated_at":  cat.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
			if len(children) > 0 {
				item["children"] = children
			}
			result = append(result, item)
		}
	}
	return result
}

// GetCategory 获取单个分类
func GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的分类ID")
		return
	}

	var category models.Category
	if err := models.DB.Preload("Parent").Preload("Children").First(&category, id).Error; err != nil {
		NotFound(c, "分类不存在")
		return
	}

	SuccessResponse(c, category)
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if req.Name == "" {
		BadRequest(c, "分类名称不能为空")
		return
	}

	// 生成 slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	// 检查 slug 是否已存在
	var count int64
	models.DB.Model(&models.Category{}).Where("slug = ?", req.Slug).Count(&count)
	if count > 0 {
		BadRequest(c, "分类别名已存在")
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
	}

	if err := models.DB.Create(&category).Error; err != nil {
		InternalError(c, "创建分类失败")
		return
	}

	SuccessWithMessage(c, "分类创建成功", category)
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的分类ID")
		return
	}

	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		NotFound(c, "分类不存在")
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 检查 slug 是否已存在（排除自己）
	if req.Slug != "" && req.Slug != category.Slug {
		var count int64
		models.DB.Model(&models.Category{}).Where("slug = ? AND id != ?", req.Slug, id).Count(&count)
		if count > 0 {
			BadRequest(c, "分类别名已存在")
			return
		}
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"slug":        req.Slug,
		"description": req.Description,
		"parent_id":   req.ParentID,
		"sort_order":  req.SortOrder,
	}

	if err := models.DB.Model(&category).Updates(updates).Error; err != nil {
		InternalError(c, "更新分类失败")
		return
	}

	SuccessWithMessage(c, "分类更新成功", category)
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的分类ID")
		return
	}

	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		NotFound(c, "分类不存在")
		return
	}

	// 检查是否有子分类
	var childCount int64
	models.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		BadRequest(c, "该分类下有子分类，无法删除")
		return
	}

	// 检查是否有文章
	var postCount int64
	models.DB.Model(&models.Post{}).Where("category_id = ?", id).Count(&postCount)
	if postCount > 0 {
		BadRequest(c, "该分类下有文章，无法删除")
		return
	}

	if err := models.DB.Delete(&category).Error; err != nil {
		InternalError(c, "删除分类失败")
		return
	}

	SuccessWithMessage(c, "分类删除成功", nil)
}
