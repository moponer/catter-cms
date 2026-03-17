package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Media 定义媒体模型
type Media struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Filename    string    `json:"filename" db:"filename"`
	Path        string    `json:"path" db:"path"`
	Size        int64     `json:"size" db:"size"`
	MimeType    string    `json:"mime_type" db:"mime_type"`
	Thumbnail   string    `json:"thumbnail,omitempty" db:"thumbnail"`
	Description string    `json:"description,omitempty" db:"description"`
	UserID      int       `json:"user_id" db:"user_id"` // 修改为int类型，与User.ID一致
	User        User      `json:"user,omitempty" belongs_to:"users"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TableName 指定表名
func (Media) TableName() string {
	return "media"
}

// Validate 验证媒体数据
func (m *Media) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Filename, Name: "Filename"},
		&validators.StringIsPresent{Field: m.Path, Name: "Path"},
		&validators.StringIsPresent{Field: m.MimeType, Name: "MimeType"},
	), nil
}

// Media 媒体集合
type MediaItems []Media