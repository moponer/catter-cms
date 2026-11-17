<template>
  <div class="settings-page">
    <div class="page-header">
      <h1>系统设置</h1>
      <el-button type="primary" @click="handleSave" :loading="saving">保存设置</el-button>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 基本设置 -->
      <el-tab-pane label="基本设置" name="basic">
        <el-card>
          <el-form :model="settings.basic" label-width="120px">
            <el-form-item label="站点名称">
              <el-input v-model="settings.basic.site_title" placeholder="请输入站点名称" />
            </el-form-item>
            <el-form-item label="站点描述">
              <el-input v-model="settings.basic.site_description" type="textarea" :rows="3" placeholder="站点描述" />
            </el-form-item>
            <el-form-item label="站点URL">
              <el-input v-model="settings.basic.site_url" placeholder="https://example.com" />
            </el-form-item>
            <el-form-item label="联系邮箱">
              <el-input v-model="settings.basic.contact_email" placeholder="admin@example.com" />
            </el-form-item>
            <el-form-item label="时区">
              <el-select v-model="settings.basic.timezone" placeholder="选择时区" style="width: 100%">
                <el-option label="UTC" value="UTC" />
                <el-option label="北京时间 (UTC+8)" value="Asia/Shanghai" />
                <el-option label="纽约 (UTC-5)" value="America/New_York" />
              </el-select>
            </el-form-item>
            <el-form-item label="站点Logo">
              <div class="upload-container">
                <el-upload
                  class="logo-uploader"
                  :show-file-list="false"
                  :before-upload="(file) => handleUpload(file, 'logo')"
                  accept="image/*"
                >
                  <img v-if="settings.basic.logo" :src="settings.basic.logo" class="logo-preview" />
                  <el-icon v-else class="logo-uploader-icon"><Plus /></el-icon>
                </el-upload>
                <el-button v-if="settings.basic.logo" type="danger" link @click="settings.basic.logo = ''">
                  移除Logo
                </el-button>
              </div>
              <div class="form-hint">推荐尺寸: 200x50px，支持 PNG/JPG/SVG</div>
            </el-form-item>
            <el-form-item label="网站图标">
              <div class="upload-container">
                <el-upload
                  class="favicon-uploader"
                  :show-file-list="false"
                  :before-upload="(file) => handleUpload(file, 'favicon')"
                  accept="image/*"
                >
                  <img v-if="settings.basic.favicon" :src="settings.basic.favicon" class="favicon-preview" />
                  <el-icon v-else class="favicon-uploader-icon"><Plus /></el-icon>
                </el-upload>
                <el-button v-if="settings.basic.favicon" type="danger" link @click="settings.basic.favicon = ''">
                  移除图标
                </el-button>
              </div>
              <div class="form-hint">推荐尺寸: 32x32px 或 16x16px，支持 ICO/PNG</div>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 邮件设置 -->
      <el-tab-pane label="邮件设置" name="email">
        <el-card>
          <el-form :model="settings.email" label-width="120px">
            <el-form-item label="SMTP服务器">
              <el-input v-model="settings.email.smtp_host" placeholder="smtp.example.com" />
            </el-form-item>
            <el-form-item label="SMTP端口">
              <el-input-number v-model="settings.email.smtp_port" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="发件人邮箱">
              <el-input v-model="settings.email.smtp_user" placeholder="noreply@example.com" />
            </el-form-item>
            <el-form-item label="SMTP密码">
              <el-input v-model="settings.email.smtp_password" type="password" show-password placeholder="SMTP密码" />
            </el-form-item>
            <el-form-item label="发件人名称">
              <el-input v-model="settings.email.sender_name" placeholder="系统通知" />
            </el-form-item>
            <el-form-item label="加密方式">
              <el-radio-group v-model="settings.email.encryption">
                <el-radio label="none">无</el-radio>
                <el-radio label="ssl">SSL</el-radio>
                <el-radio label="tls">TLS</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" plain @click="handleTestEmail" :loading="testingEmail">
                发送测试邮件
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 上传设置 -->
      <el-tab-pane label="上传设置" name="upload">
        <el-card>
          <el-form :model="settings.upload" label-width="140px">
            <el-form-item label="图片压缩">
              <el-switch v-model="settings.upload.image_compress" />
              <div class="form-hint">开启后自动压缩上传的图片，减小文件体积</div>
            </el-form-item>
            <el-form-item label="最大宽度（像素）">
              <el-input-number v-model="settings.upload.max_width" :min="400" :max="8000" :step="100" />
              <div class="form-hint">图片宽度超过此值时自动缩放。0 表示不限制</div>
            </el-form-item>
            <el-form-item label="JPEG 压缩质量">
              <el-slider v-model="settings.upload.jpeg_quality" :min="60" :max="100" :show-tooltip="true" style="max-width: 300px;" />
              <div class="form-hint">质量越高文件越大，推荐 85。60-100 范围内调节</div>
            </el-form-item>
            <el-divider content-position="left">当前效果预览</el-divider>
            <div class="upload-preview">
              <div class="preview-item">
                <span class="preview-label">压缩开关：</span>
                <el-tag :type="settings.upload.image_compress ? 'success' : 'info'">
                  {{ settings.upload.image_compress ? '已开启' : '已关闭' }}
                </el-tag>
              </div>
              <div class="preview-item">
                <span class="preview-label">最大宽度：</span>
                <span class="preview-value">{{ settings.upload.max_width }}px</span>
              </div>
              <div class="preview-item">
                <span class="preview-label">JPEG 质量：</span>
                <span class="preview-value">{{ settings.upload.jpeg_quality }}%</span>
              </div>
            </div>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 存储设置 -->
      <el-tab-pane label="存储设置" name="storage">
        <el-card>
          <el-form :model="storageConfig" label-width="150px">
            <el-form-item label="存储方式">
              <el-radio-group v-model="storageConfig.type" @change="onStorageTypeChange">
                <el-radio value="local">本地存储</el-radio>
                <el-radio value="qiniu">七牛云</el-radio>
                <el-radio value="alioss">阿里云 OSS</el-radio>
                <el-radio value="tencentcos">腾讯云 COS</el-radio>
                <el-radio value="s3">S3 兼容存储</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="CDN 域名" v-if="storageConfig.type !== 'local'">
              <el-input v-model="storageConfig.cdn" placeholder="可选，全局覆盖访问域名" />
              <div class="form-hint">填写后所有文件链接将使用此域名，留空则使用各存储商默认域名</div>
            </el-form-item>

            <!-- 本地存储说明 -->
            <template v-if="storageConfig.type === 'local'">
              <el-form-item>
                <el-alert type="info" :closable="false" show-icon>
                  <template #title>本地存储模式</template>
                  文件将保存在服务器的 <code>uploads/</code> 目录下。适合开发和测试环境。
                </el-alert>
              </el-form-item>
            </template>

            <!-- 七牛云配置 -->
            <template v-if="storageConfig.type === 'qiniu'">
              <el-divider content-position="left">七牛云配置</el-divider>
              <el-form-item label="Access Key">
                <el-input v-model="storageConfig.qiniu.access_key" placeholder="七牛云 Access Key" />
              </el-form-item>
              <el-form-item label="Secret Key">
                <el-input v-model="storageConfig.qiniu.secret_key" type="password" show-password placeholder="七牛云 Secret Key" />
              </el-form-item>
              <el-form-item label="Bucket">
                <el-input v-model="storageConfig.qiniu.bucket" placeholder="存储空间名称" />
              </el-form-item>
              <el-form-item label="存储区域">
                <el-select v-model="storageConfig.qiniu.zone" placeholder="选择区域" style="width: 100%">
                  <el-option label="华东 (z0)" value="huadong" />
                  <el-option label="华北 (z1)" value="huabei" />
                  <el-option label="华南 (z2)" value="huanan" />
                  <el-option label="北美 (na0)" value="beimei" />
                  <el-option label="东南亚 (as0)" value="xinjiapo" />
                </el-select>
              </el-form-item>
              <el-form-item label="访问域名">
                <el-input v-model="storageConfig.qiniu.domain" placeholder="如 http://cdn.example.com" />
              </el-form-item>
              <el-form-item label="使用 HTTPS">
                <el-switch v-model="storageConfig.qiniu.use_https" />
              </el-form-item>
            </template>

            <!-- 阿里云 OSS 配置 -->
            <template v-if="storageConfig.type === 'alioss'">
              <el-divider content-position="left">阿里云 OSS 配置</el-divider>
              <el-form-item label="Endpoint">
                <el-input v-model="storageConfig.alioss.endpoint" placeholder="如 https://oss-cn-hangzhou.aliyuncs.com" />
              </el-form-item>
              <el-form-item label="AccessKey ID">
                <el-input v-model="storageConfig.alioss.access_key_id" placeholder="阿里云 AccessKey ID" />
              </el-form-item>
              <el-form-item label="AccessKey Secret">
                <el-input v-model="storageConfig.alioss.access_key_secret" type="password" show-password placeholder="阿里云 AccessKey Secret" />
              </el-form-item>
              <el-form-item label="Bucket">
                <el-input v-model="storageConfig.alioss.bucket" placeholder="存储空间名称" />
              </el-form-item>
              <el-form-item label="自定义域名">
                <el-input v-model="storageConfig.alioss.custom_domain" placeholder="可选，如 https://cdn.example.com" />
              </el-form-item>
            </template>

            <!-- 腾讯云 COS 配置 -->
            <template v-if="storageConfig.type === 'tencentcos'">
              <el-divider content-position="left">腾讯云 COS 配置</el-divider>
              <el-form-item label="Bucket URL">
                <el-input v-model="storageConfig.tencentcos.bucket_url" placeholder="如 https://bucket-1234567890.cos.ap-guangzhou.myqcloud.com" />
              </el-form-item>
              <el-form-item label="Secret ID">
                <el-input v-model="storageConfig.tencentcos.secret_id" placeholder="腾讯云 Secret ID" />
              </el-form-item>
              <el-form-item label="Secret Key">
                <el-input v-model="storageConfig.tencentcos.secret_key" type="password" show-password placeholder="腾讯云 Secret Key" />
              </el-form-item>
              <el-form-item label="Region">
                <el-input v-model="storageConfig.tencentcos.region" placeholder="如 ap-guangzhou、ap-beijing" />
                <div class="form-hint">与 Bucket URL 中的 Region 一致</div>
              </el-form-item>
              <el-form-item label="自定义域名">
                <el-input v-model="storageConfig.tencentcos.custom_domain" placeholder="可选，如 https://cdn.example.com" />
              </el-form-item>
            </template>

            <!-- S3 兼容配置 -->
            <template v-if="storageConfig.type === 's3'">
              <el-divider content-position="left">S3 兼容存储配置</el-divider>
              <div class="form-hint" style="margin-bottom: 16px;">
                支持 AWS S3、Cloudflare R2、MinIO 等 S3 兼容存储
              </div>
              <el-form-item label="Endpoint">
                <el-input v-model="storageConfig.s3.endpoint" placeholder="如 https://s3.amazonaws.com 或 MinIO 地址" />
              </el-form-item>
              <el-form-item label="Region">
                <el-input v-model="storageConfig.s3.region" placeholder="如 us-east-1" />
              </el-form-item>
              <el-form-item label="Access Key ID">
                <el-input v-model="storageConfig.s3.access_key_id" placeholder="Access Key ID" />
              </el-form-item>
              <el-form-item label="Secret Access Key">
                <el-input v-model="storageConfig.s3.secret_access_key" type="password" show-password placeholder="Secret Access Key" />
              </el-form-item>
              <el-form-item label="Bucket">
                <el-input v-model="storageConfig.s3.bucket" placeholder="存储桶名称" />
              </el-form-item>
              <el-form-item label="Path Style">
                <el-switch v-model="storageConfig.s3.force_path_style" />
                <div class="form-hint">MinIO 等兼容存储需开启。AWS S3 / Cloudflare R2 保持关闭</div>
              </el-form-item>
              <el-form-item label="自定义域名">
                <el-input v-model="storageConfig.s3.custom_domain" placeholder="可选，如 Cloudflare R2 公开域名" />
              </el-form-item>
            </template>

            <!-- 操作按钮 -->
            <el-divider />
            <el-form-item>
              <el-button type="primary" @click="handleTestStorage" :loading="testingStorage" :disabled="storageConfig.type === 'local'">
                测试连接
              </el-button>
              <el-button type="success" @click="handleSaveStorage" :loading="savingStorage">
                保存并切换
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 缓存设置 -->
      <el-tab-pane label="缓存设置" name="cache">
        <el-card>
          <el-form :model="redisConfig" label-width="150px">
            <el-form-item label="启用 Redis">
              <el-switch v-model="redisConfig.enabled" />
              <div class="form-hint">开启后系统将使用 Redis 作为缓存后端，提升性能</div>
            </el-form-item>

            <template v-if="redisConfig.enabled">
              <el-divider content-position="left">Redis 连接配置</el-divider>

              <el-form-item label="Redis 地址">
                <el-input v-model="redisConfig.host" placeholder="127.0.0.1" />
              </el-form-item>

              <el-form-item label="端口">
                <el-input-number v-model="redisConfig.port" :min="1" :max="65535" :step="1" />
              </el-form-item>

              <el-form-item label="密码">
                <el-input v-model="redisConfig.password" type="password" show-password placeholder="留空表示无密码" />
              </el-form-item>

              <el-form-item label="数据库编号">
                <el-input-number v-model="redisConfig.db" :min="0" :max="15" :step="1" />
                <div class="form-hint">Redis 数据库编号（0-15），默认使用 0</div>
              </el-form-item>

              <el-form-item label="Key 前缀">
                <el-input v-model="redisConfig.key_prefix" placeholder="catter:" />
                <div class="form-hint">用于区分不同项目的缓存，避免 Key 冲突</div>
              </el-form-item>

              <!-- 连接状态 -->
              <el-form-item label="当前状态" v-if="redisConfig.status">
                <el-tag :type="redisConfig.status === 'connected' ? 'success' : redisConfig.status === 'disabled' ? 'info' : 'danger'">
                  {{ redisConfig.status === 'connected' ? '已连接' : redisConfig.status === 'disabled' ? '未启用' : '连接异常' }}
                </el-tag>
                <span v-if="redisConfig.status === 'connected' && redisConfig.keys !== undefined" style="margin-left: 12px; color: var(--el-text-color-secondary);">
                  当前缓存 Keys: {{ redisConfig.keys }}
                </span>
              </el-form-item>
            </template>

            <template v-else>
              <el-form-item>
                <el-alert type="info" :closable="false" show-icon>
                  <template #title>缓存未启用</template>
                  系统当前未使用 Redis 缓存。开启后可提升页面加载速度和 API 响应性能。
                  需要确保服务器已安装并运行 Redis 服务。
                </el-alert>
              </el-form-item>
            </template>

            <!-- 操作按钮 -->
            <el-divider />
            <el-form-item>
              <el-button type="primary" @click="handleTestRedis" :loading="testingRedis" :disabled="!redisConfig.enabled">
                测试连接
              </el-button>
              <el-button type="success" @click="handleSaveRedis" :loading="savingRedis">
                保存配置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { get, put, post } from '@/utils/request'
import { Plus } from '@element-plus/icons-vue'

const activeTab = ref('basic')
const saving = ref(false)
const testingEmail = ref(false)
const savingStorage = ref(false)
const testingStorage = ref(false)
const savingRedis = ref(false)
const testingRedis = ref(false)

const settings = reactive({
  basic: {
    site_title: '',
    site_description: '',
    site_url: '',
    contact_email: '',
    timezone: 'Asia/Shanghai',
    logo: '',
    favicon: ''
  },
  email: {
    smtp_host: '',
    smtp_port: 587,
    smtp_user: '',
    smtp_password: '',
    sender_name: '',
    encryption: 'tls'
  },
  upload: {
    image_compress: true,
    max_width: 1920,
    jpeg_quality: 85
  }
})

// 存储配置（独立于 settings，走 config.yaml）
const storageConfig = reactive({
  type: 'local',
  cdn: '',
  qiniu: {
    access_key: '',
    secret_key: '',
    bucket: '',
    zone: 'huadong',
    domain: '',
    use_https: true
  },
  alioss: {
    endpoint: '',
    access_key_id: '',
    access_key_secret: '',
    bucket: '',
    custom_domain: ''
  },
  tencentcos: {
    bucket_url: '',
    secret_id: '',
    secret_key: '',
    region: '',
    custom_domain: ''
  },
  s3: {
    endpoint: '',
    region: '',
    access_key_id: '',
    secret_access_key: '',
    bucket: '',
    force_path_style: false,
    custom_domain: ''
  }
})

// Redis 缓存配置（独立于 settings，走 config.yaml）
const redisConfig = reactive({
  enabled: false,
  host: '127.0.0.1',
  port: 6379,
  password: '',
  db: 0,
  key_prefix: 'catter:',
  status: '',
  keys: undefined
})

const loadSettings = async () => {
  try {
    const res = await get('/api/settings')
    const data = res.data || res || {}

    // 后端返回扁平结构，需要手动映射
    // basic
    if (data.site_title !== undefined) settings.basic.site_title = data.site_title
    if (data.site_description !== undefined) settings.basic.site_description = data.site_description
    if (data.site_url !== undefined) settings.basic.site_url = data.site_url
    if (data.contact_email !== undefined) settings.basic.contact_email = data.contact_email
    if (data.timezone !== undefined) settings.basic.timezone = data.timezone
    if (data.logo !== undefined) settings.basic.logo = data.logo
    if (data.favicon !== undefined) settings.basic.favicon = data.favicon

    // email
    if (data.smtp_host !== undefined) settings.email.smtp_host = data.smtp_host
    if (data.smtp_port !== undefined) settings.email.smtp_port = data.smtp_port
    if (data.smtp_user !== undefined) settings.email.smtp_user = data.smtp_user
    if (data.smtp_password !== undefined) settings.email.smtp_password = data.smtp_password
    if (data.sender_name !== undefined) settings.email.sender_name = data.sender_name
    if (data.encryption !== undefined) settings.email.encryption = data.encryption

    // upload
    if (data.image_compress !== undefined) settings.upload.image_compress = Boolean(data.image_compress)
    if (data.max_width !== undefined) settings.upload.max_width = Number(data.max_width)
    if (data.jpeg_quality !== undefined) settings.upload.jpeg_quality = Number(data.jpeg_quality)
  } catch (error) {
    ElMessage.error('加载设置失败')
  }
}

const handleUpload = async (file, type) => {
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    const res = await post('/api/media/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const url = res.data?.url || res.url
    if (type === 'logo') {
      settings.basic.logo = url
    } else if (type === 'favicon') {
      settings.basic.favicon = url
    }
    ElMessage.success('上传成功')
  } catch (error) {
    ElMessage.error('上传失败')
  }
  return false // 阻止默认上传
}

const handleSave = async () => {
  saving.value = true
  try {
    // 合并所有设置并发送
    const allSettings = {
      ...settings.basic,
      ...settings.email,
      ...settings.upload
    }
    await put('/api/settings', allSettings)
    ElMessage.success('设置保存成功')
  } catch (error) {
    ElMessage.error('保存设置失败')
  } finally {
    saving.value = false
  }
}

const handleTestEmail = async () => {
  testingEmail.value = true
  try {
    await post('/api/settings/test-email', settings.email)
    ElMessage.success('测试邮件发送成功')
  } catch (error) {
    ElMessage.error('发送测试邮件失败')
  } finally {
    testingEmail.value = false
  }
}

// ========== 存储配置 ==========

const loadStorageConfig = async () => {
  try {
    const res = await get('/api/storage/config')
    const data = res.data || res || {}

    storageConfig.type = data.type || 'local'
    storageConfig.cdn = data.cdn || ''

    if (data.qiniu) {
      Object.assign(storageConfig.qiniu, data.qiniu)
    }
    if (data.alioss) {
      Object.assign(storageConfig.alioss, data.alioss)
    }
    if (data.tencentcos) {
      Object.assign(storageConfig.tencentcos, data.tencentcos)
    }
    if (data.s3) {
      Object.assign(storageConfig.s3, data.s3)
    }
  } catch (error) {
    // 静默失败，使用默认值
  }
}

const onStorageTypeChange = () => {
  // 切换存储类型时不需要额外操作，UI 已通过 v-if 自动切换
}

const buildStoragePayload = () => {
  const payload = {
    type: storageConfig.type,
    cdn: storageConfig.cdn
  }

  switch (storageConfig.type) {
    case 'qiniu':
      payload.qiniu = { ...storageConfig.qiniu }
      break
    case 'alioss':
      payload.alioss = { ...storageConfig.alioss }
      break
    case 'tencentcos':
      payload.tencentcos = { ...storageConfig.tencentcos }
      break
    case 's3':
      payload.s3 = { ...storageConfig.s3 }
      break
  }

  return payload
}

const handleTestStorage = async () => {
  testingStorage.value = true
  try {
    const payload = buildStoragePayload()
    await post('/api/storage/config/test', payload)
    ElMessage.success('连接测试成功，存储配置有效')
  } catch (error) {
    const msg = error?.response?.data?.message || error?.message || '连接测试失败'
    ElMessage.error(msg)
  } finally {
    testingStorage.value = false
  }
}

const handleSaveStorage = async () => {
  savingStorage.value = true
  try {
    const payload = buildStoragePayload()
    await put('/api/storage/config', payload)
    ElMessage.success(`存储已切换为：${storageTypeLabel(storageConfig.type)}`)
    // 重新加载以获取脱敏后的值
    await loadStorageConfig()
  } catch (error) {
    const msg = error?.response?.data?.message || error?.message || '保存存储配置失败'
    ElMessage.error(msg)
  } finally {
    savingStorage.value = false
  }
}

const storageTypeLabel = (type) => {
  const labels = {
    local: '本地存储',
    qiniu: '七牛云',
    alioss: '阿里云 OSS',
    tencentcos: '腾讯云 COS',
    s3: 'S3 兼容存储'
  }
  return labels[type] || type
}

// ========== Redis 缓存配置 ==========

const loadRedisConfig = async () => {
  try {
    const res = await get('/api/redis/config')
    const data = res.data || res || {}

    redisConfig.enabled = data.enabled || false
    redisConfig.host = data.host || '127.0.0.1'
    redisConfig.port = data.port || 6379
    redisConfig.password = data.password || ''
    redisConfig.db = data.db || 0
    redisConfig.key_prefix = data.key_prefix || 'catter:'
    redisConfig.status = data.status || ''
    redisConfig.keys = data.keys
  } catch (error) {
    // 静默失败，使用默认值
  }
}

const handleTestRedis = async () => {
  testingRedis.value = true
  try {
    const payload = {
      enabled: true,
      host: redisConfig.host,
      port: redisConfig.port,
      password: redisConfig.password,
      db: redisConfig.db,
      key_prefix: redisConfig.key_prefix
    }
    await post('/api/redis/config/test', payload)
    ElMessage.success('Redis 连接测试成功')
  } catch (error) {
    const msg = error?.response?.data?.message || error?.message || '连接测试失败'
    ElMessage.error(msg)
  } finally {
    testingRedis.value = false
  }
}

const handleSaveRedis = async () => {
  savingRedis.value = true
  try {
    const payload = {
      enabled: redisConfig.enabled,
      host: redisConfig.host,
      port: redisConfig.port,
      password: redisConfig.password,
      db: redisConfig.db,
      key_prefix: redisConfig.key_prefix
    }
    await put('/api/redis/config', payload)
    ElMessage.success(redisConfig.enabled ? 'Redis 缓存已启用' : 'Redis 缓存已禁用')
    // 重新加载以获取最新状态
    await loadRedisConfig()
  } catch (error) {
    const msg = error?.response?.data?.message || error?.message || '保存 Redis 配置失败'
    ElMessage.error(msg)
  } finally {
    savingRedis.value = false
  }
}

onMounted(() => {
  loadSettings()
  loadStorageConfig()
  loadRedisConfig()
})
</script>

<style scoped>
.settings-page {
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
  margin: 0;
}

.upload-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-uploader,
.favicon-uploader {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration);
}

.logo-uploader:hover,
.favicon-uploader:hover {
  border-color: var(--el-color-primary);
}

.logo-uploader-icon,
.favicon-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 120px;
  height: 60px;
  text-align: center;
  line-height: 60px;
}

.favicon-uploader-icon {
  width: 60px;
  height: 60px;
  line-height: 60px;
}

.logo-preview {
  width: 120px;
  height: 60px;
  object-fit: contain;
  display: block;
}

.favicon-preview {
  width: 60px;
  height: 60px;
  object-fit: contain;
  display: block;
}

.form-hint {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  margin-top: 4px;
}

.upload-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preview-label {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  min-width: 100px;
}

.preview-value {
  font-weight: 600;
  font-size: 14px;
  color: var(--el-color-primary);
}
</style>
