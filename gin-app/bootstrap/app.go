package bootstrap

import (
	"gin-app/internal/config"
	"gin-app/internal/database"
	"gin-app/internal/theme"
	"gin-app/pkg/logger"

	"github.com/gin-gonic/gin"
)

// App 应用实例
type App struct {
	Config         *config.Config
	Logger         *logger.Logger
	DB             *database.Database
	ThemeManager   *theme.Manager
	Engine         *gin.Engine
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{}
}

// Bootstrap 初始化应用
func (a *App) Bootstrap() error {
	var err error

	// 1. 初始化配置
	a.Config, err = config.Load()
	if err != nil {
		return err
	}

	// 2. 初始化日志
	a.Logger, err = logger.New(a.Config.Log.Level, a.Config.Log.File)
	if err != nil {
		return err
	}

	// 3. 初始化数据库
	a.DB, err = database.New(a.Config.Database)
	if err != nil {
		a.Logger.Error("数据库初始化失败", "error", err)
		return err
	}

	// 4. 初始化主题管理器
	a.ThemeManager = theme.NewManager(
		a.Config.Theme.TemplateDir,
		a.Config.Theme.ConfigDir,
	)

	// 加载所有主题
	if err := a.ThemeManager.LoadThemes(); err != nil {
		a.Logger.Error("加载主题失败", "error", err)
		// 不中断启动，使用默认主题
	}

	// 5. 初始化Gin引擎
	gin.SetMode(a.Config.Server.Mode)
	a.Engine = gin.New()

	return nil
}

// Run 运行应用
func (a *App) Run() error {
	return a.Engine.Run(":" + a.Config.Server.Port)
}

// Close 关闭应用
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
