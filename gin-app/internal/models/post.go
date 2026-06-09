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
	Cover           string     `json:"cover" gorm:"size:1000"` // 封面图
	Status          string     `json:"status" gorm:"default:'draft';size:50"`
	Views           int        `json:"views" gorm:"default:0"`
	MetaTitle       string     `json:"meta_title" gorm:"size:500"`
	MetaDescription string     `json:"meta_description" gorm:"type:text"`
	// SEO 扩展字段
	OgImage         string     `json:"og_image" gorm:"size:1000"` // Open Graph 图片
	CanonicalURL    string     `json:"canonical_url" gorm:"size:1000"` // 规范 URL
	Robots          string     `json:"robots" gorm:"size:100"` // 搜索引擎索引指令 (noindex, nofollow 等)
	CategoryID      *uint      `json:"category_id"`
	Category        *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	AuthorID        uint       `json:"author_id"`
	Author          *User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Tags            []Tag      `json:"tags,omitempty" gorm:"many2many:post_tags;"`
	PublishedAt     *time.Time `json:"published_at"`
	ScheduledAt     *time.Time `json:"scheduled_at"` // 预定发布时间
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

// PublishPost 发布文章
func PublishPost(postID uint) error {
	now := time.Now()
	return DB.Model(&Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": &now,
	}).Error
}

// GetScheduledPosts 获取已到预定发布时间但未发布的文章
func GetScheduledPosts() ([]Post, error) {
	var posts []Post
	now := time.Now()
	err := DB.Where("status = ? AND scheduled_at <= ? AND published_at IS NULL", "scheduled", now).
		Find(&posts).Error
	return posts, err
}

// GetPostByID 根据 ID 获取文章
func GetPostByID(id uint) (*Post, error) {
	var post Post
	err := DB.Preload("Category").Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
