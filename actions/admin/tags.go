package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"strconv"

	"github.com/gobuffalo/buffalo"
)

// TagsListAPI 获取标签列表API
func TagsListAPI(c buffalo.Context) error {
	// 从查询字符串中获取参数
	page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(c.Request().URL.Query().Get("page_size"))
	keyword := c.Request().URL.Query().Get("keyword")

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 50
	}

	tags := &models.Tags{}
	query := models.DB.Paginate(page, pageSize)
	query = query.Order("created_at DESC")

	// 搜索
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.All(tags); err != nil {
		return api.InternalError(c, "获取标签列表失败")
	}

	// 获取总数
	total, _ := models.DB.Count(&models.Tag{})

	return api.PagedSuccess(c, tags, total, page, pageSize)
}

// TagCreateRequest 创建标签请求
type TagCreateRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TagsCreateAPI 创建标签API
func TagsCreateAPI(c buffalo.Context) error {
	req := &TagCreateRequest{}
	if err := c.Bind(req); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	if req.Name == "" {
		return api.BadRequest(c, "标签名称不能为空")
	}

	// 生成slug
	slug := req.Slug
	if slug == "" {
		slug = req.Name
	}

	tag := &models.Tag{
		Name: req.Name,
		Slug: slug,
	}

	if err := models.DB.Create(tag); err != nil {
		return api.InternalError(c, "创建标签失败")
	}

	return api.SuccessWithMessage(c, "标签创建成功", tag)
}

// TagsUpdateAPI 更新标签API
func TagsUpdateAPI(c buffalo.Context) error {
	tagID := c.Param("id")
	tag := &models.Tag{}

	if err := models.DB.Find(tag, tagID); err != nil {
		return api.NotFound(c, "标签不存在")
	}

	req := &TagCreateRequest{}
	if err := c.Bind(req); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Slug != "" {
		tag.Slug = req.Slug
	}

	if err := models.DB.Update(tag); err != nil {
		return api.InternalError(c, "更新标签失败")
	}

	return api.SuccessWithMessage(c, "标签更新成功", tag)
}

// TagsDeleteAPI 删除标签API
func TagsDeleteAPI(c buffalo.Context) error {
	tagID := c.Param("id")
	tag := &models.Tag{}

	if err := models.DB.Find(tag, tagID); err != nil {
		return api.NotFound(c, "标签不存在")
	}

	if err := models.DB.Destroy(tag); err != nil {
		return api.InternalError(c, "删除标签失败")
	}

	return api.SuccessWithMessage(c, "标签删除成功", nil)
}
