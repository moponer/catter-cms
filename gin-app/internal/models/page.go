package models

import (
	"time"

	"gorm.io/gorm"
)

// Page 页面模型
type Page struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Title          string    `json:"title" gorm:"not null;size:500"`
	Slug           string    `json:"slug" gorm:"uniqueIndex;not null;size:255"`
	Content        string    `json:"content" gorm:"type:text"`
	Excerpt        string    `json:"excerpt" gorm:"type:text"`
	Status         string    `json:"status" gorm:"default:'draft';size:50"`
	Template       string    `json:"template" gorm:"size:100"`
	TemplateData   string    `json:"template_data" gorm:"type:text"` // 模板专属数据（JSON格式）
	MetaTitle      string    `json:"meta_title" gorm:"size:500"`
	MetaDescription string   `json:"meta_description" gorm:"type:text"`
	// SEO 扩展字段
	OgImage      string    `json:"og_image" gorm:"size:1000"` // Open Graph 图片
	CanonicalURL string    `json:"canonical_url" gorm:"size:1000"` // 规范 URL
	Robots       string    `json:"robots" gorm:"size:100"` // 搜索引擎索引指令 (noindex, nofollow 等)
	Views        int       `json:"views" gorm:"default:0"`
	SortOrder    int       `json:"sort_order" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Page) TableName() string {
	return "pages"
}

// BeforeCreate 创建前钩子
func (p *Page) BeforeCreate(tx *gorm.DB) error {
	if p.Status == "" {
		p.Status = "draft"
	}
	return nil
}

// IsPublished 检查是否已发布
func (p *Page) IsPublished() bool {
	return p.Status == "published"
}
