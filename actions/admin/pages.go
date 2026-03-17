package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"catter-cms/services"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop/v6"
)

// PagesListAPI 获取页面列表API
func PagesListAPI(c buffalo.Context) error {
	// 从查询字符串中获取参数
	page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(c.Request().URL.Query().Get("page_size"))
	keyword := c.Request().URL.Query().Get("keyword")
	status := c.Request().URL.Query().Get("status")

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	pages := &models.Pages{}
	query := models.DB.Paginate(page, pageSize)
	query = query.Order("order ASC, created_at DESC")

	// 搜索
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 加载关联数据
	query = query.Eager("User")

	if err := query.All(pages); err != nil {
		return api.InternalError(c, "获取页面列表失败")
	}

	// 获取总数
	countQuery := models.DB.Q()
	if keyword != "" {
		countQuery = countQuery.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	total, _ := countQuery.Count(&models.Page{})

	return api.PagedSuccess(c, pages, total, page, pageSize)
}

// PageDetailAPI 获取页面详情API
func PageDetailAPI(c buffalo.Context) error {
	pageID := c.Param("id")
	page := &models.Page{}

	if err := models.DB.Eager("User").Find(page, pageID); err != nil {
		return api.NotFound(c, "页面不存在")
	}

	return api.Success(c, page)
}

// PageCreateAPI 创建页面API
func PageCreateAPI(c buffalo.Context) error {
	page := &models.Page{}
	if err := c.Bind(page); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 设置作者
	if currentUser, ok := c.Value("current_user").(*models.User); ok {
		page.UserID = currentUser.ID
	}

	// 验证数据
	verrs, err := page.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 保存页面
	if err := models.DB.Create(page); err != nil {
		return api.InternalError(c, "创建页面失败")
	}

	// 如果页面已发布,添加到搜索索引
	if page.Status == "published" {
		searchService := services.GetSearchService()
		content := page.Content
		if page.Excerpt != "" {
			content = page.Excerpt + "\n" + content
		}
		if err := searchService.IndexPage(
			strconv.Itoa(page.ID),
			page.Title,
			content,
		); err != nil {
			// 搜索索引失败不影响页面创建
		}
	}

	return api.SuccessWithMessage(c, "页面创建成功", page)
}

// PageUpdateAPI 更新页面API
func PageUpdateAPI(c buffalo.Context) error {
	pageID := c.Param("id")
	page := &models.Page{}

	if err := models.DB.Find(page, pageID); err != nil {
		return api.NotFound(c, "页面不存在")
	}

	if err := c.Bind(page); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 验证数据
	verrs, err := page.Validate(models.DB)
	if err != nil {
		return api.InternalError(c, "验证失败")
	}
	if verrs.HasAny() {
		return api.BadRequest(c, verrs.Error())
	}

	// 更新页面
	if err := models.DB.Update(page); err != nil {
		return api.InternalError(c, "更新页面失败")
	}

	// 更新搜索索引
	searchService := services.GetSearchService()
	if page.Status == "published" {
		content := page.Content
		if page.Excerpt != "" {
			content = page.Excerpt + "\n" + content
		}
		if err := searchService.IndexPage(
			strconv.Itoa(page.ID),
			page.Title,
			content,
		); err != nil {
			// 搜索索引失败不影响页面更新
		}
	} else {
		// 页面未发布,从索引中删除
		if err := searchService.DeleteIndex(strconv.Itoa(page.ID)); err != nil {
			// 搜索索引失败不影响页面更新
		}
	}

	return api.SuccessWithMessage(c, "页面更新成功", page)
}

// PageDeleteAPI 删除页面API
func PageDeleteAPI(c buffalo.Context) error {
	pageID := c.Param("id")
	page := &models.Page{}

	if err := models.DB.Find(page, pageID); err != nil {
		return api.NotFound(c, "页面不存在")
	}

	// 检查是否有子页面
	var childCount int
	models.DB.Where("parent_id = ?", page.ID).Count(&childCount)
	if childCount > 0 {
		return api.BadRequest(c, "该页面有子页面，无法删除")
	}

	if err := models.DB.Destroy(page); err != nil {
		return api.InternalError(c, "删除页面失败")
	}

	// 从搜索索引中删除
	searchService := services.GetSearchService()
	if err := searchService.DeleteIndex(strconv.Itoa(page.ID)); err != nil {
		// 搜索索引失败不影响页面删除
	}

	return api.SuccessWithMessage(c, "页面删除成功", nil)
}

// PageReorderAPI 页面排序API
func PageReorderAPI(c buffalo.Context) error {
	var reorderData struct {
		Pages []struct {
			ID    int `json:"id"`
			Order int `json:"order"`
		} `json:"pages"`
	}

	if err := c.Bind(&reorderData); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	err := models.DB.Transaction(func(tx *pop.Connection) error {
		for _, page := range reorderData.Pages {
			if err := tx.RawQuery("UPDATE pages SET \"order\" = ? WHERE id = ?", page.Order, page.ID).Exec(); err != nil {
				return err
			}
		}
		return nil
	})
	
	if err != nil {
		return api.InternalError(c, "更新排序失败")
	}

	return api.SuccessWithMessage(c, "排序更新成功", nil)
}

// PagesIndex 显示页面列表页面 (已弃用，SPA模式)
func PagesIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.HTML("<html><body><p>Pages SPA</p></body></html>"))
}

// PagesNew 显示创建页面页面 (已弃用，SPA模式)
func PagesNew(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.HTML("<html><body><p>Pages SPA</p></body></html>"))
}

// PagesCreate 处理创建页面 (已弃用，使用API)
func PagesCreate(c buffalo.Context) error {
	return PageCreateAPI(c)
}

// PagesEdit 显示编辑页面页面 (已弃用，SPA模式)
func PagesEdit(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.HTML("<html><body><p>Pages SPA</p></body></html>"))
}

// PagesUpdate 处理更新页面 (已弃用，使用API)
func PagesUpdate(c buffalo.Context) error {
	return PageUpdateAPI(c)
}

// PagesDestroy 处理删除页面 (已弃用，使用API)
func PagesDestroy(c buffalo.Context) error {
	return PageDeleteAPI(c)
}
