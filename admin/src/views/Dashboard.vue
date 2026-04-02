<template>
  <div class="dashboard-page">

    <!-- 欢迎横幅 -->
    <div class="welcome-banner">
      <div class="welcome-left">
        <div class="welcome-title">欢迎回来 👋</div>
        <div class="welcome-sub">今天是 {{ todayStr }}，祝工作顺利</div>
      </div>
      <div class="welcome-right">
        <div class="banner-stat">
          <span class="bs-num">{{ stats.published_posts }}</span>
          <span class="bs-label">已发布</span>
        </div>
        <div class="banner-divider"></div>
        <div class="banner-stat">
          <span class="bs-num">{{ stats.draft_posts }}</span>
          <span class="bs-label">草稿</span>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon-wrap color-primary">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-value">{{ stats.total_posts }}</div>
            <div class="stat-label">文章总数</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon-wrap color-success">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-value">{{ stats.total_users }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon-wrap color-warning">
            <el-icon><Folder /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-value">{{ stats.total_categories }}</div>
            <div class="stat-label">分类总数</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon-wrap color-danger">
            <el-icon><Picture /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-value">{{ stats.total_media }}</div>
            <div class="stat-label">媒体文件</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 主内容区 -->
    <el-row :gutter="16" class="main-row">

      <!-- 左侧：快捷操作 + 最近文章 -->
      <el-col :xs="24" :lg="15">

        <!-- 快捷操作 -->
        <el-card class="section-card quick-card" shadow="never">
          <template #header>
            <span class="card-title">快捷操作</span>
          </template>
          <div class="quick-grid">
            <div
              class="quick-item"
              v-for="action in quickActions"
              :key="action.path"
              @click="$router.push(action.path)"
            >
              <div class="qi-icon" :class="action.iconClass">
                <el-icon><component :is="action.icon" /></el-icon>
              </div>
              <span class="qi-label">{{ action.label }}</span>
              <el-icon class="qi-arrow"><ArrowRight /></el-icon>
            </div>
          </div>
        </el-card>

        <!-- 最近文章 -->
        <el-card class="section-card" shadow="never" style="margin-top:16px" v-loading="loading.posts">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">最近文章</span>
              <el-button type="primary" link @click="$router.push('/posts')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentPosts" style="width:100%" :show-header="true">
            <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip>
              <template #default="{ row }">
                <el-link type="primary" @click="$router.push(`/posts/edit/${row.id}`)">{{ row.title }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="86">
              <template #default="{ row }">
                <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
                  {{ row.status === 'published' ? '已发布' : '草稿' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="155">
              <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>

      </el-col>

      <!-- 右侧：用户信息 + 系统监控 -->
      <el-col :xs="24" :lg="9">

        <!-- 用户信息卡片 -->
        <el-card class="section-card user-card" shadow="never" body-style="padding:0">
          <!-- 顶部色条 -->
          <div class="user-card-banner"></div>
          <!-- 头像 -->
          <div class="user-avatar-wrap">
            <div class="user-avatar-btn">
              <el-avatar :size="70" :src="currentUser.avatar" class="user-avatar">
                <span class="avatar-fallback">{{ currentUser.name ? currentUser.name.charAt(0).toUpperCase() : 'U' }}</span>
              </el-avatar>
            </div>
          </div>
          <!-- 基本信息 -->
          <div class="user-info">
            <div class="user-name">{{ currentUser.name || '未设置昵称' }}</div>
            <div class="user-role-wrap">
              <span class="role-badge" :class="'role-' + (currentUser.role || 'editor')">{{ roleLabel(currentUser.role) }}</span>
            </div>
            <div class="user-email">
              <el-icon :size="12"><Message /></el-icon>
              <span>{{ currentUser.email || '-' }}</span>
            </div>
            <p class="user-bio" v-if="currentUser.bio">{{ currentUser.bio }}</p>
            <p class="user-bio user-bio--empty" v-else>暂无个人简介</p>
          </div>
          <div class="card-sep"></div>
          <!-- 统计 -->
          <div class="stats-block">
            <div class="stat-item">
              <div class="stat-num">{{ stats.published_posts }}</div>
              <div class="stat-label">已发布</div>
            </div>
            <div class="stat-item">
              <div class="stat-num">{{ stats.draft_posts }}</div>
              <div class="stat-label">草稿</div>
            </div>
            <div class="stat-item">
              <div class="stat-num">{{ stats.total_posts }}</div>
              <div class="stat-label">文章总数</div>
            </div>
          </div>
        </el-card>

        <!-- 系统监控 -->
        <el-card class="section-card" shadow="never" style="margin-top:16px" v-loading="loading.system">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">系统监控</span>
              <el-tag size="small" :type="systemInfo.database?.status === 'connected' ? 'success' : 'danger'">
                {{ systemInfo.database?.status === 'connected' ? '运行中' : '异常' }}
              </el-tag>
            </div>
          </template>
          <div class="monitor-list">
            <div class="monitor-row">
              <div class="mr-label"><el-icon><Timer /></el-icon> 运行时间</div>
              <div class="mr-value">{{ systemInfo.uptime || '-' }}</div>
            </div>
            <div class="monitor-row">
              <div class="mr-label"><el-icon><Monitor /></el-icon> 服务器</div>
              <div class="mr-value">{{ systemInfo.server?.os }} {{ systemInfo.server?.arch }}</div>
            </div>
            <div class="monitor-row">
              <div class="mr-label"><el-icon><Coin /></el-icon> 数据库</div>
              <div class="mr-value">{{ systemInfo.database?.type }} {{ systemInfo.database?.version }}</div>
            </div>
            <div class="monitor-row">
              <div class="mr-label"><el-icon><Cpu /></el-icon> CPU 核心</div>
              <div class="mr-value">{{ systemInfo.cpu?.cores || '-' }} 核</div>
            </div>
            <div class="monitor-row">
              <div class="mr-label"><el-icon><Connection /></el-icon> Goroutines</div>
              <div class="mr-value">{{ systemInfo.cpu?.goroutines || '-' }}</div>
            </div>
            <div class="monitor-row memory-row">
              <div class="mr-label"><el-icon><Histogram /></el-icon> 内存使用</div>
              <div class="mr-value">{{ formatBytes(systemInfo.memory?.used) }} / {{ formatBytes(systemInfo.memory?.total) }}</div>
            </div>
            <el-progress
              :percentage="Math.round(systemInfo.memory?.used_percent || 0)"
              :status="getMemoryStatus(systemInfo.memory?.used_percent)"
              :stroke-width="6"
              style="margin-bottom:12px"
            />
          </div>
        </el-card>

      </el-col>
    </el-row>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import {
  Document, User, Folder, Picture,
  Edit, DocumentAdd, Upload, FolderAdd,
  Setting, Timer, Monitor, Coin, Cpu,
  Connection, Histogram, PriceTag, ChatDotRound,
  Brush, ArrowRight, Message, UserFilled, View
} from '@element-plus/icons-vue'
import { get } from '@/utils/request'

const stats = reactive({
  total_posts: 0, published_posts: 0, draft_posts: 0,
  total_users: 0, total_categories: 0, total_tags: 0, total_media: 0
})
const recentPosts = ref([])
const systemInfo = ref({})
const loading = reactive({ posts: false, system: false })

// 当前用户信息
const currentUser = reactive({
  name: '',
  email: '',
  role: '',
  avatar: '',
  bio: ''
})

// 今日日期
const todayStr = computed(() => {
  return new Date().toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' })
})

// 快捷操作
const quickActions = [
  { label: '新建文章',  icon: 'Edit',         iconClass: 'ic-primary',  path: '/posts/create' },
  { label: '新建页面',  icon: 'DocumentAdd',  iconClass: 'ic-success',  path: '/pages/create' },
  { label: '上传媒体',  icon: 'Upload',       iconClass: 'ic-warning',  path: '/media' },
  { label: '添加分类',  icon: 'FolderAdd',    iconClass: 'ic-danger',   path: '/categories' },
  { label: '添加标签',  icon: 'PriceTag',     iconClass: 'ic-info',     path: '/tags' },
  { label: '评论管理',  icon: 'ChatDotRound', iconClass: 'ic-warning2', path: '/comments' },
  { label: '主题管理',  icon: 'Brush',        iconClass: 'ic-purple',   path: '/themes' },
  { label: '系统设置',  icon: 'Setting',      iconClass: 'ic-neutral',  path: '/settings' }
]

const loadDashboardData = async () => {
  try {
    const statsRes = await get('/api/dashboard/stats')
    Object.assign(stats, statsRes.data || statsRes)
    loading.posts = true
    const postsRes = await get('/api/dashboard/recent-posts')
    recentPosts.value = (postsRes.data || postsRes) || []
  } catch {}
  finally { loading.posts = false }
}

const loadSystemInfo = async () => {
  try {
    loading.system = true
    const res = await get('/api/admin/system/info')
    systemInfo.value = res.data || res || {}
  } catch {}
  finally { loading.system = false }
}

// 加载当前用户信息
const loadCurrentUser = async () => {
  try {
    const res = await get('/api/auth/me')
    const data = res.data || res
    if (data) {
      Object.assign(currentUser, data)
    }
  } catch {}
}

// 角色标签
const roleLabel = (role) => {
  const map = { admin: '管理员', editor: '编辑', author: '作者', subscriber: '订阅者' }
  return map[role] || role || '未知'
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN', { year:'numeric', month:'2-digit', day:'2-digit', hour:'2-digit', minute:'2-digit' })
}

const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024, sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const getMemoryStatus = (p) => {
  if (!p) return ''
  if (p > 90) return 'exception'
  if (p > 70) return 'warning'
  return 'success'
}

onMounted(() => { loadDashboardData(); loadSystemInfo(); loadCurrentUser() })
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ===== 欢迎横幅 ===== */
.welcome-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  background-color: var(--admin-welcome-bg, #374151);
  border-radius: 10px;
  color: #fff;
}

.welcome-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.welcome-sub {
  font-size: 13px;
  opacity: 0.85;
}

.welcome-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.banner-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.bs-num {
  font-size: 22px;
  font-weight: 700;
  line-height: 1.2;
}

.bs-label {
  font-size: 12px;
  opacity: 0.8;
  margin-top: 2px;
}

.banner-divider {
  width: 1px;
  height: 32px;
  background: rgba(255,255,255,0.3);
}

/* ===== 统计卡片 ===== */
.stats-row {
  margin: 0 !important;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 16px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
  margin-bottom: 0;
  transition: box-shadow 0.2s;
}

.stat-card:hover {
  box-shadow: var(--el-box-shadow-light);
}

.stat-icon-wrap {
  width: 46px;
  height: 46px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  flex-shrink: 0;
}

.color-primary { background: var(--el-color-primary-light-9); color: var(--el-color-primary); }
.color-success { background: var(--el-color-success-light-9); color: var(--el-color-success); }
.color-warning { background: var(--el-color-warning-light-9); color: var(--el-color-warning); }
.color-danger  { background: var(--el-color-danger-light-9);  color: var(--el-color-danger);  }

.stat-body { flex: 1; min-width: 0; }

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin-top: 3px;
}

/* ===== 主内容区 ===== */
.main-row {
  margin: 0 !important;
  align-items: flex-start;
}

.main-row .el-col {
  margin-bottom: 0;
}

.section-card {
  border-radius: 10px;
  border: 1px solid var(--el-border-color-lighter) !important;
}

.section-card :deep(.el-card__header) {
  padding: 14px 18px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.section-card :deep(.el-card__body) {
  padding: 16px 18px;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.card-header-flex {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* ===== 快捷操作 ===== */
.quick-card :deep(.el-card__body) {
  padding: 8px 10px;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 2px;
}

.quick-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.quick-item:hover {
  background: var(--el-fill-color-light);
}

.qi-icon {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

/* 快捷操作图标颜色 — 明确的灰色系 */
.ic-primary  { background: #F3F4F6;  color: #374151;  }
.ic-success  { background: #DCFCE7;  color: #16A34A;  }
.ic-warning  { background: #FEF3C7;  color: #D97706;  }
.ic-danger   { background: #FEE2E2; color: #DC2626;  }
.ic-info     { background: #F3F4F6;   color: #6B7280;  }
.ic-warning2 { background: #FEF3C7;   color: #D97706;  }
.ic-purple   { background: #F3F4F6;   color: #374151;  }
.ic-neutral  { background: #F3F4F6;  color: #374151;  }

.qi-label {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.qi-arrow {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}

/* ===== 用户信息卡片（Profile 风格） ===== */
.user-card {
  overflow: hidden;
  padding: 0;
}

/* 顶部色条 */
.user-card-banner {
  height: 56px;
  background-color: var(--el-color-primary);
  opacity: 0.12;
}

/* 头像 */
.user-avatar-wrap {
  display: flex;
  justify-content: center;
  margin-top: -36px;
  margin-bottom: 4px;
}
.user-avatar-btn {
  position: relative;
  width: 76px;
  height: 76px;
  flex-shrink: 0;
}
.user-card :deep(.user-avatar) {
  width: 70px !important;
  height: 70px !important;
  border: 3px solid var(--el-bg-color);
  box-shadow: 0 2px 12px rgba(0,0,0,0.12);
  font-size: 26px;
}
.avatar-fallback { font-size: 26px; font-weight: 700; }

/* 用户信息 */
.user-info {
  padding: 0 16px 8px;
  text-align: center;
}
.user-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.user-role-wrap { margin-bottom: 8px; }
.role-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
}
.role-admin {
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
  border: 1px solid var(--el-color-primary-light-7);
}
.role-editor {
  background: var(--el-color-success-light-9);
  color: var(--el-color-success);
  border: 1px solid var(--el-color-success-light-7);
}
.user-email {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 8px;
  max-width: 100%;
  overflow: hidden;
}
.user-email span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.user-bio {
  font-size: 12px;
  color: var(--el-text-color-regular);
  line-height: 1.6;
  margin: 4px 0 0;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}
.user-bio--empty { color: var(--el-text-color-placeholder); font-style: italic; }

/* 分割线 */
.card-sep {
  height: 1px;
  background: var(--el-border-color-lighter);
  margin: 0 0 4px;
}

/* 统计块 */
.stats-block {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  padding: 8px 12px 12px;
  gap: 4px;
}
.stat-item {
  text-align: center;
  padding: 8px 4px;
  border-radius: 8px;
  transition: background 0.15s;
}
.stat-item:hover { background: var(--el-fill-color-light); }
.stat-num {
  font-size: 20px;
  font-weight: 800;
  color: var(--el-color-primary);
  line-height: 1;
}
.stat-label {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}

/* ===== 系统监控 ===== */
.monitor-list { padding: 4px 0; }

.monitor-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--el-border-color-extra-light);
  font-size: 13px;
}

.monitor-row:last-of-type { border-bottom: none; }

.memory-row { margin-bottom: 4px; }

.mr-label {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--el-text-color-regular);
}

.mr-label .el-icon {
  font-size: 15px;
  color: var(--el-text-color-secondary);
}

.mr-value {
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

/* ===== 响应式 ===== */
@media (max-width: 992px) {
  .main-row .el-col + .el-col {
    margin-top: 16px;
  }
}

@media (max-width: 768px) {
  .welcome-banner {
    flex-direction: column;
    align-items: flex-start;
    gap: 14px;
  }
  .quick-grid {
    grid-template-columns: 1fr;
  }
  .stats-row .el-col {
    margin-bottom: 12px;
  }
}
</style>
