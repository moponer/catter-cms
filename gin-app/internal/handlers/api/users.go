package api

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"catter-cms-gin/internal/models"
)

// UserRequest 用户请求结构
type UserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// ListUsers 获取用户列表
func ListUsers(c *gin.Context) {
	var users []models.User
	query := models.DB

	// 角色筛选
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}

	// 搜索（转义 LIKE 通配符防止注入）
	if search := c.Query("search"); search != "" {
		escapedSearch := escapeLikeWildcards(search)
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+escapedSearch+"%", "%"+escapedSearch+"%")
	}

	// 获取总数
	var total int64
	query.Model(&models.User{}).Count(&total)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		InternalError(c, "获取用户列表失败")
		return
	}

	SuccessResponse(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的用户ID")
		return
	}

	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	SuccessResponse(c, user)
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 验证必填字段
	if req.Email == "" || req.Name == "" {
		BadRequest(c, "邮箱和名称不能为空")
		return
	}

	// 验证邮箱格式
	if req.Email != "" {
		var count int64
		models.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
		if count > 0 {
			BadRequest(c, "邮箱已被使用")
			return
		}
	}

	// 验证角色
	if req.Role == "" {
		req.Role = "editor"
	}
	if req.Role != "admin" && req.Role != "editor" {
		BadRequest(c, "无效的角色")
		return
	}

	// 创建用户
	user := models.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  req.Role,
	}

	// 如果提供了密码，则加密
	if req.Password != "" {
		if err := user.SetPassword(req.Password); err != nil {
			InternalError(c, "密码加密失败")
			return
		}
	} else {
		// 生成随机密码
		if err := user.SetPassword(generateRandomPassword(12)); err != nil {
			InternalError(c, "密码生成失败")
			return
		}
	}

	if err := models.DB.Create(&user).Error; err != nil {
		InternalError(c, "创建用户失败")
		return
	}

	SuccessWithMessage(c, "用户创建成功", user)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的用户ID")
		return
	}

	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	// 获取当前用户，检查权限
	currentUser, exists := c.Get("current_user")
	if !exists || currentUser == nil {
		Unauthorized(c, "未登录")
		return
	}
	cu := currentUser.(*models.User)

	// 非管理员不能修改管理员
	if cu.Role != "admin" && user.Role == "admin" {
		Forbidden(c, "权限不足")
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	// 检查邮箱是否已存在（排除自己）
	if req.Email != "" && req.Email != user.Email {
		var count int64
		models.DB.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, id).Count(&count)
		if count > 0 {
			BadRequest(c, "邮箱已被使用")
			return
		}
	}

	// 检查角色是否有效
	if req.Role != "" && req.Role != "admin" && req.Role != "editor" {
		BadRequest(c, "无效的角色")
		return
	}

	// 非管理员不能将其他用户设为管理员
	if cu.Role != "admin" && req.Role == "admin" {
		Forbidden(c, "权限不足")
		return
	}

	updates := map[string]interface{}{
		"email": req.Email,
		"name":  req.Name,
		"role":  req.Role,
	}

	// 如果提供了新密码
	if req.Password != "" {
		tmpUser := &models.User{}
		if err := tmpUser.SetPassword(req.Password); err != nil {
			InternalError(c, "密码加密失败")
			return
		}
		updates["password_hash"] = tmpUser.PasswordHash
	}

	if err := models.DB.Model(&user).Updates(updates).Error; err != nil {
		InternalError(c, "更新用户失败")
		return
	}

	SuccessWithMessage(c, "用户更新成功", user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的用户ID")
		return
	}

	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		NotFound(c, "用户不存在")
		return
	}

	// 获取当前用户，检查权限
	currentUser, exists := c.Get("current_user")
	if !exists || currentUser == nil {
		Unauthorized(c, "未登录")
		return
	}
	cu := currentUser.(*models.User)

	// 非管理员不能删除管理员
	if cu.Role != "admin" && user.Role == "admin" {
		Forbidden(c, "权限不足")
		return
	}

	// 不能删除自己
	if cu.ID == user.ID {
		BadRequest(c, "不能删除当前登录用户")
		return
	}

	if err := models.DB.Delete(&user).Error; err != nil {
		InternalError(c, "删除用户失败")
		return
	}

	SuccessWithMessage(c, "用户删除成功", nil)
}

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[r.Intn(len(charset))]
	}
	return string(password)
}
