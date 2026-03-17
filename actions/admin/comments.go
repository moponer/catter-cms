package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
)

// CommentsListAPI 获取评论列表API
func CommentsListAPI(c buffalo.Context) error {
	// 从查询字符串中获取参数
	page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(c.Request().URL.Query().Get("page_size"))
	keyword := c.Request().URL.Query().Get("keyword")
	status := c.Request().URL.Query().Get("status")
	postID := c.Request().URL.Query().Get("post_id")

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	comments := &models.Comments{}
	query := models.DB.Paginate(page, pageSize)
	query = query.Order("created_at DESC")

	// 搜索
	if keyword != "" {
		query = query.Where("content LIKE ? OR author_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 文章筛选
	if postID != "" {
		query = query.Where("post_id = ?", postID)
	}

	// 预加载关联数据 - 暂时禁用Eager避免反射错误
	// query = query.Eager("User", "Post")

	if err := query.All(comments); err != nil {
		return api.InternalError(c, "获取评论列表失败")
	}

	// 获取总数 - 应用相同的过滤条件
	countQuery := models.DB.Q()
	if keyword != "" {
		countQuery = countQuery.Where("content LIKE ? OR author_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	if postID != "" {
		countQuery = countQuery.Where("post_id = ?", postID)
	}
	total, _ := countQuery.Count(&models.Comment{})

	return api.PagedSuccess(c, comments, total, page, pageSize)
}

// CommentDetailAPI 获取评论详情API
func CommentDetailAPI(c buffalo.Context) error {
	commentID := c.Param("id")
	comment := &models.Comment{}

	// 预加载关联数据 - 暂时禁用Eager避免反射错误
	if err := models.DB.Find(comment, commentID); err != nil {
		return api.NotFound(c, "评论不存在")
	}

	return api.Success(c, comment)
}

// CommentCreateAPI 创建评论API
func CommentCreateAPI(c buffalo.Context) error {
	comment := &models.Comment{}
	if err := c.Bind(comment); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 设置IP和User-Agent
	comment.IPAddress = strings.Split(c.Request().RemoteAddr, ":")[0]
	comment.UserAgent = c.Request().Header.Get("User-Agent")

	// 如果已登录,设置用户ID
	if currentUser, ok := c.Value("current_user").(*models.User); ok {
		comment.UserID = &currentUser.ID
		comment.AuthorName = currentUser.Name
		comment.AuthorEmail = currentUser.Email
	}

	// 默认状态为待审核
	if comment.Status == "" {
		comment.Status = "pending"
	}

	// 验证数据
	verrs, err := comment.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 保存评论
	if err := models.DB.Create(comment); err != nil {
		return api.InternalError(c, "创建评论失败")
	}

	return api.SuccessWithMessage(c, "评论创建成功", comment)
}

// CommentUpdateAPI 更新评论API
func CommentUpdateAPI(c buffalo.Context) error {
	commentID := c.Param("id")
	comment := &models.Comment{}

	if err := models.DB.Find(comment, commentID); err != nil {
		return api.NotFound(c, "评论不存在")
	}

	req := map[string]interface{}{}
	if err := c.Bind(&req); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 更新字段
	if status, ok := req["status"].(string); ok {
		comment.Status = status
	}
	if content, ok := req["content"].(string); ok {
		comment.Content = content
	}
	if authorName, ok := req["author_name"].(string); ok {
		comment.AuthorName = authorName
	}
	if authorEmail, ok := req["author_email"].(string); ok {
		comment.AuthorEmail = authorEmail
	}

	// 验证数据
	verrs, err := comment.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 更新评论
	if err := models.DB.Update(comment); err != nil {
		return api.InternalError(c, "更新评论失败")
	}

	return api.SuccessWithMessage(c, "评论更新成功", comment)
}

// CommentDeleteAPI 删除评论API
func CommentDeleteAPI(c buffalo.Context) error {
	commentID := c.Param("id")
	comment := &models.Comment{}

	if err := models.DB.Find(comment, commentID); err != nil {
		return api.NotFound(c, "评论不存在")
	}

	if err := models.DB.Destroy(comment); err != nil {
		return api.InternalError(c, "删除评论失败")
	}

	return api.SuccessWithMessage(c, "评论删除成功", nil)
}

// CommentApproveAPI 批准评论API
func CommentApproveAPI(c buffalo.Context) error {
	commentID := c.Param("id")
	comment := &models.Comment{}

	if err := models.DB.Find(comment, commentID); err != nil {
		return api.NotFound(c, "评论不存在")
	}

	comment.Status = "approved"
	if err := models.DB.Update(comment); err != nil {
		return api.InternalError(c, "批准评论失败")
	}

	return api.SuccessWithMessage(c, "评论已批准", comment)
}

// CommentRejectAPI 拒绝评论API
func CommentRejectAPI(c buffalo.Context) error {
	commentID := c.Param("id")
	comment := &models.Comment{}

	if err := models.DB.Find(comment, commentID); err != nil {
		return api.NotFound(c, "评论不存在")
	}

	comment.Status = "spam"
	if err := models.DB.Update(comment); err != nil {
		return api.InternalError(c, "拒绝评论失败")
	}

	return api.SuccessWithMessage(c, "评论已拒绝", comment)
}

// CommentBatchDeleteAPI 批量删除评论API
func CommentBatchDeleteAPI(c buffalo.Context) error {
	ids := []string{}
	if err := c.Bind(&ids); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	for _, id := range ids {
		comment := &models.Comment{}
		if err := models.DB.Find(comment, id); err != nil {
			continue
		}
		models.DB.Destroy(comment)
	}

	return api.SuccessWithMessage(c, "批量删除成功", nil)
}

// CommentBatchApproveAPI 批量批准评论API
func CommentBatchApproveAPI(c buffalo.Context) error {
	ids := []string{}
	if err := c.Bind(&ids); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	comments := &models.Comments{}
	if err := models.DB.Where("id IN (?)", ids).All(comments); err != nil {
		return api.InternalError(c, "批量批准失败")
	}

	for _, comment := range *comments {
		comment.Status = "approved"
		models.DB.Update(comment)
	}

	return api.SuccessWithMessage(c, "批量批准成功", nil)
}

// CommentsIndexHandler 评论列表页面 (兼容旧API)
func CommentsIndexHandler(c buffalo.Context) error {
	return c.Render(200, render.HTML("<html><body><p>Comments SPA</p></body></html>"))
}

// unused - avoid compile error
var _ = uuid.UUID{}
