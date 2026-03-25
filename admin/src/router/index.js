import { createRouter, createWebHistory } from 'vue-router'

// 懒加载路由组件
const Login = () => import('@/views/Login.vue')
const Dashboard = () => import('@/views/Dashboard.vue')
const PostList = () => import('@/views/posts/PostList.vue')
const PostEdit = () => import('@/views/posts/PostEdit.vue')
const CategoryList = () => import('@/views/categories/CategoryList.vue')
const TagList = () => import('@/views/tags/TagList.vue')
const UserList = () => import('@/views/users/UserList.vue')
const UserEdit = () => import('@/views/users/UserEdit.vue')
const MediaLibrary = () => import('@/views/media/MediaLibrary.vue')
const PageList = () => import('@/views/pages/PageList.vue')
const PageEdit = () => import('@/views/pages/PageEdit.vue')
const ThemeList = () => import('@/views/themes/ThemeList.vue')
const ThemeSettings = () => import('@/views/themes/ThemeSettings.vue')
const SiteSettings = () => import('@/views/settings/SiteSettings.vue')
const Comments = () => import('@/views/comments/Index.vue')
const ScheduledTasks = () => import('@/views/scheduled-tasks/Index.vue')
const Sitemap = () => import('@/views/sitemap/Index.vue')
const MenuList = () => import('@/views/menus/MenuList.vue')
const Cache = () => import('@/views/cache/Index.vue')

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
    component: PostList,
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
    component: CategoryList,
    meta: { title: '分类管理', requiresAuth: true }
  },
  {
    path: '/tags',
    name: 'Tags',
    component: TagList,
    meta: { title: '标签管理', requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: UserList,
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
    component: MediaLibrary,
    meta: { title: '媒体库', requiresAuth: true }
  },
  {
    path: '/pages',
    name: 'Pages',
    component: PageList,
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
    component: ThemeList,
    meta: { title: '主题管理', requiresAuth: true }
  },
  {
    path: '/themes/:slug',
    name: 'ThemeSettings',
    component: ThemeSettings,
    meta: { title: '主题配置', requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SiteSettings,
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
    path: '/cache',
    name: 'Cache',
    component: Cache,
    meta: { title: '清理缓存', requiresAuth: true }
  },
  {
    path: '/menus',
    name: 'Menus',
    component: MenuList,
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
  } else if (to.path === '/login' || to.path === '/admin/login') {
    // 访问登录页
    if (token) {
      // 有token，验证是否有效
      try {
        const response = await fetch('/api/auth/me', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        if (response.ok) {
          // token有效，跳转到首页
          next('/')
          return
        }
      } catch (e) {
        // 网络错误，继续显示登录页
      }
    }
    // 无token或token无效，允许访问登录页
    next()
  } else if (requiresAuth && token) {
    // 需要认证且有token，可以继续访问
    // token有效性会在后续的API请求中自动验证
    next()
  } else {
    // 其他情况，允许访问
    next()
  }
})

export default router