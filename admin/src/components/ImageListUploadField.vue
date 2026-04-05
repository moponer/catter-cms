<template>
  <div class="image-list-field">
    <div class="image-list-container">
      <draggable
        v-model="imageList"
        item-key="url"
        handle=".drag-handle"
        class="image-list"
        @end="handleReorder"
      >
        <template #item="{ element, index }">
          <div class="image-item">
            <div class="image-preview">
              <img :src="element.url" alt="" />
              <div class="image-overlay">
                <el-button size="small" class="drag-handle" icon="Rank" circle />
                <el-button size="small" type="danger" icon="Delete" circle @click="removeImage(index)" />
              </div>
            </div>
            <div class="image-meta">
              <el-input v-model="element.title" placeholder="图片标题(可选)" size="small" @input="updateValue" />
              <el-input v-model="element.link" placeholder="点击链接(可选)" size="small" @input="updateValue" />
            </div>
          </div>
        </template>
      </draggable>

      <div class="add-image-actions">
        <!-- 压缩设置 -->
        <el-popover placement="bottom-start" :width="270" trigger="click">
          <template #reference>
            <el-button :type="compressEnabled ? 'warning' : 'default'" plain size="small">
              <el-icon><Setting /></el-icon>
              压缩{{ compressEnabled ? '已开' : '已关' }}
            </el-button>
          </template>
          <div class="compress-settings">
            <div class="compress-title">图片压缩设置</div>
            <el-form label-width="80px" size="small">
              <el-form-item label="自动压缩">
                <el-switch v-model="compressEnabled" />
              </el-form-item>
              <template v-if="compressEnabled">
                <el-form-item label="触发阈值">
                  <el-input-number
                    v-model="compressThresholdMB"
                    :min="0.1" :max="10" :step="0.1" :precision="1"
                    size="small" style="width:120px"
                  />
                  <span class="unit-label">MB</span>
                </el-form-item>
                <el-form-item label="压缩质量">
                  <el-slider
                    v-model="compressQuality"
                    :min="0.1" :max="1" :step="0.05"
                    :format-tooltip="v => Math.round(v*100) + '%'"
                    style="width:120px"
                  />
                </el-form-item>
                <el-form-item label="最大边长">
                  <el-input-number
                    v-model="compressMaxSize"
                    :min="400" :max="4096" :step="100"
                    size="small" style="width:120px"
                  />
                  <span class="unit-label">px</span>
                </el-form-item>
              </template>
            </el-form>
          </div>
        </el-popover>

        <el-upload
          :show-file-list="false"
          :before-upload="handleUpload"
          accept="image/*"
          multiple
        >
          <el-button type="primary">
            <el-icon><Plus /></el-icon>
            上传图片
          </el-button>
        </el-upload>
        <el-button @click="openMediaPicker">
          <el-icon><Files /></el-icon> 从媒体库选择
        </el-button>
      </div>

      <!-- 上传队列进度 -->
      <transition name="fade">
        <div v-if="uploadQueue.length > 0" class="upload-queue">
          <div class="queue-header">
            <span>上传队列 ({{ uploadDoneCount }}/{{ uploadQueue.length }})</span>
            <el-button link size="small" @click="uploadQueue = uploadQueue.filter(t => t.status === 'uploading')">
              清除已完成
            </el-button>
          </div>
          <div
            v-for="task in uploadQueue"
            :key="task.id"
            class="queue-item"
            :class="task.status"
          >
            <el-icon v-if="task.status === 'done'" color="#67c23a" style="flex-shrink:0"><CircleCheck /></el-icon>
            <el-icon v-else-if="task.status === 'error'" color="#f56c6c" style="flex-shrink:0"><CircleClose /></el-icon>
            <el-icon v-else class="rotating" style="flex-shrink:0"><Loading /></el-icon>
            <span class="queue-name">{{ task.name }}</span>
            <span v-if="task.compressed" class="compress-badge">已压缩</span>
            <span v-if="task.status === 'error'" class="error-tip">失败</span>
          </div>
        </div>
      </transition>
    </div>

    <!-- 媒体库选择弹窗 -->
    <el-dialog v-model="mediaPickerVisible" title="选择图片" width="80%" :close-on-click-modal="false">
      <div class="media-picker">
        <div class="filter-bar">
          <el-input
            v-model="keyword"
            placeholder="搜索文件名"
            clearable
            @keyup.enter="loadMedia"
            style="width: 200px;"
          />
          <el-button type="primary" @click="loadMedia">
            <el-icon><Search /></el-icon> 搜索
          </el-button>
        </div>

        <div class="media-grid" v-loading="loading">
          <div
            v-for="item in mediaList"
            :key="item.id"
            class="media-item"
            :class="{ selected: selectedUrls.includes(item.url || item.path) }"
            @click="toggleSelect(item)"
          >
            <div class="media-preview">
              <img :src="item.url || item.path" :alt="item.filename" />
              <div v-if="selectedUrls.includes(item.url || item.path)" class="selected-badge">
                <el-icon><Check /></el-icon>
              </div>
            </div>
            <div class="media-info">
              <div class="filename">{{ item.filename }}</div>
            </div>
          </div>
        </div>

        <el-empty v-if="!loading && mediaList.length === 0" description="暂无媒体文件" />

        <div class="pagination" v-if="total > pageSize">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[20, 40, 60]"
            :total="total"
            layout="total, sizes, prev, pager, next"
            @size-change="loadMedia"
            @current-change="loadMedia"
          />
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <span class="selected-count">已选择 {{ selectedUrls.length }} 张图片</span>
          <div>
            <el-button @click="mediaPickerVisible = false">取消</el-button>
            <el-button type="primary" @click="confirmSelection" :disabled="selectedUrls.length === 0">
              添加选中图片
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete, Rank, Files, Search, Check, Setting, CircleCheck, CircleClose, Loading } from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import { get, post } from '@/utils/request'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

// Parse JSON string to array
const imageList = ref([])

watch(() => props.modelValue, (val) => {
  if (val && typeof val === 'string') {
    try {
      imageList.value = JSON.parse(val)
    } catch {
      imageList.value = []
    }
  } else if (Array.isArray(val)) {
    imageList.value = val
  } else {
    imageList.value = []
  }
}, { immediate: true })

// Emit JSON string when changed
const updateValue = () => {
  emit('update:modelValue', JSON.stringify(imageList.value))
}

// Watch imageList changes to auto-save
watch(imageList, () => {
  updateValue()
}, { deep: true })

// Upload state
const uploading = ref(false)

// ====== 压缩配置 ======
const compressEnabled = ref(true)
const compressThresholdMB = ref(1)
const compressQuality = ref(0.82)
const compressMaxSize = ref(1920)

// ====== 上传队列 ======
const uploadQueue = ref([])
let taskIdCounter = 0

const uploadDoneCount = computed(() =>
  uploadQueue.value.filter(t => t.status === 'done' || t.status === 'error').length
)

// ====== 图片压缩 ======
const compressImage = (file) => {
  return new Promise((resolve) => {
    const thresholdBytes = compressThresholdMB.value * 1024 * 1024
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
        if (width > height) { height = Math.round((height * max) / width); width = max }
        else { width = Math.round((width * max) / height); height = max }
      }
      const canvas = document.createElement('canvas')
      canvas.width = width; canvas.height = height
      canvas.getContext('2d').drawImage(img, 0, 0, width, height)
      canvas.toBlob((blob) => {
        if (!blob || blob.size >= file.size) { resolve({ file, compressed: false }); return }
        resolve({ file: new File([blob], file.name, { type: 'image/jpeg', lastModified: Date.now() }), compressed: true })
      }, 'image/jpeg', compressQuality.value)
    }
    img.onerror = () => { URL.revokeObjectURL(url); resolve({ file, compressed: false }) }
    img.src = url
  })
}

// Media picker
const mediaPickerVisible = ref(false)
const loading = ref(false)
const mediaList = ref([])
const selectedUrls = ref([])
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(40)
const total = ref(0)

// Load media
const loadMedia = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      type: 'image'
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }

    const res = await get('/api/media', { params })
    mediaList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}

// Open media picker
const openMediaPicker = () => {
  selectedUrls.value = []
  mediaPickerVisible.value = true
  loadMedia()
}

// Toggle selection
const toggleSelect = (item) => {
  const url = item.url || item.path
  const index = selectedUrls.value.indexOf(url)
  if (index > -1) {
    selectedUrls.value.splice(index, 1)
  } else {
    selectedUrls.value.push(url)
  }
}

// Confirm selection
const confirmSelection = () => {
  for (const url of selectedUrls.value) {
    // Avoid duplicates
    if (!imageList.value.find(img => img.url === url)) {
      imageList.value.push({ url, title: '', link: '' })
    }
  }
  updateValue()
  mediaPickerVisible.value = false
  ElMessage.success(`已添加 ${selectedUrls.value.length} 张图片`)
}

// Handle upload
const handleUpload = async (file) => {
  const taskId = ++taskIdCounter
  const task = { id: taskId, name: file.name, status: 'uploading', compressed: false }
  uploadQueue.value.unshift(task)

  try {
    const { file: uploadFile, compressed } = await compressImage(file)
    task.compressed = compressed

    const formData = new FormData()
    formData.append('file', uploadFile)
    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    // 兼容多种响应格式
    const url = res?.data?.url || res?.data?.path || res?.url || res?.path
    if (url) {
      imageList.value.push({ url, title: '', link: '' })
      updateValue()
      const idx = uploadQueue.value.findIndex(t => t.id === taskId)
      if (idx > -1) uploadQueue.value.splice(idx, 1)
    } else {
      throw new Error('服务器返回格式异常')
    }
  } catch (error) {
    const idx = uploadQueue.value.findIndex(t => t.id === taskId)
    if (idx > -1) uploadQueue.value.splice(idx, 1, { ...uploadQueue.value[idx], status: 'error' })
  }
  return false
}

// Remove image
const removeImage = (index) => {
  imageList.value.splice(index, 1)
  updateValue()
}

// Handle reorder
const handleReorder = () => {
  updateValue()
}
</script>

<style scoped>
.image-list-field {
  width: 100%;
}

.image-list-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.image-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.image-item {
  background: #f5f7fa;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.image-preview {
  position: relative;
  aspect-ratio: 16/9;
  overflow: hidden;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}

.image-item:hover .image-overlay {
  opacity: 1;
}

.drag-handle {
  cursor: grab;
}

.drag-handle:active {
  cursor: grabbing;
}

.image-meta {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.add-image-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

/* 压缩设置 */
.compress-settings { padding: 4px 0; }
.compress-title { font-size: 13px; font-weight: 600; margin-bottom: 10px; color: var(--el-text-color-primary); }
.unit-label { margin-left: 6px; font-size: 12px; color: var(--el-text-color-secondary); }
.compress-settings .el-form-item { margin-bottom: 8px; }

/* 上传队列 */
.upload-queue {
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  overflow: hidden;
}
.queue-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  font-size: 12px;
  color: var(--el-text-color-secondary);
  font-weight: 600;
}
.queue-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  font-size: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.queue-item:last-child { border-bottom: none; }
.queue-name { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.compress-badge {
  background: var(--el-color-warning-light-8);
  color: var(--el-color-warning);
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 4px;
  font-weight: 600;
  flex-shrink: 0;
}
.error-tip { color: var(--el-color-danger); flex-shrink: 0; }
.rotating { animation: spin 1s linear infinite; color: var(--el-color-primary); }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

/* Media Picker */
.media-picker {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-bar {
  display: flex;
  gap: 12px;
  align-items: center;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
  max-height: 400px;
  overflow-y: auto;
  padding: 8px;
  background: #f5f5f5;
  border-radius: 8px;
}

.media-item {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.media-item:hover {
  border-color: #409eff;
}

.media-item.selected {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.media-preview {
  position: relative;
  width: 100%;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f9f9f9;
}

.media-preview img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.selected-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 24px;
  height: 24px;
  background: #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.media-info {
  padding: 8px;
  border-top: 1px solid #eee;
}

.filename {
  font-size: 12px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.selected-count {
  color: #909399;
  font-size: 14px;
}
</style>
