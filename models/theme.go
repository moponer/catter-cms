package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Theme 定义主题模型
type Theme struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Version     string    `json:"version" db:"version"`
	Description string    `json:"description" db:"description"`
	Author      string    `json:"author" db:"author"`
	Thumbnail   string    `json:"thumbnail" db:"thumbnail"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Theme) TableName() string {
	return "themes"
}

// Validate 验证主题数据
func (t *Theme) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
		&validators.StringIsPresent{Field: t.Slug, Name: "Slug"},
		&validators.StringIsPresent{Field: t.Version, Name: "Version"},
	), nil
}

// Themes 主题集合
type Themes []Theme