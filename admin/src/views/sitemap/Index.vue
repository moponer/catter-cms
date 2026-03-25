<template>
  <div class="sitemap-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>Sitemap 管理</span>
        </div>
      </template>

      <!-- Sitemap Statistics -->
      <div class="stats-section">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stat-card">
              <el-statistic title="总URL数" :value="stats.url_count">
                <template #suffix>
                  <span class="suffix">个</span>
                </template>
              </el-statistic>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <el-statistic title="最后更新" :value="lastUpdated" format="YYYY-MM-DD HH:mm:ss">
              </el-statistic>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- Sitemap Configuration -->
      <el-divider>Sitemap 配置</el-divider>

      <el-form :model="sitemapConfig" label-width="150px" class="sitemap-form">
        <el-form-item label="站点URL">
          <el-input
            v-model="sitemapConfig.base_url"
            placeholder="https://yourdomain.com"
            clearable
          >
            <template #prepend>
              <el-icon><Link /></el-icon>
            </template>
          </el-input>
          <div class="form-tip">
            您的网站基础URL，用于生成sitemap中的链接。请确保以http://或https://开头。
          </div>
        </el-form-item>

        <el-form-item label="更新频率">
          <el-select v-model="sitemapConfig.change_frequency" placeholder="请选择更新频率">
            <el-option label="总是" value="always"></el-option>
            <el-option label="每小时" value="hourly"></el-option>
            <el-option label="每天" value="daily"></el-option>
            <el-option label="每周" value="weekly"></el-option>
            <el-option label="每月" value="monthly"></el-option>
            <el-option label="每年" value="yearly"></el-option>
            <el-option label="从不" value="never"></el-option>
          </el-select>
          <div class="form-tip">
            默认的内容更新频率，影响搜索引擎抓取策略。
          </div>
        </el-form-item>

        <el-form-item label="包含内容类型">
          <el-checkbox-group v-model="sitemapConfig.include_types">
            <el-checkbox label="posts">文章</el-checkbox>
            <el-checkbox label="pages">页面</el-checkbox>
            <el-checkbox label="categories">分类</el-checkbox>
          </el-checkbox-group>
          <div class="form-tip">
            选择要包含在sitemap中的内容类型。
          </div>
        </el-form-item>

        <el-form-item label="文章优先级">
          <el-slider
            v-model="sitemapConfig.post_priority"
            :min="0"
            :max="1"
            :step="0.1"
            :marks="{ 0: '0', 0.5: '0.5', 1: '1' }"
          />
          <div class="form-tip">
            文章在sitemap中的优先级（0-1），影响搜索引擎索引权重。
          </div>
        </el-form-item>

        <el-form-item label="页面优先级">
          <el-slider
            v-model="sitemapConfig.page_priority"
            :min="0"
            :max="1"
            :step="0.1"
            :marks="{ 0: '0', 0.5: '0.5', 1: '1' }"
          />
          <div class="form-tip">
            页面在sitemap中的优先级（0-1），影响搜索引擎索引权重。
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveConfig" :loading="saving">
            <el-icon><Check /></el-icon>
            保存配置
          </el-button>
          <el-button @click="resetConfig">
            <el-icon><RefreshLeft /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>

      <el-divider>Sitemap 操作</el-divider>

      <div class="action-section">
        <el-button type="primary" @click="rebuildSitemap" :loading="rebuilding">
          <el-icon><Refresh /></el-icon>
          重新生成Sitemap
        </el-button>
        <el-button type="success" @click="viewSitemap" :disabled="!sitemapConfig.base_url">
          <el-icon><View /></el-icon>
          查看Sitemap
        </el-button>
        <el-button @click="copySitemapUrl" :disabled="!sitemapConfig.base_url">
          <el-icon><DocumentCopy /></el-icon>
          复制Sitemap URL
        </el-button>
      </div>

      <!-- Sitemap Preview -->
      <el-divider>Sitemap 预览</el-divider>
      <div class="sitemap-preview">
        <el-input
          v-model="sitemapUrl"
          readonly
          placeholder="配置站点URL后显示"
        >
          <template #append>
            <el-button @click="viewSitemap" :disabled="!sitemapConfig.base_url">
              访问
            </el-button>
          </template>
        </el-input>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Link,
  Check,
  RefreshLeft,
  Refresh,
  View,
  DocumentCopy
} from '@element-plus/icons-vue'
import axios from 'axios'

// 状态
const saving = ref(false)
const rebuilding = ref(false)
const lastUpdated = ref('')

// 统计数据
const stats = ref({
  url_count: 0
})

// Sitemap配置
const sitemapConfig = ref({
  base_url: '',
  change_frequency: 'weekly',
  include_types: ['posts', 'pages', 'categories'],
  post_priority: 0.8,
  page_priority: 0.6
})

// 计算属性
const sitemapUrl = computed(() => {
  if (!sitemapConfig.value.base_url) return ''
  return `${sitemapConfig.value.base_url}/sitemap.xml`
})

// 获取Sitemap统计信息
const fetchSitemapInfo = async () => {
  try {
    const response = await axios.get('/api/admin/sitemap/info', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.success) {
      stats.value = response.data.data
    }
  } catch (error) {
    ElMessage.error('获取Sitemap信息失败')
  }
}

// 保存配置
const saveConfig = async () => {
  saving.value = true
  try {
    // 验证基础URL
    if (sitemapConfig.value.base_url && !sitemapConfig.value.base_url.startsWith('http')) {
      ElMessage.error('站点URL必须以http://或https://开头')
      saving.value = false
      return
    }

    // 保存到系统设置
    const settings = {
      sitemap_base_url: sitemapConfig.value.base_url,
      sitemap_change_frequency: sitemapConfig.value.change_frequency,
      sitemap_include_types: sitemapConfig.value.include_types.join(','),
      sitemap_post_priority: sitemapConfig.value.post_priority,
      sitemap_page_priority: sitemapConfig.value.page_priority
    }

    const response = await axios.put('/api/settings', settings, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.data.success) {
      ElMessage.success('配置保存成功')
      lastUpdated.value = new Date().toLocaleString('zh-CN')
    }
  } catch (error) {
    ElMessage.error('保存配置失败')
  } finally {
    saving.value = false
  }
}

// 重置配置
const resetConfig = () => {
  ElMessageBox.confirm('确定要重置所有配置吗？', '确认重置', {
    type: 'warning'
  }).then(() => {
    sitemapConfig.value = {
      base_url: '',
      change_frequency: 'weekly',
      include_types: ['posts', 'pages', 'categories'],
      post_priority: 0.8,
      page_priority: 0.6
    }
    ElMessage.success('配置已重置')
  }).catch(() => {})
}

// 重新生成Sitemap
const rebuildSitemap = async () => {
  rebuilding.value = true
  try {
    const response = await axios.post('/api/admin/sitemap/rebuild', {}, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.data.success) {
      ElMessage.success('Sitemap重新生成成功')
      lastUpdated.value = new Date().toLocaleString('zh-CN')
      await fetchSitemapInfo()
    }
  } catch (error) {
    ElMessage.error('重新生成Sitemap失败')
  } finally {
    rebuilding.value = false
  }
}

// 查看Sitemap
const viewSitemap = () => {
  if (!sitemapConfig.value.base_url) {
    ElMessage.warning('请先配置站点URL')
    return
  }
  window.open(sitemapUrl.value, '_blank')
}

// 复制Sitemap URL
const copySitemapUrl = async () => {
  if (!sitemapConfig.value.base_url) {
    ElMessage.warning('请先配置站点URL')
    return
  }

  try {
    await navigator.clipboard.writeText(sitemapUrl.value)
    ElMessage.success('Sitemap URL已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 加载配置
const loadConfig = async () => {
  try {
    const response = await axios.get('/api/settings', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.data.success && response.data.data) {
      const settings = response.data.data
      if (settings.sitemap_base_url) {
        sitemapConfig.value.base_url = settings.sitemap_base_url
      }
      if (settings.sitemap_change_frequency) {
        sitemapConfig.value.change_frequency = settings.sitemap_change_frequency
      }
      if (settings.sitemap_include_types) {
        sitemapConfig.value.include_types = settings.sitemap_include_types.split(',')
      }
      if (settings.sitemap_post_priority) {
        sitemapConfig.value.post_priority = parseFloat(settings.sitemap_post_priority)
      }
      if (settings.sitemap_page_priority) {
        sitemapConfig.value.page_priority = parseFloat(settings.sitemap_page_priority)
      }
    }
  } catch (error) {
    ElMessage.error('加载配置失败')
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadConfig()
  fetchSitemapInfo()
  lastUpdated.value = new Date().toLocaleString('zh-CN')
})
</script>

<style scoped>
.sitemap-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
}

.stats-section {
  margin-bottom: 30px;
}

.stat-card {
  text-align: center;
}

.suffix {
  font-size: 12px;
  color: #909399;
}

.sitemap-form {
  max-width: 600px;
  margin: 20px 0;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.action-section {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.sitemap-preview {
  margin-top: 20px;
}
</style>
