package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/cache"
	"catter-cms-gin/internal/handlers/api"
	"catter-cms-gin/internal/handlers/frontend"
	"catter-cms-gin/internal/middleware"
	"catter-cms-gin/internal/models"
	"catter-cms-gin/internal/services"
)

// loadConfig 加载配置文件
func loadConfig() *config.Config {
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

	// 修复历史数据：已发布但 PublishedAt 为空的文章，用 created_at 回填
	models.DB.Exec("UPDATE posts SET published_at = created_at WHERE status = ? AND published_at IS NULL", "published")

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

	// 初始化存储服务
	uploadsPath := cfg.GetUploadsPath()
	if err := api.InitStorage(cfg, uploadsPath); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	log.Printf("Storage initialized: %s", api.GetStorage().Type())

	// 初始化 Redis 缓存
	if err := cache.Init(cfg.Redis); err != nil {
		log.Printf("Warning: Redis init failed (caching disabled): %v", err)
	} else if cfg.Redis.Enabled {
		log.Println("Redis cache initialized")
	}

	// 启动任务调度器
	log.Println("[Main] Starting scheduler...")
	scheduler := services.GetScheduler()
	if err := scheduler.Start(); err != nil {
		log.Printf("Warning: Failed to start scheduler: %v", err)
	} else {
		log.Println("Task scheduler started")
	}
	log.Println("[Main] Scheduler done, creating Gin engine...")

	// 创建 Gin 引擎
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// 恢复和日志中间件
	r.Use(middleware.CustomRecovery())
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
		apiGroup.POST("/auth/login", middleware.LoginRateLimiter(), api.LoginHandler)
		apiGroup.POST("/auth/logout", api.LogoutHandler)

		// 需要认证的 API
		protected := apiGroup.Group("/")
		protected.Use(middleware.APIAuth())
		protected.Use(middleware.CSRFProtection())
		{
			protected.GET("/auth/me", api.MeHandler)
			protected.PUT("/auth/profile", api.UpdateProfileHandler)
			protected.PUT("/auth/password", api.ChangePasswordHandler)

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
			protected.DELETE("/categories/batch", api.BatchDeleteCategories)

			// 标签 API
			protected.GET("/tags", api.ListTags)
			protected.GET("/tags/:id", api.GetTag)
			protected.POST("/tags", api.CreateTag)
			protected.PUT("/tags/:id", api.UpdateTag)
			protected.DELETE("/tags/:id", api.DeleteTag)
			protected.DELETE("/tags/batch", api.BatchDeleteTags)

			// 媒体 API
			protected.GET("/media", api.ListMedia)
			protected.GET("/media/:id", api.GetMedia)
			protected.POST("/media/upload", api.UploadMedia)
			protected.POST("/media/upload-cover", api.UploadCover) // 封面图上传（自动压缩）
			protected.DELETE("/media/:id", api.DeleteMedia)
			protected.POST("/media/batch-delete", api.BatchDeleteMedia)

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

			// 主题页面模板 API
			protected.GET("/themes/:slug/templates", api.GetPageTemplates)

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
			protected.DELETE("/pages/batch", api.BatchDeletePages)

			// 用户 API
			protected.GET("/users", api.ListUsers)
			protected.GET("/users/:id", api.GetUser)
			protected.POST("/users", api.CreateUser)
			protected.PUT("/users/:id", api.UpdateUser)
			protected.DELETE("/users/:id", api.DeleteUser)
			protected.POST("/users/:id/reset-password", api.ResetUserPassword)

			// 评论管理 API
			protected.GET("/admin/comments", api.ListComments)
			protected.DELETE("/admin/comments/:id", api.DeleteComment)
			protected.PUT("/admin/comments/:id", api.UpdateComment)
			protected.POST("/admin/comments/:id/approve", api.ApproveComment)
			protected.POST("/admin/comments/:id/reject", api.RejectComment)
			protected.POST("/admin/comments/batch/delete", api.BatchDeleteComments)
			protected.POST("/admin/comments/batch/approve", api.BatchApproveComments)

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

			// 存储配置 API
			protected.GET("/storage/config", api.GetStorageConfig)
			protected.PUT("/storage/config", api.SaveStorageConfig)
			protected.POST("/storage/config/test", api.TestStorageConnection)

			// Redis 缓存配置 API
			protected.GET("/redis/config", api.GetRedisConfig)
			protected.PUT("/redis/config", api.SaveRedisConfig)
			protected.POST("/redis/config/test", api.TestRedisConnection)

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
			protected.GET("/admin/scheduled-tasks/:id/executions", api.GetTaskExecutions)
			protected.GET("/admin/scheduled-tasks/executions/all", api.GetAllExecutions)
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

			// AI Key 管理
			protected.GET("/ai/keys", api.ListAIKeys)
			protected.POST("/ai/keys", api.CreateAIKey)
			protected.DELETE("/ai/keys/:id", api.DeleteAIKey)
			protected.GET("/ai/keys/:id/logs", api.GetAIKeyLogs)
			protected.GET("/ai/setup-commands", api.GetAISetupCommands)
		}
	}

	// 前台评论 API（公开）- 必须放在 /posts/:id 之前
	r.GET("/api/posts/comments/:post_id", api.ListPostComments)
	r.POST("/api/comments", api.CreateComment)

	// 搜索 API（公开）
	r.GET("/api/search", api.GlobalSearch)

	// ========================================
	// AI 集成 API
	// ========================================

	// AI Skill 和 Schema（公开，无需认证）
	r.GET("/api/ai/skill", api.GetAISkill)
	r.GET("/api/ai/schema", api.GetAISchema)
	r.GET("/api/ai/themes-skill", api.GetThemeDevSkill)

	// AI v1 API（使用 AI Key 认证）
	aiGroup := r.Group("/api/v1")
	aiGroup.Use(middleware.AIAuth())
	{
		// 文章
		aiGroup.GET("/posts", api.AIListPosts)
		aiGroup.GET("/posts/:id", api.AIGetPost)
		aiGroup.POST("/posts", middleware.RequireAIPermission("write"), api.AICreatePost)
		aiGroup.PUT("/posts/:id", middleware.RequireAIPermission("write"), api.AIUpdatePost)
		aiGroup.DELETE("/posts/:id", middleware.RequireAIPermission("admin"), api.AIDeletePost)

		// 页面
		aiGroup.GET("/pages", api.AIListPages)
		aiGroup.GET("/pages/:id", api.AIGetPage)

		// 分类
		aiGroup.GET("/categories", api.AIListCategories)

		// 标签
		aiGroup.GET("/tags", api.AIListTags)

		// 媒体
		aiGroup.GET("/media", api.AIListMedia)

		// 评论
		aiGroup.GET("/comments", api.AIListComments)

		// 设置
		aiGroup.GET("/settings", api.AIGetSettings)
	}

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

		// Favicon.ico 处理（避免 404 HTML 页面被当作脚本执行）
		r.GET("/favicon.ico", func(c *gin.Context) {
			faviconPath := adminDistPath + "/favicon.ico"
			if _, err := os.Stat(faviconPath); err == nil {
				c.File(faviconPath)
			} else {
				c.Status(http.StatusNoContent)
			}
		})

		// 管理后台静态文件（带 SPA fallback 中间件）
		// SPA 中间件会在请求的文件不存在时返回 index.html
		r.Use(middleware.SPAMiddleware(adminDistPath, "/admin"))
		r.Static("/admin", adminDistPath)
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
	r.GET("/", frontend.HomeHandler)
	r.GET("/article", frontend.ArticleListHandler)
	r.GET("/article/:slug", frontend.ArticleShowHandler)
	r.GET("/page/:slug", frontend.PageShowHandler)
	r.GET("/categories", frontend.CategoriesListHandler)
	r.GET("/category/:slug", frontend.CategoryHandler)
	r.GET("/tags", frontend.TagsListHandler)
	r.GET("/tag/:slug", frontend.TagHandler)
	r.GET("/archives", frontend.ArchivesHandler)
	r.GET("/search", frontend.SearchHandler)

	// ========================================
	// NoRoute 处理 - 必须在所有路由和静态文件服务之后注册
	// ========================================
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 管理后台路由
		if strings.HasPrefix(path, "/admin") || path == "/login" {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.File(cfg.GetAdminDistPath() + "/index.html")
			return
		}
		// 其他路径返回 404（主题优先）
		frontend.Render404(c)
	})

	// ========================================
	// 启动服务器（Graceful Shutdown）
	// ========================================
	log.Println("[Main] Creating HTTP server...")
	addr := cfg.GetFullPort()
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Printf("[Main] HTTP server created, addr=%s\n", addr)

	// 启动服务器 goroutine
	go func() {
		log.Println("[Main] Server goroutine starting...")
		log.Printf("Starting Catter CMS (%s) on %s", cfg.Env, addr)
		log.Printf("Admin panel: http://localhost:%s/admin", cfg.Server.Port)
		log.Printf("API: http://localhost:%s/api", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待 SIGINT / SIGTERM 信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 停止主题文件监控 watcher
	services.StopThemeWatcher()

	// 关闭 Redis 连接
	cache.Close()

	// 给正在处理的请求最多 5 秒完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}