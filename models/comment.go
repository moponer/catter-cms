package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Comment 评论模型
type Comment struct {
	ID          string     `db:"id" json:"id"`
	PostID      string     `db:"post_id" json:"post_id"`
	UserID      *int       `db:"user_id" json:"user_id"`
	ParentID    *string    `db:"parent_id" json:"parent_id"`
	AuthorName  string     `db:"author_name" json:"author_name"`
	AuthorEmail string     `db:"author_email" json:"author_email"`
	AuthorURL   string     `db:"author_url" json:"author_url"`
	Content     string     `db:"content" json:"content"`
	Status      string     `db:"status" json:"status"` // pending, approved, spam, trash
	IPAddress   string     `db:"ip_address" json:"ip_address"`
	UserAgent   string     `db:"user_agent" json:"user_agent"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	
	// 关联字段
	User    *User     `json:"user,omitempty" belongs_to:"users" fk_id:"user_id"`
	Post    *Post     `json:"post,omitempty" belongs_to:"posts" fk_id:"post_id"`
	Replies Comments  `json:"replies,omitempty" has_many:"comments" fk_id:"parent_id" order_by:"created_at asc"`
}

// Comments 评论集合
type Comments []*Comment

// Validate 验证评论数据
func (c *Comment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Content, Name: "Content"},
		&validators.StringLengthInRange{Field: c.AuthorName, Name: "AuthorName", Min: 1, Max: 255},
		&validators.EmailLike{Field: c.AuthorEmail, Name: "AuthorEmail"},
		&validators.FuncValidator{
			Name:    "Status",
			Field:   c.Status,
			Message: "无效的状态值",
			Fn: func() bool {
				validStatuses := map[string]bool{
					"pending": true,
					"approved": true,
					"spam": true,
					"trash": true,
				}
				return validStatuses[c.Status]
			},
		},
	), nil
}

// TableName 返回表名
func (c Comment) TableName() string {
	return "comments"
}
