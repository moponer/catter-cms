package services

import (
	"catter-cms/errors"
	"catter-cms/models"
	"github.com/gobuffalo/pop/v6"
)

// UserService 用户服务
type UserService struct {
	db *pop.Connection
}

// NewUserService 创建新的用户服务
func NewUserService(db *pop.Connection) *UserService {
	return &UserService{db: db}
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	if err := s.db.Find(user, id); err != nil {
		return nil, errors.NewNotFoundError("用户", id)
	}
	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := s.db.Where("email = ?", email).First(user); err != nil {
		return nil, errors.NewInternalServerError("获取用户失败")
	}
	return user, nil
}

// Create 创建新用户
func (s *UserService) Create(user *models.User) (*models.User, error) {
	if err := s.db.Create(user); err != nil {
		return nil, errors.NewInternalServerError("创建用户失败")
	}
	return user, nil
}

// Update 更新用户
func (s *UserService) Update(user *models.User) error {
	if err := s.db.Update(user); err != nil {
		return errors.NewInternalServerError("更新用户失败")
	}
	return nil
}

// Delete 删除用户
func (s *UserService) Delete(id int) error {
	user := &models.User{}
	if err := s.db.Find(user, id); err != nil {
		return errors.NewNotFoundError("用户", id)
	}
	
	if err := s.db.Destroy(user); err != nil {
		return errors.NewInternalServerError("删除用户失败")
	}
	
	return nil
}

// ValidatePassword 验证密码
func (s *UserService) ValidatePassword(user *models.User, password string) error {
	if !user.CheckPassword(password) {
		return errors.NewBadRequestError("邮箱或密码错误", nil)
	}
	return nil
}
