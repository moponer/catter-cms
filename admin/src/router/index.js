import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Dashboard from '@/views/Dashboard.vue'
import Posts from '@/views/posts/PostList.vue'
import PostEdit from '@/views/posts/PostEdit.vue'
import Categories from '@/views/categories/CategoryList.vue'
import Tags from '@/views/tags/TagList.vue'
import Users from '@/views/users/UserList.vue'
import UserEdit from '@/views/users/UserEdit.vue'
import Media from '@/views/media/MediaLibrary.vue'
import Pages from '@/views/pages/PageList.vue'
import PageEdit from '@/views/pages/PageEdit.vue'
import Themes from '@/views/themes/ThemeSettings.vue'
import Settings from '@/views/settings/SiteSettings.vue'
import Comments from '@/views/comments/Index.vue'
import ScheduledTasks from '@/views/scheduled-tasks/Index.vue'
import Sitemap from '@/views/sitemap/Index.vue'
import Menus from '@/views/menus/MenuList.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登录' }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: { title: '仪表盘', requiresAuth: true }
  },
  {
    path: '/posts',
    name: 'Posts',
    component: Posts,
    meta: { title: '文章列表', requiresAuth: true }
  },
  {
    path: '/posts/create',
    name: 'PostCreate',
    component: PostEdit,
    meta: { title: '新建文章', requiresAuth: true }
  },
  {
    path: '/posts/:id/edit',
    name: 'PostEdit',
    component: PostEdit,
    meta: { title: '编辑文章', requiresAuth: true }
  },
  {
    path: '/categories',
    name: 'Categories',
    component: Categories,
    meta: { title: '分类管理', requiresAuth: true }
  },
  {
    path: '/tags',
    name: 'Tags',
    component: Tags,
    meta: { title: '标签管理', requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: Users,
    meta: { title: '用户管理', requiresAuth: true }
  },
  {
    path: '/users/create',
    name: 'UserCreate',
    component: UserEdit,
    meta: { title: '新建用户', requiresAuth: true }
  },
  {
    path: '/users/:id/edit',
    name: 'UserEdit',
    component: UserEdit,
    meta: { title: '编辑用户', requiresAuth: true }
  },
  {
    path: '/media',
    name: 'Media',
    component: Media,
    meta: { title: '媒体库', requiresAuth: true }
  },
  {
    path: '/pages',
    name: 'Pages',
    component: Pages,
    meta: { title: '页面管理', requiresAuth: true }
  },
  {
    path: '/pages/new',
    name: 'PageCreate',
    component: PageEdit,
    meta: { title: '新建页面', requiresAuth: true }
  },
  {
    path: '/pages/:id/edit',
    name: 'PageEdit',
    component: PageEdit,
    meta: { title: '编辑页面', requiresAuth: true }
  },
  {
    path: '/themes',
    name: 'Themes',
    component: Themes,
    meta: { title: '主题配置', requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: { title: '系统设置', requiresAuth: true }
  },
  {
    path: '/comments',
    name: 'Comments',
    component: Comments,
    meta: { title: '评论管理', requiresAuth: true }
  },
  {
    path: '/scheduled-tasks',
    name: 'ScheduledTasks',
    component: ScheduledTasks,
    meta: { title: '定时任务', requiresAuth: true }
  },
  {
    path: '/sitemap',
    name: 'Sitemap',
    component: Sitemap,
    meta: { title: 'Sitemap管理', requiresAuth: true }
  },
  {
    path: '/menus',
    name: 'Menus',
    component: Menus,
    meta: { title: '菜单管理', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  const title = to.meta.title
  document.title = title ? `${title} - Catter CMS` : 'Catter CMS'
  
  // 检查登录状态
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const token = localStorage.getItem('token')
  
  if (requiresAuth && !token) {
    // 需要认证但没有token，跳转到登录页
    next('/login')
  } else if (to.path === '/login' && token) {
    // 已登录访问登录页，跳转到首页
    next('/')
  } else if (requiresAuth && token) {
    // 需要认证且有token，可以继续访问
    // 注意: 这里我们不在路由守卫中验证token有效性，
    // 因为会在后续的API请求中自动验证(通过request拦截器)
    // 如果token无效，API会返回401，拦截器会自动跳转到登录页
    next()
  } else {
    // 其他情况，允许访问
    next()
  }
})

export default router