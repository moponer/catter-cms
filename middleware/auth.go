package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/golang-jwt/jwt/v5"
	"catter-cms/config"
	"catter-cms/errors"
	"catter-cms/models"
)

// RequireAuth 要求用户已认证
func RequireAuth(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		userID, ok := c.Session().Get("user_id").(string)
		if !ok || userID == "" {
			return c.Redirect(http.StatusSeeOther, "/admin/login")
		}
		
		// 获取用户信息
		user := &models.User{}
		id, _ := strconv.Atoi(userID)
		if err := models.DB.Find(user, id); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/login")
		}
		
		// 将用户信息存储到上下文
		c.Set("current_user", user)
		return next(c)
	}
}

// RequireRole 要求用户具有指定角色
func RequireRole(roles ...string) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			user, ok := c.Value("current_user").(*models.User)
			if !ok || user == nil {
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}
			
			// 检查用户角色
			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}
			
			if !hasRole {
				return c.Error(http.StatusForbidden, errors.NewForbiddenError("权限不足"))
			}
			
			return next(c)
		}
	}
}

// RequireOwnership 要求用户拥有资源
func RequireOwnership(model interface{}) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			user, ok := c.Value("current_user").(*models.User)
			if !ok || user == nil {
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}
			
			// 检查资源所有权
			if !checkOwnership(user, model) {
				return c.Error(http.StatusForbidden, errors.NewForbiddenError("无权访问此资源"))
			}
			
			return next(c)
		}
	}
}

// checkOwnership 检查资源所有权
func checkOwnership(user *models.User, model interface{}) bool {
	switch m := model.(type) {
	case *models.Post:
		return m.UserID == user.ID
	case *models.User:
		return m.ID == user.ID
	default:
		return false
	}
}

// APIAuth API认证中间件
// 用于验证API请求的JWT令牌
func APIAuth(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// 获取Authorization头
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.Error(http.StatusUnauthorized, errors.NewUnauthorizedError("未授权"))
		}
		
		// 解析Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Error(http.StatusUnauthorized, errors.NewUnauthorizedError("无效的授权令牌"))
		}
		
		// 验证JWT token
		user, err := validateToken(tokenString)
		if err != nil {
			return c.Error(http.StatusUnauthorized, errors.NewUnauthorizedError("无效的授权令牌"))
		}
		
		// 将用户信息存储到上下文
		c.Set("current_user", user)
		return next(c)
	}
}

// validateToken 验证JWT token并返回用户信息
// 参数:
//   - tokenString: JWT令牌字符串
// 返回:
//   - *models.User: 用户对象
//   - error: 错误信息
func validateToken(tokenString string) (*models.User, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewBadRequestError("无效的签名方法", nil)
		}
		// 返回密钥 - 使用配置文件中的密钥
		return []byte(config.GetSessionSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token有效性
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 查找用户
		user := &models.User{}
		if err := models.DB.Find(user, claims.UserID); err != nil {
			return nil, errors.NewNotFoundError("用户", claims.UserID)
		}
		return user, nil
	}

	return nil, errors.NewBadRequestError("无效的令牌", nil)
}

// Claims JWT声明结构
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
