package services

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"catter-cms/models"
)

var (
	activeThemeCache string
	cacheMutex       sync.RWMutex
)

// GetActiveTheme 从数据库获取当前激活的主题
// 如果获取失败，返回默认主题 "blog"
func GetActiveTheme() string {
	// 先尝试从缓存读取
	cacheMutex.RLock()
	if activeThemeCache != "" {
		cacheMutex.RUnlock()
		return activeThemeCache
	}
	cacheMutex.RUnlock()

	// 从数据库查询
	theme := &models.Theme{}
	err := models.DB.Where("is_active = ?", true).First(theme)
	if err != nil {
		// 尝试从settings表获取
		setting := &models.Setting{}
		err = models.DB.Where("key = ?", "active_theme").First(setting)
		if err != nil {
			log.Printf("获取激活主题失败，使用默认主题: %v", err)
			return "blog"
		}
		// 检查主题是否有模板文件
		if !ThemeHasTemplates(setting.Value) {
			log.Printf("主题 %s 没有模板文件，回退到默认主题 blog", setting.Value)
			return "blog"
		}
		return setting.Value
	}

	// 检查主题是否有模板文件
	if !ThemeHasTemplates(theme.Slug) {
		log.Printf("主题 %s 没有模板文件，回退到默认主题 blog", theme.Slug)
		return "blog"
	}

	// 更新缓存
	cacheMutex.Lock()
	activeThemeCache = theme.Slug
	cacheMutex.Unlock()

	return theme.Slug
}

// ThemeHasTemplates 检查主题是否有必要的模板文件
func ThemeHasTemplates(themeSlug string) bool {
	// 检查主题目录是否存在
	themePath := filepath.Join("templates", "themes", themeSlug)
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return false
	}

	// 检查是否有pages目录
	pagesPath := filepath.Join(themePath, "pages")
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		return false
	}

	// 检查是否有home.plush.html模板
	homeTemplate := filepath.Join(pagesPath, "home.plush.html")
	if _, err := os.Stat(homeTemplate); os.IsNotExist(err) {
		return false
	}

	return true
}

// ClearThemeCache 清除主题缓存
// 当主题被切换时调用
func ClearThemeCache() {
	cacheMutex.Lock()
	activeThemeCache = ""
	cacheMutex.Unlock()
}

// ThemeExists 检查主题是否存在
func ThemeExists(themeSlug string) bool {
	theme := &models.Theme{}
	err := models.DB.Where("slug = ?", themeSlug).First(theme)
	return err == nil
}

// GetTemplatePath 获取主题模板路径
func GetTemplatePath(templateName string) string {
	activeTheme := GetActiveTheme()
	return "themes/" + activeTheme + "/pages/" + templateName + ".plush.html"
}

// GetThemeTemplatePath 获取指定主题的模板路径
func GetThemeTemplatePath(themeSlug, templateName string) string {
	return "themes/" + themeSlug + "/pages/" + templateName + ".plush.html"
}
