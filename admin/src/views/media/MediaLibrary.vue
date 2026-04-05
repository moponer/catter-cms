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
        <!-- 压缩选项 toggle -->
        <el-popover
          placement="bottom-end"
          :width="280"
          trigger="click"
        >
          <template #reference>
            <el-button :type="compressEnabled ? 'warning' : 'default'" plain>
              <el-icon><Setting /></el-icon>
              压缩{{ compressEnabled ? '已开' : '已关' }}
            </el-button>
          </template>
          <div class="compress-settings">
            <div class="compress-title">图片压缩设置</div>
            <el-form label-width="90px" size="small">
              <el-form-item label="自动压缩">
                <el-switch v-model="compressEnabled" />
              </el-form-item>
              <template v-if="compressEnabled">
                <el-form-item label="触发阈值">
                  <el-slider
                    v-model="compressThresholdMB"
                    :min="0.1"
                    :max="10"
                    :step="0.1"
                    :format-tooltip="v => v + ' MB'"
                    show-input
                    :show-input-controls="false"
                    size="small"
                    style="width:140px"
                  />
                </el-form-item>
                <el-form-item label="压缩质量">
                  <el-slider
                    v-model="compressQuality"
                    :min="0.1"
                    :max="1"
                    :step="0.05"
                    :format-tooltip="v => Math.round(v*100) + '%'"
                    show-input
                    :show-input-controls="false"
                    size="small"
                    style="width:140px"
                  />
                </el-form-item>
                <el-form-item label="最大边长">
                  <el-input-number
                    v-model="compressMaxSize"
                    :min="400"
                    :max="4096"
                    :step="100"
                    size="small"
                    style="width:140px"
                  />
                  <span style="margin-left:6px;color:var(--el-text-color-secondary);font-size:12px">px</span>
                </el-form-item>
              </template>
            </el-form>
          </div>
        </el-popover>
        <el-upload
          ref="uploadRef"
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

    <!-- 批量上传进度浮层 -->
    <transition name="upload-progress-fade">
      <div v-if="uploadQueue.length > 0" class="upload-progress-panel">
        <div class="upload-progress-header">
          <span class="upload-progress-title">
            <el-icon><Upload /></el-icon>
            上传队列 ({{ uploadDoneCount }}/{{ uploadQueue.length }})
          </span>
          <el-button type="primary" link size="small" @click="clearDoneQueue">
            清除已完成
          </el-button>
        </div>
        <div class="upload-progress-list">
          <div
            v-for="task in uploadQueue"
            :key="task.id"
            class="upload-task-item"
            :class="task.status"
          >
            <div class="upload-task-icon">
              <el-icon v-if="task.status === 'done'" color="#67c23a"><CircleCheck /></el-icon>
              <el-icon v-else-if="task.status === 'error'" color="#f56c6c"><CircleClose /></el-icon>
              <el-icon v-else class="rotating"><Loading /></el-icon>
            </div>
            <div class="upload-task-info">
              <div class="upload-task-name">{{ task.name }}</div>
              <div class="upload-task-meta">
                {{ task.originalSize }}
                <template v-if="task.compressed"> → {{ task.compressedSize }} <span class="compress-badge">已压缩</span></template>
                <span v-if="task.status === 'error'" class="error-tip">{{ task.error }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>

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
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Document, Delete, VideoPlay, ZoomIn, Select, CopyDocument, Setting, CircleCheck, CircleClose, Loading } from '@element-plus/icons-vue'
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

// ====== 压缩配置 ======
const compressEnabled = ref(true)
const compressThresholdMB = ref(1)   // 超过此大小才压缩
const compressQuality = ref(0.82)    // JPEG 质量
const compressMaxSize = ref(1920)    // 最大边长(px)

// ====== 上传队列 ======
const uploadQueue = ref([])
let taskIdCounter = 0

const uploadDoneCount = computed(() =>
  uploadQueue.value.filter(t => t.status === 'done' || t.status === 'error').length
)

const clearDoneQueue = () => {
  uploadQueue.value = uploadQueue.value.filter(t => t.status === 'uploading')
}

// ====== 图片压缩 ======
const compressImage = (file) => {
  return new Promise((resolve) => {
    const thresholdBytes = compressThresholdMB.value * 1024 * 1024
    // 非图片或未启用压缩或未超阈值 → 直接返回原文件
    if (!compressEnabled.value || !file.type.startsWith('image/') || file.size < thresholdBytes) {
      resolve({ file, compressed: false })
      return
    }
    const img = new Image()
    const url = URL.createObjectURL(file)
    img.onload = () => {
      URL.revokeObjectURL(url)
      let { width, height } = img
      const max = compressMaxSize.value
      if (width > max || height > max) {
        if (width > height) {
          height = Math.round((height * max) / width)
          width = max
        } else {
          width = Math.round((width * max) / height)
          height = max
        }
      }
      const canvas = document.createElement('canvas')
      canvas.width = width
      canvas.height = height
      const ctx = canvas.getContext('2d')
      ctx.drawImage(img, 0, 0, width, height)
      canvas.toBlob(
        (blob) => {
          if (!blob || blob.size >= file.size) {
            // 压缩后反而变大，放弃压缩
            resolve({ file, compressed: false })
            return
          }
          const compressedFile = new File([blob], file.name, { type: 'image/jpeg', lastModified: Date.now() })
          resolve({ file: compressedFile, compressed: true })
        },
        'image/jpeg',
        compressQuality.value
      )
    }
    img.onerror = () => {
      URL.revokeObjectURL(url)
      resolve({ file, compressed: false })
    }
    img.src = url
  })
}

const loadMedia = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (keyword.value) params.search = keyword.value
    if (filterType.value) params.file_type = filterType.value
    
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
  const taskId = ++taskIdCounter
  const task = {
    id: taskId,
    name: file.name,
    originalSize: formatSize(file.size),
    compressedSize: '',
    compressed: false,
    status: 'uploading',
    error: ''
  }
  uploadQueue.value.unshift(task)

  try {
    // 尝试压缩
    const { file: uploadFile, compressed } = await compressImage(file)
    task.compressed = compressed
    if (compressed) {
      task.compressedSize = formatSize(uploadFile.size)
    }

    const formData = new FormData()
    formData.append('file', uploadFile)

    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    // 兼容多种响应格式
    const url = res?.data?.url || res?.data?.path || res?.url || res?.path
    if (url) {
      // 上传成功即标记完成，loadMedia 在用户关闭队列或刷新时再调
      const idx = uploadQueue.value.findIndex(t => t.id === taskId)
      if (idx > -1) uploadQueue.value.splice(idx, 1, { ...uploadQueue.value[idx], status: 'done' })
    } else {
      throw new Error('服务器返回格式异常')
    }
  } catch (error) {
    const idx = uploadQueue.value.findIndex(t => t.id === taskId)
    if (idx > -1) uploadQueue.value.splice(idx, 1, { ...uploadQueue.value[idx], status: 'error', error: error?.message || '上传失败' })
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

/* ====== 压缩设置 Popover 内容 ====== */
.compress-settings {
  padding: 4px 0;
}
.compress-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 12px;
}
.compress-settings .el-form-item {
  margin-bottom: 10px;
}

/* ====== 上传进度浮层 ====== */
.upload-progress-panel {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 340px;
  max-height: 400px;
  background: var(--el-bg-color-overlay);
  border: 1px solid var(--el-border-color);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  z-index: 2000;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.upload-progress-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px 10px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  flex-shrink: 0;
}

.upload-progress-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.upload-progress-list {
  overflow-y: auto;
  max-height: 320px;
  padding: 8px 0;
}

.upload-task-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 8px 16px;
  transition: background 0.15s;
}

.upload-task-item:hover {
  background: var(--el-fill-color-light);
}

.upload-task-icon {
  flex-shrink: 0;
  margin-top: 2px;
  font-size: 16px;
}

.upload-task-icon .rotating {
  animation: spin 1s linear infinite;
  color: var(--el-color-primary);
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.upload-task-info {
  flex: 1;
  min-width: 0;
}

.upload-task-name {
  font-size: 13px;
  color: var(--el-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.upload-task-meta {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  margin-top: 2px;
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.compress-badge {
  background: var(--el-color-warning-light-8);
  color: var(--el-color-warning);
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 4px;
  font-weight: 600;
}

.error-tip {
  color: var(--el-color-danger);
}

/* 进度面板动画 */
.upload-progress-fade-enter-active,
.upload-progress-fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.upload-progress-fade-enter-from,
.upload-progress-fade-leave-to {
  opacity: 0;
  transform: translateY(12px);
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
