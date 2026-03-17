---
name: Catter CMS 管理后台优化计划 v2
overview: 轻量化、易部署、SEO友好的CMS管理后台，支持灵活的主题系统，前台SSR渲染，后台Vue SPA嵌入部署
design:
  architecture:
    framework: vue
  styleKeywords:
    - Modern Minimalism
    - Clean
    - Light Theme
    - Card Layout
    - Responsive
  fontSystem:
    fontFamily: Inter
    heading:
      size: 24px
      weight: 600
    subheading:
      size: 16px
      weight: 500
    body:
      size: 14px
      weight: 400
  colorSystem:
    primary:
      - "#3B82F6"
      - "#2563EB"
      - "#1D4ED8"
    background:
      - "#FFFFFF"
      - "#F8FAFC"
      - "#F1F5F9"
    text:
      - "#1E293B"
      - "#475569"
      - "#94A3B8"
    functional:
      - "#22C55E"
      - "#F59E0B"
      - "#EF4444"
      - "#06B6D4"
todos:
  - id: phase1-setup
    content: 初始化前端项目配置：更新vite.config.js构建配置、创建样式目录、配置Element Plus按需加载
    status: completed
  - id: phase1-backend-api
    content: 完善后端API：创建统一响应格式、实现仪表盘统计API、完善用户管理API
    status: completed
  - id: phase2-layout
    content: 实现布局组件：AppLayout主布局、Sidebar可折叠侧边栏、Header头部组件
    status: completed
    dependencies:
      - phase1-setup
  - id: phase2-dashboard
    content: 完善仪表盘模块：统计卡片、最近文章列表、API数据对接
    status: completed
    dependencies:
      - phase2-layout
      - phase1-backend-api
  - id: phase2-posts
    content: 完善文章管理：文章列表页（搜索筛选分页）、文章编辑页（Markdown编辑器）
    status: completed
    dependencies:
      - phase2-layout
  - id: phase2-categories-tags
    content: 实现分类和标签管理：分类树形列表、标签列表、CRUD操作
    status: completed
    dependencies:
      - phase2-layout
  - id: phase3-users-roles
    content: 实现用户和角色管理：用户列表、用户编辑、角色管理、权限配置
    status: completed
    dependencies:
      - phase2-layout
  - id: phase3-media
    content: 实现媒体库模块：网格视图、文件上传、图片预览
    status: completed
    dependencies:
      - phase2-layout
  - id: phase3-themes-settings
    content: 实现主题配置和系统设置：主题可视化配置表单、站点/SEO/邮件设置页
    status: completed
    dependencies:
      - phase2-layout
  - id: phase4-integration
    content: 集成测试和构建优化：Vite生产构建、复制到public目录、测试单二进制部署
    status: completed
    dependencies:
      - phase2-dashboard
      - phase2-posts
      - phase2-categories-tags
      - phase3-users-roles
      - phase3-media
      - phase3-themes-settings
---

## 产品概述

Catter CMS管理后台优化项目，核心目标是构建一个**轻量化、易部署、SEO友好、主题灵活**的内容管理系统。

### 核心设计目标

1. **轻量化**：二进制体积小（目标30-50MB）、内存占用低（目标50-100MB运行时）、启动快
2. **易部署**：单二进制文件、零配置启动、无外部依赖
3. **SEO友好**：前台服务端渲染（SSR）、搜索引擎可抓取
4. **主题灵活**：

- 易于切换激活（后台一键操作）
- 方便开发新主题（标准化模板结构）
- 后台可视化配置主题模块（颜色、布局、功能开关）

5. **管理后台完整**：现代化UI、功能完整、非残次品

### 核心功能

#### 前台功能（SSR，保持不变）

- 服务端渲染：使用Buffalo Plush模板引擎
- 主题系统：多主题支持，数据库配置切换
- SEO优化：完整的meta标签、sitemap、结构化数据
- 静态生成：支持生成静态HTML文件

#### 管理后台功能

- **仪表盘**：统计卡片、最近文章、快捷操作
- **文章管理**：列表/搜索/筛选、Markdown编辑、批量操作
- **分类/标签管理**：树形结构、CRUD操作
- **用户/角色管理**：用户列表、角色分配、权限配置
- **媒体库**：文件上传、图片预览、媒体管理
- **主题配置**：主题列表、一键激活、可视化配置表单
- **系统设置**：站点信息、SEO设置、邮件配置

## 技术栈

### 后端

- **框架**：Go 1.25 + Buffalo 1.1.3
- **数据库**：SQLite（轻量化首选）
- **认证**：JWT
- **API**：RESTful设计

### 前端

- **框架**：Vue 3.4 + TypeScript
- **构建工具**：Vite 5
- **UI组件库**：Element Plus 2.4（按需加载）
- **状态管理**：Pinia 2.1
- **路由**：Vue Router 4.2

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    Buffalo 服务（单二进制）                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  前台SSR    │  │  静态生成   │  │    管理后台SPA      │  │
│  │  (Plush)    │  │  (命令)     │  │  (Vue 3嵌入)        │  │
│  │  SEO友好    │  │  可选       │  │  现代化UI           │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                      API层 (RESTful)                         │
├─────────────────────────────────────────────────────────────┤
│                     业务逻辑层 (Services)                     │
├─────────────────────────────────────────────────────────────┤
│                     数据模型层 (Models)                       │
└─────────────────────────────────────────────────────────────┘
```

## 主题系统设计

### 主题目录结构

```
templates/themes/{theme-slug}/
├── theme.yaml          # 主题配置定义
├── pages/              # 页面模板
│   ├── home.plush.html
│   ├── blog_show.plush.html
│   └── page.plush.html
└── partials/           # 公共组件
```

### theme.yaml 配置示例

```
name: Catter Blog
slug: catter-blog
version: 1.0.0
config:
  - name: primary_color
    label: 主色调
    type: color
    default: "#3B82F6"
  - name: layout
    label: 布局方式
    type: select
    options:
      - value: sidebar
        label: 侧边栏布局
      - value: full
        label: 全宽布局
```

## 轻量化策略

- Element Plus按需加载（减少约60%体积）
- Vite代码分割（vendor/app分离）
- Tree Shaking移除未使用代码
- Gzip压缩（减少约70%传输体积）
- 预期二进制大小：30-50MB

## 目录结构

```
catter-cms/
├── actions/
│   ├── admin/           # [MODIFY] 管理后台API
│   ├── api/             # [KEEP] 前台API
│   └── front/           # [KEEP] 前台SSR
├── admin/               # [MODIFY] Vue前端项目
│   └── src/
│       ├── components/  # [NEW] 布局/通用组件
│       ├── views/       # [EXPAND] 页面组件
│       ├── store/       # [EXPAND] Pinia状态
│       └── utils/       # [NEW] 工具函数
├── templates/themes/    # [KEEP] 前台主题模板
├── models/              # [KEEP] 数据模型
├── services/            # [KEEP] 业务服务
└── public/              # Vue构建输出
```

## 设计风格

采用**现代极简风格**，融合Vercel和Notion的设计语言，打造专业、清爽、高效的管理后台界面。

### 设计原则

1. **留白优先**：充足的间距，内容呼吸感
2. **层次分明**：通过阴影和背景色区分层级
3. **交互反馈**：即时响应的状态变化
4. **功能导向**：减少装饰，聚焦任务

### 布局结构

```
┌──────────────────────────────────────────────────────────┐
│  Sidebar (240px/64px折叠)  │      Main Content          │
│  ┌──────────────────────┐  │  ┌──────────────────────┐  │
│  │ Catter CMS           │  │  │ Header (用户信息)     │  │
│  ├──────────────────────┤  │  ├──────────────────────┤  │
│  │ 仪表盘               │  │  │                      │  │
│  │ 文章管理             │  │  │                      │  │
│  │ 分类/标签            │  │  │     Page Content     │  │
│  │ 用户管理             │  │  │                      │  │
│  │ 媒体库               │  │  │                      │  │
│  │ 主题配置             │  │  │                      │  │
│  │ 系统设置             │  │  │                      │  │
│  └──────────────────────┘  │  └──────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

### 页面规划

| 页面 | 核心功能 | 特殊组件 |
| --- | --- | --- |
| 登录页 | 登录表单 | 品牌展示 |
| 仪表盘 | 统计卡片、最近文章 | 数据图表 |
| 文章列表 | 搜索、筛选、分页 | 数据表格 |
| 文章编辑 | Markdown编辑 | 富文本编辑器 |
| 主题管理 | 主题列表、激活、配置 | 动态配置表单 |
| 系统设置 | 站点/SEO/邮件配置 | Tab分组表单 |


## Agent Extensions

### SubAgent

- **code-explorer**: 用于在开发过程中探索代码库，查找现有实现模式、依赖关系和代码结构，确保新代码与项目风格一致。