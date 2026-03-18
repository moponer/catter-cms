package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name" gorm:"not null;size:255"`
	Slug        string      `json:"slug" gorm:"uniqueIndex;not null;size:255"`
	Description string      `json:"description" gorm:"type:text"`
	ParentID    *uint       `json:"parent_id"`
	Parent      *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []Category  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	SortOrder   int         `json:"sort_order" gorm:"default:0"`
	Posts       []Post      `json:"posts,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// BeforeDelete 删除前钩子 - 防止删除有子分类的分类
func (c *Category) BeforeDelete(tx *gorm.DB) error {
	var count int64
	tx.Model(&Category{}).Where("parent_id = ?", c.ID).Count(&count)
	if count > 0 {
		return gorm.ErrInvalidData
	}
	return nil
}

// GetPostCount 获取分类下的文章数量
func (c *Category) GetPostCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&Post{}).Where("category_id = ? AND status = ?", c.ID, "published").Count(&count)
	return count
}
