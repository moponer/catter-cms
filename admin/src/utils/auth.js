/**
 * 认证相关工具函数
 */

const TOKEN_KEY = 'token'
const USER_KEY = 'user'

/**
 * 获取 Token
 */
export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

/**
 * 设置 Token
 */
export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

/**
 * 移除 Token
 */
export function removeToken() {
  localStorage.removeItem(TOKEN_KEY)
}

/**
 * 检查是否已登录
 */
export function isAuthenticated() {
  return !!getToken()
}

/**
 * 获取当前用户信息
 */
export function getCurrentUser() {
  const userStr = localStorage.getItem(USER_KEY)
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch {
      return null
    }
  }
  return null
}

/**
 * 设置当前用户信息
 */
export function setCurrentUser(user) {
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

/**
 * 移除当前用户信息
 */
export function removeCurrentUser() {
  localStorage.removeItem(USER_KEY)
}

/**
 * 清除所有认证信息
 */
export function clearAuth() {
  removeToken()
  removeCurrentUser()
}

/**
 * 登出
 */
export function logout() {
  clearAuth()
  window.location.href = '/admin/login'
}
