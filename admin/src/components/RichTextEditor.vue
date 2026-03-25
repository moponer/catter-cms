<template>
  <div class="rich-text-editor-wrapper">
    <Editor
      v-model="content"
      :init="editorConfig"
      :disabled="disabled"
      @onBlur="handleBlur"
      @onFocus="handleFocus"
    />
    <MediaPicker v-model="mediaPickerVisible" @select="handleMediaSelect" />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import Editor from '@tinymce/tinymce-vue'
import 'tinymce/tinymce'
import 'tinymce/themes/silver'
import 'tinymce/icons/default'
import 'tinymce/models/dom'
import MediaPicker from './MediaPicker.vue'

// 插件
import 'tinymce/plugins/advlist'
import 'tinymce/plugins/autolink'
import 'tinymce/plugins/lists'
import 'tinymce/plugins/link'
import 'tinymce/plugins/image'
import 'tinymce/plugins/charmap'
import 'tinymce/plugins/preview'
import 'tinymce/plugins/anchor'
import 'tinymce/plugins/searchreplace'
import 'tinymce/plugins/visualblocks'
import 'tinymce/plugins/code'
import 'tinymce/plugins/fullscreen'
import 'tinymce/plugins/insertdatetime'
import 'tinymce/plugins/media'
import 'tinymce/plugins/table'
import 'tinymce/plugins/help'
import 'tinymce/plugins/wordcount'
import 'tinymce/plugins/emoticons'
import 'tinymce/plugins/codesample'

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
    type: [String, Number],
    default: 500
  },
  disabled: {
    type: Boolean,
    default: false
  },
  toolbar: {
    type: String,
    default: 'undo redo | blocks | bold italic underline strikethrough | forecolor backcolor | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | link image media table | codesample code fullscreen preview help'
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'blur', 'focus'])

const content = ref(props.modelValue)
const mediaPickerVisible = ref(false)
let currentFilePickerCallback = null
let currentValue = null
let currentMeta = null

// TinyMCE 配置
const editorConfig = computed(() => ({
  height: props.height,
  promotion: false,
  branding: false,
  // 添加空许可证密钥以禁用警告（社区版功能仍然可用）
  license_key: '',
  menubar: 'file edit view insert format tools table help',
  menu: {
    file: { title: '文件', items: 'newdocument restoredraft | preview' },
    edit: { title: '编辑', items: 'undo redo | cut copy paste pastetext | selectall | searchreplace' },
    view: { title: '视图', items: 'code | visualaid visualblocks visualchars | preview fullscreen' },
    insert: { title: '插入', items: 'image link media codesample inserttable | charmap emoticons hr' },
    format: { title: '格式', items: 'bold italic underline strikethrough superscript subscript codeformat | styles blocks fontfamily fontsize align lineheight | forecolor backcolor | removeformat' },
    tools: { title: '工具', items: 'code wordcount' },
    table: { title: '表格', items: 'inserttable | cell row column | tableprops deletetable' },
    help: { title: '帮助', items: 'help' }
  },
  plugins: [
    'advlist', 'autolink', 'lists', 'link', 'image', 'charmap', 'preview',
    'anchor', 'searchreplace', 'visualblocks', 'code', 'fullscreen',
    'insertdatetime', 'media', 'table', 'help', 'wordcount', 'emoticons',
    'codesample'
  ],
  toolbar: props.toolbar,
  toolbar_mode: 'sliding',
  content_style: `
    body { 
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      font-size: 14px;
      line-height: 1.6;
      padding: 16px;
    }
    img { max-width: 100%; height: auto; }
    pre { background-color: #f5f5f5; padding: 12px; border-radius: 4px; overflow-x: auto; }
    code { background-color: #f5f5f5; padding: 2px 6px; border-radius: 3px; font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace; }
    blockquote { border-left: 4px solid #e0e0e0; margin: 16px 0; padding-left: 16px; color: #666; }
    table { border-collapse: collapse; width: 100%; }
    table td, table th { border: 1px solid #ddd; padding: 8px; }
    table th { background-color: #f5f5f5; }
  `,
  skin: 'oxide',
  language: 'zh_CN',
  placeholder: props.placeholder,
  statusbar: true,
  elementpath: true,
  resize: true,
  // 图片上传配置
  images_upload_handler: async (blobInfo, progress) => {
    try {
      const formData = new FormData()
      formData.append('file', blobInfo.blob(), blobInfo.filename())
      
      const response = await fetch('/api/media/upload', {
        method: 'POST',
        body: formData
      })
      
      const result = await response.json()
      return result.data?.url || result.url || result.path
    } catch (error) {
      throw new Error('图片上传失败')
    }
  },
  // 文件选择回调
  file_picker_callback: (callback, value, meta) => {
    if (meta.filetype === 'image') {
      // 保存回调信息
      currentFilePickerCallback = callback
      currentValue = value
      currentMeta = meta
      
      // 打开媒体库选择器
      mediaPickerVisible.value = true
    }
  },
  // 禁用自动保存
  autosave_ask_before_unload: false
}))

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  if (newVal !== content.value) {
    content.value = newVal
  }
})

// 监听内部值变化
watch(content, (newVal) => {
  emit('update:modelValue', newVal)
  emit('change', newVal)
})

const handleBlur = (e) => {
  emit('blur', e)
}

const handleFocus = (e) => {
  emit('focus', e)
}

// 媒体选择处理
const handleMediaSelect = (media) => {
  if (currentFilePickerCallback) {
    const url = media.url || media.path
    currentFilePickerCallback(url, { title: media.filename })
    currentFilePickerCallback = null
    currentValue = null
    currentMeta = null
  }
}

// 暴露方法
defineExpose({
  getContent: () => content.value,
  setContent: (val) => { content.value = val },
  focus: () => {
    // TinyMCE focus logic
  }
})
</script>

<style scoped>
.rich-text-editor-wrapper {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
}

:deep(.tox-tinymce) {
  border-radius: 8px;
  border: 1px solid var(--el-border-color);
}

:deep(.tox-tinymce:focus-within) {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-9);
}
</style>
