package repositories

import (
	"catter-cms/errors"
	"catter-cms/models"
)

// UserRepository 用户仓储
type UserRepository struct {
	db *pop.Connection
}

// NewUserRepository 创建新的用户仓储
func NewUserRepository(db *pop.Connection) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Find(user, id); err != nil {
		return nil, errors.NewNotFoundError("用户", id)
	}
	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("email = ?", email).First(user); err != nil {
		return nil, errors.NewInternalServerError("获取用户失败")
	}
	return user, nil
}

// Create 创建新用户
func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.Create(user); err != nil {
		return errors.NewInternalServerError("创建用户失败")
	}
	return nil
}

// Update 更新用户
func (r *UserRepository) Update(user *models.User) error {
	if err := r.db.Update(user); err != nil {
		return errors.NewInternalServerError("更新用户失败")
	}
	return nil
}

// Delete 删除用户
func (r *UserRepository) Delete(id int) error {
	user := &models.User{}
	if err := r.db.Find(user, id); err != nil {
		return errors.NewNotFoundError("用户", id)
	}
	
	if err := r.db.Destroy(user); err != nil {
		return errors.NewInternalServerError("删除用户失败")
	}
	return nil
}
