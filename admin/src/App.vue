<template>
  <div class="app-container">
    <!-- 登录页面 - 不需要Layout -->
    <template v-if="$route.path === '/login'">
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
import { AppLayout } from '@/components/layout'
import GlobalSearch from '@/components/GlobalSearch.vue'
</script>

<style scoped>
.app-container {
  width: 100%;
  min-height: 100vh;
}
</style>
