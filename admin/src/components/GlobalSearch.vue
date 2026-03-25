<template>
  <div class="global-search">
    <!-- 全局搜索按钮 -->
    <el-button
      type="primary"
      :icon="Search"
      circle
      size="large"
      @click="openSearch"
      class="search-trigger-btn"
    />
    
    <!-- 全局搜索对话框 -->
    <el-dialog
      v-model="searchVisible"
      title="全局搜索"
      :width="'600px'"
      :before-close="handleClose"
      class="search-dialog"
      destroy-on-close
    >
      <div class="search-container">
        <!-- 搜索输入框 -->
        <div class="search-input-wrapper">
          <el-input
            v-model="searchQuery"
            placeholder="输入关键词搜索 (Ctrl + K)"
            :prefix-icon="Search"
            size="large"
            clearable
            @input="handleSearchInput"
            @keydown="handleKeydown"
            ref="searchInputRef"
            autofocus
          />
        </div>

        <!-- 搜索类型选择 -->
        <div class="search-filters">
          <el-radio-group v-model="searchType" size="small" @change="handleSearch">
            <el-radio-button label="all">全部</el-radio-button>
            <el-radio-button label="post">文章</el-radio-button>
            <el-radio-button label="page">页面</el-radio-button>
            <el-radio-button label="comment">评论</el-radio-button>
          </el-radio-group>
        </div>

        <!-- 搜索结果 -->
        <div class="search-results" v-loading="searching">
          <div v-if="!searchQuery" class="search-placeholder">
            <el-empty description="输入关键词开始搜索" />
            <div class="search-tips">
              <p>💡 提示:</p>
              <ul>
                <li>使用 Ctrl + K 快捷键快速打开搜索</li>
                <li>支持搜索文章、页面、评论内容</li>
                <li>点击结果可以直接跳转</li>
              </ul>
            </div>
          </div>

          <div v-else-if="searchResults.length === 0 && !searching" class="no-results">
            <el-empty description="未找到相关结果" />
          </div>

          <div v-else class="results-list">
            <div
              v-for="(result, index) in searchResults"
              :key="result.id"
              class="result-item"
              :class="{ 'selected': selectedIndex === index }"
              @click="handleResultClick(result)"
              @mouseenter="selectedIndex = index"
            >
              <div class="result-header">
                <el-tag :type="getTypeColor(result.type)" size="small">
                  {{ getTypeName(result.type) }}
                </el-tag>
                <span class="result-score">相关度: {{ (result.score * 100).toFixed(0) }}%</span>
              </div>
              
              <div class="result-title">{{ result.title || result.content.substring(0, 50) }}</div>
              
              <div class="result-meta" v-if="result.type === 'post' || result.type === 'page'">
                <span v-if="result.excerpt" class="result-excerpt">{{ result.excerpt }}</span>
                <span v-if="result.category" class="result-category">
                  <el-tag size="small" type="info">{{ result.category }}</el-tag>
                </span>
                <span v-if="result.tags && result.tags.length > 0" class="result-tags">
                  <el-tag
                    v-for="tag in result.tags.slice(0, 3)"
                    :key="tag"
                    size="small"
                    type="success"
                    style="margin-right: 5px"
                  >
                    {{ tag }}
                  </el-tag>
                </span>
              </div>

              <div class="result-meta" v-if="result.type === 'comment'">
                <span class="result-author">{{ result.author_name }}</span>
                <span class="result-content">{{ result.content }}</span>
                <span v-if="result.post_title" class="result-post-title">
                  文章: {{ result.post_title }}
                </span>
              </div>

              <div class="result-footer">
                <span class="result-time">{{ formatDate(result.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 搜索结果统计 -->
        <div class="search-footer" v-if="searchQuery && !searching">
          <span>找到 {{ totalResults }} 条结果</span>
          <el-pagination
            v-if="totalResults > limit"
            v-model:current-page="currentPage"
            :page-size="limit"
            :total="totalResults"
            layout="prev, pager, next"
            small
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()

// 搜索状态
const searchVisible = ref(false)
const searchQuery = ref('')
const searchType = ref('all')
const searchResults = ref([])
const searching = ref(false)
const selectedIndex = ref(0)
const currentPage = ref(1)
const limit = ref(10)
const totalResults = ref(0)

// 搜索输入框引用
const searchInputRef = ref(null)

// 搜索防抖定时器
let searchDebounceTimer = null

// 打开搜索
const openSearch = () => {
  searchVisible.value = true
  searchQuery.value = ''
  searchResults.value = []
  selectedIndex.value = 0
  nextTick(() => {
    if (searchInputRef.value) {
      searchInputRef.value.focus()
    }
  })
}

// 关闭搜索
const handleClose = () => {
  searchVisible.value = false
  searchQuery.value = ''
  searchResults.value = []
}

// 搜索输入处理
const handleSearchInput = () => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
  }
  searchDebounceTimer = setTimeout(() => {
    if (searchQuery.value.trim()) {
      handleSearch()
    }
  }, 300)
}

// 执行搜索
const handleSearch = async () => {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    totalResults.value = 0
    return
  }

  searching.value = true
  selectedIndex.value = 0

  try {
    const response = await request.get('/api/search', {
      params: {
        q: searchQuery.value,
        type: searchType.value,
        limit: limit.value,
        page: currentPage.value,
      }
    })

    if (response.success) {
      searchResults.value = response.data.results
      totalResults.value = response.data.total
    } else {
      ElMessage.error(response.message)
      searchResults.value = []
    }
  } catch (error) {
    ElMessage.error('搜索失败,请稍后重试')
    searchResults.value = []
  } finally {
    searching.value = false
  }
}

// 键盘导航
const handleKeydown = (event) => {
  if (event.key === 'ArrowDown') {
    event.preventDefault()
    if (selectedIndex.value < searchResults.value.length - 1) {
      selectedIndex.value++
    }
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    if (selectedIndex.value > 0) {
      selectedIndex.value--
    }
  } else if (event.key === 'Enter') {
    event.preventDefault()
    if (searchResults.value.length > 0 && selectedIndex.value < searchResults.value.length) {
      handleResultClick(searchResults.value[selectedIndex.value])
    }
  } else if (event.key === 'Escape') {
    handleClose()
  }
}

// 点击搜索结果
const handleResultClick = (result) => {
  if (result.type === 'post') {
    router.push(`/admin/posts/${result.id}`)
  } else if (result.type === 'page') {
    router.push(`/admin/pages/${result.id}`)
  } else if (result.type === 'comment') {
    router.push(`/admin/comments`)
  }
  handleClose()
}

// 分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  handleSearch()
}

// 获取类型名称
const getTypeName = (type) => {
  const names = {
    post: '文章',
    page: '页面',
    comment: '评论',
  }
  return names[type] || type
}

// 获取类型颜色
const getTypeColor = (type) => {
  const colors = {
    post: 'primary',
    page: 'success',
    comment: 'warning',
  }
  return colors[type] || 'info'
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

// 键盘快捷键监听
const handleGlobalKeydown = (event) => {
  // Ctrl + K 或 Cmd + K
  if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
    event.preventDefault()
    if (searchVisible.value) {
      handleClose()
    } else {
      openSearch()
    }
  }
}

// 挂载和卸载
onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
  window.addEventListener('open-global-search', handleOpenSearchEvent)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
  window.removeEventListener('open-global-search', handleOpenSearchEvent)
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer)
  }
})

// 处理打开搜索事件
const handleOpenSearchEvent = (event) => {
  openSearch()
  if (event.detail && event.detail.query) {
    searchQuery.value = event.detail.query
    handleSearch()
  }
}
</script>

<style scoped>
.global-search {
  display: inline-block;
}

.search-trigger-btn {
  position: fixed;
  bottom: 30px;
  right: 30px;
  z-index: 9999;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.search-dialog :deep(.el-dialog__header) {
  padding: 15px 20px;
  border-bottom: 1px solid #e4e7ed;
}

.search-container {
  min-height: 400px;
  display: flex;
  flex-direction: column;
}

.search-input-wrapper {
  margin-bottom: 15px;
}

.search-filters {
  margin-bottom: 15px;
  display: flex;
  justify-content: center;
}

.search-results {
  flex: 1;
  overflow-y: auto;
  max-height: 500px;
  padding: 10px 0;
}

.search-placeholder,
.no-results {
  padding: 40px 20px;
  text-align: center;
}

.search-tips {
  margin-top: 20px;
  text-align: left;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.search-tips ul {
  margin: 10px 0;
  padding-left: 20px;
  color: #606266;
}

.search-tips li {
  margin: 5px 0;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.result-item {
  padding: 12px 15px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.result-item:hover,
.result-item.selected {
  border-color: #409eff;
  background: #f0f9ff;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.result-score {
  font-size: 12px;
  color: #909399;
}

.result-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
}

.result-meta {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.result-excerpt {
  font-size: 13px;
  color: #606266;
  flex: 1;
}

.result-content {
  font-size: 13px;
  color: #606266;
  flex: 1;
}

.result-author {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.result-post-title {
  font-size: 13px;
  color: #909399;
}

.result-footer {
  margin-top: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-time {
  font-size: 12px;
  color: #909399;
}

.search-footer {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #e4e7ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-footer span {
  color: #606266;
  font-size: 14px;
}

/* 滚动条样式 */
.search-results::-webkit-scrollbar {
  width: 6px;
}

.search-results::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.search-results::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.search-results::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
