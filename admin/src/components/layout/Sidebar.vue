<template>
  <aside class="sidebar" :class="{ collapsed }">
    <!-- Logo -->
    <div class="logo">
      <router-link to="/" class="logo-link">
        <span class="logo-icon">C</span>
        <span v-show="!collapsed" class="logo-text">Catter CMS</span>
      </router-link>
    </div>

    <!-- 导航菜单 -->
    <nav class="sidebar-nav">
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

        <el-menu-item index="/menus">
          <el-icon><Menu /></el-icon>
          <template #title>菜单管理</template>
        </el-menu-item>

        <el-menu-item index="/themes">
          <el-icon><Brush /></el-icon>
          <template #title>主题配置</template>
        </el-menu-item>

        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>

        <el-menu-item index="/ai-settings">
          <el-icon><Connection /></el-icon>
          <template #title>AI 集成</template>
        </el-menu-item>

        <el-menu-item index="/scheduled-tasks">
          <el-icon><Timer /></el-icon>
          <template #title>定时任务</template>
        </el-menu-item>

        <el-menu-item index="/cache">
          <el-icon><Delete /></el-icon>
          <template #title>清理缓存</template>
        </el-menu-item>
      </el-menu>
    </nav>

    <!-- 折叠按钮 -->
    <div class="sidebar-footer">
      <button class="collapse-btn" @click="$emit('toggle')">
        <el-icon>
          <Expand v-if="collapsed" />
          <Fold v-else />
        </el-icon>
        <span v-if="!collapsed" class="collapse-text">收起</span>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
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
  Timer,
  Connection,
  Menu,
  Delete
} from '@element-plus/icons-vue'

defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle'])

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style scoped>
.sidebar {
  width: var(--admin-sidebar-width);
  height: 100vh;
  background-color: var(--admin-sidebar-bg);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  z-index: 1000;
  overflow: hidden;
  transition: width var(--admin-transition-normal);
}

.sidebar.collapsed {
  width: var(--admin-sidebar-collapsed-width);
}

/* Logo */
.logo {
  height: var(--admin-header-height);
  display: flex;
  align-items: center;
  padding: 0 var(--admin-spacing-md);
  border-bottom: 1px solid var(--admin-sidebar-border);
  flex-shrink: 0;
}

.logo-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  gap: var(--admin-spacing-sm);
}

.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--admin-logo-bg);
  border-radius: var(--admin-radius-sm);
  color: var(--admin-logo-text);
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
}

.logo-text {
  color: var(--admin-menu-text-active);
  font-weight: 600;
  font-size: var(--admin-font-size-lg);
  white-space: nowrap;
  overflow: hidden;
}

/* 导航 */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: var(--admin-spacing-sm);
  min-height: 0;
}

/* 菜单 */
.sidebar-menu {
  border: none;
  background: transparent;
  --el-menu-bg-color: transparent;
  --el-menu-text-color: var(--admin-menu-text);
  --el-menu-hover-bg-color: var(--admin-menu-bg-hover);
  --el-menu-active-color: var(--admin-menu-text-active);
  --el-menu-hover-text-color: var(--admin-menu-text-active);
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 100%;
}

/* 菜单项 */
:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 40px;
  line-height: 40px;
  border-radius: var(--admin-radius-sm);
  margin-bottom: 2px;
  color: var(--admin-menu-text);
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: var(--admin-menu-bg-hover) !important;
}

:deep(.el-menu-item.is-active) {
  background-color: var(--admin-menu-bg-active) !important;
  color: var(--admin-menu-text-active) !important;
}

:deep(.el-menu-item.is-active)::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--admin-menu-indicator);
  border-radius: 0 2px 2px 0;
}

:deep(.el-sub-menu .el-menu-item) {
  padding-left: 44px !important;
}

/* 折叠按钮 */
.sidebar-footer {
  padding: var(--admin-spacing-sm);
  border-top: 1px solid var(--admin-sidebar-border);
  flex-shrink: 0;
}

.collapse-btn {
  width: 100%;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--admin-spacing-sm);
  background: transparent;
  border: none;
  border-radius: var(--admin-radius-sm);
  color: var(--admin-menu-text);
  cursor: pointer;
  transition: all var(--admin-transition-fast);
}

.collapse-btn:hover {
  background: var(--admin-menu-bg-hover);
  color: var(--admin-menu-text-active);
}

.collapse-text {
  font-size: var(--admin-font-size-sm);
}

/* 响应式 */
@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar:not(.collapsed) {
    transform: translateX(0);
  }
}
</style>
