<template>
  <div class="field-renderer" :class="[`field-${fieldType}`]">
    <label v-if="fieldLabel && !fieldHideLabel" class="field-label">{{ fieldLabel }}</label>
    <div class="field-input-wrapper">
      <!-- Text / Email / URL / Password -->
      <el-input
        v-if="['text', 'email', 'url', 'password'].includes(fieldType)"
        v-model="fieldValue"
        :type="fieldType === 'password' ? 'password' : 'text'"
        :placeholder="fieldPlaceholder"
        :disabled="fieldDisabled"
        :clearable="clearable"
        size="default"
        class="field-input-full"
      />

      <!-- Textarea -->
      <el-input
        v-else-if="fieldType === 'textarea'"
        v-model="fieldValue"
        type="textarea"
        :rows="fieldRows || 4"
        :placeholder="fieldPlaceholder"
        :disabled="fieldDisabled"
        :maxlength="maxlength"
        show-word-limit
        class="field-input-full"
      />

      <!-- Number -->
      <el-input-number
        v-else-if="fieldType === 'number'"
        v-model="fieldValue"
        :min="fieldMin"
        :max="fieldMax"
        :step="fieldStep || 1"
        :disabled="fieldDisabled"
        controls-position="right"
        class="field-number"
      />

      <!-- Select -->
      <el-select
        v-else-if="fieldType === 'select'"
        v-model="fieldValue"
        :placeholder="fieldPlaceholder || '请选择'"
        :disabled="fieldDisabled"
        :clearable="clearable !== false"
        :multiple="multiple"
        :filterable="filterable"
        size="default"
        class="field-input-full"
      >
        <el-option
          v-for="opt in fieldOptions"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>

      <!-- Multi-Select / Checkbox -->
      <el-select
        v-else-if="fieldType === 'multi-select'"
        v-model="fieldValue"
        :placeholder="fieldPlaceholder || '请选择'"
        :disabled="fieldDisabled"
        :clearable="clearable !== false"
        :multiple="true"
        :filterable="filterable"
        size="default"
        class="field-input-full"
      >
        <el-option
          v-for="opt in fieldOptions"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>

      <!-- Radio -->
      <el-radio-group
        v-else-if="fieldType === 'radio'"
        v-model="fieldValue"
        :disabled="fieldDisabled"
        class="field-radio-group"
      >
        <el-radio
          v-for="opt in fieldOptions"
          :key="opt.value"
          :value="opt.value"
        >
          {{ opt.label }}
        </el-radio>
      </el-radio-group>

      <!-- Checkbox -->
      <el-checkbox-group
        v-else-if="fieldType === 'checkbox'"
        v-model="fieldValue"
        :disabled="fieldDisabled"
        class="field-checkbox-group"
      >
        <el-checkbox
          v-for="opt in fieldOptions"
          :key="opt.value"
          :value="opt.value"
        >
          {{ opt.label }}
        </el-checkbox>
      </el-checkbox-group>

      <!-- Switch -->
      <el-switch
        v-else-if="fieldType === 'switch'"
        v-model="fieldValue"
        :disabled="fieldDisabled"
        :active-text="fieldActiveText"
        :inactive-text="fieldInactiveText"
      />

      <!-- Icon (FontAwesome) -->
      <IconPicker
        v-else-if="fieldType === 'icon'"
        v-model="fieldValue"
        :placeholder="fieldPlaceholder"
      />

      <!-- Color -->
      <div v-else-if="fieldType === 'color'" class="color-picker-wrapper">
        <el-input
          v-model="fieldValue"
          :disabled="fieldDisabled"
          placeholder="#ffffff"
          class="color-input-main"
        >
          <template #prefix>
            <div class="color-swatch" :style="{ backgroundColor: fieldValue || '#000000' }"></div>
          </template>
        </el-input>
        <el-color-picker
          v-model="fieldValue"
          :disabled="fieldDisabled"
          :show-alpha="showAlpha !== false"
          :predefine="predefine"
          class="color-picker-btn"
        />
        <span v-if="fieldDefault" class="color-default-hint">默认: {{ fieldDefault }}</span>
      </div>

      <!-- Image -->
      <ImageUploadField
        v-else-if="fieldType === 'image'"
        v-model="fieldValue"
        :disabled="fieldDisabled"
      />

      <!-- Image List (Carousel) -->
      <ImageListUploadField
        v-else-if="fieldType === 'image-list'"
        v-model="fieldValue"
        :disabled="fieldDisabled"
      />

      <!-- RichText -->
      <RichTextEditor
        v-else-if="fieldType === 'richtext'"
        v-model="fieldValue"
        :disabled="fieldDisabled"
        :height="height || 300"
      />

      <!-- Date / DateTime / Time -->
      <el-date-picker
        v-else-if="['date', 'datetime', 'time'].includes(fieldType)"
        v-model="fieldValue"
        :type="fieldType === 'datetime' ? 'datetime' : fieldType"
        :placeholder="fieldPlaceholder || (fieldType === 'datetime' ? '选择日期时间' : fieldType === 'date' ? '选择日期' : '选择时间')"
        :disabled="fieldDisabled"
        :clearable="clearable !== false"
        size="default"
        class="field-input-full"
      />

      <!-- Date Range -->
      <el-date-picker
        v-else-if="fieldType === 'daterange'"
        v-model="fieldValue"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        :disabled="fieldDisabled"
        :clearable="clearable !== false"
        size="default"
        class="field-input-full"
      />

      <!-- Slider -->
      <div v-else-if="fieldType === 'slider'" class="field-slider-wrapper">
        <el-slider
          v-model="fieldValue"
          :min="fieldMin"
          :max="fieldMax"
          :step="fieldStep || 1"
          :disabled="fieldDisabled"
          :show-input="true"
          class="field-slider"
        />
      </div>

      <!-- Code -->
      <div v-else-if="fieldType === 'code'" class="field-code-wrapper">
        <el-input
          v-model="fieldValue"
          type="textarea"
          :rows="fieldRows || 10"
          :placeholder="fieldPlaceholder || '// 输入代码...'"
          :disabled="fieldDisabled"
          class="field-code-input"
        />
      </div>

      <!-- Repeater (嵌套) -->
      <RepeaterField
        v-else-if="fieldType === 'repeater'"
        v-model="fieldValue"
        :label="null"
        :fields="fieldFields"
        :min="fieldMin"
        :max="fieldMax"
        :item-title="fieldItemTitle"
      />

      <!-- 未知类型 -->
      <el-alert
        v-else
        :title="`不支持的字段类型: ${fieldType}`"
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
import ImageUploadField from './ImageUploadField.vue'
import ImageListUploadField from './ImageListUploadField.vue'
import IconPicker from './IconPicker.vue'

const props = defineProps({
  // 支持两种用法：
  // 1. 传入 field 对象（推荐）
  field: Object,
  // 2. 直接传入各个属性（兼容旧代码）
  modelValue: [String, Number, Boolean, Array, Object, Date],
  value: [String, Number, Boolean, Array, Object, Date], // 兼容 :value
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
  defaultValue: String, // 默认值提示

  // 图片属性
  accept: String,

  // 富文本属性
  height: Number,

  // Repeater属性
  fields: Array,
  itemTitle: String
})

// 统一处理：优先使用 field 对象，否则使用单独的 props
const fieldType = computed(() => props.field?.type || props.type || 'text')
const fieldValue = computed({
  get: () => {
    // 父组件传入 :value="configValues[field.name]"
    // 优先级：props.value (传入值) > props.modelValue > field.value > ''
    if (props.value !== undefined) return props.value
    if (props.modelValue !== undefined) return props.modelValue
    if (props.field && 'value' in props.field) return props.field.value
    return ''
  },
  set: (val) => emit('update:modelValue', val)
})

// 从 field 对象提取其他属性
const fieldLabel = computed(() => props.field?.label || props.label)
const fieldHideLabel = computed(() => props.field?._hideLabel === true || props.field?._formItem === true)
const fieldPlaceholder = computed(() => props.field?.placeholder || props.placeholder)
const fieldRequired = computed(() => props.field?.required || props.required)
const fieldDisabled = computed(() => props.field?.disabled || props.disabled)
const fieldOptions = computed(() => props.field?.options || props.options || [])
const fieldMin = computed(() => props.field?.min ?? props.min)
const fieldMax = computed(() => props.field?.max ?? props.max)
const fieldStep = computed(() => props.field?.step ?? props.step)
const fieldRows = computed(() => props.field?.rows ?? props.rows)
const fieldFields = computed(() => props.field?.fields || props.fields || [])
const fieldItemTitle = computed(() => props.field?.itemTitle || props.itemTitle)
const fieldDefault = computed(() => props.field?.default || props.defaultValue)
const fieldActiveText = computed(() => props.field?.activeText || props.activeText)
const fieldInactiveText = computed(() => props.field?.inactiveText || props.inactiveText)

const emit = defineEmits(['update:modelValue'])

</script>

<style scoped>
.field-renderer {
  width: 100%;
}

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 8px;
  line-height: 1.4;
}

.field-input-wrapper {
  width: 100%;
}

/* 通用：输入框全宽 */
.field-input-full {
  width: 100%;
}

.field-input-full :deep(.el-input__wrapper) {
  width: 100%;
}

.field-input-full :deep(.el-textarea__inner) {
  width: 100%;
}

/* Number 组件 */
.field-number {
  width: 160px;
}

.field-number :deep(.el-input__wrapper) {
  width: 100%;
}

/* Radio 组 */
.field-radio-group {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

/* Checkbox 组 */
.field-checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

/* Color Picker */
.color-picker-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.color-input-main {
  width: 140px;
  font-family: 'Fira Code', 'Cascadia Code', Consolas, monospace;
}

.color-input-main :deep(.el-input__prefix) {
  margin-right: 4px;
}

.color-swatch {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
  flex-shrink: 0;
}

.color-picker-btn {
  flex-shrink: 0;
}

.color-default-hint {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
}

/* Slider */
.field-slider-wrapper {
  width: 100%;
  padding-right: 120px;
}

.field-slider {
  width: 100%;
}

/* Code 编辑器 */
.field-code-wrapper {
  width: 100%;
}

.field-code-input {
  width: 100%;
}

.field-code-input :deep(textarea) {
  font-family: 'Fira Code', 'Cascadia Code', Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  background: #1e1e1e;
  color: #d4d4d4;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  padding: 12px;
  resize: vertical;
  min-height: 120px;
}

.field-code-input :deep(.el-textarea__inner) {
  background: #1e1e1e;
  color: #d4d4d4;
}

/* 字段描述 */
.field-description {
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

/* 嵌套 Repeater 样式调整 */
.field-repeater .field-repeater {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
}
</style>
