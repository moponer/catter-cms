package tasks

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gobuffalo/buffalo"

	"catter-cms/models"
	"catter-cms/services"
)

// GenerateStaticSite 全量生成静态站点
func GenerateStaticSite() error {
	log.Println("开始全量生成静态站点...")

	// 初始化应用上下文
	app := buffalo.New(buffalo.Options{Env: "production"})

	// 1. 获取激活的主题（从数据库读取）
	activeTheme := services.GetActiveTheme()
	log.Printf("使用主题: %s", activeTheme)

	// 2. 生成首页
	if err := generateHomePage(app, activeTheme); err != nil {
		log.Printf("生成首页失败: %v", err)
		return err
	}

	// 3. 生成所有已发布文章（从数据库读取真实数据）
	posts := models.Posts{}
	if err := models.DB.Where("status = ?", "published").Order("published_at desc").All(&posts); err != nil {
		log.Printf("获取文章列表失败: %v", err)
		return err
	}
	log.Printf("发现 %d 篇已发布文章", len(posts))

	for _, post := range posts {
		if err := generatePostPage(app, activeTheme, post); err != nil {
			log.Printf("生成文章页面失败 %s: %v", post.Slug, err)
			return err
		}
	}

	// 4. 生成分类页面
	if err := generateCategoryPages(app, activeTheme); err != nil {
		log.Printf("生成分类页面失败: %v", err)
		return err
	}

	// 5. 生成标签页面
	if err := generateTagPages(app, activeTheme); err != nil {
		log.Printf("生成标签页面失败: %v", err)
		return err
	}

	log.Println("静态站点生成完成")
	return nil
}

// getPublishedPosts 从数据库获取已发布的文章
func getPublishedPosts() models.Posts {
	posts := models.Posts{}
	if err := models.DB.Where("status = ?", "published").Order("published_at desc").All(&posts); err != nil {
		log.Printf("获取文章失败: %v", err)
		return models.Posts{}
	}
	return posts
}

// getAllCategories 从数据库获取所有分类
func getAllCategories() models.Categories {
	categories := models.Categories{}
	if err := models.DB.All(&categories); err != nil {
		log.Printf("获取分类失败: %v", err)
		return models.Categories{}
	}
	return categories
}

// getAllTags 从数据库获取所有标签
func getAllTags() models.Tags {
	tags := models.Tags{}
	if err := models.DB.All(&tags); err != nil {
		log.Printf("获取标签失败: %v", err)
		return models.Tags{}
	}
	return tags
}

// RevalidatePost 增量更新单篇文章
func RevalidatePost(slug string) error {
	log.Printf("开始增量更新文章: %s", slug)

	// 初始化应用上下文
	app := buffalo.New(buffalo.Options{Env: "production"})

	// 获取激活的主题（从数据库读取）
	activeTheme := services.GetActiveTheme()

	// 查找文章
	post := &models.Post{}
	query := models.DB.Where("slug = ?", slug)
	if err := query.First(post); err != nil {
		log.Printf("文章不存在: %s", slug)
		return err
	}

	// 生成文章页面
	if err := generatePostPage(app, activeTheme, *post); err != nil {
		log.Printf("增量更新文章失败: %v", err)
		return err
	}

	// 如果文章状态为已发布，还需要更新首页
	if post.Status == "published" {
		if err := generateHomePage(app, activeTheme); err != nil {
			log.Printf("更新首页失败: %v", err)
			return err
		}
	}

	log.Printf("文章增量更新完成: %s", slug)
	return nil
}

// generateHomePage 生成首页静态页面
func generateHomePage(app *buffalo.App, theme string) error {
	// 获取最新文章（从数据库读取真实数据）
	posts := getPublishedPosts()

	// 获取分类和标签
	categories := getAllCategories()
	tags := getAllTags()

	// 准备模板数据
	data := map[string]interface{}{
		"posts":      posts,
		"categories": categories,
		"tags":       tags,
		"theme":      theme,
		"site":       getSiteSettings(),
		"now":        func() time.Time { return time.Now() },
	}

	// 渲染模板
	tmplPath := filepath.Join("templates", "themes", theme, "pages", "home.plush.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("解析模板失败 %s: %v", tmplPath, err)
	}

	// 创建输出目录
	outputDir := filepath.Join("public", "static")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// 写入静态文件
	outputPath := filepath.Join(outputDir, "index.html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	log.Println("首页生成完成")
	return nil
}

// generatePostPage 生成单篇文章静态页面
func generatePostPage(app *buffalo.App, theme string, post models.Post) error {
	// 准备模板数据
	data := map[string]interface{}{
		"post":  post,
		"theme": theme,
		"site":  getSiteSettings(),
		"now":   func() time.Time { return time.Now() },
	}

	// 渲染模板
	tmplPath := filepath.Join("templates", "themes", theme, "pages", "blog_show.plush.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("解析模板失败 %s: %v", tmplPath, err)
	}

	// 创建输出目录
	outputDir := filepath.Join("public", "static", "blog")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// 写入静态文件
	outputPath := filepath.Join(outputDir, post.Slug+".html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	log.Printf("文章页面生成完成: %s", post.Slug)
	return nil
}

// generateCategoryPages 生成分类页面
func generateCategoryPages(app *buffalo.App, theme string) error {
	// 获取所有分类（从数据库读取）
	categories := getAllCategories()

	for _, category := range categories {
		// 获取该分类下的文章
		posts := models.Posts{}
		models.DB.Where("category_id = ? AND status = ?", category.ID, "published").All(&posts)

		// 准备模板数据
		data := map[string]interface{}{
			"category": category,
			"posts":    posts,
			"theme":    theme,
			"site":     getSiteSettings(),
			"now":      func() time.Time { return time.Now() },
		}

		// 渲染模板
		tmplPath := filepath.Join("templates", "themes", theme, "pages", "category.plush.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Printf("分类模板不存在，跳过: %s", tmplPath)
			continue
		}

		// 创建输出目录
		outputDir := filepath.Join("public", "static", "category")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}

		// 写入静态文件
		outputPath := filepath.Join(outputDir, category.Slug+".html")
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			return err
		}

		log.Printf("分类页面生成完成: %s", category.Slug)
	}

	return nil
}

// generateTagPages 生成标签页面
func generateTagPages(app *buffalo.App, theme string) error {
	// 获取所有标签（从数据库读取）
	tags := getAllTags()

	for _, tag := range tags {
		// 准备模板数据
		data := map[string]interface{}{
			"tag":   tag,
			"posts": getPublishedPosts(), // 简化处理，实际应该通过关联表查询
			"theme": theme,
			"site":  getSiteSettings(),
			"now":   func() time.Time { return time.Now() },
		}

		// 渲染模板
		tmplPath := filepath.Join("templates", "themes", theme, "pages", "tag.plush.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Printf("标签模板不存在，跳过: %s", tmplPath)
			continue
		}

		// 创建输出目录
		outputDir := filepath.Join("public", "static", "tag")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}

		// 写入静态文件
		outputPath := filepath.Join(outputDir, tag.Slug+".html")
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			return err
		}

		log.Printf("标签页面生成完成: %s", tag.Slug)
	}

	return nil
}

// getSiteSettings 获取站点设置
func getSiteSettings() map[string]interface{} {
	// 从数据库查询站点设置
	settings := models.Settings{}
	if err := models.DB.All(&settings); err != nil {
		log.Printf("获取站点设置失败: %v", err)
	}

	// 构建设置映射
	siteSettings := map[string]interface{}{
		"title":       "Catter CMS",
		"description": "轻量级内容管理系统",
		"url":         "",
	}

	// 覆盖默认设置
	for _, setting := range settings {
		siteSettings[setting.Key] = setting.Value
	}

	return siteSettings
}
