package front

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// FrontCommentsListAPI 获取文章评论列表API(前台)
func FrontCommentsListAPI(c buffalo.Context) error {
	postID := c.Param("post_id")
	page, _ := strconv.Atoi(c.Param("page"))
	pageSize, _ := strconv.Atoi(c.Param("page_size"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 50
	}

	comments := &models.Comments{}
	query := models.DB.Where("post_id = ?", postID)
	query = query.Where("status = ?", "approved") // 只显示已批准的评论
	query = query.Where("parent_id IS NULL") // 只获取顶层评论
	query = query.Order("created_at ASC") // 按时间正序

	// 分页
	query = query.Paginate(page, pageSize)

	// 预加载关联数据 - 暂时禁用Eager避免反射错误
	// query = query.Eager("User", "Replies")

	if err := query.All(comments); err != nil {
		return api.InternalError(c, "获取评论列表失败")
	}

	// 获取总数
	total, _ := models.DB.Where("post_id = ? AND status = ? AND parent_id IS NULL", postID, "approved").Count(&models.Comment{})

	return api.PagedSuccess(c, comments, total, page, pageSize)
}

// FrontCommentCreateAPI 创建评论API(前台)
func FrontCommentCreateAPI(c buffalo.Context) error {
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
	comment.Status = "pending"

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

	return api.SuccessWithMessage(c, "评论提交成功,等待审核", comment)
}

// unused - avoid compile error
var _ = render.HTML
