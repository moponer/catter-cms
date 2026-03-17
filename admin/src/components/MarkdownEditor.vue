<template>
  <div class="markdown-editor-wrapper">
    <MdEditor
      v-model="content"
      :language="language"
      :theme="theme"
      :toolbars="toolbars"
      :placeholder="placeholder"
      :style="{ height: height }"
      @change="handleChange"
      @onUploadImg="handleUploadImage"
      :preview="preview"
      noKatex
      noMermaid
    />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { post } from '@/utils/request'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '请输入内容...'
  },
  height: {
    type: String,
    default: '500px'
  },
  preview: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const content = ref(props.modelValue)
const language = ref('zh-CN')
const theme = ref('light')

// 完整的工具栏配置，支持富文本编辑
const toolbars = [
  'bold',
  'underline',
  'italic',
  'strikeThrough',
  '-',
  'title',
  'sub',
  'sup',
  'quote',
  'unorderedList',
  'orderedList',
  'task',
  '-',
  'codeRow',
  'code',
  'link',
  'image',
  'table',
  '-',
  'revoke',
  'next',
  'save',
  '=',
  'pageFullscreen',
  'fullscreen',
  'preview',
  'htmlPreview',
  'catalog'
]

// 图片上传处理
const handleUploadImage = async (files, callback) => {
  try {
    const file = files[0]
    const formData = new FormData()
    formData.append('file', file)

    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    const imageUrl = res.data?.path || res.path
    callback([imageUrl])
  } catch (error) {
    console.error('图片上传失败:', error)
    // 返回失败回调，编辑器会显示错误
    callback([])
  }
}

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  content.value = newVal
})

// 内容变化
const handleChange = (value) => {
  emit('update:modelValue', value)
  emit('change', value)
}

// 暴露方法
defineExpose({
  getHtml: () => content.value,
  setHtml: (val) => { content.value = val }
})
</script>

<style scoped>
.markdown-editor-wrapper {
  width: 100%;
}

:deep(.md-editor) {
  border-radius: 8px;
  border: 1px solid var(--el-border-color);
}

:deep(.md-editor:focus-within) {
  border-color: var(--el-color-primary);
}
</style>
