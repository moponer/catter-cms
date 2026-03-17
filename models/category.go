package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Category 定义分类模型
type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Slug        string    `json:"slug" db:"slug" validate:"required,slug"`
	Description string    `json:"description,omitempty" db:"description" validate:"max=500"`
	ParentID    *int      `json:"parent_id,omitempty" db:"parent_id"`
	Parent      *Category  `json:"parent,omitempty" belongs_to:"categories"`
	Children    Categories `json:"children,omitempty" has_many:"categories" fk:"parent_id"`
	Posts       Posts      `json:"posts,omitempty" has_many:"posts" fk:"category_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Validate 验证分类数据
func (c *Category) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
		&validators.StringIsPresent{Field: c.Slug, Name: "Slug"},
	), nil
}

// Categories 分类集合
type Categories []Category