<template>
  <div class="menu-management">
    <!-- Header -->
    <div class="page-header">
      <div class="header-left">
        <h1>菜单管理</h1>
        <p class="subtitle">管理网站导航菜单，支持拖拽排序和多级菜单</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="handleAddMenu">
          <el-icon><Plus /></el-icon>
          新建菜单
        </el-button>
      </div>
    </div>

    <!-- Menu Locations Tabs -->
    <el-card class="menu-locations-card">
      <el-tabs v-model="activeLocation" type="border-card">
        <el-tab-pane 
          v-for="location in menuLocations" 
          :key="location.id"
          :label="location.name"
          :name="location.id"
        >
          <template #label>
            <span class="tab-label">
              <el-icon><component :is="getLocationIcon(location.id)" /></el-icon>
              {{ location.name }}
            </span>
          </template>
          
          <div class="location-description">{{ location.description }}</div>
          
          <!-- Menu Items Draggable List -->
          <div class="menu-items-container">
            <div class="menu-toolbar">
              <el-button type="primary" size="small" @click="handleAddItem">
                <el-icon><Plus /></el-icon>
                添加菜单项
              </el-button>
              <el-button size="small" @click="handleExpandAll">
                <el-icon><ArrowDown /></el-icon>
                展开全部
              </el-button>
              <el-button size="small" @click="handleCollapseAll">
                <el-icon><ArrowUp /></el-icon>
                收起全部
              </el-button>
              <el-button type="success" size="small" @click="handleSaveMenu" :loading="saving">
                <el-icon><Check /></el-icon>
                保存菜单
              </el-button>
            </div>

            <draggable
              v-model="menuItems"
              group="menu"
              item-key="id"
              handle=".drag-handle"
              :animation="200"
              ghost-class="ghost"
              drag-class="dragging"
              @end="onDragEnd"
              class="menu-tree"
            >
              <template #item="{ element, index }">
                <div class="menu-item-wrapper" :class="{ 'is-expanded': element.expanded }">
                  <div class="menu-item" :class="{ 'is-nested': element.parentId }" :style="getItemStyle(element)">
                    <div class="drag-handle">
                      <el-icon><Rank /></el-icon>
                    </div>
                    <div class="item-info">
                      <div class="item-main">
                        <el-icon class="item-icon">
                          <component :is="element.icon || 'Link'" />
                        </el-icon>
                        <span class="item-label">{{ element.label }}</span>
                        <el-tag v-if="element.children && element.children.length" size="small" type="info">
                          {{ element.children.length }} 个子项
                        </el-tag>
                      </div>
                      <div class="item-meta">
                        <span class="item-url">{{ element.url }}</span>
                        <el-tag v-if="element.target === '_blank'" size="small" type="warning">新窗口</el-tag>
                      </div>
                    </div>
                    <div class="item-actions">
                      <el-button 
                        v-if="canHaveChildren(element)"
                        link 
                        size="small" 
                        @click="toggleExpand(element)"
                      >
                        <el-icon>
                          <ArrowDown v-if="!element.expanded" />
                          <ArrowUp v-else />
                        </el-icon>
                      </el-button>
                      <el-button link size="small" @click="handleEditItem(element)">
                        <el-icon><Edit /></el-icon>
                      </el-button>
                      <el-button link size="small" type="danger" @click="handleDeleteItem(element, index)">
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                  </div>

                  <!-- Nested Children -->
                  <draggable
                    v-if="element.children && element.children.length && element.expanded"
                    v-model="element.children"
                    group="menu"
                    item-key="id"
                    handle=".drag-handle"
                    :animation="200"
                    ghost-class="ghost"
                    class="menu-children"
                  >
                    <template #item="{ element: child, index: childIndex }">
                      <div class="menu-item-wrapper">
                        <div class="menu-item is-nested">
                          <div class="drag-handle">
                            <el-icon><Rank /></el-icon>
                          </div>
                          <div class="item-info">
                            <div class="item-main">
                              <el-icon class="item-icon">
                                <component :is="child.icon || 'Link'" />
                              </el-icon>
                              <span class="item-label">{{ child.label }}</span>
                            </div>
                            <div class="item-meta">
                              <span class="item-url">{{ child.url }}</span>
                              <el-tag v-if="child.target === '_blank'" size="small" type="warning">新窗口</el-tag>
                            </div>
                          </div>
                          <div class="item-actions">
                            <el-button link size="small" @click="handleEditItem(child, element)">
                              <el-icon><Edit /></el-icon>
                            </el-button>
                            <el-button link size="small" type="danger" @click="handleDeleteChild(element, childIndex)">
                              <el-icon><Delete /></el-icon>
                            </el-button>
                          </div>
                        </div>
                      </div>
                    </template>
                  </draggable>
                </div>
              </template>
            </draggable>

            <!-- Empty State -->
            <el-empty v-if="!menuItems.length" description="暂无菜单项，点击上方按钮添加">
              <el-button type="primary" @click="handleAddItem">添加菜单项</el-button>
            </el-empty>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Add/Edit Menu Dialog -->
    <el-dialog
      v-model="menuDialogVisible"
      :title="isEditing ? '编辑菜单' : '新建菜单'"
      width="500px"
    >
      <el-form :model="menuForm" label-width="100px" :rules="menuRules" ref="menuFormRef">
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="menuForm.name" placeholder="例如：主导航" />
        </el-form-item>
        <el-form-item label="菜单位置" prop="location">
          <el-select v-model="menuForm.location" placeholder="选择菜单位置" style="width: 100%">
            <el-option
              v-for="loc in availableLocations"
              :key="loc.id"
              :label="loc.name"
              :value="loc.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input 
            v-model="menuForm.description" 
            type="textarea" 
            rows="3"
            placeholder="菜单描述（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveMenuConfig" :loading="saving">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- Add/Edit Menu Item Dialog -->
    <el-dialog
      v-model="itemDialogVisible"
      :title="isEditingItem ? '编辑菜单项' : '添加菜单项'"
      width="600px"
    >
      <el-form :model="itemForm" label-width="100px" :rules="itemRules" ref="itemFormRef">
        <el-form-item label="菜单名称" prop="label">
          <el-input v-model="itemForm.label" placeholder="例如：首页" />
        </el-form-item>
        
        <el-form-item label="链接类型" prop="linkType">
          <el-radio-group v-model="itemForm.linkType">
            <el-radio label="custom">自定义链接</el-radio>
            <el-radio label="page">页面</el-radio>
            <el-radio label="category">分类</el-radio>
            <el-radio label="post">文章</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="itemForm.linkType === 'custom'" label="链接地址" prop="url">
          <el-input v-model="itemForm.url" placeholder="例如：/about 或 https://example.com">
            <template #prepend>
              <el-select v-model="urlPrefix" style="width: 100px">
                <el-option label="站内" value="" />
                <el-option label="https://" value="https://" />
                <el-option label="http://" value="http://" />
              </el-select>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item v-if="itemForm.linkType === 'page'" label="选择页面" prop="url">
          <el-select v-model="itemForm.url" placeholder="选择页面" style="width: 100%">
            <el-option
              v-for="page in pages"
              :key="page.id"
              :label="page.title"
              :value="`/page/${page.slug}`"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="itemForm.linkType === 'category'" label="选择分类" prop="url">
          <el-select v-model="itemForm.url" placeholder="选择分类" style="width: 100%">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="`/category/${cat.slug}`"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="itemForm.linkType === 'post'" label="选择文章" prop="url">
          <el-select
            v-model="itemForm.url"
            placeholder="搜索并选择文章"
            filterable
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="post in posts"
              :key="post.id"
              :label="post.title"
              :value="`/article/${post.slug}`"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="图标">
          <IconPicker v-model="itemForm.icon" />
          <div class="form-hint">选择一个图标用于菜单项显示</div>
        </el-form-item>

        <el-form-item label="打开方式">
          <el-radio-group v-model="itemForm.target">
            <el-radio label="_self">当前窗口</el-radio>
            <el-radio label="_blank">新窗口</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="itemForm.enabled" />
          <span class="form-hint" style="margin-left: 12px;">禁用后此菜单项将不在前台显示</span>
        </el-form-item>

        <el-form-item label="父级菜单" v-if="canSelectParent">
          <el-select v-model="itemForm.parentId" placeholder="作为顶级菜单" clearable style="width: 100%">
            <el-option
              v-for="item in parentOptions"
              :key="item.id"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="itemDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveItem" :loading="saving">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import draggable from 'vuedraggable'
import {
  Plus, Edit, Delete, Rank, ArrowDown, ArrowUp, Check,
  Link, Home, Document, Folder, Menu, Cellphone
} from '@element-plus/icons-vue'
import request from '@/utils/request'
import IconPicker from '@/components/IconPicker.vue'

// Types
interface MenuLocation {
  id: string
  name: string
  description: string
  maxDepth: number
}

interface MenuItem {
  id: string
  label: string
  url: string
  icon?: string
  target: string
  parentId?: string
  children?: MenuItem[]
  expanded?: boolean
  sort: number
  enabled?: boolean
}

interface Menu {
  id: string
  name: string
  location: string
  description?: string
  items: MenuItem[]
}

// State
const activeLocation = ref('header')
const menuLocations = ref<MenuLocation[]>([
  { id: 'header', name: '顶部导航', description: '网站顶部主导航菜单', maxDepth: 3 },
  { id: 'footer', name: '底部导航', description: '网站底部导航菜单', maxDepth: 2 },
  { id: 'sidebar', name: '侧边栏菜单', description: '侧边栏导航菜单', maxDepth: 2 },
  { id: 'mobile', name: '移动端菜单', description: '移动端专用菜单', maxDepth: 3 }
])

const menus = ref<Menu[]>([])
const menuItems = ref<MenuItem[]>([])
const pages = ref<any[]>([])
const categories = ref<any[]>([])
const posts = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)

// Dialog states
const menuDialogVisible = ref(false)
const itemDialogVisible = ref(false)
const isEditing = ref(false)
const isEditingItem = ref(false)
const currentMenu = ref<Menu | null>(null)
const currentParent = ref<MenuItem | null>(null)

// Forms
const menuFormRef = ref()
const itemFormRef = ref()

const menuForm = reactive({
  id: '',
  name: '',
  location: '',
  description: ''
})

const itemForm = reactive({
  id: '',
  label: '',
  url: '',
  icon: '',
  target: '_self',
  linkType: 'custom',
  parentId: '',
  enabled: true
})

const urlPrefix = ref('')

// Rules
const menuRules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  location: [{ required: true, message: '请选择菜单位置', trigger: 'change' }]
}

const itemRules = {
  label: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入或选择链接', trigger: 'blur' }]
}

// Computed
const availableLocations = computed(() => {
  const usedLocations = menus.value.map(m => m.location)
  return menuLocations.value.filter(loc => !usedLocations.includes(loc.id) || 
    (isEditing.value && currentMenu.value?.location === loc.id))
})

const postOptions = computed(() => {
  return posts.value.map(post => ({
    value: `/blog/${post.slug}`,
    label: post.title
  }))
})

const canSelectParent = computed(() => {
  const location = menuLocations.value.find(l => l.id === activeLocation.value)
  return location && location.maxDepth > 1
})

const parentOptions = computed(() => {
  return menuItems.value.filter(item => !item.parentId)
})

// Methods
const getLocationIcon = (locationId: string) => {
  const icons: Record<string, any> = {
    header: Menu,
    footer: Document,
    sidebar: Folder,
    mobile: Cellphone
  }
  return icons[locationId] || Link
}

const getItemStyle = (item: MenuItem) => {
  if (!item.parentId) return {}
  return { marginLeft: '32px' }
}

const canHaveChildren = (item: MenuItem) => {
  const location = menuLocations.value.find(l => l.id === activeLocation.value)
  if (!location) return false
  const depth = getItemDepth(item)
  return depth < location.maxDepth - 1
}

const getItemDepth = (item: MenuItem): number => {
  if (!item.parentId) return 0
  const parent = menuItems.value.find(i => i.id === item.parentId)
  return parent ? getItemDepth(parent) + 1 : 0
}

const toggleExpand = (item: MenuItem) => {
  item.expanded = !item.expanded
}

const handleExpandAll = () => {
  menuItems.value.forEach(item => {
    item.expanded = true
  })
}

const handleCollapseAll = () => {
  menuItems.value.forEach(item => {
    item.expanded = false
  })
}

const onDragEnd = () => {
  // Update sort order
  menuItems.value.forEach((item, index) => {
    item.sort = index
  })
}

// API Functions
const fetchMenus = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/menus')
    menus.value = res.data || []
    loadMenuItems()
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

const fetchPages = async () => {
  try {
    const res = await request.get('/api/pages', { params: { status: 'published', page_size: 100 } })
    pages.value = res.data?.list || []
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

const fetchCategories = async () => {
  try {
    const res = await request.get('/api/categories', { params: { page_size: 100 } })
    categories.value = res.data?.list || []
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

const fetchPosts = async () => {
  try {
    const res = await request.get('/api/posts', { params: { status: 'published', page_size: 100 } })
    posts.value = res.data?.list || []
  } catch (error) {
    // 错误已在 request 拦截器中处理
  }
}

const loadMenuItems = () => {
  const menu = menus.value.find(m => m.location === activeLocation.value)
  if (menu) {
    menuItems.value = menu.items || []
  } else {
    menuItems.value = []
  }
}

// Menu Actions
const handleAddMenu = () => {
  isEditing.value = false
  menuForm.id = ''
  menuForm.name = ''
  menuForm.location = ''
  menuForm.description = ''
  menuDialogVisible.value = true
}

const handleSaveMenuConfig = async () => {
  const valid = await menuFormRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    if (isEditing.value) {
      await request.put(`/api/menus/${menuForm.id}`, menuForm)
    } else {
      await request.post('/api/menus', menuForm)
    }
    ElMessage.success('保存成功')
    menuDialogVisible.value = false
    fetchMenus()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// Item Actions
const handleAddItem = () => {
  isEditingItem.value = false
  itemForm.id = ''
  itemForm.label = ''
  itemForm.url = ''
  itemForm.icon = ''
  itemForm.target = '_self'
  itemForm.linkType = 'custom'
  itemForm.parentId = ''
  itemForm.enabled = true
  urlPrefix.value = ''
  itemDialogVisible.value = true
}

const handleEditItem = (item: MenuItem, parent?: MenuItem) => {
  isEditingItem.value = true
  currentParent.value = parent || null
  itemForm.id = item.id
  itemForm.label = item.label
  itemForm.url = item.url
  itemForm.icon = item.icon || ''
  itemForm.target = item.target || '_self'
  itemForm.linkType = 'custom'
  itemForm.parentId = item.parentId || ''
  itemForm.enabled = item.enabled !== undefined ? item.enabled : true

  // Detect link type
  if (item.url.startsWith('/page/')) {
    itemForm.linkType = 'page'
  } else if (item.url.startsWith('/category/')) {
    itemForm.linkType = 'category'
  } else if (item.url.startsWith('/blog/')) {
    itemForm.linkType = 'post'
  } else if (item.url.startsWith('http')) {
    itemForm.linkType = 'custom'
    if (item.url.startsWith('https://')) {
      urlPrefix.value = 'https://'
      itemForm.url = item.url.replace('https://', '')
    } else if (item.url.startsWith('http://')) {
      urlPrefix.value = 'http://'
      itemForm.url = item.url.replace('http://', '')
    }
  }

  itemDialogVisible.value = true
}

const handleSaveItem = async () => {
  const valid = await itemFormRef.value?.validate().catch(() => false)
  if (!valid) return

  // Build full URL
  let fullUrl = itemForm.url
  if (itemForm.linkType === 'custom' && urlPrefix.value) {
    fullUrl = urlPrefix.value + itemForm.url
  }

  const newItem: MenuItem = {
    id: itemForm.id || generateId(),
    label: itemForm.label,
    url: fullUrl,
    icon: itemForm.icon,
    target: itemForm.target,
    parentId: itemForm.parentId || undefined,
    sort: 0,
    expanded: false,
    enabled: itemForm.enabled
  }

  if (isEditingItem.value) {
    // Update existing
    if (currentParent.value) {
      const index = currentParent.value.children?.findIndex(c => c.id === newItem.id)
      if (index !== undefined && index >= 0) {
        currentParent.value.children![index] = newItem
      }
    } else {
      const index = menuItems.value.findIndex(i => i.id === newItem.id)
      if (index >= 0) {
        menuItems.value[index] = newItem
      }
    }
  } else {
    // Add new
    if (itemForm.parentId) {
      const parent = menuItems.value.find(i => i.id === itemForm.parentId)
      if (parent) {
        if (!parent.children) parent.children = []
        parent.children.push(newItem)
        parent.expanded = true
      }
    } else {
      menuItems.value.push(newItem)
    }
  }

  itemDialogVisible.value = false
  
  // 直接保存到服务器
  await saveMenuToServer()
}

// 保存菜单到服务器
const saveMenuToServer = async () => {
  saving.value = true
  try {
    const existingMenu = menus.value.find(m => m.location === activeLocation.value)
    if (existingMenu) {
      await request.put(`/api/menus/${existingMenu.id}`, { ...existingMenu, items: menuItems.value })
    } else {
      await request.post('/api/menus', {
        name: menuLocations.value.find(l => l.id === activeLocation.value)?.name,
        location: activeLocation.value,
        items: menuItems.value
      })
    }
    
    ElMessage.success('保存成功')
    fetchMenus()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleDeleteItem = async (item: MenuItem, index: number) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜单项 "${item.label}" 吗？`,
      '确认删除',
      { type: 'warning' }
    )
    menuItems.value.splice(index, 1)
    
    // 直接保存到服务器
    await saveMenuToServer()
  } catch {
    // Cancelled
  }
}

const handleDeleteChild = async (parent: MenuItem, index: number) => {
  parent.children?.splice(index, 1)
  
  // 直接保存到服务器
  await saveMenuToServer()
}

const handleSaveMenu = async () => {
  saving.value = true
  try {
    const menuData = {
      location: activeLocation.value,
      items: menuItems.value
    }
    
    // Check if menu exists
    const existingMenu = menus.value.find(m => m.location === activeLocation.value)
    if (existingMenu) {
      await request.put(`/api/menus/${existingMenu.id}`, { ...existingMenu, items: menuItems.value })
    } else {
      await request.post('/api/menus', {
        name: menuLocations.value.find(l => l.id === activeLocation.value)?.name,
        location: activeLocation.value,
        items: menuItems.value
      })
    }
    
    ElMessage.success('菜单保存成功')
    fetchMenus()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const generateId = () => {
  return 'menu_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
}

// Watch
watch(activeLocation, () => {
  loadMenuItems()
})

// Lifecycle
onMounted(() => {
  fetchMenus()
  fetchPages()
  fetchCategories()
  fetchPosts()
})
</script>

<style scoped lang="scss">
.menu-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;

  h1 {
    margin: 0 0 8px;
    font-size: 24px;
    font-weight: 600;
  }

  .subtitle {
    margin: 0;
    color: #909399;
    font-size: 14px;
  }
}

.menu-locations-card {
  .tab-label {
    display: flex;
    align-items: center;
    gap: 6px;
  }
}

.location-description {
  color: #909399;
  font-size: 13px;
  margin-bottom: 16px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.menu-items-container {
  min-height: 300px;
}

.menu-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.menu-tree {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.menu-item-wrapper {
  .menu-item {
    display: flex;
    align-items: center;
    padding: 12px 16px;
    background: #fff;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    transition: all 0.2s;

    &:hover {
      border-color: #409eff;
      box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
    }

    &.is-nested {
      background: #f5f7fa;
    }

    .drag-handle {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      margin-right: 12px;
      color: #909399;
      cursor: grab;
      border-radius: 4px;
      transition: all 0.2s;

      &:hover {
        background: #ebeef5;
        color: #606266;
      }

      &:active {
        cursor: grabbing;
      }
    }

    .item-info {
      flex: 1;
      min-width: 0;

      .item-main {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 4px;

        .item-icon {
          color: #409eff;
          font-size: 16px;
        }

        .item-label {
          font-weight: 500;
          color: #303133;
        }
      }

      .item-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 12px;
        color: #909399;

        .item-url {
          max-width: 200px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }

    .item-actions {
      display: flex;
      gap: 4px;
      opacity: 0;
      transition: opacity 0.2s;
    }

    &:hover .item-actions {
      opacity: 1;
    }
  }
}

.menu-children {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
  margin-left: 32px;
}

.ghost {
  opacity: 0.5;
  background: #f5f7fa;
  border: 2px dashed #409eff;
}

.dragging {
  opacity: 0.8;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.form-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
