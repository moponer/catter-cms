package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"catter-cms-gin/internal/models"
)

// MenuHandler 菜单处理器
type MenuHandler struct {
	db *gorm.DB
}

// NewMenuHandler 创建菜单处理器
func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{db: db}
}

// ListMenus 获取所有菜单
func ListMenus(c *gin.Context) {
	var menus []models.Menu
	if err := models.DB.Find(&menus).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单列表失败")
		return
	}

	// 加载每个菜单的菜单项
	for i := range menus {
		items, err := models.GetMenuItemsWithChildren(menus[i].ID)
		if err != nil {
			continue
		}
		menus[i].Items = items
	}

	SuccessResponse(c, menus)
}

// GetMenu 获取单个菜单
func GetMenu(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := models.DB.First(&menu, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, http.StatusNotFound, "菜单不存在")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单失败")
		return
	}

	// 加载菜单项
	items, err := models.GetMenuItemsWithChildren(menu.ID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单项失败")
		return
	}
	menu.Items = items

	SuccessResponse(c, menu)
}

// GetMenuByLocation 根据位置获取菜单
func GetMenuByLocation(c *gin.Context) {
	location := c.Param("location")

	var menu models.Menu
	if err := models.DB.First(&menu, "location = ?", location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, http.StatusNotFound, "该位置暂无菜单")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单失败")
		return
	}

	// 加载菜单项
	items, err := models.GetMenuItemsWithChildren(menu.ID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单项失败")
		return
	}
	menu.Items = items

	SuccessResponse(c, menu)
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	var req struct {
		Name        string           `json:"name" binding:"required"`
		Location    string           `json:"location" binding:"required"`
		Description string           `json:"description"`
		Items       []models.MenuItem `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 检查位置是否已被使用
	var existingMenu models.Menu
	if err := models.DB.First(&existingMenu, "location = ?", req.Location).Error; err == nil {
		ErrorResponse(c, http.StatusConflict, "该位置已存在菜单")
		return
	}

	menu := models.Menu{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	// 开启事务
	tx := models.DB.Begin()

	// 创建菜单
	if err := tx.Create(&menu).Error; err != nil {
		tx.Rollback()
		ErrorResponse(c, http.StatusInternalServerError, "创建菜单失败")
		return
	}

	// 创建菜单项
	if len(req.Items) > 0 {
		for i := range req.Items {
			req.Items[i].ID = uuid.New().String()
			req.Items[i].MenuID = menu.ID
			req.Items[i].Sort = i

			if err := tx.Create(&req.Items[i]).Error; err != nil {
				tx.Rollback()
				ErrorResponse(c, http.StatusInternalServerError, "创建菜单项失败")
				return
			}

			// 创建子菜单项
			if len(req.Items[i].Children) > 0 {
				for j := range req.Items[i].Children {
					req.Items[i].Children[j].ID = uuid.New().String()
					req.Items[i].Children[j].MenuID = menu.ID
					req.Items[i].Children[j].ParentID = &req.Items[i].ID
					req.Items[i].Children[j].Sort = j

					if err := tx.Create(&req.Items[i].Children[j]).Error; err != nil {
						tx.Rollback()
						ErrorResponse(c, http.StatusInternalServerError, "创建子菜单项失败")
						return
					}
				}
			}
		}
	}

	tx.Commit()
	SuccessResponse(c, menu)
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string           `json:"name"`
		Location    string           `json:"location"`
		Description string           `json:"description"`
		Items       []models.MenuItem `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	var menu models.Menu
	if err := models.DB.First(&menu, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, http.StatusNotFound, "菜单不存在")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单失败")
		return
	}

	// 更新菜单信息
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Location != "" {
		menu.Location = req.Location
	}
	menu.Description = req.Description

	tx := models.DB.Begin()

	// 保存菜单
	if err := tx.Save(&menu).Error; err != nil {
		tx.Rollback()
		ErrorResponse(c, http.StatusInternalServerError, "更新菜单失败")
		return
	}

	// 如果有新的菜单项，删除旧的并创建新的
	if req.Items != nil {
		// 删除旧的菜单项
		if err := tx.Where("menu_id = ?", id).Delete(&models.MenuItem{}).Error; err != nil {
			tx.Rollback()
			ErrorResponse(c, http.StatusInternalServerError, "更新菜单项失败")
			return
		}

		// 创建新的菜单项
		for i := range req.Items {
			req.Items[i].ID = uuid.New().String()
			req.Items[i].MenuID = menu.ID
			req.Items[i].Sort = i

			if err := tx.Create(&req.Items[i]).Error; err != nil {
				tx.Rollback()
				ErrorResponse(c, http.StatusInternalServerError, "创建菜单项失败")
				return
			}

			// 创建子菜单项
			if len(req.Items[i].Children) > 0 {
				for j := range req.Items[i].Children {
					req.Items[i].Children[j].ID = uuid.New().String()
					req.Items[i].Children[j].MenuID = menu.ID
					req.Items[i].Children[j].ParentID = &req.Items[i].ID
					req.Items[i].Children[j].Sort = j

					if err := tx.Create(&req.Items[i].Children[j]).Error; err != nil {
						tx.Rollback()
						ErrorResponse(c, http.StatusInternalServerError, "创建子菜单项失败")
						return
					}
				}
			}
		}
	}

	tx.Commit()

	// 重新加载菜单项
	items, _ := models.GetMenuItemsWithChildren(menu.ID)
	menu.Items = items

	SuccessResponse(c, menu)
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	tx := models.DB.Begin()

	// 删除菜单项
	if err := tx.Where("menu_id = ?", id).Delete(&models.MenuItem{}).Error; err != nil {
		tx.Rollback()
		ErrorResponse(c, http.StatusInternalServerError, "删除菜单项失败")
		return
	}

	// 删除菜单
	if err := tx.Delete(&models.Menu{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		ErrorResponse(c, http.StatusInternalServerError, "删除菜单失败")
		return
	}

	tx.Commit()
	SuccessResponse(c, gin.H{"message": "删除成功"})
}

// GetMenuItems 获取菜单项列表
func GetMenuItems(c *gin.Context) {
	menuID := c.Param("id")

	items, err := models.GetMenuItemsWithChildren(menuID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单项失败")
		return
	}

	SuccessResponse(c, items)
}

// CreateMenuItem 创建菜单项
func CreateMenuItem(c *gin.Context) {
	menuID := c.Param("id")

	var req models.MenuItem
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 检查菜单是否存在
	var menu models.Menu
	if err := models.DB.First(&menu, "id = ?", menuID).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "菜单不存在")
		return
	}

	req.ID = uuid.New().String()
	req.MenuID = menuID

	if req.Target == "" {
		req.Target = "_self"
	}

	if err := models.DB.Create(&req).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "创建菜单项失败")
		return
	}

	SuccessResponse(c, req)
}

// UpdateMenuItem 更新菜单项
func UpdateMenuItem(c *gin.Context) {
	itemID := c.Param("itemId")

	var req models.MenuItem
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	var item models.MenuItem
	if err := models.DB.First(&item, "id = ?", itemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, http.StatusNotFound, "菜单项不存在")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "获取菜单项失败")
		return
	}

	// 更新字段
	if req.Label != "" {
		item.Label = req.Label
	}
	if req.URL != "" {
		item.URL = req.URL
	}
	item.Icon = req.Icon
	item.Target = req.Target
	if req.ParentID != nil {
		item.ParentID = req.ParentID
	}
	item.Sort = req.Sort
	item.Enabled = req.Enabled

	if err := models.DB.Save(&item).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新菜单项失败")
		return
	}

	SuccessResponse(c, item)
}

// DeleteMenuItem 删除菜单项
func DeleteMenuItem(c *gin.Context) {
	itemID := c.Param("itemId")

	tx := models.DB.Begin()

	// 删除子菜单项
	if err := tx.Where("parent_id = ?", itemID).Delete(&models.MenuItem{}).Error; err != nil {
		tx.Rollback()
		ErrorResponse(c, http.StatusInternalServerError, "删除子菜单项失败")
		return
	}

	// 删除菜单项
	if err := tx.Delete(&models.MenuItem{}, "id = ?", itemID).Error; err != nil {
		tx.Rollback()
		ErrorResponse(c, http.StatusInternalServerError, "删除菜单项失败")
		return
	}

	tx.Commit()
	SuccessResponse(c, gin.H{"message": "删除成功"})
}

// ReorderMenuItems 重新排序菜单项
func ReorderMenuItems(c *gin.Context) {
	menuID := c.Param("id")

	var req struct {
		Items []struct {
			ID       string `json:"id"`
			Sort     int    `json:"sort"`
			ParentID string `json:"parent_id"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	tx := models.DB.Begin()

	for _, item := range req.Items {
		updates := map[string]interface{}{
			"sort": item.Sort,
		}
		if item.ParentID != "" {
			updates["parent_id"] = item.ParentID
		} else {
			updates["parent_id"] = nil
		}

		if err := tx.Model(&models.MenuItem{}).Where("id = ? AND menu_id = ?", item.ID, menuID).Updates(updates).Error; err != nil {
			tx.Rollback()
			ErrorResponse(c, http.StatusInternalServerError, "更新排序失败")
			return
		}
	}

	tx.Commit()
	SuccessResponse(c, gin.H{"message": "排序更新成功"})
}
