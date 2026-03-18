package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/models"
)

// Claims JWT 声明结构
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string       `json:"token"`
	User  models.User  `json:"user"`
}

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "邮箱和密码不能为空")
		return
	}

	// 查找用户
	user := &models.User{}
	if err := models.DB.Where("email = ?", req.Email).First(user).Error; err != nil {
		Unauthorized(c, "邮箱或密码错误")
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		Unauthorized(c, "邮箱或密码错误")
		return
	}

	// 生成 JWT token
	token, err := generateToken(user)
	if err != nil {
		InternalError(c, "登录失败，请稍后重试")
		return
	}

	// 返回响应
	SuccessWithMessage(c, "登录成功", LoginResponse{
		Token: token,
		User:  *user,
	})
}

// generateToken 生成 JWT token
func generateToken(user *models.User) (string, error) {
	cfg := config.Get()

	// 创建声明
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.Expire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// MeHandler 获取当前用户信息
func MeHandler(c *gin.Context) {
	user, exists := c.Get("current_user")
	if !exists || user == nil {
		Unauthorized(c, "未登录")
		return
	}

	SuccessResponse(c, user)
}
