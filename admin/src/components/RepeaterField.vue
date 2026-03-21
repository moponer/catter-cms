<template>
  <div class="repeater-field">
    <div class="repeater-header">
      <span class="repeater-title">{{ label }}</span>
      <el-button type="primary" size="small" @click="addItem">
        <el-icon><Plus /></el-icon> 添加
      </el-button>
    </div>

    <el-empty v-if="!items.length" description="暂无数据，点击添加" />

    <draggable v-model="items" item-key="_id" handle=".drag-handle" class="repeater-list">
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
            <FieldRenderer
              v-for="field in resolvedFields"
              :key="field.name"
              v-model="element[field.name]"
              :type="field.type"
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
import { computed, watch } from 'vue'
import draggable from 'vuedraggable'
import { Plus, Delete, Rank } from '@element-plus/icons-vue'
import FieldRenderer from './FieldRenderer.vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  label: {
    type: String,
    default: ''
  },
  fields: {
    type: Array,
    default: () => []
  },
  min: {
    type: Number,
    default: 0
  },
  max: {
    type: Number,
    default: 0
  },
  itemTitle: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

// 确保每个项目都有唯一ID
const items = computed({
  get: () => {
    const list = props.modelValue || []
    return list.map((item, index) => {
      if (!item._id) {
        return { ...item, _id: generateId() }
      }
      return item
    })
  },
  set: (val) => {
    // 移除内部 _id 字段
    const cleanList = val.map(({ _id, ...rest }) => rest)
    emit('update:modelValue', cleanList)
  }
})

// 解析字段配置
const resolvedFields = computed(() => {
  return props.fields.map(field => {
    if (typeof field === 'string') {
      return { name: field, label: field, type: 'text' }
    }
    return field
  })
})

const generateId = () => {
  return Date.now().toString(36) + Math.random().toString(36).substr(2, 5)
}

const addItem = () => {
  if (props.max && items.value.length >= props.max) {
    ElMessage.warning(`最多只能添加 ${props.max} 项`)
    return
  }

  const newItem = { _id: generateId() }
  resolvedFields.value.forEach(field => {
    const defaultValue = field.type === 'switch' ? false :
                        field.type === 'number' ? 0 :
                        field.type === 'select' ? (field.multiple ? [] : null) :
                        field.type === 'repeater' ? [] :
                        null
    newItem[field.name] = field.default !== undefined ? field.default : defaultValue
  })

  items.value = [...items.value, newItem]
}

const removeItem = (index) => {
  if (props.min && items.value.length <= props.min) {
    ElMessage.warning(`至少需要保留 ${props.min} 项`)
    return
  }

  const newList = [...items.value]
  newList.splice(index, 1)
  items.value = newList
}

const getItemTitle = (item, index) => {
  if (props.itemTitle && item[props.itemTitle]) {
    return item[props.itemTitle]
  }

  // 尝试找到合适的标题字段
  const titleFields = ['title', 'name', 'label', 'heading', 'text']
  for (const field of titleFields) {
    if (item[field]) {
      return item[field]
    }
  }

  return `项目 ${index + 1}`
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
