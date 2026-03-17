package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Setting 定义系统设置模型
type Setting struct {
	ID        int       `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	Type      string    `json:"type" db:"type"`
	Category  string    `json:"category" db:"category"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Setting) TableName() string {
	return "settings"
}

// Validate 验证设置数据
func (s *Setting) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Key, Name: "Key"},
		&validators.StringIsPresent{Field: s.Value, Name: "Value"},
	), nil
}

// Settings 设置集合
type Settings []Setting