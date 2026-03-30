---
name: catter-cms-themes
description: This skill should be used whenever developing, creating, or modifying themes for the Catter CMS system. It provides the complete specification-driven workflow for building production-quality CMS themes using Go html/template engine, CSS variable systems, and theme.yml configuration. Trigger when the user asks to create a new theme, modify existing themes, fix theme bugs, or design theme configuration forms for Catter CMS.
---

# Catter CMS Theme Development Skill

这是一个**规范驱动开发（Specification-Driven Development）**的主题开发 Skill。在开发任何主题前，必须完整阅读并遵守本 Skill 中的所有规范，以确保输出符合 Catter CMS 标准。

---

## 快速索引

- **🔴 开发前必读** → [⚠️ 最常见错误 TOP 9](#⚠️-最常见错误-top-9-先读这里)（getConfig / dateFormat / API / 内联样式 / **asset_path 路径** / 硬编码导航 / **markdownHTML 误用** / **Page 字段错误** / **Emoji 禁止**）
- **⚠️ 开发流程（必须从这里开始）** → [Development Workflow](#development-workflow-sdd-pattern)
- **⚠️ Phase 1 Design Brief 格式** → [Phase 1 强制交付物格式](#phase-1-强制交付物格式)
- **⛔ 文件写入安全规范（禁止 PowerShell）** → [文件写入安全规范](#️-文件写入安全规范强制违反会破坏主题)
- **⚠️ Error Prevention Checklist** → [逐项检查清单](#error-prevention-checklist)（**开发完成后必查**）
- 文件结构规范 → [Directory Structure](#directory-structure)
- theme.yml 完整格式 → [theme.yml Reference](#themeyml-reference)
- 模板语法规则 → [Template Syntax Rules](#template-syntax-rules)
- CSS 规范 → [CSS Standards](#css-standards)
- 完整模板函数列表 → [Template Functions](#template-functions)
- 数据模型 → [Data Models](#data-models)
- 详细参考文档 → `references/`（**仅供语法对照，严禁复制**）

---

> ## ⛔ BEFORE YOU WRITE A SINGLE LINE OF CODE
>
> **必须先完成 Phase 1 设计规格（Design Brief），并等待用户确认，才能开始任何编码工作。**
>
> 跳过 Phase 1 是本 Skill 中唯一不可接受的行为。
> 立即跳到写模板 = 输出一个复制品，不是主题。
>
> → 直接跳到 [Development Workflow](#development-workflow-sdd-pattern) 查看 Phase 1 的完整执行步骤。
>
> ---

## ⚠️ 最常见错误 TOP 9（先读这里）

> 以下 8 类错误在过去开发中频繁出现，每次开发主题前请务必确认已理解。
>
> ### 🔴 #1 getConfig key 路径错误（最高频）
>
> **config 数据是扁平结构（Flat）。** `configDefinition` 中 section 的 `key`（如 `basic`、`hero`、`stats`）只是管理后台的**分组标识**，不影响数据访问路径。
>
> ```yaml
> # theme.yml
> configDefinition:
>   - section: basic          # ← 这只是分组标识，不是路径！
>     label: 基础信息
>     fields:
>       - name: phone         # ← 这才是真正的 config key
>         type: text
>       - name: address
>         type: textarea
>   - section: hero           # ← 分组标识
>     label: Hero区域
>     fields:
>       - name: heroEnabled   # ← 真正的 key
>         type: switch
> ```
>
> **正确的访问方式：**
> ```html
> {{ getConfig "phone" .config }}           ✅
> {{ getConfig "heroEnabled" true }}         ✅
> {{ getConfig "address" .config }}         ✅
> ```
>
> **错误演示（以前犯过的）：**
> ```html
> {{ getConfig "basic.phone" .config }}     ❌ 多了 "basic."
> {{ getConfig "hero.heroEnabled" .config }} ❌ 多了 "hero."
> {{ getConfig "stats.statsYears" .config }} ❌ 多了 "stats."
> {{ getConfig "seo.metaDescription" .config }} ❌ 多了 "seo."
> ```
>
> **规则：getConfig 的第一个参数 = configDefinition 中 `name` 字段的值，与 section key 完全无关。**
>
> ### 🔴 #2 markdownHTML | safe 类型冲突（高危，2026-03-27 新增）
>
> **`markdownHTML` 函数已返回 `template.HTML` 类型，绝对不能再接 `| safe`。** 否则会导致模板静默失败，回退到默认模板。
>
> ```html
> <!-- ❌ 致命错误：| safe 导致类型冲突，模板执行失败 -->
> <div class="article-body">{{ markdownHTML .post.Content | safe }}</div>
> <div class="article-body">{{ markdownHTML .page.Content | safe }}</div>
>
> <!-- ✅ 正确：直接使用 markdownHTML -->
> <div class="article-body">{{ markdownHTML .post.Content }}</div>
> <div class="article-body">{{ markdownHTML .page.Content }}</div>
> ```
>
> **⚠️ 此 Bug 的危险之处：** 不会返回 500 错误，而是**静默回退**到 `page.html`。页面看起来能打开但内容为空，极难排查。**服务器 stderr 日志中才有 `[ERROR] ExecuteTemplate failed` 记录。**
>
> **`| safe` 的唯一合法用途：**
> ```html
> {{ getConfig "customCSS" "" | safe }}   <!-- ✅ 管理员自定义 CSS，原本是 string -->
> {{ getConfig "customJS" "" | safe }}    <!-- ✅ 管理员自定义 JS -->
> ```
>
> **详细说明见 Bug 模式 [B017](references/bug-patterns.md#b017)。**

> ### 🔴 #3 Page 模型字段错误（2026-03-27 新增）
>
> **Page 模型没有 `Description` 字段**，使用 `.page.Description` 会导致模板执行失败并静默回退。
>
> ```
> Page 实际字段：ID, Title, Slug, Content, Excerpt, Status, Template,
>               TemplateData, MetaTitle, MetaDescription, Views, SortOrder,
>               CreatedAt, UpdatedAt
> ⚠️ 不存在：Description, Category, Tags, Author
> ```
>
> ```html
> <!-- ❌ 错误 -->
> {{ if .page.Description }}<p>{{ .page.Description }}</p>{{ end }}
>
> <!-- ✅ 正确：使用 Excerpt -->
> {{ if .page.Excerpt }}<p>{{ .page.Excerpt }}</p>{{ end }}
> ```
>
> **详细说明见 Bug 模式 [B018](references/bug-patterns.md#b018)。**

> ### 🔴 #4 dateFormat 参数顺序错误
>
> 函数签名是 **`dateFormat format date`**（格式在前，日期在后）：
>
> ```html
> {{ dateFormat "2006-01-02" .post.PublishedAt }}    ✅ 正确
> {{ dateFormat .post.PublishedAt "2006-01-02" }}    ❌ 顺序反了
> ```
>
> 常见错误写法：`{{ dateFormat .PublishedAt "2006-01-02" }}` — 缺少右括号、顺序反，两种错误同时出现。
>
> ### 🔴 #5 使用不存在的 API
>
> **只使用 Skill 中明确列出的公开 API**，不要猜测：
>
> ```
> ✅ POST /api/comments          ← 正确，评论提交
> ❌ POST /api/contact           ← 不存在，会永远返回 404
> ```
>
> 可用 API 列表见 [前台公开 API](#前台公开api主题可调用)。
>
> ### 🔴 #6 内联样式使用硬编码颜色
>
> ```html
> <div style="color: rgba(255,255,255,0.7)">    ❌ rgba() 硬编码
> <div style="background: #1a1a1a">               ❌ #rrggbb 硬编码
> <div style="border: 1px solid rgb(0,0,0)">     ❌ rgb() 硬编码
> ```
>
> **正确做法：** 所有颜色值必须来自 CSS 变量（`var(--text-primary)` 等），CSS 变量在 `<head>` 的 `:root` 中通过 `getConfig` 注入。
>
> ### 🔴 #7 asset_path 路径重复（高频，2026-03-30 新增）
>
> **`asset_path` 函数**会在参数前自动添加 `/themes/{themeName}/` 前缀。如果传入的参数已经带了 `/themes/xxx/`，会导致路径重复。
>
> ```html
> <!-- ❌ 致命错误：传入完整路径导致重复 -->
> <link rel="stylesheet" href="{{ asset_path "/themes/catter-cms-official/assets/css/base.css" }}">
> <!-- 生成：/themes/catter-cms-official//themes/catter-cms-official/assets/css/base.css -->
>
> <!-- ✅ 正确：只传相对路径 -->
> <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
> <!-- 生成：/themes/catter-cms-official/assets/css/base.css -->
>
> <!-- ❌ 其他错误示例 -->
> <link rel="stylesheet" href="{{ asset_path "/assets/css/base.css" }}">
> <!-- 多了开头的 /，生成：/themes/xxx//assets/css/base.css -->
>
> <link rel="stylesheet" href="{{ asset_path "themes/assets/css/base.css" }}">
> <!-- 少了 assets/，路径不完整 -->
> ```
>
> **规则：`asset_path` 参数必须是相对路径，不能以 `/` 开头，不能包含 `/themes/`。**
>
> ### 🔴 #8 Footer 导航链接硬编码
>
> 禁止在模板中硬编码导航链接，必须使用系统菜单：
>
> ```html
> ❌ <a href="/page/about">关于我们</a>
> ❌ <a href="/page/contact">联系我们</a>
> ✅ {{ range .footerMenu.Items }}<a href="{{ .URL }}">{{ .Label }}</a>{{ end }}
> ```
>
> 唯一例外：行业特定内容（如"危化品运输资质认证"）因其不适合放入 CMS 通用菜单，可以保留硬编码。
>
> ### 🔴 #9 禁止使用 Emoji 作为 UI 图标（2026-03-27 新增）
>
> **主题模板中禁止使用 Emoji 字符（🚛📞📍✉️💰 等）替代图标。** 必须使用 FontAwesome 图标库。
>
> ```html
> <!-- ❌ 致命错误：使用 Emoji 作为 UI 图素 -->
> <div class="nav-logo-icon">🚛</div>
> <span class="nav-phone-icon">📞</span>
> <div class="footer-contact-item-icon">📍</div>
> <div class="feature-icon">💰</div>
>
> <!-- ✅ 正确：使用 FontAwesome 图标 -->
> <div class="nav-logo-icon"><i class="fa-solid fa-truck"></i></div>
> <span class="nav-phone-icon"><i class="fa-solid fa-phone"></i></span>
> <div class="footer-contact-item-icon"><i class="fa-solid fa-location-dot"></i></div>
> <div class="feature-icon"><i class="fa-solid fa-coins"></i></div>
> ```
>
> **为什么禁止 Emoji：**
> 1. **跨平台渲染不一致**：同一个 Emoji 在 Windows/macOS/Linux/iOS/Android 上外观差异巨大，部分系统显示为黑白方块
> 2. **不可控尺寸和颜色**：Emoji 不受 CSS `color`、`font-size` 控制，与主题配色系统冲突
> 3. **不够专业**：Emoji 是聊天表情，不是 UI 图标；企业站/博客使用 Emoji 显得不专业
> 4. **与 CSS 变量系统冲突**：Emoji 颜色无法通过 `var(--primary)` 等 CSS 变量统一管理
>
> **FontAwesome 使用方式：**
> ```html
> <!-- 模板 head 中引入（所有主题必须） -->
> <link rel="stylesheet" href="/fontawesome/css/all.min.css">
>
> <!-- 常用图标速查 -->
> <i class="fa-solid fa-truck"></i>          <!-- 卡车/运输 -->
> <i class="fa-solid fa-phone"></i>           <!-- 电话 -->
> <i class="fa-solid fa-location-dot"></i>    <!-- 定位/地址 -->
> <i class="fa-solid fa-envelope"></i>        <!-- 邮箱/信封 -->
> <i class="fa-solid fa-newspaper"></i>       <!-- 新闻/报纸 -->
> <i class="fa-solid fa-shield-halved"></i>   <!-- 盾牌/安全 -->
> <i class="fa-solid fa-satellite-dish"></i>  <!-- 卫星/GPS -->
> <i class="fa-solid fa-clock"></i>           <!-- 时钟/时间 -->
> <i class="fa-solid fa-bolt"></i>             <!-- 闪电/快捷 -->
> <i class="fa-solid fa-globe"></i>            <!-- 地球/全球 -->
> <i class="fa-solid fa-box"></i>              <!-- 包裹 -->
> <i class="fa-solid fa-industry"></i>         <!-- 工厂 -->
> <i class="fa-solid fa-magnifying-glass"></i><!-- 搜索 -->
> <i class="fa-solid fa-calendar"></i>         <!-- 日历 -->
> <i class="fa-solid fa-eye"></i>              <!-- 浏览/阅读量 -->
> <i class="fa-solid fa-folder"></i>           <!-- 文件夹/分类 -->
> <i class="fa-solid fa-camera"></i>           <!-- 相机 -->
> <i class="fa-solid fa-fax"></i>              <!-- 传真 -->
> <i class="fa-solid fa-comments"></i>         <!-- 评论/咨询 -->
> <i class="fa-solid fa-cloud"></i>            <!-- 云/在线 -->
> <i class="fa-solid fa-check"></i>            <!-- 对勾/认证 -->
> <i class="fa-solid fa-certificate"></i>      <!-- 证书/资质 -->
> <i class="fa-solid fa-list"></i>             <!-- 列表 -->
> <i class="fa-solid fa-chevron-left"></i>    <!-- 左箭头（翻页） -->
> <i class="fa-solid fa-chevron-right"></i>   <!-- 右箭头（翻页） -->
> ```
>
> **详细说明见 Bug 模式 [B021](references/bug-patterns.md#b021)。**
>
> ---
>
> **⛔ 文件写入安全规范（同等优先级）**
>
> **禁止用 PowerShell 修改 `theme.yml` / `.css` / 任何配置文件。**
> PowerShell 的 `WriteAllText` / `Set-Content -Encoding UTF8` 会产生双 BOM，导致 Go yaml 解析失败，主题从后台消失。
>
> - ✅ 修改现有文件 → 用 `replace_in_file` 工具（最安全）
> - ✅ 必须全量覆写 → 用 Python（`encoding='utf-8-sig'` 读 + `encoding='utf-8', newline='\n'` 写）
> - ❌ 任何 PowerShell 文件写入操作 → 一律禁止
>
> 详细说明见 [文件写入安全规范](#️-文件写入安全规范强制违反会破坏主题)。

---

## Directory Structure

每个主题的目录必须放在 `templates/themes/<slug>/` 下，结构如下：

```
templates/themes/<theme-slug>/
├── theme.yml              # 【必须】元数据 + 配置定义
├── assets/
│   ├── css/
│   │   ├── base.css       # 【必须】CSS变量 + Reset + Header/Footer 公用样式
│   │   └── *.css          # 按需拆分（home.css / article.css / list.css 等）
│   ├── js/
│   │   └── *.js           # 可选，通用脚本
│   └── images/            # 可选，主题图片资源
└── pages/
    ├── home.html          # 【必须】首页
    ├── article_show.html  # 【建议】文章详情页
    ├── article_list.html  # 【建议】文章列表页
    ├── page.html          # 【建议】独立页面（管理后台创建单页时使用）
    │
    │   以下页面均为【按需开发】，不强制要求：
    ├── category.html      # 分类页
    ├── tag.html           # 标签页
    ├── tags.html          # 标签云页
    ├── archives.html      # 归档页
    ├── categories.html    # 分类列表页
    ├── search.html        # 搜索页
    ├── 404.html           # 404 页面
    └── <自定义>.html      # 行业扩展页（作品集、产品页、时间轴等，自由命名）
```

**强制要求只有两条：**
1. `theme.yml` 必须存在，且 `slug` 与目录名完全一致
2. `home.html` 必须存在（首页是所有主题的入口）
3. `base.css` 必须存在（CSS 变量系统的根基）

**其余页面由 Phase 1 的模板规划决定**，根据行业需求和用户要求按需开发，不需要的页面不必创建。

**命名规范：**
- 目录名 = `slug` 字段值，使用 kebab-case（如 `my-blog`）
- 模板文件扩展名必须是 `.html`（不是 `.plush.html`）
- 自定义页面可自由命名，如 `portfolio.html`、`products.html`

---

### modules/ 目录：可复用模块（重要机制）

**`modules/` 是 Catter CMS 支持的跨文件模板复用机制**，这是唯一合法的 `{{ template "name" . }}` 调用方式。

```
templates/themes/<theme-slug>/
├── pages/           ← 页面模板
│   ├── home.html
│   └── ...
└── modules/         ← 可复用模块
    ├── header.html  ← 定义 {{ define "header" }} ... {{ end }}
    ├── footer.html  ← 定义 {{ define "footer" }} ... {{ end }}
    └── sidebar.html ← 定义 {{ define "sidebar" }} ... {{ end }}
```

**modules/ 文件写法（以 header.html 为例）：**

```html
{{ define "header" }}
<header class="site-header">
    <div class="header-inner">
        {{ with .site.logo }}
        <a href="/" class="site-logo">
            <img src="{{ . }}" alt="{{ $.site.title }}">
        </a>
        {{ else }}
        <a href="/" class="site-title">{{ .site.title }}</a>
        {{ end }}
        <nav>
        {{ with .headerMenu }}
            {{ range .Items }}
            <a href="{{ .URL }}">{{ .Label }}</a>
            {{ end }}
        {{ end }}
        </nav>
    </div>
</header>
{{ end }}
```

**pages/ 文件调用模块：**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>...</head>
<body>
    {{ template "header" . }}  <!-- ← 调用 modules/header.html 中的 header 块 -->
    <main>
        <!-- 页面主要内容 -->
    </main>
    {{ template "footer" . }}  <!-- ← 调用 modules/footer.html 中的 footer 块 -->
</body>
</html>
```

**⚠️ modules/ 使用规则：**
- `{{ define "name" }}` 中的 name 必须全局唯一，不能重复定义
- 模块内用 `.` 访问当前数据（与调用时传入的 `.` 相同）
- 在模块内访问顶层变量时用 `$.site`、`$.headerMenu`（`$` 代表根数据）
- `modules/` 文件本身不会被直接渲染，只有被 `{{ template }}` 调用才生效
- 如果不需要跨文件复用，不必创建 modules/，直接在页面文件内写完整结构也合法

---

## theme.yml Reference

```yaml
# =============================================
# 元数据（必需）
# =============================================
name: 主题显示名称
slug: theme-slug          # 必须与目录名一致
version: 1.0.0
author: 作者名称
description: 主题描述

# =============================================
# config - 默认配置值（扁平格式）
# 可被用户在管理后台覆盖
# =============================================
config:
  primaryColor: "#2563eb"
  primaryDark: "#0f3460"
  accentColor: "#60a5fa"
  bgPrimary: "#f8fafc"
  bgSecondary: "#ffffff"
  textPrimary: "#1e293b"
  textSecondary: "#64748b"
  borderColor: "#e5e7eb"
  sidebarEnabled: true
  heroEnabled: true
  # 数组值使用 JSON 字符串格式
  footerLinks: '[{"name":"关于","url":"/page/about"}]'

# =============================================
# configDefinition - 后台表单定义
# 决定用户在管理后台能修改哪些配置项
# ⚠️ settings.yaml 已废弃，必须在这里定义
# =============================================
configDefinition:
  - section: colors          # ⚠️ 分组标识符，不是路径！模板中访问时忽略此值
    label: 颜色设置          # 后台显示名（必须有 label）
    icon: fa-solid fa-palette
    fields:
      - name: primaryColor   # ← 模板中 getConfig "primaryColor" .config（不是 "colors.primaryColor"）
        type: color
        label: 主色调
        default: "#2563eb"
      - name: sidebarEnabled
        type: switch
        label: 显示侧边栏
        default: true
      - name: heroTitle
        type: text
        label: Hero 主标题
        show_if: "heroEnabled == true"   # 条件显示
        default: 欢迎来到
```

> ### ⚠️ 重要：config 是扁平结构，section key 不是路径
>
> `configDefinition` 中每个 `section` 的 `key`（如 `colors`、`hero`、`stats`）**仅用于后台管理界面的分组**，在模板数据访问时**完全被忽略**。
>
> **数据访问规则：**
> - `configDefinition[x].fields[y].name` = 模板中 `getConfig` 的第一个参数
> - section key 与数据访问**无关**
>
> **正确示例：**
> ```yaml
> configDefinition:
>   - section: hero
>     fields:
>       - name: heroEnabled    # ← getConfig "heroEnabled" true
>       - name: heroTitle      # ← getConfig "heroTitle" .config
>   - section: basic
>     fields:
>       - name: phone          # ← getConfig "phone" .config
> ```
>
> **错误示例（绝对不要这样写）：**
> ```html
> {{ getConfig "hero.heroEnabled" .config }}   ❌ 多了 "hero."
> {{ getConfig "basic.phone" .config }}       ❌ 多了 "basic."
> ```

**支持的字段类型：**

| type | 组件 | 用途 |
|------|------|------|
| `text` | 单行输入 | 标题、名称 |
| `textarea` | 多行文本 | 描述、备注 |
| `number` | 数字输入 | 每页数量 |
| `color` | 颜色选择器 | 主色调、背景色 |
| `switch` | 开关 | 功能开关 |
| `select` | 下拉 | 布局样式选择 |
| `radio` | 单选按钮 | 排列方式 |
| `checkbox` | 复选框 | 多项选择 |
| `image` | 图片上传 | Logo、头像 |
| `image-list` | 多图上传 | 轮播图片 |
| `datetime` | 日期时间 | 定时相关 |
| `slider` | 滑块 | 透明度、数量范围 |
| `code` | 代码编辑器 | 自定义 CSS/JS |
| `repeater` | 可视化数组 | 特性列表、链接列表 |

**show_if 语法：**
```yaml
show_if: "fieldName == value"      # 等值条件
show_if: "fieldName != value"      # 不等条件
show_if: "fieldName"               # 布尔值检查
show_if: "f1 == v1 && f2 == v2"   # AND 条件
show_if: "f1 == v1 || f2 == v2"   # OR 条件
```

---

## Template Syntax Rules

### ⛔ 绝对禁止

```html
<!-- ❌ 调用未在 modules/ 中定义的模板名（找不到定义会报错） -->
{{ template "_head" . }}       <!-- 除非 modules/中有 {{ define "_head" }} -->
{{ template "partials/nav" . }} <!-- 路径式调用不支持 -->

<!-- ✅ 正确：调用已在 modules/*.html 中 define 的模板名 -->
{{ template "header" . }}      <!-- 需要 modules/header.html 中有 {{ define "header" }} -->
{{ template "footer" . }}      <!-- 需要 modules/footer.html 中有 {{ define "footer" }} -->

<!-- ❌ dot 和字段名之间有空格（会被解析为函数调用） -->
{{ . Title }}
{{ . post.Excerpt }}

<!-- ❌ 使用 JavaScript 比较运算符 -->
{{ if .count > 0 }}   <!-- 应该用 gt -->
{{ if .count >= 0 }}  <!-- 应该用 ge -->

<!-- ❌ 引用 Post 不存在的字段 -->
{{ .Cover }}        <!-- Post 无此字段 -->
{{ .ReadTime }}     <!-- Post 无此字段 -->
{{ .Category.URL }} <!-- Category 无 URL 字段 -->
```

### ✅ 正确语法

```html
<!-- ✅ 每个页面包含完整 HTML 结构 -->
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .site.title }}</title>
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <style>:root{ --primary: {{ getConfig "primaryColor" "#2563eb" }}; }</style>
</head>
<body>
    ...
</body>
</html>

<!-- ✅ 正确字段访问（无空格） -->
{{ .Title }}
{{ .post.Excerpt }}

<!-- ✅ 比较使用函数 -->
{{ if gt .count 0 }}
{{ if ge .count 0 }}
{{ if eq .Status "published" }}

<!-- ✅ URL 手动拼接 -->
<a href="/article/{{ .Slug }}">{{ .Title }}</a>
<a href="/category/{{ .Category.Slug }}">{{ .Category.Name }}</a>
<a href="/tag/{{ .Slug }}">{{ .Name }}</a>

<!-- ✅ range 循环取索引 -->
{{ range $index, $post := .posts }}
    <span>{{ add $index 1 }}</span>
    <h3>{{ $post.Title }}</h3>
{{ end }}

<!-- ✅ in 函数参数顺序：in item slice -->
{{ if in "profile" $widgets }}...{{ end }}
```

### 配置访问两种方式

```html
<!-- 方式一（推荐）：直接访问配置对象 -->
{{ .config.primaryColor }}
{{ if .config.sidebarEnabled }}...{{ end }}

<!-- 方式二：getConfig 函数（支持默认值） -->
{{ getConfig "primaryColor" "#2563eb" }}
{{ if getConfig "sidebarEnabled" true }}...{{ end }}

<!-- ⚠️ 注意：两种方式都可用，保持全文件统一使用一种 -->
```

---

## Template Functions

```html
<!-- 配置 -->
{{ getConfig "key" "default" }}              <!-- 获取配置，支持默认值 -->

<!-- 资源路径 -->
{{ asset_path "assets/css/base.css" }}       <!-- 输出 /themes/{slug}/assets/css/base.css -->
<!-- ⚠️ asset_path 自动添加主题前缀，参数不要再带 /themes/xxx/ -->

{{ current_theme }}                          <!-- 输出当前主题 slug -->

<!-- 内容处理 -->
{{ markdownHTML .post.Content }}                <!-- Markdown → HTML（返回 template.HTML，不可再接 | safe） -->
{{ safe "<b>raw html</b>" }}                 <!-- 标记安全 HTML（谨慎使用，XSS 风险） -->
{{ default .value "默认值" }}                <!-- 提供默认值 -->

<!-- 字符串 -->
{{ lower "HELLO" }}                          <!-- hello -->
{{ upper "world" }}                          <!-- WORLD -->
{{ trim "  text  " }}                        <!-- text -->
{{ substr .Excerpt 0 150 }}                  <!-- 截取字符串 -->

<!-- 数字运算 -->
{{ add 5 3 }}  {{ sub 5 3 }}                 <!-- 8 / 2 -->
{{ div 10 3 }} {{ mod 10 3 }}                <!-- 3 / 1（整数）-->
{{ int "42" }}                               <!-- 字符串转整数 -->

<!-- 集合 -->
{{ len .posts }}                             <!-- 长度 -->
{{ first 5 .posts }}                         <!-- 前5个 -->
{{ iterate 10 }}                             <!-- 生成 0-9 序列 -->
{{ in "item" .list }}                        <!-- 包含判断（注意：item 在前，slice 在后）-->
{{ range where .posts "Status" "published" }}<!-- 过滤 -->

<!-- 日期 -->
{{ dateFormat "2006-01-02" .post.PublishedAt }}        <!-- 2026-03-26 -->
{{ dateFormat "2006年01月02日" .post.PublishedAt }}    <!-- 2026年03月26日 -->

<!-- JSON -->
{{ range (json_decode .config.footerLinks) }}
    <a href="{{ .url }}">{{ .name }}</a>
{{ end }}

<!-- 处理配置值可能是切片或 JSON 字符串的情况 -->
{{ $items := ensureSlice (getConfig "sidebarWidgets" `["profile","categories","tags"]`) }}
{{ range $items }}{{ . }}{{ end }}

<!-- 比较 -->
{{ eq .a .b }} {{ ne .a .b }}                <!-- 等于 / 不等于 -->
{{ gt .a .b }} {{ lt .a .b }}                <!-- 大于 / 小于 -->
{{ ge .a .b }} {{ le .a .b }}                <!-- 大于等于 / 小于等于 -->
{{ and .a .b }} {{ or .a .b }}               <!-- 逻辑与 / 或 -->
```

---

## Data Models

> ⚠️ 以下数据来自后端源码 `gin-app/cmd/main.go` 的实际 handler，是唯一可信来源。

### 各页面真实数据结构

```
home.html         → { site, config, pageTitle, posts[], recentPosts[], categories[], tags[] }
article_show.html → { site, config, pageTitle, post, categories[], tags[], recentPosts[] }
article_list.html → { site, config, pageTitle, posts[], categories[], tags[], pagination{} }
category.html     → { site, config, pageTitle, category, posts[], categories[], tags[] }
tag.html          → { site, config, pageTitle, tag, posts[] }
categories.html   → { site, config, pageTitle, categories[] }  ← categories 含 PostCount 字段
tags.html         → { site, config, pageTitle, tags[] }        ← tags 含 PostCount 字段
archives.html     → { site, config, pageTitle, archives{}, posts[] }
search.html       → { site, config, pageTitle, query, posts[], count }
page.html         → { site, config, pageTitle, page, recentPosts[] }
```

**⚠️ 重要区别：**
- `home.html` 有 `recentPosts`（最近5篇，无分类/标签预加载），`posts` 是最新10篇（含分类/标签）
- `article_show.html` 也有 `recentPosts`、`categories`、`tags`（可用于侧边栏）
- `category.html` 有完整的 `categories` 和 `tags` 列表（可用于侧边栏）
- `tag.html` **只有** `tag` 和 `posts`，没有 `categories`/`tags` 列表
- 分页字段名是 `pagination`（不是 `paginator`）

---

### 分页器（pagination）使用方式

分页字段名为 **`pagination`**，结构如下：

```
pagination: {
    page       int    // 当前页码（从1开始）
    pageSize   int    // 每页条数
    total      int64  // 总条数
    totalPages int64  // 总页数
}
```

**⚠️ 没有 `Pages[]` 数组，也没有 `Paginator` 对象。** 分页链接需要手动用 `iterate` 生成。

**标准分页控件写法：**

```html
{{ with .pagination }}
{{ if gt .totalPages 1 }}
<nav class="pagination">
    {{ if gt .page 1 }}
    <a href="?page={{ sub .page 1 }}" class="page-prev">上一页</a>
    {{ end }}

    {{ range $i := iterate (int .totalPages) }}
    {{ $p := add $i 1 }}
    <a href="?page={{ $p }}"
       class="page-num {{ if eq $p $.pagination.page }}active{{ end }}">
        {{ $p }}
    </a>
    {{ end }}

    {{ if lt .page .totalPages }}
    <a href="?page={{ add .page 1 }}" class="page-next">下一页</a>
    {{ end }}
</nav>
{{ end }}
{{ end }}
```

---

### 分类列表（categories.html）特有字段

`categories` 页面传入的不是普通 `Category` 数组，而是带文章数量的扩展结构：

```html
{{ range .categories }}
    {{ .Name }}          <!-- 分类名 -->
    {{ .Slug }}          <!-- 分类 slug -->
    {{ .Description }}   <!-- 分类描述 -->
    {{ .PostCount }}     <!-- ⭐ 该分类的文章数量（仅 categories.html 可用）-->
{{ end }}
```

同理，`tags.html` 的 `tags` 也有 `PostCount`：

```html
{{ range .tags }}
    {{ .Name }}       {{ .Slug }}    {{ .PostCount }}
{{ end }}
```

---

### archives.html 的数据结构

`archives` 是一个以 `"2026年03月"` 为 key 的 map，**Go 的 map 遍历是无序的**：

```html
<!-- 遍历归档（注意：Go map 无序，不保证按年月排序） -->
{{ range $yearMonth, $posts := .archives }}
<section class="archive-group">
    <h2>{{ $yearMonth }}</h2>
    {{ range $posts }}
    <a href="/article/{{ .Slug }}">{{ .Title }}</a>
    {{ end }}
</section>
{{ end }}

<!-- 如果需要有序显示，用 .posts 按年月手动处理 -->
```

---

### page.html 独立页面字段

`page` 对象与 `post` 字段不同，以下是真实字段：

```
page: { ID, Title, Slug, Content, Excerpt, Status, Template,
        MetaTitle, MetaDescription, Views, SortOrder,
        CreatedAt, UpdatedAt }

⚠️ page 无 Category、Tags、Author 字段
⚠️ page.Content 是原始文本，如需 Markdown 渲染同样使用 markdownHTML
```

**page.Template 字段**是一个重要机制：用户在管理后台创建页面时可以指定模板名。系统会优先渲染 `page-{template}.html`，找不到时回退到 `page.html`。

这意味着主题可以为不同类型页面提供专属模板：

```
pages/
├── page.html              # 通用页面（默认）
├── page-contact.html      # 联系页面（用户选 template="contact" 时使用）
├── page-gallery.html      # 图片墙（用户选 template="gallery" 时使用）
└── page-timeline.html     # 时间轴（用户选 template="timeline" 时使用）
```

---

### 关键数据字段

```
site: { title, description, url, logo, favicon, year }
  ⚠️ site 没有 pages[] 字段，导航菜单通过 headerMenu/footerMenu 访问

post: { ID, Title, Slug, Content, Excerpt, Cover, Views, Status,
        PublishedAt, ScheduledAt, CreatedAt, UpdatedAt,
        CategoryID, Category{ID,Name,Slug,Description,ParentID},
        AuthorID, Author{ID,Name,Email},
        Tags[]{ID,Name,Slug},
        MetaTitle, MetaDescription }
  ⚠️ 有 Cover 字段（封面图 URL 字符串）
  ⚠️ 无 ReadTime、URL 字段（URL 需手动拼接 /article/{{ .Slug }}）
  ⚠️ Cover 为空字符串时代表无封面，渲染前需 if 判断

category: { ID, Name, Slug, Description, ParentID, SortOrder,
            Parent{...}, Children[]{...} }
  ⚠️ ParentID 支持多级分类，Children 可遍历子分类

tag:  { ID, Name, Slug }

user: { ID, Name, Email, Role }
  ⚠️ 密码字段不暴露给模板（json:"-"）
```

---

### 菜单系统（自动注入，无需 Handler 传入）

**以下变量由 `RenderTemplate` 自动注入，每个页面均可访问：**

```
headerMenu  → *Menu  顶部菜单（location="header"）
footerMenu  → *Menu  底部菜单（location="footer"）
sidebarMenu → *Menu  侧边栏菜单（location="sidebar"）
siteName    → string 站点标题（等同于 .site.title）
themeConfig → map    主题配置（等同于 .config）
```

**Menu 对象结构：**

```
Menu: {
    ID, Name, Location, Description,
    Items: []MenuItem{
        ID, Label, URL, Icon, Target, Sort, Enabled,
        Children: []MenuItem{...}   ← 支持多级菜单
    }
}
```

**菜单渲染标准写法：**

```html
<!-- 单级菜单（headerMenu / footerMenu） -->
<nav>
{{ with .headerMenu }}
    {{ range .Items }}
    <a href="{{ .URL }}" target="{{ .Target }}">
        {{ if .Icon }}<i class="{{ .Icon }}"></i>{{ end }}
        {{ .Label }}
    </a>
    {{ end }}
{{ end }}
</nav>

<!-- 多级菜单（支持 Children 子项） -->
<nav>
{{ with .headerMenu }}
    {{ range .Items }}
    <div class="menu-item{{ if .Children }} has-children{{ end }}">
        <a href="{{ .URL }}">{{ .Label }}</a>
        {{ if .Children }}
        <ul class="sub-menu">
            {{ range .Children }}
            <li><a href="{{ .URL }}">{{ .Label }}</a></li>
            {{ end }}
        </ul>
        {{ end }}
    </div>
    {{ end }}
{{ end }}
</nav>
```

**⚠️ 菜单注意事项：**
- 菜单在管理后台配置，主题不应硬编码导航链接
- `Target` 字段值：`_self`（默认，当前页跳转）或 `_blank`（新页打开）
- `Icon` 字段为 FontAwesome class（如 `fa-solid fa-home`），可能为空
- `sidebarMenu` 用于侧边栏快捷链接，与导航菜单独立
- 菜单不存在时 `headerMenu` 为 `nil`，**必须用 `{{ with .headerMenu }}` 包裹避免 nil 错误**

---

### URL 拼接规则

```html
文章:     /article/{{ .Slug }}
分类:     /category/{{ .Slug }}
标签:     /tag/{{ .Slug }}
页面:     /page/{{ .Slug }}
文章列表: /article
归档:     /archives
标签云:   /tags
分类列表: /categories
搜索:     /search?q=...
Sitemap:  /sitemap.xml
```

---

### 前台公开 API（主题可调用）

以下 API 无需认证，可在主题的 JavaScript 中直接调用：

```
GET  /api/posts/comments/:post_id  → 获取文章评论列表
POST /api/comments                  → 提交新评论
GET  /api/search?q=关键词           → 全局搜索
```

**评论提交接口（POST /api/comments）请求体：**

```json
{
    "post_id": 1,
    "author_name": "张三",
    "author_email": "xxx@example.com",
    "author_url": "",
    "content": "评论内容"
}
```

**评论响应格式：**
```json
{
    "code": 0,
    "message": "评论提交成功，等待审核",
    "data": { ... }
}
```

**⚠️ 评论功能注意事项：**
- 评论提交后默认为 `pending`（待审核）状态，需要管理员审核后才显示
- 主题 JavaScript 调用上述 API 时应处理 `code !== 0` 的错误情况
- 系统设置 `comment_moderation` 和 `comment_approval` 控制审核策略

---

### 安全规范（模板开发必读）

#### safe 函数的风险

```html
<!-- ⚠️ safe 会绕过 HTML 自动转义，存在 XSS 风险 -->
{{ safe .userInputContent }}   <!-- ❌ 危险！用户输入内容绝不能用 safe -->
{{ safe .post.Title }}         <!-- ❌ 危险！标题字段也不能用 safe -->

<!-- ✅ safe 的唯一合法用途：渲染已经确认安全的 HTML 内容 -->
{{ markdownHTML .post.Content }}             <!-- ✅ markdownHTML 已返回 template.HTML，不需要 | safe -->
{{ markdownHTML .page.Content }}             <!-- ✅ 同上 -->
```

#### 用户数据输出规范

```html
<!-- ✅ Go 模板默认自动 HTML 转义，普通输出是安全的 -->
{{ .post.Title }}         <!-- 安全：自动转义 < > & 等字符 -->
{{ .query }}               <!-- 安全：搜索关键词自动转义 -->
{{ .site.description }}   <!-- 安全 -->

<!-- ⚠️ 在 href/src/action 等 URL 属性中仍需注意 -->
<a href="{{ .post.Slug }}">   <!-- 安全：Slug 只包含字母数字和连字符 -->
<img src="{{ .post.Cover }}"> <!-- 安全：系统存储的 URL -->
<img src="{{ .userURL }}">    <!-- ⚠️ 若来自用户输入，需验证是否为合法 URL -->
```

#### 配置值安全

```html
<!-- 配置值由管理员设置，通常安全，但 type:code 字段要小心 -->
{{ getConfig "customCSS" "" | safe }}   <!-- 管理员的自定义 CSS，允许 safe -->
{{ getConfig "customJS" "" | safe }}    <!-- 管理员的自定义 JS，允许 safe -->
{{ getConfig "heroTitle" "默认标题" }}  <!-- 普通文字不需要 safe -->
```

---

### SEO 元信息标准写法

每个页面 `<head>` 中应包含完整的 SEO 信息：

```html
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <!-- 页面标题：文章/页面用自身 title，列表页用 pageTitle + 站点名 -->
    <title>{{ if .post }}{{ .post.MetaTitle | default .post.Title }}{{ else }}{{ .pageTitle }} - {{ .site.title }}{{ end }}</title>

    <!-- 描述 -->
    <meta name="description" content="{{ if .post }}{{ .post.MetaDescription | default .post.Excerpt }}{{ else }}{{ .site.description }}{{ end }}">

    <!-- Open Graph（社交分享） -->
    <meta property="og:title" content="{{ if .post }}{{ .post.Title }}{{ else }}{{ .pageTitle }}{{ end }}">
    <meta property="og:description" content="{{ if .post }}{{ .post.Excerpt }}{{ else }}{{ .site.description }}{{ end }}">
    <meta property="og:type" content="{{ if .post }}article{{ else }}website{{ end }}">

    <!-- Canonical（page.html 和 article_show.html 才需要） -->
    {{ if .post }}<link rel="canonical" href="{{ .site.url }}/article/{{ .post.Slug }}">{{ end }}
    {{ if .page }}<link rel="canonical" href="{{ .site.url }}/page/{{ .page.Slug }}">{{ end }}

    <link rel="icon" href="{{ .site.favicon }}">
</head>
```

---

### with 语法（处理可能为空的对象）

`{{ with }}` 是处理可能为 nil/空 的对象的标准方式，比 `{{ if }}` 更简洁：

```html
<!-- ✅ 推荐：with 语法，块内 . 变为该对象 -->
{{ with .post.Category }}
<a href="/category/{{ .Slug }}">{{ .Name }}</a>
{{ end }}

<!-- ✅ with 支持 else -->
{{ with .post.Category }}
<span>{{ .Name }}</span>
{{ else }}
<span>未分类</span>
{{ end }}

<!-- ✅ 需要同时访问外部变量时，先保存到 $var -->
{{ $post := .post }}
{{ with .post.Category }}
<span>{{ .Name }} · {{ $post.Title }}</span>
{{ end }}

<!-- ❌ 不推荐（等价但啰嗦）-->
{{ if .post.Category }}
<a href="/category/{{ .post.Category.Slug }}">{{ .post.Category.Name }}</a>
{{ end }}
```

---

### repeater 字段的模板读取方式

`configDefinition` 中 `type: repeater` 的值存储为 JSON 字符串，需要 `json_decode` + `ensureSlice` 组合读取：

```yaml
# theme.yml 中定义
config:
  features: '[{"icon":"fa-solid fa-rocket","title":"快速","desc":"极速加载"}]'
configDefinition:
  - section: home
    fields:
      - name: features
        type: repeater
        label: 特性列表
        fields:
          - name: icon
            type: text
            label: 图标（FontAwesome class）
          - name: title
            type: text
            label: 标题
          - name: desc
            type: textarea
            label: 描述
```

```html
<!-- 模板中读取 repeater 字段 -->
{{ $features := ensureSlice (getConfig "features" "[]") }}
{{ range $features }}
{{ $item := json_decode . }}   <!-- 每个元素是 JSON 字符串，需再次解码 -->
<div class="feature-card">
    <i class="{{ $item.icon }}"></i>
    <h3>{{ $item.title }}</h3>
    <p>{{ $item.desc }}</p>
</div>
{{ end }}
```

**⚠️ repeater 注意事项：**
- `ensureSlice` 处理"值是 JSON 字符串"还是"值已经是 slice"两种情况
- 每个数组元素本身也是 JSON 字符串，需要 `json_decode` 再访问字段
- 字段访问用小写（`$item.icon`），与 `configDefinition` 中 `name` 一致

---

## CSS Standards

### 核心规范

```
✅ 所有样式封装在 CSS 文件中（base.css/home.css/list.css/article.css）
✅ 所有颜色使用 CSS 变量（var(--primary) 等）
✅ CSS 变量默认值在 base.css 的 :root 中定义
✅ CSS 变量名与 configDefinition 的 name 保持一致（camelCase → kebab-case）
   例：name: primaryDark → CSS: --primary-dark

❌ 禁止 style="..." 内联样式
❌ 禁止 HTML 中嵌入 <style> 块（CSS 变量注入到 :root 除外）
❌ 禁止硬编码颜色值（#xxx、rgb()、rgba()）
❌ 禁止使用 Google Fonts（国内网络问题），改用系统字体栈
❌ 禁止使用 Emoji 字符作为 UI 图标（🚛📞📍 等），必须使用 FontAwesome
```

### 标准 CSS 变量体系（base.css）

```css
:root {
    /* 颜色 */
    --primary: #2563eb;
    --primary-hover: #1d4ed8;
    --primary-dark: #0f3460;
    --primary-darker: #16213e;
    --accent: #60a5fa;
    --success: #10b981;
    --warning: #f59e0b;
    --danger: #ef4444;

    /* 背景 */
    --bg-primary: #f8fafc;
    --bg-secondary: #ffffff;

    /* 文字 */
    --text-primary: #1e293b;
    --text-secondary: #64748b;
    --text-muted: #94a3b8;

    /* 其他 */
    --border-color: #e5e7eb;
    --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    --font-mono: 'SF Mono', Monaco, 'Fira Code', monospace;
    --max-width: 1200px;
    --radius-sm: 0.25rem;
    --radius-md: 0.5rem;
    --radius-lg: 0.75rem;
    --shadow-sm: 0 1px 2px rgba(0,0,0,.05);
    --shadow-md: 0 4px 6px -1px rgba(0,0,0,.1);
}
```

### 配置驱动 CSS 变量注入模式

每个模板的 `<head>` 中，在 CSS 引用后注入用户配置覆盖：

```html
<link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
<style>
:root {
    --primary: {{ getConfig "primaryColor" "#2563eb" }};
    --primary-dark: {{ getConfig "primaryDark" "#0f3460" }};
    --accent: {{ getConfig "accentColor" "#60a5fa" }};
}
</style>
```

### 静态资源路径

```html
<!-- FontAwesome（全局可用，路径固定） -->
<link rel="stylesheet" href="/fontawesome/css/all.min.css">

<!-- 主题 CSS（使用 asset_path，参数不含主题名前缀） -->
<link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
<link rel="stylesheet" href="{{ asset_path "assets/css/home.css" }}">

<!-- ⚠️ 错误示例 -->
<link href="/themes/fontawesome/css/all.min.css">                                        <!-- ❌ 路径错误，FontAwesome 不在 themes 下 -->
<link href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">                  <!-- ❌ 路径重复：传入 /themes/my-theme/... 会被 asset_path 再添加一次前缀 → /themes/my-theme//themes/my-theme/... -->
<link href="{{ asset_path "themes/assets/css/base.css" }}">                           <!-- ❌ 少了 assets/，路径不完整 -->
<link href="{{ asset_path "/assets/css/base.css" }}">                                  <!-- ❌ 多了开头的 /，会变成 /themes/xxx//assets/... -->
```

---

## Design Philosophy

在开发主题时，遵循以下设计原则。这些原则不是建议，**是强制约束**。

### 强制规则：行业驱动设计

每个主题的设计决策必须从**目标行业的视觉文化**出发，而非从"一个通用网站"出发。

| 行业 | 视觉特征 | 禁止借鉴的通用模式 |
|------|---------|----------------|
| 美食/生活 | 暖色调、手写感、大图留白 | 蓝白科技感布局 |
| 法律/金融 | 深色庄重、无衬线+衬线混排 | 活泼的圆角卡片堆叠 |
| 科技/开发者 | 深色背景、等宽字体点缀、边框感 | 柔和渐变色块 |
| 教育/知识 | 高对比度、强排版层级 | 扁平化大色块 |
| 摄影/艺术 | 暗黑/全屏图、极少文字 | 多列卡片列表 |

### 6 条核心设计原则

1. **行业优先**：色彩、版式、模块选择必须服务于行业定位，不服从"好看但通用"的原则
2. **配置优先**：任何用户可能修改的内容都必须可配置，不硬编码文本/颜色/URL
3. **禁用通用 AI 配色**：明确禁止以下配色组合在新主题中出现：
   - 白底 + 蓝紫渐变
   - 白底 + 橙色点缀
   - 白底 + 蓝色系（`#2563eb` 蓝色已被参考模板使用，新主题不得沿用）
4. **系统字体优先**：不使用 Google Fonts，使用系统字体栈（国内访问友好）
5. **FontAwesome 图标**：所有 UI 图标必须使用 FontAwesome（`<i class="fa-solid fa-xxx"></i>`），禁止使用 Emoji 字符
6. **CSS 变量一致**：主题色彩必须通过 CSS 变量系统实现，确保用户可自定义
7. **响应式必须**：所有页面模板必须支持移动端（768px 断点）

### ⛔ 反复制规则（仅适用于开发全新主题）

以下规则**仅在创建一个全新主题时**生效。二次开发已有主题、查阅语法用法不受此限制。

**开发全新主题时，禁止：**

- **禁止将参考模板（`references/template-examples.md`）的 HTML 骨架/模块结构照搬**（看懂语法再自己写是允许的，直接复制不是）
- **禁止沿用参考模板的 CSS 类名体系**（`.post-card`、`.hero-inner`、`.sidebar` 等类名是参考模板的，新主题必须为自己的行业/风格设计类名）
- **禁止以参考模板的布局作为起点**（应以 Phase 1 Design Brief 的模块规划为起点）
- **禁止在未完成 Phase 1 Design Brief 的情况下开始写代码**

**判断标准**：如果把新主题和参考模板放在一起，一眼能看出是"换了颜色的同一套结构"，说明违反了此规则。

---

## Error Prevention Checklist

开发完成后逐项检查：

### ⚠️ 高频错误专项检查（必查，每次开发都要验证）

> 以下错误在过去开发中出现频率最高，每次主题开发完成后必须逐项确认。

#### getConfig key 路径
- [ ] **所有 `getConfig` 调用不带 section 前缀**（如 `"phone"`、`"heroEnabled"`），不带 `"basic."` / `"hero."` / `"stats."` / `"seo."` / `"contact."` 等前缀
- [ ] 搜索所有文件确认没有 `getConfig "xxx.yyy"` 形式的嵌套 key
- [ ] `theme.yml` 中 `configDefinition` 的 `section.key` 只用于后台分组，不参与任何模板数据访问

#### dateFormat 语法
- [ ] 所有 `dateFormat` 调用格式为 `{{ dateFormat "2006-01-02" .dateField }}`（格式字符串在前，日期字段在后）
- [ ] 无 `{{ dateFormat .dateField "2006-01-02" }}`（顺序反了）

#### API 调用
- [ ] JavaScript 中只调用 `POST /api/comments`（评论提交）
- [ ] 确认没有调用 `/api/contact`（不存在）
- [ ] 搜索 `fetch(` 和 `GET /api/` 确认只有 `POST /api/comments` 和 `GET /api/search`

#### CSS / 内联样式
- [ ] 无 `style="...rgba(...)..."` 内联样式
- [ ] 无 `style="...#rrggbb..."` 内联样式
- [ ] 所有颜色使用 `var(--xxx)` CSS 变量
- [ ] `style="var(--xxx)"` 是允许的（CSS 变量引用）

#### markdownHTML 语法（高危）
- [ ] 所有 `markdownHTML` 调用**不接** `| safe`（`markdownHTML` 已返回 `template.HTML`，`| safe` 会导致类型冲突和模板静默失败）
- [ ] 文章内容：`{{ markdownHTML .post.Content }}`
- [ ] 页面内容：`{{ markdownHTML .page.Content }}`
- [ ] `| safe` 仅用于管理员自定义代码：`{{ getConfig "customCSS" "" | safe }}`

#### Page 模型字段
- [ ] **未使用** `.page.Description`（Page 模型无此字段，用 `.page.Excerpt` 替代）
- [ ] Page 无 Category、Tags、Author 字段

#### 导航链接
- [ ] footer 中的通用导航链接（关于、联系、服务）使用 `{{ range .footerMenu.Items }}`
- [ ] 搜索确认没有 `href="/page/` 类硬编码链接（行业特定内容如"危化品运输资质"除外）

#### Emoji 禁止（2026-03-27 新增）
- [ ] **模板中无任何 Emoji 字符**（🚛📞📍✉️💰📅👁🔍📁📋📱 等均禁止）
- [ ] 所有图标使用 FontAwesome `<i class="fa-solid fa-xxx"></i>` 格式
- [ ] 占位符/空状态图标也使用 FontAwesome（不用 📰📰📷）
- [ ] 灯箱翻页等 UI 控件使用 FontAwesome（不用 ❮❯ ←→）

---

### 模板文件
- [ ] `home.html` 存在（必须）
- [ ] `theme.yml` 存在，`slug` 与目录名完全一致（必须）
- [ ] `base.css` 存在，包含完整 CSS 变量定义（必须）
- [ ] Phase 1 规划的所有页面均已实现（按需，非强制全集）
- [ ] 每个 `.html` 文件都有完整的 `<!DOCTYPE html>` + `<head>` + `<body>` 结构
- [ ] `{{ template "xxx" . }}` 调用的 name 必须在 `modules/*.html` 中有对应 `{{ define "xxx" }}` 定义
- [ ] modules/ 中的 define 名称全局唯一，无重复定义
- [ ] 所有 `{{ .字段名 }}` 中间无空格
- [ ] 未引用不存在的字段（ReadTime、URL 不存在；Page.Description 不存在，用 Excerpt；Cover 存在但可能为空需判断）
- [ ] URL 使用 `/article/{{ .Slug }}` 等手动拼接
- [ ] 比较使用 `eq/gt/lt/ge/le/ne` 函数（不用 >/< 符号）
- [ ] `range` 循环需要索引时使用 `$index, $item :=` 语法
- [ ] `in` 函数参数顺序正确（`in item slice`）
- [ ] 文章内容使用 `{{ markdownHTML .post.Content }}`（**禁止**再加 `| safe`，会导致类型冲突）
- [ ] 页面内容使用 `{{ markdownHTML .page.Content }}`（**禁止**再加 `| safe`）
- [ ] `safe` 函数**仅用于管理员自定义代码字段**（如 `getConfig "customCSS" "" | safe`），不用于 markdownHTML 输出

### 菜单渲染
- [ ] 导航菜单通过 `headerMenu` / `footerMenu` 访问，不硬编码导航链接
- [ ] 菜单访问必须用 `{{ with .headerMenu }}` 包裹（菜单可能为 nil）
- [ ] 菜单链接用 `{{ .URL }}`、`{{ .Label }}`、`{{ .Target }}`（大写，Go 结构体字段）
- [ ] 如有多级菜单，遍历 `{{ range .Items }}` 内的 `{{ if .Children }}` 处理子菜单

### 图片与封面
- [ ] `post.Cover` 可能为空字符串，渲染前用 `{{ if .Cover }}` 判断
- [ ] 图片资源用 `{{ asset_path "assets/images/xxx.jpg" }}` 引用（主题内部图片）
- [ ] 上传图片用 `/uploads/` 前缀路径（从数据库字段直接输出）
- [ ] Logo 从 `{{ .site.logo }}` 取，可能为空，需有 fallback 处理

### theme.yml
- [ ] `slug` 与目录名完全一致
- [ ] `configDefinition` 中每个 section 都有 `label`
- [ ] **`section.key` 只是后台分组标识，不影响模板数据访问路径**
- [ ] 每个字段 `name` 与 `config` 中的 key 完全一致
- [ ] `settings.yaml` 不存在（已废弃）
- [ ] 颜色字段 `type: color`，值格式为 `"#rrggbb"`
- [ ] 所有默认值在 `config:` 节中都有对应定义（避免 getConfig 返回 nil）
- [ ] 颜色配置有对应 CSS 变量（如 `primaryColor` → `var(--primary)`）

### ⛔ 文件写入安全规范（强制，违反会破坏主题）

> **YAML / CSS / 配置文件的写入方式决定主题能否被正常加载。编码问题是 Catter CMS 主题开发中最高频的隐性故障来源。**

#### 绝对禁止的操作

| 操作 | 后果 | 替代方案 |
|------|------|----------|
| PowerShell `[IO.File]::WriteAllText(path, content)` | 在已有 BOM 的文件上叠加 BOM，产生**双 BOM**（`EF BB BF EF BB BF`），Go yaml 库无法解析 → 主题整个被跳过，后台消失 | 用 Python 写入 |
| PowerShell `Set-Content -Encoding UTF8` | 同上，PowerShell 5.x 的 UTF8 = UTF-8 with BOM，会叠加 BOM | 用 Python 写入 |
| PowerShell `-split "\`n"` + `-join "\`n"` 处理 YAML | 丢失 Windows `\r\n`，产生换行符混乱，yaml 解析报错 | 用 Python 处理 |
| 任何读-改-写 YAML/CSS 的 PowerShell 脚本 | 同上 | 用 Python 处理 |

#### 唯一正确的文件修改方式

**优先使用编辑工具（`replace_in_file`）** — 不接触编码，最安全。

**必须生成新文件时，用 Python**：

```python
# 读取（自动处理 BOM）
content = open('theme.yml', 'r', encoding='utf-8-sig').read()

# 修改 content ...

# 写回（无 BOM，LF 换行）
with open('theme.yml', 'w', encoding='utf-8', newline='\n') as f:
    f.write(content)
```

#### 如何验证文件没有双 BOM

```powershell
$bytes = [IO.File]::ReadAllBytes("theme.yml")
$hex = ($bytes[0..5] | % { $_.ToString("X2") }) -join " "
Write-Host $hex
# 正常（无BOM）：23 20 3D ...（# = ...）
# 正常（单BOM）：EF BB BF 23 ...
# 异常（双BOM）：EF BB BF EF BB BF ...  ← 主题加载会失败
```

#### 历史案例

2026-03-26：用 PowerShell `WriteAllText` 修复 `duocai-blog/theme.yml` 乱码时，在原有单 BOM 文件上叠加了第二个 BOM，导致：
- Go yaml 报 `line 8: did not find expected key`
- 后台主题列表消失，前台继续显示旧激活主题
- 排查耗时 15 分钟，用 Python `utf-8-sig` 读取 + `utf-8` 写回修复

### CSS 文件
- [ ] `base.css` 包含完整 CSS 变量定义（`:root`）
- [ ] **无任何硬编码颜色值**：`#xxx`、`rgb()`、`rgba()`、`hsl()` 均禁止
- [ ] 无内联样式（`style="..."`）在 HTML 中
- [ ] 所有颜色来自 CSS 变量（`var(--primary)` 等）
- [ ] CSS 变量名与 configDefinition name 对应（camelCase↔kebab-case）
- [ ] 不使用 Google Fonts
- [ ] 有 `@media (max-width: 768px)` 响应式规则
- [ ] **`rgba()` 即使在 CSS 文件中也是禁止的**（如 `color: rgba(0,0,0,0.5)` → 应改为 `var(--text-muted)`）

### 资源路径
- [ ] FontAwesome 路径：`/fontawesome/css/all.min.css`（固定路径，不要用 asset_path）
- [ ] 主题 CSS/JS/图片：`{{ asset_path "assets/css/xxx.css" }}`（**正确：相对路径，不带 /themes/xxx/**）
- [ ] **asset_path 参数不以 `/themes/` 开头**（否则会路径重复变成 `/themes/xxx//themes/xxx/...`）
- [ ] **asset_path 参数不以 `/` 开头**（否则会多出一个斜杠）
- [ ] 上传文件路径：`/uploads/xxx`（不用 asset_path）

---

## Development Workflow (SDD Pattern)

遵循**规范驱动开发（Specification-Driven Development）**模式，分为四个阶段。

> ⚠️ **硬性规则：必须完成当前阶段的交付物，才能进入下一阶段。不得跳过 Phase 1。**

---

### Phase 1：设计规格（强制交付 DESIGN_BRIEF）

**在写任何代码之前，必须先产出一份完整的 Design Brief，经用户确认后方可进入 Phase 2。**

以下是 Phase 1 的完整执行步骤：

#### 步骤 1-A：行业调研与风格定位

根据用户需求，明确以下信息：

- **目标行业**：如美食博客、科技媒体、摄影作品集、教育机构、律师事务所等
- **目标用户**：浏览此网站的主要人群画像
- **视觉风格关键词**：3～5 个形容词，如「沉静、工业感、极简、有机」
- **竞品/参考**：同行业代表性网站的视觉特征（布局、颜色、排版偏好）

#### 步骤 1-B：独特色彩方案设计

- **明确禁止**：白底+渐变紫色/蓝色、白底+橙色点缀等常见 AI 通用配色
- 必须根据行业特性选择专属色彩，例如：
  - 美食类 → 暖土色系（#5c3317, #f0a500）
  - 科技媒体 → 深墨绿 + 高亮青（#0d1f1c, #00e5c0）
  - 律所 → 深靛蓝 + 烫金（#1a2744, #c9a84c）
  - 自然/摄影 → 炭黑 + 松石绿（#1c1c1e, #4a9b8e）
- 输出 5～8 个 CSS 变量色值（主色、辅色、背景、文字、边框）

#### 步骤 1-C：页面模板规划

**根据行业需求和用户要求**，规划该主题需要开发哪些页面，并详细描述每个页面的内容模块。

**不要默认列出所有可能存在的页面**，只列出这个主题真正需要的。例如：
- 纯展示型企业站 → 不需要 `archives.html`、`tags.html`
- 摄影作品集 → 不需要 `category.html`，但需要自定义的 `gallery.html`
- 知识型博客 → 可能需要 `search.html`，但不需要 `404.html`（系统有默认）

输出格式示例：

```
必须开发的页面（根据行业需求确定）：

1. home.html（首页）—— 【必须】
   模块：
   - 顶部导航栏（Logo + 菜单 + 搜索图标）
   - [行业特有] Hero 区域：全屏背景图 + 标语 + CTA 按钮
   - [行业特有] 精选文章横向滚动卡片（3列）
   - [行业特有] 分类导航（图标卡片形式，6个分类）
   - 最新文章列表（带缩略图、时间、标签）
   - 侧边栏：博主介绍 / 热门标签
   - 页脚：版权信息 + 快速导航

2. article_show.html（文章阅读页）—— 【建议】
   模块：
   - 文章头部：标题 + 分类/标签/时间/作者
   - 正文内容区（Markdown 渲染）
   - 侧边栏：目录（TOC）/ 相关文章
   - 底部：上一篇/下一篇导航

3. article_list.html（文章列表页）—— 【建议】
   模块：
   - 页面标题 + 文章数量提示
   - 列表项：标题 + 摘要 + 分类标签 + 发布时间
   - 分页控件

4. page.html（单页模板）—— 【建议，供管理后台创建普通页面使用】
   模块：
   - 页面标题
   - 富文本内容区

... （其他按需开发的页面）
```

还应规划**行业扩展页面**（如果需要）：
- 作品集/产品展示 → `portfolio.html`（图文卡片网格）
- 品牌故事/关于我们 → `timeline.html`（大事记时间轴）
- 联系信息 → `contact.html`（地图 + 表单）

#### 步骤 1-D：可配置项清单

列出所有计划放进 `theme.yml` 的 config 字段及其类型、用途。这一步必须先于写 `theme.yml`。

---

#### Phase 1 强制交付物格式

将以上内容整理为如下的 **DESIGN_BRIEF 输出**，并展示给用户确认：

```markdown
## [主题名称] - 设计规格（Design Brief）

### 定位
- 行业：[xxx]
- 目标用户：[xxx]
- 风格关键词：[A / B / C / D]

### 色彩方案
| 变量名 | 色值 | 用途 |
|--------|------|------|
| --primary | #xxxxxx | 主色调，用于按钮/链接 |
| --bg-primary | #xxxxxx | 页面背景色 |
| ... | | |

### 字体方案
- 标题：[字体栈，系统字体，无 Google Fonts]
- 正文：[字体栈]

### 页面模板清单
1. **home.html** — 包含：[模块列表]
2. **article_show.html** — 包含：[模块列表]
3. **article_list.html** — 包含：[模块列表]
4. **category.html** — 包含：[模块列表]
5. **page.html** — 包含：[模块列表]
...（其余模板）

### 配置项清单（theme.yml config）
| key | type | 用途 |
|-----|------|------|
| primaryColor | color | 主色调 |
| ... | | |

---
**⏸ 等待用户确认后，进入 Phase 2 开发。**
```

> **用户确认 Design Brief 后，才能开始写任何代码或 YAML 文件。**

---

### Phase 2：构建基础结构

1. 创建目录结构（严格按 [Directory Structure](#directory-structure)）
2. 编写 `theme.yml`（以 Phase 1 确认的可配置项清单为准）
3. 编写 `base.css`（CSS 变量以 Phase 1 色彩方案为准；写 Reset、Header、Footer 公用样式）

---

### Phase 3：开发页面模板

按 Phase 1 确认的模板清单逐页开发：

1. 严格按照各模板的**模块规划**实现，不得随意省略或简化模块
2. **参考文档（`references/`）可用于核对语法和查阅函数用法；但不得以参考模板的 HTML 结构/布局作为新主题的实现起点**
3. 每个模板的 UI 布局、色彩运用、视觉节奏必须与 Phase 1 设计方向一致
4. 开发顺序（优先级）：`home → article_show → article_list → page → 其他 Phase 1 规划页面`
5. 每个模板完成后，对照 [Error Prevention Checklist](#error-prevention-checklist) 逐项检查

---

### Phase 4：质量自检与总结

1. 用 [Error Prevention Checklist](#error-prevention-checklist) 逐项检查所有文件
2. 确认所有模板与 Phase 1 设计方案一致（色彩、布局、模块内容）
3. 输出交付清单，列出所有已创建的文件

---

## Reference Files

详细参考文档在 `references/` 目录：

- `references/template-examples.md` — 完整页面模板示例（首页/文章/分类/搜索/归档等）
- `references/theme-yml-examples.md` — 多种类型主题的 theme.yml 完整示例
- `references/css-base-template.md` — base.css 标准模板（含完整 CSS 变量和组件样式）
- `references/bug-patterns.md` — 已知 Bug 模式与修复方案速查
