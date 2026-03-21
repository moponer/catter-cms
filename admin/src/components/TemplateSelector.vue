<template>
  <div class="template-selector">
    <div class="selector-grid" :class="`grid-${column}`">
      <div
        v-for="template in templates"
        :key="template.value"
        class="template-card"
        :class="{ active: modelValue === template.value, disabled: disabled }"
        @click="selectTemplate(template.value)"
      >
        <div class="template-preview">
          <el-image
            :src="template.preview || defaultPreview"
            fit="cover"
            class="preview-image"
          >
            <template #error>
              <div class="preview-placeholder">
                <el-icon><Picture /></el-icon>
                <span>暂无预览</span>
              </div>
            </template>
          </el-image>

          <div v-if="modelValue === template.value" class="selected-badge">
            <el-icon><Check /></el-icon>
          </div>
        </div>

        <div class="template-info">
          <div class="template-name">{{ template.label }}</div>
          <div v-if="template.description" class="template-desc">
            {{ template.description }}
          </div>
        </div>
      </div>
    </div>

    <div v-if="!templates || templates.length === 0" class="empty-state">
      <el-empty description="暂无可用模板" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Picture, Check } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  templates: {
    type: Array,
    default: () => []
  },
  column: {
    type: Number,
    default: 3
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const defaultPreview = '/images/template-placeholder.jpg'

const selectTemplate = (value) => {
  if (props.disabled) return

  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
.template-selector {
  width: 100%;
}

.selector-grid {
  display: grid;
  gap: 20px;
}

.grid-2 {
  grid-template-columns: repeat(2, 1fr);
}

.grid-3 {
  grid-template-columns: repeat(3, 1fr);
}

.grid-4 {
  grid-template-columns: repeat(4, 1fr);
}

@media (max-width: 1200px) {
  .grid-4 {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 992px) {
  .grid-3, .grid-4 {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .grid-2, .grid-3, .grid-4 {
    grid-template-columns: 1fr;
  }
}

.template-card {
  background: #fff;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.template-card:hover {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
}

.template-card.active {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.template-card.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.template-preview {
  position: relative;
  width: 100%;
  height: 160px;
  background: #f5f7fa;
  overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
}

.preview-placeholder .el-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.selected-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  background: #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
}

.template-info {
  padding: 12px;
}

.template-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.template-desc {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.empty-state {
  padding: 40px;
  text-align: center;
}
</style>
