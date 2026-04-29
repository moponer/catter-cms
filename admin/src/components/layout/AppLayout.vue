<template>
  <div class="app-layout">
    <!-- 移动端遮罩 -->
    <div
      v-if="isMobile && !sidebarCollapsed"
      class="sidebar-overlay"
      @click="sidebarCollapsed = true"
    />

    <!-- 侧边栏 -->
    <Sidebar :collapsed="sidebarCollapsed" @toggle="toggleSidebar" />

    <!-- 主内容区 -->
    <div class="main-container" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import Sidebar from './Sidebar.vue'
import Header from './Header.vue'

const sidebarCollapsed = ref(false)
const windowWidth = ref(window.innerWidth)

// 移动端判断
const isMobile = computed(() => windowWidth.value < 768)

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

// 监听窗口大小变化
const handleResize = () => {
  windowWidth.value = window.innerWidth
  // 切换到移动端时自动折叠侧边栏
  if (isMobile.value && !sidebarCollapsed.value) {
    sidebarCollapsed.value = true
  }
  // 切换到桌面端时自动展开侧边栏
  if (!isMobile.value && sidebarCollapsed.value) {
    sidebarCollapsed.value = false
  }
}

onMounted(() => {
  // 初始化：移动端默认折叠
  if (isMobile.value) {
    sidebarCollapsed.value = true
  }
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background-color: var(--admin-bg-page);
}

/* 遮罩层 */
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.main-content {
  flex: 1;
  padding: var(--admin-spacing-lg);
  padding-bottom: 0;
  overflow-y: auto;
  min-height: 0;
  box-sizing: border-box;
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

/* 响应式 */
@media (max-width: 768px) {
  .main-content {
    padding: var(--admin-spacing-sm);
    padding-bottom: 0;
  }
}
</style>
