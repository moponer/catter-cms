package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// ContextKeyAIToken AI Token 上下文键
const ContextKeyAIToken = "ai_token"

// AIAuth AI Key 认证中间件
// 验证 Authorization: Bearer cat_xxx 格式的 catter-cms-key
// 与现有 JWT 认证完全隔离
func AIAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少 Authorization 头",
			})
			c.Abort()
			return
		}

		// 提取 Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization 格式错误，应为 Bearer <key>",
			})
			c.Abort()
			return
		}

		// 必须是 cat_ 前缀
		if !strings.HasPrefix(tokenString, "cat_") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的 API Key 格式",
			})
			c.Abort()
			return
		}

		// 查找 token
		var aiToken models.AIToken
		if err := models.DB.Where("token = ? AND is_active = ?", tokenString, true).First(&aiToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的 API Key",
			})
			c.Abort()
			return
		}

		// 检查过期
		if aiToken.IsExpired() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "API Key 已过期",
			})
			c.Abort()
			return
		}

		// 更新最后使用时间（异步，不影响响应速度）
		now := time.Now()
		go models.DB.Model(&aiToken).Update("last_used_at", now)

		// 将 token 信息存入上下文
		c.Set(ContextKeyAIToken, &aiToken)
		c.Next()
	}
}

// GetAIToken 从上下文获取 AI Token
func GetAIToken(c *gin.Context) (*models.AIToken, bool) {
	token, exists := c.Get(ContextKeyAIToken)
	if !exists || token == nil {
		return nil, false
	}
	aiToken, ok := token.(*models.AIToken)
	return aiToken, ok
}

// RequireAIPermission 检查 AI Token 是否拥有指定权限的中间件
func RequireAIPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		aiToken, ok := GetAIToken(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未认证",
			})
			c.Abort()
			return
		}

		if !aiToken.HasPermission(permission) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足，需要 " + permission + " 权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
