<template>
  <div class="ai-settings-page">
    <div class="page-header">
      <h1>AI 集成</h1>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 快速设置 -->
      <el-tab-pane label="CMS 集成" name="setup">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">一键配置 AI 客户端</span>
              <el-button type="primary" plain size="small" @click="loadSetupCommands" :loading="loadingSetup">
                刷新
              </el-button>
            </div>
          </template>

          <div v-if="!setupData.has_key" class="empty-state">
            <el-icon class="empty-icon"><Connection /></el-icon>
            <p>还没有可用的 Catter CMS Key</p>
            <p class="empty-hint">请先在「密钥管理」中创建一个 Key</p>
          </div>

          <div v-else>
            <el-alert
              v-if="setupData.key_name"
              :title="'当前使用密钥：' + setupData.key_name + '（' + setupData.key_prefix + '）'"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 20px"
            />

            <!-- Step 1: 环境变量 -->
            <div class="setup-step">
              <div class="step-header">
                <span class="step-number">1</span>
                <span class="step-title">设置环境变量</span>
              </div>
              <p class="step-desc">复制以下命令到终端执行，将 CMS 地址和 API Key 持久化保存到系统环境变量（仅需执行一次，重启终端后仍然生效）：</p>

              <div class="code-block-wrapper">
                <div class="code-tabs">
                  <button
                    v-for="tab in codeTabs"
                    :key="tab.key"
                    :class="['code-tab', { active: activeCodeTab === tab.key }]"
                    @click="activeCodeTab = tab.key"
                  >{{ tab.label }}</button>
                </div>
                <div class="code-block">
                  <code>{{ setupData.commands[activeCodeTab] }}</code>
                  <button class="copy-btn" @click="copyCode(setupData.commands[activeCodeTab])">
                    <el-icon><DocumentCopy /></el-icon>
                  </button>
                </div>
              </div>
            </div>

            <!-- Step 2: Skill 提示词 -->
            <div class="setup-step">
              <div class="step-header">
                <span class="step-number">2</span>
                <span class="step-title">添加 AI Skill</span>
              </div>
              <p class="step-desc">复制以下提示词，粘贴到 AI 客户端的系统规则中或发送给AI，AI 会自动读取 Skill 文档并学会操作 CMS：</p>

              <div class="code-block-wrapper">
                <div class="code-block multiline">
                  <pre>{{ setupData.skill_prompt }}</pre>
                  <button class="copy-btn" @click="copyCode(setupData.skill_prompt)">
                    <el-icon><DocumentCopy /></el-icon>
                  </button>
                </div>
              </div>
            </div>

            <!-- 权限信息 -->
            <div class="setup-step">
              <div class="step-header">
                <span class="step-number">3</span>
                <span class="step-title">完成</span>
              </div>
              <p class="step-desc">配置完成！AI 客户端将自动获取 Skill 文档并了解所有 API 用法。</p>
              <div class="perm-tags">
                <el-tag v-for="perm in setupData.permissions" :key="perm" size="small" type="info">
                  {{ perm }}
                </el-tag>
                <el-tag v-for="scope in setupData.scopes" :key="scope" size="small">{{ scope }}</el-tag>
              </div>
              <div v-if="setupData.expires_at && !isNeverExpires(setupData.expires_at)" class="expire-info">
                <el-icon><Warning /></el-icon>
                <span>密钥有效期至：{{ formatDate(setupData.expires_at) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-tab-pane>

      <!-- 主题开发 Skill -->
      <el-tab-pane label="主题开发" name="theme-dev">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">主题开发 AI Skill</span>
              <el-button type="primary" plain size="small" @click="loadThemeSkillPrompt" :loading="loadingThemeSkill">
                刷新
              </el-button>
            </div>
          </template>

          <div class="theme-skill-intro">
            <p>将以下提示词发送给 AI 编程助手（如 WorkBuddy、Cursor、Claude 等），AI 会自动获取完整的主题开发规范并学会开发 Catter CMS 主题。</p>
            <p class="hint">开发或优化主题时，AI 会自动调用该 Skill，确保符合 Catter CMS 的模板语法、配置系统和设计规范。</p>
          </div>

          <div class="setup-step">
            <div class="step-header">
              <span class="step-number">1</span>
              <span class="step-title">复制主题开发提示词</span>
            </div>
            <p class="step-desc">点击下方按钮复制完整提示词，然后粘贴到 AI 助手的对话中：</p>

            <div class="code-block-wrapper">
              <div class="code-block multiline">
                <pre>{{ themeSkillPrompt }}</pre>
                <button class="copy-btn" @click="copyCode(themeSkillPrompt)">
                  <el-icon><DocumentCopy /></el-icon>
                </button>
              </div>
            </div>
          </div>

          <div class="setup-step">
            <div class="step-header">
              <span class="step-number">2</span>
              <span class="step-title">开始开发主题</span>
            </div>
            <p class="step-desc">安装完成后，告诉 AI 你的需求，例如：</p>
            <ul class="example-list">
              <li>"帮我开发一个新的博客主题，风格要简洁清新"</li>
              <li>"优化现有主题的移动端适配"</li>
              <li>"给主题增加一个轮播图模块"</li>
              <li>"修复主题配置表单的显示问题"</li>
            </ul>
          </div>

          <div class="skill-features">
            <h3>Skill 包含的开发规范</h3>
            <div class="feature-grid">
              <div class="feature-item">
                <el-icon><Document /></el-icon>
                <span>模板语法规范</span>
              </div>
              <div class="feature-item">
                <el-icon><Setting /></el-icon>
                <span>主题配置系统</span>
              </div>
              <div class="feature-item">
                <el-icon><Brush /></el-icon>
                <span>CSS 设计规范</span>
              </div>
              <div class="feature-item">
                <el-icon><DataAnalysis /></el-icon>
                <span>模板数据参考</span>
              </div>
              <div class="feature-item">
                <el-icon><Lightning /></el-icon>
                <span>快速入门教程</span>
              </div>
              <div class="feature-item">
                <el-icon><CircleCheck /></el-icon>
                <span>常见错误排查</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-tab-pane>

      <!-- 密钥管理 -->
      <el-tab-pane label="密钥管理" name="keys">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">Catter CMS Key 管理</span>
              <el-button type="primary" size="small" @click="showCreateDialog = true">
                创建密钥
              </el-button>
            </div>
          </template>

          <el-table :data="keys" stripe style="width: 100%" v-loading="loadingKeys">
            <el-table-column prop="name" label="名称" min-width="120" />
            <el-table-column label="密钥" min-width="150">
              <template #default="{ row }">
                <code class="key-prefix">{{ row.token_prefix }}</code>
              </template>
            </el-table-column>
            <el-table-column label="权限" min-width="120">
              <template #default="{ row }">
                <el-tag v-for="p in parseJSON(row.permissions)" :key="p" size="small" style="margin-right: 4px">
                  {{ p }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="有效期" min-width="120">
              <template #default="{ row }">
                <span v-if="isExpired(row.expires_at)" class="expired-text">已过期</span>
                <span v-else-if="isNeverExpires(row.expires_at)">永不过期</span>
                <span v-else>{{ formatDate(row.expires_at) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="最后使用" min-width="120">
              <template #default="{ row }">
                <span v-if="row.last_used_at">{{ formatDate(row.last_used_at) }}</span>
                <span v-else class="muted-text">未使用</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
                  {{ row.is_active ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="viewLogs(row)">
                  日志
                </el-button>
                <el-button type="danger" link size="small" @click="handleDelete(row)">
                  撤销
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- 审计日志 -->
      <el-tab-pane label="审计日志" name="logs" v-if="selectedToken">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">{{ selectedToken.name }} - 调用日志</span>
              <el-button size="small" @click="selectedToken = null; activeTab = 'keys'">
                返回
              </el-button>
            </div>
          </template>

          <el-table :data="logs" stripe style="width: 100%" v-loading="loadingLogs">
            <el-table-column prop="created_at" label="时间" min-width="160">
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="action" label="方法" width="80">
              <template #default="{ row }">
                <el-tag :type="getMethodType(row.action)" size="small">{{ row.action }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="resource" label="资源" width="100" />
            <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 200 ? 'success' : 'danger'" size="small">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ip_address" label="IP" width="130" />
          </el-table>

          <div class="pagination-wrapper" v-if="logTotal > logPageSize">
            <el-pagination
              v-model:current-page="logPage"
              :page-size="logPageSize"
              :total="logTotal"
              layout="prev, pager, next"
              @current-change="loadLogs"
            />
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 创建密钥对话框 -->
    <el-dialog v-model="showCreateDialog" title="创建 Catter CMS Key" width="500px" :close-on-click-modal="false">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="为这个 Key 起个名字" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="权限">
          <el-checkbox-group v-model="createForm.permissions">
            <el-checkbox label="read">读取</el-checkbox>
            <el-checkbox label="write">写入</el-checkbox>
            <el-checkbox label="admin">管理</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="资源范围">
          <el-checkbox-group v-model="createForm.scopes">
            <el-checkbox label="posts">文章</el-checkbox>
            <el-checkbox label="pages">页面</el-checkbox>
            <el-checkbox label="categories">分类</el-checkbox>
            <el-checkbox label="tags">标签</el-checkbox>
            <el-checkbox label="media">媒体</el-checkbox>
            <el-checkbox label="comments">评论</el-checkbox>
            <el-checkbox label="settings">设置</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="有效期">
          <el-radio-group v-model="createForm.expiryType">
            <el-radio label="never">永不过期</el-radio>
            <el-radio label="custom">自定义</el-radio>
          </el-radio-group>
          <el-input-number
            v-if="createForm.expiryType === 'custom'"
            v-model="createForm.expiresIn"
            :min="1"
            :max="365"
            style="margin-left: 12px"
          />
          <span v-if="createForm.expiryType === 'custom'" style="margin-left: 4px">天</span>
        </el-form-item>
      </el-form>

      <!-- 创建成功后显示完整 Token -->
      <div v-if="createdToken" class="created-token-box">
        <el-alert type="warning" :closable="false" show-icon>
          <template #title>
            <strong>请立即复制此 Key！关闭后将无法再次查看完整值。</strong>
          </template>
        </el-alert>
        <div class="token-display">
          <code>{{ createdToken.token }}</code>
          <el-button type="primary" size="small" @click="copyCode(createdToken.token)">复制</el-button>
        </div>
      </div>

      <template #footer>
        <el-button @click="showCreateDialog = false">{{ createdToken ? '关闭' : '取消' }}</el-button>
        <el-button v-if="!createdToken" type="primary" @click="handleCreate" :loading="creating">
          创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { get, post, del } from '@/utils/request'
import { Connection, Warning, Document, Setting, Brush, DataAnalysis, Lightning, CircleCheck } from '@element-plus/icons-vue'

const activeTab = ref('setup')
const activeCodeTab = ref('windows')
const codeTabs = [
  { key: 'windows', label: 'Windows (永久)' },
  { key: 'mac', label: 'macOS (zsh)' },
]

// Setup 数据
const loadingSetup = ref(false)
const setupData = reactive({
  has_key: false,
  key_name: '',
  key_prefix: '',
  permissions: [],
  scopes: [],
  expires_at: null,
  commands: { windows: '', mac: '', linux: '' },
  skill_prompt: '',
})

// 密钥列表
const loadingKeys = ref(false)
const keys = ref([])

// 创建密钥
const showCreateDialog = ref(false)
const creating = ref(false)
const createdToken = ref(null)
const createForm = reactive({
  name: '',
  permissions: ['read'],
  scopes: ['posts', 'pages', 'categories', 'tags', 'media', 'comments', 'settings'],
  expiryType: 'never',
  expiresIn: 30,
})

// 审计日志
const selectedToken = ref(null)
const loadingLogs = ref(false)
const logs = ref([])
const logPage = ref(1)
const logPageSize = 20
const logTotal = ref(0)

// 主题开发 Skill
const loadingThemeSkill = ref(false)
const themeSkillPrompt = ref('')

// 加载设置命令
const loadSetupCommands = async () => {
  loadingSetup.value = true
  try {
    const res = await get('/api/ai/setup-commands')
    const data = res.data || res
    Object.assign(setupData, data)
  } catch (error) {
    // ignore
  } finally {
    loadingSetup.value = false
  }
}

// 加载主题开发 Skill 提示词
const loadThemeSkillPrompt = async () => {
  loadingThemeSkill.value = true
  try {
    const res = await get('/api/ai/themes-skill')
    // API 返回的是纯文本，直接赋值
    themeSkillPrompt.value = typeof res === 'string' ? res : JSON.stringify(res, null, 2)
  } catch (error) {
    ElMessage.error('加载主题开发 Skill 失败')
  } finally {
    loadingThemeSkill.value = false
  }
}

// 加载密钥列表
const loadKeys = async () => {
  loadingKeys.value = true
  try {
    const res = await get('/api/ai/keys')
    keys.value = res.data || res || []
  } catch (error) {
    // ignore
  } finally {
    loadingKeys.value = false
  }
}

// 创建密钥
const handleCreate = async () => {
  if (!createForm.name) {
    ElMessage.warning('请输入密钥名称')
    return
  }

  creating.value = true
  try {
    const payload = {
      name: createForm.name,
      permissions: createForm.permissions,
      scopes: createForm.scopes,
      expires_in: createForm.expiryType === 'never' ? null : createForm.expiresIn,
    }
    const res = await post('/api/ai/keys', payload)
    createdToken.value = res.data || res
    loadKeys()
    loadSetupCommands()
  } catch (error) {
    // handled by interceptor
  } finally {
    creating.value = false
  }
}

// 删除密钥
const handleDelete = (row) => {
  ElMessageBox.confirm(
    '撤销后将无法恢复，使用此 Key 的 AI 客户端将失去访问权限。',
    '确认撤销密钥',
    { confirmButtonText: '撤销', cancelButtonText: '取消', type: 'warning' }
  ).then(async () => {
    try {
      await del('/api/ai/keys/' + row.id)
      ElMessage.success('密钥已撤销')
      loadKeys()
      loadSetupCommands()
    } catch (error) {
      // handled
    }
  }).catch(() => {})
}

// 查看日志
const viewLogs = (row) => {
  selectedToken.value = row
  logPage.value = 1
  activeTab.value = 'logs'
  loadLogs()
}

// 加载日志
const loadLogs = async () => {
  if (!selectedToken.value) return
  loadingLogs.value = true
  try {
    const res = await get('/api/ai/keys/' + selectedToken.value.id + '/logs', {
      page: logPage.value,
      page_size: logPageSize,
    })
    const data = res.data || res
    logs.value = data.list || []
    logTotal.value = data.total || 0
  } catch (error) {
    // handled
  } finally {
    loadingLogs.value = false
  }
}

// 复制到剪贴板
const copyCode = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (e) {
    // fallback
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success('已复制到剪贴板')
  }
}

// 工具函数
const parseJSON = (str) => {
  if (!str) return []
  try { return JSON.parse(str) } catch { return [] }
}

const isExpired = (dateStr) => {
  if (!dateStr) return false
  return new Date(dateStr) < new Date()
}

const isNeverExpires = (dateStr) => {
  if (!dateStr) return false
  const d = new Date(dateStr)
  return d.getFullYear() >= 9999
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (d.getFullYear() >= 9999) return '永不过期'
  return d.toLocaleDateString('zh-CN')
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const getMethodType = (action) => {
  const map = { GET: 'info', POST: 'success', PUT: 'warning', DELETE: 'danger' }
  return map[action] || 'info'
}

// 对话框关闭时重置
const handleCloseDialog = () => {
  if (!createdToken.value) {
    createForm.name = ''
    createForm.permissions = ['read']
    createForm.scopes = ['posts', 'pages', 'categories', 'tags', 'media', 'comments', 'settings']
    createForm.expiryType = 'never'
    createForm.expiresIn = 30
  }
}

onMounted(() => {
  loadSetupCommands()
  loadKeys()
  loadThemeSkillPrompt()
})
</script>

<style scoped>
.ai-settings-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-weight: 600;
  color: var(--admin-text-primary);
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 40px 0;
  color: var(--admin-text-secondary);
}

.empty-icon {
  font-size: 48px;
  color: var(--admin-text-muted);
  margin-bottom: 12px;
}

.empty-hint {
  font-size: var(--admin-font-size-sm);
  color: var(--admin-text-muted);
  margin-top: 8px;
}

/* 设置步骤 */
.setup-step {
  margin-bottom: 24px;
}

.setup-step:last-child {
  margin-bottom: 0;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.step-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: var(--admin-primary);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.step-title {
  font-weight: 600;
  font-size: var(--admin-font-size-lg);
  color: var(--admin-text-primary);
}

.step-desc {
  color: var(--admin-text-secondary);
  font-size: var(--admin-font-size-sm);
  margin-bottom: 12px;
  line-height: 1.6;
}

/* 代码块 */
.code-block-wrapper {
  margin-top: 8px;
  padding: 0 60px;
}

.code-tabs {
  display: flex;
  gap: 2px;
  background-color: var(--admin-bg-hover);
  border-radius: var(--admin-radius-sm) var(--admin-radius-sm) 0 0;
  padding: 2px 2px 0;
}

.code-tab {
  padding: 6px 14px;
  border: none;
  background: transparent;
  color: var(--admin-text-secondary);
  font-size: 12px;
  cursor: pointer;
  border-radius: var(--admin-radius-sm) var(--admin-radius-sm) 0 0;
  transition: all var(--admin-transition-fast);
}

.code-tab.active {
  background-color: var(--admin-bg-surface);
  color: var(--admin-text-primary);
  font-weight: 500;
}

.code-block {
  position: relative;
  background-color: var(--admin-bg-hover);
  border: 1px solid var(--admin-border);
  border-radius: 0 0 var(--admin-radius-sm) var(--admin-radius-sm);
  padding: 12px 40px 12px 14px;
  overflow-x: auto;
}

.code-block.multiline {
  border-radius: var(--admin-radius-sm);
}

.code-block code {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  color: var(--admin-text-primary);
  white-space: pre;
  word-break: break-all;
}

.code-block pre {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  color: var(--admin-text-primary);
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.6;
}

.copy-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 8px;
  border: 1px solid var(--admin-border);
  border-radius: var(--admin-radius-sm);
  background-color: var(--admin-bg-surface);
  color: var(--admin-text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  transition: all var(--admin-transition-fast);
}

.copy-btn:hover {
  color: var(--admin-primary);
  border-color: var(--admin-primary);
}

/* 密钥前缀 */
.key-prefix {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 12px;
  color: var(--admin-text-secondary);
  background-color: var(--admin-bg-hover);
  padding: 2px 8px;
  border-radius: var(--admin-radius-sm);
}

/* 状态文字 */
.expired-text {
  color: var(--admin-danger);
  font-weight: 500;
}

.muted-text {
  color: var(--admin-text-muted);
  font-size: var(--admin-font-size-sm);
}

/* 权限标签 */
.perm-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

/* 过期提示 */
.expire-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  color: var(--admin-warning);
  font-size: var(--admin-font-size-sm);
}

/* 创建成功 Token */
.created-token-box {
  margin-top: 16px;
}

.token-display {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
  padding: 12px;
  background-color: var(--admin-bg-hover);
  border-radius: var(--admin-radius-sm);
}

.token-display code {
  flex: 1;
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 12px;
  color: var(--admin-danger);
  word-break: break-all;
}

/* 快速设置面板宽度限制 */

/* 主题开发 Skill 页面 */
.theme-skill-intro {
  margin-bottom: 20px;
  color: var(--admin-text-secondary);
  line-height: 1.6;
}

.theme-skill-intro .hint {
  font-size: var(--admin-font-size-sm);
  color: var(--admin-text-muted);
  margin-top: 8px;
}

.example-list {
  margin: 12px 0 0 20px;
  color: var(--admin-text-secondary);
  font-size: var(--admin-font-size-sm);
  line-height: 1.8;
}

.example-list li {
  margin-bottom: 4px;
}

.skill-features {
  margin-top: 24px;
  padding: 16px;
  background-color: var(--admin-bg-hover);
  border-radius: var(--admin-radius-sm);
}

.skill-features h3 {
  margin: 0 0 16px 0;
  font-size: var(--admin-font-size-base);
  color: var(--admin-text-primary);
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 12px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background-color: var(--admin-bg-surface);
  border-radius: var(--admin-radius-sm);
  color: var(--admin-text-secondary);
  font-size: var(--admin-font-size-sm);
  transition: all var(--admin-transition-fast);
}

.feature-item:hover {
  color: var(--admin-primary);
  background-color: var(--admin-bg-hover);
}

.feature-item .el-icon {
  font-size: 16px;
  flex-shrink: 0;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
