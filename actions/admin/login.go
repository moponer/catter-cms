package admin

import (
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"

	"catter-cms/models"
)

// LoginHandler 显示登录页面 (SPA模式，返回简单页面)
func LoginHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, render.HTML("<html><head><title>Login</title></head><body><p>Admin SPA Login</p></body></html>"))
}

// LoginPostHandler 处理登录提交
func LoginPostHandler(c buffalo.Context) error {
	// 解析表单数据
	if err := c.Request().ParseForm(); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	email := c.Request().FormValue("email")
	password := c.Request().FormValue("password")

	// 查找用户
	user := &models.User{}
	query := models.DB.Where("email = ?", email)
	if err := query.First(user); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return c.Redirect(http.StatusSeeOther, "/admin/login")
	}

	// 设置会话
	c.Session().Set("user_id", strconv.Itoa(user.ID))

	// 重定向到仪表盘
	return c.Redirect(http.StatusSeeOther, "/admin")
}

// LogoutHandler 处理登出
func LogoutHandler(c buffalo.Context) error {
	// 清除会话
	c.Session().Clear()

	// 重定向到登录页
	return c.Redirect(http.StatusSeeOther, "/admin/login")
}
