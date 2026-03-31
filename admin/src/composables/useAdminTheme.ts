/**
 * Admin Theme Composable
 * 管理后台双主题切换：Dark Pro / Light Minimal
 */

import { ref, watch, onMounted } from 'vue'
import { useLocalStorage } from '@vueuse/core'

export interface AdminTheme {
  id: string
  name: string
  description: string
  icon: string // moon 或 sun
}

// 双主题定义
export const ADMIN_THEMES: AdminTheme[] = [
  {
    id: 'dark',
    name: 'Dark Pro',
    description: '深色专业版',
    icon: 'moon'
  },
  {
    id: 'light',
    name: 'Light Minimal',
    description: '浅色简约版',
    icon: 'sun'
  }
]

// 默认主题
const defaultTheme = 'dark'

// 持久化
const currentThemeId = useLocalStorage('catter-admin-theme', defaultTheme)

// 过渡状态
const isTransitioning = ref(false)

// 应用主题到 DOM
const applyTheme = (themeId: string) => {
  isTransitioning.value = true

  const html = document.documentElement

  // 移除旧主题
  html.removeAttribute('data-admin-theme')

  // 应用新主题（默认dark不需要属性）
  if (themeId !== defaultTheme) {
    html.setAttribute('data-admin-theme', themeId)
  }

  // 触发主题变化事件
  window.dispatchEvent(new CustomEvent('theme-change', { detail: { themeId } }))

  // 移除过渡状态
  setTimeout(() => {
    isTransitioning.value = false
  }, 200)
}

// 主题切换 Hook
export const useAdminTheme = () => {
  const currentTheme = ref<AdminTheme>(ADMIN_THEMES[0])

  // 初始化
  onMounted(() => {
    const savedTheme = ADMIN_THEMES.find(t => t.id === currentThemeId.value)
    if (savedTheme) {
      currentTheme.value = savedTheme
      applyTheme(savedTheme.id)
    }
  })

  // 监听变化
  watch(currentThemeId, (newThemeId) => {
    const theme = ADMIN_THEMES.find(t => t.id === newThemeId)
    if (theme) {
      currentTheme.value = theme
      applyTheme(theme.id)
    }
  })

  // 切换主题
  const setTheme = (themeId: string) => {
    if (ADMIN_THEMES.some(t => t.id === themeId)) {
      currentThemeId.value = themeId
    }
  }

  // 获取图标类型
  const getThemeIcon = (theme: AdminTheme) => theme.icon

  return {
    currentThemeId,
    currentTheme,
    isTransitioning,
    setTheme,
    getThemeIcon,
    themes: ADMIN_THEMES
  }
}
