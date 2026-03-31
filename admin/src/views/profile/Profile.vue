<template>
  <div class="profile-root">

    <!-- 顶栏 -->
    <div class="page-header">
      <div>
        <h1 class="page-title">个人中心</h1>
        <p class="page-subtitle">管理你的账号信息和安全设置</p>
      </div>
    </div>

    <!-- 双栏主体 -->
    <div class="profile-layout">

      <!-- ▌左侧：身份卡 ▌ -->
      <aside class="profile-sidebar">

        <!-- 用户卡 -->
        <div class="user-card">

          <!-- 顶部色条（无渐变，实色） -->
          <div class="user-card-banner"></div>

          <!-- 头像 -->
          <div class="user-avatar-wrap">
            <div class="user-avatar-btn" @click="triggerAvatarUpload" :title="'点击换头像'">
              <el-avatar :size="80" :src="form.avatar" class="user-avatar">
                <span class="avatar-fallback">{{ (form.name || '?').charAt(0).toUpperCase() }}</span>
              </el-avatar>
              <div class="avatar-hover-mask">
                <el-icon :size="18"><Camera /></el-icon>
              </div>
              <div class="avatar-uploading" v-if="uploadingAvatar">
                <el-icon class="spin"><Loading /></el-icon>
              </div>
            </div>
            <input ref="avatarInput" type="file" accept="image/*" style="display:none" @change="handleAvatarChange" />
          </div>

          <!-- 基本信息 -->
          <div class="user-info">
            <div class="user-name">{{ form.name || '未设置昵称' }}</div>
            <div class="user-role-wrap">
              <span class="role-badge" :class="'role-' + (userInfo?.role || 'editor')">{{ roleLabel }}</span>
            </div>
            <div class="user-email">
              <el-icon :size="12"><Message /></el-icon>
              <span>{{ userInfo?.email || '-' }}</span>
            </div>
            <p class="user-bio" v-if="form.bio">{{ form.bio }}</p>
            <p class="user-bio user-bio--empty" v-else>暂无个人简介</p>
          </div>

          <div class="card-sep"></div>

          <!-- 统计 -->
          <div class="stats-block">
            <div class="stat-item">
              <div class="stat-num">{{ stats.posts }}</div>
              <div class="stat-label">已发布</div>
            </div>
            <div class="stat-item">
              <div class="stat-num">{{ stats.drafts }}</div>
              <div class="stat-label">草稿</div>
            </div>
            <div class="stat-item">
              <div class="stat-num">{{ stats.days }}</div>
              <div class="stat-label">在站天</div>
            </div>
          </div>

          <div class="card-sep"></div>

          <!-- 快捷操作 -->
          <nav class="quick-nav">
            <router-link to="/change-password" class="quick-nav-item quick-nav-item--accent">
              <el-icon><Key /></el-icon>
              <span>修改密码</span>
              <el-icon class="nav-arrow"><ArrowRight /></el-icon>
            </router-link>
            <router-link to="/posts/new" class="quick-nav-item">
              <el-icon><EditPen /></el-icon>
              <span>写新文章</span>
              <el-icon class="nav-arrow"><ArrowRight /></el-icon>
            </router-link>
            <a href="/" target="_blank" class="quick-nav-item">
              <el-icon><View /></el-icon>
              <span>查看前台</span>
              <el-icon class="nav-arrow"><ArrowRight /></el-icon>
            </a>
          </nav>
        </div>

        <!-- 账号元信息 -->
        <div class="meta-card">
          <div class="meta-card-title">
            <el-icon><InfoFilled /></el-icon>
            账号详情
          </div>
          <div class="meta-rows">
            <div class="meta-row">
              <span class="meta-key">账号 ID</span>
              <span class="meta-val meta-mono">#{{ userInfo?.id || '-' }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-key">注册时间</span>
              <span class="meta-val">{{ formatDate(userInfo?.created_at) }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-key">最近更新</span>
              <span class="meta-val">{{ formatDate(userInfo?.updated_at) }}</span>
            </div>
          </div>
        </div>

      </aside>

      <!-- ▌右侧：编辑区 ▌ -->
      <main class="profile-main">

        <!-- 表单卡 -->
        <div class="form-card">
          <div class="form-card-header">
            <div class="form-card-icon">
              <el-icon :size="18"><UserFilled /></el-icon>
            </div>
            <div class="form-card-title-wrap">
              <span class="form-card-title">编辑资料</span>
              <span class="form-card-desc">更新的信息会同步到文章署名和后台导航</span>
            </div>
          </div>

          <div class="form-body">
            <el-form :model="form" :rules="rules" ref="formRef" label-position="top">

              <!-- 头像行 -->
              <div class="avatar-field">
                <div class="avatar-field-label">头像</div>
                <div class="avatar-field-body">
                  <el-avatar :size="52" :src="form.avatar" class="preview-avatar">
                    {{ (form.name || '?').charAt(0).toUpperCase() }}
                  </el-avatar>
                  <div class="avatar-field-actions">
                    <el-button size="small" @click="triggerAvatarUpload" :loading="uploadingAvatar">
                      <el-icon><Camera /></el-icon>
                      上传图片
                    </el-button>
                    <el-button size="small" plain v-if="form.avatar" @click="form.avatar = ''">
                      移除
                    </el-button>
                  </div>
                  <span class="avatar-field-hint">建议 200×200px 以上，JPG / PNG</span>
                </div>
              </div>

              <!-- 昵称 + 邮箱 两列 -->
              <div class="field-row-2col">
                <el-form-item label="显示昵称" prop="name">
                  <el-input v-model="form.name" placeholder="你的显示名称" :prefix-icon="UserFilled" maxlength="50" show-word-limit size="large" />
                </el-form-item>
                <el-form-item label="邮箱地址" prop="email">
                  <el-input v-model="form.email" placeholder="your@email.com" :prefix-icon="Message" size="large">
                    <template #suffix>
                      <el-tag size="small" type="info" effect="plain">登录账号</el-tag>
                    </template>
                  </el-input>
                </el-form-item>
              </div>

              <!-- 简介 -->
              <el-form-item label="个人简介">
                <el-input v-model="form.bio" type="textarea" :rows="4" placeholder="简单介绍一下自己，会显示在博客作者简介中..." maxlength="200" show-word-limit resize="none" size="large" />
              </el-form-item>

              <!-- 操作行 -->
              <div class="form-actions">
                <div class="form-actions-right">
                  <el-button @click="resetForm" :disabled="saving">重置</el-button>
                  <el-button type="primary" :loading="saving" @click="handleSave" size="large" style="min-width:120px">
                    {{ saving ? '保存中...' : '保存更改' }}
                  </el-button>
                </div>
              </div>
            </el-form>
          </div>
        </div>

        <!-- 安全设置 -->
        <div class="security-card">
          <div class="security-card-header">
            <el-icon><Lock /></el-icon>
            安全设置
          </div>
          <div class="security-item">
            <div class="security-item-info">
              <div class="security-item-title">登录密码</div>
              <div class="security-item-desc">建议定期更换密码，提升账号安全性</div>
            </div>
            <router-link to="/change-password">
              <el-button size="small">
                <el-icon><Key /></el-icon>
                修改密码
              </el-button>
            </router-link>
          </div>
        </div>

      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Camera, Message, UserFilled, InfoFilled,
  Key, ArrowRight, View, Loading,
  EditPen, Lock
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store'
import { get, put, upload } from '@/utils/request'

const userStore = useUserStore()
const formRef = ref(null)
const saving = ref(false)
const uploadingAvatar = ref(false)
const avatarInput = ref(null)
const userInfo = ref(null)

const stats = reactive({ posts: 0, drafts: 0, days: 0 })

const form = reactive({ name: '', email: '', avatar: '', bio: '' })

const rules = {
  name: [{ required: true, message: '昵称不能为空', trigger: 'blur' }],
  email: [
    { required: true, message: '邮箱不能为空', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱格式', trigger: 'blur' }
  ]
}

const roleLabel = computed(() => {
  const r = userInfo.value?.role
  return r === 'admin' ? '管理员' : r === 'editor' ? '编辑' : r || '用户'
})

const formatDate = (d) => {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' })
}

const loadUserInfo = async () => {
  try {
    const res = await get('/api/auth/me')
    userInfo.value = res.data || res
    Object.assign(form, {
      name: userInfo.value.name || '',
      email: userInfo.value.email || '',
      avatar: userInfo.value.avatar || '',
      bio: userInfo.value.bio || ''
    })
    if (userInfo.value.created_at) {
      stats.days = Math.max(1, Math.floor((Date.now() - new Date(userInfo.value.created_at)) / 86400000))
    }
  } catch {
    ElMessage.error('加载用户信息失败')
  }
}

const loadStats = async () => {
  try {
    const [pub, draft] = await Promise.allSettled([
      get('/api/posts', { page: 1, page_size: 1, status: 'published' }),
      get('/api/posts', { page: 1, page_size: 1, status: 'draft' })
    ])
    if (pub.status === 'fulfilled') stats.posts = pub.value.data?.total || 0
    if (draft.status === 'fulfilled') stats.drafts = draft.value.data?.total || 0
  } catch { /* ignore */ }
}

const triggerAvatarUpload = () => avatarInput.value?.click()

const handleAvatarChange = async (e) => {
  const file = e.target.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) { ElMessage.error('请选择图片文件'); return }
  uploadingAvatar.value = true
  try {
    const res = await upload('/api/media/upload', file)
    form.avatar = res.data?.url || res.url || ''
    if (form.avatar) ElMessage.success('头像上传成功')
    else ElMessage.warning('上传成功但未获取到地址')
  } catch {
    ElMessage.error('头像上传失败，请检查后端服务')
  } finally {
    uploadingAvatar.value = false
    e.target.value = ''
  }
}

const handleSave = async () => {
  try { await formRef.value.validate() } catch { return }
  saving.value = true
  try {
    await put('/api/auth/profile', { name: form.name, email: form.email, avatar: form.avatar, bio: form.bio })
    await userStore.getCurrentUser()
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const resetForm = () => loadUserInfo()

onMounted(() => { loadUserInfo(); loadStats() })
</script>

<style scoped>
/* ─── 根容器 ─── */
.profile-root {
  max-width: 1080px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ─── 顶栏 ─── */
.page-header { padding: 0 0 4px; }
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

/* ─── 双栏 ─── */
.profile-layout {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 20px;
  align-items: start;
}

/* ══ 左侧 ══ */
.profile-sidebar {
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: sticky;
  top: 16px;
}

/* 用户卡 */
.user-card {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  overflow: hidden;
}

/* 顶部色条 — 实色，无渐变 */
.user-card-banner {
  height: 56px;
  background-color: var(--el-color-primary);
  opacity: 0.12;
}

/* 头像 */
.user-avatar-wrap {
  display: flex;
  justify-content: center;
  margin-top: -40px;
  margin-bottom: 4px;
}
.user-avatar-btn {
  position: relative;
  width: 84px;
  height: 84px;
  cursor: pointer;
  flex-shrink: 0;
}
.user-avatar {
  width: 80px !important;
  height: 80px !important;
  border: 3px solid var(--el-bg-color);
  box-shadow: 0 2px 12px rgba(0,0,0,0.12);
  font-size: 28px;
}
.avatar-fallback { font-size: 28px; font-weight: 700; }
.avatar-hover-mask {
  position: absolute;
  inset: 3px;
  border-radius: 50%;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  opacity: 0;
  transition: opacity 0.18s;
}
.user-avatar-btn:hover .avatar-hover-mask { opacity: 1; }
.avatar-uploading {
  position: absolute;
  inset: 3px;
  border-radius: 50%;
  background: rgba(0,0,0,0.55);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}
.spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* 用户信息 */
.user-info {
  padding: 0 16px 8px;
  text-align: center;
}
.user-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.user-role-wrap { margin-bottom: 8px; }
.role-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
}
.role-admin {
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
  border: 1px solid var(--el-color-primary-light-7);
}
.role-editor {
  background: var(--el-color-success-light-9);
  color: var(--el-color-success);
  border: 1px solid var(--el-color-success-light-7);
}
.user-email {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 8px;
  max-width: 100%;
  overflow: hidden;
}
.user-email span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.user-bio {
  font-size: 12px;
  color: var(--el-text-color-regular);
  line-height: 1.6;
  margin: 4px 0 0;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}
.user-bio--empty { color: var(--el-text-color-placeholder); font-style: italic; }

/* 分割线 */
.card-sep {
  height: 1px;
  background: var(--el-border-color-lighter);
  margin: 0 0 4px;
}

/* 统计块 */
.stats-block {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  padding: 8px 12px 12px;
  gap: 4px;
}
.stat-item {
  text-align: center;
  padding: 8px 4px;
  border-radius: 8px;
  transition: background 0.15s;
}
.stat-item:hover { background: var(--el-fill-color-light); }
.stat-num {
  font-size: 22px;
  font-weight: 800;
  color: var(--el-color-primary);
  line-height: 1;
}
.stat-label {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}

/* 快捷导航 */
.quick-nav {
  padding: 0 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.quick-nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 8px;
  font-size: 13px;
  color: var(--el-text-color-regular);
  text-decoration: none;
  transition: background 0.15s, color 0.15s;
}
.quick-nav-item:hover { background: var(--el-fill-color-light); color: var(--el-text-color-primary); }
.quick-nav-item--accent { color: var(--el-color-primary); }
.quick-nav-item--accent:hover { background: var(--el-color-primary-light-9); }
.quick-nav-item .el-icon:first-child { font-size: 14px; flex-shrink: 0; }
.quick-nav-item > span { flex: 1; }
.nav-arrow { font-size: 11px; opacity: 0.4; }

/* 账号元信息卡 */
.meta-card {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  overflow: hidden;
}
.meta-card-title {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 16px;
  font-size: 12px;
  font-weight: 600;
  color: var(--el-text-color-secondary);
  border-bottom: 1px solid var(--el-border-color-lighter);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.meta-card-title .el-icon { color: var(--el-color-primary); }
.meta-rows { padding: 4px 0; }
.meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  font-size: 12px;
}
.meta-row:not(:last-child) { border-bottom: 1px solid var(--el-fill-color-light); }
.meta-key { color: var(--el-text-color-secondary); }
.meta-val { color: var(--el-text-color-primary); }
.meta-mono {
  font-family: 'Courier New', monospace;
  background: var(--el-fill-color-light);
  padding: 1px 6px;
  border-radius: 4px;
}

/* ══ 右侧 ══ */
.profile-main {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 表单卡 */
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

/* 头像字段 */
.avatar-field { margin-bottom: 24px; }
.avatar-field-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--el-text-color-regular);
  margin-bottom: 10px;
}
.avatar-field-body {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}
.preview-avatar {
  flex-shrink: 0;
  border: 2px solid var(--el-border-color-light);
}
.avatar-field-actions { display: flex; gap: 8px; }
.avatar-field-hint {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  flex-basis: 100%;
  margin-top: -4px;
}

/* 两列表单 */
.field-row-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 20px;
}

/* 操作栏 */
.form-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 20px;
  margin-top: 4px;
  border-top: 1px solid var(--el-border-color-lighter);
}
.form-actions-right { display: flex; gap: 10px; }

/* 安全设置卡 */
.security-card {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  overflow: hidden;
}
.security-card-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 14px 20px;
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color-lighter);
  background: var(--el-fill-color-blank);
}
.security-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  gap: 16px;
}
.security-item-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  margin-bottom: 3px;
}
.security-item-desc {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

/* ─── 响应式 ─── */
@media (max-width: 860px) {
  .profile-layout { grid-template-columns: 1fr; }
  .profile-sidebar { position: static; }
  .field-row-2col { grid-template-columns: 1fr; }
}
@media (max-width: 560px) {
  .form-body { padding: 16px; }
  .form-card-header { padding: 16px; }
}
</style>
