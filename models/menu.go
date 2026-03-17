package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Menu 定义菜单模型（支持多级菜单）
type Menu struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Location    string    `json:"location" db:"location"`
	Description string    `json:"description,omitempty" db:"description"`
	Order       int       `json:"order" db:"order"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Items       MenuItems `json:"items,omitempty" has_many:"menu_items" order_by:"order asc"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// Validate 验证菜单数据
func (m *Menu) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Name, Name: "Name"},
		&validators.StringIsPresent{Field: m.Location, Name: "Location"},
	), nil
}

// Menus 菜单集合
type Menus []Menu

// MenuItem 定义菜单项模型（支持多级）
type MenuItem struct {
	ID          int        `json:"id" db:"id"`
	MenuID      int        `json:"menu_id" db:"menu_id"`
	ParentID    *int       `json:"parent_id,omitempty" db:"parent_id"`
	Title       string     `json:"title" db:"title"`
	URL         string     `json:"url" db:"url"`
	Target      string     `json:"target" db:"target"`
	Classes     string     `json:"classes,omitempty" db:"classes"`
	Icon        string     `json:"icon,omitempty" db:"icon"`
	Order       int        `json:"order" db:"order"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Menu        Menu       `json:"menu,omitempty" belongs_to:"menus"`
	Parent      *MenuItem  `json:"parent,omitempty" belongs_to:"menu_items"`
	Children    MenuItems  `json:"children,omitempty" has_many:"menu_items" order_by:"order asc"`
}

// TableName 指定表名
func (MenuItem) TableName() string {
	return "menu_items"
}

// Validate 验证菜单项数据
func (mi *MenuItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: mi.Title, Name: "Title"},
		&validators.StringIsPresent{Field: mi.URL, Name: "URL"},
	), nil
}

// MenuItems 菜单项集合
type MenuItems []MenuItem
