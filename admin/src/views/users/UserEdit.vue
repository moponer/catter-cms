<template>
  <div class="user-edit-page">
    <div class="page-header">
      <h1>编辑用户</h1>
    </div>
    <el-card>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="form.password" type="password" placeholder="留空则不修改密码" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="选择角色">
            <el-option label="管理员" value="admin" />
            <el-option label="编辑" value="editor" />
            <el-option label="作者" value="author" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { get, put } from '@/utils/request'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const saving = ref(false)

const form = reactive({
  name: '',
  email: '',
  password: '',
  role: 'editor'
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  email: [{ required: true, type: 'email', message: '请输入有效的邮箱', trigger: 'blur' }]
}

const loadUser = async () => {
  try {
    const res = await get(`/api/users/${route.params.id}`)
    const data = res.data || res
    form.name = data.name || ''
    form.email = data.email || ''
    form.role = data.role || 'editor'
    form.password = ''
  } catch (error) {
    console.error('加载用户失败:', error)
    ElMessage.error('加载用户失败')
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    const data = { name: form.name, email: form.email, role: form.role }
    if (form.password) data.password = form.password
    await put(`/api/users/${route.params.id}`, data)
    ElMessage.success('保存成功')
    router.push('/users')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadUser()
})
</script>

<style scoped>
.user-edit-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.page-header h1 {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}
</style>
