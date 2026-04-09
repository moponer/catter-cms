<template>
  <div class="tiptap-editor-wrapper">
    <!-- 工具栏 -->
    <div class="editor-toolbar" v-if="editor">
      <div class="toolbar-group">
        <button type="button" @click="editor.chain().focus().undo().run()" :disabled="!editor.can().undo()" title="撤销">
          ↩
        </button>
        <button type="button" @click="editor.chain().focus().redo().run()" :disabled="!editor.can().redo()" title="重做">
          ↪
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <div class="toolbar-group">
        <button type="button" @click="editor.chain().focus().toggleHeading({ level: 1 }).run()" :class="{ active: editor.isActive('heading', { level: 1 }) }" title="标题1">
          H1
        </button>
        <button type="button" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ active: editor.isActive('heading', { level: 2 }) }" title="标题2">
          H2
        </button>
        <button type="button" @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ active: editor.isActive('heading', { level: 3 }) }" title="标题3">
          H3
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <div class="toolbar-group">
        <button type="button" @click="editor.chain().focus().toggleBold().run()" :class="{ active: editor.isActive('bold') }" title="粗体">
          B
        </button>
        <button type="button" @click="editor.chain().focus().toggleItalic().run()" :class="{ active: editor.isActive('italic') }" title="斜体">
          I
        </button>
        <button type="button" @click="editor.chain().focus().toggleUnderline().run()" :class="{ active: editor.isActive('underline') }" title="下划线">
          U
        </button>
        <button type="button" @click="editor.chain().focus().toggleStrike().run()" :class="{ active: editor.isActive('strike') }" title="删除线">
          S
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <div class="toolbar-group">
        <button type="button" @click="editor.chain().focus().toggleBulletList().run()" :class="{ active: editor.isActive('bulletList') }" title="无序列表">
          •
        </button>
        <button type="button" @click="editor.chain().focus().toggleOrderedList().run()" :class="{ active: editor.isActive('orderedList') }" title="有序列表">
          1.
        </button>
        <button type="button" @click="editor.chain().focus().toggleTaskList().run()" :class="{ active: editor.isActive('taskList') }" title="任务列表">
          ☑
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <div class="toolbar-group">
        <button type="button" @click="editor.chain().focus().toggleBlockquote().run()" :class="{ active: editor.isActive('blockquote') }" title="引用">
          "
        </button>
        <button type="button" @click="editor.chain().focus().toggleCodeBlock().run()" :class="{ active: editor.isActive('codeBlock') }" title="代码块">
          &lt;/&gt;
        </button>
        <button type="button" @click="editor.chain().focus().toggleCode().run()" :class="{ active: editor.isActive('code') }" title="行内代码">
          &lt;&gt;
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <div class="toolbar-group">
        <button type="button" @click="editor.chain().focus().setHorizontalRule().run()" title="分割线">
          ―
        </button>
        <button type="button" @click="editor.chain().focus().setHardBreak().run()" title="换行">
          ↵
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <div class="toolbar-group">
        <button type="button" @click="addLink" :class="{ active: editor.isActive('link') }" title="添加链接">
          🔗
        </button>
        <button type="button" @click="triggerImageUpload" title="添加图片">
          🖼️
        </button>
      </div>

      <div class="toolbar-divider"></div>

      <div class="toolbar-group">
        <button type="button" @click="editor.chain().focus().setTextAlign('left').run()" :class="{ active: editor.isActive({ textAlign: 'left' }) }" title="左对齐">
          ⫷
        </button>
        <button type="button" @click="editor.chain().focus().setTextAlign('center').run()" :class="{ active: editor.isActive({ textAlign: 'center' }) }" title="居中">
          ⫿
        </button>
        <button type="button" @click="editor.chain().focus().setTextAlign('right').run()" :class="{ active: editor.isActive({ textAlign: 'right' }) }" title="右对齐">
          ⫸
        </button>
      </div>
    <!-- 源码/可视化切换 -->
    <div class="toolbar-spacer"></div>
    <div class="toolbar-group">
      <button type="button" @click="toggleSourceMode" :class="{ active: sourceMode }" :title="sourceMode ? '可视化编辑' : '源码编辑'">
        {{ sourceMode ? '&lt;/&gt; 可视化' : '&lt;/&gt; 源码' }}
      </button>
    </div>
  </div>

    <!-- 编辑器内容区 / 源码区 -->
    <EditorContent v-show="!sourceMode" :editor="editor" class="editor-content" />
    <textarea
      v-if="sourceMode"
      v-model="sourceCode"
      class="source-editor"
      :style="{ minHeight: height + 'px', maxHeight: (height * 1.5) + 'px' }"
      spellcheck="false"
      placeholder="在此编辑 HTML 源码..."
    />

    <!-- 隐藏的图片上传输入 -->
    <input
      ref="imageInput"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleImageUpload"
    />

    <!-- 链接对话框 -->
    <el-dialog v-model="linkDialogVisible" title="添加链接" width="400px">
      <el-form :model="linkForm" label-width="80px">
        <el-form-item label="链接地址">
          <el-input v-model="linkForm.url" placeholder="https://example.com" />
        </el-form-item>
        <el-form-item label="打开方式">
          <el-select v-model="linkForm.target" style="width: 100%">
            <el-option label="当前窗口" value="_self" />
            <el-option label="新窗口" value="_blank" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="linkDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmLink">确定</el-button>
      </template>
    </el-dialog>

    <!-- 媒体选择器 -->
    <MediaPicker v-model="mediaPickerVisible" @select="handleMediaSelect" />
  </div>
</template>

<script setup>
import { ref, watch, onBeforeUnmount, computed } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import { ElMessage } from 'element-plus'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import ResizeImage from 'tiptap-extension-resize-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import TextAlign from '@tiptap/extension-text-align'
import Underline from '@tiptap/extension-underline'
import { TextStyle } from '@tiptap/extension-text-style'
import Color from '@tiptap/extension-color'
import Highlight from '@tiptap/extension-highlight'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Subscript from '@tiptap/extension-subscript'
import Superscript from '@tiptap/extension-superscript'
import { Table, TableRow, TableCell, TableHeader } from '@tiptap/extension-table'
import { post } from '@/utils/request'
import MediaPicker from './MediaPicker.vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '开始输入内容...'
  },
  height: {
    type: [String, Number],
    default: 500
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

// 源码模式
const sourceMode = ref(false)
const sourceCode = ref('')

// 清理内容中的 AI 引用标记（如 [1](@context-ref?id=0)）和多余空白
const sanitizeContent = (html) => {
  if (!html) return html
  return html
    .replace(/\[\d+\]\(@context-ref\?id=\d+\)/g, '')
    .replace(/\[1\]\(@context-ref\?id=0\)/g, '')
    .replace(/ style="list-style-type:\s*[^"]*"/g, '')
    .replace(/ style="padding-left:\s*[^"]*"/g, '')
    .replace(/ style="border:\s*[^"]*"/g, '')
}

const imageInput = ref(null)
const linkDialogVisible = ref(false)
const mediaPickerVisible = ref(false)
const linkForm = ref({
  url: '',
  target: '_blank'
})

// 创建编辑器
const editor = useEditor({
  content: sanitizeContent(props.modelValue),
  editable: !props.disabled,
  extensions: [
    StarterKit.configure({
      heading: {
        levels: [1, 2, 3, 4, 5, 6]
      },
      // 禁用 StarterKit 内置的 link 和 underline，使用自定义配置
      link: false,
      underline: false
    }),
    ResizeImage.configure({
      inline: true,
      allowBase64: false,
      HTMLAttributes: {
        class: 'editor-image'
      }
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        target: '_blank',
        rel: 'noopener noreferrer'
      }
    }),
    Placeholder.configure({
      placeholder: props.placeholder
    }),
    TextAlign.configure({
      types: ['heading', 'paragraph']
    }),
    Underline,
    TextStyle,
    Color,
    Highlight.configure({
      multicolor: true
    }),
    TaskList,
    TaskItem.configure({
      nested: true
    }),
    Subscript,
    Superscript,
    Table.configure({
      resizable: true
    }),
    TableRow,
    TableCell,
    TableHeader
  ],
  onUpdate: ({ editor }) => {
    const html = editor.getHTML()
    emit('update:modelValue', html)
    emit('change', html)
  }
})

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  if (!editor.value) return
  const cleaned = sanitizeContent(newVal)
  const current = editor.value.getHTML()
  if (cleaned && cleaned !== current) {
    editor.value.commands.setContent(cleaned, false)
  }
})

// 监听 disabled 变化
watch(() => props.disabled, (newVal) => {
  if (editor.value) {
    editor.value.setEditable(!newVal)
  }
})

// 切换源码模式
const toggleSourceMode = () => {
  if (!sourceMode.value) {
    // 进入源码模式：把编辑器 HTML 取出来
    sourceCode.value = editor.value ? editor.value.getHTML() : props.modelValue
    sourceMode.value = true
  } else {
    // 退出源码模式：把源码写回编辑器
    if (editor.value && sourceCode.value) {
      editor.value.commands.setContent(sanitizeContent(sourceCode.value), false)
      emit('update:modelValue', editor.value.getHTML())
      emit('change', editor.value.getHTML())
    }
    sourceMode.value = false
  }
}

// 监听源码变化，同步到外部
watch(sourceCode, (val) => {
  if (sourceMode.value) {
    emit('update:modelValue', val)
    emit('change', val)
  }
})

// 添加链接
const addLink = () => {
  if (editor.value) {
    const previousUrl = editor.value.getAttributes('link').href
    linkForm.value.url = previousUrl || ''
    linkForm.value.target = editor.value.getAttributes('link').target || '_blank'
    linkDialogVisible.value = true
  }
}

const confirmLink = () => {
  if (editor.value) {
    if (linkForm.value.url) {
      editor.value.chain().focus().extendMarkRange('link').setLink({
        href: linkForm.value.url,
        target: linkForm.value.target
      }).run()
    } else {
      editor.value.chain().focus().unsetLink().run()
    }
  }
  linkDialogVisible.value = false
}

// 触发图片上传
const triggerImageUpload = () => {
  mediaPickerVisible.value = true
}

// 处理媒体选择
const handleMediaSelect = (media) => {
  if (editor.value && media) {
    const url = media.url || media.path
    editor.value.chain().focus().setImage({ src: url, alt: media.filename }).run()
  }
}

// 处理图片上传（备用）
const handleImageUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    const response = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const url = response.data?.url || response.url || response.path
    if (editor.value && url) {
      editor.value.chain().focus().setImage({ src: url, alt: file.name }).run()
    }
  } catch (error) {
    ElMessage.error('图片上传失败')
  }

  // 清空 input
  event.target.value = ''
}

// 组件销毁时销毁编辑器
onBeforeUnmount(() => {
  if (editor.value) {
    editor.value.destroy()
  }
})

// 暴露方法
defineExpose({
  getContent: () => editor.value?.getHTML() || '',
  setContent: (content) => {
    if (editor.value) {
      editor.value.commands.setContent(sanitizeContent(content), false)
    }
  },
  focus: () => editor.value?.commands.focus()
})
</script>

<style scoped>
.tiptap-editor-wrapper {
  width: 100%;
  border: 1px solid var(--el-border-color);
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}

.tiptap-editor-wrapper:focus-within {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-9);
}

.editor-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 8px 12px;
  background: #f8f9fa;
  border-bottom: 1px solid var(--el-border-color);
}

.toolbar-group {
  display: flex;
  gap: 2px;
}

.toolbar-divider {
  width: 1px;
  background: var(--el-border-color);
  margin: 0 4px;
}

.editor-toolbar button {
  min-width: 32px;
  height: 32px;
  padding: 4px 8px;
  border: 1px solid transparent;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  transition: all 0.2s;
}

.editor-toolbar button:hover {
  background: #e5e7eb;
  border-color: #d1d5db;
}

.editor-toolbar button.active {
  background: var(--el-color-primary);
  color: white;
  border-color: var(--el-color-primary);
}

.editor-toolbar button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.editor-content {
  min-height: v-bind(height + 'px');
  max-height: calc(v-bind(height) * 1.5 + 'px');
  overflow-y: auto;
}

.editor-content :deep(.ProseMirror) {
  min-height: v-bind(height + 'px');
  padding: 16px;
  outline: none;
}

.editor-content :deep(.ProseMirror p) {
  margin: 0 0 0.5em 0;
}

.editor-content :deep(.ProseMirror h1),
.editor-content :deep(.ProseMirror h2),
.editor-content :deep(.ProseMirror h3),
.editor-content :deep(.ProseMirror h4),
.editor-content :deep(.ProseMirror h5),
.editor-content :deep(.ProseMirror h6) {
  margin: 1em 0 0.5em 0;
  font-weight: 600;
}

.editor-content :deep(.ProseMirror h1) { font-size: 2em; }
.editor-content :deep(.ProseMirror h2) { font-size: 1.5em; }
.editor-content :deep(.ProseMirror h3) { font-size: 1.25em; }

.editor-content :deep(.ProseMirror ul),
.editor-content :deep(.ProseMirror ol) {
  padding-left: 1.5em;
  margin: 0.5em 0;
}

.editor-content :deep(.ProseMirror li) {
  margin: 0.25em 0;
}

.editor-content :deep(.ProseMirror blockquote) {
  border-left: 3px solid #e5e7eb;
  padding-left: 1em;
  margin: 0.5em 0;
  color: #6b7280;
}

.editor-content :deep(.ProseMirror pre) {
  background: #1f2937;
  color: #f9fafb;
  padding: 12px 16px;
  border-radius: 6px;
  overflow-x: auto;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  margin: 0.5em 0;
}

.editor-content :deep(.ProseMirror code) {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9em;
}

.editor-content :deep(.ProseMirror pre code) {
  background: transparent;
  padding: 0;
}

.editor-content :deep(.ProseMirror img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 0.5em 0;
}

/* 图片选中状态样式 */
.editor-content :deep(.ProseMirror img.ProseMirror-selectednode) {
  outline: 3px solid var(--el-color-primary);
  outline-offset: 2px;
}

/* 图片容器样式 */
.editor-content :deep(.ProseMirror .image-view) {
  display: inline-block;
  position: relative;
  margin: 0.5em 0;
}

.editor-content :deep(.ProseMirror .image-view img) {
  display: block;
  cursor: pointer;
}

/* 缩放手柄样式 */
.editor-content :deep(.ProseMirror .image-view .resize-handle) {
  position: absolute;
  width: 12px;
  height: 12px;
  background: white;
  border: 2px solid var(--el-color-primary);
  border-radius: 2px;
  z-index: 10;
}

.editor-content :deep(.ProseMirror .image-view .resize-handle-top-left) {
  top: -6px;
  left: -6px;
  cursor: nwse-resize;
}

.editor-content :deep(.ProseMirror .image-view .resize-handle-top-right) {
  top: -6px;
  right: -6px;
  cursor: nesw-resize;
}

.editor-content :deep(.ProseMirror .image-view .resize-handle-bottom-left) {
  bottom: -6px;
  left: -6px;
  cursor: nesw-resize;
}

.editor-content :deep(.ProseMirror .image-view .resize-handle-bottom-right) {
  bottom: -6px;
  right: -6px;
  cursor: nwse-resize;
}

.editor-content :deep(.ProseMirror a) {
  color: var(--el-color-primary);
  text-decoration: underline;
}

.editor-content :deep(.ProseMirror hr) {
  border: none;
  border-top: 2px solid #e5e7eb;
  margin: 1em 0;
}

.editor-content :deep(.ProseMirror .is-empty::before) {
  content: attr(data-placeholder);
  float: left;
  color: #9ca3af;
  pointer-events: none;
  height: 0;
}

/* 任务列表样式 */
.editor-content :deep(.ProseMirror ul[data-type="taskList"]) {
  list-style: none;
  padding-left: 0;
}

.editor-content :deep(.ProseMirror ul[data-type="taskList"] li) {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.editor-content :deep(.ProseMirror ul[data-type="taskList"] li label) {
  flex-shrink: 0;
  margin-top: 4px;
}

.editor-content :deep(.ProseMirror ul[data-type="taskList"] li div) {
  flex: 1;
}

/* 表格样式 */
.toolbar-spacer {
  flex: 1;
}

.source-editor {
  width: 100%;
  padding: 16px;
  border: none;
  outline: none;
  resize: vertical;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: #e5e7eb;
  background: #1f2937;
  border-radius: 0 0 8px 8px;
  box-sizing: border-box;
}

.source-editor::placeholder {
  color: #6b7280;
}

.editor-content :deep(.ProseMirror table) {
  border-collapse: collapse;
  width: 100%;
  margin: 0.5em 0;
}

.editor-content :deep(.ProseMirror td),
.editor-content :deep(.ProseMirror th) {
  border: 1px solid #e5e7eb;
  padding: 8px 12px;
  text-align: left;
}

.editor-content :deep(.ProseMirror th) {
  background: #f9fafb;
  font-weight: 600;
}

.editor-content :deep(.ProseMirror .selectedCell) {
  background: var(--el-color-primary-light-9);
}
</style>
