<template>
  <div class="page-edit-page">
    <div class="page-header">
      <h1>{{ isEdit ? '编辑页面' : '新建页面' }}</h1>
      <div class="header-actions">
        <el-button @click="$router.back()">取消</el-button>
        <el-button type="info" @click="handleSaveDraft" :loading="saving">
          保存草稿
        </el-button>
        <el-button type="primary" @click="handlePublish" :loading="saving">
          {{ form.status === 'published' ? '更新发布' : '发布页面' }}
        </el-button>
      </div>
    </div>

    <div class="edit-container">
      <!-- 左侧主内容区 -->
      <div class="main-content">
        <el-card class="editor-card">
          <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
            <el-form-item prop="title">
              <el-input
                v-model="form.title"
                placeholder="请输入页面标题"
                class="title-input"
                size="large"
              />
            </el-form-item>

            <el-form-item prop="slug">
              <el-input
                v-model="form.slug"
                placeholder="页面别名（URL路径，如 about-us）"
                size="default"
              >
                <template #prepend>/</template>
              </el-input>
            </el-form-item>

            <el-form-item prop="content">
              <TiptapEditor
                ref="editorRef"
                v-model="form.content"
                placeholder="请输入页面内容..."
                :height="500"
              />
            </el-form-item>
          </el-form>
        </el-card>
      </div>

      <!-- 右侧设置面板 -->
      <div class="side-panel">
        <el-card class="settings-card">
          <template #header>
            <span>页面设置</span>
          </template>

          <el-form :model="form" label-position="top">
            <el-form-item label="状态">
              <el-tag :type="form.status === 'published' ? 'success' : 'info'">
                {{ form.status === 'published' ? '已发布' : '草稿' }}
              </el-tag>
            </el-form-item>

            <el-form-item label="模板">
              <el-select v-model="form.template" placeholder="选择模板" style="width: 100%">
                <el-option label="默认模板" value="default" />
                <el-option label="关于页面" value="about" />
                <el-option label="联系页面" value="contact" />
                <el-option label="全宽页面" value="full-width" />
              </el-select>
            </el-form-item>

            <el-form-item label="父页面">
              <el-select 
                v-model="form.parent_id" 
                placeholder="选择父页面（可选）" 
                clearable 
                style="width: 100%"
              >
                <el-option
                  v-for="page in pageOptions"
                  :key="page.id"
                  :label="page.title"
                  :value="page.id"
                  :disabled="page.id === form.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="排序">
              <el-input-number 
                v-model="form.order" 
                :min="0" 
                style="width: 100%"
              />
            </el-form-item>

            <el-form-item label="摘要">
              <el-input
                v-model="form.excerpt"
                type="textarea"
                :rows="3"
                placeholder="页面摘要（可选）"
              />
            </el-form-item>

            <el-form-item label="特色图片">
              <div class="featured-image">
                <div v-if="form.featured_image" class="image-preview">
                  <img :src="form.featured_image" alt="特色图片" />
                  <el-button type="danger" size="small" @click="form.featured_image = ''">
                    移除
                  </el-button>
                </div>
                <el-upload
                  v-else
                  :show-file-list="false"
                  :before-upload="handleImageUpload"
                  accept="image/*"
                >
                  <el-button type="primary" plain>上传图片</el-button>
                </el-upload>
              </div>
            </el-form-item>

            <el-form-item label="SEO设置">
              <el-input
                v-model="form.meta_title"
                placeholder="SEO标题（可选）"
                style="margin-bottom: 10px"
              />
              <el-input
                v-model="form.meta_description"
                type="textarea"
                :rows="2"
                placeholder="SEO描述（可选）"
              />
            </el-form-item>
          </el-form>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import TiptapEditor from '@/components/TiptapEditor.vue'
import { get, post, put } from '@/utils/request'

const route = useRoute()
const router = useRouter()

const formRef = ref(null)
const editorRef = ref(null)
const saving = ref(false)
const pageOptions = ref([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  slug: '',
  content: '',
  excerpt: '',
  status: 'draft',
  template: 'default',
  parent_id: null,
  order: 0,
  featured_image: '',
  meta_title: '',
  meta_description: ''
})

const rules = {
  title: [{ required: true, message: '请输入页面标题', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入页面别名', trigger: 'blur' }],
  content: [{ required: true, message: '请输入页面内容', trigger: 'blur' }]
}

// 加载页面列表（用于父页面选择）
const loadPageOptions = async () => {
  try {
    const res = await get('/api/pages', { page_size: 100 })
    pageOptions.value = res.data?.list || res.list || []
  } catch (error) {
    console.error('加载页面列表失败:', error)
  }
}

// 加载页面数据
const loadPage = async () => {
  if (!isEdit.value) return

  try {
    const res = await get(`/api/pages/${route.params.id}`)
    const data = res.data || res
    Object.assign(form, {
      title: data.title || '',
      slug: data.slug || '',
      content: data.content || '',
      excerpt: data.excerpt || '',
      status: data.status || 'draft',
      template: data.template || 'default',
      parent_id: data.parent_id || null,
      order: data.order || 0,
      featured_image: data.featured_image || '',
      meta_title: data.meta_title || '',
      meta_description: data.meta_description || ''
    })
  } catch (error) {
    console.error('加载页面失败:', error)
    ElMessage.error('加载页面失败')
  }
}

// 保存草稿
const handleSaveDraft = async () => {
  form.status = 'draft'
  await savePage()
}

// 发布页面
const handlePublish = async () => {
  form.status = 'published'
  await savePage()
}

// 保存页面
const savePage = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    if (isEdit.value) {
      await put(`/api/pages/${route.params.id}`, form)
      ElMessage.success('页面更新成功')
    } else {
      await post('/api/pages', form)
      ElMessage.success('页面创建成功')
    }
    router.push('/pages')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 上传特色图片
const handleImageUpload = async (file) => {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    form.featured_image = res.data?.path || res.path
    ElMessage.success('图片上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('图片上传失败')
  }
  return false
}

onMounted(() => {
  loadPageOptions()
  loadPage()
})
</script>

<style scoped>
.page-edit-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.page-header h1 {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.edit-container {
  display: flex;
  gap: var(--spacing-md);
  flex: 1;
  min-height: 0;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.editor-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.editor-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: var(--spacing-md);
}

.title-input :deep(.el-input__wrapper) {
  font-size: 24px;
  font-weight: 600;
  padding: 12px 0;
  box-shadow: none !important;
  border-bottom: 1px solid var(--el-border-color);
  border-radius: 0;
}

.title-input :deep(.el-input__wrapper):focus-within {
  border-bottom-color: var(--el-color-primary);
}

.title-input :deep(.el-input__inner) {
  font-size: 24px;
}

.side-panel {
  width: 320px;
  flex-shrink: 0;
}

.settings-card {
  position: sticky;
  top: var(--spacing-md);
}

.featured-image .image-preview {
  text-align: center;
}

.featured-image .image-preview img {
  max-width: 100%;
  max-height: 150px;
  border-radius: 4px;
  margin-bottom: 10px;
}

@media (max-width: 1024px) {
  .edit-container {
    flex-direction: column;
  }

  .side-panel {
    width: 100%;
  }
}
</style>
