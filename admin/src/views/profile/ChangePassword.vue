<template>
  <div class="changepwd-root">

    <!-- 顶栏 -->
    <div class="page-header">
      <router-link to="/profile" class="back-btn">
        <el-icon><ArrowLeft /></el-icon>
        返回个人中心
      </router-link>
      <h1 class="page-title">修改密码</h1>
      <p class="page-subtitle">为保障安全，请定期更新你的登录密码</p>
    </div>

    <!-- 成功覆盖层 -->
    <transition name="success-overlay">
      <div v-if="success" class="success-overlay">
        <div class="success-box">
          <div class="success-icon-wrap">
            <el-icon :size="36"><CircleCheck /></el-icon>
          </div>
          <div class="success-title">密码修改成功</div>
          <div class="success-desc">即将跳转到个人中心...</div>
          <div class="success-bar">
            <div class="success-bar-fill" :style="{ width: countdownPct + '%' }"></div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 主体 -->
    <div class="changepwd-layout">

      <!-- 左侧：强度面板 -->
      <div class="strength-panel">

        <!-- 密码强度 -->
        <div class="strength-card">
          <div class="strength-card-title">密码强度</div>
          <!-- 四格强度条 -->
          <div class="strength-bars">
            <div
              v-for="n in 4"
              :key="n"
              class="strength-bar-seg"
              :class="getStrengthSegClass(n)"
            ></div>
          </div>
          <div class="strength-label" :class="'strength-' + strengthLevel">
            {{ strengthText }}
          </div>

          <!-- 规则列表 -->
          <div class="strength-rules">
            <div class="rule-item" :class="{ 'rule-ok': rules_check.length }">
              <span class="rule-dot"></span>
              至少 8 位字符
            </div>
            <div class="rule-item" :class="{ 'rule-ok': rules_check.hasUpper }">
              <span class="rule-dot"></span>
              包含大写字母
            </div>
            <div class="rule-item" :class="{ 'rule-ok': rules_check.hasLower }">
              <span class="rule-dot"></span>
              包含小写字母
            </div>
            <div class="rule-item" :class="{ 'rule-ok': rules_check.hasNumber }">
              <span class="rule-dot"></span>
              包含数字
            </div>
            <div class="rule-item" :class="{ 'rule-ok': rules_check.hasSymbol }">
              <span class="rule-dot"></span>
              包含特殊字符
            </div>
          </div>
        </div>

        <!-- 安全建议 -->
        <div class="tips-card">
          <div class="tips-card-title">安全建议</div>
          <ul class="tips-list">
            <li>不要使用生日、手机号等个人信息</li>
            <li>不同账号使用不同密码</li>
            <li>建议使用密码管理器</li>
          </ul>
        </div>

      </div>

      <!-- 右侧：表单 -->
      <div class="form-card">
        <div class="form-card-header">
          <div class="form-card-icon">
            <el-icon :size="18"><Lock /></el-icon>
          </div>
          <div>
            <span class="form-card-title">设置新密码</span>
            <span class="form-card-desc">请依次填写当前密码和新密码</span>
          </div>
        </div>

        <div class="form-body">
          <el-form :model="form" ref="formRef" label-position="top">

            <!-- 当前密码 -->
            <div class="field-section">
              <div class="field-step">
                <span class="step-num" :class="{ 'step-done': form.oldPassword }">
                  {{ form.oldPassword ? '✓' : '1' }}
                </span>
                <span class="step-label">当前密码</span>
              </div>
              <el-form-item prop="oldPassword">
                <el-input
                  v-model="form.oldPassword"
                  type="password"
                  show-password
                  placeholder="输入当前登录密码"
                  size="large"
                  :prefix-icon="Key"
                  :class="{ 'input-filled': form.oldPassword }"
                />
              </el-form-item>
            </div>

            <div class="field-divider">
              <span></span>
              <el-icon :size="12"><ArrowDown /></el-icon>
              <span></span>
            </div>

            <!-- 新密码 -->
            <div class="field-section">
              <div class="field-step">
                <span class="step-num" :class="{ 'step-active': form.oldPassword && !form.newPassword, 'step-done': form.newPassword }">
                  {{ form.newPassword ? '✓' : '2' }}
                </span>
                <span class="step-label">新密码</span>
              </div>
              <el-form-item prop="newPassword">
                <el-input
                  v-model="form.newPassword"
                  type="password"
                  show-password
                  placeholder="输入新密码（建议 8 位以上）"
                  size="large"
                  :prefix-icon="Lock"
                  @input="onNewPasswordInput"
                />
              </el-form-item>
              <!-- 内联强度条 -->
              <div class="inline-strength" v-if="form.newPassword">
                <div class="inline-bars">
                  <div v-for="n in 4" :key="n" class="inline-bar" :class="getStrengthSegClass(n)"></div>
                </div>
                <span class="inline-strength-text" :class="'strength-' + strengthLevel">{{ strengthText }}</span>
              </div>
            </div>

            <div class="field-divider">
              <span></span>
              <el-icon :size="12"><ArrowDown /></el-icon>
              <span></span>
            </div>

            <!-- 确认新密码 -->
            <div class="field-section">
              <div class="field-step">
                <span class="step-num" :class="{ 'step-active': form.newPassword && !form.confirmPassword, 'step-done': form.confirmPassword }">
                  {{ form.confirmPassword ? '✓' : '3' }}
                </span>
                <span class="step-label">确认新密码</span>
              </div>
              <el-form-item prop="confirmPassword">
                <el-input
                  v-model="form.confirmPassword"
                  type="password"
                  show-password
                  placeholder="再次输入新密码"
                  size="large"
                  :prefix-icon="Check"
                  :class="{
                    'input-match': form.confirmPassword && form.confirmPassword === form.newPassword,
                    'input-nomatch': form.confirmPassword && form.confirmPassword !== form.newPassword
                  }"
                />
              </el-form-item>
              <!-- 匹配状态 -->
              <div class="match-status" v-if="form.confirmPassword">
                <template v-if="form.confirmPassword === form.newPassword">
                  <el-icon color="var(--el-color-success)"><CircleCheck /></el-icon>
                  <span class="match-ok">两次密码一致</span>
                </template>
                <template v-else>
                  <el-icon color="var(--el-color-danger)"><CircleClose /></el-icon>
                  <span class="match-fail">两次密码不一致</span>
                </template>
              </div>
            </div>

            <!-- 操作栏 -->
            <div class="form-actions">
              <el-button @click="resetForm" :disabled="submitting">重置</el-button>
              <el-button
                type="primary"
                :loading="submitting"
                @click="handleSubmit"
                size="large"
                :disabled="!canSubmit"
                style="min-width: 140px"
              >
                {{ submitting ? '提交中...' : '确认修改' }}
              </el-button>
            </div>

          </el-form>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, ArrowDown, Key, Lock, Check, CircleCheck, CircleClose } from '@element-plus/icons-vue'
import { put } from '@/utils/request'

const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)
const success = ref(false)
const countdownPct = ref(100)

const form = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })

// ── 密码规则检查 ──
const rules_check = computed(() => {
  const p = form.newPassword
  return {
    length:    p.length >= 8,
    hasUpper:  /[A-Z]/.test(p),
    hasLower:  /[a-z]/.test(p),
    hasNumber: /[0-9]/.test(p),
    hasSymbol: /[^A-Za-z0-9]/.test(p)
  }
})

const strengthScore = computed(() => {
  const r = rules_check.value
  return [r.length, r.hasUpper, r.hasLower, r.hasNumber, r.hasSymbol].filter(Boolean).length
})

const strengthLevel = computed(() => {
  const s = strengthScore.value
  if (s <= 1) return 'weak'
  if (s === 2) return 'fair'
  if (s === 3) return 'good'
  return 'strong'
})

const strengthText = computed(() => {
  return { weak: '弱', fair: '一般', good: '强', strong: '非常强' }[strengthLevel.value]
})

const getStrengthSegClass = (n) => {
  const score = strengthScore.value
  if (score === 0) return ''
  if (n > score) return ''
  return `seg-${strengthLevel.value}`
}

const canSubmit = computed(() => {
  return form.oldPassword && form.newPassword && form.confirmPassword === form.newPassword && strengthScore.value >= 2
})

const onNewPasswordInput = () => { /* 实时响应 */ }

const handleSubmit = async () => {
  if (!canSubmit.value) {
    if (!form.oldPassword) { ElMessage.warning('请输入当前密码'); return }
    if (!form.newPassword) { ElMessage.warning('请输入新密码'); return }
    if (form.confirmPassword !== form.newPassword) { ElMessage.warning('两次密码不一致'); return }
    if (strengthScore.value < 2) { ElMessage.warning('密码强度太低，请设置更复杂的密码'); return }
    return
  }
  submitting.value = true
  try {
    await put('/api/auth/password', {
      old_password: form.oldPassword,
      new_password: form.newPassword
    })
    success.value = true
    // 倒计时进度
    const step = 100 / 20
    const timer = setInterval(() => {
      countdownPct.value = Math.max(0, countdownPct.value - step)
      if (countdownPct.value <= 0) {
        clearInterval(timer)
        router.push('/profile')
      }
    }, 100)
  } catch (e) {
    ElMessage.error('修改失败：' + (e?.message || '当前密码可能不正确'))
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  Object.assign(form, { oldPassword: '', newPassword: '', confirmPassword: '' })
}
</script>

<style scoped>
/* ─── 根容器 ─── */
.changepwd-root {
  max-width: 920px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
  position: relative;
}

/* ─── 顶栏 ─── */
.page-header { display: flex; flex-direction: column; gap: 4px; }
.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
  text-decoration: none;
  margin-bottom: 8px;
  transition: color 0.15s;
}
.back-btn:hover { color: var(--el-color-primary); }
.page-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin: 0 0 4px;
}
.page-subtitle {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

/* ─── 成功覆盖层 ─── */
.success-overlay {
  position: fixed;
  inset: 0;
  background: rgba(255,255,255,0.92);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.dark .success-overlay { background: rgba(0,0,0,0.85); }
.success-box {
  text-align: center;
  padding: 48px 40px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 16px;
  box-shadow: 0 8px 40px rgba(0,0,0,0.12);
  min-width: 300px;
}
.success-icon-wrap {
  color: var(--el-color-success);
  margin-bottom: 16px;
}
.success-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin-bottom: 8px;
}
.success-desc {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin-bottom: 24px;
}
.success-bar {
  height: 4px;
  background: var(--el-fill-color-darker);
  border-radius: 2px;
  overflow: hidden;
}
.success-bar-fill {
  height: 100%;
  background: var(--el-color-success);
  transition: width 0.1s linear;
}
.success-overlay-enter-active, .success-overlay-leave-active { transition: opacity 0.25s; }
.success-overlay-enter-from, .success-overlay-leave-to { opacity: 0; }

/* ─── 主布局 ─── */
.changepwd-layout {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 20px;
  align-items: start;
}

/* ══ 左侧强度面板 ══ */
.strength-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: sticky;
  top: 16px;
}

.strength-card, .tips-card {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  padding: 18px 16px;
}

.strength-card-title, .tips-card-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--el-text-color-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-bottom: 14px;
}

/* 四格强度条 */
.strength-bars {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 4px;
  margin-bottom: 8px;
}
.strength-bar-seg {
  height: 6px;
  border-radius: 3px;
  background: var(--el-fill-color-darker);
  transition: background 0.25s;
}
.seg-weak    { background: var(--el-color-danger); }
.seg-fair    { background: var(--el-color-warning); }
.seg-good    { background: var(--el-color-primary); }
.seg-strong  { background: var(--el-color-success); }

.strength-label {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--el-text-color-placeholder);
}
.strength-weak   { color: var(--el-color-danger); }
.strength-fair   { color: var(--el-color-warning); }
.strength-good   { color: var(--el-color-primary); }
.strength-strong { color: var(--el-color-success); }

/* 规则列表 */
.strength-rules {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.rule-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  transition: color 0.2s;
}
.rule-item.rule-ok { color: var(--el-color-success); }
.rule-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--el-border-color);
  flex-shrink: 0;
  transition: background 0.2s;
}
.rule-item.rule-ok .rule-dot { background: var(--el-color-success); }

/* 安全建议 */
.tips-list {
  margin: 0;
  padding-left: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.tips-list li {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}

/* ══ 右侧表单 ══ */
.form-card {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  overflow: hidden;
}
.form-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.form-card-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: var(--el-color-primary-light-9);
  border: 1px solid var(--el-color-primary-light-7);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--el-color-primary);
  flex-shrink: 0;
}
.form-card-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  display: block;
}
.form-card-desc {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 2px;
  display: block;
}

.form-body { padding: 24px; }

/* 字段步骤 */
.field-section { margin-bottom: 4px; }
.field-step {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}
.step-num {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1.5px solid var(--el-border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--el-text-color-secondary);
  flex-shrink: 0;
  transition: all 0.2s;
}
.step-num.step-active {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}
.step-num.step-done {
  border-color: var(--el-color-success);
  color: var(--el-color-success);
  background: var(--el-color-success-light-9);
}
.step-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-regular);
}

/* 步骤间分隔 */
.field-divider {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0 12px 11px;
  color: var(--el-border-color);
}
.field-divider span {
  flex: 1;
  height: 1px;
  background: var(--el-border-color-lighter);
}

/* 内联强度条 */
.inline-strength {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 6px;
  margin-bottom: 4px;
}
.inline-bars {
  display: flex;
  gap: 3px;
  flex-shrink: 0;
}
.inline-bar {
  width: 28px;
  height: 4px;
  border-radius: 2px;
  background: var(--el-fill-color-darker);
  transition: background 0.25s;
}
.inline-strength-text {
  font-size: 12px;
  font-weight: 500;
  color: var(--el-text-color-placeholder);
}

/* 匹配状态 */
.match-status {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-top: 6px;
  font-size: 12px;
}
.match-ok  { color: var(--el-color-success); }
.match-fail { color: var(--el-color-danger); }

/* 输入框状态 */
:deep(.input-filled .el-input__wrapper) { box-shadow: 0 0 0 1px var(--el-color-primary) inset; }
:deep(.input-match .el-input__wrapper)  { box-shadow: 0 0 0 1px var(--el-color-success) inset; }
:deep(.input-nomatch .el-input__wrapper){ box-shadow: 0 0 0 1px var(--el-color-danger)  inset; }

/* 操作栏 */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  margin-top: 8px;
  border-top: 1px solid var(--el-border-color-lighter);
}

/* ─── 响应式 ─── */
@media (max-width: 720px) {
  .changepwd-layout { grid-template-columns: 1fr; }
  .strength-panel { position: static; }
  .form-body { padding: 16px; }
}
</style>
