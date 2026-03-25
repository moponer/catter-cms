package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/handlers/api"
	"catter-cms-gin/internal/middleware"
	"catter-cms-gin/internal/models"
	"catter-cms-gin/internal/services"
)

// getSiteData 获取站点基础数据
func getSiteData() map[string]interface{} {
	siteTitle := models.GetSettingAsString("site_title")
	if siteTitle == "" {
		siteTitle = "Catter CMS"
	}
	return map[string]interface{}{
		"title":       siteTitle,
		"description": models.GetSettingAsString("site_description"),
		"url":         models.GetSettingAsString("site_url"),
		"logo":        models.GetSettingAsString("logo"),
		"favicon":     models.GetSettingAsString("favicon"),
		"year":        2026,
	}
}

// loadConfig 加载配置文件
func loadConfig() *config.Config {
	// 使用新的 InitConfig 方法
	cfg, err := config.InitConfig(
		"./config.yaml",
		"../config.yaml",
		"../../config.yaml",
	)
	if err != nil {
		log.Printf("Config initialization error: %v, using defaults", err)
		cfg = config.Get()
	}

	// 输出配置信息（不包含敏感信息）
	log.Printf("Config loaded: %s", cfg.String())

	// 检查安全警告
	if warnings := cfg.ValidateSecurity(); len(warnings) > 0 {
		for _, w := range warnings {
			log.Printf("Security warning: %s", w)
		}
	}

	return cfg
}

func main() {
	// 将日志输出重定向到 stderr，避免与 HTTP 响应混在一起
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 加载配置
	cfg := loadConfig()

	// 初始化数据库
	dbConfig := &config.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Database: cfg.Database.Database,
	}
	if err := models.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库表
	if err := models.AutoMigrate(); err != nil {
		log.Printf("Warning: Failed to auto migrate: %v", err)
	}

	// 创建默认管理员用户
	if err := models.CreateDefaultAdmin(); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	} else {
		log.Println("Default admin user created or already exists")
	}

	// 初始化示例数据
	if err := models.SeedData(); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	} else {
		log.Println("Sample data initialized successfully")
	}

	// 初始化系统设置
	if err := models.InitSettings(); err != nil {
		log.Printf("Warning: Failed to init settings: %v", err)
	} else {
		log.Println("Settings initialized")
	}

	// 初始化主题服务
	if err := services.InitThemeService(); err != nil {
		log.Printf("Warning: Failed to init theme service: %v", err)
	} else {
		log.Println("Theme service initialized")
	}

	// 启动任务调度器
	scheduler := services.GetScheduler()
	if err := scheduler.Start(); err != nil {
		log.Printf("Warning: Failed to start scheduler: %v", err)
	} else {
		log.Println("Task scheduler started")
	}

	// 创建 Gin 引擎
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// 恢复和日志中间件
	r.Use(gin.Recovery())
	if !cfg.IsProduction() {
		r.Use(gin.Logger())
	}

	// 安全头中间件
	r.Use(middleware.SecurityHeaders())

	// CORS 中间件
	r.Use(middleware.CORS())

	// 请求 ID
	r.Use(middleware.RequestIDHandler())

	// Gzip 压缩（排除管理后台和静态资源）
	r.Use(middleware.GzipExcludingPaths([]string{"/admin", "/uploads", "/themes", "/fontawesome", "/public"}))

	// 设置信任的代理
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Catter CMS is running",
			"env":     cfg.Env,
		})
	})

	// API 路由
	apiGroup := r.Group("/api")
	{
		// 公开 API
		apiGroup.POST("/auth/login", api.LoginHandler)

		// 需要认证的 API
		protected := apiGroup.Group("/")
		protected.Use(middleware.APIAuth())
		{
			protected.GET("/auth/me", api.MeHandler)

			// 仪表盘 API
			protected.GET("/dashboard/stats", api.GetDashboardStats)
			protected.GET("/dashboard/recent-posts", api.GetRecentPosts)

			// 文章 API
			protected.GET("/posts", api.ListPosts)
			protected.GET("/posts/:id", api.GetPost)
			protected.POST("/posts", api.CreatePost)
			protected.PUT("/posts/:id", api.UpdatePost)
			protected.DELETE("/posts/:id", api.DeletePost)
			protected.DELETE("/posts/batch", api.BatchDeletePosts)
			protected.PUT("/posts/batch/status", api.BatchUpdateStatus)

			// 分类 API
			protected.GET("/categories", api.ListCategories)
			protected.GET("/categories/tree", api.GetCategoryTree)
			protected.GET("/categories/:id", api.GetCategory)
			protected.POST("/categories", api.CreateCategory)
			protected.PUT("/categories/:id", api.UpdateCategory)
			protected.DELETE("/categories/:id", api.DeleteCategory)

			// 标签 API
			protected.GET("/tags", api.ListTags)
			protected.GET("/tags/:id", api.GetTag)
			protected.POST("/tags", api.CreateTag)
			protected.PUT("/tags/:id", api.UpdateTag)
			protected.DELETE("/tags/:id", api.DeleteTag)

			// 媒体 API
			protected.GET("/media", api.ListMedia)
			protected.GET("/media/:id", api.GetMedia)
			protected.POST("/media/upload", api.UploadMedia)
			protected.DELETE("/media/:id", api.DeleteMedia)

			// 主题 API
			protected.GET("/themes", api.ListThemes)
			protected.POST("/themes/upload", api.UploadTheme)
			protected.GET("/themes/:slug/config", api.GetThemeConfig)
			protected.PUT("/themes/:slug/config", api.SaveThemeConfig)
			protected.GET("/themes/:slug/config/full", api.GetThemeConfigFull)
			protected.PUT("/themes/:slug/config/full", api.SaveThemeConfigFull)
			protected.POST("/themes/:slug/config/reset", api.ResetThemeConfig)
			protected.POST("/themes/:slug/activate", api.ActivateTheme)
			protected.DELETE("/themes/:slug", api.DeleteTheme)
			protected.GET("/themes/:slug", api.GetTheme)

			// 主题模块 API
			protected.GET("/themes/:slug/modules", api.GetThemeModules)
			protected.PUT("/themes/:slug/modules/:moduleId", api.UpdateModuleConfig)
			protected.PUT("/themes/:slug/modules/order", api.UpdateModuleOrder)

			// 主题设置 Schema API
			protected.GET("/themes/:slug/schema", api.GetThemeSettingsSchema)

			// 主题配置导入/导出 API
			protected.GET("/themes/:slug/config/export", api.ExportThemeConfig)
			protected.POST("/themes/:slug/config/import", api.ImportThemeConfig)

			// 主题缓存清理 API
			protected.POST("/themes/:slug/clear-cache", api.ClearThemeCache)

			// 页面 API
			protected.GET("/pages", api.ListPages)
			protected.GET("/pages/:id", api.GetPage)
			protected.POST("/pages", api.CreatePage)
			protected.PUT("/pages/:id", api.UpdatePage)
			protected.DELETE("/pages/:id", api.DeletePage)

			// 用户 API
			protected.GET("/users", api.ListUsers)
			protected.GET("/users/:id", api.GetUser)
			protected.POST("/users", api.CreateUser)
			protected.PUT("/users/:id", api.UpdateUser)
			protected.DELETE("/users/:id", api.DeleteUser)

			// 评论 API
			protected.GET("/comments", api.ListComments)
			protected.DELETE("/comments/:id", api.DeleteComment)
			protected.POST("/comments/:id/approve", api.ApproveComment)
			protected.POST("/comments/:id/reject", api.RejectComment)

			// Sitemap管理 API
			protected.GET("/admin/sitemap/info", api.GetSitemapInfo)
			protected.POST("/admin/sitemap/rebuild", api.RebuildSitemap)

			// 缓存管理 API
			protected.GET("/admin/cache/status", api.GetCacheStatus)
			protected.POST("/admin/cache/clear", api.ClearCache)

			// 系统信息 API
			protected.GET("/admin/system/info", api.GetSystemInfo)

			// 设置 API
			protected.GET("/settings", api.GetSettings)
			protected.PUT("/settings", api.SaveSettings)

			// 定时任务 API
			protected.GET("/admin/scheduled-tasks", api.ListScheduledTasks)
			protected.GET("/admin/scheduled-tasks/:id", api.GetScheduledTask)
			protected.POST("/admin/scheduled-tasks", api.CreateScheduledTask)
			protected.PUT("/admin/scheduled-tasks/:id", api.UpdateScheduledTask)
			protected.DELETE("/admin/scheduled-tasks/:id", api.DeleteScheduledTask)
			protected.POST("/admin/scheduled-tasks/:id/toggle", api.ToggleScheduledTask)
			protected.POST("/admin/scheduled-tasks/:id/run", api.RunScheduledTask)
			protected.POST("/admin/scheduled-tasks/:id/enable", api.EnableScheduledTask)
			protected.POST("/admin/scheduled-tasks/:id/disable", api.DisableScheduledTask)
			protected.GET("/admin/scheduler/status", api.GetSchedulerStatus)
			protected.POST("/admin/scheduler/restart", api.RestartScheduler)

			// 菜单管理 API
			protected.GET("/menus", api.ListMenus)
			protected.GET("/menus/:id", api.GetMenu)
			protected.POST("/menus", api.CreateMenu)
			protected.PUT("/menus/:id", api.UpdateMenu)
			protected.DELETE("/menus/:id", api.DeleteMenu)
			protected.GET("/menus/location/:location", api.GetMenuByLocation)
			protected.GET("/menus/:id/items", api.GetMenuItems)
			protected.POST("/menus/:id/items", api.CreateMenuItem)
			protected.PUT("/menus/:id/items/:itemId", api.UpdateMenuItem)
			protected.DELETE("/menus/:id/items/:itemId", api.DeleteMenuItem)
			protected.POST("/menus/:id/items/reorder", api.ReorderMenuItems)
		}
	}

	// 前台评论 API（公开）- 必须放在 /posts/:id 之前
	r.GET("/api/posts/comments/:post_id", api.ListPostComments)
	r.POST("/api/comments", api.CreateComment)

	// 搜索 API（公开）
	r.GET("/api/search", api.GlobalSearch)

	// Sitemap（公开）
	r.GET("/sitemap.xml", api.SitemapHandler)

	// ========================================
	// 管理后台静态文件服务 - SPA 模式
	// ========================================
	adminDistPath := cfg.GetAdminDistPath()
	if _, err := os.Stat(adminDistPath); err == nil {
		// 登录页面
		r.GET("/login", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.File(adminDistPath + "/index.html")
		})

		// 管理后台静态文件（带 SPA fallback 中间件）
		// SPA 中间件会在请求的文件不存在时返回 index.html
		r.Use(middleware.SPAMiddleware(adminDistPath, "/admin"))
		r.Static("/admin", adminDistPath)

		// NoRoute 处理 - SPA fallback
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// 管理后台路由
			if strings.HasPrefix(path, "/admin") || path == "/login" {
				c.Header("Content-Type", "text/html; charset=utf-8")
				c.File(adminDistPath + "/index.html")
				return
			}
			// 其他路径返回 404
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Status(http.StatusNotFound)
			c.File("./templates/404.html")
		})
	} else {
		log.Printf("Warning: Admin dist not found at %s, admin UI won't be available", adminDistPath)
	}

	// ========================================
	// 其他静态文件服务
	// ========================================

	// FontAwesome
	fontawesomePath := "./templates/fontawesome"
	if _, err := os.Stat(fontawesomePath); err != nil {
		fontawesomePath = "../templates/fontawesome"
	}
	r.Static("/fontawesome", fontawesomePath)

	// Uploads
	uploadsPath := cfg.GetUploadsPath()
	if _, err := os.Stat(uploadsPath); err != nil {
		uploadsPath = "../uploads"
	}
	r.Static("/uploads", uploadsPath)

	// Public
	publicPath := "./public"
	if _, err := os.Stat(publicPath); err != nil {
		publicPath = "../public"
	}
	r.Static("/public", publicPath)

	// Themes
	themesPath := cfg.GetThemePath()
	if _, err := os.Stat(themesPath); err != nil {
		themesPath = "../templates/themes"
	}
	r.Static("/themes", themesPath)

	// ========================================
	// 前台页面路由
	// ========================================
	r.GET("/", homeHandler)
	r.GET("/article", articleListHandler)
	r.GET("/article/:slug", articleShowHandler)
	r.GET("/page/:slug", pageShowHandler)
	r.GET("/categories", categoriesListHandler)
	r.GET("/category/:slug", categoryHandler)
	r.GET("/tags", tagsListHandler)
	r.GET("/tag/:slug", tagHandler)
	r.GET("/archives", archivesHandler)
	r.GET("/search", searchHandler)

	// ========================================
	// 启动服务器
	// ========================================
	addr := cfg.GetFullPort()
	log.Printf("Starting Catter CMS (%s) on %s", cfg.Env, addr)
	log.Printf("Admin panel: http://localhost:%s/admin", cfg.Server.Port)
	log.Printf("API: http://localhost:%s/api", cfg.Server.Port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ========== 前台页面处理函数 ==========

func homeHandler(c *gin.Context) {
	var posts []models.Post
	models.DB.Where("status = ?", "published").
		Preload("Category").
		Preload("Tags").
		Order("COALESCE(published_at, created_at) DESC").
		Limit(10).
		Find(&posts)

	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("COALESCE(published_at, created_at) DESC").
		Limit(5).
		Find(&recentPosts)

	var categories []models.Category
	models.DB.Find(&categories)

	var tags []models.Tag
	models.DB.Find(&tags)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   "首页",
		"posts":       posts,
		"recentPosts": recentPosts,
		"categories":  categories,
		"tags":        tags,
	}

	theme := services.GetActiveTheme()
	if theme == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "没有激活的主题",
		})
		return
	}

	content, err := services.RenderTemplate("home", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "主题渲染失败",
			"error":   err.Error(),
			"theme":   theme.Name,
			"path":    theme.Path,
		})
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(content))
}

func articleShowHandler(c *gin.Context) {
	slug := c.Param("slug")

	var post models.Post
	if err := models.DB.Where("slug = ? AND status = ?", slug, "published").
		Preload("Category").
		Preload("Tags").
		First(&post).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Post Not Found</h1></body></html>"))
		return
	}

	models.DB.Model(&post).Update("views", post.Views+1)

	var categories []models.Category
	models.DB.Find(&categories)
	var tags []models.Tag
	models.DB.Find(&tags)

	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentPosts)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   post.Title,
		"post":        post,
		"categories":  categories,
		"tags":        tags,
		"recentPosts": recentPosts,
	}

	if content, err := services.RenderTemplate("article_show", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// 默认页面 - XSS 安全修复：使用 html.EscapeString 转义
	defaultHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Catter CMS</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
               background: #f5f5f5; min-height: 100vh; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; margin-bottom: 20px; }
        .meta { color: #666; font-size: 0.9em; margin-bottom: 20px; }
        .content { line-height: 1.8; color: #555; }
        .back-link { display: inline-block; margin-top: 20px; color: #667eea; text-decoration: none; }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <div class="meta">Views: %d | Published: %s</div>
        <div class="content">%s</div>
        <a href="/" class="back-link">Back to Home</a>
    </div>
</body>
</html>`,
		html.EscapeString(post.Title),
		html.EscapeString(post.Title),
		post.Views,
		post.PublishedAt.Format("2006-01-02"),
		html.EscapeString(post.Content))

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

func pageShowHandler(c *gin.Context) {
	slug := c.Param("slug")

	var page models.Page
	if err := models.DB.Where("slug = ? AND status = ?", slug, "published").First(&page).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Page Not Found</h1></body></html>"))
		return
	}

	models.DB.Model(&page).Update("views", page.Views+1)

	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentPosts)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   page.Title,
		"page":        page,
		"recentPosts": recentPosts,
	}

	templateName := "page"
	if page.Template != "" && page.Template != "default" {
		templateName = "page-" + page.Template
	}

	if content, err := services.RenderTemplate(templateName, data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	if templateName != "page" {
		if content, err := services.RenderTemplate("page", data); err == nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, string(content))
			return
		}
	}

	defaultHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Catter CMS</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
               background: #f5f5f5; min-height: 100vh; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; margin-bottom: 20px; }
        .content { line-height: 1.8; color: #555; }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <div class="content">%s</div>
    </div>
</body>
</html>`,
		html.EscapeString(page.Title),
		html.EscapeString(page.Title),
		html.EscapeString(page.Content))

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

func categoryHandler(c *gin.Context) {
	slug := c.Param("slug")

	var category models.Category
	if err := models.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Category Not Found</h1></body></html>"))
		return
	}

	var posts []models.Post
	models.DB.Where("category_id = ? AND status = ?", category.ID, "published").
		Preload("Category").
		Order("published_at DESC").
		Find(&posts)

	var categories []models.Category
	models.DB.Find(&categories)
	var tags []models.Tag
	models.DB.Find(&tags)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  category.Name,
		"category":   category,
		"posts":      posts,
		"categories": categories,
		"tags":       tags,
	}

	if content, err := services.RenderTemplate("category", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// 默认页面 - XSS 安全修复
	postsHTML := ""
	for _, post := range posts {
		postsHTML += fmt.Sprintf(`<div class="post-item"><h3><a href="/article/%s">%s</a></h3></div>`,
			html.EscapeString(post.Slug),
			html.EscapeString(post.Title))
	}
	if postsHTML == "" {
		postsHTML = "<p>No posts in this category.</p>"
	}

	defaultHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Catter CMS</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
               background: #f5f5f5; min-height: 100vh; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; margin-bottom: 20px; }
        .post-item { padding: 10px 0; border-bottom: 1px solid #eee; }
        .post-item a { color: #667eea; text-decoration: none; }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s</h1>
        <p>%s</p>
        <div class="posts">%s</div>
    </div>
</body>
</html>`,
		html.EscapeString(category.Name),
		html.EscapeString(category.Name),
		html.EscapeString(category.Description),
		postsHTML)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

func articleListHandler(c *gin.Context) {
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("pageSize"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	var posts []models.Post
	var total int64
	models.DB.Where("status = ?", "published").
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts)
	models.DB.Model(&models.Post{}).Where("status = ?", "published").Count(&total)

	var categories []models.Category
	var tags []models.Tag
	models.DB.Find(&categories)
	models.DB.Find(&tags)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  "文章",
		"posts":      posts,
		"categories": categories,
		"tags":       tags,
		"pagination": map[string]interface{}{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	content, err := services.RenderTemplate("article_list", data)
	if err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	log.Printf("[ERROR] article_list template render failed: %v", err)
	c.String(http.StatusOK, "<h1>文章列表</h1><p>请创建 article_list.html 模板</p>")
}

func categoriesListHandler(c *gin.Context) {
	var categories []models.Category
	models.DB.Find(&categories)

	type CategoryWithCount struct {
		models.Category
		PostCount int `json:"post_count"`
	}
	var categoriesWithCount []CategoryWithCount
	for _, cat := range categories {
		var count int64
		models.DB.Model(&models.Post{}).Where("category_id = ? AND status = ?", cat.ID, "published").Count(&count)
		categoriesWithCount = append(categoriesWithCount, CategoryWithCount{
			Category:  cat,
			PostCount: int(count),
		})
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  "分类",
		"categories": categoriesWithCount,
	}

	if content, err := services.RenderTemplate("categories", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>分类列表页</h1><p>请创建 categories.plush.html 模板</p>")
}

func tagsListHandler(c *gin.Context) {
	var tags []models.Tag
	models.DB.Find(&tags)

	type TagWithCount struct {
		models.Tag
		PostCount int `json:"post_count"`
	}
	var tagsWithCount []TagWithCount
	for _, tag := range tags {
		var count int64
		models.DB.Table("post_tags").Where("tag_id = ?", tag.ID).Count(&count)
		tagsWithCount = append(tagsWithCount, TagWithCount{
			Tag:       tag,
			PostCount: int(count),
		})
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "标签",
		"tags":      tagsWithCount,
	}

	if content, err := services.RenderTemplate("tags", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>标签列表页</h1><p>请创建 tags.plush.html 模板</p>")
}

func tagHandler(c *gin.Context) {
	slug := c.Param("slug")

	var tag models.Tag
	if err := models.DB.Where("slug = ?", slug).First(&tag).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Tag Not Found</h1></body></html>"))
		return
	}

	var posts []models.Post
	models.DB.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = ?", tag.ID, "published").
		Preload("Category").
		Order("posts.published_at DESC").
		Find(&posts)

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": tag.Name,
		"tag":       tag,
		"posts":     posts,
	}

	if content, err := services.RenderTemplate("tag", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// XSS 安全修复
	c.String(http.StatusOK, fmt.Sprintf("<h1>标签: %s</h1><p>请创建 tag.plush.html 模板</p>",
		html.EscapeString(tag.Name)))
}

func archivesHandler(c *gin.Context) {
	var posts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Find(&posts)

	archives := make(map[string][]models.Post)
	for _, post := range posts {
		if post.PublishedAt != nil {
			key := post.PublishedAt.Format("2006年01月")
			archives[key] = append(archives[key], post)
		}
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "文章归档",
		"archives":  archives,
		"posts":     posts,
	}

	if content, err := services.RenderTemplate("archives", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>文章归档</h1><p>请创建 archives.plush.html 模板</p>")
}

func searchHandler(c *gin.Context) {
	query := c.Query("q")

	var posts []models.Post
	if query != "" {
		models.DB.Where("status = ? AND (title LIKE ? OR content LIKE ?)",
			"published",
			"%"+query+"%",
			"%"+query+"%").
			Preload("Category").
			Order("published_at DESC").
			Find(&posts)
	}

	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "搜索",
		"query":     query,
		"posts":     posts,
		"count":     len(posts),
	}

	if content, err := services.RenderTemplate("search", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// XSS 安全修复
	c.String(http.StatusOK, fmt.Sprintf("<h1>搜索: %s</h1><p>请创建 search.plush.html 模板</p>",
		html.EscapeString(query)))
}
