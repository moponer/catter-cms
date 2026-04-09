package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
	"catter-cms-gin/internal/services"
)

// TagRequest 标签请求结构
type TagRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ListTags 获取标签列表
func ListTags(c *gin.Context) {
	var tags []models.Tag
	query := models.DB

	// 搜索
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Model(&models.Tag{}).Count(&total)

	// 查询数据
	if err := query.Order("name ASC").Find(&tags).Error; err != nil {
		InternalError(c, "获取标签列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":  tags,
		"total": total,
	})
}

// GetTag 获取单个标签
func GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的标签ID")
		return
	}

	var tag models.Tag
	if err := models.DB.First(&tag, id).Error; err != nil {
		NotFound(c, "标签不存在")
		return
	}

	SuccessResponse(c, tag)
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if req.Name == "" {
		BadRequest(c, "标签名称不能为空")
		return
	}

	// 生成 slug
	if req.Slug == "" {
		req.Slug = services.GenerateSlug(req.Name)
	}

	// 检查名称是否已存在
	var count int64
	models.DB.Model(&models.Tag{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		BadRequest(c, "标签名称已存在")
		return
	}

	// 检查 slug 是否已存在
	models.DB.Model(&models.Tag{}).Where("slug = ?", req.Slug).Count(&count)
	if count > 0 {
		BadRequest(c, "标签别名已存在")
		return
	}

	tag := models.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := models.DB.Create(&tag).Error; err != nil {
		InternalError(c, "创建标签失败")
		return
	}

	SuccessWithMessage(c, "标签创建成功", tag)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的标签ID")
		return
	}

	var tag models.Tag
	if err := models.DB.First(&tag, id).Error; err != nil {
		NotFound(c, "标签不存在")
		return
	}

	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 检查名称是否已存在（排除自己）
	if req.Name != "" && req.Name != tag.Name {
		var count int64
		models.DB.Model(&models.Tag{}).Where("name = ? AND id != ?", req.Name, id).Count(&count)
		if count > 0 {
			BadRequest(c, "标签名称已存在")
			return
		}
	}

	// 检查 slug 是否已存在（排除自己）
	if req.Slug != "" && req.Slug != tag.Slug {
		var count int64
		models.DB.Model(&models.Tag{}).Where("slug = ? AND id != ?", req.Slug, id).Count(&count)
		if count > 0 {
			BadRequest(c, "标签别名已存在")
			return
		}
	}

	updates := map[string]interface{}{
		"name": req.Name,
		"slug": req.Slug,
	}

	if err := models.DB.Model(&tag).Updates(updates).Error; err != nil {
		InternalError(c, "更新标签失败")
		return
	}

	SuccessWithMessage(c, "标签更新成功", tag)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的标签ID")
		return
	}

	var tag models.Tag
	if err := models.DB.First(&tag, id).Error; err != nil {
		NotFound(c, "标签不存在")
		return
	}

	if err := models.DB.Delete(&tag).Error; err != nil {
		InternalError(c, "删除标签失败")
		return
	}

	SuccessWithMessage(c, "标签删除成功", nil)
}

// BatchDeleteTags 批量删除标签
func BatchDeleteTags(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	if len(req.IDs) == 0 {
		BadRequest(c, "请选择要删除的标签")
		return
	}

	if err := models.DB.Where("id IN ?", req.IDs).Delete(&models.Tag{}).Error; err != nil {
		InternalError(c, "批量删除标签失败")
		return
	}

	// 清理标签关联
	models.DB.Where("tag_id IN ?", req.IDs).Delete(&models.PostTag{})

	SuccessWithMessage(c, "批量删除成功", gin.H{"deleted_count": len(req.IDs)})
}
