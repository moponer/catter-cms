package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"catter-cms/services"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// PostsListAPI 获取文章列表API (管理后台)
func PostsListAPI(c buffalo.Context) error {
	// 从查询字符串中获取参数
	page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(c.Request().URL.Query().Get("page_size"))
	keyword := c.Request().URL.Query().Get("keyword")
	status := c.Request().URL.Query().Get("status")
	categoryID := c.Request().URL.Query().Get("category_id")

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	posts := &models.Posts{}
	query := models.DB.Paginate(page, pageSize)
	query = query.Order("created_at DESC")

	// 搜索
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 分类筛选
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.All(posts); err != nil {
		return api.InternalError(c, "获取文章列表失败")
	}

	// 获取总数 - 应用相同的过滤条件
	countQuery := models.DB.Q()
	if keyword != "" {
		countQuery = countQuery.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	if categoryID != "" {
		countQuery = countQuery.Where("category_id = ?", categoryID)
	}
	total, _ := countQuery.Count(&models.Post{})

	return api.PagedSuccess(c, posts, total, page, pageSize)
}

// PostDetailAPI 获取文章详情API
func PostDetailAPI(c buffalo.Context) error {
	postID := c.Param("id")
	post := &models.Post{}

	if err := models.DB.Find(post, postID); err != nil {
		return api.NotFound(c, "文章不存在")
	}

	return api.Success(c, post)
}

// PostCreateAPI 创建文章API
func PostCreateAPI(c buffalo.Context) error {
	post := &models.Post{}
	if err := c.Bind(post); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 设置作者
	if currentUser, ok := c.Value("current_user").(*models.User); ok {
		post.UserID = currentUser.ID
	}

	// 验证数据
	verrs, err := post.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 保存文章
	if err := models.DB.Create(post); err != nil {
		return api.InternalError(c, "创建文章失败")
	}

	// 如果文章已发布,添加到搜索索引
	if post.Status == "published" {
		searchService := services.GetSearchService()
		content := post.Content
		if post.Excerpt != "" {
			content = post.Excerpt + "\n" + content
		}
		if err := searchService.IndexPost(
			strconv.Itoa(post.ID),
			post.Title,
			content,
			"",
		); err != nil {
			// 搜索索引失败不影响文章创建
			// 可以记录日志
		}
	}

	return api.SuccessWithMessage(c, "文章创建成功", post)
}

// PostUpdateAPI 更新文章API
func PostUpdateAPI(c buffalo.Context) error {
	postID := c.Param("id")
	post := &models.Post{}

	if err := models.DB.Find(post, postID); err != nil {
		return api.NotFound(c, "文章不存在")
	}

	if err := c.Bind(post); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 验证数据
	verrs, err := post.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 更新文章
	if err := models.DB.Update(post); err != nil {
		return api.InternalError(c, "更新文章失败")
	}

	// 更新搜索索引
	searchService := services.GetSearchService()
	if post.Status == "published" {
		content := post.Content
		if post.Excerpt != "" {
			content = post.Excerpt + "\n" + content
		}
		if err := searchService.IndexPost(
			strconv.Itoa(post.ID),
			post.Title,
			content,
			"",
		); err != nil {
			// 搜索索引失败不影响文章更新
		}
	} else {
		// 文章未发布,从索引中删除
		if err := searchService.DeleteIndex(strconv.Itoa(post.ID)); err != nil {
			// 搜索索引失败不影响文章更新
		}
	}

	return api.SuccessWithMessage(c, "文章更新成功", post)
}

// PostDeleteAPI 删除文章API
func PostDeleteAPI(c buffalo.Context) error {
	postID := c.Param("id")
	post := &models.Post{}

	if err := models.DB.Find(post, postID); err != nil {
		return api.NotFound(c, "文章不存在")
	}

	if err := models.DB.Destroy(post); err != nil {
		return api.InternalError(c, "删除文章失败")
	}

	// 从搜索索引中删除
	searchService := services.GetSearchService()
	if err := searchService.DeleteIndex(strconv.Itoa(post.ID)); err != nil {
		// 搜索索引失败不影响文章删除
	}

	return api.SuccessWithMessage(c, "文章删除成功", nil)
}

// PostsIndex 显示文章列表页面 (已弃用，SPA模式)
func PostsIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.HTML("<html><body><p>Posts SPA</p></body></html>"))
}

// PostsNew 显示创建文章页面 (已弃用，SPA模式)
func PostsNew(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.HTML("<html><body><p>Posts SPA</p></body></html>"))
}

// PostsCreate 处理创建文章 (已弃用，使用API)
func PostsCreate(c buffalo.Context) error {
	return PostCreateAPI(c)
}

// PostsEdit 显示编辑文章页面 (已弃用，SPA模式)
func PostsEdit(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.HTML("<html><body><p>Posts SPA</p></body></html>"))
}

// PostsUpdate 处理更新文章 (已弃用，使用API)
func PostsUpdate(c buffalo.Context) error {
	return PostUpdateAPI(c)
}

// PostsDestroy 处理删除文章 (已弃用，使用API)
func PostsDestroy(c buffalo.Context) error {
	return PostDeleteAPI(c)
}
