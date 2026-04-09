<template>
  <div class="posts-container">
    <div class="page-header">
      <h1>文章管理</h1>
      <el-button type="primary" @click="createPost">新建文章</el-button>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="filter-card">
      <el-row :gutter="16">
        <el-col :xs="24" :sm="8" :md="6">
          <el-input
            v-model="filters.keyword"
            placeholder="搜索标题或内容"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :xs="12" :sm="8" :md="4">
          <el-select
            v-model="filters.status"
            placeholder="文章状态"
            clearable
            @change="handleFilter"
            style="width: 100%"
          >
            <el-option label="全部" value="" />
            <el-option label="已发布" value="published" />
            <el-option label="草稿" value="draft" />
            <el-option label="待审核" value="pending" />
          </el-select>
        </el-col>
        <el-col :xs="12" :sm="8" :md="4">
          <el-select
            v-model="filters.category_id"
            placeholder="选择分类"
            clearable
            @change="handleFilter"
            style="width: 100%"
          >
            <el-option label="全部分类" value="" />
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6">
          <el-date-picker
            v-model="filters.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="handleFilter"
            style="width: 100%"
          />
        </el-col>
        <el-col :xs="24" :sm="8" :md="4">
          <el-button @click="resetFilters">重置筛选</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 批量操作 -->
    <div v-if="selectedPosts.length > 0" class="bulk-actions">
      <span class="selected-count">已选择 {{ selectedPosts.length }} 项</span>
      <el-button type="danger" size="small" @click="bulkDelete">批量删除</el-button>
      <el-button type="primary" size="small" @click="bulkPublish">批量发布</el-button>
      <el-button size="small" @click="bulkDraft">批量设为草稿</el-button>
    </div>

    <!-- 文章列表 -->
    <el-table
      :data="posts"
      style="width: 100%"
      v-loading="loading"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="author" label="作者" width="100">
        <template #default="{ row }">
          {{ row.author_name || row.Author?.name || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="category" label="分类" width="120">
        <template #default="{ row }">
          {{ row.category_name || row.Category?.name || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag
            :type="getStatusType(row.status)"
            size="small"
          >
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="更新时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.updated_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            link
            @click="editPost(row.id)"
          >
            编辑
          </el-button>
          <el-button
            type="danger"
            size="small"
            link
            @click="deletePost(row.id)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { get, del, delBody, put } from '@/utils/request'

const router = useRouter()

// 加载状态
const loading = ref(false)

// 分页数据
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 文章列表
const posts = ref([])

// 分类列表
const categories = ref([])

// 选中的文章
const selectedPosts = ref([])

// 筛选条件
const filters = reactive({
  keyword: '',
  status: '',
  category_id: '',
  date_range: null
})

onMounted(() => {
  loadCategories()
  loadPosts()
})

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await get('/api/categories')
    categories.value = res.data || res || []
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

// 加载文章列表
const loadPosts = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: filters.keyword || undefined,
      status: filters.status || undefined,
      category_id: filters.category_id || undefined
    }

    // 处理日期范围
    if (filters.date_range && filters.date_range.length === 2) {
      params.start_date = filters.date_range[0]
      params.end_date = filters.date_range[1]
    }

    const res = await get('/api/posts', { params })
    posts.value = res.data?.list || res.list || res.data || res || []
    total.value = res.data?.total || res.total || 0
  } catch (error) {
    ElMessage.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
let searchTimer = null
const handleSearch = () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    loadPosts()
  }, 500)
}

// 筛选处理
const handleFilter = () => {
  currentPage.value = 1
  loadPosts()
}

// 重置筛选
const resetFilters = () => {
  filters.keyword = ''
  filters.status = ''
  filters.category_id = ''
  filters.date_range = null
  currentPage.value = 1
  loadPosts()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedPosts.value = selection
}

// 新建文章
const createPost = () => {
  router.push('/posts/create')
}

// 编辑文章
const editPost = (id) => {
  router.push(`/posts/${id}/edit`)
}

// 删除文章
const deletePost = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await del(`/api/posts/${id}`)
    ElMessage.success('删除成功')
    loadPosts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const bulkDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedPosts.value.length} 篇文章吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const ids = selectedPosts.value.map(item => item.id).join(',')
    await delBody('/api/posts/batch', { ids })
    ElMessage.success('批量删除成功')
    selectedPosts.value = []
    loadPosts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 批量发布
const bulkPublish = async () => {
  try {
    const ids = selectedPosts.value.map(item => item.id).join(',')
    await put('/api/posts/batch/status', { ids, status: 'published' })
    ElMessage.success('批量发布成功')
    selectedPosts.value = []
    loadPosts()
  } catch (error) {
    ElMessage.error('批量发布失败')
  }
}

// 批量设为草稿
const bulkDraft = async () => {
  try {
    const ids = selectedPosts.value.map(item => item.id).join(',')
    await put('/api/posts/batch/status', { ids, status: 'draft' })
    ElMessage.success('操作成功')
    selectedPosts.value = []
    loadPosts()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 分页大小变化
const handleSizeChange = (size) => {
  pageSize.value = size
  loadPosts()
}

// 当前页变化
const handleCurrentChange = (current) => {
  currentPage.value = current
  loadPosts()
}

// 获取状态标签类型
const getStatusType = (status) => {
  const typeMap = {
    published: 'success',
    draft: 'info',
    pending: 'warning'
  }
  return typeMap[status] || 'info'
}

// 获取状态标签文字
const getStatusLabel = (status) => {
  const labelMap = {
    published: '已发布',
    draft: '草稿',
    pending: '待审核'
  }
  return labelMap[status] || status
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  try {
    const date = new Date(dateStr)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}
</script>

<style scoped>
.posts-container {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.filter-card {
  margin-bottom: var(--spacing-md);
}

.bulk-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background-color: var(--color-bg-hover);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-md);
}

.selected-count {
  font-weight: 500;
  color: var(--color-text-primary);
}

.pagination {
  display: flex;
  justify-content: flex-end;
  padding-top: var(--spacing-md);
}

@media (max-width: 768px) {
  .bulk-actions {
    flex-direction: column;
    align-items: flex-start;
  }

  .pagination {
    justify-content: center;
  }
}
</style>
