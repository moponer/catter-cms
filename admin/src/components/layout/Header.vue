<template>
  <header class="header">
    <div class="header-left">
      <!-- 折叠按钮（移动端） -->
      <button class="toggle-btn" @click="$emit('toggle-sidebar')">
        <el-icon><Expand v-if="collapsed" /><Fold v-else /></el-icon>
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
      <!-- 清缓存 -->
      <el-tooltip content="清理全站缓存，使主题模板修改立即生效" placement="bottom">
        <el-button
          :icon="Refresh"
          :loading="cacheClearing"
          class="clear-cache-btn"
          @click="handleClearCache"
        >
          清缓存
        </el-button>
      </el-tooltip>

      <!-- 搜索 -->
      <el-input
        v-model="searchKeyword"
        placeholder="搜索..."
        prefix-icon="Search"
        class="search-input"
        clearable
        @keyup.enter="handleSearch"
      />

      <!-- 用户菜单 -->
      <el-dropdown trigger="click" @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="32" :src="userAvatar" class="header-avatar">
            {{ userName?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <span class="user-name">{{ userName }}</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <div class="user-dropdown-header">
            <el-avatar :size="44" :src="userAvatar" class="dropdown-avatar">
              {{ userName?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <div class="dropdown-user-info">
              <div class="dropdown-name">{{ userName }}</div>
              <div class="dropdown-email">{{ userEmail }}</div>
            </div>
          </div>
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
  Refresh
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store'
import { logout } from '@/utils/auth'
import { clearAllCache } from '@/api/theme'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle-sidebar'])

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const searchKeyword = ref('')
const userName = computed(() => userStore.currentUser?.name || '管理员')
const userAvatar = computed(() => userStore.currentUser?.avatar || '')
const userEmail = computed(() => userStore.currentUser?.email || '')
const cacheClearing = ref(false)

// 清缓存
const handleClearCache = async () => {
  cacheClearing.value = true
  try {
    await clearAllCache()
    ElMessage.success('缓存已清理，页面将使用最新模板')
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
  height: var(--header-height);
  background-color: var(--color-bg-header);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-lg);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.toggle-btn {
  display: none;
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
}

.toggle-btn:hover {
  background-color: var(--color-bg-hover);
  color: var(--color-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.search-input {
  width: 240px;
}

/* 用户信息 */
.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-md);
  transition: background-color var(--transition-fast);
}

.user-info:hover {
  background-color: var(--color-bg-hover);
}

.user-name {
  font-size: var(--font-size-base);
  color: var(--color-text-primary);
  font-weight: 500;
}

/* 下拉菜单头部 */
.user-dropdown-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px 10px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  min-width: 200px;
}
.dropdown-avatar {
  flex-shrink: 0;
}
.dropdown-user-info {
  overflow: hidden;
}
.dropdown-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.dropdown-email {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-top: 2px;
}

/* 响应式 */
@media (max-width: 768px) {
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

  .clear-cache-btn span:not(.el-icon) {
    display: none;
  }
}
</style>
