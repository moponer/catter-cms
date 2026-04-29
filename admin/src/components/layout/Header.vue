<template>
  <header class="header">
    <div class="header-left">
      <!-- 折叠按钮 -->
      <button class="toggle-btn" @click="$emit('toggle-sidebar')">
        <el-icon v-if="collapsed"><Expand /></el-icon>
        <el-icon v-else><Fold /></el-icon>
      </button>

      <!-- 面包屑 -->
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="header-right">
      <!-- 主题切换 - VS Code 风格 -->
      <el-tooltip :content="currentTheme.id === 'dark' ? '切换到浅色模式' : '切换到深色模式'" placement="bottom">
        <button class="theme-toggle" @click="toggleTheme">
          <el-icon size="18">
            <Moon v-if="currentTheme.id === 'dark'" />
            <Sunny v-else />
          </el-icon>
        </button>
      </el-tooltip>

      <!-- 清缓存 -->
      <el-tooltip content="清理缓存" placement="bottom">
        <el-button
          :icon="Refresh"
          :loading="cacheClearing"
          class="icon-btn"
          @click="handleClearCache"
        />
      </el-tooltip>

      <!-- 搜索 -->
      <el-input
        v-model="searchKeyword"
        placeholder="搜索..."
        :prefix-icon="Search"
        class="search-input"
        clearable
        @keyup.enter="handleSearch"
      />

      <!-- 用户菜单 -->
      <el-dropdown trigger="click" @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="32" class="header-avatar">
            {{ userName?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <span class="user-name">{{ userName }}</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item command="password">
              <el-icon><Key /></el-icon>
              修改密码
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  Expand,
  Fold,
  ArrowDown,
  User,
  Key,
  SwitchButton,
  Refresh,
  Search,
  Sunny,
  Moon
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store'
import { logout } from '@/utils/auth'
import { clearAllCache } from '@/api/theme'
import { useAdminTheme } from '@/composables/useAdminTheme'

defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle-sidebar'])

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 主题系统
const { currentTheme, setTheme, themes } = useAdminTheme()

const searchKeyword = ref('')
const userName = computed(() => userStore.currentUser?.name || '管理员')
const cacheClearing = ref(false)

// 切换主题
const toggleTheme = () => {
  const newTheme = currentTheme.value.id === 'dark' ? 'light' : 'dark'
  setTheme(newTheme)
}

// 清缓存
const handleClearCache = async () => {
  cacheClearing.value = true
  try {
    await clearAllCache()
    ElMessage.success('缓存已清理')
  } catch (err) {
    ElMessage.error('清缓存失败：' + (err?.message || '未知错误'))
  } finally {
    cacheClearing.value = false
  }
}

// 面包屑导航
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    path: item.path,
    title: item.meta.title
  }))
})

// 处理搜索
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    ElMessage.info(`搜索: ${searchKeyword.value}`)
  }
}

// 处理下拉菜单命令
const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'password':
      router.push('/change-password')
      break
    case 'logout':
      handleLogout()
      break
  }
}

// 处理退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    logout()
  } catch {
    // 用户取消
  }
}

onMounted(async () => {
  await userStore.getCurrentUser()
})
</script>

<style scoped>
.header {
  height: var(--admin-header-height);
  background-color: var(--admin-header-bg);
  border-bottom: 1px solid var(--admin-header-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--admin-spacing-lg);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--admin-spacing-md);
}

.toggle-btn {
  display: none;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: var(--admin-radius-sm);
  cursor: pointer;
  color: var(--admin-text-secondary);
  transition: all var(--admin-transition-fast);
}

.toggle-btn:hover {
  background-color: var(--admin-bg-hover);
  color: var(--admin-text-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--admin-spacing-md);
}

/* 主题切换按钮 */
.theme-toggle {
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: var(--admin-radius-sm);
  cursor: pointer;
  color: var(--admin-text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--admin-transition-fast);
}

.theme-toggle:hover {
  background-color: var(--admin-bg-hover);
  color: var(--admin-text-primary);
}

/* 图标按钮 */
.icon-btn {
  width: 36px;
  height: 36px;
  padding: 0;
  border: none;
  background: transparent;
  color: var(--admin-text-secondary);
}

.icon-btn:hover {
  background-color: var(--admin-bg-hover);
  color: var(--admin-text-primary);
}

/* 搜索框 */
.search-input {
  width: 200px;
}

/* 用户信息 */
.user-info {
  display: flex;
  align-items: center;
  gap: var(--admin-spacing-sm);
  cursor: pointer;
  padding: var(--admin-spacing-xs) var(--admin-spacing-sm);
  border-radius: var(--admin-radius-sm);
  transition: background-color var(--admin-transition-fast);
}

.user-info:hover {
  background-color: var(--admin-bg-hover);
}

.user-name {
  font-size: var(--admin-font-size-base);
  color: var(--admin-text-primary);
  font-weight: 500;
}

.header-avatar {
  background: var(--admin-primary);
  color: white;
  font-weight: 500;
}

/* 响应式 */
@media (max-width: 768px) {
  .header {
    padding: 0 var(--admin-spacing-sm);
  }

  .toggle-btn {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .search-input {
    display: none;
  }

  .user-name {
    display: none;
  }
}
</style>
