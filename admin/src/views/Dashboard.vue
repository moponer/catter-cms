<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h1>仪表盘</h1>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-icon articles">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_posts }}</div>
            <div class="stat-label">文章总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-icon users">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_users }}</div>
            <div class="stat-label">用户总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-icon categories">
            <el-icon><Folder /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_categories }}</div>
            <div class="stat-label">分类总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-icon media">
            <el-icon><Picture /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_media }}</div>
            <div class="stat-label">媒体文件</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作和系统介绍 - 使用14:10比例与底部对齐 -->
    <el-row :gutter="20" class="section-row">
      <el-col :xs="24" :lg="14">
        <el-card class="quick-actions-card">
          <template #header>
            <div class="card-header">
              <span>快捷操作</span>
            </div>
          </template>
          <el-row :gutter="16">
            <el-col :xs="12" :sm="8" :md="6" v-for="action in quickActions" :key="action.path">
              <div class="quick-action" @click="$router.push(action.path)">
                <div class="quick-action-icon" :class="action.iconClass">
                  <el-icon>
                    <component :is="action.icon" />
                  </el-icon>
                </div>
                <div class="quick-action-label">{{ action.label }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="10">
        <el-card class="system-intro-card">
          <template #header>
            <div class="card-header">
              <span>系统介绍</span>
            </div>
          </template>
          <div class="intro-content">
            <div class="intro-logo">
              <el-avatar :size="64" class="logo-avatar">C</el-avatar>
              <h3>Catter CMS</h3>
              <p class="version">v{{ systemInfo.version || '2.0.0' }}</p>
            </div>
            <p class="intro-text">
              Catter CMS 是一个现代化的内容管理系统，基于 Go + Vue 构建。
              提供文章管理、页面管理、媒体管理、主题管理等核心功能。
            </p>
            <div class="intro-links">
              <el-link type="primary" href="https://github.com" target="_blank">
                <el-icon><Link /></el-icon> GitHub
              </el-link>
              <el-link type="primary" href="#" target="_blank">
                <el-icon><Document /></el-icon> 文档
              </el-link>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近文章和系统监控 - 使用14:10比例 -->
    <el-row :gutter="20" class="section-row">
      <el-col :xs="24" :lg="14">
        <el-card class="recent-posts-card">
          <template #header>
            <div class="card-header">
              <span>最近文章</span>
              <el-button type="primary" link @click="$router.push('/posts')">
                查看全部
              </el-button>
            </div>
          </template>
          <el-table :data="recentPosts" style="width: 100%" v-loading="loading.posts">
            <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip>
              <template #default="{ row }">
                <el-link type="primary" @click="$router.push(`/posts/edit/${row.id}`)">
                  {{ row.title }}
                </el-link>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
                  {{ row.status === 'published' ? '已发布' : '草稿' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="10">
        <el-card class="system-monitor-card">
          <template #header>
            <div class="card-header">
              <span>系统监控</span>
              <el-tag size="small" :type="systemInfo.database?.status === 'connected' ? 'success' : 'danger'">
                {{ systemInfo.database?.status === 'connected' ? '运行中' : '异常' }}
              </el-tag>
            </div>
          </template>
          <div class="monitor-content" v-loading="loading.system">
            <!-- 运行时间 -->
            <div class="monitor-item">
              <div class="monitor-label">
                <el-icon><Timer /></el-icon>
                <span>运行时间</span>
              </div>
              <div class="monitor-value">{{ systemInfo.uptime || '-' }}</div>
            </div>

            <!-- 服务器信息 -->
            <div class="monitor-item">
              <div class="monitor-label">
                <el-icon><Monitor /></el-icon>
                <span>服务器</span>
              </div>
              <div class="monitor-value">
                {{ systemInfo.server?.os }} {{ systemInfo.server?.arch }}
              </div>
            </div>

            <!-- 数据库 -->
            <div class="monitor-item">
              <div class="monitor-label">
                <el-icon><Coin /></el-icon>
                <span>数据库</span>
              </div>
              <div class="monitor-value">
                {{ systemInfo.database?.type }} {{ systemInfo.database?.version }}
              </div>
            </div>

            <!-- CPU -->
            <div class="monitor-item">
              <div class="monitor-label">
                <el-icon><Cpu /></el-icon>
                <span>CPU 核心</span>
              </div>
              <div class="monitor-value">{{ systemInfo.cpu?.cores || '-' }} 核</div>
            </div>

            <!-- Goroutines -->
            <div class="monitor-item">
              <div class="monitor-label">
                <el-icon><Connection /></el-icon>
                <span>Goroutines</span>
              </div>
              <div class="monitor-value">{{ systemInfo.cpu?.goroutines || '-' }}</div>
            </div>

            <!-- 内存使用 -->
            <div class="monitor-item memory-item">
              <div class="monitor-label">
                <el-icon><Histogram /></el-icon>
                <span>内存使用</span>
              </div>
              <div class="monitor-value">
                {{ formatBytes(systemInfo.memory?.used) }} / {{ formatBytes(systemInfo.memory?.total) }}
              </div>
            </div>
            <el-progress 
              :percentage="Math.round(systemInfo.memory?.used_percent || 0)" 
              :status="getMemoryStatus(systemInfo.memory?.used_percent)"
            />

            <!-- CMS版本 -->
            <div class="monitor-item version-item">
              <div class="monitor-label">
                <el-icon><InfoFilled /></el-icon>
                <span>Catter CMS</span>
              </div>
              <div class="monitor-value version">v{{ systemInfo.version || '2.0.0' }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import {
  Document,
  User,
  Folder,
  Picture,
  Edit,
  DocumentAdd,
  Upload,
  FolderAdd,
  UserFilled,
  Setting,
  Link,
  Timer,
  Monitor,
  Coin,
  Cpu,
  Connection,
  Histogram,
  InfoFilled,
  Collection,
  PriceTag,
  ChatDotRound,
  Brush
} from '@element-plus/icons-vue'
import { get } from '@/utils/request'

const stats = reactive({
  total_posts: 0,
  published_posts: 0,
  draft_posts: 0,
  total_users: 0,
  total_categories: 0,
  total_tags: 0,
  total_media: 0
})

const recentPosts = ref([])
const systemInfo = ref({})

const loading = reactive({
  posts: false,
  system: false
})

// 快捷操作配置 - 使用纯色图标替代渐变色
const quickActions = [
  { label: '新建文章', icon: 'Edit', iconClass: 'articles', path: '/posts/create' },
  { label: '新建页面', icon: 'DocumentAdd', iconClass: 'pages', path: '/pages/create' },
  { label: '上传媒体', icon: 'Upload', iconClass: 'media', path: '/media' },
  { label: '添加分类', icon: 'FolderAdd', iconClass: 'categories', path: '/categories' },
  { label: '添加标签', icon: 'PriceTag', iconClass: 'tags', path: '/tags' },
  { label: '评论管理', icon: 'ChatDotRound', iconClass: 'comments', path: '/comments' },
  { label: '主题管理', icon: 'Brush', iconClass: 'themes', path: '/themes' },
  { label: '系统设置', icon: 'Setting', iconClass: 'settings', path: '/settings' }
]

const loadDashboardData = async () => {
  try {
    // 加载统计数据
    const statsRes = await get('/api/dashboard/stats')
    const statsData = statsRes.data || statsRes
    Object.assign(stats, statsData)

    // 加载最近文章
    loading.posts = true
    const postsRes = await get('/api/dashboard/recent-posts')
    const postsData = postsRes.data || postsRes
    recentPosts.value = postsData || []
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.posts = false
  }
}

const loadSystemInfo = async () => {
  try {
    loading.system = true
    const res = await get('/api/admin/system/info')
    systemInfo.value = res.data || res || {}
  } catch (error) {
    console.error('加载系统信息失败:', error)
  } finally {
    loading.system = false
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 格式化字节
const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 获取内存状态
const getMemoryStatus = (percent) => {
  if (!percent) return ''
  if (percent > 90) return 'exception'
  if (percent > 70) return 'warning'
  return 'success'
}

onMounted(() => {
  loadDashboardData()
  loadSystemInfo()
})
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0;
}

/* 统计卡片 */
.stats-row {
  margin-bottom: 10px;
}

.stat-card {
  display: flex;
  align-items: center;
}

.stat-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  padding: 16px;
  width: 100%;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  margin-right: 12px;
  flex-shrink: 0;
}

/* 统计图标 - 使用纯色替代渐变色 */
.stat-icon.articles {
  background-color: #409eff;
}

.stat-icon.users {
  background-color: #67c23a;
}

.stat-icon.categories {
  background-color: #e6a23c;
}

.stat-icon.media {
  background-color: #f56c6c;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}

/* 分栏布局 - 统一使用14:10比例 */
.section-row {
  margin: 0 !important;
}

.section-row .el-col {
  margin-bottom: 20px;
}

/* 快捷操作卡片 */
.quick-actions-card,
.system-intro-card,
.recent-posts-card,
.system-monitor-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.quick-action {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 8px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.3s;
  margin-bottom: 12px;
}

.quick-action:hover {
  background-color: var(--el-fill-color-light);
  transform: translateY(-2px);
}

.quick-action-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  margin-bottom: 8px;
  transition: all 0.3s;
}

/* 快捷操作图标 - 使用纯色替代渐变色 */
.quick-action-icon.articles {
  background-color: #409eff;
}

.quick-action-icon.pages {
  background-color: #67c23a;
}

.quick-action-icon.media {
  background-color: #e6a23c;
}

.quick-action-icon.categories {
  background-color: #f56c6c;
}

.quick-action-icon.tags {
  background-color: #909399;
}

.quick-action-icon.comments {
  background-color: #b88230;
}

.quick-action-icon.themes {
  background-color: #8e44ad;
}

.quick-action-icon.settings {
  background-color: #606266;
}

.quick-action-label {
  font-size: 14px;
  color: var(--el-text-color-primary);
  font-weight: 500;
  text-align: center;
}

/* 系统介绍卡片 */
.intro-content {
  text-align: center;
}

.intro-logo {
  margin-bottom: 16px;
}

.logo-avatar {
  background: linear-gradient(135deg, #667eea, #764ba2);
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 8px;
}

.intro-logo h3 {
  margin: 8px 0 4px;
  font-size: 18px;
  color: var(--el-text-color-primary);
}

.version {
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin: 0;
}

.intro-text {
  color: var(--el-text-color-regular);
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 16px;
}

.intro-links {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.intro-links .el-link {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 系统监控卡片 */
.monitor-content {
  padding: 8px 0;
}

.monitor-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.monitor-item:last-child {
  border-bottom: none;
}

.monitor-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.monitor-label .el-icon {
  font-size: 16px;
  color: var(--el-text-color-secondary);
}

.monitor-value {
  font-size: 14px;
  color: var(--el-text-color-primary);
  font-weight: 500;
}

.monitor-value.version {
  color: var(--el-color-primary);
}

.memory-item {
  margin-bottom: 8px;
}

.version-item {
  margin-top: 8px;
  padding-top: 12px;
  border-top: 1px solid var(--el-border-color-lighter);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .stat-card :deep(.el-card__body) {
    flex-direction: column;
    text-align: center;
  }

  .stat-icon {
    margin-right: 0;
    margin-bottom: 8px;
  }

  .quick-action {
    padding: 12px 4px;
  }

  .quick-action-icon {
    width: 40px;
    height: 40px;
    font-size: 20px;
  }
}
</style>
