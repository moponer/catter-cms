<template>
  <div class="image-upload-field">
    <div class="image-preview" v-if="modelValue">
      <img :src="modelValue" alt="预览" class="preview-img" />
      <div class="preview-actions">
        <el-button size="small" @click="openMediaPicker">
          <el-icon><Files /></el-icon> 媒体库
        </el-button>
        <el-button size="small" type="danger" @click="removeImage">
          <el-icon><Delete /></el-icon> 移除
        </el-button>
      </div>
    </div>
    <div class="image-actions" v-else>
      <el-upload
        :show-file-list="false"
        :before-upload="handleUpload"
        accept="image/*"
        :disabled="uploading"
      >
        <el-button type="primary" :loading="uploading">
          <el-icon v-if="!uploading"><Upload /></el-icon>
          {{ uploading ? '上传中...' : '上传图片' }}
        </el-button>
      </el-upload>
      <el-button @click="openMediaPicker" :disabled="uploading">
        <el-icon><Files /></el-icon> 从媒体库选择
      </el-button>
    </div>

    <!-- 媒体库选择弹窗 -->
    <el-dialog v-model="mediaPickerVisible" title="选择图片" width="80%" :close-on-click-modal="false">
      <div class="media-picker">
        <!-- 筛选栏 -->
        <div class="filter-bar">
          <el-select v-model="filterType" placeholder="全部类型" clearable @change="loadMedia">
            <el-option label="图片" value="image" />
          </el-select>
          <el-input
            v-model="keyword"
            placeholder="搜索文件名"
            clearable
            @keyup.enter="loadMedia"
            style="width: 200px; margin-left: 12px;"
          />
          <el-upload
            :show-file-list="false"
            :before-upload="handleUpload"
            accept="image/*"
            style="margin-left: 12px;"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              上传新文件
            </el-button>
          </el-upload>
        </div>

        <!-- 媒体网格 -->
        <div class="media-grid" v-loading="loading">
          <div
            v-for="item in mediaList"
            :key="item.id"
            class="media-item"
            :class="{ selected: selectedId === item.id }"
            @click="handleSelect(item)"
            @dblclick="handleConfirm"
          >
            <div class="media-preview">
              <img v-if="isImage(item)" :src="item.url || item.path" :alt="item.filename" />
              <div v-else class="file-icon">
                <el-icon :size="32"><Document /></el-icon>
              </div>
            </div>
            <div class="media-info">
              <div class="filename">{{ item.filename }}</div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <el-empty v-if="!loading && mediaList.length === 0" description="暂无媒体文件" />

        <!-- 分页 -->
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
        <el-button @click="mediaPickerVisible = false">取消</el-button>
        <el-button type="primary" @click="handleConfirm" :disabled="!selectedItem">
          确定选择
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Files, Document, Delete } from '@element-plus/icons-vue'
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

// Upload state
const uploading = ref(false)

// Media picker state
const mediaPickerVisible = ref(false)
const loading = ref(false)
const mediaList = ref([])
const selectedId = ref(null)
const selectedItem = ref(null)
const filterType = ref('')
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// Watch for media picker visibility
watch(mediaPickerVisible, (val) => {
  if (val) {
    loadMedia()
    // If we already have a value, try to select it
    if (props.modelValue) {
      const found = mediaList.value.find(item => (item.url || item.path) === props.modelValue)
      if (found) {
        selectedId.value = found.id
        selectedItem.value = found
      }
    }
  }
})

// Open media picker
const openMediaPicker = () => {
  selectedId.value = null
  selectedItem.value = null
  mediaPickerVisible.value = true
}

// Load media list
const loadMedia = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (filterType.value) {
      params.type = filterType.value
    } else {
      params.type = 'image'
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }

    const res = await get('/api/media', { params })
    mediaList.value = res.data?.list || []
    total.value = res.data?.total || 0

    // Auto-select current image if in list
    if (props.modelValue) {
      const found = mediaList.value.find(item => (item.url || item.path) === props.modelValue)
      if (found) {
        selectedId.value = found.id
        selectedItem.value = found
      }
    }
  } catch (error) {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}

// Check if item is image
const isImage = (item) => {
  const type = item.file_type || item.type || ''
  return type.startsWith('image/') || /\.(jpg|jpeg|png|gif|webp|svg)$/i.test(item.filename)
}

// Select media item
const handleSelect = (item) => {
  selectedId.value = item.id
  selectedItem.value = item
}

// Confirm selection
const handleConfirm = () => {
  if (selectedItem.value) {
    const url = selectedItem.value.url || selectedItem.value.path
    emit('update:modelValue', url)
    mediaPickerVisible.value = false
  }
}

// Handle file upload
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
      emit('update:modelValue', url)
    }
    // Refresh media list
    loadMedia()
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
  return false
}

// Remove image
const removeImage = () => {
  emit('update:modelValue', '')
}
</script>

<style scoped>
.image-upload-field {
  width: 100%;
}

.image-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-img {
  max-width: 300px;
  max-height: 200px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.preview-actions {
  display: flex;
  gap: 8px;
}

.image-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* Media Picker Styles */
.media-picker {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
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

.file-icon {
  color: #909399;
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
</style>
