package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Page 定义页面模型（用于静态页面如关于我们、联系我们等）
type Page struct {
	ID            int        `json:"id" db:"id"`
	Title         string     `json:"title" db:"title"`
	Slug          string     `json:"slug" db:"slug"`
	Content       string     `json:"content" db:"content"`
	Excerpt       string     `json:"excerpt,omitempty" db:"excerpt"`
	FeaturedImage string     `json:"featured_image,omitempty" db:"featured_image"`
	Template      string     `json:"template,omitempty" db:"template"`
	Status        string     `json:"status" db:"status"`
	UserID        int        `json:"user_id" db:"user_id"`
	User          User       `json:"user,omitempty" belongs_to:"users"`
	ParentID      *int       `json:"parent_id,omitempty" db:"parent_id"`
	Parent        *Page      `json:"parent,omitempty" belongs_to:"pages"`
	Children      Pages      `json:"children,omitempty" has_many:"pages"`
	Order         int        `json:"order" db:"order"`
	MetaTitle     string     `json:"meta_title,omitempty" db:"meta_title"`
	MetaDescription string   `json:"meta_description,omitempty" db:"meta_description"`
	PublishedAt   *time.Time `json:"published_at,omitempty" db:"published_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Page) TableName() string {
	return "pages"
}

// Validate 验证页面数据
func (p *Page) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Slug, Name: "Slug"},
		&validators.StringIsPresent{Field: p.Content, Name: "Content"},
		&validators.StringInclusion{Field: p.Status, Name: "Status", List: []string{"draft", "published", "archived"}},
	), nil
}

// Pages 页面集合
type Pages []Page
