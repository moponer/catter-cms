<template>
  <div class="post-edit-page">
    <div class="page-header">
      <h1>{{ isEdit ? '编辑文章' : '新建文章' }}</h1>
      <div class="header-actions">
        <el-button @click="$router.back()">取消</el-button>
        <el-button type="info" @click="handleSaveDraft" :loading="saving">
          保存草稿
        </el-button>
        <el-button type="primary" @click="handlePublish" :loading="saving">
          {{ form.status === 'published' ? '更新发布' : '发布文章' }}
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
                placeholder="请输入文章标题"
                class="title-input"
                size="large"
              />
            </el-form-item>

            <el-form-item label="URL别名（slug）">
              <el-input
                v-model="form.slug"
                placeholder="文章URL，如：my-first-post（留空则自动生成）"
              />
            </el-form-item>

            <el-form-item prop="content">
              <TiptapEditor
                ref="editorRef"
                v-model="form.content"
                placeholder="请输入文章内容..."
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
            <span>文章设置</span>
          </template>

          <el-form :model="form" label-position="top">
            <el-form-item label="状态">
              <el-tag :type="form.status === 'published' ? 'success' : 'info'">
                {{ form.status === 'published' ? '已发布' : '草稿' }}
              </el-tag>
            </el-form-item>

            <el-form-item label="分类">
              <el-select v-model="form.category_id" placeholder="选择分类" clearable style="width: 100%">
                <el-option
                  v-for="cat in categories"
                  :key="cat.id"
                  :label="cat.name"
                  :value="cat.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="标签">
              <el-select
                v-model="form.tag_names"
                multiple
                filterable
                allow-create
                default-first-option
                :reserve-keyword="false"
                placeholder="选择或输入标签"
                style="width: 100%"
              >
                <el-option
                  v-for="tag in tags"
                  :key="tag.id"
                  :label="tag.name"
                  :value="tag.name"
                />
              </el-select>
              <div class="tag-hint">可选择已有标签，或输入新标签名称后按回车创建</div>
            </el-form-item>

            <el-form-item label="摘要">
              <el-input
                v-model="form.excerpt"
                type="textarea"
                :rows="3"
                placeholder="文章摘要（可选）"
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
const categories = ref([])
const tags = ref([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  slug: '',
  content: '',
  excerpt: '',
  status: 'draft',
  category_id: null,
  tag_names: [],
  featured_image: '',
  meta_title: '',
  meta_description: ''
})

const rules = {
  title: [{ required: true, message: '请输入文章标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入文章内容', trigger: 'blur' }]
}

// 加载分类和标签
const loadCategories = async () => {
  try {
    const res = await get('/api/categories')
    categories.value = res.data || res || []
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

const loadTags = async () => {
  try {
    const res = await get('/api/tags')
    // API 返回 {code: 0, data: {list: [...], total: N}}
    tags.value = res.data?.list || res.list || []
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

// 加载文章数据
const loadPost = async () => {
  if (!isEdit.value) return

  try {
    const res = await get(`/api/posts/${route.params.id}`)
    const data = res.data || res
    Object.assign(form, {
      title: data.title || '',
      slug: data.slug || '',
      content: data.content || '',
      excerpt: data.excerpt || '',
      status: data.status || 'draft',
      category_id: data.category_id || data.CategoryID || null,
      // 从 Tag 对象数组提取名称
      tag_names: Array.isArray(data.tags)
        ? data.tags.map(t => typeof t === 'object' ? t.name : t)
        : [],
      // 字段映射：cover → featured_image
      featured_image: data.cover || data.featured_image || '',
      meta_title: data.meta_title || '',
      meta_description: data.meta_description || ''
    })
  } catch (error) {
    ElMessage.error('加载文章失败')
  }
}

// 保存草稿
const handleSaveDraft = async () => {
  form.status = 'draft'
  await savePost()
}

// 发布文章
const handlePublish = async () => {
  form.status = 'published'
  await savePost()
}

// 保存文章
const savePost = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    // 构建请求数据，映射字段名
    const postData = {
      title: form.title,
      slug: form.slug,
      content: form.content,
      excerpt: form.excerpt,
      status: form.status,
      category_id: form.category_id,
      // 使用标签名称（支持创建新标签）
      tag_names: form.tag_names,
      // 字段映射：featured_image → cover
      cover: form.featured_image,
      meta_title: form.meta_title,
      meta_description: form.meta_description
    }

    if (isEdit.value) {
      await put(`/api/posts/${route.params.id}`, postData)
      ElMessage.success('文章更新成功')
    } else {
      await post('/api/posts', postData)
      ElMessage.success('文章创建成功')
    }
    router.push('/posts')
  } catch (error) {
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
    ElMessage.error('图片上传失败')
  }
  return false
}

onMounted(() => {
  loadCategories()
  loadTags()
  loadPost()
})
</script>

<style scoped>
.post-edit-page {
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

.tag-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
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
