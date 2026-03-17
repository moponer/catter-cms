package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Tag 定义标签模型
type Tag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=50"`
	Slug      string    `json:"slug" db:"slug" validate:"required,slug"`
	Posts     Posts     `json:"posts,omitempty" many_to_many:"post_tags"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// Validate 验证标签数据
func (t *Tag) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
		&validators.StringIsPresent{Field: t.Slug, Name: "Slug"},
	), nil
}

// Tags 标签集合
type Tags []Tag