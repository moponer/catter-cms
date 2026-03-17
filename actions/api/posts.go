package api

import (
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"

	"catter-cms/models"
)

// PostsIndex 获取文章列表
func PostsIndex(c buffalo.Context) error {
	posts := models.Posts{}
	query := models.DB.Where("status = ?", "published").Order("published_at desc")

	// 分页参数
	limit := c.Param("limit")
	if limit == "" {
		limit = "10"
	}
	page := c.Param("page")
	if page == "" {
		page = "1"
	}

	// 执行分页查询
	paginatedQuery := query.PaginateFromParams(c.Params())
	if err := paginatedQuery.All(&posts); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// 返回JSON
	return JSON(c, http.StatusOK, map[string]interface{}{
		"posts": posts,
		"meta": map[string]interface{}{
			"total_pages": 1,
			"total_count": len(posts),
			"page":        1,
			"per_page":    10,
		},
	})
}

// PostsShow 获取文章详情
func PostsShow(c buffalo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	post := &models.Post{}
	if err := models.DB.Find(post, id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return JSON(c, http.StatusOK, post)
}

// PostsCreate 创建文章
func PostsCreate(c buffalo.Context) error {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// 设置作者
	user := c.Value("current_user").(*models.User)
	post.UserID = user.ID

	// 验证数据
	verrs, err := post.Validate(models.DB)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	if verrs.HasAny() {
		return c.Error(http.StatusUnprocessableEntity, verrs)
	}

	// 保存文章
	if err := models.DB.Create(post); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusCreated, post)
}

// PostsUpdate 更新文章
func PostsUpdate(c buffalo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	post := &models.Post{}
	if err := models.DB.Find(post, id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(post); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// 验证数据
	verrs, err := post.Validate(models.DB)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	if verrs.HasAny() {
		return c.Error(http.StatusUnprocessableEntity, verrs)
	}

	// 更新文章
	if err := models.DB.Update(post); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusOK, post)
}

// PostsDestroy 删除文章
func PostsDestroy(c buffalo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	post := &models.Post{}
	if err := models.DB.Find(post, id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := models.DB.Destroy(post); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusOK, map[string]string{
		"message": "文章删除成功",
	})
}