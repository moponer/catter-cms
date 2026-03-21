import request from '@/utils/request'

// 获取主题列表
export function getThemes(params) {
  return request({
    url: '/themes',
    method: 'get',
    params
  })
}

// 获取当前主题
export function getCurrentTheme() {
  return request({
    url: '/themes/current',
    method: 'get'
  })
}

// 切换主题
export function switchTheme(themeId) {
  return request({
    url: `/themes/${themeId}/activate`,
    method: 'post'
  })
}

// 获取主题配置
export function getThemeConfig(themeId) {
  const url = themeId ? `/themes/${themeId}/config` : '/themes/current/config'
  return request({
    url,
    method: 'get'
  })
}

// 保存主题配置
export function saveThemeConfig(data, themeId) {
  const url = themeId ? `/themes/${themeId}/config` : '/themes/current/config'
  return request({
    url,
    method: 'put',
    data
  })
}

// 获取主题模块配置
export function getThemeModules(themeId) {
  const url = themeId ? `/themes/${themeId}/modules` : '/themes/current/modules'
  return request({
    url,
    method: 'get'
  })
}

// 更新模块配置
export function updateModuleConfig(themeId, moduleId, data) {
  const url = themeId
    ? `/themes/${themeId}/modules/${moduleId}`
    : `/themes/current/modules/${moduleId}`
  return request({
    url,
    method: 'put',
    data
  })
}

// 更新模块排序
export function updateModuleOrder(themeId, moduleIds) {
  const url = themeId
    ? `/themes/${themeId}/modules/order`
    : '/themes/current/modules/order'
  return request({
    url,
    method: 'put',
    data: { moduleIds }
  })
}

// 获取页面模板列表
export function getPageTemplates(themeId) {
  const url = themeId
    ? `/themes/${themeId}/templates`
    : '/themes/current/templates'
  return request({
    url,
    method: 'get'
  })
}

// 获取主题设置定义
export function getThemeSettingsSchema(themeId) {
  const url = themeId
    ? `/themes/${themeId}/schema`
    : '/themes/current/schema'
  return request({
    url,
    method: 'get'
  })
}

// 重置主题配置
export function resetThemeConfig(themeId) {
  const url = themeId
    ? `/themes/${themeId}/config/reset`
    : '/themes/current/config/reset'
  return request({
    url,
    method: 'post'
  })
}

// 导入主题配置
export function importThemeConfig(themeId, data) {
  const url = themeId
    ? `/themes/${themeId}/config/import`
    : '/themes/current/config/import'
  return request({
    url,
    method: 'post',
    data
  })
}

// 导出主题配置
export function exportThemeConfig(themeId) {
  const url = themeId
    ? `/themes/${themeId}/config/export`
    : '/themes/current/config/export'
  return request({
    url,
    method: 'get',
    responseType: 'blob'
  })
}
