---
name: catter-cms-themes-v3
description: Catter CMS 主题开发 Skill v3。开发主题时使用。包含强制规则前置检查和原子化规则集，确保 AI 不再犯重复错误。
version: 3.0.0
---

# Catter CMS 主题开发 Skill v3

## 使用说明

本 Skill 保证输出正确主题。所有主题开发任务必须遵循本 Skill。

**核心原则**：规则原子化，对错分明。写代码前必做前置检查。

---

## 强制前置检查（写任何代码前必须完成）

在开始写任何代码之前，必须在上下文中确认以下所有事项：

```
【开发前确认清单 - v3 新增】

1. 禁止蓝色：
   - 严禁使用 #2563EB、#1E40AF、#3B82F6 等蓝色系
   - 严禁使用 blue、#0000FF、rgb(37, 99, 235) 等
   - 可用色：绿色 #10B981、琥珀色 #F59E0B、深灰 #1E293B

2. 模块化开发：
   - 每个页面区块（Hero、特性、Footer 等）必须独立
   - 必须提供对应的配置开关（如 heroEnabled）
   - 可通过 theme.yml 开启/关闭，不影响其他模块

3. 数据解耦：
   - 所有内容数据必须来自 .config 或 CMS 数据模型
   - 严禁在模板中硬编码文字内容（标题、描述等）
   - 区块内容通过 getConfig 获取，可后台配置

4. 样式统一：
   - 所有 CSS 类名遵循统一命名规范
   - 禁止同一个主题内相同元素使用不同样式
   - 所有颜色使用 CSS 变量，不得硬编码

5. getConfig 路径：
   - 所有 getConfig 调用只传字段名，不带 section 前缀
   - 搜索确认没有任何 getConfig "xxx.yyy" 形式
   - 示例：getConfig "heroEnabled" true ✅ / getConfig "hero.heroEnabled" true ❌

6. markdownHTML：
   - 确认所有 markdownHTML 调用不接 | safe
   - 示例：{{ markdownHTML .post.Content }} ✅ / {{ markdownHTML .post.Content | safe }} ❌

7. asset_path：
   - 确认所有 asset_path 参数都是相对路径，不带 /themes/
   - 示例：asset_path "assets/css/base.css" ✅ / asset_path "/themes/xxx/..." ❌

8. Page 字段：
   - 确认没有使用 .page.Description（应用 .page.Excerpt）
   - 确认没有使用 .page.Category、.page.Tags

9. Emoji 检查：
   - 确认模板中没有 Emoji 字符

10. 文件写入：
    - 修改 YAML/CSS 文件用 replace_in_file 工具
    - 不用 PowerShell 写入任何文件

【确认以上全部后，开始写代码】
```

---

## 原子化规则集（写代码时对照）

### R0: 禁止蓝色系

```
所有主题严禁使用蓝色系颜色。

❌ 错误
#2563EB / #1E40AF / #3B82F6 / #60A5FA
blue / #0000FF / rgb(37, 99, 235)

✅ 正确（推荐配色）
--primary: #10B981;       /* 品牌绿 */
--primary-dark: #059669;
--accent: #F59E0B;        /* 琥珀强调 */
--bg-dark: #0F172A;       /* 深色区域 */
--text-primary: #1E293B;
--text-secondary: #64748B;
```

### R1: getConfig 路径 - 最高频错误

```
config 是扁平结构。section 只是后台分组标识，模板访问时完全忽略。

❌ 错误
{{ getConfig "hero.heroEnabled" true }}
{{ getConfig "basic.phone" .config }}
{{ getConfig "colors.primaryColor" }}

✅ 正确
{{ getConfig "heroEnabled" true }}
{{ getConfig "phone" .config }}
{{ getConfig "primaryColor" }}
```

**规则：getConfig 的第一个参数 = theme.yml 中 fields 的 name 值。**

### R2: markdownHTML 禁止 | safe

```
markdownHTML 返回 template.HTML，已标记安全。

❌ 错误（致命：模板静默失败）
{{ markdownHTML .post.Content | safe }}
{{ markdownHTML .page.Content | safe }}

✅ 正确
{{ markdownHTML .post.Content }}
{{ markdownHTML .page.Content }}
```

**唯一合法用途：**
```html
{{ getConfig "customCSS" "" | safe }}
{{ getConfig "customJs" "" | safe }}
```

### R3: asset_path 路径

```
asset_path 自动添加 /themes/{slug}/ 前缀。

❌ 错误
{{ asset_path "/themes/my-theme/assets/css/base.css" }}  → 双斜杠
{{ asset_path "/assets/css/base.css" }}                  → 多了开头的 /
{{ asset_path "themes/assets/css/base.css" }}            → 少了 assets/

✅ 正确
{{ asset_path "assets/css/base.css" }}
```

### R4: FontAwesome 路径（固定）

```
<link rel="stylesheet" href="/fontawesome/css/all.min.css">
<!-- 不是 /themes/fontawesome/，FontAwesome 不在 themes 下 -->
```

### R5: Page 模型字段

```
Page 对象没有 Description、Category、Tags 字段。

❌ 错误
{{ .page.Description }}
{{ .page.Category.Name }}
{{ .page.Tags }}

✅ 正确
{{ .page.Excerpt }}
{{ range .post.Tags }}...{{ end }}  <!-- Tags 在 post 上，不在 page 上 -->
```

### R6: Post 模型字段

```
Post 对象有 Cover（封面图），无 ReadTime、URL（需手动拼接）。

❌ 错误
{{ .post.ReadTime }}
{{ .post.Url }}

✅ 正确
{{ .post.Cover }}
<a href="/article/{{ .post.Slug }}">{{ .post.Title }}</a>
<!-- Cover 可能为空，渲染前需判断：{{ if .post.Cover }} -->
```

### R7: 禁止 Emoji 作为图标

```
❌ 错误
<span>📞</span> <span>✉️</span> <span>🚛</span>

✅ 正确
<i class="fa-solid fa-phone"></i>
<i class="fa-solid fa-envelope"></i>
<i class="fa-solid fa-truck"></i>

<!-- 常用图标 -->
fa-solid fa-home / fa-search / fa-user / fa-envelope / fa-phone
fa-solid fa-calendar / fa-eye / fa-folder / fa-tag
fa-solid fa-chevron-left / fa-chevron-right
fa-solid fa-check / fa-times / fa-bars / fa-ellipsis
fa-solid fa-truck / fa-location-dot / fa-globe
```

### R8: CSS 禁止内联和硬编码

```
❌ 错误
<div style="color: #333">
<div style="background: rgba(0,0,0,0.5)">
<div style="border: 1px solid #eee">

✅ 正确
<div class="my-element">  <!-- CSS 中定义 .my-element { color: var(--text-primary); } -->
```

### R9: 日期格式化参数顺序

```
❌ 错误
{{ dateFormat .post.PublishedAt "2006-01-02" }}

✅ 正确
{{ dateFormat "2006-01-02" .post.PublishedAt }}
```

### R10: 比较运算必须用函数

```
❌ 错误
{{ if .count > 0 }}
{{ if .count >= 10 }}

✅ 正确
{{ if gt .count 0 }}
{{ if ge .count 10 }}
{{ if eq .status "published" }}
```

### R11: range 循环取索引

```
❌ 错误
{{ range .posts }}
    {{ index $.posts . }}  <!-- . 是对象，不是索引 -->

✅ 正确
{{ range $index, $post := .posts }}
    {{ add $index 1 }}
```

### R12: in 函数参数顺序

```
❌ 错误
{{ if in $widgets "profile" }}

✅ 正确
{{ if in "profile" $widgets }}
<!-- 记忆：in "要找的东西" "在哪些里面" -->
```

### R13: 菜单必须判空

```
❌ 错误
{{ range .headerMenu.Items }}

✅ 正确
{{ if .headerMenu }}
{{ range .headerMenu.Items }}
```

### R14: 点和字段之间不能有空格

```
❌ 错误
{{ . Title }}
{{ . post.Excerpt }}

✅ 正确
{{ .Title }}
{{ .post.Excerpt }}
```

### R15: 禁止硬编码导航链接

```
❌ 错误
<a href="/page/about">关于我们</a>
<a href="/page/contact">联系我们</a>

✅ 正确
{{ range .footerMenu.Items }}
<a href="{{ .URL }}">{{ .Label }}</a>
{{ end }}
```

### R16: 文件写入安全

```
修改 theme.yml 或 CSS 文件：
✅ 用 replace_in_file 工具
❌ 不用 PowerShell 的任何写入命令（WriteAllText / Set-Content）

例外：可以用 PowerShell 执行命令，但不能写入文件。
```

### R17: 模块化开发

```
每个页面区块必须独立，可配置化。

❌ 错误
<!-- Hero 内容硬编码 -->
<div class="hero">
    <h1>Catter CMS</h1>
    <p>开源CMS</p>
</div>

✅ 正确
<!-- 通过配置开关控制 -->
{{ if getConfig "heroEnabled" true }}
<section class="hero">
    <h1>{{ getConfig "heroTitle" "Catter CMS" }}</h1>
    <p>{{ getConfig "heroSubtitle" "" }}</p>
</section>
{{ end }}

配置项：
heroEnabled: true/false  # 开关
heroTitle: "标题"        # 内容
heroSubtitle: "副标题"
heroBgImage: "背景图"
```

**必须提供的模块开关：**
- heroEnabled - Hero 区域
- featuresEnabled - 功能特性
- quickstartEnabled - 快速开始
- themeShowcaseEnabled - 主题展示
- footerEnabled - Footer

### R18: 数据解耦

```
所有内容必须来自配置或 CMS 数据模型，严禁硬编码。

❌ 错误
<a href="/page/about">关于我们</a>           <!-- 硬编码链接 -->
<p>这是一个开源项目</p>                     <!-- 硬编码文字 -->
<i class="fa-solid fa-star"></i>             <!-- 硬编码图标 -->

✅ 正确
{{ range .footerMenu.Items }}               <!-- 来自菜单系统 -->
<a href="{{ .URL }}">{{ .Label }}</a>
{{ end }}

{{ getConfig "featureDesc" "默认描述" }}     <!-- 来自配置 -->

<i class="{{ getConfig "featureIcon" "fa-solid fa-star" }}"></i>
```

**数据来源优先级：**
1. `.config` - theme.yml 中的配置
2. `.site` - 站点信息（标题、描述、logo）
3. `.headerMenu` / `.footerMenu` - 菜单系统
4. `.posts` / `.categories` / `.tags` - CMS 内容数据
5. 硬编码 - 仅允许无意义的占位符（如空状态提示）

### R19: 样式统一

```
同一主题内相同元素必须使用相同样式类。

❌ 错误
.card { background: white; }
.feature-card { background: #fff; }  /* 重复样式，不同类名 */

✅ 正确
.card { background: var(--bg-primary); }

/* 所有卡片统一使用 .card 类 */
<div class="card feature-card">
<div class="card post-card">
```

## 系统主题渲染机制

### 模板引擎

Catter CMS 使用 **Go 标准模板引擎** `html/template`，不是 Buffalo Plush。

```
❌ 错误（Plush 语法）
<%= config.heroEnabled %>
<%= for (post) in posts { %>
<%= if (post.Status == "published") { %>

✅ 正确（Go 模板语法）
{{ getConfig "heroEnabled" true }}
{{ range .posts }}
{{ if eq .Status "published" }}
```

### 主题目录结构

```
templates/themes/<slug>/
├── theme.yml           # 主题配置（元数据 + configDefinition + pageTemplates）
├── theme_thumb.png     # 缩略图（自动检测）
├── screenshot.png      # 截图（自动检测）
├── assets/             # 静态资源
│   ├── css/
│   │   ├── base.css    # 必须：CSS 变量 + Reset
│   │   └── *.css
│   └── js/
├── modules/            # 可选：可复用模板模块
│   ├── header.html
│   ├── footer.html
│   ├── sidebar.html
│   ├── comments.html
│   ├── carousel.html
│   └── lightbox.html
└── pages/              # 必须目录
    ├── home.html           # 必须：首页
    ├── article_show.html   # 必须：文章详情
    ├── article_list.html   # 必须：文章列表
    ├── category.html       # 分类页
    ├── tag.html            # 标签页
    ├── page.html           # 必须：通用单页模板
    ├── archives.html       # 归档页
    ├── search.html         # 搜索页
    └── ...
```

### 页面模板系统（单页模板）

CMS 管理后台创建页面时，用户可选择页面模板。系统会在**已激活主题的 pages/ 目录下**寻找模板文件。

**主题必须提供的页面模板：**

| 文件 | 说明 | 是否必须 |
|------|------|---------|
| `home.html` | 首页 | ✅ 必须 |
| `article_show.html` | 文章详情页 | ✅ 必须 |
| `article_list.html` | 文章列表页 | ✅ 必须 |
| `page.html` | 通用单页模板 | ✅ 必须 |
| `category.html` | 分类页 | 推荐 |
| `tag.html` | 标签页 | 推荐 |
| `archives.html` | 归档页 | 可选 |
| `search.html` | 搜索页 | 可选 |
| `tags.html` | 标签云页 | 可选 |
| `categories.html` | 分类列表页 | 可选 |

**自定义页面模板：**

主题开发者可以在 pages/ 目录下创建自定义单页模板，例如：
- `page_fullwidth.html` - 全宽页面（无侧边栏）
- `page_landing.html` - 着陆页模板
- `page_contact.html` - 联系页面模板

这些模板会在管理后台创建页面时的"页面模板"下拉列表中显示，供用户选择。

### 模块系统（可复用组件）

主题开发者可以在 `modules/` 目录下创建可复用的模板模块，通过 `{{template "模块名" .}}` 调用。

**常见模块示例：**

| 模块文件 | 说明 | 调用方式 |
|---------|------|---------|
| `header.html` | 页头/导航 | `{{template "header" .}}` |
| `footer.html` | 页脚 | `{{template "footer" .}}` |
| `sidebar.html` | 侧边栏 | `{{template "sidebar" .}}` |
| `comments.html` | 评论模块 | `{{template "comments" .}}` |
| `carousel.html` | 图片走马灯 | `{{template "carousel" .}}` |
| `lightbox.html` | 图片灯箱 | `{{template "lightbox" .}}` |
| `breadcrumb.html` | 面包屑导航 | `{{template "breadcrumb" .}}` |
| `pagination.html` | 分页组件 | `{{template "pagination" .}}` |

**模块配置化：**

每个模块都应该支持配置开关，用户可在管理后台启用/禁用：

```yaml
# theme.yml 配置项
modules:
  - section: 模块设置
    fields:
      - name: commentsEnabled
        type: switch
        label: 启用评论模块
        default: true
      - name: sidebarEnabled
        type: switch
        label: 启用侧边栏
        default: true
      - name: carouselEnabled
        type: switch
        label: 启用图片走马灯
        default: false
```

**模块模板示例：**

```html
<!-- modules/comments.html -->
{{ if getConfig "commentsEnabled" true }}
<section class="comments">
    <h3>评论</h3>
    <!-- 评论列表 -->
</section>
{{ end }}
```

**在页面中调用模块：**

```html
<!-- pages/article_show.html -->
{{template "header" .}}

<main class="article-detail">
    <article>
        <h1>{{ .post.Title }}</h1>
        {{ markdownHTML .post.Content }}
    </article>
    
    {{template "sidebar" .}}
    {{template "comments" .}}
</main>

{{template "footer" .}}
```

### 主题加载流程

```
1. 启动时加载所有主题
   └─→ 读取 theme.yml → 解析 ConfigDefinition

2. configDefinition 的 default 值初始化 config
   └─→ theme.Config[field.Name] = field.Default

3. 合并数据库配置（覆盖 YAML 默认值）
   └─→ DeepMerge(theme.Config, dbConfig)

4. 模板渲染时注入数据
   └─→ site、config、themeConfig、headerMenu、footerMenu、sidebarMenu
```

### 模板数据注入（每个页面自动获得）

| 变量 | 类型 | 说明 |
|------|------|------|
| `.site` | map | 站点信息 {title, description, url, logo, favicon, year} |
| `.config` | map | 主题配置（与 themeConfig 相同，扁平结构） |
| `.themeConfig` | map | 主题配置（与 config 相同） |
| `.headerMenu` | Menu | 顶部菜单 |
| `.footerMenu` | Menu | 底部菜单 |
| `.sidebarMenu` | Menu | 侧边栏菜单 |

**config 扁平结构：**
- theme.yml 中的嵌套 config 会被系统自动扁平化
- `getConfig "heroTitle"` 直接访问，无需考虑 section 前缀

### 可用的模板函数（完整列表）

**配置相关：**
```go
getConfig "key" "default"    // 获取配置，自动解析 JSON 字符串
asset_path "css/base.css"   // 生成 /themes/{slug}/css/base.css
current_theme               // 返回当前主题名
```

**内容相关：**
```go
markdownHTML .post.Content  // Markdown → HTML，禁止 | safe
safe "<b>raw</b>"          // 标记 HTML 安全（仅用于管理员代码）
default .value "默认"       // 提供默认值
```

**字符串相关：**
```go
lower "A"                  // 转小写
upper "a"                  // 转大写
trim " x "                 // 去空格
substr .Excerpt 0 100      // 截取字符串
split "a,b,c" ","          // 分割字符串
```

**数字相关：**
```go
add 5 3 / sub 5 3 / mul 4 3   // 加减乘
div 10 3 / mod 10 3           // 除，取模
int "42"                       // 转整数
divide100 50                   // 0-100 转 0.00-1.00
```

**集合相关：**
```go
len .posts                     // 长度
first 5 .posts                 // 取前5个
in "item" .list                // 判断存在（item 在前）
where .posts "Status" "pub"   // 过滤
iterate 10                     // 生成 0-9 序列
ensureSlice (getConfig "x")   // 确保是切片
```

**日期相关：**
```go
dateFormat "2006-01-02" .post.PublishedAt  // 格式在前，日期在后
```

**JSON 相关：**
```go
json_decode .jsonStr          // JSON 字符串解析为对象
```

**比较运算：**
```go
eq "a" "a" / ne "a" "b"      // 等于/不等于
gt 5 3 / lt 3 5              // 大于/小于
ge 5 5 / le 3 5              // 大于等于/小于等于
```

**Map/Slice 构造：**
```go
dict "key" "value" "key2" "value2"  // 构造 map
makeSlice 1 2 3                      // 构造 slice
```

### CSS 资源路径

```html
<!-- 正确：asset_path 自动添加 /themes/{slug}/ 前缀 -->
<link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">

<!-- 正确：FontAwesome 路径固定（不在 themes 下）-->
<link rel="stylesheet" href="/fontawesome/css/all.min.css">

<!-- 错误：不要手动加 /themes/ 前缀 -->
<link rel="stylesheet" href="/themes/my-theme/assets/css/base.css">
```

### 配置存储格式

**theme.yml（原始格式）：**
```yaml
config:
  heroTitle: "Catter CMS"
  features: '[{"icon":"fa-star","title":"标题"}]'
```

**数据库存储（扁平格式）：**
```json
{
  "heroTitle": "Catter CMS",
  "features": "[{\"icon\":\"fa-star\",\"title\":\"标题\"}]"
}
```

**getConfig 返回值类型（重要）：**

`getConfig` 会自动解析 JSON 字符串，**返回的是已解析的对象**，不是字符串。

| 字段类型 | getConfig 返回值类型 | 是否需要 json_decode |
|---------|---------------------|---------------------|
| text / color / number | `string` | ❌ 不需要 |
| switch | `bool` | ❌ 不需要 |
| repeater / image-list | `[]interface{}` | ❌ **不需要** |
| select / radio | `string` | ❌ 不需要 |

**❌ 错误（重复解析）：**
```go
{{ $features := getConfig "features" "[]" }}
{{ range (json_decode $features) }}  ❌ 错误：getConfig 已返回 []interface{}
  {{ .title }}
{{ end }}
```

**✅ 正确（直接使用）：**
```go
{{ $features := getConfig "features" "[]" }}
{{ range $features }}  ✅ 正确：直接 range
  {{ .title }}
{{ end }}
```

**ensureSlice 配合（防御性编程）：**
```go
{{ $features := getConfig "features" "[]" }}
{{ range ensureSlice $features }}  ✅ 安全：确保是 slice
  {{ .title }}
{{ end }}
```

**重要规则：任何 getConfig 获取的 repeater/image-list 字段，禁止使用 json_decode。**

### 菜单数据结构

```go
type Menu struct {
  ID        uint
  Name      string
  Location  string  // "header" / "footer" / "sidebar"
  Items     []MenuItem
}

type MenuItem struct {
  Label     string
  URL       string
  Icon      string
  Target    string  // "_blank" 等
  Children  []MenuItem
}
```

**使用前必须判空：**
```go
{{ if .headerMenu }}
{{ range .headerMenu.Items }}
  <a href="{{ .URL }}" {{ if .Target }}target="{{ .Target }}"{{ end }}>{{ .Label }}</a>
{{ end }}
{{ end }}
```

---

## 模板函数速查

```
【配置】
getConfig "key" "default"          # 获取配置，无 section 前缀
asset_path "assets/css/xxx.css"     # 生成资源路径，相对路径
current_theme                       # 当前主题名

【内容】
markdownHTML .post.Content          # Markdown → HTML，禁止 | safe
safe "<b>raw</b>"                  # 仅用于管理员自定义代码
default .value "默认"               # 提供默认值

【字符串】
lower "A" / upper "a"             # 大小写
trim " x "                         # 去空格
substr .Excerpt 0 100             # 截取字符串

【数字】
add 5 3 / sub 5 3 / mul 4 3      # 加减乘
div 10 3 / mod 10 3              # 除，取模
int "42"                          # 转整数

【集合】
len .posts                         # 长度
first 5 .posts                    # 取前5个
in "item" .list                   # 判断存在（item 在前）
where .posts "Status" "pub"       # 过滤
iterate 10                        # 生成 0-9 序列
ensureSlice (getConfig "x" "[]")  # 确保是切片

【日期】
dateFormat "2006-01-02" .post.PublishedAt  # 格式在前，日期在后

【JSON】
json_decode .jsonStr              # JSON 字符串解析
```

---

## 数据模型速查

```
home.html         → { site, config, pageTitle, posts[], recentPosts[], categories[], tags[] }
article_show.html → { site, config, pageTitle, post, categories[], tags[], recentPosts[] }
article_list.html → { site, config, pageTitle, posts[], categories[], tags[], pagination{} }
category.html     → { site, config, pageTitle, category, posts[], categories[], tags[] }
tag.html          → { site, config, pageTitle, tag, posts[] }
page.html         → { site, config, pageTitle, page, recentPosts[] }
search.html       → { site, config, pageTitle, query, posts[], count }
archives.html     → { site, config, pageTitle, archives{}, posts[] }

site: { title, description, url, logo, favicon, year }
post: { ID, Title, Slug, Content, Excerpt, Cover, Views, Status,
        PublishedAt, Category{ID,Name,Slug}, Author{ID,Name},
        Tags[]{ID,Name,Slug} }
page: { ID, Title, Slug, Content, Excerpt, Status, Template,
        MetaTitle, MetaDescription, Views, SortOrder, CreatedAt, UpdatedAt }
        ⚠️ 无 Description，用 Excerpt；无 Category/Tags
pagination: { page, pageSize, total, totalPages }
        ⚠️ 不是 Paginator，无 Pages[]

menu: headerMenu / footerMenu / sidebarMenu
     Items: [{ Label, URL, Icon, Target, Children[] }]
     ⚠️ 必须用 {{ if .headerMenu }} 包裹判空
```

---

## URL 规则

```
文章      → /article/{{ .Slug }}
分类      → /category/{{ .Slug }}
标签      → /tag/{{ .Slug }}
页面      → /page/{{ .Slug }}
文章列表  → /article
归档      → /archives
标签云    → /tags
分类列表  → /categories
搜索      → /search?q=关键词
```

---

## theme.yml 结构

```yaml
name: 主题名称
slug: theme-slug        # 必须与目录名一致
version: 1.0.0
author: 作者
description: 描述

config:
  primaryColor: "#2563eb"
  heroEnabled: true
  # 数组值用 JSON 字符串
  features: '[{"icon":"fa-solid fa-star","title":"标题","desc":"描述"}]'

configDefinition:
  - section: basic      # ⚠️ 只是后台分组标识，不影响模板访问
    label: 基础设置
    fields:
      - name: primaryColor  # ← getConfig "primaryColor"
        type: color
        label: 主色调
        default: "#2563eb"
      - name: heroEnabled  # ← getConfig "heroEnabled"
        type: switch
        label: 显示 Hero
        default: true
        show_if: "heroEnabled == true"
```

**支持的字段类型：**
```
text / textarea / number / color / switch
select / radio / checkbox
image / image-list
datetime / slider / code / repeater
```

---

## 目录结构

```
templates/themes/<slug>/
├── theme.yml
├── assets/
│   ├── css/
│   │   ├── base.css    # 必须：CSS 变量 + Reset
│   │   └── *.css
│   └── js/
└── pages/
    ├── home.html           # 必须
    ├── article_show.html
    ├── article_list.html
    ├── category.html
    ├── tag.html
    ├── page.html
    └── ...
```

---

## 开发流程

```
Phase 1: 确认需求和设计方案
         - 确认配色方案（无蓝色）
         - 确认模块列表
         - 确认数据来源

Phase 2: 创建 theme.yml + base.css
         - 定义所有模块配置项
         - 定义 CSS 变量（颜色、间距、字体）
         - 定义统一组件类

Phase 3: 开发页面模板（home → article_show → 其他）
         - 每个区块独立开发
         - 使用配置开关控制显示
         - 统一组件样式

Phase 4: 对照规则集自检
         - R0 蓝色检查
         - R17 模块化检查
         - R18 数据解耦检查
         - R19 样式统一检查
```

---

## 完成后自检清单

开发完成后，逐项确认：

```
【P0 致命错误】
□ 没有 getConfig "xxx.yyy" 形式（section 前缀）
□ 没有 markdownHTML | safe
□ 没有 asset_path "/themes/..."

【P1 常见错误】
□ 没有 .page.Description（应用 .page.Excerpt）
□ 没有 Emoji
□ 没有内联 style="color:#xxx"
□ 菜单用 {{ if .headerMenu }} 包裹

【P2 语法错误】
□ 比较用 eq/gt/lt/ge/le/ne
□ range 取索引用 $index, $item
□ dateFormat 格式在前
□ in 函数参数顺序正确
□ 无 . 和字段之间的空格

【文件】
□ 用 replace_in_file 修改文件
□ slug 与目录名一致
```
