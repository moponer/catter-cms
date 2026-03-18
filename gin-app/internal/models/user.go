package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Name         string    `json:"name" gorm:"not null;size:255"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	Role         string    `json:"role" gorm:"default:'editor';size:50"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = "editor"
	}
	return nil
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin 检查是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// CreateDefaultAdmin 创建默认管理员用户
func CreateDefaultAdmin() error {
	var count int64
	DB.Model(&User{}).Count(&count)

	// 如果已存在用户，则不创建
	if count > 0 {
		return nil
	}

	// 创建默认管理员
	admin := &User{
		Email: "admin@example.com",
		Name:  "Administrator",
		Role:  "admin",
	}

	if err := admin.SetPassword("admin123"); err != nil {
		return err
	}

	if err := DB.Create(admin).Error; err != nil {
		return err
	}

	log.Println("Default admin user created: admin@example.com / admin123")
	return nil
}
