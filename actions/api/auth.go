package api

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"catter-cms/config"
	"catter-cms/models"
)

// Claims JWT 声明结构
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// LoginHandler 处理登录请求
func LoginHandler(c buffalo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "请求参数错误")
	}

	// 验证必填字段
	if req.Email == "" || req.Password == "" {
		return BadRequest(c, "邮箱和密码不能为空")
	}

	// 查找用户
	user := &models.User{}
	query := models.DB.Where("email = ?", req.Email)
	if err := query.First(user); err != nil {
		return Unauthorized(c, "邮箱或密码错误")
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return Unauthorized(c, "邮箱或密码错误")
	}

	// 生成 JWT token
	token, err := generateToken(user)
	if err != nil {
		return InternalError(c, "登录失败,请稍后重试")
	}

	// 返回统一格式的响应
	return SuccessWithMessage(c, "登录成功", LoginResponse{
		Token: token,
		User:  *user,
	})
}

// MeHandler 获取当前用户信息
func MeHandler(c buffalo.Context) error {
	user, ok := c.Value("current_user").(*models.User)
	if !ok || user == nil {
		return c.Error(http.StatusUnauthorized, nil)
	}

	return JSON(c, http.StatusOK, user)
}

// generateToken 生成 JWT token
func generateToken(user *models.User) (string, error) {
	// 创建声明
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 这里可以设置过期时间
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token - 使用配置文件中的密钥
	secretKey := config.GetSessionSecret()
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
