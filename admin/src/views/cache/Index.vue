<template>
  <div class="cache-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>系统缓存管理</span>
        </div>
      </template>

      <!-- 缓存状态 -->
      <div class="stats-section">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card class="stat-card">
              <el-statistic title="模板缓存" :value="cacheStatus.template">
                <template #suffix>
                  <span class="suffix">个</span>
                </template>
              </el-statistic>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <el-statistic title="数据缓存" :value="cacheStatus.data">
                <template #suffix>
                  <span class="suffix">个</span>
                </template>
              </el-statistic>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <el-statistic title="最后清理" :value="lastCleared" />
            </el-card>
          </el-col>
        </el-row>
      </div>

      <el-divider>一键清理</el-divider>

      <div class="action-section">
        <el-button type="danger" @click="clearAllCache" :loading="clearing">
          <el-icon><Delete /></el-icon>
          清理所有缓存
        </el-button>
        <el-button type="warning" @click="clearTemplateCache" :loading="clearingTemplate">
          <el-icon><Document /></el-icon>
          仅清理模板缓存
        </el-button>
        <el-button type="info" @click="clearDataCache" :loading="clearingData">
          <el-icon><Grid /></el-icon>
          仅清理数据缓存
        </el-button>
      </div>

      <el-divider>操作说明</el-divider>

      <div class="info-section">
        <el-alert type="info" :closable="false">
          <template #title>
            <strong>缓存类型说明：</strong>
            <ul class="cache-info-list">
              <li><strong>模板缓存</strong>：主题模板编译后的缓存，修改模板后建议清理</li>
              <li><strong>数据缓存</strong>：数据库查询结果的缓存，修改内容后建议清理</li>
              <li><strong>所有缓存</strong>：同时清理模板和数据缓存</li>
            </ul>
          </template>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Document, Grid } from '@element-plus/icons-vue'
import axios from 'axios'

// 状态
const clearing = ref(false)
const clearingTemplate = ref(false)
const clearingData = ref(false)
const lastCleared = ref('从未')
const cacheStatus = ref({
  template: 0,
  data: 0
})

// 获取缓存状态
const fetchCacheStatus = async () => {
  try {
    const response = await axios.get('/api/admin/cache/status', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.data.success) {
      cacheStatus.value = response.data.data
    }
  } catch (error) {
    console.error('获取缓存状态失败', error)
  }
}

// 清理所有缓存
const clearAllCache = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理所有缓存吗？这将清除模板和数据缓存。',
      '确认清理',
      { type: 'warning' }
    )

    clearing.value = true
    const response = await axios.post('/api/admin/cache/clear', { type: 'all' }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.data.code === 0) {
      ElMessage.success(response.data.message || '所有缓存已清理完成')
      lastCleared.value = new Date().toLocaleString('zh-CN')
      await fetchCacheStatus()
    } else {
      ElMessage.error(response.data.message || '清理缓存失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清理缓存失败')
    }
  } finally {
    clearing.value = false
  }
}

// 仅清理模板缓存
const clearTemplateCache = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理模板缓存吗？',
      '确认清理',
      { type: 'warning' }
    )

    clearingTemplate.value = true
    const response = await axios.post('/api/admin/cache/clear', { type: 'template' }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.data.code === 0) {
      ElMessage.success(response.data.message || '模板缓存已清理完成')
      lastCleared.value = new Date().toLocaleString('zh-CN')
      await fetchCacheStatus()
    } else {
      ElMessage.error(response.data.message || '清理缓存失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清理缓存失败')
    }
  } finally {
    clearingTemplate.value = false
  }
}

// 仅清理数据缓存
const clearDataCache = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理数据缓存吗？',
      '确认清理',
      { type: 'warning' }
    )

    clearingData.value = true
    const response = await axios.post('/api/admin/cache/clear', { type: 'data' }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.data.code === 0) {
      ElMessage.success(response.data.message || '数据缓存已清理完成')
      lastCleared.value = new Date().toLocaleString('zh-CN')
      await fetchCacheStatus()
    } else {
      ElMessage.error(response.data.message || '清理缓存失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清理缓存失败')
    }
  } finally {
    clearingData.value = false
  }
}

// 组件挂载时获取状态
onMounted(() => {
  fetchCacheStatus()
})
</script>

<style scoped>
.cache-container {
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
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
}

.suffix {
  font-size: 12px;
  color: #909399;
}

.action-section {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
  margin: 20px 0;
}

.info-section {
  margin-top: 20px;
}

.cache-info-list {
  margin: 10px 0 0 0;
  padding-left: 20px;
}

.cache-info-list li {
  margin-bottom: 5px;
  line-height: 1.6;
}
</style>
