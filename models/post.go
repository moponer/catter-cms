package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Post 定义文章模型
type Post struct {
	ID            int        `json:"id" db:"id"`
	Title         string     `json:"title" db:"title" validate:"required,min=1,max=255"`
	Slug          string     `json:"slug" db:"slug" validate:"required,slug"`
	Content       string     `json:"content" db:"content" validate:"required"`
	Excerpt       string     `json:"excerpt,omitempty" db:"excerpt" validate:"max=500"`
	FeaturedImage string     `json:"featured_image,omitempty" db:"featured_image" validate:"omitempty,url"`
	Status        string     `json:"status" db:"status" validate:"required,oneof=draft published archived"`
	UserID        int        `json:"user_id" db:"user_id" validate:"required"`
	User          User       `json:"user,omitempty" belongs_to:"users"`
	CategoryID    *int       `json:"category_id,omitempty" db:"category_id"`
	Category      *Category  `json:"category,omitempty" belongs_to:"categories"`
	Tags          Tags       `json:"tags,omitempty" many_to_many:"post_tags"`
	Views         int        `json:"views" db:"views" validate:"min=0"`
	PublishedAt   *time.Time `json:"published_at,omitempty" db:"published_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// Validate 验证文章数据
func (p *Post) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Slug, Name: "Slug"},
		&validators.StringIsPresent{Field: p.Content, Name: "Content"},
		&validators.StringInclusion{Field: p.Status, Name: "Status", List: []string{"draft", "published", "archived"}},
	), nil
}

// ReadingTime 计算阅读时间（按每分钟300字计算）
func (p *Post) ReadingTime() int {
	if len(p.Content) == 0 {
		return 1
	}
	// 中文按字符数，英文按单词数，这里简化处理
	minutes := len(p.Content) / 300
	if minutes < 1 {
		return 1
	}
	return minutes
}

// Posts 文章集合
type Posts []Post