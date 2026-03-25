<template>
  <div class="page-container">
    <div class="page-header">
      <h1>标签管理</h1>
      <el-button type="primary" @click="openDialog()">新建标签</el-button>
    </div>
    <el-card>
      <el-table :data="tags" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="slug" label="别名" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新建/编辑标签对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑标签' : '新建标签'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="别名" prop="slug">
          <el-input v-model="form.slug" placeholder="请输入别名（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { get, post, put, del } from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const tags = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref()
const currentId = ref(null)

const form = reactive({
  name: '',
  slug: ''
})

const rules = {
  name: [{ required: true, message: '请输入标签名称', trigger: 'blur' }]
}

// 加载标签列表
const loadTags = async () => {
  loading.value = true
  try {
    const res = await get('/api/tags')
    tags.value = res.data?.list || res.list || []
  } catch (error) {
    ElMessage.error('加载标签失败')
  } finally {
    loading.value = false
  }
}

// 打开对话框
const openDialog = (tag = null) => {
  isEdit.value = !!tag
  if (tag) {
    currentId.value = tag.id
    form.name = tag.name
    form.slug = tag.slug
  } else {
    currentId.value = null
    form.name = ''
    form.slug = ''
  }
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await put(`/api/tags/${currentId.value}`, form)
      ElMessage.success('标签更新成功')
    } else {
      await post('/api/tags', form)
      ElMessage.success('标签创建成功')
    }
    dialogVisible.value = false
    loadTags()
  } catch (error) {
    // 错误已在 request 拦截器中处理
  } finally {
    submitting.value = false
  }
}

// 删除标签
const handleDelete = async (tag) => {
  try {
    await ElMessageBox.confirm(`确定要删除标签 "${tag.name}" 吗？`, '提示', {
      type: 'warning'
    })
    await del(`/api/tags/${tag.id}`)
    ElMessage.success('标签删除成功')
    loadTags()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: var(--spacing-md); }
.page-header { display: flex; justify-content: space-between; align-items: center; }
.page-header h1 { font-size: var(--font-size-2xl); font-weight: 600; margin: 0; }
</style>
