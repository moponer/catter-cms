package models

import (
	"time"

	"gorm.io/gorm"
)

// Media 媒体模型
type Media struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Filename    string    `json:"filename" gorm:"not null;size:255"`
	Path        string    `json:"path" gorm:"not null;size:500"`
	URL         string    `json:"url" gorm:"size:500"`
	FileType    string    `json:"file_type" gorm:"size:50"`
	FileSize    int64     `json:"file_size"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Alt         string    `json:"alt" gorm:"size:255"`
	Title       string    `json:"title" gorm:"size:255"`
	Description string    `json:"description" gorm:"type:text"`
	UploadedBy  uint      `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Media) TableName() string {
	return "media"
}

// BeforeCreate 创建前钩子
func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.FileType == "" {
		m.FileType = getFileType(m.Filename)
	}
	return nil
}

// getFileType 根据文件名获取文件类型
func getFileType(filename string) string {
	ext := ""
	if len(filename) > 4 {
		ext = filename[len(filename)-4:]
	}
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp":
		return "image"
	case ".mp4", ".webm":
		return "video"
	case ".mp3", ".wav", ".ogg":
		return "audio"
	case ".pdf":
		return "pdf"
	default:
		return "other"
	}
}

// IsImage 检查是否为图片
func (m *Media) IsImage() bool {
	return m.FileType == "image"
}
