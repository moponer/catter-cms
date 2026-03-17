<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h1>仪表盘</h1>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea, #764ba2)">
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
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb, #f5576c)">
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
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe, #00f2fe)">
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
          <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b, #38f9d7)">
            <el-icon><Picture /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_media }}</div>
            <div class="stat-label">媒体文件</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-card class="quick-actions-card">
      <template #header>
        <div class="card-header">
          <span>快捷操作</span>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col :xs="12" :sm="6" :md="4">
          <div class="quick-action" @click="$router.push('/posts/create')">
            <div class="quick-action-icon" style="background: linear-gradient(135deg, #667eea, #764ba2)">
              <el-icon><Edit /></el-icon>
            </div>
            <div class="quick-action-label">新建文章</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="4">
          <div class="quick-action" @click="$router.push('/pages/create')">
            <div class="quick-action-icon" style="background: linear-gradient(135deg, #f093fb, #f5576c)">
              <el-icon><DocumentAdd /></el-icon>
            </div>
            <div class="quick-action-label">新建页面</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="4">
          <div class="quick-action" @click="$router.push('/media')">
            <div class="quick-action-icon" style="background: linear-gradient(135deg, #4facfe, #00f2fe)">
              <el-icon><Upload /></el-icon>
            </div>
            <div class="quick-action-label">上传媒体</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="4">
          <div class="quick-action" @click="$router.push('/categories')">
            <div class="quick-action-icon" style="background: linear-gradient(135deg, #43e97b, #38f9d7)">
              <el-icon><FolderAdd /></el-icon>
            </div>
            <div class="quick-action-label">添加分类</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="4">
          <div class="quick-action" @click="$router.push('/users')">
            <div class="quick-action-icon" style="background: linear-gradient(135deg, #fa709a, #fee140)">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="quick-action-label">添加用户</div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6" :md="4">
          <div class="quick-action" @click="$router.push('/settings')">
            <div class="quick-action-icon" style="background: linear-gradient(135deg, #a18cd1, #fbc2eb)">
              <el-icon><Setting /></el-icon>
            </div>
            <div class="quick-action-label">系统设置</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 最近文章 -->
    <el-card class="recent-posts-card">
      <template #header>
        <div class="card-header">
          <span>最近文章</span>
          <el-button type="primary" link @click="$router.push('/posts')">
            查看全部
          </el-button>
        </div>
      </template>
      <el-table :data="recentPosts" style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
      </el-table>
    </el-card>
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
  Setting
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

const loadDashboardData = async () => {
  try {
    // 加载统计数据
    const statsRes = await get('/api/dashboard/stats')
    // 处理响应格式: { code: 0, message: "...", data: {...} }
    const statsData = statsRes.data || statsRes
    Object.assign(stats, statsData)

    // 加载最近文章
    const postsRes = await get('/api/dashboard/recent-posts')
    const postsData = postsRes.data || postsRes
    recentPosts.value = postsData || []
  } catch (error) {
    console.error('加载数据失败:', error)
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.page-header h1 {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

/* 统计卡片 */
.stats-row {
  margin-bottom: var(--spacing-md);
}

.stat-card {
  display: flex;
  align-items: center;
  padding: var(--spacing-md);
}

.stat-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  padding: var(--spacing-md);
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
  margin-right: var(--spacing-md);
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin-top: 4px;
}

/* 快捷操作卡片 */
.quick-actions-card {
  margin-bottom: var(--spacing-lg);
}

.quick-action {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-md);
  cursor: pointer;
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}

.quick-action:hover {
  background-color: var(--color-bg-hover);
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
  margin-bottom: var(--spacing-sm);
}

.quick-action-label {
  font-size: 14px;
  color: var(--color-text-primary);
  font-weight: 500;
  text-align: center;
}

/* 最近文章卡片 */
.recent-posts-card {
  flex: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>