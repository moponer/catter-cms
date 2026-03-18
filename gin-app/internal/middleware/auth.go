package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/models"
)

// ContextKeys 定义上下文键
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUser    = "current_user"
	ContextKeyToken   = "token"
)

// Claims JWT声明结构
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RequireAuth 要求用户已认证
// 如果用户未认证，重定向到登录页面
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Cookie获取session token
		tokenString, err := c.Cookie("session_token")
		if err != nil || tokenString == "" {
			// 尝试从Header获取
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		// 验证token
		user, err := ValidateToken(tokenString)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUser, user)
		c.Set(ContextKeyToken, tokenString)

		c.Next()
	}
}

// RequireRole 要求用户具有指定角色
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(ContextKeyUser)
		if !exists || user == nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		currentUser, ok := user.(*models.User)
		if !ok || currentUser == nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		// 检查用户角色
		hasRole := false
		for _, role := range roles {
			if currentUser.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIAuth API认证中间件
// 用于验证API请求的JWT令牌
func APIAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权",
			})
			c.Abort()
			return
		}

		// 解析Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" || tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的授权令牌格式",
			})
			c.Abort()
			return
		}

		// 验证JWT token
		user, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的授权令牌",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set(ContextKeyUserID, user.ID)
		c.Set(ContextKeyUser, user)
		c.Set(ContextKeyToken, tokenString)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件
// 如果用户提供了有效的token，则设置用户信息；否则继续执行
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从多个来源获取token
		tokenString := ""

		// 1. 从Cookie获取
		if cookie, err := c.Cookie("session_token"); err == nil && cookie != "" {
			tokenString = cookie
		}

		// 2. 从Header获取
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 如果有token，验证并设置用户信息
		if tokenString != "" {
			if user, err := ValidateToken(tokenString); err == nil {
				c.Set(ContextKeyUserID, user.ID)
				c.Set(ContextKeyUser, user)
			}
		}

		c.Next()
	}
}

// ValidateToken 验证JWT token并返回用户信息
func ValidateToken(tokenString string) (*models.User, error) {
	cfg := config.Get()

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token有效性
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 查找用户
		user := &models.User{}
		db := models.DB
		if err := db.First(user, claims.UserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, jwt.ErrSignatureInvalid
			}
			return nil, err
		}
		return user, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GenerateToken 生成JWT token
func GenerateToken(user *models.User) (string, error) {
	cfg := config.Get()

	// 创建token
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.JWT.Expire))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "catter-cms",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GetCurrentUser 获取当前用户
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get(ContextKeyUser)
	if !exists || user == nil {
		return nil, false
	}
	currentUser, ok := user.(*models.User)
	return currentUser, ok
}

// GetUserID 获取当前用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists || userID == nil {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}
