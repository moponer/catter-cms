---
name: Catter CMS 继续开发计划
overview: 修复编译问题、完善后端API、优化前端组件，确保项目可以正常运行
todos:
  - id: fix-backend-renderer
    content: 修复后端渲染器引用问题：更新 actions/api/response.go 和 actions/admin/ 中的渲染器调用方式
    status: completed
  - id: add-markdown-editor
    content: 集成Markdown编辑器到文章编辑页：安装 bytemd、创建编辑器组件、集成到 PostEdit.vue
    status: completed
    dependencies:
      - fix-backend-renderer
  - id: complete-category-page
    content: 完善分类管理页面：添加新建/编辑对话框、实现CRUD操作、支持树形结构
    status: completed
    dependencies:
      - fix-backend-renderer
  - id: complete-user-pages
    content: 完善用户管理页面：添加用户编辑页、实现密码重置、完善CRUD操作
    status: completed
    dependencies:
      - fix-backend-renderer
  - id: complete-media-library
    content: 完善媒体库页面：实现文件上传、图片预览、批量删除
    status: completed
    dependencies:
      - fix-backend-renderer
  - id: complete-theme-settings
    content: 完善主题配置页面：读取主题配置定义、渲染动态配置表单
    status: completed
    dependencies:
      - fix-backend-renderer
  - id: complete-site-settings
    content: 完善系统设置页面：站点设置、SEO设置、邮件配置表单
    status: completed
    dependencies:
      - fix-backend-renderer
  - id: integration-test
    content: 集成测试：运行前后端、测试API对接、验证构建流程
    status: completed
    dependencies:
      - add-markdown-editor
      - complete-category-page
      - complete-user-pages
      - complete-media-library
      - complete-theme-settings
      - complete-site-settings
---

## 产品概述

继续完善 Catter CMS 管理后台，修复已发现的问题并完善未完成的功能模块。

## 核心问题

1. **后端编译错误**：`actions/api/response.go` 和 `actions/admin/` 中的文件使用了 `r.JSON`/`r.HTML`，但 `r` 变量定义在 `actions` 包中，子包无法直接访问
2. **前端页面不完整**：分类管理、用户管理、媒体库、主题配置等页面功能简陋，缺少完整的CRUD操作
3. **文章编辑页缺少Markdown编辑器**：当前只是简单的textarea

## 需要完成的任务

### 后端修复

- 修复 `actions/api/response.go` 的渲染器引用问题
- 修复 `actions/admin/` 包中所有文件的渲染器引用问题

### 前端完善

- 完善文章编辑页（集成Markdown编辑器）
- 完善分类管理页面（完整的CRUD、树形结构）
- 完善用户管理页面（完整的CRUD、密码重置）
- 完善媒体库页面（上传、预览、删除）
- 完善主题配置页面（配置表单）
- 完善系统设置页面

### 集成测试

- 确保前后端可以正常对接
- 测试完整的构建流程

## 技术栈

- **后端**：Go 1.25 + Buffalo 1.1.3 + SQLite
- **前端**：Vue 3 + Element Plus + Vite

## 关键技术问题修复

### 渲染器引用问题

问题根源：`r` 变量在 `actions/render.go` 的 `actions` 包中定义，但 `api` 和 `admin` 子包无法直接访问。

解决方案：

1. 在 `actions/api/response.go` 中直接使用 `c.Render()` 和 `render.JSON()` 
2. 在 `actions/admin/` 中，API 函数使用 `api.Success()` 等方法，HTML渲染函数需要导入 `actions` 包的渲染器

### 前端Markdown编辑器

推荐使用 **ByteMD** 或 **@kangc/v-md-editor**，它们轻量且功能完善。

## 目录结构

```
F:/Catter-CMS/
├── actions/
│   ├── api/
│   │   └── response.go       # [MODIFY] 修复渲染器引用
│   ├── admin/
│   │   ├── dashboard.go      # [MODIFY] 修复渲染器引用
│   │   ├── users.go          # [MODIFY] 修复渲染器引用
│   │   └── ...其他文件
│   └── render.go             # [KEEP] 渲染器定义
├── admin/src/
│   ├── views/
│   │   ├── posts/
│   │   │   └── PostEdit.vue  # [MODIFY] 添加Markdown编辑器
│   │   ├── categories/
│   │   │   └── CategoryList.vue  # [MODIFY] 完善CRUD
│   │   ├── users/
│   │   │   ├── UserList.vue  # [MODIFY] 完善功能
│   │   │   └── UserEdit.vue  # [MODIFY] 完善编辑页
│   │   ├── media/
│   │   │   └── MediaLibrary.vue  # [MODIFY] 完善上传功能
│   │   ├── themes/
│   │   │   └── ThemeSettings.vue # [MODIFY] 添加配置表单
│   │   └── settings/
│   │       └── SiteSettings.vue  # [MODIFY] 完善设置表单
│   └── utils/
│       └── request.js        # [KEEP] API请求封装
```

## SubAgent

- **code-explorer**: 用于探索代码库中的渲染器使用模式和现有的API实现方式，确保修复方案与项目风格一致