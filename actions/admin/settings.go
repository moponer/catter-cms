package admin

import (
	"strconv"

	"catter-cms/actions/api"
	"catter-cms/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// SettingsGetAPI 获取系统设置API
func SettingsGetAPI(c buffalo.Context) error {
	settings := make(map[string]interface{})

	basicSettings := make(map[string]string)
	seoSettings := make(map[string]string)
	emailSettings := make(map[string]string)
	themeSettings := make(map[string]string)

	dbSettings := []models.Setting{}
	if err := models.DB.All(&dbSettings); err == nil {
		for _, setting := range dbSettings {
			switch setting.Category {
			case "basic":
				basicSettings[setting.Key] = setting.Value
			case "seo":
				seoSettings[setting.Key] = setting.Value
			case "email":
				emailSettings[setting.Key] = setting.Value
			case "theme":
				themeSettings[setting.Key] = setting.Value
			}
		}
	}

	settings["basic"] = basicSettings
	settings["seo"] = seoSettings
	settings["email"] = emailSettings
	settings["theme"] = themeSettings

	return api.Success(c, settings)
}

// SettingsSaveAPI 保存系统设置API
func SettingsSaveAPI(c buffalo.Context) error {
	requestData := map[string]interface{}{}
	if err := c.Bind(&requestData); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	category, ok := requestData["category"].(string)
	if !ok {
		return api.BadRequest(c, "缺少分类参数")
	}

	settingsData, ok := requestData["settings"].(map[string]interface{})
	if !ok {
		return api.BadRequest(c, "缺少设置数据")
	}

	for key, value := range settingsData {
		setting := &models.Setting{}
		if err := models.DB.Where("key = ? AND category = ?", key, category).First(setting); err != nil {
			// 创建新设置
			setting = &models.Setting{
				Key:      key,
				Value:    toString(value),
				Category: category,
				Type:     "string",
			}
			if err := models.DB.Create(setting); err != nil {
				return api.InternalError(c, "保存设置失败")
			}
		} else {
			// 更新设置
			setting.Value = toString(value)
			if err := models.DB.Update(setting); err != nil {
				return api.InternalError(c, "更新设置失败")
			}
		}
	}

	return api.SuccessWithMessage(c, "设置保存成功", nil)
}

// SettingsTestEmailAPI 测试邮件发送API
func SettingsTestEmailAPI(c buffalo.Context) error {
	emailSettings := map[string]interface{}{}
	if err := c.Bind(&emailSettings); err != nil {
		return api.BadRequest(c, "参数错误")
	}

	// TODO: 实际发送测试邮件
	// 这里只是模拟成功

	return api.SuccessWithMessage(c, "测试邮件发送成功", nil)
}

// SettingsIndexHandler 系统设置页面 (已弃用，SPA模式)
func SettingsIndexHandler(c buffalo.Context) error {
	return c.Render(200, render.HTML("<html><body><p>Settings SPA</p></body></html>"))
}

// SettingsGetHandler 获取系统设置 (兼容旧API)
func SettingsGetHandler(c buffalo.Context) error {
	return SettingsGetAPI(c)
}

// SettingsSaveHandler 保存系统设置 (兼容旧API)
func SettingsSaveHandler(c buffalo.Context) error {
	return SettingsSaveAPI(c)
}

// SettingsTestEmailHandler 测试邮件发送 (兼容旧API)
func SettingsTestEmailHandler(c buffalo.Context) error {
	return SettingsTestEmailAPI(c)
}

// toString 将各种类型转换为字符串
func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		// 处理 JSON 解析后的整数（会被解析为 float64）
		if v == float64(int(v)) {
			return strconv.Itoa(int(v))
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}
