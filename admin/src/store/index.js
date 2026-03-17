import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post, put, del } from '@/utils/request'
import { getToken } from '@/utils/auth'

// 用户Store
export const useUserStore = defineStore('user', () => {
  const currentUser = ref(null)
  const token = ref(getToken())

  const getCurrentUser = async () => {
    if (!token.value) return null
    
    try {
      const res = await get('/api/auth/me')
      currentUser.value = res.data || res
      return currentUser.value
    } catch (error) {
      console.error('获取用户信息失败:', error)
      return null
    }
  }

  const logout = () => {
    currentUser.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    currentUser,
    token,
    getCurrentUser,
    logout
  }
})

// 文章Store
export const usePostStore = defineStore('post', () => {
  const posts = ref([])
  const currentPost = ref(null)
  const loading = ref(false)
  const total = ref(0)

  const fetchPosts = async (params = {}) => {
    loading.value = true
    try {
      const res = await get('/api/posts', params)
      const data = res.data || res
      posts.value = data.list || data.posts || []
      total.value = data.total || 0
    } catch (error) {
      console.error('获取文章列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchPost = async (id) => {
    loading.value = true
    try {
      const res = await get(`/api/posts/${id}`)
      currentPost.value = res.data || res
    } catch (error) {
      console.error('获取文章详情失败:', error)
    } finally {
      loading.value = false
    }
  }

  const createPost = async (postData) => {
    try {
      const res = await post('/api/posts', postData)
      return res.data || res
    } catch (error) {
      console.error('创建文章失败:', error)
      return null
    }
  }

  const updatePost = async (id, postData) => {
    try {
      const res = await put(`/api/posts/${id}`, postData)
      return res.data || res
    } catch (error) {
      console.error('更新文章失败:', error)
      return null
    }
  }

  const deletePost = async (id) => {
    try {
      await del(`/api/posts/${id}`)
      return true
    } catch (error) {
      console.error('删除文章失败:', error)
      return false
    }
  }

  return {
    posts,
    currentPost,
    loading,
    total,
    fetchPosts,
    fetchPost,
    createPost,
    updatePost,
    deletePost
  }
})

// 分类Store
export const useCategoryStore = defineStore('category', () => {
  const categories = ref([])
  const loading = ref(false)

  const fetchCategories = async () => {
    loading.value = true
    try {
      const res = await get('/api/categories')
      categories.value = res.data || res || []
    } catch (error) {
      console.error('获取分类列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  return {
    categories,
    loading,
    fetchCategories
  }
})
