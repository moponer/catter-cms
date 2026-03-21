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
        controls-position="right"
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
