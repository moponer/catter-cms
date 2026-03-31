package services

import (
	"strings"
	"time"

	"catter-cms-gin/internal/models"
)

// GenerateSlug 生成 URL 友好的 slug
func GenerateSlug(title string) string {
	slug := ""
	hasAlphaNum := false

	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			slug += string(r)
			hasAlphaNum = true
		} else if r == ' ' || r == '-' {
			slug += "-"
		}
	}

	// 如果标题全是中文或其他非拉丁字符，使用时间戳
	if !hasAlphaNum || slug == "" {
		slug = "post-" + time.Now().Format("20060102150405")
	}

	return slug
}

// SyncPostTags 同步文章标签（支持 tag_names 和 tag_ids）
func SyncPostTags(post *models.Post, tagNames []string, tagIDs []uint) error {
	var postTags []models.Tag

	if len(tagNames) > 0 {
		// 使用标签名称，自动创建不存在的标签
		for _, name := range tagNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}

			var tag models.Tag
			result := models.DB.Where("name = ?", name).First(&tag)
			if result.Error != nil {
				// 标签不存在，创建新标签
				tag = models.Tag{
					Name: name,
					Slug: GenerateSlug(name),
				}
				models.DB.Create(&tag)
			}
			postTags = append(postTags, tag)
		}
	} else if len(tagIDs) > 0 {
		// 使用标签 ID
		models.DB.Find(&postTags, tagIDs)
	}

	// 替换标签关联
	if len(postTags) > 0 {
		models.DB.Model(post).Association("Tags").Replace(&postTags)
	} else {
		// 如果没有标签，清空关联
		models.DB.Model(post).Association("Tags").Clear()
	}

	return nil
}