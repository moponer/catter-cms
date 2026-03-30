<template>
  <div class="repeater-field">
    <div class="repeater-header">
      <span class="repeater-title">{{ label }}</span>
      <el-button type="primary" size="small" @click="addItem">
        <el-icon><Plus /></el-icon> 添加
      </el-button>
    </div>

    <el-empty v-if="!items || items.length === 0" description="暂无数据，点击添加" />

    <draggable v-model="items" item-key="_id" handle=".drag-handle" class="repeater-list">
      <template #item="{ element, index }">
        <div class="repeater-item">
          <div class="item-header" @click="toggleExpand(index)">
            <el-icon class="drag-handle" @click.stop><Rank /></el-icon>
            <span class="item-title">{{ getItemTitle(element, index) }}</span>
            <el-icon class="expand-icon" :class="{ expanded: expandedItems.includes(index) }">
              <ArrowRight />
            </el-icon>
            <el-button type="danger" link @click.stop="removeItem(index)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
          <div class="item-fields" v-show="expandedItems.includes(index)">
            <FieldRenderer
              v-for="f in resolvedFields"
              :key="f.name"
              :field="f"
              :value="element[f.name]"
              @update:modelValue="updateItemField(index, f.name, $event)"
            />
          </div>
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import draggable from 'vuedraggable'
import { Plus, Delete, Rank, ArrowRight } from '@element-plus/icons-vue'
import FieldRenderer from './FieldRenderer.vue'
import { ElMessage } from 'element-plus'

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

// 使用 ref 替代 computed，避免响应式初始化问题
const items = ref([])

// 防止循环更新的标志
let isInternalUpdate = false

// 解析 items 的函数
const parseItems = (val) => {
  let list = val

  // 如果已经是数组，直接使用
  if (Array.isArray(list)) {
    return list.map((item) => {
      if (!item._id) {
        return { ...item, _id: generateId() }
      }
      return item
    })
  }

  // 如果是字符串，尝试 JSON 解析
  if (typeof list === 'string' && list.trim()) {
    try {
      const parsed = JSON.parse(list)
      if (Array.isArray(parsed)) {
        list = parsed
      } else {
        list = []
      }
    } catch {
      list = []
    }
  } else {
    list = []
  }

  return list.map((item) => {
    if (!item._id) {
      return { ...item, _id: generateId() }
    }
    return item
  })
}

// 在 onMounted 中初始化，避免 immediate watch 在 setup 阶段的响应式问题
onMounted(() => {
  items.value = parseItems(props.modelValue)
})

// 监听 props.modelValue 变化
watch(() => props.modelValue, (newVal) => {
  // 如果是内部更新触发的变化，跳过
  if (isInternalUpdate) {
    isInternalUpdate = false
    return
  }
  items.value = parseItems(newVal)
}, { deep: true })

// 展开状态管理
const expandedItems = ref([])

const toggleExpand = (index) => {
  const idx = expandedItems.value.indexOf(index)
  if (idx >= 0) {
    expandedItems.value.splice(idx, 1)
  } else {
    expandedItems.value.push(index)
  }
}

// 解析字段配置
const resolvedFields = computed(() => {
  return props.fields.map(field => {
    if (typeof field === 'string') {
      return { name: field, label: field, type: 'text' }
    }
    return field
  })
})

// 监听 items 长度变化，自动展开第一个项目
watch(() => items.value.length, (newLen) => {
  if (newLen > 0 && expandedItems.value.length === 0) {
    expandedItems.value = [0]
  }
})

// 监听 items 内部变化，emit 更新到父组件
watch(items, (newItems) => {
  const cleanList = newItems.map(({ _id, ...rest }) => rest)
  isInternalUpdate = true
  emit('update:modelValue', cleanList)
}, { deep: true })

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

// 更新某个项目的字段
const updateItemField = (index, fieldName, value) => {
  const newItems = [...items.value]
  newItems[index] = {
    ...newItems[index],
    [fieldName]: value
  }
  items.value = newItems
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

.expand-icon {
  transition: transform 0.2s;
  cursor: pointer;
  color: #909399;
  margin-right: 8px;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.item-header {
  cursor: pointer;
}
</style>
