/**
 * Axios 请求封装
 * 统一处理请求拦截、响应拦截、错误处理
 */
import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
// 注意：baseURL 设置为空，因为请求路径已经包含 /api 前缀
// 例如：post('/api/auth/login') 会请求 /api/auth/login
const request = axios.create({
  baseURL: '',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

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
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data
    
    // 如果响应直接是数据，返回
    if (response.config.returnRaw) {
      return res
    }
    
    // 检查业务状态码
    if (res.code && res.code !== 0 && res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    
    return res
  },
  (error) => {
    const { response } = error
    
    if (response) {
      switch (response.status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          // 清除 token 并跳转登录页
          localStorage.removeItem('token')
          window.location.href = '/admin/login'
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器错误')
          break
        default:
          ElMessage.error(response.data?.message || '请求失败')
      }
    } else {
      ElMessage.error('网络连接失败')
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
