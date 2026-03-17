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
          <el-avatar :size="32" :src="userAvatar">
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
  SwitchButton
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store'
import { logout } from '@/utils/auth'

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
const userName = ref('管理员')
const userAvatar = ref('')

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
      router.push('/password')
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
  if (userStore.currentUser) {
    userName.value = userStore.currentUser.name || '管理员'
    userAvatar.value = userStore.currentUser.avatar || ''
  }
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
}
</style>
