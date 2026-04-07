package main

import (
	"fmt"
	"log"

	"catter-cms-gin/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("../../catter.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	models.DB = db

	// 创建菜单
	menus := []models.Menu{
		{
			ID:          "menu-header",
			Name:        "顶部导航菜单",
			Location:    "header",
			Description: "网站顶部导航菜单",
		},
		{
			ID:          "menu-footer",
			Name:        "底部导航菜单",
			Location:    "footer",
			Description: "网站底部导航菜单",
		},
		{
			ID:          "menu-sidebar",
			Name:        "侧边栏菜单",
			Location:    "sidebar",
			Description: "侧边栏导航菜单",
		},
	}

	// 创建菜单项
	menuItems := []models.MenuItem{
		// Header菜单项
		{ID: "item-home", MenuID: "menu-header", Label: "首页", URL: "/", Sort: 1, Enabled: true},
		{ID: "item-blog", MenuID: "menu-header", Label: "文章", URL: "/blog", Sort: 2, Enabled: true},
		{ID: "item-categories", MenuID: "menu-header", Label: "分类", URL: "/categories", Sort: 3, Enabled: true},
		{ID: "item-about", MenuID: "menu-header", Label: "关于", URL: "/page/about", Sort: 4, Enabled: true},

		// Footer菜单项
		{ID: "item-footer-home", MenuID: "menu-footer", Label: "首页", URL: "/", Sort: 1, Enabled: true},
		{ID: "item-footer-blog", MenuID: "menu-footer", Label: "文章", URL: "/blog", Sort: 2, Enabled: true},
		{ID: "item-footer-contact", MenuID: "menu-footer", Label: "联系我们", URL: "/page/contact", Sort: 3, Enabled: true},

		// Sidebar菜单项
		{ID: "item-sidebar-all", MenuID: "menu-sidebar", Label: "全部文章", URL: "/blog", Sort: 1, Enabled: true},
		{ID: "item-sidebar-tech", MenuID: "menu-sidebar", Label: "技术文章", URL: "/category/tech", Sort: 2, Enabled: true},
		{ID: "item-sidebar-life", MenuID: "menu-sidebar", Label: "生活随笔", URL: "/category/life", Sort: 3, Enabled: true},
	}

	// 清空现有菜单数据
	db.Exec("DELETE FROM menu_items")
	db.Exec("DELETE FROM menus")

	// 插入菜单
	for _, menu := range menus {
		if err := db.Create(&menu).Error; err != nil {
			log.Printf("Failed to create menu %s: %v", menu.ID, err)
		} else {
			fmt.Printf("✓ Created menu: %s (%s)\n", menu.Name, menu.Location)
		}
	}

	// 插入菜单项
	for _, item := range menuItems {
		if err := db.Create(&item).Error; err != nil {
			log.Printf("Failed to create menu item %s: %v", item.ID, err)
		} else {
			fmt.Printf("✓ Created menu item: %s -> %s\n", item.Label, item.URL)
		}
	}

	fmt.Println("\n✅ 菜单数据初始化完成！")
}
