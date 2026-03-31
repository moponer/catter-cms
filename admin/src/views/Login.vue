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
    <!-- 粒子效果 -->
    <div class="particles">
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
      <div class="particle"></div>
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
  background: linear-gradient(135deg, #0f1923 0%, #1a2f4a 50%, #0d2137 100%);
}

/* 背景动画形状 - 统一的蓝色调 */
.bg-shapes {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 25s ease-in-out infinite;
}

.shape-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(34, 211, 238, 0.6) 0%, rgba(59, 130, 246, 0.3) 70%);
  top: -150px;
  left: -100px;
  animation-delay: 0s;
  animation-duration: 28s;
}

.shape-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.5) 0%, rgba(59, 130, 246, 0.2) 70%);
  top: 60%;
  right: -100px;
  animation-delay: -8s;
  animation-duration: 32s;
}

.shape-3 {
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, rgba(34, 211, 238, 0.5) 0%, rgba(6, 182, 212, 0.2) 70%);
  bottom: -100px;
  left: 25%;
  animation-delay: -15s;
  animation-duration: 26s;
}

.shape-4 {
  width: 280px;
  height: 280px;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.4) 0%, rgba(59, 130, 246, 0.15) 70%);
  top: 15%;
  right: 25%;
  animation-delay: -12s;
  animation-duration: 30s;
}

.shape-5 {
  width: 320px;
  height: 320px;
  background: radial-gradient(circle, rgba(6, 182, 212, 0.45) 0%, rgba(34, 211, 238, 0.15) 70%);
  bottom: 20%;
  right: 15%;
  animation-delay: -5s;
  animation-duration: 24s;
}

/* 光晕粒子效果 */
.particles {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: rgba(34, 211, 238, 0.6);
  border-radius: 50%;
  animation: particle-float 15s linear infinite;
}

.particle:nth-child(1) { left: 10%; animation-delay: 0s; animation-duration: 18s; }
.particle:nth-child(2) { left: 20%; animation-delay: -3s; animation-duration: 15s; }
.particle:nth-child(3) { left: 35%; animation-delay: -6s; animation-duration: 20s; }
.particle:nth-child(4) { left: 50%; animation-delay: -9s; animation-duration: 16s; }
.particle:nth-child(5) { left: 65%; animation-delay: -12s; animation-duration: 22s; }
.particle:nth-child(6) { left: 80%; animation-delay: -2s; animation-duration: 17s; }
.particle:nth-child(7) { left: 90%; animation-delay: -7s; animation-duration: 19s; }

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  20% {
    transform: translate(40px, -50px) scale(1.05);
  }
  40% {
    transform: translate(-30px, 30px) scale(0.95);
  }
  60% {
    transform: translate(50px, 40px) scale(1.02);
  }
  80% {
    transform: translate(-20px, -30px) scale(1.08);
  }
}

@keyframes particle-float {
  0% {
    transform: translateY(100vh) scale(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
    transform: translateY(90vh) scale(1);
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-10vh) scale(0.5);
    opacity: 0;
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
  animation: icon-float 4s ease-in-out infinite;
  filter: drop-shadow(0 0 30px rgba(34, 211, 238, 0.5));
}

.brand-icon svg {
  width: 100%;
  height: 100%;
}

@keyframes icon-float {
  0%, 100% {
    transform: translateY(0) scale(1);
    filter: drop-shadow(0 0 30px rgba(34, 211, 238, 0.5));
  }
  50% {
    transform: translateY(-8px) scale(1.03);
    filter: drop-shadow(0 0 50px rgba(34, 211, 238, 0.7));
  }
}

.brand-title {
  font-size: 48px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 12px;
  letter-spacing: -1px;
  text-shadow: 0 0 40px rgba(34, 211, 238, 0.3);
  background: linear-gradient(135deg, #fff 0%, rgba(34, 211, 238, 0.9) 50%, #fff 100%);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: shimmer 4s linear infinite;
}

@keyframes shimmer {
  0% { background-position: 0% center; }
  100% { background-position: 200% center; }
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
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(34, 211, 238, 0.2);
  transition: all 0.3s ease;
}

.feature-item:hover {
  background: rgba(34, 211, 238, 0.15);
  border-color: rgba(34, 211, 238, 0.4);
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(34, 211, 238, 0.2);
}

.feature-item i {
  font-size: 16px;
  color: #22d3ee;
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
  box-shadow: 0 0 0 1px rgba(34, 211, 238, 0.2);
  padding: 4px 15px;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.95);
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px rgba(34, 211, 238, 0.5);
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(34, 211, 238, 0.15), 0 0 0 1px #22d3ee;
  background: #fff;
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
  border-radius: 12px;
  background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 15px rgba(14, 165, 233, 0.35);
  position: relative;
  overflow: hidden;
}

.login-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: left 0.5s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(14, 165, 233, 0.45);
  background: linear-gradient(135deg, #38bdf8 0%, #0ea5e9 100%);
}

.login-btn:hover::before {
  left: 100%;
}

.login-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(14, 165, 233, 0.35);
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
