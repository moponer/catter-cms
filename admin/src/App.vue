<template>
  <div class="app-container">
    <!-- 登录页面 - 不需要Layout (支持 /login 和 /admin/login) -->
    <template v-if="isLoginPage">
      <router-view />
    </template>

    <!-- 管理后台页面 - 需要Layout -->
    <template v-else>
      <AppLayout>
        <router-view />
      </AppLayout>
      <!-- 全局搜索组件 -->
      <GlobalSearch />
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { AppLayout } from '@/components/layout'
import GlobalSearch from '@/components/GlobalSearch.vue'
import { useAdminTheme } from '@/composables/useAdminTheme'

const route = useRoute()

// 初始化主题
const { currentTheme } = useAdminTheme()

const isLoginPage = computed(() => {
  const path = route.path
  return path === '/login' || path === '/admin/login'
})

// 确保主题在挂载时应用
onMounted(() => {
  // 主题已在 composable 的 onMounted 中自动应用
})
</script>

<style scoped>
.app-container {
  width: 100%;
  min-height: 100vh;
}
</style>
