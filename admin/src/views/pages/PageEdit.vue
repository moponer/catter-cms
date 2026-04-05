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
        <!-- 动态模板编辑区 -->
        <el-card class="editor-card">
          <!-- 默认图文编辑器（默认模板 或 没有 fields 的自定义模板） -->
          <div v-if="showDefaultEditor" class="default-editor">
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
          </div>

          <!-- 自定义模板编辑器（有 fields 配置的模板，动态渲染） -->
          <div v-else class="custom-template-editor">
            <el-form :model="form" label-position="top">
              <!-- 标题 & Slug 永远显示 -->
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
                  :placeholder="`页面别名（URL路径，如 ${currentTemplateConfig.type}）`"
                  size="default"
                >
                  <template #prepend>/</template>
                </el-input>
              </el-form-item>

              <!-- 模板专属字段区 -->
              <div class="template-fields">
                <h3 class="fields-title">
                  <i :class="currentTemplateConfig.icon" style="margin-right: 6px;"></i>
                  {{ currentTemplateConfig.name }}内容
                </h3>

                <el-form-item
                  v-for="field in currentTemplateConfig.fields"
                  :key="field.name"
                  :label="field.label"
                >
                  <!-- image-list：图片上传 + 媒体库 -->
                  <ImageListUploadField
                    v-if="field.type === 'image-list'"
                    v-model="form.templateDataObj[field.name]"
                    :placeholder="field.placeholder || '点击上传图片'"
                  />

                  <!-- select：下拉选择 -->
                  <el-select
                    v-else-if="field.type === 'select'"
                    v-model="form.templateDataObj[field.name]"
                    :placeholder="field.label"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="opt in field.options"
                      :key="opt.value"
                      :label="opt.label || opt.value"
                      :value="opt.value"
                    />
                  </el-select>

                  <!-- switch：开关 -->
                  <el-switch
                    v-else-if="field.type === 'switch'"
                    v-model="form.templateDataObj[field.name]"
                  />

                  <!-- textarea：多行文本 -->
                  <el-input
                    v-else-if="field.type === 'textarea'"
                    v-model="form.templateDataObj[field.name]"
                    type="textarea"
                    :rows="4"
                    :placeholder="field.placeholder || field.label"
                  />

                  <!-- richtext：富文本编辑器 -->
                  <TiptapEditor
                    v-else-if="field.type === 'richtext'"
                    v-model="form.templateDataObj[field.name]"
                    :placeholder="field.placeholder || field.label"
                    :height="400"
                  />

                  <!-- 默认：单行文本 -->
                  <el-input
                    v-else
                    v-model="form.templateDataObj[field.name]"
                    :placeholder="field.placeholder || field.label"
                  />
                </el-form-item>
              </div>
            </el-form>
          </div>
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
              <el-select 
                v-model="form.template" 
                placeholder="选择模板" 
                style="width: 100%"
                @change="handleTemplateChange"
              >
                <el-option label="默认模板" value="default" />
                <el-option
                  v-for="tpl in pageTemplates"
                  :key="tpl.type"
                  :label="tpl.name"
                  :value="tpl.type"
                >
                  <span>{{ tpl.name }}</span>
                  <span class="template-desc">{{ tpl.description }}</span>
                </el-option>
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

        <!-- 侧边栏设置 -->
        <el-card class="sidebar-card" v-if="currentTemplateConfig?.hasSidebar || form.template === 'default'">
          <template #header>
            <span>侧边栏设置</span>
          </template>
          <el-form :model="form" label-position="top">
            <el-form-item label="启用侧边栏">
              <el-switch v-model="sidebarConfig.enabled" @change="handleSidebarChange" />
              <div class="form-tip">启用后页面将显示侧边栏模块</div>
            </el-form-item>

            <template v-if="sidebarConfig.enabled">
              <el-form-item label="侧边栏位置">
                <el-radio-group v-model="sidebarConfig.position" @change="handleSidebarChange">
                  <el-radio value="right">右侧</el-radio>
                  <el-radio value="left">左侧</el-radio>
                </el-radio-group>
              </el-form-item>

              <el-form-item label="侧边栏模块">
                <div class="sidebar-modules">
                  <div v-if="!sidebarModules.length" class="empty-tip">
                    暂无可用模块
                  </div>
                  <div
                    v-for="mod in sidebarModules"
                    :key="mod.id"
                    class="sidebar-module-item"
                  >
                    <el-checkbox
                      :model-value="isModuleEnabled(mod.id)"
                      @change="(val) => toggleModule(mod.id, val)"
                    >
                      <i :class="mod.icon" style="margin-right: 4px;"></i>
                      {{ mod.name }}
                    </el-checkbox>
                  </div>
                </div>
              </el-form-item>
            </template>
          </el-form>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import TiptapEditor from '@/components/TiptapEditor.vue'
import ImageListUploadField from '@/components/ImageListUploadField.vue'
import { get, post, put } from '@/utils/request'
import { getPageTemplates } from '@/api/theme'

const route = useRoute()
const router = useRouter()

const formRef = ref(null)
const editorRef = ref(null)
const saving = ref(false)
const pageOptions = ref([])
const pageTemplates = ref([])

const isEdit = computed(() => !!route.params.id)

// 模板数据对象（用于编辑）
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
  meta_description: '',
  template_data: '',
  templateDataObj: {}
})

const rules = {
  title: [{ required: true, message: '请输入页面标题', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入页面别名', trigger: 'blur' }],
  content: [{ required: true, message: '请输入页面内容', trigger: 'blur' }]
}

// 当前模板配置
const currentTemplateConfig = computed(() => {
  if (!form.template || form.template === 'default') return null
  return pageTemplates.value.find(t => t.type === form.template) || null
})

// 是否显示默认图文编辑器
// 条件：默认模板 或 当前模板没有 fields（如联系我们）
const showDefaultEditor = computed(() => {
  if (!currentTemplateConfig.value) return true
  const fields = currentTemplateConfig.value.fields
  return !fields || fields.length === 0
})

// 侧边栏模块列表（从主题配置获取）
const sidebarModules = ref([])

// 侧边栏配置
const sidebarConfig = reactive({
  enabled: false,
  position: 'right',
  modules: [] // [{id, enabled, sort}]
})

// 检查模块是否启用
const isModuleEnabled = (moduleId) => {
  const mod = sidebarConfig.modules.find(m => m.id === moduleId)
  return mod ? mod.enabled : false
}

// 切换模块启用状态
const toggleModule = (moduleId, enabled) => {
  const existingMod = sidebarConfig.modules.find(m => m.id === moduleId)
  if (existingMod) {
    existingMod.enabled = enabled
  } else {
    sidebarConfig.modules.push({
      id: moduleId,
      enabled,
      sort: sidebarConfig.modules.length + 1
    })
  }
  handleSidebarChange()
}

// 侧边栏配置变更处理
const handleSidebarChange = () => {
  // 将侧边栏配置保存到 templateDataObj.sidebar 中
  form.templateDataObj.sidebar = {
    enabled: sidebarConfig.enabled,
    position: sidebarConfig.position,
    modules: sidebarConfig.modules.filter(m => m.enabled)
  }
}

// 加载侧边栏模块列表
const loadSidebarModules = async () => {
  try {
    const res = await getPageTemplates()
    if (res.code === 0 && res.data) {
      sidebarModules.value = res.data.sidebarModules || []
      // 如果有模板默认配置，应用它
      const defaults = res.data.templateSidebarDefaults || {}
      const templateType = form.template === 'default' ? 'page' : form.template
      if (defaults[templateType]) {
        const defaultConfig = defaults[templateType]
        sidebarConfig.position = defaultConfig.position || 'right'
        // 从模块列表初始化默认启用的模块
        if (defaultConfig.modules && defaultConfig.modules.length > 0) {
          sidebarConfig.modules = defaultConfig.modules.map(m => ({
            id: m.id,
            enabled: m.enabled,
            sort: m.sort
          }))
        }
      }
    }
  } catch (error) {
    console.error('加载侧边栏模块失败:', error)
  }
}

// 初始化侧边栏配置（从 templateData 中读取）
const initSidebarConfig = () => {
  const sidebarData = form.templateDataObj.sidebar
  if (sidebarData) {
    sidebarConfig.enabled = sidebarData.enabled || false
    sidebarConfig.position = sidebarData.position || 'right'
    if (sidebarData.modules && sidebarData.modules.length > 0) {
      sidebarConfig.modules = sidebarData.modules
    }
  }
}

// 加载页面模板列表
const loadPageTemplates = async () => {
  try {
    const res = await getPageTemplates()
    if (res.code === 0 && res.data) {
      pageTemplates.value = res.data.templates || []
    }
  } catch (error) {
    console.error('加载页面模板失败:', error)
    pageTemplates.value = []
  }
}

// 模板切换时重置 templateDataObj
const handleTemplateChange = (newTemplate) => {
  // 获取新模板的字段默认值
  const tpl = pageTemplates.value.find(t => t.type === newTemplate)
  if (tpl && tpl.fields && tpl.fields.length > 0) {
    const defaults = {}
    tpl.fields.forEach(field => {
      if (field.type === 'image-list') {
        defaults[field.name] = ''  // ImageListUploadField 接受 JSON 字符串
      } else if (field.type === 'select' || field.type === 'radio') {
        defaults[field.name] = field.default || (field.options && field.options[0]?.value) || ''
      } else if (field.type === 'switch') {
        defaults[field.name] = field.default === true || field.default === 'true'
      } else {
        defaults[field.name] = field.default || ''
      }
    })
    Object.assign(form.templateDataObj, defaults)
  } else {
    // 无 fields 的模板（如联系我们），重置为空对象
    form.templateDataObj = {}
  }
}

// 序列化 templateDataObj 为 JSON 字符串
const serializeTemplateData = () => {
  try {
    // 只序列化有值的字段
    const data = {}
    Object.keys(form.templateDataObj).forEach(key => {
      const val = form.templateDataObj[key]
      // 跳过 null/undefined
      if (val === null || val === undefined) return
      // image-list 字段的值是 JSON 字符串（由 ImageListUploadField 产出），解析成数组再存
      if (typeof val === 'string' && val.startsWith('[')) {
        try {
          const arr = JSON.parse(val)
          if (Array.isArray(arr) && arr.length > 0) {
            data[key] = arr
          }
          // 空数组直接跳过
        } catch {
          // 解析失败，跳过
        }
        return
      }
      // 跳过空字符串
      if (val === '') return
      // 跳过空数组
      if (Array.isArray(val) && val.length === 0) return
      data[key] = val
    })
    form.template_data = Object.keys(data).length > 0 ? JSON.stringify(data) : ''
  } catch (e) {
    form.template_data = ''
  }
}

// 反序列化 template_data 字符串为对象
// image-list 字段在 DB 中存的是数组，但 ImageListUploadField 需要 JSON 字符串
const deserializeTemplateData = (jsonStr) => {
  if (!jsonStr) {
    form.templateDataObj = {}
    return
  }
  try {
    const obj = typeof jsonStr === 'string' ? JSON.parse(jsonStr) : jsonStr
    const result = {}
    Object.keys(obj).forEach(key => {
      const val = obj[key]
      // 数组类型（image-list），转成 JSON 字符串供 ImageListUploadField 使用
      if (Array.isArray(val)) {
        result[key] = JSON.stringify(val)
      } else {
        result[key] = val
      }
    })
    form.templateDataObj = result
  } catch (e) {
    form.templateDataObj = {}
  }
}

// 加载页面列表（用于父页面选择）
const loadPageOptions = async () => {
  try {
    const res = await get('/api/pages', { page_size: 100 })
    pageOptions.value = res.data?.list || res.list || []
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

// 加载页面数据
const loadPage = async () => {
  if (!isEdit.value) {
    // 新建页面时检查 URL 参数是否有模板
    if (route.query.template) {
      form.template = route.query.template
      handleTemplateChange(form.template)
    }
    return
  }

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
      meta_description: data.meta_description || '',
      template_data: data.template_data || ''
    })
    // 反序列化 template_data
    deserializeTemplateData(form.template_data)
    // 初始化侧边栏配置
    initSidebarConfig()
  } catch (error) {
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
    await formRef.value?.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    // 序列化模板数据
    serializeTemplateData()

    const payload = {
      title: form.title,
      slug: form.slug,
      content: form.content,
      excerpt: form.excerpt,
      status: form.status,
      template: form.template,
      template_data: form.template_data,
      parent_id: form.parent_id,
      sort_order: form.order,
      featured_image: form.featured_image,
      meta_title: form.meta_title,
      meta_description: form.meta_description
    }

    if (isEdit.value) {
      await put(`/api/pages/${route.params.id}`, payload)
      ElMessage.success('页面更新成功')
    } else {
      await post('/api/pages', payload)
      ElMessage.success('页面创建成功')
    }
    router.push('/pages')
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
  loadPageTemplates()
  loadPageOptions()
  loadSidebarModules()
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
  overflow-y: auto;
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

/* 模板字段区域 */
.template-fields {
  margin-top: 20px;
  padding: 16px;
  background: var(--el-fill-color-light);
  border-radius: 8px;
}

.fields-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.template-desc {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-left: 8px;
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

/* 侧边栏设置卡片 */
.sidebar-card {
  margin-top: var(--spacing-md);
}

.sidebar-card :deep(.el-card__header) {
  background: var(--el-fill-color-light);
}

.sidebar-modules {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sidebar-module-item {
  padding: 8px 12px;
  background: var(--el-fill-color-light);
  border-radius: 6px;
  transition: background var(--transition-fast);
}

.sidebar-module-item:hover {
  background: var(--el-fill-color);
}

.empty-tip {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  text-align: center;
  padding: 20px 0;
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
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
