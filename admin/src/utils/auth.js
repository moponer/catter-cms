/**
 * 认证相关工具函数
 */

const TOKEN_KEY = 'token'
const USER_KEY = 'user'

/**
 * 获取 Token (已废弃，改用 httpOnly Cookie)
 */
export function getToken() {
  // httpOnly cookie 由浏览器自动发送，前端 JS 无法也无需获取
  return null
}

/**
 * 设置 Token (已废弃，改用 httpOnly Cookie)
 */
export function setToken(token) {
  // 空操作
}

/**
 * 移除 Token (已废弃，改用 httpOnly Cookie)
 */
export function removeToken() {
  // 登出时 Cookie 会由后端清除
}

/**
 * 检查是否已登录 (前端无法直接判断 httpOnly cookie)
 * 依赖 API 请求返回 401 来判断会话是否有效
 */
export function isAuthenticated() {
  // 这个函数不再可靠，因为前端无法读取 httpOnly cookie。
  // 真实的认证状态由每次 API 请求的结果决定。
  // 暂时返回 true，让路由守卫放行，具体页面再根据API返回的401做跳转。
  return true
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
 * 登出 - 调用后端接口清除 Cookie，然后跳转
 * 注意：需要在 Vue 应用中调用
 */
export async function logout() {
  try {
    // 调用后端接口清除 httpOnly Cookie
    await fetch('/api/auth/logout', {
      method: 'POST',
      credentials: 'include'
    })
  } catch (e) {
    // 即使接口调用失败，也继续执行清理和跳转
    console.warn('Logout API call failed:', e)
  } finally {
    // 清除本地认证信息
    clearAuth()
    // 跳转到登录页
    if (typeof window !== 'undefined') {
      try {
        import('@/router').then(({ default: router }) => {
          router.push('/login').catch(() => {
            window.location.href = '/admin/login'
          })
        }).catch(() => {
          window.location.href = '/admin/login'
        })
      } catch (e) {
        window.location.href = '/admin/login'
      }
    }
  }
}

/**
 * 跳转到登录页（带返回）
 */
export function redirectToLogin(returnUrl = null) {
  const url = returnUrl || window.location.href
  window.location.href = `/admin/login?redirect=${encodeURIComponent(url)}`
}
