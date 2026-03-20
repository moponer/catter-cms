package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 文章模型
type Post struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Title           string     `json:"title" gorm:"not null;size:500"`
	Slug            string     `json:"slug" gorm:"uniqueIndex;not null;size:255"`
	Content         string     `json:"content" gorm:"type:text"`
	Excerpt         string     `json:"excerpt" gorm:"type:text"`
	Status          string     `json:"status" gorm:"default:'draft';size:50"`
	Views           int        `json:"views" gorm:"default:0"`
	MetaTitle       string     `json:"meta_title" gorm:"size:500"`
	MetaDescription string     `json:"meta_description" gorm:"type:text"`
	CategoryID      *uint      `json:"category_id"`
	Category        *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	AuthorID        uint       `json:"author_id"`
	Author          *User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Tags            []Tag      `json:"tags,omitempty" gorm:"many2many:post_tags;"`
	PublishedAt     *time.Time `json:"published_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate 创建前钩子
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Status == "" {
		p.Status = "draft"
	}
	return nil
}

// IsPublished 检查是否已发布
func (p *Post) IsPublished() bool {
	return p.Status == "published"
}
