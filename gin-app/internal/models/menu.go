package models

import (
	"time"
)

// Menu 菜单模型
type Menu struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Location    string     `json:"location" gorm:"uniqueIndex;not null"` // header, footer, sidebar, mobile
	Description string     `json:"description"`
	Items       []MenuItem `json:"items,omitempty" gorm:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// MenuItem 菜单项模型
type MenuItem struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	MenuID    string     `json:"menu_id" gorm:"index"`
	Label     string     `json:"label" gorm:"not null"`
	URL       string     `json:"url" gorm:"not null"`
	Icon      string     `json:"icon"`
	Target    string     `json:"target" gorm:"default:'_self'"`
	ParentID  *string    `json:"parent_id" gorm:"index"`
	Sort      int        `json:"sort" gorm:"default:0"`
	Enabled   bool       `json:"enabled" gorm:"default:true"`
	Children  []MenuItem `json:"children,omitempty" gorm:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// TableName 指定表名
func (MenuItem) TableName() string {
	return "menu_items"
}

// GetMenuByLocation 根据位置获取菜单
func GetMenuByLocation(location string) (*Menu, error) {
	var menu Menu
	if err := DB.Where("location = ?", location).First(&menu).Error; err != nil {
		return nil, err
	}

	// 加载菜单项
	items, err := GetMenuItemsWithChildren(menu.ID)
	if err != nil {
		return nil, err
	}
	menu.Items = items

	return &menu, nil
}

// GetMenuItemsWithChildren 获取带层级结构的菜单项（只获取启用的项）
func GetMenuItemsWithChildren(menuID string) ([]MenuItem, error) {
	var allItems []MenuItem
	if err := DB.Where("menu_id = ? AND enabled = ?", menuID, true).Order("sort ASC").Find(&allItems).Error; err != nil {
		return nil, err
	}

	// 构建层级结构
	itemMap := make(map[string]*MenuItem)
	var rootItems []MenuItem

	// 先创建映射
	for i := range allItems {
		itemMap[allItems[i].ID] = &allItems[i]
		allItems[i].Children = nil
	}

	// 构建树形结构
	for i := range allItems {
		if allItems[i].ParentID == nil {
			rootItems = append(rootItems, allItems[i])
		} else {
			parent := itemMap[*allItems[i].ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, allItems[i])
			}
		}
	}

	return rootItems, nil
}
