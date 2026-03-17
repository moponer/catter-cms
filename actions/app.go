package actions

import (
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"catter-cms/models"
	"catter-cms/middleware"
	"catter-cms/services"
	admin "catter-cms/actions/admin"
	front "catter-cms/actions/front"
	api "catter-cms/actions/api"
)

// ENV 环境变量
var ENV = envy.Get("GO_ENV", "development")

// App 应用实例
var App *buffalo.App

// T 翻译函数
var T *i18n.Translator

// serveAdminSPA 返回管理后台 SPA 入口页面
func serveAdminSPA(c buffalo.Context) error {
	c.Response().Header().Set("Content-Type", "text/html")
	http.ServeFile(c.Response(), c.Request(), "public/index.html")
	return nil
}

// init 初始化应用
func init() {
	// 设置端口环境变量（Buffalo 优先使用 PORT 环境变量）
	os.Setenv("PORT", "3001")
	
	// 创建应用
	App = buffalo.New(buffalo.Options{
		Env:         ENV,
		SessionName: "_catter_cms_session",
		Host:        "0.0.0.0:3001",
	})
	
	// 全局中间件
	App.Use(paramlogger.ParameterLogger)
	
	// CORS 配置
	App.Use(func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			w := c.Response()
			r := c.Request()
			
			// 设置CORS头
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			
			// 处理OPTIONS请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return nil
			}
			
			return next(c)
		}
	})
	
	// 数据库事务中间件
	App.Use(popmw.Transaction(models.DB))
	
	// API 路由
	apiGroup := App.Group("/api")
	
	// 认证API（不需要认证保护）
	authAPI := apiGroup.Group("/auth")
	authAPI.POST("/login", api.LoginHandler)
	
	// 需要认证保护的API
	protectedAPI := apiGroup.Group("/")
	protectedAPI.Use(middleware.APIAuth)
	protectedAPI.GET("/auth/me", api.MeHandler)
	
	// 文章API
	protectedAPI.GET("/posts", admin.PostsListAPI)
	protectedAPI.GET("/posts/{id}", admin.PostDetailAPI)
	protectedAPI.POST("/posts", admin.PostCreateAPI)
	protectedAPI.PUT("/posts/{id}", admin.PostUpdateAPI)
	protectedAPI.DELETE("/posts/{id}", admin.PostDeleteAPI)
	
	// 分类API（前台）
	protectedAPI.GET("/categories", api.CategoriesIndex)
	protectedAPI.GET("/categories/{id}", api.CategoriesShow)
	
	// 分类管理API（后台）
	adminCategories := protectedAPI.Group("/admin/categories")
	adminCategories.GET("/", admin.CategoriesListAPI)
	adminCategories.GET("/tree", admin.CategoriesTreeAPI)
	adminCategories.GET("/{id}", admin.CategoryDetailAPI)
	adminCategories.POST("/", admin.CategoryCreateAPI)
	adminCategories.PUT("/{id}", admin.CategoryUpdateAPI)
	adminCategories.DELETE("/{id}", admin.CategoryDeleteAPI)
	adminCategories.POST("/reorder", admin.CategoryReorderAPI)
	
	// 标签API
	protectedAPI.GET("/tags", admin.TagsListAPI)
	protectedAPI.POST("/tags", admin.TagsCreateAPI)
	protectedAPI.PUT("/tags/{id}", admin.TagsUpdateAPI)
	protectedAPI.DELETE("/tags/{id}", admin.TagsDeleteAPI)
	
	// 媒体API
	protectedAPI.GET("/media", admin.MediaListAPI)
	protectedAPI.GET("/media/{id}", admin.MediaDetailAPI)
	protectedAPI.POST("/media/upload", admin.MediaUploadAPI)
	protectedAPI.DELETE("/media/{id}", admin.MediaDeleteAPI)
	protectedAPI.POST("/media/batch-delete", admin.MediaBatchDeleteAPI)
	
	// 主题API
	protectedAPI.GET("/themes", admin.ThemesListAPI)
	protectedAPI.GET("/themes/{id}", admin.ThemeDetailAPI)
	protectedAPI.GET("/themes/{slug}/config", admin.ThemesGetConfigAPI)
	protectedAPI.PUT("/themes/{slug}/config", admin.ThemesSaveConfigAPI)
	protectedAPI.POST("/themes/{id}/activate", admin.ThemesActivateAPI)
	protectedAPI.DELETE("/themes/{id}", admin.ThemesDeleteHandler)
	
	// 新主题配置管理API
	protectedAPI.GET("/themes/{slug}/config/full", admin.GetThemeConfigAPI)
	protectedAPI.PUT("/themes/{slug}/config/full", admin.UpdateThemeConfigAPI)
	protectedAPI.GET("/themes/{slug}/modules", admin.GetThemeModulesAPI)
	protectedAPI.GET("/themes/{slug}/modules/{moduleId}", admin.GetThemeModuleAPI)
	protectedAPI.PUT("/themes/{slug}/modules/{moduleId}", admin.SetThemeModuleAPI)
	protectedAPI.POST("/themes/{slug}/modules/{moduleId}/enable", admin.EnableThemeModuleAPI)
	protectedAPI.POST("/themes/{slug}/modules/{moduleId}/disable", admin.DisableThemeModuleAPI)
	protectedAPI.GET("/themes/{slug}/modules/enabled", admin.GetEnabledThemeModulesAPI)
	protectedAPI.GET("/themes/{slug}/widgets", admin.GetThemeWidgetsAPI)
	protectedAPI.GET("/themes/{slug}/widgets/{widgetId}", admin.GetThemeWidgetAPI)
	protectedAPI.PUT("/themes/{slug}/widgets/{widgetId}", admin.SetThemeWidgetAPI)
	protectedAPI.GET("/themes/{slug}/widgets/area/{area}", admin.GetThemeWidgetsByAreaAPI)
	protectedAPI.GET("/themes/{slug}/config/form", admin.GetThemeConfigFormAPI)
	protectedAPI.GET("/themes/{slug}/custom-fields/{contentType}", admin.GetThemeCustomFieldsAPI)
	protectedAPI.GET("/themes/{slug}/widget-areas", admin.GetThemeWidgetAreasAPI)
	protectedAPI.GET("/themes/{slug}/menu-locations", admin.GetThemeMenuLocationsAPI)
	protectedAPI.POST("/themes/{slug}/config/reset", admin.ResetThemeConfigAPI)
	protectedAPI.POST("/themes/{slug}/config/preview", admin.PreviewThemeConfigAPI)
	protectedAPI.POST("/themes/upload", admin.UploadThemeAPI)
	
	// 页面API
	protectedAPI.GET("/pages", admin.PagesListAPI)
	protectedAPI.GET("/pages/{id}", admin.PageDetailAPI)
	protectedAPI.POST("/pages", admin.PageCreateAPI)
	protectedAPI.PUT("/pages/{id}", admin.PageUpdateAPI)
	protectedAPI.DELETE("/pages/{id}", admin.PageDeleteAPI)
	protectedAPI.POST("/pages/reorder", admin.PageReorderAPI)
	
	// 仪表盘API
	protectedAPI.GET("/dashboard/stats", admin.DashboardStatsAPI)
	protectedAPI.GET("/dashboard/recent-posts", admin.DashboardRecentPostsAPI)
	
	// 用户管理API
	protectedAPI.GET("/users", admin.UsersListAPI)
	protectedAPI.GET("/users/{id}", admin.UserDetailAPI)
	protectedAPI.POST("/users", admin.UserCreateAPI)
	protectedAPI.PUT("/users/{id}", admin.UserUpdateAPI)
	protectedAPI.DELETE("/users/{id}", admin.UserDeleteAPI)
	protectedAPI.POST("/users/{id}/reset-password", admin.UserResetPasswordAPI)
	
	// 系统设置API
	protectedAPI.GET("/settings", admin.SettingsGetAPI)
	protectedAPI.PUT("/settings", admin.SettingsSaveAPI)
	protectedAPI.POST("/settings/test-email", admin.SettingsTestEmailAPI)

	// 评论API
	protectedAPI.GET("/admin/comments", admin.CommentsListAPI)
	protectedAPI.GET("/admin/comments/{id}", admin.CommentDetailAPI)
	protectedAPI.POST("/admin/comments", admin.CommentCreateAPI)
	protectedAPI.PUT("/admin/comments/{id}", admin.CommentUpdateAPI)
	protectedAPI.DELETE("/admin/comments/{id}", admin.CommentDeleteAPI)
	protectedAPI.POST("/admin/comments/{id}/approve", admin.CommentApproveAPI)
	protectedAPI.POST("/admin/comments/{id}/reject", admin.CommentRejectAPI)
	protectedAPI.POST("/admin/comments/batch/delete", admin.CommentBatchDeleteAPI)
	protectedAPI.POST("/admin/comments/batch/approve", admin.CommentBatchApproveAPI)

	// 搜索API
	protectedAPI.POST("/admin/search/rebuild", admin.RebuildSearchIndexAPI)

	// 定时任务API
	protectedAPI.GET("/admin/scheduled-tasks", admin.ScheduledTasksListAPI)
	protectedAPI.GET("/admin/scheduled-tasks/{id}", admin.ScheduledTaskDetailAPI)
	protectedAPI.POST("/admin/scheduled-tasks", admin.ScheduledTaskCreateAPI)
	protectedAPI.PUT("/admin/scheduled-tasks/{id}", admin.ScheduledTaskUpdateAPI)
	protectedAPI.DELETE("/admin/scheduled-tasks/{id}", admin.ScheduledTaskDeleteAPI)
	protectedAPI.POST("/admin/scheduled-tasks/{id}/run-now", admin.ScheduledTaskRunNowAPI)
	protectedAPI.POST("/admin/scheduled-tasks/{id}/enable", admin.ScheduledTaskEnableAPI)
	protectedAPI.POST("/admin/scheduled-tasks/{id}/disable", admin.ScheduledTaskDisableAPI)
	protectedAPI.GET("/admin/scheduler/status", admin.SchedulerStatusAPI)
	protectedAPI.POST("/admin/scheduler/restart", admin.SchedulerRestartAPI)
	
	// 管理后台静态资源 - 自定义 handler 处理 /admin/assets/*
	App.GET("/admin/assets/{filename}", func(c buffalo.Context) error {
		filename := c.Param("filename")
		filePath := "public/assets/" + filename
		http.ServeFile(c.Response(), c.Request(), filePath)
		return nil
	})
	
	// 管理后台其他静态文件（favicon等）
	App.GET("/admin/favicon.svg", func(c buffalo.Context) error {
		http.ServeFile(c.Response(), c.Request(), "public/favicon.svg")
		return nil
	})
	
	// 管理后台 SPA 入口 - 所有 /admin/* 页面路由都返回 index.html
	// Vue Router 在前端处理具体路由
	spaRoutes := []string{
		"/admin",
		"/admin/",
		"/admin/login",
		"/admin/dashboard",
		"/admin/posts",
		"/admin/posts/new",
		"/admin/posts/{id}",
		"/admin/posts/{id}/edit",
		"/admin/categories",
		"/admin/tags",
		"/admin/users",
		"/admin/users/new",
		"/admin/users/{id}",
		"/admin/users/{id}/edit",
		"/admin/media",
		"/admin/pages",
		"/admin/pages/new",
		"/admin/pages/{id}",
		"/admin/pages/{id}/edit",
		"/admin/themes",
		"/admin/settings",
		"/admin/comments",
		"/admin/scheduled-tasks",
	}
	for _, route := range spaRoutes {
		App.GET(route, serveAdminSPA)
	}
	
	// 前台路由
	App.GET("/", front.HomeHandler)
	App.GET("/blog/{slug}", front.BlogShowHandler)
	App.GET("/{slug}", front.PageShowHandler)

	// 前台评论API（不需要认证）
	App.GET("/api/posts/{post_id}/comments", front.FrontCommentsListAPI)
	App.POST("/api/comments", front.FrontCommentCreateAPI)
	
	// 全局搜索API（不需要认证）
	App.GET("/api/search", admin.GlobalSearchAPI)

	// Sitemap API（不需要认证，公开访问）
	// 注意：必须在静态文件服务之前定义，避免被静态文件拦截
	App.GET("/sitemap.xml", admin.SitemapHandler)
	protectedAPI.GET("/admin/sitemap/info", admin.SitemapInfoHandler)
	protectedAPI.POST("/admin/sitemap/rebuild", admin.RebuildSitemapHandler)

	// 静态文件服务
	App.ServeFiles("/", http.Dir("public"))
	
	// 主题静态文件服务（用于访问主题资源如缩略图、CSS、JS等）
	App.ServeFiles("/themes/", http.Dir("templates/themes"))
	
	// 初始化搜索服务（Windows 上可能卡住，暂时禁用）
	// log.Println("Initializing search service...")
	// searchService := services.GetSearchService()
	// if err := searchService.Init(); err != nil {
	// 	log.Printf("Failed to initialize search service: %v", err)
	// } else {
	// 	log.Println("Search service initialized successfully")
	// }
	log.Println("Search service disabled on Windows")

	// 初始化定时任务调度器
	log.Println("Initializing scheduler...")
	scheduler := services.GetScheduler()
	if err := scheduler.Init(models.DB); err != nil {
		log.Printf("Failed to initialize scheduler: %v", err)
	} else {
		log.Println("Scheduler initialized successfully")
	}
	
	log.Println("Application initialization completed!")
}
