<template>
  <div class="module-manager">
    <div class="manager-header">
      <el-input
        v-model="searchQuery"
        placeholder="搜索模块..."
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="enableAll">全部启用</el-button>
      <el-button @click="disableAll">全部禁用</el-button>
    </div>

    <draggable
      v-model="moduleList"
      item-key="id"
      handle=".drag-handle"
      class="module-list"
      @end="onDragEnd"
    >
      <template #item="{ element, index }">
        <div
          class="module-card"
          :class="{ disabled: !element.enabled, expanded: expandedModules.includes(element.id) }"
        >
          <div class="module-main">
            <el-icon class="drag-handle"><Rank /></el-icon>

            <div class="module-icon">
              <el-icon :size="24">
                <component :is="getModuleIcon(element.id)" />
              </el-icon>
            </div>

            <div class="module-info">
              <div class="module-name">{{ element.name }}</div>
              <div class="module-desc">{{ element.description }}</div>
            </div>

            <div class="module-actions">
              <el-switch
                v-model="element.enabled"
                @change="(val) => onEnableChange(element, val)"
              />
              <el-button
                type="primary"
                link
                @click="toggleExpand(element.id)"
              >
                <el-icon>
                  <ArrowDown v-if="expandedModules.includes(element.id)" />
                  <ArrowRight v-else />
                </el-icon>
                配置
              </el-button>
            </div>
          </div>

          <div v-if="expandedModules.includes(element.id)" class="module-config">
            <el-divider />
            <div class="config-content">
              <h4>模块配置</h4>
              <el-form label-position="top" size="small">
                <el-row :gutter="16">
                  <el-col :span="12">
                    <el-form-item label="显示标题">
                      <el-input v-model="element.config.title" placeholder="模块标题" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="副标题">
                      <el-input v-model="element.config.subtitle" placeholder="副标题" />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="背景图片">
                  <MediaPicker v-model="element.config.backgroundImage" type="image" />
                </el-form-item>

                <el-row :gutter="16">
                  <el-col :span="12">
                    <el-form-item label="背景颜色">
                      <el-color-picker v-model="element.config.backgroundColor" show-alpha />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="文字颜色">
                      <el-color-picker v-model="element.config.textColor" show-alpha />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="内边距">
                  <el-slider v-model="element.config.padding" :max="100" show-input />
                </el-form-item>
              </el-form>
            </div>
          </div>
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import draggable from 'vuedraggable'
import { Rank, Search, ArrowDown, ArrowRight } from '@element-plus/icons-vue'
import {
  Monitor, Service, TrendCharts, UserFilled, News, Phone,
  Picture, Star, Grid, Collection, Calendar, Link
} from '@element-plus/icons-vue'
import MediaPicker from './MediaPicker.vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  availableModules: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const searchQuery = ref('')
const expandedModules = ref([])

const iconMap = {
  hero: Monitor,
  services: Service,
  stats: TrendCharts,
  team: UserFilled,
  news: News,
  contact: Phone,
  gallery: Picture,
  features: Star,
  portfolio: Grid,
  testimonials: Collection,
  blog: Calendar,
  social: Link
}

const moduleList = computed({
  get: () => props.modelValue || [],
  set: (val) => {
    emit('update:modelValue', val)
  }
})

const filteredModules = computed(() => {
  if (!searchQuery.value) return moduleList.value
  const query = searchQuery.value.toLowerCase()
  return moduleList.value.filter(m =>
    m.name.toLowerCase().includes(query) ||
    m.description?.toLowerCase().includes(query)
  )
})

const getModuleIcon = (id) => {
  return iconMap[id] || Monitor
}

const toggleExpand = (id) => {
  const index = expandedModules.value.indexOf(id)
  if (index > -1) {
    expandedModules.value.splice(index, 1)
  } else {
    expandedModules.value.push(id)
  }
}

const onEnableChange = (module, enabled) => {
  emit('change', { module: module.id, enabled, type: 'enable' })
  ElMessage.success(`${module.name}已${enabled ? '启用' : '禁用'}`)
}

const onDragEnd = () => {
  emit('change', { type: 'reorder', modules: moduleList.value.map(m => m.id) })
}

const enableAll = () => {
  moduleList.value.forEach(m => {
    m.enabled = true
  })
  emit('change', { type: 'enableAll' })
  ElMessage.success('所有模块已启用')
}

const disableAll = () => {
  moduleList.value.forEach(m => {
    m.enabled = false
  })
  emit('change', { type: 'disableAll' })
  ElMessage.success('所有模块已禁用')
}

watch(() => props.availableModules, (newVal) => {
  if (newVal && newVal.length > 0 && moduleList.value.length === 0) {
    moduleList.value = newVal.map(m => ({
      ...m,
      enabled: m.enabled !== false,
      config: m.config || {}
    }))
  }
}, { immediate: true })
</script>

<style scoped>
.module-manager {
  background: #fff;
  border-radius: 8px;
}

.manager-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.search-input {
  width: 240px;
}

.module-list {
  padding: 16px;
}

.module-card {
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.3s;
}

.module-card:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.module-card.disabled {
  opacity: 0.6;
  background: #f5f7fa;
}

.module-card.expanded {
  border-color: #409eff;
}

.module-main {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
}

.drag-handle {
  cursor: move;
  color: #909399;
  font-size: 18px;
}

.module-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
  border-radius: 8px;
  color: #606266;
}

.module-info {
  flex: 1;
}

.module-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.module-desc {
  font-size: 13px;
  color: #909399;
}

.module-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.module-config {
  padding: 0 16px 16px;
}

.config-content {
  background: #f5f7fa;
  border-radius: 6px;
  padding: 16px;
}

.config-content h4 {
  margin: 0 0 16px;
  font-size: 14px;
  color: #606266;
}
</style>
