package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"golang.org/x/crypto/bcrypt"
)

// User 定义用户模型
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	PasswordHash string    `json:"-" db:"password_hash"` // 不返回给前端
	Role         string    `json:"role" db:"role"` // admin / editor
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Validate 验证用户数据
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		&validators.StringInclusion{Field: u.Role, Name: "Role", List: []string{"admin", "editor"}},
	), nil
}

// BeforeSave 保存前加密密码
func (u *User) BeforeSave(tx *pop.Connection) error {
	// 这里暂时不处理密码加密，因为我们直接使用password_hash字段
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Users 用户集合
type Users []User