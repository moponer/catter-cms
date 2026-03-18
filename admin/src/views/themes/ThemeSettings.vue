<template>
  <div class="themes-page">
    <div class="page-header">
      <div class="header-left">
        <h1>主题管理</h1>
      </div>
      <div class="header-right">
        <!-- 主题上传 -->
        <el-upload
          class="theme-uploader"
          action="/api/themes/upload"
          :show-file-list="false"
          accept=".zip"
          :headers="uploadHeaders"
          :before-upload="beforeUpload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
        >
          <el-button type="primary" :icon="Upload" :loading="uploading">
            {{ uploading ? '上传中...' : '上传主题' }}
          </el-button>
        </el-upload>
      </div>
    </div>

    <!-- 主题列表 -->
    <el-row :gutter="20" class="theme-grid">
      <el-col :xs="24" :sm="12" :md="8" v-for="theme in themes" :key="theme.id">
        <div class="theme-card-compact" :class="{ active: theme.is_active }">
          <!-- 主题头部 -->
          <div class="theme-header">
            <div class="theme-icon">
              <el-icon :size="28"><Picture /></el-icon>
            </div>
            <div class="theme-title">
              <h3>{{ theme.name }}</h3>
              <el-tag v-if="theme.is_active" type="success" size="small">当前使用</el-tag>
            </div>
          </div>

          <!-- 主题信息 -->
          <div class="theme-body">
            <p class="description">{{ theme.description || '暂无描述' }}</p>
            <div class="theme-meta">
              <span class="version">v{{ theme.version || '1.0.0' }}</span>
              <span class="divider">|</span>
              <span class="author">{{ theme.author || '未知作者' }}</span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="theme-actions-compact">
            <el-button
              v-if="!theme.is_active"
              type="primary"
              size="small"
              @click="handleActivate(theme)"
              :loading="activating === theme.id"
            >
              激活
            </el-button>
            <template v-else>
              <el-button
                type="success"
                size="small"
                @click="handleConfigure(theme)"
              >
                <el-icon class="mr-2"><Setting /></el-icon>
                配置
              </el-button>
            </template>
            <el-button
              size="small"
              @click="handlePreview(theme)"
            >
              <el-icon class="mr-2"><View /></el-icon>
              预览
            </el-button>
            <el-button
              v-if="!theme.is_active"
              type="danger"
              size="small"
              plain
              @click="handleDelete(theme)"
            >
              删除
            </el-button>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 主题配置对话框 - 标签切换布局 -->
    <el-dialog
      v-model="configVisible"
      :title="`${currentTheme?.name} - 主题配置`"
      width="900px"
      class="theme-config-dialog"
      destroy-on-close
      :fullscreen="false"
      top="5vh"
    >
      <div v-if="currentTheme" class="config-container">
        <!-- 标签页切换 - 横向布局 -->
        <el-tabs 
          v-if="configGroups.length > 0" 
          v-model="activeTab" 
          type="border-card" 
          class="config-tabs"
          stretch
          style="display: flex; flex-direction: column;"
        >
          <el-tab-pane
            v-for="group in configGroups"
            :key="group.id"
            :label="group.label"
            :name="group.id"
          >
            <template #label>
              <div class="tab-label">
                <el-icon v-if="group.icon" class="tab-icon">
                  <component :is="group.icon" />
                </el-icon>
                <span>{{ group.label }}</span>
                <el-badge
                  v-if="group.fields && group.fields.length > 0"
                  :value="group.fields.length"
                  class="tab-badge"
                  type="info"
                />
              </div>
            </template>

            <!-- 配置表单 -->
            <el-scrollbar max-height="50vh" class="tab-content">
              <el-form
                :model="configForm"
                label-position="top"
                class="config-form"
              >
                <el-row :gutter="24">
                  <el-col
                    v-for="field in group.fields"
                    :key="field.name"
                    :span="getFieldSpan(field)"
                  >
                    <div class="form-field">
                      <div class="field-header">
                        <label class="field-label">
                          {{ field.label }}
                          <el-tag v-if="field.required" type="danger" size="small">必填</el-tag>
                        </label>
                        <div v-if="field.description" class="field-description">
                          {{ field.description }}
                        </div>
                      </div>

                      <!-- 文本输入 -->
                      <el-input
                        v-if="field.type === 'text'"
                        v-model="configForm[field.name]"
                        :placeholder="field.placeholder || '请输入' + field.label"
                        clearable
                      />

                      <!-- 文本域 -->
                      <el-input
                        v-else-if="field.type === 'textarea'"
                        v-model="configForm[field.name]"
                        type="textarea"
                        :rows="field.rows || 4"
                        :placeholder="field.placeholder || '请输入' + field.label"
                      />

                      <!-- 数字输入 -->
                      <el-input-number
                        v-else-if="field.type === 'number'"
                        v-model="configForm[field.name]"
                        :min="field.min"
                        :max="field.max"
                        :step="field.step || 1"
                        style="width: 100%"
                      />

                      <!-- 颜色选择器 -->
                      <div v-else-if="field.type === 'color'" class="color-picker-wrapper">
                        <el-color-picker
                          v-model="configForm[field.name]"
                          :show-alpha="field.showAlpha !== false"
                          size="large"
                        />
                        <el-input
                          v-model="configForm[field.name]"
                          placeholder="选择颜色"
                          style="margin-left: 12px; flex: 1"
                        />
                      </div>

                      <!-- 开关 -->
                      <div v-else-if="field.type === 'switch'" class="switch-wrapper">
                        <el-switch
                          v-model="configForm[field.name]"
                          :active-text="field.activeText || '开启'"
                          :inactive-text="field.inactiveText || '关闭'"
                        />
                      </div>

                      <!-- 下拉选择 -->
                      <el-select
                        v-else-if="field.type === 'select'"
                        v-model="configForm[field.name]"
                        placeholder="请选择"
                        style="width: 100%"
                        clearable
                      >
                        <el-option
                          v-for="opt in field.options"
                          :key="opt.value"
                          :label="opt.label"
                          :value="opt.value"
                        >
                          <div v-if="opt.description" class="option-item">
                            <span class="option-label">{{ opt.label }}</span>
                            <span class="option-desc">{{ opt.description }}</span>
                          </div>
                          <span v-else>{{ opt.label }}</span>
                        </el-option>
                      </el-select>

                      <!-- 多选 -->
                      <el-select
                        v-else-if="field.type === 'multi-select'"
                        v-model="configForm[field.name]"
                        multiple
                        placeholder="请选择（可多选）"
                        style="width: 100%"
                        collapse-tags
                      >
                        <el-option
                          v-for="opt in field.options"
                          :key="opt.value"
                          :label="opt.label"
                          :value="opt.value"
                        />
                      </el-select>

                      <!-- 单选 -->
                      <el-radio-group
                        v-else-if="field.type === 'radio'"
                        v-model="configForm[field.name]"
                      >
                        <el-radio
                          v-for="opt in field.options"
                          :key="opt.value"
                          :value="opt.value"
                        >
                          {{ opt.label }}
                        </el-radio>
                      </el-radio-group>

                      <!-- 复选框 -->
                      <el-checkbox-group
                        v-else-if="field.type === 'checkbox'"
                        v-model="configForm[field.name]"
                      >
                        <el-checkbox
                          v-for="opt in field.options"
                          :key="opt.value"
                          :label="opt.value"
                        >
                          {{ opt.label }}
                        </el-checkbox>
                      </el-checkbox-group>

                      <!-- 图片上传 -->
                      <div v-else-if="field.type === 'image'" class="image-upload-wrapper">
                        <div v-if="configForm[field.name]" class="image-preview">
                          <img :src="configForm[field.name]" :alt="field.label" />
                          <div class="image-actions">
                            <el-button size="small" @click="previewImage(configForm[field.name])">
                              查看大图
                            </el-button>
                            <el-button size="small" type="danger" @click="configForm[field.name] = ''">
                              删除
                            </el-button>
                          </div>
                        </div>
                        <el-upload
                          v-else
                          :show-file-list="false"
                          :before-upload="(file) => handleConfigImageUpload(file, field.name)"
                          accept="image/*"
                          drag
                        >
                          <el-icon :size="40"><Upload /></el-icon>
                          <div class="upload-text">点击或拖拽上传图片</div>
                        </el-upload>
                      </div>

                      <!-- 日期时间 -->
                      <el-date-picker
                        v-else-if="field.type === 'datetime'"
                        v-model="configForm[field.name]"
                        type="datetime"
                        placeholder="选择日期时间"
                        style="width: 100%"
                        format="YYYY-MM-DD HH:mm:ss"
                        value-format="YYYY-MM-DD HH:mm:ss"
                      />

                      <!-- 日期范围 -->
                      <el-date-picker
                        v-else-if="field.type === 'daterange'"
                        v-model="configForm[field.name]"
                        type="daterange"
                        range-separator="至"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        style="width: 100%"
                        format="YYYY-MM-DD"
                        value-format="YYYY-MM-DD"
                      />

                      <!-- 滑块 -->
                      <div v-else-if="field.type === 'slider'" class="slider-wrapper">
                        <el-slider
                          v-model="configForm[field.name]"
                          :min="field.min || 0"
                          :max="field.max || 100"
                          :step="field.step || 1"
                          :show-stops="field.showStops"
                          :marks="field.marks"
                        />
                        <div class="slider-value">
                          当前值: <strong>{{ configForm[field.name] }}</strong>
                        </div>
                      </div>

                      <!-- 代码编辑器 -->
                      <div v-else-if="field.type === 'code'" class="code-editor-wrapper">
                        <el-input
                          v-model="configForm[field.name]"
                          type="textarea"
                          :rows="field.rows || 10"
                          :placeholder="field.placeholder || '请输入代码'"
                        />
                      </div>

                      <!-- 默认文本输入 -->
                      <el-input
                        v-else
                        v-model="configForm[field.name]"
                        :placeholder="field.placeholder || '请输入'"
                      />
                    </div>
                  </el-col>
                </el-row>
              </el-form>
            </el-scrollbar>
          </el-tab-pane>
        </el-tabs>

        <!-- 无配置项 -->
        <el-empty
          v-else
          description="该主题没有可配置的选项"
          :image-size="200"
        >
          <template #description>
            <div class="empty-config">
              <el-icon :size="60" color="#909399"><Box /></el-icon>
              <p>当前主题暂无配置项</p>
              <p class="sub-hint">您可以联系主题作者添加更多配置选项</p>
            </div>
          </template>
        </el-empty>
      </div>

      <template #footer>
        <div class="dialog-footer-simple">
          <el-button @click="configVisible = false">取消</el-button>
          <el-button @click="handleResetConfig">重置默认</el-button>
          <el-button type="primary" @click="handleSaveConfig" :loading="savingConfig">
            <el-icon class="mr-2"><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 图片预览对话框 -->
    <el-dialog v-model="imagePreviewVisible" title="图片预览" width="60%">
      <img :src="previewImageUrl" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Picture,
  Upload,
  Setting,
  View,
  Box,
  Check,
  User,
  Document,
  Setting as SettingIcon,
  Brush,
  Link
} from '@element-plus/icons-vue'
import { get, post, put, del } from '@/utils/request'

const themes = ref([])
const activating = ref(null)
const configVisible = ref(false)
const currentTheme = ref(null)
const configForm = ref({})
const configFields = ref([])
const configGroups = ref([])
const activeTab = ref('')
const savingConfig = ref(false)
const uploading = ref(false)
const imagePreviewVisible = ref(false)
const previewImageUrl = ref('')

// 上传请求头
const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return {
    Authorization: token ? `Bearer ${token}` : ''
  }
})

// 加载主题列表
const loadThemes = async () => {
  try {
    const res = await get('/api/themes')
    // 后端返回 { list: [...], total: n } 格式
    if (res.data && Array.isArray(res.data.list)) {
      themes.value = res.data.list
    } else if (Array.isArray(res.data)) {
      themes.value = res.data
    } else if (Array.isArray(res)) {
      themes.value = res
    } else {
      themes.value = []
    }
  } catch (error) {
    console.error('加载主题列表失败:', error)
    themes.value = []
  }
}

// 激活主题
const handleActivate = async (theme) => {
  activating.value = theme.id
  try {
    await post(`/api/themes/${theme.slug}/activate`)
    ElMessage.success('主题激活成功')
    loadThemes()
  } catch (error) {
    console.error('激活失败:', error)
    ElMessage.error(error.response?.data?.message || '激活失败')
  } finally {
    activating.value = null
  }
}

// 配置主题
const handleConfigure = async (theme) => {
  currentTheme.value = theme
  try {
    // 获取完整配置，包括配置表单定义
    const res = await get(`/api/themes/${theme.slug}/config/full`)
    
    // axios 拦截器返回的是响应体: { code: 0, message: "success", data: {...} }
    // res 就是 API 返回的数据，res.data 包含 configForm, appearance 等
    console.log('API响应:', res)
    const config = res.data || res || {}

    // 使用主题开发者定义的配置表单（configForm）
    // configForm 对应 theme.yml 中的 configDefinition
    if (config.configForm && Array.isArray(config.configForm) && config.configForm.length > 0) {
      // 直接使用配置表单定义，已经是分组好的
      configGroups.value = config.configForm.map((group, index) => ({
        id: group.section || `group-${index}`,
        label: group.section || '未命名分组',
        icon: getGroupIcon(group.section),
        description: group.description || '',
        fields: group.fields || []
      }))

      // 默认激活第一个标签
      if (configGroups.value.length > 0) {
        activeTab.value = configGroups.value[0].id
      }
      
      ElMessage.success(`成功加载 ${config.configForm.length} 个配置分组`)
    } else {
      // 如果主题没有定义配置表单，显示空状态
      configGroups.value = []
      ElMessage.warning('该主题没有定义配置表单（configForm为空）')
    }

    // 初始化表单值 - 首先获取已保存的配置值，然后用默认值填充未设置的字段
    const formValues = {}
    
    // 1. 首先加载已保存的配置值（从后端返回的config对象中，排除configForm等元数据）
    const savedConfig = { ...config }
    delete savedConfig.configForm  // 删除配置表单定义
    delete savedConfig.name
    delete savedConfig.slug
    delete savedConfig.version
    delete savedConfig.description
    delete savedConfig.author
    
    // 将已保存的配置值放入表单
    Object.keys(savedConfig).forEach(key => {
      if (savedConfig[key] !== undefined && savedConfig[key] !== null) {
        formValues[key] = savedConfig[key]
      }
    })
    
    // 2. 从 configDefinition 的 fields 中提取默认值，填充未设置的字段
    if (config.configForm && Array.isArray(config.configForm)) {
      config.configForm.forEach(group => {
        if (group.fields && Array.isArray(group.fields)) {
          group.fields.forEach(field => {
            // 只在字段没有已保存值时使用默认值
            if (formValues[field.name] === undefined && field.default !== undefined) {
              formValues[field.name] = field.default
            }
          })
        }
      })
    }
    configForm.value = formValues

    configVisible.value = true
  } catch (error) {
    console.error('加载配置失败:', error)
    configGroups.value = []
    configForm.value = {}
    configVisible.value = true
  }
}

// 获取字段布局宽度
const getFieldSpan = (field) => {
  // 根据字段类型返回不同的布局宽度
  const fullWidthTypes = ['textarea', 'code', 'image', 'daterange']
  const halfWidthTypes = ['text', 'number', 'color', 'switch', 'select', 'datetime', 'slider']
  
  if (fullWidthTypes.includes(field.type)) {
    return 24 // 整行
  } else if (halfWidthTypes.includes(field.type)) {
    return 12 // 半行
  } else {
    return 12 // 默认半行
  }
}

// 获取分组图标
const getGroupIcon = (groupName) => {
  const iconMap = {
    '基本设置': SettingIcon,
    '外观设置': Brush,
    '导航设置': Link,
    'SEO设置': Document,
    '作者信息': User
  }
  return iconMap[groupName] || SettingIcon
}

// 预览主题
const handlePreview = (theme) => {
  const previewUrl = `${window.location.origin}/?theme_preview=${theme.slug}`
  window.open(previewUrl, '_blank')
}

// 删除主题
const handleDelete = async (theme) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除主题 "${theme.name}" 吗？删除后无法恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await del(`/api/themes/${theme.slug}`)
    ElMessage.success('主题删除成功')
    loadThemes()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

// 保存配置
const handleSaveConfig = async () => {
  savingConfig.value = true
  try {
    await put(`/api/themes/${currentTheme.value.slug}/config/full`, configForm.value)
    ElMessage.success('配置保存成功')
    configVisible.value = false
  } catch (error) {
    console.error('保存配置失败:', error)
    ElMessage.error(error.response?.data?.message || '保存配置失败')
  } finally {
    savingConfig.value = false
  }
}

// 重置配置
const handleResetConfig = async () => {
  try {
    await ElMessageBox.confirm('确定要重置为默认配置吗？所有自定义配置将丢失。', '重置确认', {
      confirmButtonText: '确定重置',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await post(`/api/themes/${currentTheme.value.slug}/config/reset`)
    ElMessage.success('配置已重置')
    handleConfigure(currentTheme.value)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置失败:', error)
      ElMessage.error('重置失败')
    }
  }
}

// 上传配置图片
const handleConfigImageUpload = async (file, fieldName) => {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    configForm.value[fieldName] = res.data?.path || res.path
    ElMessage.success('图片上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('图片上传失败')
  }
  return false
}

// 预览图片
const previewImage = (url) => {
  previewImageUrl.value = url
  imagePreviewVisible.value = true
}

// 上传前验证
const beforeUpload = (file) => {
  const isZip = file.type === 'application/zip' || file.name.endsWith('.zip')
  if (!isZip) {
    ElMessage.error('只能上传ZIP格式的主题文件')
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }
  uploading.value = true
  return true
}

// 上传成功处理
const handleUploadSuccess = (response) => {
  uploading.value = false
  if (response.code === 0 || response.code === 200) {
    ElMessage.success('主题上传成功')
    loadThemes()
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传失败处理
const handleUploadError = (error) => {
  uploading.value = false
  console.error('上传失败:', error)
  ElMessage.error('主题上传失败')
}

onMounted(() => {
  loadThemes()
})
</script>

<style scoped>
.themes-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.theme-grid {
  margin-top: 16px;
}

.theme-card {
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.theme-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-4px);
}

.theme-card.active {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
}

.theme-preview {
  position: relative;
  height: 200px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.theme-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: rgba(255, 255, 255, 0.8);
  gap: 12px;
}

.theme-overlay {
  position: absolute;
  top: 12px;
  right: 12px;
}

.theme-info {
  padding: 20px;
}

.theme-info h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.theme-info .description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.theme-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #9ca3af;
}

.theme-actions {
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mr-2 {
  margin-right: 6px;
}

/* 紧凑型主题卡片 */
.theme-card-compact {
  background: #ffffff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 2px solid transparent;
  margin-bottom: 20px;
}

.theme-card-compact:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.theme-card-compact.active {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.theme-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.theme-icon {
  width: 48px;
  height: 48px;
  background: #f5f7fa;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
}

.theme-title {
  flex: 1;
}

.theme-title h3 {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.theme-body {
  margin-bottom: 16px;
}

.theme-body .description {
  margin: 0 0 8px 0;
  font-size: 13px;
  color: #6b7280;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.theme-body .theme-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #9ca3af;
}

.theme-body .theme-meta .divider {
  color: #dcdfe6;
}

.theme-actions-compact {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.theme-actions-compact .el-button {
  flex: 1;
  min-width: 60px;
}

/* 主题配置对话框 - 强制横向布局 */
.theme-config-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.theme-config-dialog :deep(.el-dialog) {
  max-width: 95vw !important;
}

.config-container {
  min-height: 400px;
  display: flex;
  flex-direction: column;
}

/* 标签页样式 - 简洁横向布局 */
.config-tabs {
  border-radius: 8px;
  overflow: hidden;
  display: flex !important;
  flex-direction: column !important;
  border: 1px solid #e4e7ed;
}

/* 强制覆盖任何可能的折叠面板样式 */
.config-tabs :deep(.el-collapse),
.config-tabs :deep(.el-collapse-item__header),
.config-tabs :deep(.el-collapse-item__wrap) {
  display: none !important;
}

/* 标签导航 - 简洁背景 */
.config-tabs :deep(.el-tabs__header) {
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
  margin: 0;
  padding: 0;
  flex-shrink: 0;
}

.config-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0 16px;
  overflow: hidden;
}

.config-tabs :deep(.el-tabs__nav-scroll) {
  overflow: hidden;
}

/* 横向排列 */
.config-tabs :deep(.el-tabs__nav) {
  display: flex !important;
  flex-direction: row !important;
  flex-wrap: nowrap !important;
  border: none !important;
}

/* 标签项 - 宽度自适应 */
.config-tabs :deep(.el-tabs__item) {
  flex: 0 0 auto !important;
  background: transparent;
  color: #606266;
  border: none;
  border-bottom: 2px solid transparent;
  padding: 0 24px;
  height: 48px;
  line-height: 48px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  text-align: center;
  white-space: nowrap;
  margin: 0 4px;
}

.config-tabs :deep(.el-tabs__item:first-child) {
  margin-left: 0;
}

.config-tabs :deep(.el-tabs__item:last-child) {
  margin-right: 0;
}

.config-tabs :deep(.el-tabs__item:hover) {
  color: #409eff;
  background: rgba(64, 158, 255, 0.05);
}

.config-tabs :deep(.el-tabs__item.is-active) {
  background: #fff;
  color: #409eff;
  border-bottom-color: #409eff;
  font-weight: 600;
}

.config-tabs :deep(.el-tabs__active-bar) {
  display: none !important;
}

/* 内容区域 */
.config-tabs :deep(.el-tabs__content) {
  flex: 1;
  overflow: auto;
  background: #fff;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-icon {
  font-size: 18px;
}

.tab-badge :deep(.el-badge__content) {
  background-color: rgba(255, 255, 255, 0.9);
  color: #667eea;
  font-weight: 600;
}

.config-tabs :deep(.el-tabs__item.is-active) .tab-badge :deep(.el-badge__content) {
  background-color: #667eea;
  color: #ffffff;
}

.tab-content {
  padding: 24px;
  background: #f9fafb;
  min-height: 300px;
}

.config-groups {
  background: #f9fafb;
  border-radius: 8px;
  padding: 16px;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.group-icon {
  font-size: 20px;
  color: #3b82f6;
}

.group-title {
  font-weight: 600;
  font-size: 16px;
}

.group-badge {
  margin-left: auto;
}

.group-content {
  padding: 16px 0;
}

.config-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-field {
  background: #ffffff;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.field-header {
  margin-bottom: 12px;
}

.field-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  font-size: 14px;
  color: #374151;
}

.field-description {
  margin-top: 4px;
  font-size: 13px;
  color: #6b7280;
  line-height: 1.5;
}

.color-picker-wrapper {
  display: flex;
  align-items: center;
}

.switch-wrapper {
  padding: 8px 0;
}

.image-upload-wrapper {
  margin-top: 8px;
}

.image-upload-wrapper .image-preview {
  position: relative;
  display: inline-block;
}

.image-upload-wrapper .image-preview img {
  max-width: 300px;
  max-height: 200px;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}

.image-actions {
  margin-top: 8px;
  display: flex;
  gap: 8px;
}

.upload-text {
  margin-top: 8px;
  color: #6b7280;
  font-size: 14px;
}

.slider-wrapper {
  padding: 16px 0;
}

.slider-value {
  margin-top: 8px;
  font-size: 13px;
  color: #6b7280;
}

.code-editor-wrapper :deep(.el-textarea__inner) {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
}

.option-item {
  display: flex;
  flex-direction: column;
}

.option-label {
  font-weight: 500;
}

.option-desc {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 2px;
}

.empty-config {
  text-align: center;
  padding: 40px 20px;
}

.empty-config p {
  margin: 12px 0 0 0;
  color: #6b7280;
}

.empty-config .sub-hint {
  font-size: 13px;
  color: #9ca3af;
  margin-top: 8px;
}

.dialog-footer-simple {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .theme-card.active {
    border-width: 1px;
  }

  .theme-preview {
    height: 150px;
  }
}
</style>
