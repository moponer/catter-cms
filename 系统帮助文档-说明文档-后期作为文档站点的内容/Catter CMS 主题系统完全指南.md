# Catter CMS 主题系统完全指南

## 一、Halo 风格主题系统概述

### 1.1 什么是主题系统？

主题系统是 CMS（内容管理系统）的核心功能之一，它允许用户：
- **切换外观**：一键更换网站的整体风格和布局
- **自定义配置**：通过后台界面修改颜色、字体、布局等
- **扩展功能**：添加新的小工具、页面模板、自定义字段

### 1.2 Halo 主题系统的特点

Halo 是一个流行的开源博客/CMS系统，其主题系统设计非常优秀：

| 特性 | 说明 |
|------|------|
| **YAML配置驱动** | 主题的所有配置都在 `theme.yml` 中定义 |
| **模块化配置** | 配置按功能分组（外观、SEO、布局等） |
| **表单自动生成** | 根据配置定义自动生成后台配置表单 |
| **多页面模板** | 支持首页、文章页、分类页、标签页等不同模板 |
| **小工具系统** | 侧边栏可拖拽配置各种小工具 |
| **条件显示** | 配置项可根据其他配置的值显示/隐藏 |

### 1.3 Catter CMS 的实现

Catter CMS 参考 Halo 的设计，实现了类似的主题系统：

```
┌─────────────────────────────────────────────────────────────┐
│                    主题系统架构                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │   theme.yml  │───▶│  后台配置界面 │───▶│  数据库保存  │  │
│  │  (配置定义)   │    │  (动态表单)   │    │  (持久化)    │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│          │                                            │     │
│          ▼                                            ▼     │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Plush 模板引擎渲染                       │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │ home.plush  │  │blog_show.   │  │category.    │  │  │
│  │  │   .html     │  │  plush.html │  │ plush.html  │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
│                            │                                │
│                            ▼                                │
│                    ┌──────────────┐                        │
│                    │   前台页面    │                        │
│                    │  (HTML输出)   │                        │
│                    └──────────────┘                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 二、主题目录结构

### 2.1 最小主题结构

一个最基本的主题只需要两个文件：

```
my-theme/
├── theme.yml          # 主题配置文件（必需）
└── pages/
    └── home.plush.html    # 首页模板（必需）
```

### 2.2 完整主题结构

```
my-theme/
├── theme.yml              # 主题配置定义
├── README.md              # 主题说明文档
├── LICENSE                # 许可证
│
├── pages/                 # 页面模板目录
│   ├── home.plush.html        # 首页
│   ├── blog_list.plush.html   # 文章列表页
│   ├── blog_show.plush.html   # 文章详情页
│   ├── categories.plush.html  # 分类列表页
│   ├── category.plush.html    # 分类文章列表页
│   ├── tags.plush.html        # 标签列表页
│   ├── tag.plush.html         # 标签文章列表页
│   ├── archives.plush.html    # 归档页
│   ├── search.plush.html      # 搜索页
│   └── page.plush.html        # 独立页面
│
├── assets/                # 静态资源目录
│   ├── css/
│   │   ├── style.css
│   │   └── responsive.css
│   ├── js/
│   │   └── main.js
│   ├── images/
│   │   ├── logo.png
│   │   └── banner.jpg
│   └── fonts/
│       └── iconfont.woff
│
├── partials/              # 可复用模板片段
│   ├── header.plush.html
│   ├── footer.plush.html
│   ├── sidebar.plush.html
│   └── pagination.plush.html
│
└── screenshots/           # 主题预览图
    ├── screenshot1.jpg
    └── thumbnail.png
```

---

## 三、theme.yml 配置详解

### 3.1 基础信息

```yaml
# 主题基本信息
name: "My Theme"                    # 主题名称（显示用）
version: "1.0.0"                    # 版本号
apiVersion: "2.0"                   # API版本
slug: "my-theme"                    # 主题标识（目录名）
description: "一个简洁的博客主题"    # 描述
author: "Your Name"                 # 作者
authorEmail: "you@example.com"      # 作者邮箱
authorUrl: "https://your-site.com"  # 作者网站
license: "MIT"                      # 许可证
tags: ["blog", "minimal", "responsive"]  # 标签

# 预览图
screenshots:
  - "screenshots/home.jpg"
  - "screenshots/post.jpg"
thumbnail: "screenshots/thumb.png"
```

### 3.2 配置定义（configDefinition）

这是最重要的部分，定义了后台配置表单的结构：

```yaml
configDefinition:
  # 配置分组
  - section: "外观"                   # 分组名称
    icon: "palette"                   # 分组图标
    description: "配置颜色和字体"      # 分组描述
    fields:                           # 字段列表
      
      # 文本输入
      - name: "site.title"            # 配置键名（支持点号路径）
        label: "网站标题"              # 显示标签
        type: "text"                  # 字段类型
        default: "我的博客"            # 默认值
        placeholder: "输入网站标题"     # 占位符
        required: true                # 是否必填
      
      # 多行文本
      - name: "site.description"
        label: "网站描述"
        type: "textarea"
        rows: 3                       # 行数
        default: ""
      
      # 数字输入
      - name: "posts.perPage"
        label: "每页文章数"
        type: "number"
        default: 10
        min: 1                        # 最小值
        max: 50                       # 最大值
        step: 1                       # 步长
      
      # 颜色选择器
      - name: "colors.primary"
        label: "主色调"
        type: "color"
        default: "#3B82F6"
        showAlpha: true               # 是否支持透明度
      
      # 开关
      - name: "features.showSidebar"
        label: "显示侧边栏"
        type: "switch"
        default: true
      
      # 下拉选择
      - name: "layout.sidebarPosition"
        label: "侧边栏位置"
        type: "select"
        options:
          - label: "左侧"
            value: "left"
          - label: "右侧"
            value: "right"
          - label: "不显示"
            value: "none"
        default: "right"
      
      # 多选
      - name: "sidebar.widgets"
        label: "侧边栏小工具"
        type: "multi-select"
        options:
          - label: "搜索"
            value: "search"
          - label: "最近文章"
            value: "recent"
          - label: "分类"
            value: "categories"
          - label: "标签云"
            value: "tags"
        default: ["search", "recent"]
      
      # 单选按钮
      - name: "posts.layout"
        label: "文章列表布局"
        type: "radio"
        options:
          - label: "列表"
            value: "list"
          - label: "网格"
            value: "grid"
          - label: "瀑布流"
            value: "masonry"
        default: "list"
      
      # 复选框
      - name: "posts.showMeta"
        label: "显示元信息"
        type: "checkbox"
        options:
          - label: "作者"
            value: "author"
          - label: "日期"
            value: "date"
          - label: "分类"
            value: "category"
          - label: "标签"
            value: "tags"
        default: ["date", "category"]
      
      # 图片上传
      - name: "header.logo"
        label: "网站Logo"
        type: "image"
        description: "建议尺寸 200x60"
      
      # 日期时间
      - name: "campaign.startDate"
        label: "活动开始时间"
        type: "datetime"
      
      # 日期范围
      - name: "campaign.dateRange"
        label: "活动日期"
        type: "daterange"
      
      # 滑块
      - name: "layout.contentWidth"
        label: "内容宽度"
        type: "slider"
        min: 800
        max: 1400
        step: 50
        default: 1200
      
      # 代码编辑器
      - name: "custom.css"
        label: "自定义CSS"
        type: "code"
        language: "css"               # 代码语言
        rows: 10
        description: "会覆盖主题默认样式"
      
      # 标签输入
      - name: "seo.keywords"
        label: "关键词"
        type: "tags"
        placeholder: "输入后按回车"
        default: ["博客", "技术"]
      
      # 条件显示
      - name: "sidebar.customHtml"
        label: "自定义HTML"
        type: "textarea"
        show_if: "sidebar.widgets contains 'custom'"  # 条件表达式
```

### 3.3 支持的字段类型

| 类型 | 说明 | 适用场景 |
|------|------|----------|
| `text` | 单行文本 | 标题、URL、简短文字 |
| `textarea` | 多行文本 | 描述、简介 |
| `number` | 数字 | 数量、尺寸 |
| `color` | 颜色选择器 | 主题色、背景色 |
| `switch` | 开关 | 启用/禁用功能 |
| `select` | 下拉选择 | 单选场景 |
| `multi-select` | 多选下拉 | 多选场景 |
| `radio` | 单选按钮 | 少量选项 |
| `checkbox` | 复选框 | 多选场景 |
| `image` | 图片上传 | Logo、Banner |
| `datetime` | 日期时间 | 活动时间 |
| `daterange` | 日期范围 | 活动周期 |
| `slider` | 滑块 | 数值范围 |
| `code` | 代码编辑器 | CSS/JS/HTML |
| `tags` | 标签输入 | 关键词 |

### 3.4 条件显示（show_if）

条件显示允许根据其他字段的值动态显示/隐藏配置项：

```yaml
fields:
  # 启用Hero模块
  - name: "hero.enabled"
    label: "启用Hero"
    type: "switch"
    default: true
  
  # 只有启用Hero时才显示标题配置
  - name: "hero.title"
    label: "Hero标题"
    type: "text"
    show_if: "hero.enabled"           # 简单条件
  
  # 背景类型选择
  - name: "hero.backgroundType"
    label: "背景类型"
    type: "select"
    options:
      - { label: "颜色", value: "color" }
      - { label: "图片", value: "image" }
  
  # 只有选择颜色背景时才显示颜色选择器
  - name: "hero.backgroundColor"
    label: "背景颜色"
    type: "color"
    show_if: "hero.enabled && hero.backgroundType == 'color'"
  
  # 只有选择图片背景时才显示图片上传
  - name: "hero.backgroundImage"
    label: "背景图片"
    type: "image"
    show_if: "hero.enabled && hero.backgroundType == 'image'"
```

支持的运算符：
- `==` / `!=`：等于/不等于
- `&&` / `||`：与/或
- `contains`：包含（用于数组）

---

## 四、模板开发

### 4.1 Plush 模板语法

Catter CMS 使用 Plush 模板引擎（Go语言），语法类似 Ruby ERB：

#### 基础语法

```html
<!-- 输出变量 -->
<%= variable %>
<%= user.Name %>
<%= config["site.title"] %>

<!-- 条件判断 -->
<%= if (condition) { %>
  <p>条件为真时显示</p>
<%= } %>

<%= if (user.IsAdmin) { %>
  <p>管理员界面</p>
<%= } else { %>
  <p>普通用户界面</p>
<%= } %>

<!-- 循环 -->
<%= for (item) in items { %>
  <li><%= item.Name %></li>
<%= } %>

<!-- 带索引的循环 -->
<%= for (index, item) in items { %>
  <li><%= index + 1 %>. <%= item.Name %></li>
<%= } %>

<!-- 变量赋值 -->
<%= let title = config["site.title"] %>
<h1><%= title %></h1>
```

#### 常用辅助函数

```html
<!-- 主题资源路径 -->
<link rel="stylesheet" href="<%= asset_path("css/style.css") %>">
<img src="<%= asset_path("images/logo.png") %>">

<!-- 当前主题名称 -->
<p>当前主题: <%= current_theme() %></p>

<!-- 检查小工具是否启用 -->
<%= if (hasWidget("search")) { %>
  <div class="widget-search">...</div>
<%= } %>
```

### 4.2 模板中可用的数据

服务端渲染时会传递以下数据到模板：

```go
// 基础数据
site          map[string]string   // 网站信息（title, description等）
config        map[string]interface{}  // 主题配置

// 内容数据
posts         []Post              // 文章列表
post          Post                // 单篇文章
categories    []Category          // 分类列表
category      Category            // 单个分类
tags          []Tag               // 标签列表
tag           Tag                 // 单个标签
pages         []Page              // 页面列表

// 分页数据
pagination    Pagination          // 分页信息
  .CurrentPage   int
  .TotalPages    int
  .HasNext       bool
  .HasPrev       bool
  .NextURL       string
  .PrevURL       string
```

### 4.3 首页模板示例

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title><%= site["title"] %></title>
    <meta name="description" content="<%= site["description"] %>">
    
    <!-- 使用主题配置的颜色 -->
    <style>
        :root {
            --primary-color: <%= config["colors.primary"] || "#3B82F6" %>;
            --text-color: <%= config["colors.text"] || "#1f2937" %>;
            --bg-color: <%= config["colors.background"] || "#ffffff" %>;
        }
        
        body {
            font-family: <%= config["typography.font"] || "system-ui" %>;
            color: var(--text-color);
            background: var(--bg-color);
        }
        
        /* 自定义CSS注入 */
        <%= config["custom.css"] %>
    </style>
</head>
<body>
    <!-- 导航 -->
    <nav class="navbar">
        <div class="logo">
            <%= if (config["header.logo"]) { %>
                <img src="<%= config["header.logo"] %>" alt="Logo">
            <%= } else { %>
                <%= site["title"] %>
            <%= } %>
        </div>
        <div class="nav-links">
            <a href="/">首页</a>
            <a href="/blog">文章</a>
            <a href="/categories">分类</a>
            <%= for (item) in menus["header"] { %>
                <a href="<%= item.URL %>"><%= item.Name %></a>
            <%= } %>
        </div>
    </nav>
    
    <!-- Hero区域（根据配置显示/隐藏） -->
    <%= if (config["hero.enabled"]) { %>
    <section class="hero" style="background: <%= config["hero.backgroundColor"] %>">
        <h1><%= config["hero.title"] %></h1>
        <p><%= config["hero.subtitle"] %></p>
    </section>
    <%= } %>
    
    <!-- 主内容区 -->
    <div class="container <%= config["layout.sidebarPosition"] %>">
        <main class="content">
            <%= for (post) in posts { %>
            <article class="post-card">
                <%= if (config["posts.showThumbnail"] && post.Thumbnail) { %>
                <img src="<%= post.Thumbnail %>" alt="<%= post.Title %>">
                <%= } %>
                
                <h2><a href="/blog/<%= post.Slug %>"><%= post.Title %></a></h2>
                
                <div class="meta">
                    <%= if (config["posts.showDate"]) { %>
                    <time><%= post.PublishedAt %></time>
                    <%= } %>
                    <%= if (config["posts.showCategory"] && post.Category) { %>
                    <a href="/category/<%= post.Category.Slug %>">
                        <%= post.Category.Name %>
                    </a>
                    <%= } %>
                </div>
                
                <%= if (config["posts.showExcerpt"]) { %>
                <p class="excerpt"><%= post.Excerpt %></p>
                <%= } %>
            </article>
            <%= } %>
            
            <!-- 分页 -->
            <%= if (pagination.TotalPages > 1) { %>
            <div class="pagination">
                <%= if (pagination.HasPrev) { %>
                <a href="<%= pagination.PrevURL %>">上一页</a>
                <%= } %>
                <span><%= pagination.CurrentPage %> / <%= pagination.TotalPages %></span>
                <%= if (pagination.HasNext) { %>
                <a href="<%= pagination.NextURL %>">下一页</a>
                <%= } %>
            </div>
            <%= } %>
        </main>
        
        <!-- 侧边栏（根据配置显示/隐藏） -->
        <%= if (config["features.showSidebar"]) { %>
        <aside class="sidebar">
            <%= if (hasWidget("search")) { %>
            <div class="widget">
                <h4>搜索</h4>
                <form action="/search" method="get">
                    <input type="text" name="q" placeholder="搜索...">
                </form>
            </div>
            <%= } %>
            
            <%= if (hasWidget("categories")) { %>
            <div class="widget">
                <h4>分类</h4>
                <ul>
                    <%= for (cat) in categories { %>
                    <li>
                        <a href="/category/<%= cat.Slug %>">
                            <%= cat.Name %>
                            <%= if (cat.PostCount > 0) { %>
                            <span class="count">(<%= cat.PostCount %>)</span>
                            <%= } %>
                        </a>
                    </li>
                    <%= } %>
                </ul>
            </div>
            <%= } %>
            
            <%= if (hasWidget("tags")) { %>
            <div class="widget">
                <h4>标签云</h4>
                <div class="tag-cloud">
                    <%= for (tag) in tags { %>
                    <a href="/tag/<%= tag.Slug %>" class="tag">
                        <%= tag.Name %>
                    </a>
                    <%= } %>
                </div>
            </div>
            <%= } %>
        </aside>
        <%= } %>
    </div>
    
    <!-- 页脚 -->
    <footer>
        <p>&copy; <%= site["title"] %>. Powered by Catter CMS</p>
    </footer>
    
    <!-- 自定义JS -->
    <%= if (config["custom.js"]) { %>
    <script>
        <%= config["custom.js"] %>
    </script>
    <%= } %>
</body>
</html>
```

### 4.4 文章详情页模板

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title><%= post.Title %> - <%= site["title"] %></title>
    <meta name="description" content="<%= post.Excerpt %>">
</head>
<body>
    <article class="post-detail">
        <!-- 文章头部 -->
        <header class="post-header">
            <h1><%= post.Title %></h1>
            <div class="post-meta">
                <%= if (post.Category) { %>
                <a href="/category/<%= post.Category.Slug %>">
                    <%= post.Category.Name %>
                </a>
                <%= } %>
                <time><%= post.PublishedAt %></time>
                <%= if (post.Author) { %>
                <span class="author"><%= post.Author.Name %></span>
                <%= } %>
            </div>
        </header>
        
        <!-- 文章封面 -->
        <%= if (post.CoverImage) { %>
        <img src="<%= post.CoverImage %>" alt="<%= post.Title %>" class="post-cover">
        <%= } %>
        
        <!-- 文章内容 -->
        <div class="post-content">
            <%= post.ContentHTML %>
        </div>
        
        <!-- 文章标签 -->
        <%= if (post.Tags && len(post.Tags) > 0) { %>
        <div class="post-tags">
            <%= for (tag) in post.Tags { %>
            <a href="/tag/<%= tag.Slug %>" class="tag">#<%= tag.Name %></a>
            <%= } %>
        </div>
        <%= } %>
        
        <!-- 上一篇/下一篇 -->
        <nav class="post-nav">
            <%= if (post.PrevPost) { %>
            <a href="/blog/<%= post.PrevPost.Slug %>" class="prev">
                ← <%= post.PrevPost.Title %>
            </a>
            <%= } %>
            <%= if (post.NextPost) { %>
            <a href="/blog/<%= post.NextPost.Slug %>" class="next">
                <%= post.NextPost.Title %> →
            </a>
            <%= } %>
        </nav>
    </article>
</body>
</html>
```

---

## 五、服务端渲染流程

### 5.1 渲染流程图

```
用户请求 ──▶ 路由匹配 ──▶ 处理器函数 ──▶ 查询数据 ──▶ 渲染模板 ──▶ 返回HTML
                                              │
                                              ▼
                                    ┌─────────────────────┐
                                    │   Plush 模板引擎     │
                                    │                     │
                                    │  1. 读取模板文件      │
                                    │  2. 创建上下文        │
                                    │  3. 注入数据          │
                                    │  4. 执行渲染          │
                                    │  5. 返回HTML          │
                                    └─────────────────────┘
```

### 5.2 处理器函数示例

```go
// 首页处理器
func HomeHandler(c *gin.Context) {
    // 1. 获取查询参数
    page := c.DefaultQuery("page", "1")
    pageNum, _ := strconv.Atoi(page)
    
    // 2. 查询数据
    posts := models.GetPublishedPosts(pageNum, 10)
    categories := models.GetAllCategories()
    tags := models.GetAllTags()
    
    // 3. 准备模板数据
    data := map[string]interface{}{
        "site": map[string]string{
            "title":       models.GetSetting("site.title", "My Site"),
            "description": models.GetSetting("site.description", ""),
        },
        "posts":      posts,
        "categories": categories,
        "tags":       tags,
        "pagination": models.NewPagination(pageNum, 10, models.GetPostCount()),
    }
    
    // 4. 渲染模板
    html, err := services.RenderTemplate("home", data)
    if err != nil {
        c.String(500, "Template error: "+err.Error())
        return
    }
    
    // 5. 返回HTML
    c.Data(200, "text/html; charset=utf-8", []byte(html))
}
```

### 5.3 主题服务核心代码

```go
// RenderTemplate 渲染模板
func RenderTemplate(templateName string, data map[string]interface{}) (template.HTML, error) {
    theme := GetActiveTheme()
    if theme == nil {
        return "", fmt.Errorf("no active theme")
    }
    
    // 构建模板路径
    templatePath := filepath.Join(theme.Path, "pages", templateName+".plush.html")
    
    // 读取模板内容
    content, err := os.ReadFile(templatePath)
    if err != nil {
        return "", err
    }
    
    // 创建 Plush 上下文
    ctx := plush.NewContext()
    
    // 注入数据
    for key, value := range data {
        ctx.Set(key, value)
    }
    
    // 注入主题配置
    ctx.Set("config", theme.Config)
    
    // 注入辅助函数
    ctx.Set("asset_path", func(path string) string {
        return "/themes/" + theme.Name + "/" + path
    })
    
    // 执行渲染
    rendered, err := plush.Render(string(content), ctx)
    if err != nil {
        return "", err
    }
    
    return template.HTML(rendered), nil
}
```

---

## 六、开发新主题的步骤

### 6.1 快速开始

```bash
# 1. 创建主题目录
mkdir -p templates/themes/my-theme/pages

# 2. 创建 theme.yml
cat > templates/themes/my-theme/theme.yml << 'EOF'
name: "My Theme"
version: "1.0.0"
description: "我的自定义主题"
author: "Your Name"

configDefinition:
  - section: "基础设置"
    icon: "setting"
    fields:
      - name: "site.title"
        label: "网站标题"
        type: "text"
        default: "我的博客"
EOF

# 3. 创建首页模板
cat > templates/themes/my-theme/pages/home.plush.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title><%= site["title"] %></title>
</head>
<body>
    <h1><%= site["title"] %></h1>
    <%= for (post) in posts { %>
    <article>
        <h2><a href="/blog/<%= post.Slug %>"><%= post.Title %></a></h2>
    </article>
    <%= } %>
</body>
</html>
EOF

# 4. 重启服务，在后台激活主题
```

### 6.2 开发检查清单

- [ ] 创建 `theme.yml` 并定义基础信息和配置
- [ ] 创建 `pages/home.plush.html` 首页模板
- [ ] 创建 `pages/blog_show.plush.html` 文章详情页
- [ ] 创建 `pages/category.plush.html` 分类页
- [ ] 添加 `assets/css/style.css` 样式文件
- [ ] 测试所有配置项是否生效
- [ ] 测试响应式布局
- [ ] 添加预览图 `screenshots/thumbnail.png`

---

## 七、最佳实践

### 7.1 配置设计原则

1. **合理的默认值**：确保主题开箱即用
2. **渐进式配置**：基础配置简单，高级配置可选
3. **条件显示**：相关配置项使用 `show_if` 关联
4. **命名规范**：使用点号分隔的层级命名，如 `layout.sidebar.position`

### 7.2 模板开发原则

1. **防御性编程**：所有配置项都要有默认值
   ```html
   <%= config["title"] || "默认标题" %>
   ```

2. **条件渲染**：根据配置显示/隐藏元素
   ```html
   <%= if (config["showSidebar"]) { %>
     <aside>...</aside>
   <%= } %>
   ```

3. **响应式设计**：使用媒体查询适配移动端
   ```css
   @media (max-width: 768px) {
     .sidebar { display: none; }
   }
   ```

4. **性能优化**：图片懒加载、CSS/JS压缩

### 7.3 目录结构建议

```
my-theme/
├── theme.yml              # 配置定义
├── pages/                 # 页面模板
│   ├── home.plush.html
│   ├── blog_show.plush.html
│   └── ...
├── assets/                # 静态资源
│   ├── css/
│   ├── js/
│   └── images/
├── partials/              # 可复用片段
│   ├── head.plush.html
│   ├── header.plush.html
│   ├── footer.plush.html
│   └── sidebar.plush.html
└── screenshots/           # 预览图
```

---

## 八、常见问题

### Q1: 配置保存后不生效？

**A**: 检查以下几点：
1. 配置键名是否与 `theme.yml` 中定义的 `name` 一致
2. 后端是否正确从数据库加载配置并合并到主题配置
3. 模板中是否使用了正确的键名访问配置

### Q2: 模板渲染报错？

**A**: 常见原因：
1. Plush 语法错误（如缺少闭合括号）
2. 访问了不存在的变量
3. 模板文件路径错误

### Q3: 如何调试模板？

**A**: 可以临时添加调试输出：
```html
<!-- 查看所有可用数据 -->
<pre><%= inspect(data) %></pre>

<!-- 查看配置 -->
<pre><%= inspect(config) %></pre>
```

### Q4: 如何支持深色模式？

**A**: 使用 CSS 变量和 `prefers-color-scheme`：
```css
:root {
  --bg-color: #ffffff;
  --text-color: #1f2937;
}

@media (prefers-color-scheme: dark) {
  :root {
    --bg-color: #1f2937;
    --text-color: #f9fafb;
  }
}
```

---

## 九、参考资源

- [Halo 主题开发文档](https://docs.halo.run/developer-guide/theme/)
- [Plush 模板引擎文档](https://github.com/gobuffalo/plush)
- [Catter CMS 示例主题](templates/themes/catter-blog/)

---

**总结**：Catter CMS 的主题系统参考 Halo 设计，通过 `theme.yml` 定义配置，Plush 引擎渲染模板，实现了灵活、可配置的主题机制。开发新主题只需创建配置文件和模板文件即可。
