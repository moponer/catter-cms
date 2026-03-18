<template>
  <aside class="sidebar" :class="{ collapsed }">
    <!-- Logo -->
    <div class="logo">
      <router-link to="/" class="logo-link">
        <span class="logo-icon">C</span>
        <span v-show="!collapsed" class="logo-text">Catter CMS</span>
      </router-link>
    </div>

    <!-- 快速搜索 -->
    <div class="quick-search" v-show="!collapsed">
      <el-input
        v-model="searchQuery"
        placeholder="全局搜索 (Ctrl+K)"
        :prefix-icon="Search"
        clearable
        @focus="openGlobalSearch"
        @keyup.enter="openGlobalSearch"
      />
    </div>

    <!-- 导航菜单 -->
    <el-menu
      :default-active="activeMenu"
      :collapse="collapsed"
      :collapse-transition="false"
      class="sidebar-menu"
      router
    >
      <el-menu-item index="/">
        <el-icon><House /></el-icon>
        <template #title>仪表盘</template>
      </el-menu-item>

      <el-sub-menu index="/content">
        <template #title>
          <el-icon><Document /></el-icon>
          <span>内容管理</span>
        </template>
        <el-menu-item index="/posts">文章列表</el-menu-item>
        <el-menu-item index="/pages">页面管理</el-menu-item>
        <el-menu-item index="/categories">分类管理</el-menu-item>
        <el-menu-item index="/tags">标签管理</el-menu-item>
        <el-menu-item index="/comments">评论管理</el-menu-item>
      </el-sub-menu>

      <el-menu-item index="/users">
        <el-icon><User /></el-icon>
        <template #title>用户管理</template>
      </el-menu-item>

      <el-menu-item index="/media">
        <el-icon><Picture /></el-icon>
        <template #title>媒体库</template>
      </el-menu-item>

      <el-menu-item index="/themes">
        <el-icon><Brush /></el-icon>
        <template #title>主题配置</template>
      </el-menu-item>

      <el-menu-item index="/settings">
        <el-icon><Setting /></el-icon>
        <template #title>系统设置</template>
      </el-menu-item>

      <el-menu-item index="/scheduled-tasks">
        <el-icon><Timer /></el-icon>
        <template #title>定时任务</template>
      </el-menu-item>

      <el-menu-item index="/sitemap">
        <el-icon><Connection /></el-icon>
        <template #title>Sitemap</template>
      </el-menu-item>
    </el-menu>

    <!-- 折叠按钮 -->
    <div class="sidebar-footer">
      <button class="collapse-btn" @click="$emit('toggle')">
        <el-icon>
          <Expand v-if="collapsed" />
          <Fold v-else />
        </el-icon>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  House,
  Document,
  User,
  Picture,
  Brush,
  Setting,
  Expand,
  Fold,
  Search,
  Timer,
  Connection
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle'])

const route = useRoute()
const activeMenu = computed(() => route.path)

// 搜索状态
const searchQuery = ref('')

// 打开全局搜索
const openGlobalSearch = () => {
  // 触发全局搜索事件,通过自定义事件通知父组件
  const event = new CustomEvent('open-global-search', {
    detail: { query: searchQuery.value }
  })
  window.dispatchEvent(event)
  searchQuery.value = ''
}
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: 240px;
  background-color: var(--color-bg-sidebar);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-normal);
  z-index: 1000;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 64px;
}

/* Logo */
.logo {
  height: var(--header-height);
  display: flex;
  align-items: center;
  padding: 0 var(--spacing-md);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: var(--color-text-inverse);
  font-weight: 600;
  font-size: var(--font-size-lg);
}

.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--color-primary), var(--color-primary-dark));
  border-radius: var(--radius-md);
  color: white;
  font-weight: 700;
  font-size: 18px;
  flex-shrink: 0;
}

.logo-text {
  margin-left: var(--spacing-sm);
  white-space: nowrap;
  overflow: hidden;
}

/* 快速搜索 */
.quick-search {
  padding: var(--spacing-md);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.quick-search :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: none;
  transition: all var(--transition-fast);
}

.quick-search :deep(.el-input__wrapper:hover),
.quick-search :deep(.el-input__wrapper.is-focus) {
  background: rgba(255, 255, 255, 0.15);
  border-color: var(--color-primary-light);
}

.quick-search :deep(.el-input__inner) {
  color: rgba(255, 255, 255, 0.9);
}

.quick-search :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.5);
}

.quick-search :deep(.el-input__icon) {
  color: rgba(255, 255, 255, 0.7);
}

/* 菜单 */
.sidebar-menu {
  flex: 1;
  border-right: none;
  background-color: transparent;
  --el-menu-bg-color: transparent;
  --el-menu-text-color: rgba(255, 255, 255, 0.7);
  --el-menu-hover-bg-color: rgba(255, 255, 255, 0.1);
  --el-menu-active-color: var(--color-primary-light);
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 240px;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 48px;
  line-height: 48px;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

:deep(.el-menu-item.is-active) {
  background-color: rgba(59, 130, 246, 0.2) !important;
  color: var(--color-primary-light) !important;
}

:deep(.el-sub-menu .el-menu-item) {
  padding-left: 50px !important;
}

/* 底部折叠按钮 */
.sidebar-footer {
  padding: var(--spacing-md);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.collapse-btn {
  width: 100%;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: var(--radius-md);
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}
</style>
