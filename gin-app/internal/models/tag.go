package models

import (
	"time"
)

// Tag 标签模型
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null;size:255"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;not null;size:255"`
	Posts     []Post    `json:"posts,omitempty" gorm:"many2many:post_tags;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// PostTag 文章标签关联表
type PostTag struct {
	PostID uint `json:"post_id" gorm:"primaryKey"`
	TagID  uint `json:"tag_id" gorm:"primaryKey"`
}

// TableName 指定表名
func (PostTag) TableName() string {
	return "post_tags"
}
