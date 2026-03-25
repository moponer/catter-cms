package main

import (
	"fmt"
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
	// 尝试多个路径加载配置文件
	paths := []string{
		"./config.yaml",     // 当前目录
		"../config.yaml",    // 上级目录（从gin-app运行）
		"../../config.yaml", // 上两级目录
	}

	for _, path := range paths {
		cfg, err := config.Load(path)
		if err == nil {
			log.Printf("Config loaded from: %s", path)
			return cfg
		}
	}

	log.Printf("Config file not found, using default config")
	return config.Get()
}

func main() {
	// 将日志输出重定向到 stderr，避免与 HTTP 响应混在一起
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 加载配置 - 优先从根目录，其次从gin-app目录
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

	// 创建 Gin 引擎 - 强制使用 ReleaseMode 避免调试日志干扰 HTTP 响应
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	// 禁用默认的Logger中间件，使用自定义的
	// r.Use(gin.Logger())
	r.Use(gin.Recovery())                 // 错误恢复
	r.Use(middleware.SecurityHeaders())   // 安全头
	r.Use(middleware.CORS())              // CORS
	r.Use(middleware.RequestIDHandler())  // 请求ID
	// Gzip 压缩 - 排除管理后台和静态资源路径
	// 管理后台：SPA 模式，HTML 通过 NoRoute 返回，不需要 Gzip；静态资源已是预压缩格式
	// 前台主题：Go 模板渲染的 HTML，不经过静态文件服务，应该启用 Gzip
	r.Use(middleware.GzipExcludingPaths([]string{"/admin", "/uploads", "/themes", "/fontawesome", "/public"}))

	// 设置信任的代理，避免安全警告
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Catter CMS is running",
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

	// 管理后台静态文件 - 使用相对路径
	adminDistPath := "../admin/dist"
	if _, err := os.Stat(adminDistPath); err != nil {
		// 尝试当前目录下的 admin/dist
		adminDistPath = "./admin/dist"
	}

	if _, err := os.Stat(adminDistPath); err == nil {
		// 登录页面路由
		r.GET("/login", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.File(adminDistPath + "/index.html")
		})

		// 注册 /admin 的静态文件服务（通过 Group + 中间件实现）
		// 中间件：只允许磁盘上存在的文件通过，否则让 NoRoute 处理（SPA fallback）
		adminGroup := r.Group("/admin")
		adminGroup.Use(func(c *gin.Context) {
			filePath := adminDistPath + c.Request.URL.Path[len("/admin"):]
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				c.Abort() // 阻止后续 handler，让 NoRoute 生效
			}
		})
		adminGroup.Static("/", adminDistPath)

		// SPA fallback - 处理所有未被静态文件处理的 /admin 路径
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/admin") {
				c.Header("Content-Type", "text/html; charset=utf-8")
				c.File(adminDistPath + "/index.html")
				return
			}
			// 对于其他路径，返回404页面
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Status(http.StatusNotFound)
			c.File("./templates/404.html")
		})
	} else {
		log.Printf("Warning: Admin dist not found at %s, admin UI won't be available", adminDistPath)
	}

	// 静态文件必须放在通配符路由之前
	// 检测静态文件路径（支持从根目录或 gin-app 目录运行）

	// 检测 templates/fontawesome 路径
	fontawesomePath := "./templates/fontawesome"
	if _, err := os.Stat(fontawesomePath); err != nil {
		fontawesomePath = "../templates/fontawesome"
	}
	r.Static("/fontawesome", fontawesomePath)

	// 检测 uploads 路径
	uploadsPath := "./uploads"
	if _, err := os.Stat(uploadsPath); err != nil {
		uploadsPath = "../uploads"
	}
	r.Static("/uploads", uploadsPath)

	// 检测 public 路径
	publicPath := "./public"
	if _, err := os.Stat(publicPath); err != nil {
		publicPath = "../public"
	}
	r.Static("/public", publicPath)

	// 检测 themes 路径（主题静态资源）
	themesPath := "./templates/themes"
	if _, err := os.Stat(themesPath); err != nil {
		themesPath = "../templates/themes"
	}
	r.Static("/themes", themesPath)

	// 前台页面路由（放在静态文件之后）
	r.GET("/", homeHandler)
	r.GET("/article", articleListHandler)          // 文章列表页
	r.GET("/article/:slug", articleShowHandler)    // 文章详情页
	r.GET("/page/:slug", pageShowHandler)       // 独立页面
	r.GET("/categories", categoriesListHandler) // 分类列表页
	r.GET("/category/:slug", categoryHandler)   // 分类文章列表页
	r.GET("/tags", tagsListHandler)             // 标签列表页
	r.GET("/tag/:slug", tagHandler)             // 标签文章列表页
	r.GET("/archives", archivesHandler)         // 文章归档页
	r.GET("/search", searchHandler)             // 搜索页

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ========== 前台页面处理函数 ==========

func homeHandler(c *gin.Context) {
	// 获取已发布的文章
	var posts []models.Post
	models.DB.Where("status = ?", "published").
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC").
		Limit(10).
		Find(&posts)

	// 获取最近文章（用于侧边栏）
	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentPosts)

	// 获取分类
	var categories []models.Category
	models.DB.Find(&categories)

	// 获取标签
	var tags []models.Tag
	models.DB.Find(&tags)

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   "首页",
		"posts":       posts,
		"recentPosts": recentPosts,
		"categories":  categories,
		"tags":        tags,
	}

	// 尝试渲染主题模板
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
		// 渲染失败，返回详细错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "主题渲染失败",
			"error":   err.Error(),
			"theme":    theme.Name,
			"path":     theme.Path,
		})
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(content))
}

func articleShowHandler(c *gin.Context) {
	slug := c.Param("slug")

	// 获取文章
	var post models.Post
	if err := models.DB.Where("slug = ? AND status = ?", slug, "published").
		Preload("Category").
		Preload("Tags").
		First(&post).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Post Not Found</h1></body></html>"))
		return
	}

	// 增加浏览量
	models.DB.Model(&post).Update("views", post.Views+1)

	// 获取分类
	var categories []models.Category
	models.DB.Find(&categories)

	// 获取标签
	var tags []models.Tag
	models.DB.Find(&tags)

	// 获取最近文章（用于侧边栏）
	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentPosts)

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   post.Title,
		"post":        post,
		"categories":  categories,
		"tags":        tags,
		"recentPosts": recentPosts,
	}

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("article_show", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// 如果主题渲染失败，使用默认页面
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
</html>`, post.Title, post.Title, post.Views, post.PublishedAt.Format("2006-01-02"), post.Content)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

func pageShowHandler(c *gin.Context) {
	slug := c.Param("slug")

	// 获取页面
	var page models.Page
	if err := models.DB.Where("slug = ? AND status = ?", slug, "published").First(&page).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Page Not Found</h1></body></html>"))
		return
	}

	// 增加浏览量
	models.DB.Model(&page).Update("views", page.Views+1)

	// 获取最近文章（用于侧边栏）
	var recentPosts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentPosts)

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":        siteData,
		"pageTitle":   page.Title,
		"page":        page,
		"recentPosts": recentPosts,
	}

	// 使用页面指定的模板，如果没有则使用默认的 "page" 模板
	templateName := "page"
	if page.Template != "" && page.Template != "default" {
		templateName = "page-" + page.Template
	}

	if content, err := services.RenderTemplate(templateName, data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// 如果指定模板不存在，回退到默认模板
	if templateName != "page" {
		if content, err := services.RenderTemplate("page", data); err == nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, string(content))
			return
		}
	}

	// 默认页面
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
</html>`, page.Title, page.Title, page.Content)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

func categoryHandler(c *gin.Context) {
	slug := c.Param("slug")

	// 获取分类
	var category models.Category
	if err := models.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Category Not Found</h1></body></html>"))
		return
	}

	// 获取该分类下的已发布文章
	var posts []models.Post
	models.DB.Where("category_id = ? AND status = ?", category.ID, "published").
		Preload("Category").
		Order("published_at DESC").
		Find(&posts)

	// 获取所有分类和标签（用于侧边栏）
	var categories []models.Category
	models.DB.Find(&categories)
	var tags []models.Tag
	models.DB.Find(&tags)

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  category.Name,
		"category":   category,
		"posts":      posts,
		"categories": categories,
		"tags":       tags,
	}

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("category", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	// 默认页面
	postsHTML := ""
	for _, post := range posts {
		postsHTML += fmt.Sprintf(`<div class="post-item"><h3><a href="/blog/%s">%s</a></h3></div>`, post.Slug, post.Title)
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
</html>`, category.Name, category.Name, category.Description, postsHTML)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, defaultHTML)
}

// articleListHandler 文章列表页
func articleListHandler(c *gin.Context) {
	// 获取分页参数
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("pageSize"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	// 获取文章列表
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

	// 获取分类和标签
	var categories []models.Category
	var tags []models.Tag
	models.DB.Find(&categories)
	models.DB.Find(&tags)

	// 准备模板数据
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

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("article_list", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	} else {
		log.Printf("[ERROR] article_list template render failed: %v", err)
	}

	// 默认页面
	c.String(http.StatusOK, "<h1>文章列表</h1><p>请创建 article_list.html 模板</p>")
}

// categoriesListHandler 分类列表页
func categoriesListHandler(c *gin.Context) {
	// 获取所有分类及文章数
	var categories []models.Category
	models.DB.Find(&categories)

	// 为每个分类计算文章数
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

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":       siteData,
		"pageTitle":  "分类",
		"categories": categoriesWithCount,
	}

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("categories", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>分类列表页</h1><p>请创建 categories.plush.html 模板</p>")
}

// tagsListHandler 标签列表页
func tagsListHandler(c *gin.Context) {
	// 获取所有标签及文章数
	var tags []models.Tag
	models.DB.Find(&tags)

	// 为每个标签计算文章数
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

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "标签",
		"tags":      tagsWithCount,
	}

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("tags", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>标签列表页</h1><p>请创建 tags.plush.html 模板</p>")
}

// tagHandler 标签文章列表页
func tagHandler(c *gin.Context) {
	slug := c.Param("slug")

	// 获取标签
	var tag models.Tag
	if err := models.DB.Where("slug = ?", slug).First(&tag).Error; err != nil {
		c.HTML(http.StatusNotFound, "text/html", []byte("<html><body><h1>404 - Tag Not Found</h1></body></html>"))
		return
	}

	// 获取该标签下的文章
	var posts []models.Post
	models.DB.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = ?", tag.ID, "published").
		Preload("Category").
		Order("posts.published_at DESC").
		Find(&posts)

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": tag.Name,
		"tag":       tag,
		"posts":     posts,
	}

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("tag", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>标签: "+tag.Name+"</h1><p>请创建 tag.plush.html 模板</p>")
}

// archivesHandler 文章归档页
func archivesHandler(c *gin.Context) {
	// 获取所有已发布文章，按日期分组
	var posts []models.Post
	models.DB.Where("status = ?", "published").
		Order("published_at DESC").
		Find(&posts)

	// 按年月分组
	archives := make(map[string][]models.Post)
	for _, post := range posts {
		if post.PublishedAt != nil {
			key := post.PublishedAt.Format("2006年01月")
			archives[key] = append(archives[key], post)
		}
	}

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "文章归档",
		"archives":  archives,
		"posts":     posts,
	}

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("archives", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>文章归档</h1><p>请创建 archives.plush.html 模板</p>")
}

// searchHandler 搜索页
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

	// 准备模板数据
	siteData := getSiteData()
	data := map[string]interface{}{
		"site":      siteData,
		"pageTitle": "搜索",
		"query":     query,
		"posts":     posts,
		"count":     len(posts),
	}

	// 尝试渲染主题模板
	if content, err := services.RenderTemplate("search", data); err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(content))
		return
	}

	c.String(http.StatusOK, "<h1>搜索: "+query+"</h1><p>请创建 search.plush.html 模板</p>")
}
