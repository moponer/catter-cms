package models

import (
	"time"

	"gorm.io/gorm"
)

// SEOSettings 全站 SEO 设置
type SEOSettings struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	SiteURL          string    `json:"site_url" gorm:"size:1000;not null"`           // 站点 URL (用于 canonical)
	IndexNowKey      string    `json:"index_now_key" gorm:"size:1000"`               // IndexNow API Key
	IndexNowPushed   bool      `json:"index_now_pushed" gorm:"default:false"`        // 是否已完成首次 IndexNow 注册
	SitemapEnabled   bool      `json:"sitemap_enabled" gorm:"default:true"`          // 是否启用 sitemap
	RobotsEnabled    bool      `json:"robots_enabled" gorm:"default:true"`           // 是否启用 robots.txt
	FeedEnabled      bool      `json:"feed_enabled" gorm:"default:true"`             // 是否启用 RSS/Atom feed
	FeedLimit        int       `json:"feed_limit" gorm:"default:20"`                 // Feed 文章数量
	DefaultOgImage   string    `json:"default_og_image" gorm:"size:1000"`            // 默认 OG 图片（无封面时使用）
	DefaultMetaDesc  string    `json:"default_meta_desc" gorm:"size:500"`            // 默认 Meta 描述
	DefaultRobots    string    `json:"default_robots" gorm:"size:100;default:index,nofollow"` // 默认 robots 指令
	Description      string    `json:"description" gorm:"size:1000"`                 // 设置描述
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName 指定表名
func (SEOSettings) TableName() string {
	return "seo_settings"
}

// GetSEOSettings 获取 SEO 设置（单例）
func GetSEOSettings() (*SEOSettings, error) {
	var settings SEOSettings
	err := DB.First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// 初始化默认设置
		settings = SEOSettings{
			SitemapEnabled:  true,
			RobotsEnabled:   true,
			FeedEnabled:     true,
			FeedLimit:       20,
			DefaultRobots:   "index, follow",
		}
		err = DB.Create(&settings).Error
	}
	return &settings, err
}

// UpdateSEOSettings 更新 SEO 设置
func UpdateSEOSettings(updates map[string]interface{}) error {
	var settings SEOSettings
	if err := DB.First(&settings).Error; err != nil {
		return err
	}
	return DB.Model(&settings).Updates(updates).Error
}
