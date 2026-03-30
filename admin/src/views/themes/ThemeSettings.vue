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

                <!-- 特殊处理：模块排序 -->
                <template v-if="section.section === 'moduleOrder'">
                  <div class="modules-section">
                    <el-alert type="info" :closable="false" style="margin-bottom: 16px;">
                      <i class="fa-solid fa-info-circle"></i> 拖拽调整模块在首页的显示顺序，点击模块编辑内容
                    </el-alert>
                    
                    <draggable 
                      v-model="moduleOrder" 
                      item-key="id"
                      handle=".drag-handle"
                      animation="200"
                      @end="onModuleDragEnd"
                    >
                      <template #item="{ element, index }">
                        <div class="module-card" @click="openModuleEditor(element)">
                          <div class="module-drag drag-handle">
                            <i class="fa-solid fa-grip-vertical"></i>
                          </div>
                          <div class="module-icon">
                            <i :class="getModuleIcon(element.id)"></i>
                          </div>
                          <div class="module-info">
                            <h4>{{ element.name }}</h4>
                            <p class="module-status">
                              <el-tag v-if="getModuleEnabled(element.id)" type="success" size="small">已启用</el-tag>
                              <el-tag v-else type="info" size="small">已禁用</el-tag>
                            </p>
                          </div>
                          <div class="module-actions">
                            <el-button text circle @click.stop="openModuleEditor(element)">
                              <i class="fa-solid fa-edit"></i>
                            </el-button>
                          </div>
                        </div>
                      </template>
                    </draggable>
                  </div>
                </template>

                <!-- 普通字段渲染 -->
                <template v-else>
                  <div class="form-container" role="form" aria-label="主题配置表单">
                    <el-form label-position="top" size="default" class="config-form">
                      <el-form-item
                        v-for="field in getVisibleFields(section.fields)"
                        :key="field.name"
                        :label="field.label"
                        :required="field.required"
                        class="config-form-item"
                      >
                        <FieldRenderer
                          :field="{ ...field, _formItem: true }"
                          :value="configValues[field.name]"
                          @update:modelValue="configValues[field.name] = $event"
                        />
                      </el-form-item>
                    </el-form>
                  </div>
                </template>
              </div>
            </el-tab-pane>
          </el-tabs>
        </template>

        <!-- 无 schema 时显示空状态 -->
        <template v-else>
          <el-empty description="该主题暂无可配置项" />
        </template>

        <!-- 保存按钮 -->
        <div class="save-actions">
          <el-button @click="resetConfig" :disabled="loading">
            <i class="fa-solid fa-undo"></i> 重置
          </el-button>
          <el-button type="primary" @click="saveConfig" :loading="saving" :disabled="loading || !hasChanges">
            <i class="fa-solid fa-save"></i> 保存配置
          </el-button>
        </div>
      </template>
    </div>

    <!-- 模块编辑弹窗 -->
    <el-dialog
      v-model="moduleEditorVisible"
      :title="currentModule ? currentModule.name + ' - 编辑' : '模块编辑'"
      width="800px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <div v-if="currentModule" class="module-editor">
        <div class="form-container">
          <!-- 如果有 repeater 字段，直接渲染整个 repeater -->
          <FieldRenderer
            v-if="currentModuleRepeaterField"
            :field="currentModuleRepeaterField"
            :value="currentModuleValue"
            @update:modelValue="configValues[currentModuleRepeaterField.name] = $event"
          />
          <!-- 如果没有 repeater 字段，显示提示 -->
          <el-empty v-else description="该模块暂无可编辑内容" />
        </div>
      </div>
      
      <template #footer>
        <el-button @click="moduleEditorVisible = false">取消</el-button>
        <el-button type="primary" @click="saveModuleConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import draggable from 'vuedraggable'
import FieldRenderer from '@/components/FieldRenderer.vue'
import { getThemeConfigFull, saveThemeConfig, listThemes } from '@/api/theme'

const route = useRoute()
const router = useRouter()

const slug = computed(() => route.params.slug)
const loading = ref(false)
const saving = ref(false)
const currentTheme = ref(null)
const isActive = ref(false)
const schemaSections = ref([])
const configValues = ref({})
const originalConfig = ref({})
const activeTab = ref('')

// 模块相关
const moduleOrder = ref([])
const originalModuleOrder = ref([])
const moduleEditorVisible = ref(false)
const currentModule = ref(null)

// 判断是否有更改
const hasChanges = computed(() => {
  return JSON.stringify(configValues.value) !== JSON.stringify(originalConfig.value) ||
         JSON.stringify(moduleOrder.value) !== JSON.stringify(originalModuleOrder.value)
})

// 当前模块的 repeater 字段对象（根据 dataKey 动态查找）
const currentModuleRepeaterField = computed(() => {
  if (!currentModule.value) return null

  // 优先使用 dataKey（来自 homeModuleOrder 条目）
  const dataKey = currentModule.value.dataKey
  if (!dataKey) {
    // Hero 和 news 没有独立的 repeater 数据
    return null
  }

  // 在所有 sections 中查找这个 dataKey 对应的 repeater 字段
  for (const section of schemaSections.value) {
    if (section.fields) {
      const field = section.fields.find(f => f.name === dataKey && f.type === 'repeater')
      if (field) {
        return field
      }
    }
  }

  return null
})

// 当前模块的值
const currentModuleValue = computed(() => {
  if (!currentModuleRepeaterField.value) return null
  return configValues.value[currentModuleRepeaterField.value.name]
})

// 获取模块图标
const getModuleIcon = (moduleId) => {
  const icons = {
    hero: 'fa-solid fa-image',
    advantages: 'fa-solid fa-star',
    services: 'fa-solid fa-truck',
    stats: 'fa-solid fa-chart-bar',
    news: 'fa-solid fa-newspaper'
  }
  return icons[moduleId] || 'fa-solid fa-puzzle-piece'
}

// 获取模块启用状态
const getModuleEnabled = (moduleId) => {
  const enabledKey = `${moduleId}Enabled`
  return configValues.value[enabledKey] === true
}

// 打开模块编辑器
const openModuleEditor = (module) => {
  currentModule.value = module
  moduleEditorVisible.value = true
}

// 保存模块配置
const saveModuleConfig = () => {
  moduleEditorVisible.value = false
  currentModule.value = null
  ElMessage.success('模块配置已暂存，请点击「保存配置」按钮提交所有更改')
}

// 模块拖拽结束
const onModuleDragEnd = () => {
  // 自动更新到 configValues
  configValues.value.homeModuleOrder = moduleOrder.value
}

// 评估 show_if 条件
const evaluateCondition = (condition, values) => {
  if (!condition) return true

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

      let compareValue = compareValueRaw.trim()
      if ((compareValue.startsWith('"') && compareValue.endsWith('"')) ||
          (compareValue.startsWith("'") && compareValue.endsWith("'"))) {
        compareValue = compareValue.slice(1, -1)
      } else if (compareValue === 'true') {
        compareValue = true
      } else if (compareValue === 'false') {
        compareValue = false
      } else if (!isNaN(Number(compareValue))) {
        compareValue = Number(compareValue)
      }

      let compareResult = false
      switch (compareOp) {
        case '==':
          compareResult = fieldValue == compareValue
          break
        case '!=':
          compareResult = fieldValue != compareValue
          break
        case '>':
          compareResult = fieldValue > compareValue
          break
        case '<':
          compareResult = fieldValue < compareValue
          break
        case '>=':
          compareResult = fieldValue >= compareValue
          break
        case '<=':
          compareResult = fieldValue <= compareValue
          break
      }

      if (operator === '&&') {
        result = result && compareResult
      } else {
        result = result || compareResult
      }
    }
  }

  return result
}

// 获取可见字段
const getVisibleFields = (fields) => {
  if (!fields) return []
  return fields.filter(field => {
    if (!field.show_if) return true
    return evaluateCondition(field.show_if, configValues.value)
  })
}

// 加载主题数据
const loadThemeData = async () => {
  loading.value = true
  try {
    const fullRes = await getThemeConfigFull(slug.value)
    const fullData = fullRes?.data || {}

    const schema = fullRes?.data?.configForm || []
    schemaSections.value = schema

    // 设置默认 tab
    if (schema.length > 0 && !activeTab.value) {
      activeTab.value = schema[0].section
    }

    currentTheme.value = {
      name: fullData.name || slug.value,
      description: fullData.description || '',
      version: fullData.version || '',
      author: fullData.author || '',
      thumbnail: fullData.thumbnail || ''
    }

    try {
      const listRes = await listThemes()
      const themeList = listRes.data?.list || listRes.list || []
      const activeTheme = themeList.find(t => t.slug === slug.value)
      isActive.value = activeTheme ? activeTheme.is_active : false
    } catch {
      isActive.value = false
    }

    // 保存 homeModuleOrder 字段的默认值（用于空数组时的回退）
    let homeModuleOrderDefault = null

    const values = {}
    for (const section of schema) {
      for (const field of section.fields) {
        // 保存 homeModuleOrder 的默认值
        if (field.name === 'homeModuleOrder' && field.default) {
          homeModuleOrderDefault = field.default
        }

        const savedVal = fullData[field.name]
        if (savedVal !== undefined && savedVal !== null) {
          // repeater 类型需要解析 JSON 字符串
          if (field.type === 'repeater' && typeof savedVal === 'string') {
            try {
              values[field.name] = JSON.parse(savedVal)
            } catch {
              values[field.name] = field.default || []
            }
          } else {
            values[field.name] = savedVal
          }
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

    // 加载模块顺序：如果 homeModuleOrder 是空数组，使用默认值
    if (values.homeModuleOrder && Array.isArray(values.homeModuleOrder) && values.homeModuleOrder.length > 0) {
      moduleOrder.value = values.homeModuleOrder
      originalModuleOrder.value = JSON.parse(JSON.stringify(values.homeModuleOrder))
    } else if (homeModuleOrderDefault && Array.isArray(homeModuleOrderDefault) && homeModuleOrderDefault.length > 0) {
      // 空数组时使用 theme.yml 中的默认值
      moduleOrder.value = homeModuleOrderDefault
      originalModuleOrder.value = JSON.parse(JSON.stringify(homeModuleOrderDefault))
    }
  } catch (error) {
    ElMessage.error('加载主题数据失败')
  } finally {
    loading.value = false
  }
}

// 保存配置
const saveConfig = async () => {
  saving.value = true
  try {
    // 确保 moduleOrder 已同步到 configValues
    configValues.value.homeModuleOrder = moduleOrder.value
    
    await saveThemeConfig(configValues.value, slug.value)
    originalConfig.value = JSON.parse(JSON.stringify(configValues.value))
    originalModuleOrder.value = JSON.parse(JSON.stringify(moduleOrder.value))
    ElMessage.success('配置保存成功')
  } catch (error) {
    ElMessage.error('保存配置失败')
  } finally {
    saving.value = false
  }
}

// 重置配置
const resetConfig = async () => {
  try {
    await ElMessageBox.confirm('确定要重置所有配置吗？这将恢复到上次保存的状态。', '确认重置', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    configValues.value = JSON.parse(JSON.stringify(originalConfig.value))
    if (originalConfig.value.homeModuleOrder) {
      moduleOrder.value = JSON.parse(JSON.stringify(originalConfig.value.homeModuleOrder))
    }
    ElMessage.success('配置已重置')
  } catch {
    // 用户取消
  }
}

// 返回主题列表
const goBack = () => {
  router.push('/themes')
}

onMounted(() => {
  loadThemeData()
})
</script>

<style scoped lang="scss">
.theme-settings {
  padding: 20px;
}

.theme-content {
  margin-top: 20px;
}

.theme-info-card {
  margin-bottom: 20px;
  
  .theme-header {
    display: flex;
    align-items: flex-start;
    gap: 16px;
  }
  
  .theme-meta {
    flex: 1;
    
    h2 {
      margin: 0 0 8px;
      font-size: 20px;
    }
    
    .theme-description {
      margin: 0 0 12px;
      color: #666;
    }
    
    .theme-tags {
      display: flex;
      gap: 8px;
    }
  }
}

.config-tabs {
  margin-bottom: 20px;
}

.tab-content {
  padding: 20px 0;
}

.section-header {
  margin-bottom: 20px;
  
  h3 {
    margin: 0 0 8px;
    font-size: 16px;
  }
  
  .section-description {
    margin: 0;
    color: #666;
    font-size: 14px;
  }
}

.modules-section {
  .module-card {
    display: flex;
    align-items: center;
    padding: 12px;
    margin-bottom: 12px;
    background: #f9fafb;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      border-color: #409eff;
      box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
    }
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .module-drag {
      padding: 8px;
      color: #9ca3af;
      cursor: move;
      
      &:hover {
        color: #409eff;
      }
    }
    
    .module-icon {
      width: 48px;
      height: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #fff;
      border-radius: 8px;
      margin-right: 12px;
      
      i {
        font-size: 24px;
        color: #409eff;
      }
    }
    
    .module-info {
      flex: 1;
      
      h4 {
        margin: 0 0 4px;
        font-size: 15px;
      }
      
      .module-status {
        margin: 0;
      }
    }
    
    .module-actions {
      opacity: 0;
      transition: opacity 0.2s;
    }
    
    &:hover .module-actions {
      opacity: 1;
    }
  }
}

.form-container {
  max-width: 800px;
}

.config-form {
  width: 100%;
}

.config-form-item {
  margin-bottom: 20px;
}

.save-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e5e7eb;
  
  button {
    i {
      margin-right: 6px;
    }
  }
}

.module-editor {
  max-height: 60vh;
  overflow-y: auto;
}
</style>
