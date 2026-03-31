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

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
}

// UpdateProfileHandler 更新当前用户个人资料
func UpdateProfileHandler(c *gin.Context) {
	currentUser, exists := c.Get("current_user")
	if !exists || currentUser == nil {
		Unauthorized(c, "未登录")
		return
	}
	cu := currentUser.(*models.User)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 验证名称
	if req.Name == "" {
		BadRequest(c, "昵称不能为空")
		return
	}

	// 如果修改了邮箱，检查是否已存在
	if req.Email != "" && req.Email != cu.Email {
		var count int64
		models.DB.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, cu.ID).Count(&count)
		if count > 0 {
			BadRequest(c, "邮箱已被使用")
			return
		}
	}

	updates := map[string]interface{}{
		"name":   req.Name,
		"avatar": req.Avatar,
		"bio":    req.Bio,
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	if err := models.DB.Model(cu).Updates(updates).Error; err != nil {
		InternalError(c, "更新失败")
		return
	}

	// 重新查询最新数据
	var updatedUser models.User
	models.DB.First(&updatedUser, cu.ID)
	SuccessWithMessage(c, "个人资料更新成功", updatedUser)
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePasswordHandler 修改当前用户密码
func ChangePasswordHandler(c *gin.Context) {
	currentUser, exists := c.Get("current_user")
	if !exists || currentUser == nil {
		Unauthorized(c, "未登录")
		return
	}
	cu := currentUser.(*models.User)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误，新密码至少6位")
		return
	}

	// 验证旧密码
	if !cu.CheckPassword(req.OldPassword) {
		BadRequest(c, "当前密码不正确")
		return
	}

	// 新旧密码不能相同
	if req.OldPassword == req.NewPassword {
		BadRequest(c, "新密码不能与当前密码相同")
		return
	}

	// 设置新密码
	if err := cu.SetPassword(req.NewPassword); err != nil {
		InternalError(c, "密码加密失败")
		return
	}

	if err := models.DB.Model(cu).Update("password_hash", cu.PasswordHash).Error; err != nil {
		InternalError(c, "密码更新失败")
		return
	}

	SuccessWithMessage(c, "密码修改成功", nil)
}
