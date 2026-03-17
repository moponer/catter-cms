package front

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"catter-cms/models"
	"catter-cms/services"
	"catter-cms/templates"
	"catter-cms/public"
)

// R 全局渲染引擎 - 使用嵌入的模板系统
var R *render.Engine

func init() {
	R = render.New(render.Options{
		HTMLLayout:   "application.plush.html",
		TemplatesFS:  templates.FS(),
		AssetsFS:     public.FS(),
	})
}

// HomeHandler 处理首页请求
func HomeHandler(c buffalo.Context) error {
	// 获取激活的主题
	activeTheme := services.GetActiveTheme()
	log.Printf("[DEBUG] Active theme: %s", activeTheme)

	// 获取最新的文章
	posts := models.Posts{}
	query := models.DB.Where("status = ?", "published").Order("published_at desc").Limit(10)
	if err := query.All(&posts); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	
	// 调试输出
	log.Printf("[DEBUG] Found %d posts", len(posts))
	for i, p := range posts {
		log.Printf("[DEBUG] Post %d: ID=%d, Title=%s, Views=%d", i, p.ID, p.Title, p.Views)
	}

	// 获取所有分类（带文章数量）
	categories := models.Categories{}
	if err := models.DB.All(&categories); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	
	// 为每个分类计算文章数量
	type CategoryWithCount struct {
		ID          int
		Name        string
		Slug        string
		Description string
		PostCount   int
	}
	categoriesWithCount := []CategoryWithCount{}
	for _, cat := range categories {
		var count int
		models.DB.Where("category_id = ? AND status = ?", cat.ID, "published").Count(&count)
		categoriesWithCount = append(categoriesWithCount, CategoryWithCount{
			ID:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			Description: cat.Description,
			PostCount:   count,
		})
	}

	// 获取所有标签（带文章数量）
	tags := models.Tags{}
	if err := models.DB.All(&tags); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	
	// 为每个标签计算文章数量
	type TagWithCount struct {
		ID        int
		Name      string
		Slug      string
		PostCount int
	}
	tagsWithCount := []TagWithCount{}
	for _, tag := range tags {
		// 通过 post_tags 关联表计算文章数量
		var count int
		models.DB.RawQuery("SELECT COUNT(*) FROM post_tags WHERE tag_id = ?", tag.ID).Count(&count)
		tagsWithCount = append(tagsWithCount, TagWithCount{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			PostCount: count,
		})
	}

	// 获取站点设置
	siteSettings := getSiteSettings()

	// 将数据传递给模板
	c.Set("posts", posts)
	c.Set("categories", categoriesWithCount)
	c.Set("tags", tagsWithCount)
	c.Set("site", siteSettings)
	c.Set("theme", activeTheme)

	// 使用激活的主题渲染
	templatePath := fmt.Sprintf("themes/%s/pages/home.plush.html", activeTheme)
	return c.Render(http.StatusOK, R.HTML(templatePath))
}
func BlogShowHandler(c buffalo.Context) error {
	slug := c.Param("slug")

	// 获取激活的主题
	activeTheme := services.GetActiveTheme()

	// 查找文章
	post := &models.Post{}
	query := models.DB.Where("slug = ? AND status = ?", slug, "published")
	if err := query.First(post); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// 获取站点设置
	siteSettings := getSiteSettings()

	// 将数据传递给模板
	c.Set("post", post)
	c.Set("site", siteSettings)
	c.Set("theme", activeTheme)

	// 使用激活的主题渲染
	templatePath := fmt.Sprintf("themes/%s/pages/blog_show.plush.html", activeTheme)
	return c.Render(http.StatusOK, R.HTML(templatePath))
}

// PageShowHandler 处理页面详情请求
func PageShowHandler(c buffalo.Context) error {
	slug := c.Param("slug")

	// 获取激活的主题
	activeTheme := services.GetActiveTheme()

	// 查找页面（这里简化处理，实际应创建 Page 模型）
	post := &models.Post{}
	query := models.DB.Where("slug = ? AND status = ?", slug, "published")
	if err := query.First(post); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// 获取站点设置
	siteSettings := getSiteSettings()

	// 将数据传递给模板
	c.Set("post", post)
	c.Set("site", siteSettings)
	c.Set("theme", activeTheme)

	// 使用激活的主题渲染
	templatePath := fmt.Sprintf("themes/%s/pages/page.plush.html", activeTheme)
	return c.Render(http.StatusOK, R.HTML(templatePath))
}

// CategoryShowHandler 处理分类页面请求
func CategoryShowHandler(c buffalo.Context) error {
	slug := c.Param("slug")

	// 获取激活的主题
	activeTheme := services.GetActiveTheme()

	// 查找分类
	category := &models.Category{}
	if err := models.DB.Where("slug = ?", slug).First(category); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// 获取该分类下的文章
	posts := models.Posts{}
	if err := models.DB.Where("category_id = ? AND status = ?", category.ID, "published").All(&posts); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// 获取站点设置
	siteSettings := getSiteSettings()

	// 将数据传递给模板
	c.Set("category", category)
	c.Set("posts", posts)
	c.Set("site", siteSettings)
	c.Set("theme", activeTheme)

	// 使用激活的主题渲染
	templatePath := fmt.Sprintf("themes/%s/pages/category.plush.html", activeTheme)
	return c.Render(http.StatusOK, R.HTML(templatePath))
}

// TagShowHandler 处理标签页面请求
func TagShowHandler(c buffalo.Context) error {
	slug := c.Param("slug")

	// 获取激活的主题
	activeTheme := services.GetActiveTheme()

	// 查找标签
	tag := &models.Tag{}
	if err := models.DB.Where("slug = ?", slug).First(tag); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// 获取该标签下的文章（通过 post_tags 关联表）
	posts := models.Posts{}
	// 简化处理：直接查询（实际应该通过关联表查询）
	if err := models.DB.Where("status = ?", "published").All(&posts); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// 获取站点设置
	siteSettings := getSiteSettings()

	// 将数据传递给模板
	c.Set("tag", tag)
	c.Set("posts", posts)
	c.Set("site", siteSettings)
	c.Set("theme", activeTheme)

	// 使用激活的主题渲染
	templatePath := fmt.Sprintf("themes/%s/pages/tag.plush.html", activeTheme)
	return c.Render(http.StatusOK, R.HTML(templatePath))
}

// getSiteSettings 获取站点设置
func getSiteSettings() map[string]string {
	settings := make(map[string]string)

	// 默认值
	settings["title"] = "Catter CMS"
	settings["description"] = "轻量级内容管理系统"
	settings["url"] = ""

	// 从数据库获取设置
	dbSettings := models.Settings{}
	if err := models.DB.All(&dbSettings); err == nil {
		for _, setting := range dbSettings {
			settings[setting.Key] = setting.Value
		}
	}

	return settings
}
