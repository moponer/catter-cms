<template>
  <div class="media-page">
    <div class="page-header">
      <h1>媒体库</h1>
      <div class="header-actions">
        <el-button
          v-if="selectedIds.length > 0"
          type="danger"
          @click="handleBatchDelete"
        >
          <el-icon><Delete /></el-icon>
          删除选中 ({{ selectedIds.length }})
        </el-button>
        <el-upload
          :show-file-list="false"
          :before-upload="handleUpload"
          accept="image/*,video/*,.pdf,.doc,.docx"
          multiple
        >
          <el-button type="primary">
            <el-icon><Upload /></el-icon>
            上传文件
          </el-button>
        </el-upload>
      </div>
    </div>

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-form :inline="true">
        <el-form-item label="类型">
          <el-select v-model="filterType" placeholder="全部类型" clearable @change="loadMedia">
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="文档" value="document" />
          </el-select>
        </el-form-item>
        <el-form-item label="搜索">
          <el-input v-model="keyword" placeholder="文件名" clearable @keyup.enter="loadMedia" />
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 媒体网格 -->
    <el-card class="media-card">
      <div class="media-grid" v-loading="loading">
        <div
          v-for="item in mediaList"
          :key="item.id"
          class="media-item"
          :class="{ selected: selectedIds.includes(item.id) }"
          @click="handleSelect(item)"
        >
          <div class="media-preview">
            <img v-if="isImage(item)" :src="item.path" :alt="item.filename" />
            <div v-else class="file-icon">
              <el-icon :size="48"><Document /></el-icon>
            </div>
          </div>
          <div class="media-info">
            <div class="filename">{{ item.filename }}</div>
            <div class="file-size">{{ formatSize(item.size) }}</div>
          </div>
          <div class="media-actions">
            <el-button type="primary" link size="small" @click.stop="handlePreview(item)">
              查看
            </el-button>
            <el-button type="danger" link size="small" @click.stop="handleDelete(item)">
              删除
            </el-button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <el-empty v-if="!loading && mediaList.length === 0" description="暂无媒体文件" />

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[30, 60, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @size-change="loadMedia"
          @current-change="loadMedia"
        />
      </div>
    </el-card>

    <!-- 预览对话框 -->
    <el-dialog v-model="previewVisible" title="文件预览" width="80%">
      <div class="preview-content">
        <img v-if="isImage(previewItem)" :src="previewItem?.path" style="max-width: 100%" />
        <video v-else-if="isVideo(previewItem)" :src="previewItem?.path" controls style="max-width: 100%" />
        <div v-else class="file-info">
          <p><strong>文件名:</strong> {{ previewItem?.filename }}</p>
          <p><strong>类型:</strong> {{ previewItem?.mime_type }}</p>
          <p><strong>大小:</strong> {{ formatSize(previewItem?.size) }}</p>
          <p><strong>路径:</strong> {{ previewItem?.path }}</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Document, Delete } from '@element-plus/icons-vue'
import { get, post, del } from '@/utils/request'

const loading = ref(false)
const mediaList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(30)
const keyword = ref('')
const filterType = ref('')
const selectedIds = ref([])
const previewVisible = ref(false)
const previewItem = ref(null)

const loadMedia = async () => {
  loading.value = true
  try {
    const res = await get('/api/media', {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: keyword.value,
      type: filterType.value
    })
    mediaList.value = res.data?.list || res.list || []
    total.value = res.data?.total || res.total || 0
  } catch (error) {
    console.error('加载媒体列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleUpload = async (file) => {
  const formData = new FormData()
  formData.append('file', file)

  try {
    await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    ElMessage.success('上传成功')
    loadMedia()
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  }
  return false
}

const handleSelect = (item) => {
  const index = selectedIds.value.indexOf(item.id)
  if (index > -1) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(item.id)
  }
}

const handlePreview = (item) => {
  previewItem.value = item
  previewVisible.value = true
}

const handleDelete = async (item) => {
  try {
    await ElMessageBox.confirm('确定要删除该文件吗？', '提示', { type: 'warning' })
    await del(`/api/media/${item.id}`)
    ElMessage.success('删除成功')
    loadMedia()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const handleBatchDelete = async () => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择要删除的文件')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedIds.value.length} 个文件吗？`,
      '批量删除确认',
      { type: 'warning' }
    )
    await post('/api/media/batch-delete', { ids: selectedIds.value })
    ElMessage.success('批量删除成功')
    selectedIds.value = []
    loadMedia()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

const isImage = (item) => item?.mime_type?.startsWith('image/')
const isVideo = (item) => item?.mime_type?.startsWith('video/')

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  while (bytes >= 1024 && i < units.length - 1) {
    bytes /= 1024
    i++
  }
  return `${bytes.toFixed(1)} ${units[i]}`
}

onMounted(() => {
  loadMedia()
})
</script>

<style scoped>
.media-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
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

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
}

.media-item {
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
}

.media-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.media-item.selected {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-5);
}

.media-preview {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color-light);
  overflow: hidden;
}

.media-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-icon {
  color: var(--el-text-color-secondary);
}

.media-info {
  padding: 8px;
}

.filename {
  font-size: 13px;
  color: var(--el-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}

.media-actions {
  display: flex;
  justify-content: center;
  gap: 8px;
  padding: 8px;
  border-top: 1px solid var(--el-border-color-lighter);
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--spacing-md);
}

.preview-content {
  text-align: center;
}

.file-info {
  text-align: left;
  padding: 20px;
}

.file-info p {
  margin: 10px 0;
}
</style>
