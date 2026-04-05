<template>
  <div class="icon-picker">
    <!-- 当前选中的图标显示 -->
    <div class="icon-preview" @click="dialogVisible = true">
      <i v-if="modelValue" :class="modelValue"></i>
      <span v-else class="placeholder-text">点击选择图标</span>
    </div>

    <!-- 图标选择对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="选择图标"
      width="700px"
      :destroy-on-close="true"
    >
      <!-- 搜索框 -->
      <el-input
        v-model="searchKeyword"
        placeholder="搜索图标..."
        prefix-icon="Search"
        clearable
        class="icon-search"
      />

      <!-- 分类标签 -->
      <div class="category-tabs">
        <el-tag
          v-for="cat in categoryList"
          :key="cat"
          :type="selectedCategory === cat ? 'primary' : 'info'"
          class="category-tag"
          @click="selectedCategory = cat"
        >
          {{ cat }}
        </el-tag>
      </div>

      <!-- 图标网格 -->
      <div class="icon-grid">
        <div
          v-for="icon in filteredIcons"
          :key="icon.id"
          class="icon-item"
          :class="{ selected: isIconSelected(icon.id) }"
          @click="selectIcon(icon)"
          :title="icon.label"
        >
          <i :class="getIconClass(icon.id)"></i>
        </div>
      </div>

      <!-- 空状态 -->
      <el-empty v-if="filteredIcons.length === 0" description="未找到匹配的图标" />

      <template #footer>
        <div class="dialog-footer">
          <span v-if="modelValue" class="selected-info">
            当前: <i :class="modelValue"></i> {{ modelValue }}
          </span>
          <el-button @click="dialogVisible = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import iconsData from '@/assets/fontawesome-icons.json'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const dialogVisible = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref('常用')

// 判断是否为品牌图标（id 以 brands/ 开头）
const isBrandIcon = (iconId) => iconId && iconId.startsWith('brands/')

// 获取图标的完整 class（包含前缀）
const getIconClass = (iconId) => {
  if (!iconId) return ''
  if (isBrandIcon(iconId)) {
    // brands/xxx → fa-brands fa-xxx
    return `fa-brands ${iconId.replace('brands/', 'fa-')}`
  }
  // 如果已经包含 fa- 前缀，直接返回
  if (iconId.startsWith('fa-')) {
    return `fa-solid ${iconId}`
  }
  return `fa-solid fa-${iconId}`
}

// 判断 modelValue 是否选中（兼容普通 fa-xxx 和品牌 fa-brands fa-xxx 格式）
const isIconSelected = (iconId) => {
  if (isBrandIcon(iconId)) {
    // brands/weixin → 比较 fa-brands fa-weixin
    return props.modelValue === `fa-brands ${iconId.replace('brands/', 'fa-')}`
  }
  // 普通 → 比较 fa-xxx
  return props.modelValue === `fa-${iconId}`
}

// 分类列表
const categoryList = computed(() => {
  return Object.keys(iconsData)
})

// 常用图标（优先显示）
const commonIcons = computed(() => {
  const common = iconsData['常用'] || []
  return common.map(i => ({ ...i, id: i.id }))
})

// 所有其他图标（排除常用）
const otherIcons = computed(() => {
  const result = []
  for (const [cat, icons] of Object.entries(iconsData)) {
    if (cat !== '常用') {
      for (const icon of icons) {
        // 排除已经在常用中的
        if (!(commonIcons.value.some(c => c.id === icon.id))) {
          result.push({ ...icon, category: cat })
        }
      }
    }
  }
  return result
})

// 过滤后的图标列表
const filteredIcons = computed(() => {
  let icons = []

  // 如果选择了分类
  if (selectedCategory.value === '常用') {
    icons = commonIcons.value
  } else if (selectedCategory.value === '全部') {
    icons = [...commonIcons.value, ...otherIcons.value]
  } else if (iconsData[selectedCategory.value]) {
    icons = iconsData[selectedCategory.value].map(i => ({ ...i }))
  }

  // 搜索过滤
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.toLowerCase()
    icons = icons.filter(icon => {
      const idMatch = icon.id.toLowerCase().includes(keyword)
      const labelMatch = (icon.label || '').toLowerCase().includes(keyword)
      const searchMatch = (icon.searchTerms || []).some(term =>
        term.toLowerCase().includes(keyword)
      )
      return idMatch || labelMatch || searchMatch
    })
  }

  return icons
})

// 选择图标
const selectIcon = (icon) => {
  if (isBrandIcon(icon.id)) {
    // brands/xxx → 输出 fa-brands fa-xxx（完整格式，模板中可直接用作 class）
    emit('update:modelValue', `fa-brands ${icon.id.replace('brands/', 'fa-')}`)
  } else {
    // 普通 → 输出 fa-xxx（保持向后兼容）
    emit('update:modelValue', `fa-${icon.id}`)
  }
  dialogVisible.value = false
  searchKeyword.value = ''
}
</script>

<style scoped>
.icon-picker {
  width: 100%;
}

.icon-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  cursor: pointer;
  min-height: 32px;
  transition: border-color 0.2s;
}

.icon-preview:hover {
  border-color: #409eff;
}

.icon-preview i {
  font-size: 18px;
  color: #409eff;
}

.placeholder-text {
  color: #909399;
  font-size: 13px;
}

.icon-search {
  margin-bottom: 12px;
}

.category-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.category-tag {
  cursor: pointer;
  user-select: none;
}

.icon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(50px, 1fr));
  gap: 8px;
  max-height: 350px;
  overflow-y: auto;
  padding: 4px;
}

.icon-item {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-item:hover {
  border-color: #409eff;
  background-color: #f5f7fa;
}

.icon-item.selected {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.icon-item i {
  font-size: 20px;
  color: #606266;
}

.icon-item.selected i,
.icon-item:hover i {
  color: #409eff;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.selected-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
}

.selected-info i {
  color: #409eff;
}
</style>
