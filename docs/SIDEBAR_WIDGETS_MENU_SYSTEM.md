# 侧边栏小工具与菜单管理系统

## 概述

本次更新为 Catter CMS 添加了完整的侧边栏小工具系统和独立的菜单管理界面，参考 Halo CMS 的设计精髓，提供更灵活的主题配置能力。

## 功能清单

### 1. 单页面模板 (page.plush.html)
- 支持侧边栏开关和位置配置（左/右）
- 完整的页面元信息展示（发布日期、更新日期、作者）
- 响应式设计，移动端自动隐藏侧边栏
- 支持自定义CSS样式注入

### 2. 侧边栏小工具系统 (sidebar.plush.html)
支持9种小工具类型：
1. **搜索框** - 带搜索按钮和占位符配置
2. **博主信息** - 头像、名称、简介、社交链接
3. **最近文章** - 可配置数量和缩略图显示
4. **热门文章** - 按浏览量排序，带排名标识
5. **分类列表** - 显示文章数量，支持层级
6. **标签云** - 动态字号，可配置数量
7. **站点统计** - 文章/分类/标签/评论数量
8. **公告** - 富文本内容
9. **自定义HTML** - 完全自定义内容

### 3. 主题配置扩展 (theme.yml)
新增配置分组：
- 侧边栏基础设置（启用、位置、宽度、固定）
- 小工具排序（widget-sort类型，支持拖拽排序和启用/禁用）
- 各小工具独立配置分组（9个分组，共30+配置项）
- 条件显示逻辑（show_if）

### 4. 菜单管理界面 (MenuList.vue)
功能特性：
- 4个菜单位置：顶部导航、底部导航、侧边栏菜单、移动端菜单
- 拖拽排序（vuedraggable）
- 多级菜单支持（最多3级）
- 多种链接类型：自定义链接、页面、分类、文章
- 图标选择（Element Plus图标）
- 打开方式设置（当前窗口/新窗口）

### 5. 菜单管理后端API (menus.go)
API端点：
- GET /api/menus - 获取所有菜单
- POST /api/menus - 创建菜单
- GET /api/menus/:id - 获取单个菜单
- PUT /api/menus/:id - 更新菜单
- DELETE /api/menus/:id - 删除菜单
- GET /api/menus/location/:location - 根据位置获取菜单
- GET /api/menus/:id/items - 获取菜单项
- POST /api/menus/:id/items - 创建菜单项
- PUT /api/menus/:id/items/:itemId - 更新菜单项
- DELETE /api/menus/:id/items/:itemId - 删除菜单项
- POST /api/menus/:id/items/reorder - 重新排序

### 6. 数据库模型 (menu.go)
- Menu模型：id, name, location, description
- MenuItem模型：id, menu_id, label, url, icon, target, parent_id, sort, enabled
- 支持层级关系（parent_id关联）
- 自动迁移支持

## 文件变更

### 新增文件
1. `templates/themes/catter-blog/pages/page.plush.html` - 单页面模板
2. `templates/themes/catter-blog/pages/sidebar.plush.html` - 侧边栏小工具模板
3. `admin/src/views/menus/MenuList.vue` - 菜单管理界面
4. `gin-app/internal/handlers/api/menus.go` - 菜单API处理器
5. `gin-app/internal/models/menu.go` - 菜单模型

### 修改文件
1. `templates/themes/catter-blog/theme.yml` - 添加侧边栏小工具配置
2. `admin/src/router/index.js` - 添加菜单管理路由
3. `admin/src/components/layout/Sidebar.vue` - 添加菜单管理入口
4. `gin-app/cmd/main.go` - 注册菜单API路由
5. `gin-app/internal/models/db.go` - 添加菜单表自动迁移

## 使用方法

### 配置侧边栏小工具
1. 进入后台 -> 主题配置
2. 选择"侧边栏"分组
3. 启用/禁用侧边栏，设置位置和宽度
4. 在"小工具排序"中拖拽调整顺序，取消勾选可禁用
5. 展开各小工具分组进行详细配置

### 管理菜单
1. 进入后台 -> 菜单管理
2. 选择顶部导航/底部导航/侧边栏菜单/移动端菜单标签
3. 点击"添加菜单项"创建新菜单
4. 拖拽调整顺序
5. 点击保存按钮保存更改

### 在模板中使用
```html
<!-- 在页面模板中包含侧边栏 -->
<%= partial("sidebar.html") %>

<!-- 获取菜单数据 -->
<%= let menu = getMenu("header") %>
```

## 技术亮点

1. **widget-sort配置类型** - 创新的配置字段类型，支持拖拽排序和启用/禁用
2. **条件显示** - 使用show_if实现配置项的动态显示/隐藏
3. **层级菜单** - 支持无限层级（受maxDepth限制）的嵌套菜单
4. **事务处理** - 菜单保存使用数据库事务，确保数据一致性
5. **响应式设计** - 所有模板都支持移动端适配

## 与Halo的对比

| 功能 | Catter CMS | Halo CMS |
|------|-----------|----------|
| 小工具排序 | 拖拽排序 | 拖拽排序 |
| 小工具类型 | 9种 | 10+种 |
| 多级菜单 | 支持3级 | 支持多级 |
| 菜单位置 | 4个 | 多个 |
| 条件显示 | 支持 | 支持 |
| 配置分组 | 支持 | 支持 |

## 后续优化建议

1. 添加更多小工具类型（最新评论、文章归档、友情链接等）
2. 支持小工具自定义CSS类
3. 添加菜单导入/导出功能
4. 支持菜单权限控制（根据用户角色显示不同菜单）
5. 添加小工具缓存机制，提升性能
