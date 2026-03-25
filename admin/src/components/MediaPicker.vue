<template>
  <el-dialog
    v-model="visible"
    title="媒体库"
    width="80%"
    :close-on-click-modal="false"
  >
    <div class="media-picker">
      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-select v-model="filterType" placeholder="全部类型" clearable @change="loadMedia">
          <el-option label="图片" value="image" />
          <el-option label="视频" value="video" />
          <el-option label="文档" value="document" />
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
          accept="image/*,video/*"
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
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleConfirm" :disabled="!selectedItem">
        确定选择
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Document } from '@element-plus/icons-vue'
import { get, post } from '@/utils/request'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  fileType: {
    type: String,
    default: 'image' // image, video, all
  }
})

const emit = defineEmits(['update:modelValue', 'select'])

const visible = ref(false)
const loading = ref(false)
const mediaList = ref([])
const selectedId = ref(null)
const selectedItem = ref(null)
const filterType = ref('')
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 监听 visible 变化
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    loadMedia()
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

// 加载媒体列表
const loadMedia = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (filterType.value) {
      params.type = filterType.value
    } else if (props.fileType !== 'all') {
      params.type = props.fileType
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }

    const res = await get('/api/media', { params })
    mediaList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

// 判断是否为图片
const isImage = (item) => {
  const type = item.file_type || item.type || ''
  return type.startsWith('image/') || /\.(jpg|jpeg|png|gif|webp|svg)$/i.test(item.filename)
}

// 选择媒体
const handleSelect = (item) => {
  selectedId.value = item.id
  selectedItem.value = item
}

// 确认选择
const handleConfirm = () => {
  if (selectedItem.value) {
    emit('select', selectedItem.value)
    visible.value = false
  }
}

// 上传新文件
const handleUpload = async (file) => {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    ElMessage.success('上传成功')
    loadMedia()
    // 自动选中刚上传的文件
    if (res.data) {
      selectedId.value = res.data.id
      selectedItem.value = res.data
    }
  } catch (error) {
    ElMessage.error('上传失败')
  }
  return false
}
</script>

<style scoped>
.media-picker {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-bar {
  display: flex;
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
