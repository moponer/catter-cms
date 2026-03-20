<template>
  <div class="icon-picker">
    <el-popover
      v-model:visible="popoverVisible"
      placement="bottom-start"
      :width="400"
      trigger="click"
    >
      <template #reference>
        <el-input
          v-model="iconName"
          placeholder="选择图标"
          clearable
          @clear="handleClear"
        >
          <template #prefix>
            <el-icon v-if="iconName"><component :is="iconName" /></el-icon>
            <el-icon v-else><Search /></el-icon>
          </template>
        </el-input>
      </template>
      
      <div class="icon-picker-content">
        <el-input
          v-model="searchText"
          placeholder="搜索图标..."
          :prefix-icon="Search"
          size="small"
          style="margin-bottom: 12px;"
        />
        
        <div class="icon-grid">
          <div
            v-for="icon in filteredIcons"
            :key="icon"
            class="icon-item"
            :class="{ selected: iconName === icon }"
            @click="handleSelect(icon)"
            :title="icon"
          >
            <el-icon :size="20"><component :is="icon" /></el-icon>
          </div>
        </div>
        
        <el-empty v-if="filteredIcons.length === 0" description="没有找到匹配的图标" />
      </div>
    </el-popover>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { Search } from '@element-plus/icons-vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const popoverVisible = ref(false)
const searchText = ref('')
const iconName = ref(props.modelValue)

// 获取所有图标名称
const allIcons = Object.keys(ElementPlusIconsVue).filter(key => key !== 'default')

// 过滤图标
const filteredIcons = computed(() => {
  if (!searchText.value) {
    return allIcons.slice(0, 100) // 默认显示前100个
  }
  const keyword = searchText.value.toLowerCase()
  return allIcons.filter(icon => icon.toLowerCase().includes(keyword))
})

// 监听外部值变化
watch(() => props.modelValue, (val) => {
  iconName.value = val
})

// 选择图标
const handleSelect = (icon) => {
  iconName.value = icon
  emit('update:modelValue', icon)
  popoverVisible.value = false
}

// 清除图标
const handleClear = () => {
  iconName.value = ''
  emit('update:modelValue', '')
}
</script>

<style scoped>
.icon-picker {
  width: 100%;
}

.icon-picker-content {
  max-height: 400px;
}

.icon-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
  padding: 8px;
  background: #f5f5f5;
  border-radius: 8px;
}

.icon-item {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.icon-item:hover {
  border-color: #409eff;
  color: #409eff;
}

.icon-item.selected {
  border-color: #409eff;
  background: #ecf5ff;
  color: #409eff;
}
</style>
