package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	PostID       uint       `json:"post_id" gorm:"not null;index"`
	Post         *Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	UserID       *uint      `json:"user_id"`
	User         *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ParentID     *uint      `json:"parent_id"`
	Parent       *Comment   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies      []Comment  `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	AuthorName   string     `json:"author_name" gorm:"size:255"`
	AuthorEmail  string     `json:"author_email" gorm:"size:255"`
	AuthorURL    string     `json:"author_url" gorm:"size:500"`
	Content      string     `json:"content" gorm:"type:text;not null"`
	Status       string     `json:"status" gorm:"default:'pending';size:50"` // pending, approved, spam, trash
	IPAddress    string     `json:"ip_address" gorm:"size:100"`
	UserAgent    string     `json:"user_agent" gorm:"size:500"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// IsApproved 检查是否已批准
func (c *Comment) IsApproved() bool {
	return c.Status == "approved"
}
