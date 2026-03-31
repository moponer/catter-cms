<template>
  <div class="login-page">
    <!-- 背景动画 -->
    <div class="bg-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
      <div class="shape shape-4"></div>
      <div class="shape shape-5"></div>
    </div>

    <!-- 左侧装饰区域 -->
    <div class="login-left">
      <div class="brand-content">
        <div class="brand-icon">
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M24 4L42 14V34L24 44L6 34V14L24 4Z" stroke="currentColor" stroke-width="2.5" stroke-linejoin="round"/>
            <path d="M24 14L33 19V29L24 34L15 29V19L24 14Z" stroke="currentColor" stroke-width="2.5" stroke-linejoin="round"/>
            <circle cx="24" cy="24" r="4" fill="currentColor"/>
          </svg>
        </div>
        <h1 class="brand-title">Catter CMS</h1>
        <p class="brand-subtitle">轻量级 · 高性能 · 易扩展</p>
        <div class="feature-list">
          <div class="feature-item">
            <i class="fa-solid fa-bolt"></i>
            <span>极速响应</span>
          </div>
          <div class="feature-item">
            <i class="fa-solid fa-shield-halved"></i>
            <span>安全可靠</span>
          </div>
          <div class="feature-item">
            <i class="fa-solid fa-palette"></i>
            <span>主题丰富</span>
          </div>
        </div>
      </div>
      <!-- 装饰性代码背景 -->
      <div class="code-bg">
        <pre><code>{{ codeArt }}</code></pre>
      </div>
    </div>

    <!-- 右侧登录区域 -->
    <div class="login-right">
      <div class="login-card">
        <div class="login-header">
          <h2>欢迎回来</h2>
          <p>请登录您的账户以继续</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" class="login-form">
          <el-form-item prop="email">
            <div class="input-wrapper">
              <label class="input-label">邮箱地址</label>
              <el-input
                v-model="form.email"
                placeholder="请输入邮箱"
                size="large"
                class="custom-input"
              >
                <template #prefix>
                  <i class="fa-solid fa-envelope input-icon"></i>
                </template>
              </el-input>
            </div>
          </el-form-item>

          <el-form-item prop="password">
            <div class="input-wrapper">
              <label class="input-label">密码</label>
              <el-input
                v-model="form.password"
                type="password"
                placeholder="请输入密码"
                size="large"
                show-password
                class="custom-input"
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <i class="fa-solid fa-lock input-icon"></i>
                </template>
              </el-input>
            </div>
          </el-form-item>

          <el-form-item class="remember-row">
            <el-checkbox v-model="rememberMe" class="remember-checkbox">
              记住我
            </el-checkbox>
            <a href="#" class="forgot-link">忘记密码?</a>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              class="login-btn"
              :loading="loading"
              @click="handleLogin"
            >
              <span v-if="!loading">登 录</span>
              <span v-else>登录中...</span>
              <i v-if="!loading" class="fa-solid fa-arrow-right"></i>
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          <p>Powered by <strong>Catter CMS</strong></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { post } from '@/utils/request'
import { setToken, setCurrentUser } from '@/utils/auth'

const router = useRouter()
const route = useRoute()
const formRef = ref()
const loading = ref(false)
const rememberMe = ref(false)

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

// 装饰性代码艺术
const codeArt = `func main() {
  router := gin.Default()
  router.Run(":3000")
}

// 快速、简洁、优雅`

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

      // 优先使用 redirect 参数，否则回到首页
      const redirect = route.query.redirect || '/'
      router.push(redirect)
    } else {
      ElMessage.error(res.message || '登录失败')
    }
  } catch (error) {
    if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error('登录失败,请检查网络连接')
    }
  } finally {
    loading.value = false
  }
}

// 检查是否已登录，已登录则直接跳转
onMounted(() => {
  const token = localStorage.getItem('token')
  if (token) {
    // 有 token，尝试跳转回之前的页面
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  }
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}

/* 背景动画形状 */
.bg-shapes {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.3;
  animation: float 20s ease-in-out infinite;
}

.shape-1 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.shape-2 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, #f093fb, #f5576c);
  top: 50%;
  right: -50px;
  animation-delay: -5s;
}

.shape-3 {
  width: 250px;
  height: 250px;
  background: linear-gradient(135deg, #4facfe, #00f2fe);
  bottom: -50px;
  left: 30%;
  animation-delay: -10s;
}

.shape-4 {
  width: 200px;
  height: 200px;
  background: linear-gradient(135deg, #43e97b, #38f9d7);
  top: 20%;
  left: 60%;
  animation-delay: -15s;
}

.shape-5 {
  width: 350px;
  height: 350px;
  background: linear-gradient(135deg, #fa709a, #fee140);
  bottom: 10%;
  right: 20%;
  animation-delay: -7s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(30px, -30px) rotate(90deg);
  }
  50% {
    transform: translate(-20px, 20px) rotate(180deg);
  }
  75% {
    transform: translate(20px, 30px) rotate(270deg);
  }
}

/* 左侧装饰区域 */
.login-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px;
  position: relative;
  z-index: 1;
}

.brand-content {
  text-align: center;
  position: relative;
  z-index: 2;
}

.brand-icon {
  width: 100px;
  height: 100px;
  margin: 0 auto 30px;
  color: #fff;
  animation: pulse 3s ease-in-out infinite;
}

.brand-icon svg {
  width: 100%;
  height: 100%;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.9;
  }
}

.brand-title {
  font-size: 48px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 12px;
  letter-spacing: -1px;
  text-shadow: 0 4px 30px rgba(102, 126, 234, 0.4);
}

.brand-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 40px;
  letter-spacing: 2px;
}

.feature-list {
  display: flex;
  gap: 30px;
  justify-content: center;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  transition: all 0.3s ease;
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

.feature-item i {
  font-size: 16px;
  color: #4facfe;
}

/* 代码背景装饰 */
.code-bg {
  position: absolute;
  bottom: 60px;
  left: 60px;
  right: 60px;
  opacity: 0.15;
  font-family: 'Fira Code', 'Monaco', monospace;
  font-size: 13px;
  color: #fff;
  pointer-events: none;
}

.code-bg pre {
  margin: 0;
  white-space: pre-wrap;
}

/* 右侧登录区域 */
.login-right {
  width: 620px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  padding: 40px 60px 40px 40px;
  position: relative;
  z-index: 1;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 50px 45px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
  box-shadow: 
    0 25px 50px rgba(0, 0, 0, 0.25),
    0 0 0 1px rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.login-header h2 {
  margin: 0 0 10px;
  font-size: 28px;
  font-weight: 600;
  color: #1a1a2e;
}

.login-header p {
  margin: 0;
  color: #6b7280;
  font-size: 15px;
}

.login-form {
  margin-top: 0;
}

.input-wrapper {
  width: 100%;
}

.input-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.custom-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 0 0 1px #e5e7eb;
  padding: 4px 15px;
  transition: all 0.3s ease;
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #667eea;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2), 0 0 0 1px #667eea;
}

.custom-input :deep(.el-input__inner) {
  height: 44px;
  line-height: 44px;
}

.input-icon {
  font-size: 16px;
  color: #9ca3af;
}

.remember-row {
  margin-bottom: 24px;
}

.remember-row :deep(.el-form-item__content) {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.remember-checkbox {
  font-size: 13px;
  color: #6b7280;
}

.forgot-link {
  font-size: 13px;
  color: #374151;
  text-decoration: none;
  transition: color 0.3s ease;
}

.forgot-link:hover {
  color: #1a1a2e;
}

.login-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  background: #1e3a5f;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(30, 58, 95, 0.3);
}

.login-btn:hover {
  transform: translateY(-1px);
  background: #152a45;
  box-shadow: 0 4px 16px rgba(30, 58, 95, 0.4);
}

.login-btn:active {
  transform: translateY(0);
}

.login-btn i {
  font-size: 14px;
  transition: transform 0.3s ease;
}

.login-btn:hover i {
  transform: translateX(4px);
}

.login-footer {
  margin-top: 30px;
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #f3f4f6;
}

.login-footer p {
  margin: 0;
  font-size: 12px;
  color: #9ca3af;
}

.login-footer strong {
  color: #1a1a2e;
}

/* 响应式 */
@media (max-width: 1024px) {
  .login-left {
    display: none;
  }

  .login-right {
    width: 100%;
    padding: 20px;
  }

  .login-card {
    padding: 40px 30px;
  }
}

@media (max-width: 480px) {
  .login-card {
    padding: 30px 25px;
    border-radius: 20px;
  }

  .brand-title {
    font-size: 36px;
  }

  .feature-list {
    flex-direction: column;
    gap: 15px;
  }
}
</style>
