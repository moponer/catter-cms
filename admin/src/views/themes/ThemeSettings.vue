<template>
  <div class="theme-settings">
    <el-page-header title="返回主题列表" :content="currentTheme ? currentTheme.name + ' - 主题配置' : '主题配置'" @back="goBack" />

    <div class="theme-content">
      <!-- 主题信息卡片 -->
      <el-card class="theme-info-card" v-if="currentTheme">
        <div class="theme-header">
          <el-avatar v-if="currentTheme.thumbnail" :size="80" :src="currentTheme.thumbnail" shape="square" />
          <div class="theme-meta">
            <h2>{{ currentTheme.name }}</h2>
            <p class="theme-description">{{ currentTheme.description || '暂无描述' }}</p>
            <div class="theme-tags">
              <el-tag v-if="currentTheme.version" size="small">版本 {{ currentTheme.version }}</el-tag>
              <el-tag v-if="currentTheme.author" size="small" type="info">{{ currentTheme.author }}</el-tag>
              <el-tag v-if="isActive" size="small" type="success">当前使用</el-tag>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 加载中 -->
      <el-skeleton v-if="loading" :rows="8" animated />

      <!-- 配置表单 -->
      <template v-else>
        <!-- 有 schema 时：动态渲染 -->
        <template v-if="schemaSections.length > 0">
          <el-tabs v-model="activeTab" type="border-card" class="config-tabs">
            <el-tab-pane
              v-for="section in schemaSections"
              :key="section.section"
              :label="section.label || section.section"
              :name="section.section"
            >
              <div class="tab-content">
                <div class="section-header" v-if="section.description">
                  <h3>{{ section.label || section.section }}</h3>
                  <p v-if="section.description" class="section-description">{{ section.description }}</p>
                </div>

                <!-- 表单容器 - 固定宽度，避免撑满屏幕 -->
                <div class="form-container" role="form" aria-label="主题配置表单">
                  <el-form label-position="top" size="default" class="config-form">
                    <el-form-item
                      v-for="field in getVisibleFields(section.fields)"
                      :key="field.name"
                      :label="field.label"
                      :required="field.required"
                      class="config-form-item"
                      :aria-describedby="field.description ? `desc-${field.name}` : undefined"
                    >
                      <template #label>
                        <span class="field-label">{{ field.label }}</span>
                        <span v-if="field.required" class="required-star" aria-hidden="true">*</span>
                      </template>
                      <FieldRenderer
                        v-model="configValues[field.name]"
                        :type="field.type"
                        :label="field.label"
                        :description="field.description"
                        :placeholder="field.placeholder"
                        :required="field.required"
                        :options="field.options"
                        :min="field.min"
                        :max="field.max"
                        :step="field.step"
                        :rows="field.rows"
                        :show-alpha="field.showAlpha"
                        :fields="field.type === 'repeater' ? (field.fields || field.children || []) : []"
                        :item-title="field.itemTitle"
                      />
                      <div v-if="field.description" :id="`desc-${field.name}`" class="field-help" role="tooltip">
                        {{ field.description }}
                      </div>
                    </el-form-item>
                  </el-form>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </template>

        <!-- 无 schema 时：显示占位 -->
        <template v-else>
          <el-empty description="该主题没有可配置的设置项">
            <el-button type="primary" @click="goBack">返回</el-button>
          </el-empty>
        </template>

        <!-- 保存按钮 -->
        <div class="actions-bar" v-if="schemaSections.length > 0">
          <div class="actions-left">
            <el-button @click="handleExport">
              <el-icon><Download /></el-icon> 导出配置
            </el-button>
            <el-button @click="importDialogVisible = true">
              <el-icon><Upload /></el-icon> 导入配置
            </el-button>
            <el-button @click="handleClearCache" :loading="clearingCache">
              <el-icon><Refresh /></el-icon> 清理缓存
            </el-button>
          </div>
          <div class="actions-right">
            <el-button @click="resetConfig">重置</el-button>
            <el-button type="primary" :loading="saving" @click="saveConfig">
              <el-icon><Check /></el-icon> 保存配置
            </el-button>
          </div>
        </div>
      </template>
    </div>

    <!-- 导入配置对话框 -->
    <el-dialog v-model="importDialogVisible" title="导入配置" width="500px">
      <el-upload
        ref="importUploadRef"
        class="config-import-upload"
        drag
        :auto-upload="false"
        :limit="1"
        accept=".json"
        :on-change="handleImportFileChange"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">拖拽 JSON 文件到此处，或 <em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">请上传之前导出的配置文件（.json 格式）</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="handleImport">导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Download, Upload, UploadFilled, Refresh } from '@element-plus/icons-vue'
import { getThemeConfig, saveThemeConfig, getThemeSettingsSchema, resetThemeConfig, listThemes, getThemeConfigFull, exportThemeConfig, importThemeConfig, clearThemeCache } from '@/api/theme'
import FieldRenderer from '@/components/FieldRenderer.vue'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const saving = ref(false)
const clearingCache = ref(false)
const activeTab = ref('')
const currentTheme = ref(null)
const isActive = ref(false)
const slug = computed(() => route.params.slug || '')

// Schema 数据
const schemaSections = ref([])

// 配置值 - 扁平 key-value
const configValues = ref({})

// 原始配置（用于重置）
const originalConfig = ref({})

// 导入相关
const importDialogVisible = ref(false)
const importing = ref(false)
const importFile = ref(null)
const importUploadRef = ref(null)

onMounted(async () => {
  await loadThemeData()

  // 添加键盘快捷键
  document.addEventListener('keydown', handleKeyDown)
})

// 清理键盘事件监听
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})

// 键盘快捷键处理
function handleKeyDown(e) {
  // Ctrl+S 保存配置
  if (e.ctrlKey && e.key === 's') {
    e.preventDefault()
    if (!loading.value && !saving.value) {
      saveConfig()
    }
  }
  // Ctrl+R 重置配置
  if (e.ctrlKey && e.key === 'r') {
    e.preventDefault()
    if (!loading.value) {
      resetConfig()
    }
  }
  // Esc 关闭导入对话框
  if (e.key === 'Escape' && importDialogVisible.value) {
    importDialogVisible.value = false
  }
}

const loadThemeData = async () => {
  loading.value = true
  try {
    // 1. 获取完整配置（含 schema 定义 + 已保存值 + 主题信息）
    const fullRes = await getThemeConfigFull(slug.value)
    const fullData = fullRes.data || {}

    // 2. 设置 schema（配置表单定义）
    const schema = fullRes.data?.configForm || []
    schemaSections.value = schema

    // 设置默认 tab
    if (schema.length > 0 && !activeTab.value) {
      activeTab.value = schema[0].section
    }

    // 3. 设置主题基本信息
    currentTheme.value = {
      name: fullData.name || slug.value,
      description: fullData.description || '',
      version: fullData.version || '',
      author: fullData.author || '',
      thumbnail: fullData.thumbnail || ''
    }

    // 4. 获取 is_active 状态
    try {
      const listRes = await listThemes()
      const themeList = listRes.data?.list || listRes.list || []
      const activeTheme = themeList.find(t => t.slug === slug.value)
      isActive.value = activeTheme ? activeTheme.is_active : false
    } catch {
      isActive.value = false
    }

    // 5. 初始化表单值：从 schema 的默认值开始，合并已保存的值
    // 注意：后端 theme.Config 现在是纯扁平格式（dot-notation key），与 field.name 格式一致
    const values = {}

    // 先建立字段名到默认值的映射
    for (const section of schema) {
      for (const field of section.fields) {
        // 直接用 field.name（dot-notation）作为 key 从 fullData 获取值
        // 因为 theme.Config 已经是扁平格式
        const savedVal = fullData[field.name]
        if (savedVal !== undefined && savedVal !== null) {
          values[field.name] = savedVal
        } else if (field.default !== undefined && field.default !== null) {
          values[field.name] = field.default
        } else {
          switch (field.type) {
            case 'switch':
              values[field.name] = false
              break
            case 'number':
              values[field.name] = 0
              break
            case 'repeater':
              values[field.name] = []
              break
            default:
              values[field.name] = ''
          }
        }
      }
    }

    configValues.value = values
    originalConfig.value = JSON.parse(JSON.stringify(values))
  } catch (error) {
    ElMessage.error('加载主题数据失败')
  } finally {
    loading.value = false
  }
}

// 评估 show_if 条件，决定字段是否显示
const evaluateCondition = (condition, values) => {
  if (!condition) return true

  // 支持 && 和 || 运算符
  const parts = condition.split(/\s*(\&&\s*|\|\|\s*)/)
  let result = true
  let operator = '&&'

  for (const part of parts) {
    if (part === '&&' || part === '||') {
      operator = part.trim()
    } else if (part.trim()) {
      const match = part.trim().match(/(\w+)\s*(==|!=|<|>|>=|<=)\s*(.+)/)
      if (!match) continue

      const [, fieldName, compareOp, compareValueRaw] = match
      const fieldValue = values[fieldName]

      // 处理字符串值的引号
      let compareValue = compareValueRaw.trim()
      if ((compareValue.startsWith("'") && compareValue.endsWith("'")) ||
          (compareValue.startsWith('"') && compareValue.endsWith('"'))) {
        compareValue = compareValue.slice(1, -1)
      }

      let condResult = false
      if (compareOp === '==') {
        condResult = String(fieldValue) === compareValue
      } else if (compareOp === '!=') {
        condResult = String(fieldValue) !== compareValue
      } else if (compareOp === '>') {
        condResult = Number(fieldValue) > Number(compareValue)
      } else if (compareOp === '>=') {
        condResult = Number(fieldValue) >= Number(compareValue)
      } else if (compareOp === '<') {
        condResult = Number(fieldValue) < Number(compareValue)
      } else if (compareOp === '<=') {
        condResult = Number(fieldValue) <= Number(compareValue)
      }

      if (operator === '&&') {
        result = result && condResult
      } else {
        result = result || condResult
      }
    }
  }

  return result
}

// 获取可见字段列表（过滤掉 show_if 条件不满足的字段）
const getVisibleFields = (fields) => {
  return fields.filter(field => {
    if (!field.show_if) return true
    return evaluateCondition(field.show_if, configValues.value)
  })
}

const saveConfig = async () => {
  saving.value = true
  try {
    // 直接发送扁平 key 配置到后端（后端 theme.Config 现在是纯扁平格式）
    await saveThemeConfig(configValues.value, slug.value)
    ElMessage.success('配置保存成功')
    originalConfig.value = JSON.parse(JSON.stringify(configValues.value))
  } catch (error) {
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

const resetConfig = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重置所有配置吗？这将恢复到系统默认值。',
      '确认重置',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    // 重置为默认值（从 schema 的 default 值构建）
    const defaultValues = {}
    for (const section of schemaSections.value) {
      for (const field of section.fields) {
        if (field.default !== undefined && field.default !== null) {
          defaultValues[field.name] = field.default
        } else {
          switch (field.type) {
            case 'switch':
              defaultValues[field.name] = false
              break
            case 'number':
              defaultValues[field.name] = 0
              break
            case 'repeater':
              defaultValues[field.name] = []
              break
            default:
              defaultValues[field.name] = ''
          }
        }
      }
    }
    configValues.value = defaultValues
    ElMessage.success('配置已重置为默认值')
  } catch {
    // 用户取消
  }
}

// 导出配置
async function handleExport() {
  try {
    const res = await exportThemeConfig(slug.value)
    if (res.data) {
      // 创建下载
      const dataStr = JSON.stringify(res.data, null, 2)
      const blob = new Blob([dataStr], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `theme-config-${slug.value}-${Date.now()}.json`
      link.click()
      URL.revokeObjectURL(url)
      ElMessage.success('配置导出成功')
    }
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 导入文件选择
function handleImportFileChange(file) {
  importFile.value = file.raw
}

// 执行导入
async function handleImport() {
  if (!importFile.value) {
    ElMessage.warning('请选择要导入的配置文件')
    return
  }

  try {
    // 读取文件内容
    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const configData = JSON.parse(e.target.result)
        if (!configData.config || typeof configData.config !== 'object') {
          ElMessage.error('无效的配置文件格式')
          return
        }

        // 确认导入
        await ElMessageBox.confirm(
          '导入将覆盖当前所有配置，是否继续？',
          '确认导入',
          { type: 'warning' }
        )

        importing.value = true
        await importThemeConfig(slug.value, configData.config)
        ElMessage.success('配置导入成功')
        importDialogVisible.value = false
        importFile.value = null
        importUploadRef.value?.clearFiles()
        // 重新加载配置
        await loadThemeData()
      } catch (parseError) {
        ElMessage.error('配置文件解析失败')
      }
    }
    reader.readAsText(importFile.value)
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('导入失败')
    }
  } finally {
    importing.value = false
  }
}

const goBack = () => {
  router.push('/themes')
}

// 清理主题缓存
const handleClearCache = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理主题缓存吗？这将清除模板缓存并重新加载。',
      '确认清理缓存',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    clearingCache.value = true
    await clearThemeCache(slug.value)
    ElMessage.success('缓存清理成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清理缓存失败')
    }
  } finally {
    clearingCache.value = false
  }
}
</script>

<style scoped>
.theme-settings {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.theme-content {
  margin-top: 20px;
}

.theme-info-card {
  margin-bottom: 20px;
}

.theme-header {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.theme-meta h2 {
  margin: 0 0 8px;
  font-size: 20px;
  color: #303133;
}

.theme-description {
  margin: 0 0 12px;
  color: #606266;
  line-height: 1.5;
}

.theme-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.config-tabs {
  margin-bottom: 20px;
}

.tab-content {
  padding: 20px 16px;
}

/* 表单容器 - 限制最大宽度，避免输入框过宽 */
.form-container {
  max-width: 900px;
  margin: 0 auto;
}

.section-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.section-header h3 {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.section-description {
  margin: 0;
  color: #909399;
  font-size: 13px;
}

/* 表单样式 */
.config-form {
  width: 100%;
}

.config-form-item {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.config-form-item:last-child {
  border-bottom: none;
}

/* 字段标签样式 */
.field-label {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.required-star {
  color: #f56c6c;
  margin-left: 4px;
}

.tab-content :deep(.el-form-item__label) {
  font-weight: 500;
  color: #303133;
  line-height: 1.6;
  padding-bottom: 8px !important;
}

.tab-content :deep(.el-form-item__content) {
  line-height: 1.6;
  width: 100%;
}

.tab-content :deep(.el-form-item) {
  display: block;
}

.tab-content :deep(.el-form-item__label-wrap) {
  display: block;
  width: 100%;
}

.field-help {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

.actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.actions-left,
.actions-right {
  display: flex;
  gap: 12px;
}

.config-import-upload {
  width: 100%;
}

/* 焦点管理 - 确保键盘导航时焦点状态清晰可见 */
:focus-visible {
  outline: 2px solid var(--el-color-primary);
  outline-offset: 2px;
}

/* 按钮焦点样式 */
.el-button:focus-visible {
  outline: 2px solid var(--el-color-primary);
  outline-offset: 2px;
}

/* 输入框焦点样式 */
.el-input:focus-visible,
.el-textarea:focus-visible,
.el-select:focus-visible {
  outline: none;
}

/* 表单项目焦点 */
.config-form-item:focus-within {
  background: rgba(64, 158, 255, 0.04);
  border-radius: 4px;
}

/* 跳过链接（可访问性） */
.skip-link {
  position: absolute;
  top: -40px;
  left: 0;
  background: var(--el-color-primary);
  color: white;
  padding: 8px 16px;
  z-index: 9999;
  transition: top 0.3s;
}

.skip-link:focus {
  top: 0;
}
</style>
