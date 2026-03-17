package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"

	"github.com/gobuffalo/buffalo"
)

// CategoriesListAPI 获取分类列表API
func CategoriesListAPI(c buffalo.Context) error {
	categories := &models.Categories{}
	query := models.DB.Order("created_at DESC")

	// 搜索 - 从查询字符串中获取参数
	keyword := c.Request().URL.Query().Get("keyword")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.All(categories); err != nil {
		return api.InternalError(c, "获取分类列表失败")
	}

	return api.Success(c, categories)
}

// CategoriesTreeAPI 获取分类树形结构API
func CategoriesTreeAPI(c buffalo.Context) error {
	categories := &models.Categories{}
	if err := models.DB.Order("created_at ASC").All(categories); err != nil {
		return api.InternalError(c, "获取分类列表失败")
	}

	// 构建树形结构
	tree := buildCategoryTree(*categories, nil)

	return api.Success(c, tree)
}

// buildCategoryTree 构建分类树形结构
func buildCategoryTree(categories models.Categories, parentID *int) []map[string]interface{} {
	tree := []map[string]interface{}{}

	for _, cat := range categories {
		// 比较ParentID
		isMatch := (parentID == nil && cat.ParentID == nil) ||
			(parentID != nil && cat.ParentID != nil && *cat.ParentID == *parentID)

		if isMatch {
			node := map[string]interface{}{
				"id":          cat.ID,
				"name":        cat.Name,
				"slug":        cat.Slug,
				"description": cat.Description,
				"parent_id":   cat.ParentID,
				"children":    buildCategoryTree(categories, intPtr(cat.ID)),
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// intPtr 创建int指针
func intPtr(i int) *int {
	return &i
}

// CategoryDetailAPI 获取分类详情API
func CategoryDetailAPI(c buffalo.Context) error {
	categoryID := c.Param("id")
	category := &models.Category{}

	if err := models.DB.Find(category, categoryID); err != nil {
		return api.NotFound(c, "分类不存在")
	}

	return api.Success(c, category)
}

// CategoryRequest 分类请求结构
type CategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id"`
}

// CategoryCreateAPI 创建分类API
func CategoryCreateAPI(c buffalo.Context) error {
	req := &CategoryRequest{}
	if err := c.Bind(req); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	if req.Name == "" {
		return api.BadRequest(c, "分类名称不能为空")
	}

	// 生成slug
	slug := req.Slug
	if slug == "" {
		slug = req.Name
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	if err := models.DB.Create(category); err != nil {
		return api.InternalError(c, "创建分类失败")
	}

	return api.SuccessWithMessage(c, "分类创建成功", category)
}

// CategoryUpdateAPI 更新分类API
func CategoryUpdateAPI(c buffalo.Context) error {
	categoryID := c.Param("id")
	category := &models.Category{}

	if err := models.DB.Find(category, categoryID); err != nil {
		return api.NotFound(c, "分类不存在")
	}

	req := &CategoryRequest{}
	if err := c.Bind(req); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	category.ParentID = req.ParentID

	if err := models.DB.Update(category); err != nil {
		return api.InternalError(c, "更新分类失败")
	}

	return api.SuccessWithMessage(c, "分类更新成功", category)
}

// CategoryDeleteAPI 删除分类API
func CategoryDeleteAPI(c buffalo.Context) error {
	categoryID := c.Param("id")
	category := &models.Category{}

	if err := models.DB.Find(category, categoryID); err != nil {
		return api.NotFound(c, "分类不存在")
	}

	// 检查是否有子分类
	children := &models.Categories{}
	if err := models.DB.Where("parent_id = ?", category.ID).All(children); err == nil && len(*children) > 0 {
		return api.BadRequest(c, "该分类下有子分类，无法删除")
	}

	// 检查是否有文章
	posts := &models.Posts{}
	if err := models.DB.Where("category_id = ?", category.ID).All(posts); err == nil && len(*posts) > 0 {
		return api.BadRequest(c, "该分类下有文章，无法删除")
	}

	if err := models.DB.Destroy(category); err != nil {
		return api.InternalError(c, "删除分类失败")
	}

	return api.SuccessWithMessage(c, "分类删除成功", nil)
}

// CategoryReorderAPI 分类排序API
func CategoryReorderAPI(c buffalo.Context) error {
	type ReorderItem struct {
		ID       int  `json:"id"`
		ParentID *int `json:"parent_id"`
	}

	items := []ReorderItem{}
	if err := c.Bind(&items); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	for _, item := range items {
		category := &models.Category{}
		if err := models.DB.Find(category, item.ID); err == nil {
			category.ParentID = item.ParentID
			models.DB.Update(category)
		}
	}

	return api.SuccessWithMessage(c, "排序更新成功", nil)
}
