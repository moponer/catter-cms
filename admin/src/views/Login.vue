<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">C</div>
        <h1>Catter CMS</h1>
        <p>轻量级内容管理系统</p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" class="login-form">
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="邮箱"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { post } from '@/utils/request'
import { setToken, setCurrentUser } from '@/utils/auth'

const router = useRouter()
const formRef = ref()
const loading = ref(false)

const form = reactive({
  email: '',
  password: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    const res = await post('/api/auth/login', form)
    
    // 处理响应 - 标准格式是 { code: 0, message: "...", data: {...} }
    if (res.code === 0 || res.code === 200) {
      const data = res.data
      setToken(data.token)
      setCurrentUser(data.user)
      
      ElMessage.success(res.message || '登录成功')
      router.push('/')
    } else {
      ElMessage.error(res.message || '登录失败')
    }
  } catch (error) {
    console.error('登录失败:', error)
    if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error('登录失败,请检查网络连接')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
}

.login-card {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  background: #409eff;
  border-radius: 16px;
  color: white;
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 16px;
}

.login-header h1 {
  margin: 0 0 8px;
  font-size: 24px;
  color: var(--color-text-primary);
}

.login-header p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 14px;
}

.login-form {
  margin-top: 24px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
}
</style>