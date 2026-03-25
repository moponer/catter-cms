<template>
  <div class="theme-list-page">
    <div class="page-header">
      <h2>主题管理</h2>
      <el-button type="primary" @click="handleUpload">
        <el-icon><Upload /></el-icon>
        上传主题
      </el-button>
    </div>

    <!-- 加载状态 -->
    <el-skeleton v-if="loading" :rows="8" animated />

    <!-- 主题网格 -->
    <div v-else class="theme-grid">
      <el-card
        v-for="theme in themes"
        :key="theme.slug"
        class="theme-card"
        :class="{ 'is-active': theme.is_active }"
        shadow="hover"
      >
        <!-- 缩略图 -->
        <div class="theme-thumbnail">
          <img
            v-if="theme.thumbnail"
            :src="theme.thumbnail"
            :alt="theme.name"
          />
          <div v-else class="thumbnail-placeholder">
            <el-icon :size="40"><Monitor /></el-icon>
          </div>
          <div v-if="theme.is_active" class="active-badge">使用中</div>
        </div>

        <!-- 主题信息 -->
        <div class="theme-info">
          <h3 class="theme-name">{{ theme.name }}</h3>
          <p class="theme-version">v{{ theme.version }}</p>
          <p class="theme-description">{{ theme.description || '暂无描述' }}</p>
          <div class="theme-meta">
            <span v-if="theme.author" class="author">
              <el-icon><User /></el-icon>
              {{ theme.author }}
            </span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="theme-actions">
          <el-button
            v-if="!theme.is_active"
            type="primary"
            size="small"
            @click="handleActivate(theme)"
            :loading="activating === theme.slug"
          >
            激活
          </el-button>
          <el-button
            v-else
            type="success"
            size="small"
            disabled
          >
            当前主题
          </el-button>
          <el-button
            size="small"
            @click="handleConfigure(theme)"
          >
            配置
          </el-button>
          <el-button
            size="small"
            @click="handlePreview(theme)"
          >
            预览
          </el-button>
          <el-button
            v-if="!theme.is_active"
            size="small"
            type="danger"
            @click="handleDelete(theme)"
          >
            删除
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 空状态 -->
    <el-empty v-if="!loading && themes.length === 0" description="暂无主题" />

    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传主题" width="500px">
      <el-upload
        ref="uploadRef"
        class="theme-upload"
        drag
        :action="uploadUrl"
        :headers="uploadHeaders"
        :auto-upload="false"
        :limit="1"
        accept=".zip"
        :before-upload="handleBeforeUpload"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :on-change="handleFileChange"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">拖拽 ZIP 文件到此处，或 <em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">主题包格式为 ZIP，不超过 10MB</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitUpload" :loading="uploading">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Monitor, User, Upload, UploadFilled } from '@element-plus/icons-vue'
import { listThemes, switchTheme, deleteTheme } from '@/api/theme'

const router = useRouter()

const themes = ref([])
const loading = ref(true)
const activating = ref(null)
const uploadDialogVisible = ref(false)
const uploading = ref(false)
const uploadRef = ref(null)
const uploadFile = ref(null)

const uploadUrl = '/api/themes/upload'
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 加载主题列表
async function fetchThemes() {
  try {
    loading.value = true
    const res = await listThemes()
    themes.value = res.data?.list || res.list || []
  } catch (err) {
    ElMessage.error('加载主题列表失败')
  } finally {
    loading.value = false
  }
}

// 激活主题
async function handleActivate(theme) {
  try {
    await ElMessageBox.confirm(
      `确定要激活主题「${theme.name}」吗？激活后将替换当前主题。`,
      '激活主题',
      { type: 'warning' }
    )
    activating.value = theme.slug
    await switchTheme(theme.slug)
    ElMessage.success('主题激活成功')
    await fetchThemes()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('激活失败')
    }
  } finally {
    activating.value = null
  }
}

// 跳转配置页
function handleConfigure(theme) {
  router.push(`/themes/${theme.slug}`)
}

// 新窗口预览主题
function handlePreview(theme) {
  // 打开主题预览：前台首页 + theme 参数，后端会根据参数渲染对应主题
  window.open('/?theme=' + theme.slug, '_blank')
}

// 删除主题
async function handleDelete(theme) {
  try {
    await ElMessageBox.confirm(
      `确定要删除主题「${theme.name}」吗？此操作不可恢复。`,
      '删除主题',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await deleteTheme(theme.slug)
    ElMessage.success('主题删除成功')
    await fetchThemes()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 上传主题
function handleUpload() {
  uploadFile.value = null
  uploadDialogVisible.value = true
}

function handleBeforeUpload(file) {
  const isZip = file.name.endsWith('.zip')
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isZip) {
    ElMessage.error('只能上传 ZIP 格式的主题包')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('主题包不能超过 10MB')
    return false
  }
  return true
}

function handleFileChange(file) {
  uploadFile.value = file.raw
}

function submitUpload() {
  if (!uploadFile.value) {
    ElMessage.warning('请选择主题文件')
    return
  }
  uploadRef.value?.submit()
}

function handleUploadSuccess(res) {
  ElMessage.success('主题上传成功')
  uploadDialogVisible.value = false
  fetchThemes()
}

function handleUploadError(err) {
  ElMessage.error('上传失败：' + (err.message || '未知错误'))
}

onMounted(() => {
  fetchThemes()
})
</script>

<style scoped>
.theme-list-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.theme-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.theme-card {
  display: flex;
  flex-direction: column;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.theme-card.is-active {
  border-color: var(--el-color-success);
}

.theme-thumbnail {
  position: relative;
  height: 180px;
  overflow: hidden;
  border-radius: 4px 4px 0 0;
  background: #f5f7fa;
}

.theme-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumbnail-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #c0c4cc;
}

.active-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 12px;
  background: var(--el-color-success);
  color: #fff;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.theme-info {
  padding: 16px;
  flex: 1;
}

.theme-name {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 600;
}

.theme-version {
  margin: 0 0 8px;
  font-size: 12px;
  color: #909399;
}

.theme-description {
  margin: 0 0 12px;
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.theme-meta {
  font-size: 12px;
  color: #909399;
}

.theme-meta .author {
  display: flex;
  align-items: center;
  gap: 4px;
}

.theme-actions {
  padding: 12px 16px;
  border-top: 1px solid #f0f2f5;
  display: flex;
  gap: 8px;
}

.theme-upload {
  width: 100%;
}
</style>
