/**
 * 通用类型定义
 */

// API 响应格式
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页请求参数
export interface PaginationParams {
  page?: number
  page_size?: number
}

// 分页响应数据
export interface PaginationData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 用户
export interface User {
  id: number
  name: string
  email: string
  role: string
  avatar?: string
  created_at: string
  updated_at: string
}

// 文章
export interface Post {
  id: number
  title: string
  slug: string
  content: string
  summary?: string
  status: 'draft' | 'published' | 'archived'
  author_id: number
  author?: User
  category_id?: number
  category?: Category
  tags?: Tag[]
  created_at: string
  updated_at: string
  published_at?: string
}

// 分类
export interface Category {
  id: number
  name: string
  slug: string
  description?: string
  parent_id?: number
  children?: Category[]
  created_at: string
  updated_at: string
}

// 标签
export interface Tag {
  id: number
  name: string
  slug: string
  created_at: string
  updated_at: string
}

// 媒体
export interface Media {
  id: number
  name: string
  path: string
  url: string
  type: string
  size: number
  created_at: string
}

// 主题
export interface Theme {
  id: string
  name: string
  slug: string
  version: string
  description: string
  author: string
  thumbnail: string
  is_active: boolean
  config?: ThemeConfig[]
}

// 主题配置项
export interface ThemeConfig {
  name: string
  label: string
  type: 'text' | 'number' | 'color' | 'switch' | 'select' | 'textarea' | 'group'
  default?: any
  options?: { value: string; label: string }[]
  fields?: ThemeConfig[]
}

// 系统设置
export interface Setting {
  key: string
  value: string
  type: string
  category: string
}

// 仪表盘统计
export interface DashboardStats {
  total_posts: number
  published_posts: number
  draft_posts: number
  total_users: number
  total_categories: number
  total_tags: number
  total_media: number
}

// 登录请求
export interface LoginRequest {
  email: string
  password: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: User
}
