<template>
  <div class="posts-page">
    <div class="page-header">
      <h1>文章管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新建文章
      </el-button>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="filter-card">
      <el-form :inline="true">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="搜索标题" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="已发布" value="published" />
            <el-option label="草稿" value="draft" />
            <el-option label="已归档" value="archived" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 文章列表 -->
    <el-card class="table-card">
      <el-table :data="posts" v-loading="loading" style="width: 100%">
        <el-table-column label="封面" width="80" align="center">
          <template #default="{ row }">
            <div class="cover-cell">
              <el-image
                v-if="row.cover"
                :src="row.cover"
                fit="cover"
                class="cover-thumb"
                :preview-src-list="[row.cover]"
                :preview-teleported="true"
              >
                <template #error>
                  <div class="cover-fallback"><el-icon><Picture /></el-icon></div>
                </template>
              </el-image>
              <div v-else class="cover-empty"><el-icon><Picture /></el-icon></div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="author" label="作者" width="120">
          <template #default="{ row }">
            {{ row.user?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            {{ row.category?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row.id)">
              编辑
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row.id)">
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
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @size-change="loadPosts"
          @current-change="loadPosts"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Picture } from '@element-plus/icons-vue'
import { get, del } from '@/utils/request'

const router = useRouter()

const loading = ref(false)
const posts = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const loadPosts = async () => {
  loading.value = true
  try {
    const res = await get('/api/posts', {
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchForm
    })
    posts.value = res.data?.list || res.list || []
    total.value = res.data?.total || res.total || 0
  } catch (error) {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  router.push('/posts/create')
}

const handleEdit = (id) => {
  router.push(`/posts/${id}/edit`)
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
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

const handleSearch = () => {
  currentPage.value = 1
  loadPosts()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  handleSearch()
}

const getStatusType = (status) => {
  const types = {
    published: 'success',
    draft: 'info',
    archived: 'warning'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    published: '已发布',
    draft: '草稿',
    archived: '已归档'
  }
  return texts[status] || status
}

onMounted(() => {
  loadPosts()
})
</script>

<style scoped>
.posts-page {
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
  margin-bottom: 0;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--spacing-md);
}

/* 封面图 */
.cover-cell {
  display: flex;
  justify-content: center;
  align-items: center;
}
.cover-thumb {
  width: 54px;
  height: 36px;
  border-radius: 4px;
  cursor: pointer;
  border: 1px solid var(--el-border-color-lighter);
}
.cover-empty,
.cover-fallback {
  width: 54px;
  height: 36px;
  border-radius: 4px;
  background: var(--el-fill-color-light);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-placeholder);
  font-size: 16px;
  border: 1px dashed var(--el-border-color-lighter);
}
</style>
