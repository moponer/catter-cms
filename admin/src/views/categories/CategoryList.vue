<template>
  <div class="category-page">
    <div class="page-header">
      <h1>分类管理</h1>
      <div style="display:flex;gap:8px">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建分类
        </el-button>
        <el-button
          v-if="selectedIds.length > 0"
          type="danger"
          @click="handleBatchDelete"
        >
          <el-icon><Delete /></el-icon>
          批量删除 ({{ selectedIds.length }})
        </el-button>
      </div>
    </div>

    <!-- 分类树形表格 -->
    <el-card class="table-card">
      <el-table
        :data="categories"
        row-key="id"
        border
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="45" />
        <el-table-column prop="name" label="分类名称" min-width="200" />
        <el-table-column prop="slug" label="Slug" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCatAction(cmd, row)">
              <el-button type="primary" link size="small">更多<el-icon class="dropdown-arrow"><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="child">添加子分类</el-dropdown-item>
                  <el-dropdown-item command="delete" style="color: var(--el-color-danger)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="Slug" prop="slug">
          <el-input v-model="form.slug" placeholder="URL友好名称（可选）">
            <template #append>
              <el-button @click="generateSlug">自动生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="父分类">
          <el-cascader
            v-model="form.parent_path"
            :options="categoryOptions"
            :props="{ checkStrictly: true, value: 'id', label: 'name', emitPath: false }"
            placeholder="选择父分类（可选）"
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="分类描述（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, ArrowDown } from '@element-plus/icons-vue'
import { get, post, put, del, delBody } from '@/utils/request'

const categories = ref([])
const dialogVisible = ref(false)
const selectedIds = ref([])

const handleSelectionChange = (selection) => {
  selectedIds.value = selection.map(item => item.id)
}
const editingId = ref(null)
const formRef = ref(null)
const submitting = ref(false)

const form = reactive({
  name: '',
  slug: '',
  description: '',
  parent_id: null,
  parent_path: null
})

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const dialogTitle = computed(() => editingId.value ? '编辑分类' : '新建分类')

// 分类选项（用于选择父分类）
const categoryOptions = computed(() => {
  const buildOptions = (cats, excludeId = null) => {
    return cats
      .filter(cat => cat.id !== excludeId)
      .map(cat => ({
        id: cat.id,
        name: cat.name,
        children: cat.children ? buildOptions(cat.children, excludeId) : undefined
      }))
  }
  return buildOptions(categories.value, editingId.value)
})

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await get('/api/categories/tree')
    categories.value = res.data || res || []
  } catch (error) {
    // 如果树形API不存在，使用普通列表
    try {
      const res = await get('/api/categories')
      categories.value = res.data || res || []
    } catch (e) {
      // 错误已在 request 拦截器中处理
    }
  }
}

// 新建分类
const handleCreate = () => {
  editingId.value = null
  Object.assign(form, {
    name: '',
    slug: '',
    description: '',
    parent_id: null,
    parent_path: null
  })
  dialogVisible.value = true
}

// 添加子分类
// 操作下拉菜单
const handleCatAction = (command, row) => {
  switch (command) {
    case 'child':
      handleCreateChild(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleCreateChild = (parent) => {
  editingId.value = null
  Object.assign(form, {
    name: '',
    slug: '',
    description: '',
    parent_id: parent.id,
    parent_path: parent.id
  })
  dialogVisible.value = true
}

// 编辑分类
const handleEdit = (row) => {
  editingId.value = row.id
  Object.assign(form, {
    name: row.name,
    slug: row.slug,
    description: row.description || '',
    parent_id: row.parent_id,
    parent_path: row.parent_id
  })
  dialogVisible.value = true
}

// 删除分类
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该分类吗？如果该分类下有子分类或文章，将无法删除。', '提示', {
      type: 'warning'
    })
    await del(`/api/categories/${row.id}`)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除分类
const handleBatchDelete = async () => {
  if (selectedIds.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.value.length} 个分类吗？有子分类或文章的分类将无法删除。`, '批量删除', {
      type: 'warning'
    })
    await delBody('/api/categories/batch', { ids: selectedIds.value })
    ElMessage.success('批量删除成功')
    selectedIds.value = []
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 自动生成Slug
const generateSlug = () => {
  // 简单的拼音转换（实际项目中可以使用更完善的拼音库）
  form.slug = form.name
    .toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^\w\u4e00-\u9fa5-]+/g, '')
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
    const data = {
      name: form.name,
      slug: form.slug || form.name,
      description: form.description,
      parent_id: form.parent_path || null
    }

    if (editingId.value) {
      await put(`/api/categories/${editingId.value}`, data)
      ElMessage.success('更新成功')
    } else {
      await post('/api/categories', data)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadCategories()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    submitting.value = false
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
  loadCategories()
})
</script>

<style scoped>
.category-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
