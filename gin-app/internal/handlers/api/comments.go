package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// CommentRequest 评论请求结构
type CommentRequest struct {
	PostID      uint   `json:"post_id"`
	ParentID    *uint  `json:"parent_id"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	AuthorURL   string `json:"author_url"`
	Content     string `json:"content"`
}

// ListComments 获取评论列表
func ListComments(c *gin.Context) {
	var comments []models.Comment
	query := models.DB.Preload("Post").Preload("User").Preload("Parent")

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 文章筛选
	if postID := c.Query("post_id"); postID != "" {
		query = query.Where("post_id = ?", postID)
	}

	// 获取总数
	var total int64
	query.Model(&models.Comment{}).Count(&total)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

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

// GetComment 获取单个评论
func GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的评论ID")
		return
	}

	var comment models.Comment
	if err := models.DB.Preload("Post").Preload("User").Preload("Parent").First(&comment, id).Error; err != nil {
		NotFound(c, "评论不存在")
		return
	}

	SuccessResponse(c, comment)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的评论ID")
		return
	}

	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		NotFound(c, "评论不存在")
		return
	}

	// 删除该评论及其所有回复
	if err := models.DB.Delete(&comment).Error; err != nil {
		InternalError(c, "删除评论失败")
		return
	}

	// 删除子评论
	models.DB.Where("parent_id = ?", id).Delete(&models.Comment{})

	SuccessWithMessage(c, "评论删除成功", nil)
}

// ApproveComment 批准评论
func ApproveComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的评论ID")
		return
	}

	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		NotFound(c, "评论不存在")
		return
	}

	if err := models.DB.Model(&comment).Update("status", "approved").Error; err != nil {
		InternalError(c, "批准评论失败")
		return
	}

	SuccessWithMessage(c, "评论已批准", comment)
}

// RejectComment 拒绝评论
func RejectComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的评论ID")
		return
	}

	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		NotFound(c, "评论不存在")
		return
	}

	if err := models.DB.Model(&comment).Update("status", "spam").Error; err != nil {
		InternalError(c, "拒绝评论失败")
		return
	}

	SuccessWithMessage(c, "评论已拒绝", comment)
}

// ListPostComments 获取文章评论（公开）
func ListPostComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的文章ID")
		return
	}

	var comments []models.Comment
	query := models.DB.Where("post_id = ? AND status = ?", postID, "approved")

	// 获取总数
	var total int64
	query.Model(&models.Comment{}).Count(&total)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		InternalError(c, "获取评论列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"comments":  comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateComment 创建评论（公开）
func CreateComment(c *gin.Context) {
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 验证必填字段
	if req.Content == "" || req.AuthorName == "" || req.AuthorEmail == "" {
		BadRequest(c, "评论内容、名称和邮箱不能为空")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := models.DB.First(&post, req.PostID).Error; err != nil {
		NotFound(c, "文章不存在")
		return
	}

	// 获取当前用户（如果已登录）
	var userID *uint
	if user, exists := c.Get("current_user"); exists && user != nil {
		currentUser := user.(*models.User)
		userID = &currentUser.ID
	}

	comment := models.Comment{
		PostID:      req.PostID,
		ParentID:    req.ParentID,
		UserID:      userID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		AuthorURL:   req.AuthorURL,
		Content:     req.Content,
		Status:      "pending", // 新评论需要审核
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}

	if err := models.DB.Create(&comment).Error; err != nil {
		InternalError(c, "创建评论失败")
		return
	}

	SuccessWithMessage(c, "评论提交成功，等待审核", comment)
}
