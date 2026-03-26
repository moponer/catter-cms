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
              <el-input v-model="element.title" placeholder="图片标题(可选)" size="small" />
              <el-input v-model="element.link" placeholder="点击链接(可选)" size="small" />
            </div>
          </div>
        </template>
      </draggable>

      <div class="add-image-actions">
        <el-upload
          :show-file-list="false"
          :before-upload="handleUpload"
          accept="image/*"
          :disabled="uploading"
        >
          <el-button type="primary" :loading="uploading">
            <el-icon v-if="!uploading"><Plus /></el-icon>
            {{ uploading ? '上传中...' : '上传图片' }}
          </el-button>
        </el-upload>
        <el-button @click="openMediaPicker">
          <el-icon><Files /></el-icon> 从媒体库选择
        </el-button>
      </div>
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
import { Plus, Delete, Rank, Files, Search, Check } from '@element-plus/icons-vue'
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

// Upload state
const uploading = ref(false)

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
  const formData = new FormData()
  formData.append('file', file)

  uploading.value = true
  try {
    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    ElMessage.success('上传成功')
    if (res.data) {
      const url = res.data.url || res.data.path
      imageList.value.push({ url, title: '', link: '' })
      updateValue()
    }
    loadMedia()
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
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
}

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
