# Phase 3: 管理后台 - Repeater和配置UI 原子计划

## 概述
基于Phase 2完成的主题系统，开发管理后台的Repeater字段组件和主题配置UI。

---

## 任务1: Repeater字段Vue组件

### 1.1 创建基础Repeater组件
**文件**: `admin/src/components/RepeaterField.vue`

**功能**:
- 动态添加/删除项目
- 拖拽排序
- 嵌套字段渲染

**实现代码**:
```vue
<template>
  <div class="repeater-field">
    <div class="repeater-header">
      <span class="repeater-title">{{ label }}</span>
      <el-button type="primary" size="small" @click="addItem">
        <el-icon><Plus /></el-icon> 添加
      </el-button>
    </div>

    <el-empty v-if="!items.length" description="暂无数据，点击添加" />

    <draggable v-model="items" item-key="id" handle=".drag-handle" class="repeater-list">
      <template #item="{ element, index }">
        <div class="repeater-item">
          <div class="item-header">
            <el-icon class="drag-handle"><Rank /></el-icon>
            <span class="item-title">{{ getItemTitle(element, index) }}</span>
            <el-button type="danger" link @click="removeItem(index)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
          <div class="item-fields">
            <component
              v-for="field in fields"
              :key="field.name"
              :is="getFieldComponent(field.type)"
              v-model="element[field.name]"
              :label="field.label"
              :config="field"
            />
          </div>
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import draggable from 'vuedraggable'
import { Plus, Delete, Rank } from '@element-plus/icons-vue'
import FieldRenderer from './FieldRenderer.vue'

const props = defineProps({
  modelValue: Array,
  label: String,
  fields: Array,
  min: { type: Number, default: 0 },
  max: { type: Number, default: 0 },
  itemTitle: String
})

const emit = defineEmits(['update:modelValue'])

const items = computed({
  get: () => props.modelValue || [],
  set: (val) => emit('update:modelValue', val)
})

const generateId = () => Date.now() + Math.random().toString(36).substr(2, 9)

const addItem = () => {
  if (props.max && items.value.length >= props.max) {
    ElMessage.warning(`最多只能添加 ${props.max} 项`)
    return
  }
  const newItem = { id: generateId() }
  props.fields.forEach(field => {
    newItem[field.name] = field.default !== undefined ? field.default : null
  })
  items.value.push(newItem)
}

const removeItem = (index) => {
  if (props.min && items.value.length <= props.min) {
    ElMessage.warning(`至少需要保留 ${props.min} 项`)
    return
  }
  items.value.splice(index, 1)
}

const getItemTitle = (item, index) => {
  if (props.itemTitle && item[props.itemTitle]) {
    return item[props.itemTitle]
  }
  return `项目 ${index + 1}`
}

const getFieldComponent = (type) => {
  const componentMap = {
    text: 'el-input',
    textarea: 'el-input',
    number: 'el-input-number',
    select: 'el-select',
    switch: 'el-switch',
    image: 'MediaPicker',
    richtext: 'RichTextEditor',
    color: 'el-color-picker',
    date: 'el-date-picker',
    repeater: 'RepeaterField'
  }
  return componentMap[type] || 'el-input'
}
</script>

<style scoped>
.repeater-field {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  padding: 16px;
  background: #fff;
}

.repeater-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.repeater-title {
  font-weight: 500;
  font-size: 14px;
  color: #606266;
}

.repeater-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.repeater-item {
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background: #f5f7fa;
  overflow: hidden;
}

.item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}

.drag-handle {
  cursor: move;
  color: #909399;
}

.item-title {
  flex: 1;
  font-weight: 500;
  font-size: 14px;
  color: #303133;
}

.item-fields {
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
</style>
```

### 1.2 创建字段渲染器组件
**文件**: `admin/src/components/FieldRenderer.vue`

```vue
<template>
  <div class="field-renderer" :class="[`field-${type}`]">
    <label v-if="label" class="field-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>

    <div class="field-input">
      <!-- Text / Email / URL -->
      <el-input
        v-if="['text', 'email', 'url', 'password'].includes(type)"
        v-model="modelValue"
        :type="type === 'password' ? 'password' : 'text'"
        :placeholder="placeholder"
        :disabled="disabled"
        :clearable="clearable"
      />

      <!-- Textarea -->
      <el-input
        v-else-if="type === 'textarea'"
        v-model="modelValue"
        type="textarea"
        :rows="rows || 4"
        :placeholder="placeholder"
        :disabled="disabled"
        :maxlength="maxlength"
        show-word-limit
      />

      <!-- Number -->
      <el-input-number
        v-else-if="type === 'number'"
        v-model="modelValue"
        :min="min"
        :max="max"
        :step="step || 1"
        :disabled="disabled"
        :controls-position="controlsPosition || 'right'"
      />

      <!-- Select -->
      <el-select
        v-else-if="type === 'select'"
        v-model="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :clearable="clearable !== false"
        :multiple="multiple"
        :filterable="filterable"
        style="width: 100%"
      >
        <el-option
          v-for="opt in options"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>

      <!-- Switch -->
      <el-switch
        v-else-if="type === 'switch'"
        v-model="modelValue"
        :disabled="disabled"
        :active-text="activeText"
        :inactive-text="inactiveText"
      />

      <!-- Color -->
      <el-color-picker
        v-else-if="type === 'color'"
        v-model="modelValue"
        :disabled="disabled"
        :show-alpha="showAlpha !== false"
        :predefine="predefine"
      />

      <!-- Image -->
      <MediaPicker
        v-else-if="type === 'image'"
        v-model="modelValue"
        :disabled="disabled"
        :multiple="multiple"
        :accept="accept || 'image/*'"
      />

      <!-- RichText -->
      <RichTextEditor
        v-else-if="type === 'richtext'"
        v-model="modelValue"
        :disabled="disabled"
        :height="height || 300"
      />

      <!-- Date / DateTime -->
      <el-date-picker
        v-else-if="['date', 'datetime', 'time'].includes(type)"
        v-model="modelValue"
        :type="type === 'datetime' ? 'datetime' : type"
        :placeholder="placeholder"
        :disabled="disabled"
        :clearable="clearable !== false"
        style="width: 100%"
      />

      <!-- Repeater (嵌套) -->
      <RepeaterField
        v-else-if="type === 'repeater'"
        v-model="modelValue"
        :label="null"
        :fields="fields"
        :min="min"
        :max="max"
        :item-title="itemTitle"
      />

      <!-- 未知类型 -->
      <el-alert
        v-else
        :title="`不支持的字段类型: ${type}`"
        type="warning"
        :closable="false"
      />
    </div>

    <div v-if="description" class="field-description">
      {{ description }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import MediaPicker from './MediaPicker.vue'
import RichTextEditor from './RichTextEditor.vue'
import RepeaterField from './RepeaterField.vue'

const props = defineProps({
  modelValue: [String, Number, Boolean, Array, Object, Date],
  type: { type: String, default: 'text' },
  label: String,
  description: String,
  placeholder: String,
  required: Boolean,
  disabled: Boolean,

  // 通用属性
  rows: Number,
  maxlength: Number,
  min: Number,
  max: Number,
  step: Number,
  controlsPosition: String,

  // 选择器属性
  options: Array,
  multiple: Boolean,
  filterable: Boolean,
  clearable: { type: Boolean, default: true },

  // 开关属性
  activeText: String,
  inactiveText: String,

  // 颜色属性
  showAlpha: { type: Boolean, default: true },
  predefine: Array,

  // 图片属性
  accept: String,

  // 富文本属性
  height: Number,

  // Repeater属性
  fields: Array,
  itemTitle: String
})

const emit = defineEmits(['update:modelValue'])

const modelValue = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})
</script>

<style scoped>
.field-renderer {
  margin-bottom: 16px;
}

.field-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #606266;
}

.required {
  color: #f56c6c;
  margin-left: 4px;
}

.field-input {
  width: 100%;
}

.field-input :deep(.el-input__wrapper),
.field-input :deep(.el-textarea__inner) {
  width: 100%;
}

.field-description {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

/* 嵌套Repeater样式调整 */
.field-repeater .field-repeater {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
}
</style>
</code>
</thinking>

<function_results>
为了继续 Phase 3 的开发，我需要了解现有代码结构。让我深入探索管理后台和主题系统的代码。

### 现有代码探索结果

**管理后台 (Vue 3 + Element Plus)**:
- 位于 `admin/src/`
- 使用 Element Plus 组件库
- 已有 MediaPicker, RichTextEditor 组件
- 主题设置页面: `views/themes/ThemeSettings.vue`

**后端主题系统 (Go)**:
- 位于 `gin-app/internal/theme/`
- 核心文件:
  - `module.go` - 模块定义
  - `registry.go` - 模块注册管理
  - `template_funcs.go` - 模板函数
  - `repeater.go` - Repeater字段处理
  - `condition.go` - 条件评估
  - `settings.go` - 主题设置管理
  - `loader.go` - 主题加载
  - `manager.go` - 主题管理器
- API Handler: `gin-app/internal/handlers/theme.go`
- 路由: `gin-app/internal/routes/theme.go`

**主题配置 (YAML)**:
- `theme.yaml` - 主题元数据
- `settings.yaml` - 主题设置定义

现在我将开始 Phase 3 的开发，首先创建 RepeaterField 组件。

```vue
<template>
  <div class="repeater-field">
    <div class="repeater-header">
      <span class="repeater-title">{{ label }}</span>
      <el-button type="primary" size="small" @click="addItem">
        <el-icon><Plus /></el-icon> 添加
      </el-button>
    </div>

    <el-empty v-if="!items.length" description="暂无数据，点击添加" />

    <draggable v-model="items" item-key="id" handle=".drag-handle" class="repeater-list">
      <template #item="{ element, index }">
        <div class="repeater-item">
          <div class="item-header">
            <el-icon class="drag-handle"><Rank /></el-icon>
            <span class="item-title">{{ getItemTitle(element, index) }}</span>
            <el-button type="danger" link @click="removeItem(index)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
          <div class="item-fields">
            <component
              v-for="field in fields"
              :key="field.name"
              :is="getFieldComponent(field.type)"
              v-model="element[field.name]"
              :label="field.label"
              :config="field"
            />
          </div>
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import draggable from 'vuedraggable'
import { Plus, Delete, Rank } from '@element-plus/icons-vue'
import FieldRenderer from './FieldRenderer.vue'

const props = defineProps({
  modelValue: Array,
  label: String,
  fields: Array,
  min: { type: Number, default: 0 },
  max: { type: Number, default: 0 },
  itemTitle: String
})

const emit = defineEmits(['update:modelValue'])

const items = computed({
  get: () => props.modelValue || [],
  set: (val) => emit('update:modelValue', val)
})

const generateId = () => Date.now() + Math.random().toString(36).substr(2, 9)

const addItem = () => {
  if (props.max && items.value.length >= props.max) {
    ElMessage.warning(`最多只能添加 ${props.max} 项`)
    return
  }
  const newItem = { id: generateId() }
  props.fields.forEach(field => {
    newItem[field.name] = field.default !== undefined ? field.default : null
  })
  items.value.push(newItem)
}

const removeItem = (index) => {
  if (props.min && items.value.length <= props.min) {
    ElMessage.warning(`至少需要保留 ${props.min} 项`)
    return
  }
  items.value.splice(index, 1)
}

const getItemTitle = (item, index) => {
  if (props.itemTitle && item[props.itemTitle]) {
    return item[props.itemTitle]
  }
  return `项目 ${index + 1}`
}

const getFieldComponent = (type) => {
  const componentMap = {
    text: 'el-input',
    textarea: 'el-input',
    number: 'el-input-number',
    select: 'el-select',
    switch: 'el-switch',
    image: 'MediaPicker',
    richtext: 'RichTextEditor',
    color: 'el-color-picker',
    date: 'el-date-picker',
    repeater: 'RepeaterField'
  }
  return componentMap[type] || 'el-input'
}
</script>

<style scoped>
.repeater-field {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  padding: 16px;
  background: #fff;
}

.repeater-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.repeater-title {
  font-weight: 500;
  font-size: 14px;
  color: #606266;
}

.repeater-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.repeater-item {
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background: #f5f7fa;
  overflow: hidden;
}

.item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}

.drag-handle {
  cursor: move;
  color: #909399;
}

.item-title {
  flex: 1;
  font-weight: 500;
  font-size: 14px;
  color: #303133;
}

.item-fields {
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
</style>
