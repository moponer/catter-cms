<template>
  <div class="theme-settings">
    <el-page-header title="返回" content="主题设置" @back="goBack" />

    <div class="theme-content">
      <!-- 主题信息卡片 -->
      <el-card class="theme-info-card" v-if="currentTheme">
        <div class="theme-header">
          <el-avatar v-if="currentTheme.thumbnail" :size="80" :src="currentTheme.thumbnail" shape="square" />
          <div class="theme-meta">
            <h2>{{ currentTheme.name }}</h2>
            <p class="theme-description">{{ currentTheme.description || '暂无描述' }}</p>
            <div class="theme-tags">
              <el-tag v-if="currentTheme.version" size="small">版本 {{ currentTheme.version }}</el-tag>
              <el-tag v-if="currentTheme.author" size="small" type="info">{{ currentTheme.author }}</el-tag>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 加载中 -->
      <el-skeleton v-if="loading" :rows="6" animated />

      <!-- 配置表单 -->
      <template v-else>
        <el-tabs v-model="activeTab" type="border-card" class="config-tabs">
          <!-- 通用设置 -->
          <el-tab-pane label="通用设置" name="general">
            <div class="tab-content">
              <h3>网站基本信息</h3>
              <el-form :model="config.general" label-position="top">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="网站标题">
                      <el-input v-model="config.general.siteTitle" placeholder="请输入网站标题" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="网站副标题">
                      <el-input v-model="config.general.siteSubtitle" placeholder="请输入网站副标题" />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="网站描述">
                  <el-input v-model="config.general.siteDescription" type="textarea" :rows="3" placeholder="请输入网站描述" />
                </el-form-item>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="Logo">
                      <MediaPicker v-model="config.general.logo" type="image" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="Favicon">
                      <MediaPicker v-model="config.general.favicon" type="image" />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-form>
            </div>
          </el-tab-pane>

          <!-- 模块管理 -->
          <el-tab-pane label="模块管理" name="modules">
            <div class="tab-content">
              <h3>首页模块配置</h3>
              <p class="tab-description">启用/禁用模块，拖拽排序，点击配置模块详情</p>

              <ModuleManager v-model="config.modules" :available-modules="availableModules" />
            </div>
          </el-tab-pane>

          <!-- 页面配置 -->
          <el-tab-pane label="页面配置" name="pages">
            <div class="tab-content">
              <h3>页面模板配置</h3>
              <p class="tab-description">为不同页面选择合适的模板</p>

              <el-form :model="config.pages" label-position="top">
                <el-form-item label="首页模板">
                  <TemplateSelector v-model="config.pages.homeTemplate" :templates="templates.home" />
                </el-form-item>

                <el-form-item label="列表页模板">
                  <TemplateSelector v-model="config.pages.listTemplate" :templates="templates.list" />
                </el-form-item>

                <el-form-item label="详情页模板">
                  <TemplateSelector v-model="config.pages.detailTemplate" :templates="templates.detail" />
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>

          <!-- 高级设置 -->
          <el-tab-pane label="高级设置" name="advanced">
            <div class="tab-content">
              <h3>高级配置</h3>

              <el-form :model="config.advanced" label-position="top">
                <el-form-item label="自定义CSS">
                  <el-input v-model="config.advanced.customCSS" type="textarea" :rows="6" placeholder="输入自定义CSS代码" />
                </el-form-item>

                <el-form-item label="自定义JavaScript">
                  <el-input v-model="config.advanced.customJS" type="textarea" :rows="6" placeholder="输入自定义JavaScript代码" />
                </el-form-item>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="加载动画">
                      <el-switch v-model="config.advanced.loadingAnimation" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="返回顶部按钮">
                      <el-switch v-model="config.advanced.backToTop" />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-form>
            </div>
          </el-tab-pane>
        </el-tabs>

        <!-- 保存按钮 -->
        <div class="actions-bar">
          <el-button @click="resetConfig">重置</el-button>
          <el-button type="primary" :loading="saving" @click="saveConfig">
            <el-icon><Check /></el-icon> 保存配置
          </el-button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Check } from '@element-plus/icons-vue'
import { getThemeConfig, saveThemeConfig, getCurrentTheme } from '@/api/theme'
import MediaPicker from '@/components/MediaPicker.vue'
import ModuleManager from '@/components/ModuleManager.vue'
import TemplateSelector from '@/components/TemplateSelector.vue'

const router = useRouter()

const loading = ref(false)
const saving = ref(false)
const activeTab = ref('general')
const currentTheme = ref(null)

// 配置数据
const config = ref({
  general: {
    siteTitle: '',
    siteSubtitle: '',
    siteDescription: '',
    logo: '',
    favicon: ''
  },
  modules: [],
  pages: {
    homeTemplate: '',
    listTemplate: '',
    detailTemplate: ''
  },
  advanced: {
    customCSS: '',
    customJS: '',
    loadingAnimation: true,
    backToTop: true
  }
})

// 可用模块定义
const availableModules = ref([
  { id: 'hero', name: '首屏大图', icon: 'Monitor', description: '首页顶部的大图展示区域' },
  { id: 'services', name: '服务展示', icon: 'Service', description: '展示服务项目的模块' },
  { id: 'stats', name: '数据统计', icon: 'TrendCharts', description: '展示关键数据的统计模块' },
  { id: 'team', name: '团队介绍', icon: 'UserFilled', description: '展示团队成员的模块' },
  { id: 'news', name: '新闻动态', icon: 'News', description: '展示最新新闻的模块' },
  { id: 'contact', name: '联系我们', icon: 'Phone', description: '联系信息和表单模块' }
])

// 页面模板选项
const templates = ref({
  home: [
    { value: 'default', label: '默认首页', preview: '/templates/previews/home-default.jpg' },
    { value: 'fullwidth', label: '全宽首页', preview: '/templates/previews/home-fullwidth.jpg' },
    { value: 'sidebar', label: '侧边栏首页', preview: '/templates/previews/home-sidebar.jpg' }
  ],
  list: [
    { value: 'grid', label: '网格布局', preview: '/templates/previews/list-grid.jpg' },
    { value: 'list', label: '列表布局', preview: '/templates/previews/list-list.jpg' },
    { value: 'masonry', label: '瀑布流布局', preview: '/templates/previews/list-masonry.jpg' }
  ],
  detail: [
    { value: 'default', label: '默认详情', preview: '/templates/previews/detail-default.jpg' },
    { value: 'fullwidth', label: '全宽详情', preview: '/templates/previews/detail-fullwidth.jpg' },
    { value: 'sidebar', label: '侧边栏详情', preview: '/templates/previews/detail-sidebar.jpg' }
  ]
})

// 初始化
onMounted(async () => {
  await loadThemeData()
})

const loadThemeData = async () => {
  loading.value = true
  try {
    // 获取当前主题信息
    const themeRes = await getCurrentTheme()
    currentTheme.value = themeRes.data

    // 获取主题配置
    const configRes = await getThemeConfig()
    if (configRes.data) {
      mergeConfig(configRes.data)
    }
  } catch (error) {
    ElMessage.error('加载主题数据失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const mergeConfig = (data) => {
  if (data.general) Object.assign(config.value.general, data.general)
  if (data.modules) config.value.modules = data.modules
  if (data.pages) Object.assign(config.value.pages, data.pages)
  if (data.advanced) Object.assign(config.value.advanced, data.advanced)
}

const saveConfig = async () => {
  saving.value = true
  try {
    await saveThemeConfig(config.value)
    ElMessage.success('配置保存成功')
  } catch (error) {
    ElMessage.error('保存失败，请重试')
    console.error(error)
  } finally {
    saving.value = false
  }
}

const resetConfig = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重置所有配置吗？这将恢复到上次保存的状态。',
      '确认重置',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await loadThemeData()
    ElMessage.success('配置已重置')
  } catch {
    // 用户取消
  }
}

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.theme-settings {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.theme-content {
  margin-top: 20px;
}

.theme-info-card {
  margin-bottom: 20px;
}

.theme-header {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.theme-meta h2 {
  margin: 0 0 8px;
  font-size: 20px;
  color: #303133;
}

.theme-description {
  margin: 0 0 12px;
  color: #606266;
  line-height: 1.5;
}

.theme-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.config-tabs {
  margin-bottom: 20px;
}

.tab-content {
  padding: 20px;
}

.tab-content h3 {
  margin: 0 0 8px;
  font-size: 16px;
  color: #303133;
}

.tab-description {
  margin: 0 0 20px;
  color: #909399;
  font-size: 14px;
}

.actions-bar {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
</style>
