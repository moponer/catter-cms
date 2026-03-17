<template>
  <div class="app-layout">
    <!-- 侧边栏 -->
    <Sidebar :collapsed="sidebarCollapsed" @toggle="toggleSidebar" />
    
    <!-- 主内容区 -->
    <div class="main-container" :style="{ marginLeft: sidebarCollapsed ? '64px' : '240px' }">
      <!-- 头部 -->
      <Header 
        :collapsed="sidebarCollapsed" 
        @toggle-sidebar="toggleSidebar"
      />
      
      <!-- 内容区 -->
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from './Sidebar.vue'
import Header from './Header.vue'

const sidebarCollapsed = ref(false)

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background-color: var(--color-bg-page);
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: margin-left var(--transition-normal);
}

.main-content {
  flex: 1;
  padding: var(--spacing-lg);
  overflow-y: auto;
}

/* 页面过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
