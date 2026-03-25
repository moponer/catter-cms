package middleware

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SPAConfig SPA 中间件配置
type SPAConfig struct {
	// 静态文件目录路径
	StaticPath string
	// SPA 的基础路径（如 "/admin"）
	BasePath string
	// 是否启用（如果静态文件目录不存在则自动禁用）
	Enabled bool
}

// DefaultSPAConfig 创建默认 SPA 配置
func DefaultSPAConfig(staticPath, basePath string) SPAConfig {
	// 检查静态文件目录是否存在
	if _, err := os.Stat(staticPath); err != nil {
		return SPAConfig{
			StaticPath: staticPath,
			BasePath:   basePath,
			Enabled:    false,
		}
	}
	return SPAConfig{
		StaticPath: staticPath,
		BasePath:   basePath,
		Enabled:    true,
	}
}

// SPA SPA 单页应用中间件
// 处理 SPA 路由：当请求的文件不存在时，返回 index.html 让前端路由处理
// 这解决了 Vue/React 等 SPA 框架在直接访问子路由时显示空白的问题
func SPA(config SPAConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果未启用，直接跳过
		if !config.Enabled {
			c.Next()
			return
		}

		path := c.Request.URL.Path

		// 检查请求路径是否属于 SPA 管辖范围
		if !strings.HasPrefix(path, config.BasePath) && path != config.BasePath {
			c.Next()
			return
		}

		// 计算实际文件路径
		// 例如: /admin/themes -> /admin/dist/themes
		relativePath := path
		if strings.HasPrefix(path, config.BasePath) {
			relativePath = path[len(config.BasePath):]
		}
		// 确保有前导斜杠
		if !strings.HasPrefix(relativePath, "/") {
			relativePath = "/" + relativePath
		}

		filePath := filepath.Join(config.StaticPath, filepath.FromSlash(relativePath))

		// 检查文件是否存在
		info, err := os.Stat(filePath)
		if err == nil && !info.IsDir() {
			// 文件存在，直接让后续 handler 处理（Static 中间件会处理）
			c.Next()
			return
		}

		// 目录或文件不存在，检查是否是目录请求
		// 尝试查找 index.html
		indexPath := filepath.Join(config.StaticPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			// 返回 index.html 让前端路由处理
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.File(indexPath)
			c.Abort()
			return
		}

		// 没有 index.html，继续让其他 handler 处理
		c.Next()
	}
}

// SPAMiddleware 创建 SPA 中间件（兼容旧接口）
// staticPath: 静态文件根目录（如 "../admin/dist"）
// basePath: SPA 的 URL 基础路径（如 "/admin"）
func SPAMiddleware(staticPath, basePath string) gin.HandlerFunc {
	config := DefaultSPAConfig(staticPath, basePath)
	return SPA(config)
}
