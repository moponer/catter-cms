import request from '@/utils/request'

// 获取主题列表
export function listThemes() {
  return request({
    url: '/api/themes',
    method: 'get'
  })
}

export function getThemes(params) {
  return request({
    url: '/api/themes',
    method: 'get',
    params
  })
}

// 获取当前主题
export function getCurrentTheme() {
  return request({
    url: '/api/themes/current',
    method: 'get'
  })
}

// 切换主题
export function switchTheme(themeId) {
  return request({
    url: `/api/themes/${themeId}/activate`,
    method: 'post'
  })
}

// 获取主题配置
export function getThemeConfig(themeId) {
  const url = themeId ? `/api/themes/${themeId}/config` : '/api/themes/current/config'
  return request({
    url,
    method: 'get'
  })
}

// 获取主题完整配置（含 schema 定义和已保存值）
export function getThemeConfigFull(themeId) {
  return request({
    url: `/api/themes/${themeId}/config/full`,
    method: 'get'
  })
}

// 保存主题配置
export function saveThemeConfig(data, themeId) {
  const url = themeId ? `/api/themes/${themeId}/config` : '/api/themes/current/config'
  return request({
    url,
    method: 'put',
    data
  })
}

// 获取主题模块配置
export function getThemeModules(themeId) {
  const url = themeId ? `/api/themes/${themeId}/modules` : '/api/themes/current/modules'
  return request({
    url,
    method: 'get'
  })
}

// 更新模块配置
export function updateModuleConfig(themeId, moduleId, data) {
  const url = themeId
    ? `/api/themes/${themeId}/modules/${moduleId}`
    : `/api/themes/current/modules/${moduleId}`
  return request({
    url,
    method: 'put',
    data
  })
}

// 更新模块排序
export function updateModuleOrder(themeId, moduleIds) {
  const url = themeId
    ? `/api/themes/${themeId}/modules/order`
    : '/api/themes/current/modules/order'
  return request({
    url,
    method: 'put',
    data: { moduleIds }
  })
}

// 获取页面模板列表
export function getPageTemplates(themeId) {
  const url = themeId
    ? `/api/themes/${themeId}/templates`
    : '/api/themes/current/templates'
  return request({
    url,
    method: 'get'
  })
}

// 获取主题设置定义
export function getThemeSettingsSchema(themeId) {
  const url = themeId
    ? `/api/themes/${themeId}/schema`
    : '/api/themes/current/schema'
  return request({
    url,
    method: 'get'
  })
}

// 重置主题配置
export function resetThemeConfig(themeId) {
  const url = themeId
    ? `/api/themes/${themeId}/config/reset`
    : '/api/themes/current/config/reset'
  return request({
    url,
    method: 'post'
  })
}

// 导入主题配置
export function importThemeConfig(themeId, data) {
  const url = themeId
    ? `/api/themes/${themeId}/config/import`
    : '/api/themes/current/config/import'
  return request({
    url,
    method: 'post',
    data
  })
}

// 导出主题配置
export function exportThemeConfig(themeId) {
  const url = themeId
    ? `/api/themes/${themeId}/config/export`
    : '/api/themes/current/config/export'
  return request({
    url,
    method: 'get',
    responseType: 'blob'
  })
}

// 删除主题
export function deleteTheme(themeId) {
  return request({
    url: `/api/themes/${themeId}`,
    method: 'delete'
  })
}

// 清理主题缓存
export function clearThemeCache(themeId) {
  return request({
    url: `/api/themes/${themeId}/clear-cache`,
    method: 'post'
  })
}

// 清理所有缓存（全局）
export function clearAllCache() {
  return request({
    url: '/api/admin/cache/clear',
    method: 'post',
    data: { type: 'all' }
  })
}
