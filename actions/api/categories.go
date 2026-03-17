package api

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"

	"catter-cms/models"
)

// CategoriesIndex 获取分类列表
func CategoriesIndex(c buffalo.Context) error {
	categories := models.Categories{}
	query := models.DB.Order("name asc")

	if err := query.All(&categories); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusOK, categories)
}

// CategoriesShow 获取分类详情
func CategoriesShow(c buffalo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	category := &models.Category{}
	if err := models.DB.Find(category, id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// 加载子分类
	if err := models.DB.Where("parent_id = ?", category.ID).All(&category.Children); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusOK, category)
}

// CategoriesCreate 创建分类
func CategoriesCreate(c buffalo.Context) error {
	category := &models.Category{}
	if err := c.Bind(category); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// 验证数据
	verrs, err := category.Validate(models.DB)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	if verrs.HasAny() {
		return c.Error(http.StatusUnprocessableEntity, verrs)
	}

	// 保存分类
	if err := models.DB.Create(category); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusCreated, category)
}

// CategoriesUpdate 更新分类
func CategoriesUpdate(c buffalo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	category := &models.Category{}
	if err := models.DB.Find(category, id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(category); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// 验证数据
	verrs, err := category.Validate(models.DB)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	if verrs.HasAny() {
		return c.Error(http.StatusUnprocessableEntity, verrs)
	}

	// 更新分类
	if err := models.DB.Update(category); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusOK, category)
}

// CategoriesDestroy 删除分类
func CategoriesDestroy(c buffalo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	category := &models.Category{}
	if err := models.DB.Find(category, id); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := models.DB.Destroy(category); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return JSON(c, http.StatusOK, map[string]string{
		"message": "分类删除成功",
	})
}
