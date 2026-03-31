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
        <el-form-item>
          <el-button @click="loadMedia">搜索</el-button>
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
          <!-- 预览区 -->
          <div class="media-preview" @click.stop="handlePreview(item)">
            <img
              v-if="isImage(item)"
              :src="item.url"
              :alt="item.filename"
              loading="lazy"
              @error="onImgError($event)"
            />
            <div v-else-if="isVideo(item)" class="file-icon video-icon">
              <el-icon :size="40"><VideoPlay /></el-icon>
              <span>视频</span>
            </div>
            <div v-else class="file-icon">
              <el-icon :size="40"><Document /></el-icon>
              <span>{{ getExtension(item.filename) }}</span>
            </div>
            <!-- 悬浮遮罩 -->
            <div class="preview-overlay">
              <el-icon :size="24"><ZoomIn /></el-icon>
            </div>
          </div>

          <!-- 信息区 -->
          <div class="media-info">
            <div class="filename" :title="item.filename">{{ item.filename }}</div>
            <div class="meta">
              <span class="file-size">{{ formatSize(item.file_size) }}</span>
              <span v-if="item.width && item.height" class="dimensions">{{ item.width }}×{{ item.height }}</span>
            </div>
          </div>

          <!-- 操作区 -->
          <div class="media-actions">
            <el-button type="primary" link size="small" @click.stop="handlePreview(item)">
              <el-icon><ZoomIn /></el-icon> 查看
            </el-button>
            <el-button type="danger" link size="small" @click.stop="handleDelete(item)">
              <el-icon><Delete /></el-icon> 删除
            </el-button>
          </div>

          <!-- 选中勾 -->
          <div v-if="selectedIds.includes(item.id)" class="selected-badge">
            <el-icon><Select /></el-icon>
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

    <!-- 预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      :title="previewItem?.filename"
      width="780px"
      class="media-preview-dialog"
      destroy-on-close
    >
      <div class="preview-body">
        <!-- 图片预览 -->
        <div v-if="isImage(previewItem)" class="preview-image-wrap">
          <img :src="previewItem?.url" :alt="previewItem?.filename" />
        </div>
        <!-- 视频预览 -->
        <div v-else-if="isVideo(previewItem)" class="preview-video-wrap">
          <video :src="previewItem?.url" controls />
        </div>
        <!-- 其他文件 -->
        <div v-else class="preview-file-wrap">
          <el-icon :size="64" color="#909399"><Document /></el-icon>
          <p class="file-type-label">{{ getExtension(previewItem?.filename) }} 文件</p>
        </div>

        <!-- 文件信息 -->
        <div class="preview-meta">
          <div class="meta-row">
            <span class="meta-label">文件名</span>
            <span class="meta-value">{{ previewItem?.filename }}</span>
          </div>
          <div class="meta-row" v-if="previewItem?.file_size">
            <span class="meta-label">文件大小</span>
            <span class="meta-value">{{ formatSize(previewItem?.file_size) }}</span>
          </div>
          <div class="meta-row" v-if="previewItem?.width && previewItem?.height">
            <span class="meta-label">图片尺寸</span>
            <span class="meta-value">{{ previewItem?.width }} × {{ previewItem?.height }} px</span>
          </div>
          <div class="meta-row" v-if="previewItem?.file_type">
            <span class="meta-label">文件类型</span>
            <span class="meta-value">{{ previewItem?.file_type }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">访问地址</span>
            <span class="meta-value url-value">
              <span class="url-text">{{ previewItem?.url }}</span>
              <el-button type="primary" link size="small" @click="copyUrl(previewItem?.url)">
                <el-icon><CopyDocument /></el-icon> 复制
              </el-button>
            </span>
          </div>
          <div class="meta-row" v-if="previewItem?.created_at">
            <span class="meta-label">上传时间</span>
            <span class="meta-value">{{ formatDate(previewItem?.created_at) }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="previewVisible = false">关闭</el-button>
        <el-button type="primary" @click="copyUrl(previewItem?.url)">
          <el-icon><CopyDocument /></el-icon> 复制链接
        </el-button>
        <el-button type="danger" @click="handleDeleteFromPreview(previewItem)">
          <el-icon><Delete /></el-icon> 删除文件
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Document, Delete, VideoPlay, ZoomIn, Select, CopyDocument } from '@element-plus/icons-vue'
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
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (keyword.value) params.search = keyword.value       // 后端用 search
    if (filterType.value) params.file_type = filterType.value  // 后端用 file_type
    
    const res = await get('/api/media', params)
    mediaList.value = res.data?.list || res.list || []
    total.value = res.data?.total || res.total || 0
  } catch (error) {
    // 错误已在 request 拦截器中处理
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
      ElMessage.error('删除失败')
    }
  }
}

const handleDeleteFromPreview = async (item) => {
  try {
    await ElMessageBox.confirm('确定要删除该文件吗？', '提示', { type: 'warning' })
    await del(`/api/media/${item.id}`)
    ElMessage.success('删除成功')
    previewVisible.value = false
    loadMedia()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
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
      ElMessage.error('批量删除失败')
    }
  }
}

const copyUrl = async (url) => {
  if (!url) return
  try {
    await navigator.clipboard.writeText(url)
    ElMessage.success('链接已复制')
  } catch {
    ElMessage.error('复制失败，请手动复制')
  }
}

// 用 file_type 字段判断（后端 BeforeCreate 钩子自动填充）
const isImage = (item) => item?.file_type === 'image'
const isVideo = (item) => item?.file_type === 'video'

const getExtension = (filename) => {
  if (!filename) return ''
  const parts = filename.split('.')
  return parts.length > 1 ? ('.' + parts[parts.length - 1]).toUpperCase() : ''
}

const formatSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return `${size.toFixed(1)} ${units[i]}`
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

const onImgError = (e) => {
  // 图片加载失败时替换为占位背景
  e.target.style.display = 'none'
  e.target.parentElement.classList.add('img-error')
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

/* 网格布局 */
.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
  min-height: 80px;
}

/* 媒体卡片 */
.media-item {
  border: 1.5px solid var(--el-border-color);
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--el-bg-color);
  position: relative;
}

.media-item:hover {
  border-color: var(--el-color-primary-light-5);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.media-item.selected {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-7);
}

/* 预览区 */
.media-preview {
  height: 130px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color-lighter);
  overflow: hidden;
  position: relative;
}

.media-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.media-preview.img-error::after {
  content: '图片加载失败';
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

/* 悬浮遮罩 */
.preview-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  opacity: 0;
  transition: opacity 0.2s;
}

.media-preview:hover .preview-overlay {
  opacity: 1;
}

/* 文件图标 */
.file-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  color: var(--el-text-color-secondary);
}

.file-icon span {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.video-icon {
  color: var(--el-color-primary);
}

/* 信息区 */
.media-info {
  padding: 8px 10px 4px;
}

.filename {
  font-size: 12px;
  color: var(--el-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.4;
}

.meta {
  display: flex;
  gap: 8px;
  margin-top: 3px;
}

.file-size,
.dimensions {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}

/* 操作区 */
.media-actions {
  display: flex;
  justify-content: center;
  gap: 4px;
  padding: 6px 8px;
  border-top: 1px solid var(--el-border-color-lighter);
}

/* 选中角标 */
.selected-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--el-color-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.2);
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--spacing-md);
}

/* ====== 预览弹窗 ====== */
.preview-body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.preview-image-wrap {
  background: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  max-height: 420px;
}

.preview-image-wrap img {
  max-width: 100%;
  max-height: 420px;
  object-fit: contain;
  display: block;
}

.preview-video-wrap video {
  width: 100%;
  max-height: 360px;
  border-radius: 8px;
  background: #000;
}

.preview-file-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  gap: 12px;
}

.file-type-label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

/* 文件信息表格 */
.preview-meta {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  overflow: hidden;
}

.meta-row {
  display: flex;
  align-items: baseline;
  padding: 10px 16px;
  gap: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  font-size: 13px;
}

.meta-row:last-child {
  border-bottom: none;
}

.meta-label {
  flex-shrink: 0;
  width: 72px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.meta-value {
  color: var(--el-text-color-primary);
  word-break: break-all;
  flex: 1;
}

.url-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.url-text {
  flex: 1;
  font-family: monospace;
  font-size: 12px;
  color: var(--el-color-primary);
  word-break: break-all;
}
</style>
