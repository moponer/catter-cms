package admin

import (
	"catter-cms/actions/api"
	"catter-cms/models"
	"math/rand"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"golang.org/x/crypto/bcrypt"
)

// UserRequest 用户请求结构
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UsersListAPI 获取用户列表API
func UsersListAPI(c buffalo.Context) error {
	// 从查询字符串中获取参数
	page, _ := strconv.Atoi(c.Request().URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(c.Request().URL.Query().Get("page_size"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	users := &models.Users{}
	query := models.DB.Paginate(page, pageSize)
	query = query.Order("created_at DESC")

	if err := query.All(users); err != nil {
		return api.InternalError(c, "获取用户列表失败")
	}

	// 获取总数
	total, _ := models.DB.Count(&models.User{})

	return api.PagedSuccess(c, users, total, page, pageSize)
}

// UserDetailAPI 获取用户详情API
func UserDetailAPI(c buffalo.Context) error {
	userID := c.Param("id")
	user := &models.User{}

	if err := models.DB.Find(user, userID); err != nil {
		return api.NotFound(c, "用户不存在")
	}

	return api.Success(c, user)
}

// UserCreateAPI 创建用户API
func UserCreateAPI(c buffalo.Context) error {
	req := &UserRequest{}
	if err := c.Bind(req); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 验证必填字段
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return api.BadRequest(c, "邮箱、姓名和密码不能为空")
	}

	// 检查邮箱是否已存在
	existingUser := &models.User{}
	if err := models.DB.Where("email = ?", req.Email).First(existingUser); err == nil {
		return api.BadRequest(c, "邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return api.InternalError(c, "密码加密失败")
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}

	if user.Role == "" {
		user.Role = "editor"
	}

	if err := models.DB.Create(user); err != nil {
		return api.InternalError(c, "创建用户失败")
	}

	return api.SuccessWithMessage(c, "用户创建成功", user)
}

// UserUpdateAPI 更新用户API
func UserUpdateAPI(c buffalo.Context) error {
	userID := c.Param("id")
	user := &models.User{}

	if err := models.DB.Find(user, userID); err != nil {
		return api.NotFound(c, "用户不存在")
	}

	req := &UserRequest{}
	if err := c.Bind(req); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// 更新字段
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// 检查邮箱是否被其他用户使用
		existingUser := &models.User{}
		if err := models.DB.Where("email = ? AND id != ?", req.Email, user.ID).First(existingUser); err == nil {
			return api.BadRequest(c, "邮箱已被其他用户使用")
		}
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return api.InternalError(c, "密码加密失败")
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := models.DB.Update(user); err != nil {
		return api.InternalError(c, "更新用户失败")
	}

	return api.SuccessWithMessage(c, "用户更新成功", user)
}

// UserDeleteAPI 删除用户API
func UserDeleteAPI(c buffalo.Context) error {
	userID := c.Param("id")
	user := &models.User{}

	if err := models.DB.Find(user, userID); err != nil {
		return api.NotFound(c, "用户不存在")
	}

	// 不允许删除自己
	currentUser, ok := c.Value("current_user").(*models.User)
	if ok && currentUser.ID == user.ID {
		return api.BadRequest(c, "不能删除自己的账号")
	}

	if err := models.DB.Destroy(user); err != nil {
		return api.InternalError(c, "删除用户失败")
	}

	return api.SuccessWithMessage(c, "用户删除成功", nil)
}

// UserResetPasswordAPI 重置用户密码API
func UserResetPasswordAPI(c buffalo.Context) error {
	userID := c.Param("id")
	user := &models.User{}

	if err := models.DB.Find(user, userID); err != nil {
		return api.NotFound(c, "用户不存在")
	}

	newPassword := generateRandomPassword()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return api.InternalError(c, "密码加密失败")
	}

	user.PasswordHash = string(hashedPassword)

	if err := models.DB.Update(user); err != nil {
		return api.InternalError(c, "密码重置失败")
	}

	return api.Success(c, map[string]string{
		"message":      "密码重置成功",
		"new_password": newPassword,
	})
}

// generateRandomPassword 生成随机密码
func generateRandomPassword() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
	password := ""
	for i := 0; i < 12; i++ {
		password += string(chars[rand.Intn(len(chars))])
	}
	return password
}
