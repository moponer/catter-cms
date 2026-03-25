<template>
  <div class="field-renderer" :class="[`field-${type}`]">
    <div class="field-input-wrapper">
      <!-- Text / Email / URL / Password -->
      <el-input
        v-if="['text', 'email', 'url', 'password'].includes(type)"
        v-model="modelValue"
        :type="type === 'password' ? 'password' : 'text'"
        :placeholder="placeholder"
        :disabled="disabled"
        :clearable="clearable"
        size="default"
        class="field-input-full"
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
        class="field-input-full"
      />

      <!-- Number -->
      <el-input-number
        v-else-if="type === 'number'"
        v-model="modelValue"
        :min="min"
        :max="max"
        :step="step || 1"
        :disabled="disabled"
        controls-position="right"
        class="field-number"
      />

      <!-- Select -->
      <el-select
        v-else-if="type === 'select'"
        v-model="modelValue"
        :placeholder="placeholder || '请选择'"
        :disabled="disabled"
        :clearable="clearable !== false"
        :multiple="multiple"
        :filterable="filterable"
        size="default"
        class="field-input-full"
      >
        <el-option
          v-for="opt in options"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>

      <!-- Multi-Select / Checkbox -->
      <el-select
        v-else-if="type === 'multi-select'"
        v-model="modelValue"
        :placeholder="placeholder || '请选择'"
        :disabled="disabled"
        :clearable="clearable !== false"
        :multiple="true"
        :filterable="filterable"
        size="default"
        class="field-input-full"
      >
        <el-option
          v-for="opt in options"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>

      <!-- Radio -->
      <el-radio-group
        v-else-if="type === 'radio'"
        v-model="modelValue"
        :disabled="disabled"
        class="field-radio-group"
      >
        <el-radio
          v-for="opt in options"
          :key="opt.value"
          :value="opt.value"
        >
          {{ opt.label }}
        </el-radio>
      </el-radio-group>

      <!-- Checkbox -->
      <el-checkbox-group
        v-else-if="type === 'checkbox'"
        v-model="modelValue"
        :disabled="disabled"
        class="field-checkbox-group"
      >
        <el-checkbox
          v-for="opt in options"
          :key="opt.value"
          :value="opt.value"
        >
          {{ opt.label }}
        </el-checkbox>
      </el-checkbox-group>

      <!-- Switch -->
      <el-switch
        v-else-if="type === 'switch'"
        v-model="modelValue"
        :disabled="disabled"
        :active-text="activeText"
        :inactive-text="inactiveText"
      />

      <!-- Color -->
      <div v-else-if="type === 'color'" class="color-picker-wrapper">
        <el-input
          v-model="modelValue"
          :disabled="disabled"
          placeholder="#ffffff"
          class="color-input-main"
        >
          <template #prefix>
            <div class="color-swatch" :style="{ backgroundColor: modelValue || '#000000' }"></div>
          </template>
        </el-input>
        <el-color-picker
          v-model="modelValue"
          :disabled="disabled"
          :show-alpha="showAlpha !== false"
          :predefine="predefine"
          class="color-picker-btn"
        />
        <span v-if="defaultValue" class="color-default-hint">默认: {{ defaultValue }}</span>
      </div>

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

      <!-- Date / DateTime / Time -->
      <el-date-picker
        v-else-if="['date', 'datetime', 'time'].includes(type)"
        v-model="modelValue"
        :type="type === 'datetime' ? 'datetime' : type"
        :placeholder="placeholder || (type === 'datetime' ? '选择日期时间' : type === 'date' ? '选择日期' : '选择时间')"
        :disabled="disabled"
        :clearable="clearable !== false"
        size="default"
        class="field-input-full"
      />

      <!-- Date Range -->
      <el-date-picker
        v-else-if="type === 'daterange'"
        v-model="modelValue"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        :disabled="disabled"
        :clearable="clearable !== false"
        size="default"
        class="field-input-full"
      />

      <!-- Slider -->
      <div v-else-if="type === 'slider'" class="field-slider-wrapper">
        <el-slider
          v-model="modelValue"
          :min="min"
          :max="max"
          :step="step || 1"
          :disabled="disabled"
          :show-input="true"
          class="field-slider"
        />
      </div>

      <!-- Code -->
      <div v-else-if="type === 'code'" class="field-code-wrapper">
        <el-input
          v-model="modelValue"
          type="textarea"
          :rows="rows || 10"
          :placeholder="placeholder || '// 输入代码...'"
          :disabled="disabled"
          class="field-code-input"
        />
      </div>

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
  defaultValue: String, // 默认值提示

  // 图片属性
  accept: String,

  // 富文本属性
  height: Number,

  // Repeater属性
  fields: Array,
  itemTitle: String
})

const emit = defineEmits(['update:modelValue'])

// 颜色类型的空值处理：确保始终返回有效颜色值
const colorValue = computed(() => {
  if (props.type === 'color') {
    const val = props.modelValue
    if (!val || typeof val !== 'string' || val.trim() === '') {
      return '#000000'
    }
    // 确保以 # 开头
    return val.startsWith('#') ? val : '#' + val
  }
  return props.modelValue
})

const modelValue = computed({
  get: () => props.type === 'color' ? colorValue.value : props.modelValue,
  set: (val) => emit('update:modelValue', val)
})
</script>

<style scoped>
.field-renderer {
  width: 100%;
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
