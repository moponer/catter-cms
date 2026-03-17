<template>
  <div class="comments-page">
    <div class="page-header">
      <h1>评论管理</h1>
      <el-button type="primary" @click="refreshList">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 搜索筛选 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="搜索">
          <el-input
            v-model="filters.keyword"
            placeholder="搜索评论内容或作者"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button :icon="Search" @click="handleSearch" />
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable @change="handleSearch">
            <el-option label="待审核" value="pending" />
            <el-option label="已批准" value="approved" />
            <el-option label="垃圾评论" value="spam" />
          </el-select>
        </el-form-item>
        <el-form-item label="文章">
          <el-select v-model="filters.post_id" placeholder="选择文章" clearable filterable @change="handleSearch">
            <el-option
              v-for="post in posts"
              :key="post.id"
              :label="post.title"
              :value="post.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 评论列表 -->
    <el-card class="comments-card">
      <template #header>
        <div class="card-header">
          <span>评论列表</span>
          <div class="header-actions">
            <el-checkbox v-model="selectAll" @change="handleSelectAll">全选</el-checkbox>
            <el-button
              type="success"
              size="small"
              :disabled="selectedComments.length === 0"
              @click="handleBatchApprove"
            >
              批量批准
            </el-button>
            <el-button
              type="danger"
              size="small"
              :disabled="selectedComments.length === 0"
              @click="handleBatchDelete"
            >
              批量删除
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="comments-list">
        <el-empty v-if="!loading && comments.length === 0" description="暂无评论" />

        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <el-checkbox v-model="selectedComments" :value="comment.id" class="comment-checkbox" />
          <div class="comment-content">
            <div class="comment-header">
              <div class="comment-author">
                <el-avatar :size="32" class="author-avatar">
                  {{ comment.author_name?.charAt(0).toUpperCase() }}
                </el-avatar>
                <div class="author-info">
                  <span class="author-name">{{ comment.author_name }}</span>
                  <el-tag :type="getStatusType(comment.status)" size="small" class="status-tag">
                    {{ getStatusText(comment.status) }}
                  </el-tag>
                </div>
              </div>
              <div class="comment-meta">
                <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
                <span class="comment-post">{{ comment.post?.title }}</span>
              </div>
            </div>

            <div class="comment-text">{{ comment.content }}</div>

            <div class="comment-actions">
              <el-button
                v-if="comment.status === 'pending'"
                type="success"
                size="small"
                link
                @click="handleApprove(comment.id)"
              >
                <el-icon><Select /></el-icon>
                批准
              </el-button>
              <el-button
                v-if="comment.status === 'pending'"
                type="danger"
                size="small"
                link
                @click="handleReject(comment.id)"
              >
                <el-icon><Close /></el-icon>
                拒绝
              </el-button>
              <el-button type="primary" size="small" link @click="handleEdit(comment)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button type="danger" size="small" link @click="handleDelete(comment.id)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>

            <!-- 回复列表 -->
            <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
              <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
                <el-avatar :size="24" class="reply-avatar">
                  {{ reply.author_name?.charAt(0).toUpperCase() }}
                </el-avatar>
                <div class="reply-content">
                  <span class="reply-author">{{ reply.author_name }}</span>
                  <span class="reply-text">{{ reply.content }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="loadComments"
          @size-change="loadComments"
        />
      </div>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑评论" width="600px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="作者姓名">
          <el-input v-model="editForm.author_name" />
        </el-form-item>
        <el-form-item label="作者邮箱">
          <el-input v-model="editForm.author_email" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status">
            <el-option label="待审核" value="pending" />
            <el-option label="已批准" value="approved" />
            <el-option label="垃圾评论" value="spam" />
          </el-select>
        </el-form-item>
        <el-form-item label="评论内容">
          <el-input
            v-model="editForm.content"
            type="textarea"
            :rows="4"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh,
  Search,
  Select,
  Close,
  Edit,
  Delete
} from '@element-plus/icons-vue'
import { get, post, put, del } from '@/utils/request'

const loading = ref(false)
const comments = ref([])
const posts = ref([])
const selectedComments = ref([])
const selectAll = ref(false)
const editDialogVisible = ref(false)
const editForm = reactive({
  id: '',
  author_name: '',
  author_email: '',
  status: '',
  content: ''
})

const filters = reactive({
  keyword: '',
  status: '',
  post_id: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const loadPosts = async () => {
  try {
    const res = await get('/api/posts', { page_size: 100 })
    posts.value = res.data || res
  } catch (error) {
    console.error('加载文章列表失败:', error)
  }
}

const loadComments = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...filters
    }
    const res = await get('/api/admin/comments', params)
    comments.value = res.data || res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.error('加载评论列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  selectedComments.value = []
  selectAll.value = false
  loadComments()
}

const refreshList = () => {
  loadComments()
}

const handleSelectAll = () => {
  if (selectAll.value) {
    selectedComments.value = comments.value.map(c => c.id)
  } else {
    selectedComments.value = []
  }
}

const handleApprove = async (id) => {
  try {
    await post(`/api/admin/comments/${id}/approve`)
    ElMessage.success('评论已批准')
    loadComments()
  } catch (error) {
    ElMessage.error('批准评论失败')
  }
}

const handleReject = async (id) => {
  try {
    await post(`/api/admin/comments/${id}/reject`)
    ElMessage.success('评论已拒绝')
    loadComments()
  } catch (error) {
    ElMessage.error('拒绝评论失败')
  }
}

const handleEdit = (comment) => {
  Object.assign(editForm, {
    id: comment.id,
    author_name: comment.author_name,
    author_email: comment.author_email,
    status: comment.status,
    content: comment.content
  })
  editDialogVisible.value = true
}

const handleSaveEdit = async () => {
  try {
    await put(`/api/admin/comments/${editForm.id}`, editForm)
    ElMessage.success('评论更新成功')
    editDialogVisible.value = false
    loadComments()
  } catch (error) {
    ElMessage.error('更新评论失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗?', '提示', {
      type: 'warning'
    })
    await del(`/api/admin/comments/${id}`)
    ElMessage.success('评论删除成功')
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除评论失败')
    }
  }
}

const handleBatchApprove = async () => {
  try {
    await post('/api/admin/comments/batch/approve', selectedComments.value)
    ElMessage.success('批量批准成功')
    selectedComments.value = []
    selectAll.value = false
    loadComments()
  } catch (error) {
    ElMessage.error('批量批准失败')
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedComments.value.length} 条评论吗?`, '提示', {
      type: 'warning'
    })
    await post('/api/admin/comments/batch/delete', selectedComments.value)
    ElMessage.success('批量删除成功')
    selectedComments.value = []
    selectAll.value = false
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    approved: 'success',
    spam: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    pending: '待审核',
    approved: '已批准',
    spam: '垃圾评论'
  }
  return textMap[status] || status
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadPosts()
  loadComments()
})
</script>

<style scoped>
.comments-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  margin: 0;
}

.filter-card {
  margin-bottom: var(--spacing-md);
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.filter-form .el-form-item {
  margin-bottom: 0;
}

.comments-card {
  flex: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.comments-list {
  min-height: 400px;
}

.comment-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--color-border);
  transition: background-color var(--transition-fast);
}

.comment-item:hover {
  background-color: var(--color-bg-hover);
}

.comment-checkbox {
  margin-top: var(--spacing-md);
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-sm);
}

.comment-author {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.author-avatar {
  flex-shrink: 0;
}

.author-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.author-name {
  font-weight: 500;
  color: var(--color-text-primary);
}

.status-tag {
  align-self: flex-start;
}

.comment-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.comment-date {
  color: var(--color-text-secondary);
}

.comment-post {
  color: var(--color-primary);
}

.comment-text {
  margin: var(--spacing-md) 0;
  line-height: 1.6;
  color: var(--color-text-primary);
}

.comment-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.replies-list {
  margin-top: var(--spacing-md);
  padding-left: var(--spacing-lg);
  border-left: 2px solid var(--color-border);
}

.reply-item {
  display: flex;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
}

.reply-avatar {
  flex-shrink: 0;
}

.reply-content {
  flex: 1;
  background-color: var(--color-bg-light);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-sm);
}

.reply-author {
  font-weight: 500;
  color: var(--color-text-primary);
  margin-right: var(--spacing-sm);
}

.reply-text {
  color: var(--color-text-secondary);
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-lg);
}
</style>
