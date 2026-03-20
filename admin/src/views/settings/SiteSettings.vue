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
              <el-input v-model="settings.basic.site_name" placeholder="请输入站点名称" />
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

      <!-- SEO设置 -->
      <el-tab-pane label="SEO设置" name="seo">
        <el-card>
          <el-form :model="settings.seo" label-width="120px">
            <el-form-item label="网站标题">
              <el-input v-model="settings.seo.meta_title" placeholder="SEO标题" />
            </el-form-item>
            <el-form-item label="网站描述">
              <el-input v-model="settings.seo.meta_description" type="textarea" :rows="3" placeholder="SEO描述" />
            </el-form-item>
            <el-form-item label="网站关键词">
              <el-input v-model="settings.seo.meta_keywords" placeholder="关键词，用逗号分隔" />
            </el-form-item>
            <el-form-item label="启用Sitemap">
              <el-switch v-model="settings.seo.enable_sitemap" />
            </el-form-item>
            <el-form-item label="启用Robots.txt">
              <el-switch v-model="settings.seo.enable_robots" />
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

const settings = reactive({
  basic: {
    site_name: '',
    site_description: '',
    site_url: '',
    contact_email: '',
    timezone: 'Asia/Shanghai',
    logo: '',
    favicon: ''
  },
  seo: {
    meta_title: '',
    meta_description: '',
    meta_keywords: '',
    enable_sitemap: true,
    enable_robots: true
  },
  email: {
    smtp_host: '',
    smtp_port: 587,
    smtp_user: '',
    smtp_password: '',
    sender_name: '',
    encryption: 'tls'
  }
})

const loadSettings = async () => {
  try {
    const res = await get('/api/settings')
    const data = res.data || res || {}
    
    if (data.basic) {
      Object.assign(settings.basic, data.basic)
    }
    if (data.seo) {
      Object.assign(settings.seo, data.seo)
    }
    if (data.email) {
      Object.assign(settings.email, data.email)
    }
  } catch (error) {
    console.error('加载设置失败:', error)
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
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  }
  return false // 阻止默认上传
}

const handleSave = async () => {
  saving.value = true
  try {
    // 保存各个分类的设置
    await put('/api/settings', { category: 'basic', settings: settings.basic })
    await put('/api/settings', { category: 'seo', settings: settings.seo })
    await put('/api/settings', { category: 'email', settings: settings.email })
    ElMessage.success('设置保存成功')
  } catch (error) {
    console.error('保存设置失败:', error)
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
    console.error('发送测试邮件失败:', error)
    ElMessage.error('发送测试邮件失败')
  } finally {
    testingEmail.value = false
  }
}

onMounted(() => {
  loadSettings()
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
</style>
