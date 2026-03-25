# Catter CMS 主题开发文档 v3.0

> 本文档面向零基础开发者，通过手把手教程，带您从零开发一套完整的主题。

**更新日期：** 2026-03-25
**版本：** 3.0
**适配系统：** Catter CMS 2.0+

---

## 目录

### 第一部分：入门指南
1. [快速开始](#1-快速开始) - 5分钟创建第一个主题
2. [主题目录结构](#2-主题目录结构) - 认识主题的每个文件
3. [主题配置文件](#3-主题配置文件) - theme.yml 完整指南

### 第二部分：核心原则（重要！）
4. [CSS 样式规范](#4-css-样式规范) - ✅ 不得使用内联样式
5. [CSS 变量系统](#5-css-变量系统) - 主题配置的核心
6. [Go 模板限制](#6-go-模板限制) - ⚠️ 不得使用跨文件模板包含

### 第三部分：模板开发
7. [模板基础](#7-模板基础) - HTML模板语法入门
8. [模板函数](#8-模板函数) - 数据处理神器
9. [首页模板开发](#9-首页模板开发) - 实战演练
10. [文章详情页](#10-文章详情页) - 展示完整文章内容
11. [列表页模板](#11-列表页模板) - 分类/标签/归档页

### 第四部分：高级配置
12. [主题配置定义](#12-主题配置定义) - 让用户自定义主题
13. [Repeater字段](#13-repeater字段) - 数组配置的可视化编辑
14. [条件显示](#14-条件显示) - show_if语法

### 第五部分：实战案例
15. [案例：博客主题](#15-案例博客主题) - 从设计到实现
16. [案例：企业官网](#16-案例企业官网) - 完整的模块化网站

### 第六部分：参考手册
17. [API参考](#17-api参考) - 管理后台接口
18. [数据模型](#18-数据模型) - Post/Category/Tag等
19. [常见问题](#19-常见问题) - 疑难解答

---

# 第一部分：入门指南

## 1. 快速开始

### 1.1 创建主题的最低要求

一个可运行的主题只需两个文件：

```
templates/themes/my-first-theme/
├── theme.yml      # 主题元数据（必需）
└── pages/
    └── home.html  # 首页模板（必需）
```

### 1.2 最小化 theme.yml

```yaml
# theme.yml
name: 我的第一个主题
slug: my-first-theme
version: 1.0.0
author: 你的名字
description: 这是我的第一个Catter CMS主题
```

### 1.3 最小化 home.html

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .site.title }}</title>
</head>
<body>
    <h1>欢迎访问 {{ .site.title }}</h1>
    <p>这是我的第一个主题！</p>
</body>
</html>
```

### 1.4 激活主题

1. 将主题文件夹放入 `templates/themes/`
2. 访问管理后台 → 主题管理
3. 找到你的主题，点击"激活"
4. 访问网站首页查看效果

**恭喜！您已成功创建第一个主题。**

---

## 2. 主题目录结构

### 2.1 完整目录结构

```
templates/themes/<主题slug>/
│
├── theme.yml              # 【必需】主题元数据配置文件
│
├── assets/                # 【推荐】静态资源目录
│   ├── css/               # 样式文件
│   │   ├── base.css       # 基础样式（reset、布局、通用组件）
│   │   ├── home.css       # 首页特定样式
│   │   ├── article.css    # 文章页特定样式
│   │   └── list.css       # 列表页特定样式
│   ├── js/                # JavaScript 文件
│   └── images/            # 图片资源
│
└── pages/                # 【必需】模板页面目录
    ├── home.html          # 【必需】首页
    ├── article_show.html  # 文章详情页
    ├── article_list.html  # 文章列表页
    ├── category.html      # 分类页
    ├── tag.html           # 标签页
    ├── search.html        # 搜索页
    ├── archives.html      # 归档页
    ├── categories.html    # 分类列表页
    ├── tags.html          # 标签云页
    ├── page.html          # 独立页面
    └── 404.html            # 404错误页
```

### 2.2 目录命名规范

| 规则 | 正确示例 | 错误示例 | 原因 |
|------|----------|----------|------|
| 使用英文 | `my-blog-theme` | `我的博客主题` | 避免编码问题 |
| 使用kebab-case | `catter-blog` | `catterBlog` | 保持一致性 |
| 与slug一致 | slug: `my-theme` | 文件夹: `different-name` | 资源路径依赖 |

### 2.3 实际主题参考

**catter-cms-official 主题结构：**
```
templates/themes/catter-cms-official/
├── assets/
│   ├── css/
│   │   ├── base.css       # 共享基础样式（包含CSS变量默认值）
│   │   ├── home.css       # 首页模块样式
│   │   ├── article.css    # 文章/Markdown样式
│   │   └── list.css       # 列表页样式
│   └── images/
├── pages/
│   ├── home.html          # 首页（Hero+多模块）
│   ├── article_show.html  # 文章详情
│   ├── article_list.html  # 文章列表
│   ├── category.html      # 分类页
│   ├── tag.html           # 标签页
│   ├── archives.html      # 归档页
│   ├── categories.html    # 全部分类
│   ├── tags.html          # 标签云
│   ├── search.html        # 搜索页
│   └── page.html          # 独立页面
└── theme.yml
```

**dianyan-blog 主题结构：**
```
templates/themes/dianyan-blog/
├── assets/
│   └── css/
│       └── base.css       # 共享样式（含响应式）
├── pages/                 # 10个页面模板
│   ├── home.html
│   ├── article_show.html
│   ├── article_list.html
│   ├── category.html
│   ├── tag.html
│   ├── tags.html
│   ├── archives.html
│   ├── categories.html
│   ├── search.html
│   └── page.html
└── theme.yml
```

---

## 3. 主题配置文件

### 3.1 theme.yml 完整结构

```yaml
# ============================================
# 基础信息（必需）
# ============================================
name: 主题显示名称
slug: theme-slug                    # 必须与文件夹名一致
version: 1.0.0
author: 作者名称
description: 主题描述

# ============================================
# config - 默认配置值
# 这些值可以被用户在后台修改
# 注意：config 扁平格式存储
# ============================================
config:
  # 颜色配置（对应CSS变量）
  primaryColor: "#2563eb"
  primaryDark: "#1e40af"
  accentColor: "#60a5fa"

  # 简单值
  siteName: 我的网站
  sidebarEnabled: true

  # 数组值（JSON字符串格式）
  features: '[{"icon":"star","title":"特性1"},{"icon":"heart","title":"特性2"}]'

# ============================================
# configDefinition - 配置表单定义
# 定义用户在后台可以修改哪些配置项
# ============================================
configDefinition:
  - section: basic
    label: 基础设置
    fields:
      - name: siteName
        type: text
        label: 站点名称
        default: 我的网站

  - section: colors
    label: 颜色设置
    fields:
      - name: primaryColor
        type: color
        label: 主色调
        default: "#2563eb"
```

### 3.2 字段类型详解

| type 值 | 表单组件 | 说明 | 示例 |
|---------|---------|------|------|
| `text` | 单行输入框 | 文本输入 | 网站名称 |
| `textarea` | 多行文本域 | 长文本 | 备案号 |
| `number` | 数字输入 | 数值 | 每页文章数 |
| `color` | 颜色选择器 | 颜色选择 | 主色调 |
| `switch` | 开关 | true/false | 显示侧边栏 |
| `select` | 下拉选择 | 单选 | 布局风格 |
| `radio` | 单选按钮 | 单选 | 列表样式 |
| `checkbox` | 复选框 | 多选 | 显示哪些模块 |
| `image` | 图片上传 | 图片地址 | Logo |
| `datetime` | 日期时间 | 日期+时间 | 定时发布 |
| `daterange` | 日期范围 | 日期区间 | 筛选日期 |
| `slider` | 滑块 | 范围值 | 透明度 |
| `code` | 代码编辑 | 代码输入 | 自定义CSS |
| `repeater` | 可视化数组 | 数组编辑 | 特性列表 |

### 3.3 ⚠️ configDefinition 完整结构（重要）

```yaml
configDefinition:
  - section: basic                    # 唯一标识符
    label: 基础设置                   # 后台显示的分组名
    icon: fa-solid fa-gear            # Font Awesome 图标
    fields:
      - name: siteName               # 字段名（唯一）
        type: text                  # 字段类型
        label: 站点名称              # 显示名称
        description: 说明文字        # 可选，字段说明
        default: 我的网站            # 可选默认值
        required: true              # 可选，是否必填
        placeholder: 占位符          # 可选
        options:                    # 可选，select/radio/checkbox的选项
          - value: v1
            label: 选项1
          - value: v2
            label: 选项2
        show_if: "fieldName == value"  # 可选，条件显示
```

### 3.4 ⚠️ settings.yaml 已废弃

> **重要提醒：** 系统不再使用 `settings.yaml`。所有配置表单定义必须放在 `theme.yml` 的 `configDefinition` 中。这确保了配置项与主题元数据统一管理，也保证了配置能正确加载到管理后台。

---

# 第二部分：核心原则（重要！）

## 4. CSS 样式规范

### 4.1 核心原则

> **✅ 统一封装调用：所有样式必须封装在 CSS 文件中（base.css/home.css/list.css/article.css）**
>
> **❌ 禁止在 HTML 模板中使用内联 `style="..."` 属性**
>
> **❌ 禁止在 HTML 模板中嵌入 `<style>` 块（CSS 变量除外，见下节）**

### 4.2 正确做法

```html
<!-- ✅ 正确：在HTML中使用CSS类 -->
<header class="header">
    <nav class="nav-menu">
        <a href="/" class="nav-link active">首页</a>
    </nav>
</header>

<!-- ❌ 错误：内联样式 -->
<header style="background: white; padding: 15px;">
    <nav style="display: flex; gap: 20px;">
        <a href="/" style="color: #2563eb;">首页</a>
    </nav>
</header>
```

### 4.3 CSS 文件组织

```
assets/css/
├── base.css       # 共享基础样式（reset、header、footer、通用组件、空状态）
├── home.css       # 首页特定样式
├── article.css    # 文章/Markdown样式
└── list.css       # 列表页样式（分类、标签、归档、搜索）
```

### 4.4 base.css 包含内容

```css
/* ============================================
   1. CSS Variables - 设计Token（默认值）
   ============================================ */
:root {
    /* 颜色 - 默认值，可被主题配置覆盖 */
    --primary: #2563eb;
    --primary-hover: #1d4ed8;
    --accent: #2563eb;
    --bg-primary: #f8fafc;
    --bg-secondary: #ffffff;
    --text-primary: #1e293b;
    --text-secondary: #64748b;
    --text-muted: #94a3b8;
    --border-color: #e2e8f0;

    /* 字体 */
    --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --font-mono: 'SF Mono', Monaco, Inconsolata, monospace;

    /* 间距 */
    --spacing-xs: 0.25rem;
    --spacing-sm: 0.5rem;
    --spacing-md: 1rem;
    --spacing-lg: 1.5rem;
    --spacing-xl: 2rem;
    --spacing-2xl: 3rem;

    /* 圆角 */
    --radius-sm: 0.25rem;
    --radius-md: 0.5rem;
    --radius-lg: 0.75rem;
    --radius-xl: 1rem;

    /* 阴影 */
    --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);

    /* 布局 */
    --max-width: 1200px;
}

/* ============================================
   2. Reset & Base
   ============================================ */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
html { font-size: 16px; scroll-behavior: smooth; }
body {
    font-family: var(--font-sans);
    font-size: 1rem;
    line-height: 1.6;
    color: var(--text-primary);
    background: var(--bg-primary);
    min-height: 100vh;
}
a { color: var(--accent); text-decoration: none; transition: color 0.2s; }
a:hover { color: var(--primary-hover); }
img { max-width: 100%; height: auto; display: block; }

/* ============================================
   3. Header / Navbar
   ============================================ */
.header { background: var(--bg-secondary); border-bottom: 1px solid var(--border-color); }
.header-container { max-width: var(--max-width); margin: 0 auto; padding: 15px 24px; display: flex; justify-content: space-between; align-items: center; }
.logo { font-size: 1.5rem; font-weight: 700; color: var(--text-primary); }
.nav-menu { display: flex; gap: 24px; }
.nav-link { color: var(--text-secondary); font-weight: 500; padding: 8px 0; }
.nav-link:hover, .nav-link.active { color: var(--primary); }

/* ============================================
   4. Footer
   ============================================ */
.footer { background: var(--primary); color: rgba(255,255,255,0.9); padding: 40px 24px; }
.footer-container { max-width: var(--max-width); margin: 0 auto; text-align: center; }

/* ============================================
   5. Empty State（空状态）
   ============================================ */
.empty-state { text-align: center; padding: 60px 24px; color: var(--text-secondary); }
.empty-state i { font-size: 48px; margin-bottom: 16px; color: var(--text-muted); }

/* ============================================
   6. Page Header（页面标题区）
   ============================================ */
.page-header { background: linear-gradient(135deg, var(--primary-dark, #0f3460) 0%, var(--primary-darker, #16213e) 100%); color: white; padding: 80px 30px; text-align: center; }
.page-header h1 { font-size: 2.5em; margin-bottom: 15px; }

/* ============================================
   7. Responsive
   ============================================ */
@media (max-width: 768px) {
    .header-container { flex-direction: column; gap: 12px; }
    .nav-menu { gap: 16px; }
}
```

### 4.5 CSS 类命名规范

使用语义化命名，避免使用 `style="..."` 内联样式：

| 场景 | 类名示例 | 说明 |
|------|----------|------|
| 分类卡片 | `.category-card`, `.category-link` | 分类相关 |
| 标签云 | `.tags-cloud`, `.tag-item`, `.tag-count` | 标签相关 |
| 归档列表 | `.archive-year`, `.archive-list`, `.archive-item`, `.archive-date`, `.archive-title` | 归档相关 |
| 搜索表单 | `.search-form`, `.search-input`, `.search-btn`, `.results-info` | 搜索相关 |
| 空状态 | `.empty-state`, `.empty-state-full` | 无内容提示 |
| 布局工具 | `.main-full`, `.main-padded`, `.full-width`, `.span-full` | 布局辅助 |

---

## 5. CSS 变量系统

### 5.1 为什么使用 CSS 变量

CSS 变量让主题配置成为可能，用户在后台修改颜色，前台自动更新。

### 5.2 变量定义位置

CSS 变量的**默认值**定义在 `base.css` 的 `:root` 中：

```css
/* assets/css/base.css */
:root {
    --primary: #2563eb;           /* 默认值 */
    --primary-dark: #0f3460;
    --primary-darker: #16213e;
    --accent: #e94560;
    --bg-primary: #f8fafc;
    --bg-secondary: #ffffff;
    --text-primary: #1e293b;
    --text-secondary: #64748b;
    --text-muted: #9ca3af;
    --border-color: #e5e7eb;
    --success: #10b981;
}
```

### 5.3 ⚠️ CSS 变量名与 config key 必须一致

> **重要：** 为了确保配置系统正常工作，CSS 变量名与 `theme.yml` 中 `configDefinition` 的 `name` 必须一致。

例如：`theme.yml` 中配置：
```yaml
configDefinition:
  - section: colors
    fields:
      - name: primaryDark     # ← 配置 key
        type: color
        label: 深主色调
        default: "#0f3460"
```

对应 CSS：
```css
:root {
    --primary-dark: {{ getConfig "primaryDark" "#0f3460" }};  /* ← 变量名一致 */
}
```

### 5.4 完整 CSS 变量清单

```css
:root {
    /* ===== 颜色系统 ===== */
    --primary: #2563eb;           /* 主色 */
    --primary-hover: #1d4ed8;      /* 主色悬停 */
    --primary-dark: #0f3460;       /* 深主色 */
    --primary-darker: #16213e;     /* 更深主色 */
    --accent: #2563eb;             /* 强调色 */
    --success: #10b981;            /* 成功色 */
    --warning: #f59e0b;           /* 警告色 */
    --danger: #ef4444;           /* 危险色 */

    /* ===== 背景色 ===== */
    --bg-primary: #f8fafc;        /* 主背景 */
    --bg-secondary: #ffffff;       /* 次背景 */

    /* ===== 文字色 ===== */
    --text-primary: #1e293b;       /* 主文字 */
    --text-secondary: #64748b;    /* 次文字 */
    --text-muted: #94a3b8;        /* 弱化文字 */

    /* ===== 边框 ===== */
    --border-color: #e5e7eb;       /* 边框色 */

    /* ===== 字体 ===== */
    --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --font-mono: 'SF Mono', Monaco, 'Fira Code', monospace;

    /* ===== 间距 ===== */
    --spacing-xs: 0.25rem;
    --spacing-sm: 0.5rem;
    --spacing-md: 1rem;
    --spacing-lg: 1.5rem;
    --spacing-xl: 2rem;
    --spacing-2xl: 3rem;

    /* ===== 圆角 ===== */
    --radius-sm: 0.25rem;
    --radius-md: 0.5rem;
    --radius-lg: 0.75rem;
    --radius-xl: 1rem;

    /* ===== 阴影 ===== */
    --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);

    /* ===== 布局 ===== */
    --max-width: 1200px;
}
```

### 5.5 使用 CSS 变量

```css
/* 按钮 */
.btn-primary {
    background: var(--primary);
    color: white;
    border-radius: var(--radius-md);
}
.btn-primary:hover { background: var(--primary-hover); }

/* 卡片 */
.card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
}
```

---

## 6. Go 模板限制

### 6.1 ⚠️ 重要限制

> Catter CMS 使用 Go 标准库 `html/template` 引擎。**不支持跨文件模板包含**。
>
> `{{ template "模块名" . }}` 语法在 Go 模板中有特殊要求：所有被引用的模板必须在同一个 `*template.Template` 对象中通过 `ParseFiles` 解析。
>
> 因此，**不要尝试创建 `_head.html`、`_header.html`、`_footer.html` 等 partial 模板并在其他模板中引用**，这会导致 `no such template` 错误。

### 6.2 正确做法

**每个页面模板都包含完整的 HTML 结构：**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .site.title }}</title>
    <link rel="stylesheet" href="/themes/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/home.css" }}">
</head>
<body>
    <!-- Header -->
    <header class="header">...</header>

    <!-- Content -->
    <main>...</main>

    <!-- Footer -->
    <footer class="footer">...</footer>
</body>
</html>
```

### 6.3 避免的做法

```html
<!-- ❌ 错误：尝试引用 partial 模板 -->
<head>
    {{ template "_head" . }}   <!-- 会导致 "no such template" 错误 -->
</head>

<!-- ✅ 正确：每个页面包含完整的 head -->
<head>
    <meta charset="UTF-8">
    <title>{{ .site.title }}</title>
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
</head>
```

### 6.4 CSS/JS 引用规范

```html
<!-- Font Awesome（全局可用） -->
<link rel="stylesheet" href="/themes/fontawesome/css/all.min.css">

<!-- 主题CSS（使用 asset_path） -->
<link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
<link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/home.css" }}">
```

---

# 第三部分：模板开发

## 7. 模板基础

### 7.1 Go 模板语法快速入门

Catter CMS 使用 Go 的 `html/template` 引擎，语法简单但功能强大。

#### 7.1.1 变量输出

```html
<!-- 输出字符串 -->
{{ .site.title }}

<!-- 输出数字 -->
{{ .post.Views }}
```

#### 7.1.2 条件判断

```html
<!-- 基本 if -->
{{ if .post }}
    <p>文章存在</p>
{{ end }}

<!-- if-else -->
{{ if .site.description }}
    <p>{{ .site.description }}</p>
{{ else }}
    <p>暂无描述</p>
{{ end }}

<!-- 嵌套条件 -->
{{ if and .post .post.Category }}
    <span>分类：{{ .post.Category.Name }}</span>
{{ end }}
```

#### 7.1.3 循环遍历

```html
<!-- 遍历文章列表 -->
{{ range .posts }}
    <article>
        <h2><a href="/article/{{ .Slug }}">{{ .Title }}</a></h2>
        <p>{{ .Excerpt }}</p>
    </article>
{{ end }}

<!-- 带 else 的遍历（空列表时显示） -->
{{ range .posts }}
    <li>{{ .Title }}</li>
{{ else }}
    <li>暂无文章</li>
{{ end }}
```

### 7.2 模板中的数据结构

#### 7.2.1 site 对象（全局）

```go
site = {
    title: "网站名称",
    description: "网站描述",
    year: 2026,
    pages: [...]  // 所有页面列表
}
```

#### 7.2.2 posts 数组（首页/列表页）

```go
posts = [
    {
        ID: 1,
        Title: "文章标题",
        Slug: "article-slug",
        Content: "Markdown内容...",
        Excerpt: "文章摘要...",
        Views: 100,
        PublishedAt: 2026-03-24,
        Category: { ID, Name, Slug },
        Tags: [{ ID, Name, Slug }, ...]
    },
    ...
]
```

### 7.3 模板中的 URL 生成

```html
<!-- 首页 -->
<a href="/">首页</a>

<!-- 文章详情页 -->
<a href="/article/{{ .post.Slug }}">{{ .post.Title }}</a>

<!-- 分类页 -->
<a href="/category/{{ .category.Slug }}">{{ .category.Name }}</a>

<!-- 标签页 -->
<a href="/tag/{{ .tag.Slug }}">{{ .tag.Name }}</a>

<!-- 文章列表页 -->
<a href="/article">全部文章</a>

<!-- 独立页面 -->
<a href="/page/about">关于我们</a>

<!-- 资源路径（推荐） -->
<link href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
```

---

## 8. 模板函数

### 8.1 常用函数速查表

| 函数 | 用法 | 说明 | 示例 |
|------|------|------|------|
| `getConfig` | `getConfig "key" default` | 获取主题配置 | `getConfig "primaryColor" "#2563eb"` |
| `len` | `len .posts` | 返回长度 | `len .posts` |
| `asset_path` | `asset_path "css/base.css"` | 资源路径 | `asset_path "images/logo.png"` |
| `current_theme` | `current_theme` | 当前主题名 | `current_theme` |
| `safe` | `safe .htmlContent` | 标记安全HTML | `markdownHTML .Content \| safe` |
| `default` | `default "value" "default"` | 设置默认值 | `default .value "N/A"` |
| `lower` | `lower .Title` | 转小写 | `lower .Title` |
| `upper` | `upper .Title` | 转大写 | `upper .Title` |
| `trim` | `trim .Content` | 去除首尾空白 | `trim .Content` |
| `substr` | `substr .str 0 100` | 截取字符串 | `substr .Content 0 200` |
| `markdownHTML` | `markdownHTML .md` | Markdown转HTML | `markdownHTML .Content` |
| `json_decode` | `json_decode .jsonStr` | JSON解析 | `json_decode .features` |
| `dateFormat` | `dateFormat "2006-01-02" .Date` | 日期格式化 | `dateFormat "2006-01-02" .PublishedAt` |
| `add` | `add 1 2` | 加法 | `add .Views 1` |
| `sub` | `sub 5 2` | 减法 | `sub .Views 1` |
| `mod` | `mod 5 2` | 模运算 | `mod .index 2` |
| `div` | `div 10 3` | 整数除法 | `div .index 10` |
| `first` | `first .list` | 第一个元素 | `first .posts` |
| `int` | `int .floatVal` | 转整数 | `int .value` |
| `iterate` | `iterate 5` | 生成序列 | `range (iterate 5)` |
| `eq` | `eq .a .b` | 相等判断 | `eq .Status "published"` |
| `gt` | `gt .a .b` | 大于 | `gt .Views 100` |
| `lt` | `lt .a .b` | 小于 | `lt .Views 100` |
| `and` | `and .a .b` | 逻辑与 | `and .post .post.Category` |

### 8.2 getConfig 函数详解（最重要）

这是主题配置的核心函数，必须掌握。

```html
<!-- 基本用法 -->
{{ getConfig "configKey" "默认值" }}

<!-- 颜色默认值 -->
{{ getConfig "primaryColor" "#2563eb" }}
```

> ⚠️ **注意：** `getConfig` 的 key 必须与 `theme.yml` 中 `configDefinition` 的 `name` 一致。

### 8.3 markdownHTML + safe 组合

富文本内容必须使用这个组合：

```html
<!-- 文章内容（Markdown 转 HTML） -->
<article class="article-content">
    {{ markdownHTML .post.Content | safe }}
</article>
```

### 8.4 日期格式化

```html
<!-- 标准日期 -->
{{ dateFormat "2006-01-02" .post.PublishedAt }}
<!-- 输出: 2026-03-24 -->

<!-- 友好格式 -->
{{ dateFormat "Jan 02, 2006" .post.PublishedAt }}
<!-- 输出: Mar 24, 2026 -->

<!-- 中文格式 -->
{{ dateFormat "2006年1月2日" .post.PublishedAt }}
<!-- 输出: 2026年3月24日 -->
```

### 8.5 json_decode 处理数组配置

主题配置中的数组以 JSON 字符串存储，需要解码：

```yaml
# theme.yml 中的配置
config:
  features: '[{"icon":"star","title":"特性1"},{"icon":"heart","title":"特性2"}]'
```

```html
<!-- 模板中使用 -->
{{ range (json_decode .themeConfig.features) }}
    <div class="feature">
        <i class="fas fa-{{ .icon }}"></i>
        <h3>{{ .title }}</h3>
    </div>
{{ end }}
```

---

## 9. 首页模板开发

### 9.1 首页数据结构

```go
{
    site: { title, description, year, pages },
    themeConfig: { ...主题配置... },
    posts: [ ...最新文章列表... ],
    categories: [ ...分类列表... ],
    tags: [ ...标签列表... ],
    pagination: { page, totalPages, total }
}
```

### 9.2 完整首页模板示例

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .site.title }}</title>
    <meta name="description" content="{{ .site.description }}">

    <!-- 引入Font Awesome -->
    <link rel="stylesheet" href="/themes/fontawesome/css/all.min.css">

    <!-- 引入主题CSS -->
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/home.css" }}">
</head>
<body>
    <!-- 头部导航 -->
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                <a href="/article" class="nav-link">博客</a>
                <a href="/categories" class="nav-link">分类</a>
                <a href="/tags" class="nav-link">标签</a>
                <a href="/page/about" class="nav-link">关于</a>
            </nav>
        </div>
    </header>

    <!-- Hero区域（可配置开关） -->
    {{ if getConfig "heroEnabled" true }}
    <section class="hero">
        <h1>{{ getConfig "heroTitle" "欢迎来到" }} {{ .site.title }}</h1>
        <p>{{ getConfig "heroSubtitle" "分享技术与生活" }}</p>
        {{ if getConfig "heroButtonText" "" }}
        <a href="{{ getConfig "heroButtonLink" "/" }}" class="hero-btn">
            {{ getConfig "heroButtonText" "了解更多" }}
        </a>
        {{ end }}
    </section>
    {{ end }}

    <!-- 主内容区 -->
    <main class="main-container">
        <!-- 文章列表 -->
        <section class="posts-section">
            <h2 class="section-title">最新文章</h2>

            {{ range .posts }}
            <article class="post-card">
                <h3><a href="/article/{{ .Slug }}">{{ .Title }}</a></h3>
                <div class="post-meta">
                    {{ if .PublishedAt }}
                    <span>{{ dateFormat "2006-01-02" .PublishedAt }}</span>
                    {{ end }}
                    {{ if .Category }}
                    <span>{{ .Category.Name }}</span>
                    {{ end }}
                    <span>{{ .Views }} 阅读</span>
                </div>
                {{ if .Excerpt }}
                <p class="post-excerpt">{{ substr .Excerpt 0 150 }}...</p>
                {{ end }}
                <a href="/article/{{ .Slug }}" class="read-more">阅读全文 →</a>
            </article>
            {{ else }}
            <div class="empty-state">
                <p>暂无文章</p>
            </div>
            {{ end }}
        </section>

        <!-- 侧边栏 -->
        {{ if getConfig "sidebarEnabled" true }}
        <aside class="sidebar">
            <div class="widget">
                <h4 class="widget-title">分类</h4>
                <ul class="widget-list">
                    {{ range .categories }}
                    <li><a href="/category/{{ .Slug }}">{{ .Name }}</a></li>
                    {{ end }}
                </ul>
            </div>

            <div class="widget">
                <h4 class="widget-title">标签</h4>
                <div class="tags-cloud">
                    {{ range .tags }}
                    <a href="/tag/{{ .Slug }}" class="tag-item">{{ .Name }}</a>
                    {{ end }}
                </div>
            </div>
        </aside>
        {{ end }}
    </main>

    <!-- 页脚 -->
    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}. All rights reserved.</p>
        </div>
    </footer>
</body>
</html>
```

---

## 10. 文章详情页

### 10.1 数据结构

```go
{
    site: { title, description, year },
    themeConfig: { ... },
    post: {
        ID, Title, Slug, Content, Excerpt,
        Views, PublishedAt, CreatedAt,
        Category: { ID, Name, Slug },
        Tags: [{ ID, Name, Slug }, ...]
    }
}
```

### 10.2 完整文章详情模板

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .post.Title }} - {{ .site.title }}</title>
    <meta name="description" content="{{ .post.Excerpt }}">

    <link rel="stylesheet" href="/themes/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/article.css" }}">
</head>
<body>
    <!-- 头部 -->
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                <a href="/article" class="nav-link">博客</a>
                <a href="/categories" class="nav-link">分类</a>
            </nav>
        </div>
    </header>

    <!-- 文章头部 -->
    <div class="article-header">
        <div class="article-meta">
            {{ if .post.PublishedAt }}
            <span>{{ dateFormat "2006-01-02" .post.PublishedAt }}</span>
            {{ end }}
            {{ if .post.Category }}
            <a href="/category/{{ .post.Category.Slug }}">{{ .post.Category.Name }}</a>
            {{ end }}
            <span>{{ .post.Views }} 阅读</span>
        </div>
        <h1>{{ .post.Title }}</h1>
    </div>

    <!-- 文章内容 -->
    <article class="article-content">
        {{ markdownHTML .post.Content | safe }}
    </article>

    <!-- 文章标签 -->
    {{ if .post.Tags }}
    <div class="article-tags">
        <span>标签：</span>
        {{ range .post.Tags }}
        <a href="/tag/{{ .Slug }}" class="tag-item">{{ .Name }}</a>
        {{ end }}
    </div>
    {{ end }}

    <!-- 页脚 -->
    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```

---

## 11. 列表页模板

### 11.1 分类页（category.html）

**URL:** `/category/:slug`

**数据:** `{ site, themeConfig, category, posts, pagination }`

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .category.Name }} - {{ .site.title }}</title>
    <link rel="stylesheet" href="/themes/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/list.css" }}">
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                <a href="/article" class="nav-link">博客</a>
                <a href="/categories" class="nav-link">分类</a>
            </nav>
        </div>
    </header>

    <div class="page-header">
        <h1>{{ .category.Name }}</h1>
    </div>

    <main class="list-container">
        {{ range .posts }}
        <article class="article-card">
            <h2><a href="/article/{{ .Slug }}">{{ .Title }}</a></h2>
            <div class="article-meta">
                {{ if .PublishedAt }}
                <span>{{ dateFormat "2006-01-02" .PublishedAt }}</span>
                {{ end }}
                <span>{{ .Views }} 阅读</span>
            </div>
            {{ if .Excerpt }}
            <p>{{ substr .Excerpt 0 150 }}...</p>
            {{ end }}
        </article>
        {{ else }}
        <div class="empty-state">
            <p>该分类下暂无文章</p>
        </div>
        {{ end }}
    </main>

    <footer class="footer">
        <p>&copy; {{ .site.year }} {{ .site.title }}</p>
    </footer>
</body>
</html>
```

### 11.2 标签页（tag.html）

**URL:** `/tag/:slug`

**数据:** `{ site, themeConfig, tag, posts, pagination }`

### 11.3 归档页（archives.html）

**URL:** `/archives`

**数据:** `{ site, themeConfig, archives: { year: [posts...] } }`

```html
<!-- 按年份分组的归档 -->
{{ range $year, $items := .archives }}
<div class="archive-year">
    <h2>{{ $year }}</h2>
    <ul class="archive-list">
        {{ range $items }}
        <li class="archive-item">
            <span class="archive-date">{{ dateFormat "01-02" .PublishedAt }}</span>
            <a href="/article/{{ .Slug }}" class="archive-title">{{ .Title }}</a>
        </li>
        {{ end }}
    </ul>
</div>
{{ end }}
```

### 11.4 标签云页（tags.html）

**URL:** `/tags`

**数据:** `{ site, themeConfig, tags: [{ Name, Slug, PostCount }] }`

```html
<div class="tags-cloud">
    {{ range .tags }}
    <a href="/tag/{{ .Slug }}" class="tag-item">
        <span>{{ .Name }}</span>
        <span class="tag-count">{{ .PostCount }}</span>
    </a>
    {{ end }}
</div>
```

---

# 第四部分：高级配置

## 12. 主题配置定义

### 12.1 configDefinition 完整结构

```yaml
configDefinition:
  - section: section_id           # 唯一标识符
    label: 分组显示名称            # 后台显示的分组名
    icon: fa-solid fa-gear       # Font Awesome 图标类名
    fields:                       # 该分组下的字段列表
      - name: field_name         # 字段名（唯一）
        type: text               # 字段类型
        label: 字段标签           # 显示名称
        description: 说明文字     # 可选，字段说明
        default: 默认值          # 可选
        required: true           # 可选，是否必填
        placeholder: 占位符       # 可选
        options:                 # 可选，select/radio/checkbox的选项
          - value: v1
            label: 选项1
          - value: v2
            label: 选项2
        show_if: "condition"    # 可选，条件显示
```

### 12.2 完整配置示例

```yaml
# theme.yml
config:
  # 颜色配置
  primaryColor: "#2563eb"
  primaryDark: "#0f3460"
  accent: "#e94560"
  bgPrimary: "#f8f9fa"
  bgSecondary: "#ffffff"
  textPrimary: "#1e293b"
  textSecondary: "#6b7280"
  textMuted: "#9ca3af"
  borderColor: "#e5e7eb"

configDefinition:
  # 分组1：颜色设置
  - section: colors
    label: 颜色设置
    icon: fa-solid fa-palette
    fields:
      - name: primaryColor
        type: color
        label: 主色调
        default: "#2563eb"
      - name: primaryDark
        type: color
        label: 深主色调
        default: "#0f3460"

  # 分组2：Hero区域
  - section: hero
    label: Hero区域
    icon: fa-solid fa-image
    fields:
      - name: heroEnabled
        type: switch
        label: 显示Hero
        default: true
      - name: heroTitle
        type: text
        label: 主标题
        show_if: "heroEnabled == true"
        default: 欢迎来到
```

---

## 13. Repeater字段

### 13.1 repeater vs textarea

| 对比项 | textarea (JSON) | repeater |
|--------|-----------------|----------|
| 用户体验 | 手动编辑JSON，易错 | 可视化编辑，直观 |
| 适用场景 | 简单数组 | 复杂数组 |
| 配置复杂度 | 低 | 高 |
| 数据格式 | JSON字符串 | JSON字符串（相同） |

### 13.2 常用 repeater 场景

**场景1：功能特性列表**
```yaml
features:
  - icon: "fa-solid fa-bolt"
    title: 高性能
    description: 基于Go语言开发
  - icon: "fa-solid fa-shield"
    title: 安全可靠
    description: 多重安全防护
```

**场景2：统计数据**
```yaml
stats:
  - value: "10K+"
    label: 用户
  - value: "50+"
    label: 主题
```

### 13.3 repeater 模板渲染

```html
{{ range (json_decode .themeConfig.features) }}
<div class="feature-card">
    <i class="{{ .icon }}"></i>
    <h3>{{ .title }}</h3>
    <p>{{ .description }}</p>
</div>
{{ end }}
```

---

## 14. 条件显示

### 14.1 show_if 语法

`show_if` 用于根据其他配置的值决定是否显示某个字段。

```yaml
# 基本语法
show_if: "fieldName == value"
show_if: "fieldName != value"
show_if: "fieldName"           # 布尔值检查

# AND 条件
show_if: "field1 == value1 && field2 == value2"

# OR 条件
show_if: "field1 == value1 || field2 == value2"
```

### 14.2 使用示例

```yaml
configDefinition:
  - section: hero
    label: Hero区域
    fields:
      - name: heroEnabled
        type: switch
        label: 显示Hero区域
        default: true

      # 仅当 heroEnabled == true 时显示
      - name: heroTitle
        type: text
        label: 主标题
        show_if: "heroEnabled == true"
        default: 欢迎来到

      - name: heroSubtitle
        type: textarea
        label: 副标题
        show_if: "heroEnabled == true"
```

---

# 第五部分：实战案例

## 15. 案例：博客主题

### 15.1 主题设计

**定位：** 个人技术博客
**风格：** 简洁、现代、注重阅读体验
**配色：** 蓝色主色调 (#2563eb)，白底灰字

### 15.2 目录结构

```
templates/themes/my-blog/
├── theme.yml
├── assets/
│   └── css/
│       └── base.css
└── pages/
    ├── home.html
    ├── article_show.html
    ├── article_list.html
    ├── category.html
    ├── tag.html
    └── search.html
```

### 15.3 theme.yml

```yaml
name: My Blog
slug: my-blog
version: 1.0.0
author: Your Name
description: 简洁现代的个人博客主题

config:
  primaryColor: "#2563eb"
  sidebarEnabled: true
  postsPerPage: 10
  showViews: true

configDefinition:
  - section: display
    label: 显示设置
    fields:
      - name: sidebarEnabled
        type: switch
        label: 显示侧边栏
        default: true
      - name: showViews
        type: switch
        label: 显示阅读量
        default: true
```

---

## 16. 案例：企业官网

### 16.1 主题设计

**定位：** 企业官方网站
**风格：** 专业、商务、模块化
**配色：** 深蓝色主调 (#0f3460)，强调色红色 (#e94560)

### 16.2 目录结构

```
templates/themes/my-enterprise/
├── theme.yml
├── assets/
│   └── css/
│       ├── base.css        # 共享样式（包含CSS变量默认值）
│       ├── home.css        # 首页模块
│       ├── article.css     # 文章/Markdown
│       └── list.css        # 列表页
└── pages/
    ├── home.html           # 首页（多模块）
    ├── article_show.html   # 新闻详情
    ├── article_list.html   # 新闻列表
    └── page.html           # 独立页面
```

### 16.3 首页模板结构

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .site.title }}</title>
    <link rel="stylesheet" href="/themes/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-enterprise/assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "/themes/my-enterprise/assets/css/home.css" }}">
</head>
<body>
    <!-- Header -->
    <header>...</header>

    <!-- Hero -->
    {{ if .themeConfig.heroEnabled }}
    <section class="hero">...</section>
    {{ end }}

    <!-- Stats -->
    {{ if .themeConfig.moduleStatsEnabled }}
    <section class="stats">...</section>
    {{ end }}

    <!-- Features -->
    {{ if .themeConfig.moduleFeaturesEnabled }}
    <section class="features">...</section>
    {{ end }}

    <!-- CTA -->
    {{ if .themeConfig.ctaEnabled }}
    <section class="cta">...</section>
    {{ end }}

    <!-- Footer -->
    <footer>...</footer>
</body>
</html>
```

---

# 第六部分：参考手册

## 17. API参考

### 17.1 主题管理 API

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/themes` | 获取主题列表 |
| POST | `/api/themes/:id/activate` | 激活主题 |
| GET | `/api/themes/:slug/config` | 获取主题配置 |
| PUT | `/api/themes/:slug/config` | 保存主题配置 |
| GET | `/api/themes/:slug/config/full` | 获取完整配置（含表单定义） |
| POST | `/api/themes/upload` | 上传主题（ZIP） |
| DELETE | `/api/themes/:id` | 删除主题 |

### 17.2 API 响应格式

**成功响应：**
```json
{
    "code": 0,
    "message": "操作成功",
    "data": { ... }
}
```

**错误响应：**
```json
{
    "code": 1,
    "message": "错误信息",
    "data": null
}
```

---

## 18. 数据模型

### 18.1 Post（文章）

```go
type Post struct {
    ID          uint       // 主键
    Title       string     // 标题
    Slug        string     // URL友好名称
    Content     string     // Markdown内容
    Excerpt     string     // 摘要
    Views       int        // 阅读量
    Status      string     // published, draft
    PublishedAt *time.Time // 发布时间
    CreatedAt   time.Time  // 创建时间
    UpdatedAt   time.Time  // 更新时间
    CategoryID  uint       // 分类ID
    Category    *Category  // 关联分类
    Tags        []*Tag     // 关联标签
}
```

### 18.2 Category（分类）

```go
type Category struct {
    ID          uint       // 主键
    Name        string     // 分类名称
    Slug        string     // URL友好名称
    Description string     // 描述
    ParentID    uint       // 父分类ID
    Sort        int        // 排序
}
```

### 18.3 Tag（标签）

```go
type Tag struct {
    ID        uint      // 主键
    Name      string    // 标签名称
    Slug      string    // URL友好名称
}
```

---

## 19. 常见问题

### Q1: 模板修改后没有生效？

**原因：** 浏览器缓存
**解决：** 强制刷新（Ctrl+F5）或清除缓存

### Q2: 主题激活后页面显示空白？

**排查步骤：**
1. 检查服务器日志是否有错误
2. 检查浏览器控制台（F12）
3. 验证 `theme.yml` 的 `slug` 与文件夹名一致
4. 验证 `pages/home.html` 存在

### Q3: CSS 样式没有加载？

**排查步骤：**
1. 检查 CSS 文件路径是否正确
2. 确认文件放在 `assets/css/` 目录
3. 使用 `asset_path` 函数生成路径
4. 检查浏览器 Network 标签页是否有 404

### Q4: 配置保存后没有生效？

**原因：** 缓存问题
**解决：** 重启服务器或重新激活主题

### Q5: "no such template" 错误？

**原因：** 尝试使用 `{{ template "_head" . }}` 引用 partial 模板，但 Go 模板引擎不支持跨文件模板包含。

**解决：** 每个页面模板都包含完整的 HTML 结构，不要尝试拆分 partial 模板。

### Q6: 配置项在管理后台不显示？

**排查步骤：**
1. 确认配置定义写在 `theme.yml` 的 `configDefinition` 中（不是 `settings.yaml`）
2. 检查 `name` 字段是否与模板中 `getConfig` 的 key 一致
3. 检查字段类型是否正确

### Q7: 模板中出现 "unexpected EOF" 错误？

**常见原因：**
1. 模板语法错误，如 `{{ if` 没有对应的 `{{ end }}`
2. 文件编码问题（确保是 UTF-8）

### Q8: 中文显示乱码？

**解决：**
1. 确保模板文件是 UTF-8 编码（无 BOM）
2. 使用 `<meta charset="UTF-8">`

### Q9: 颜色配置不生效？

**排查步骤：**
1. 检查 `theme.yml` 中 `configDefinition` 的 `name` 与 CSS 变量名是否一致
2. 例如：配置 key 为 `primaryDark`，CSS 中对应 `--primary-dark`

### Q10: "function not defined" 错误？

**原因：** Go 模板语法错误，把字段访问误解析为函数调用。

**常见错误：**
```html
<!-- ❌ 错误：dot 和字段名之间有空格 -->
<p>{{ . Excerpt }}</p>

<!-- ✅ 正确：紧挨着写 -->
<p>{{ .Excerpt }}</p>
```

**说明：** Go 模板引擎将 `{{ . Excerpt }}` 解析为调用名为 `Excerpt` 的函数（传入 dot 作为参数），而非访问 `Excerpt` 字段。字段访问时 `.` 和字段名之间不能有任何空白字符。

### Q11: "can't evaluate field Cover" 错误？

**原因：** 模板引用了 Post 模型中不存在的字段。

**排查步骤：**
1. 检查 `Post` 结构体是否确实有该字段（参考本文档第18节数据模型）
2. 常见不存在字段：`Cover`（文章封面图）、`ReadTime`（阅读时间）、`PostCount`（标签文章数，非字段而是方法）

**Post 模型可用字段：**
```go
type Post struct {
    ID, Title, Slug, Content, Excerpt, Views, Status
    PublishedAt, CreatedAt, UpdatedAt
    CategoryID, Category, Tags
    // 注意：没有 Cover、ReadTime 等字段
}
```

### Q12: 模板加载时报 "template not found"？

**原因：** 模板文件路径或命名不正确。

**排查步骤：**
1. 确认文件位于 `templates/themes/<slug>/pages/` 目录
2. 确认 home.html 存在（激活主题的最低要求）
3. 检查文件扩展名是否为 `.html`（不是 `.plush.html`）
4. 查看服务器日志中 `Found X page templates` 确认文件被识别

---

## 附录：重要更新说明

### v2.0 → v3.0 更新内容

1. **CSS 样式规范强化**：明确禁止内联样式，统一封装到 CSS 文件
2. **Go 模板限制说明**：新增第6节，说明跨文件模板包含不可用
3. **settings.yaml 废弃**：`configDefinition` 已成为唯一的配置表单定义方式
4. **CSS 变量命名规范**：强调 CSS 变量名与 config key 必须一致
5. **共享 CSS 类规范**：新增 `.category-*`、`.archive-*`、`.empty-state` 等标准类名

### v3.1 更新内容 (2026-03-25)

**问题修复记录：**

1. **修复了 dianyan-blog 主题的模板语法错误**
   - 问题：`{{ . Excerpt }}` 中 dot 和字段名之间有空格，导致 `function "Excerpt" not defined` 错误
   - 原因：Go 模板将 `. Excerpt` 解析为函数调用而非字段访问
   - 修复：统一为 `{{ .Excerpt }}`（无空格）

2. **修复了 Post 模型不存在的字段引用**
   - 问题：模板引用了 `{{ .Cover }}`，但 Post 结构体没有 Cover 字段
   - 修复：从模板中移除缩略图判断块，或使用 `{{ if .Excerpt }}` 等实际存在的字段

3. **前端编译错误修复**
   - 问题：`ThemeSettings.vue` 中重复导入 `onUnmounted`
   - 修复：删除第177行的重复 import 语句

**经验教训：**
- Go 模板中 `.` 和字段名之间**绝对不能有空格**
- 模板引用字段前，先确认该字段在 Post/Category/Tag 等模型中确实存在
- 激活新主题后，检查服务器日志中的 `Found X page templates` 和任何 `ERROR` 输出

### 配置更新说明

1. **修改 `theme.yml` 配置后**
   - 需要在管理后台点击 **"重载主题配置"** 或重新激活主题

2. **修改模板文件（.html）后**
   - 自动生效，无需重启

3. **修改 CSS/JS/图片后**
   - 强制刷新浏览器缓存（Ctrl+F5）

---

**文档版本：** 3.1
**最后更新：** 2026-03-25
**编写者：** Catter CMS 开发团队
