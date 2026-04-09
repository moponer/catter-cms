/**
 * Axios 请求封装
 * 统一处理请求拦截、响应拦截、错误处理
 */
import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { clearAuth } from '@/utils/auth'

// 创建 axios 实例
const request = axios.create({
  baseURL: '',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 是否正在处理登出，防止重复跳转
let isHandlingAuthError = false

// 处理认证错误
const handleAuthError = () => {
  if (isHandlingAuthError) return
  isHandlingAuthError = true

  // 清除认证信息
  clearAuth()

  // 提示用户
  ElMessage.error('登录已过期，请重新登录')

  // 延迟跳转，让提示先显示
  setTimeout(() => {
    // 使用 router 跳转，不用 window.location
    router.push('/login').then(() => {
      // 跳转完成后重置标志
      isHandlingAuthError = false
    }).catch(() => {
      isHandlingAuthError = false
    })
  }, 1500)
}

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 从 localStorage 获取 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    console.warn('Request config error:', error.message || error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    // 检查 Content-Type，确保是 JSON
    const contentType = response.headers['content-type'] || ''
    if (contentType.includes('application/json')) {
      const res = response.data

      // 检查业务状态码
      if (res.code && res.code !== 0 && res.code !== 200) {
        // 如果是未授权错误，交由错误拦截器处理
        if (res.code === 401 || res.code === 403) {
          handleAuthError()
          return Promise.reject(new Error(res.message || '登录已过期'))
        }
        ElMessage.error(res.message || '请求失败')
        return Promise.reject(new Error(res.message || '请求失败'))
      }

      return res
    }

    // 非 JSON 响应，如果是成功状态码直接返回
    if (response.status >= 200 && response.status < 300) {
      return response.data
    }

    // 其他情况当作错误处理
    throw new Error('Invalid response type')
  },
  (error) => {
    const { response, message } = error

    if (response) {
      // 服务器返回了响应
      const status = response.status
      const contentType = response.headers['content-type'] || ''

      // 检查是否是 JSON 响应
      if (contentType.includes('application/json')) {
        const data = response.data || {}

        switch (status) {
          case 401:
            handleAuthError()
            return Promise.reject(new Error(data.message || data.error || '登录已过期'))
          case 403:
            ElMessage.error(data.message || data.error || '没有权限访问')
            break
          case 404:
            ElMessage.error(data.message || data.error || '请求的资源不存在')
            break
          case 500:
            ElMessage.error(data.message || data.error || '服务器错误')
            break
          default:
            ElMessage.error(data.message || data.error || '请求失败')
        }
      } else {
        // 非 JSON 响应
        switch (status) {
          case 401:
          case 403:
            handleAuthError()
            return Promise.reject(new Error('登录已过期'))
          case 404:
            ElMessage.error('请求的资源不存在')
            break
          case 500:
            ElMessage.error('服务器错误')
            break
          default:
            ElMessage.error('请求失败')
        }
      }
    } else if (error.request) {
      // 请求已发送但没有收到响应（后端未启动或网络问题）
      ElMessage.error('无法连接到服务器，请检查后端服务是否运行')
      // 清除认证状态，让用户重新登录
      handleAuthError()
    } else {
      // 请求配置错误
      ElMessage.error(message || '请求配置错误')
    }

    return Promise.reject(error)
  }
)

export default request

/**
 * GET 请求
 */
export function get(url, params, config = {}) {
  return request({
    method: 'GET',
    url,
    params,
    ...config,
  })
}

/**
 * POST 请求
 */
export function post(url, data, config = {}) {
  return request({
    method: 'POST',
    url,
    data,
    ...config,
  })
}

/**
 * PUT 请求
 */
export function put(url, data, config = {}) {
  return request({
    method: 'PUT',
    url,
    data,
    ...config,
  })
}

/**
 * DELETE 请求
 * params 为对象时作为 URL query string，config.data 时作为 request body
 */
export function del(url, params, config = {}) {
  return request({
    method: 'DELETE',
    url,
    params,
    ...config,
  })
}

/**
 * DELETE 请求（带 body）
 * 用于需要 JSON body 的 DELETE 接口（如批量删除）
 */
export function delBody(url, data, config = {}) {
  return request({
    method: 'DELETE',
    url,
    data,
    ...config,
  })
}

/**
 * 上传文件
 */
export function upload(url, file, onProgress) {
  const formData = new FormData()
  formData.append('file', file)

  return request({
    method: 'POST',
    url,
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress) {
        const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        onProgress(percent)
      }
    },
  })
}
