# Catter CMS 页面模板语法参考（仅供对照，严禁复制）

> **本文档的使用场景说明：**
>
> ✅ **以下用法完全合理：**
> - 查阅模板函数的调用语法（`getConfig`、`asset_path`、`markdownHTML` 等）
> - 确认控制结构的写法（`range`、`if`、`with`）
> - 核对数据字段名拼写（`.post.Title`、`.Category.Slug` 等）
> - 二次开发已有主题时，参考类似页面的现有结构
>
> ⛔ **仅在「开发全新主题」时禁止以下行为：**
> - 将本文档的 HTML 骨架/模块结构直接照搬到新主题（换皮不是新主题）
> - 沿用 `.post-card`、`.hero-inner`、`.sidebar` 等参考模板类名（新主题应有自己的类名体系）
> - 以参考模板的布局作为新主题的默认起点（应从 Phase 1 Design Brief 出发）
>
> **区别在于意图**：看懂语法再自己写 ≠ 复制结构。前者是学习参考，后者是换皮。

---

## 首页模板（home.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .site.title }}</title>
    <meta name="description" content="{{ .site.description }}">
    <link rel="icon" href="{{ .site.favicon }}">
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "assets/css/home.css" }}">
    <style>
    :root {
        --primary: {{ getConfig "primaryColor" "#2563eb" }};
        --primary-dark: {{ getConfig "primaryDark" "#0f3460" }};
        --accent: {{ getConfig "accentColor" "#60a5fa" }};
        --bg-primary: {{ getConfig "bgPrimary" "#f8fafc" }};
        --text-primary: {{ getConfig "textPrimary" "#1e293b" }};
    }
    </style>
</head>
<body>
    <!-- ===== 导航 ===== -->
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">
                {{ if .site.logo }}<img src="{{ .site.logo }}" alt="{{ .site.title }}" class="logo-img">{{ end }}
                <span>{{ .site.title }}</span>
            </a>
            <nav class="nav-menu">
                {{ range .headerMenu.Items }}
                <a href="{{ .URL }}" class="nav-link"
                   {{ if eq .Target "_blank" }}target="_blank" rel="noopener"{{ end }}>
                    {{ .Label }}
                </a>
                {{ end }}
            </nav>
        </div>
    </header>

    <!-- ===== Hero 区域 ===== -->
    {{ if getConfig "heroEnabled" true }}
    <section class="hero">
        <div class="hero-inner">
            <h1 class="hero-title">{{ getConfig "heroTitle" "欢迎来到" }} {{ .site.title }}</h1>
            <p class="hero-subtitle">{{ getConfig "heroSubtitle" .site.description }}</p>
            {{ if getConfig "heroButtonText" "" }}
            <a href="{{ getConfig "heroButtonLink" "/article" }}" class="hero-btn">
                {{ getConfig "heroButtonText" "开始阅读" }}
            </a>
            {{ end }}
        </div>
    </section>
    {{ end }}

    <!-- ===== 主内容 ===== -->
    <div class="main-container {{ if getConfig "sidebarEnabled" true }}has-sidebar{{ end }}">

        <!-- 文章列表 -->
        <main class="content-area">
            <h2 class="section-title">最新文章</h2>
            {{ if gt (len .posts) 0 }}
            <div class="posts-list style-{{ getConfig "postListStyle" "card" }}">
                {{ range .posts }}
                <article class="post-card">
                    <div class="post-body">
                        {{ if .Category }}
                        <a href="/category/{{ .Category.Slug }}" class="post-category">{{ .Category.Name }}</a>
                        {{ end }}
                        <h3 class="post-title">
                            <a href="/article/{{ .Slug }}">{{ .Title }}</a>
                        </h3>
                        <div class="post-meta">
                            {{ if .PublishedAt }}
                            <span class="post-date">
                                <i class="fa-regular fa-calendar"></i>
                                {{ dateFormat "2006-01-02" .PublishedAt }}
                            </span>
                            {{ end }}
                            {{ if getConfig "showViews" true }}
                            <span class="post-views">
                                <i class="fa-regular fa-eye"></i>
                                {{ .Views }}
                            </span>
                            {{ end }}
                        </div>
                        {{ if .Excerpt }}
                        <p class="post-excerpt">{{ substr .Excerpt 0 150 }}...</p>
                        {{ end }}
                        <a href="/article/{{ .Slug }}" class="read-more">阅读全文 →</a>
                    </div>
                </article>
                {{ end }}
            </div>

            <!-- 分页 -->
            {{ if gt .paginator.TotalPages 1 }}
            <nav class="pagination">
                {{ if gt .paginator.Page 1 }}
                <a href="?page={{ sub .paginator.Page 1 }}" class="page-btn">‹ 上一页</a>
                {{ end }}
                {{ range .paginator.Pages }}
                <a href="?page={{ . }}"
                   class="page-btn {{ if eq . $.paginator.Page }}active{{ end }}">{{ . }}</a>
                {{ end }}
                {{ if lt .paginator.Page .paginator.TotalPages }}
                <a href="?page={{ add .paginator.Page 1 }}" class="page-btn">下一页 ›</a>
                {{ end }}
            </nav>
            {{ end }}

            {{ else }}
            <div class="empty-state">
                <i class="fa-regular fa-file-lines"></i>
                <p>暂无文章</p>
            </div>
            {{ end }}
        </main>

        <!-- 侧边栏 -->
        {{ if getConfig "sidebarEnabled" true }}
        <aside class="sidebar">
            <!-- 分类小部件 -->
            {{ if gt (len .categories) 0 }}
            <div class="widget">
                <h4 class="widget-title">分类</h4>
                <ul class="widget-category-list">
                    {{ range .categories }}
                    <li>
                        <a href="/category/{{ .Slug }}" class="category-link">
                            <span>{{ .Name }}</span>
                        </a>
                    </li>
                    {{ end }}
                </ul>
            </div>
            {{ end }}

            <!-- 标签云小部件 -->
            {{ if gt (len .tags) 0 }}
            <div class="widget">
                <h4 class="widget-title">标签</h4>
                <div class="tags-cloud">
                    {{ range .tags }}
                    <a href="/tag/{{ .Slug }}" class="tag-item">{{ .Name }}</a>
                    {{ end }}
                </div>
            </div>
            {{ end }}
        </aside>
        {{ end }}
    </div>

    <!-- ===== 页脚 ===== -->
    <footer class="footer">
        <div class="footer-container">
            <nav class="footer-nav">
                {{ range .footerMenu.Items }}
                <a href="{{ .URL }}" class="footer-link">{{ .Label }}</a>
                {{ end }}
            </nav>
            <p class="footer-copy">
                &copy; {{ .site.year }} {{ .site.title }}.
                {{ if getConfig "footerText" "" }}{{ getConfig "footerText" "" }}{{ end }}
            </p>
        </div>
    </footer>
</body>
</html>
```

---

## 文章详情页（article_show.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .post.Title }} - {{ .site.title }}</title>
    <meta name="description" content="{{ default .post.MetaDescription .post.Excerpt }}">
    <link rel="icon" href="{{ .site.favicon }}">
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "assets/css/article.css" }}">
    <style>
    :root {
        --primary: {{ getConfig "primaryColor" "#2563eb" }};
        --primary-dark: {{ getConfig "primaryDark" "#0f3460" }};
    }
    </style>
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                {{ range .headerMenu.Items }}
                <a href="{{ .URL }}" class="nav-link">{{ .Label }}</a>
                {{ end }}
            </nav>
        </div>
    </header>

    <div class="article-wrapper">
        <!-- 文章头部信息 -->
        <div class="article-header">
            {{ if .post.Category }}
            <a href="/category/{{ .post.Category.Slug }}" class="article-category">
                {{ .post.Category.Name }}
            </a>
            {{ end }}
            <h1 class="article-title">{{ .post.Title }}</h1>
            <div class="article-meta">
                {{ if .post.PublishedAt }}
                <span><i class="fa-regular fa-calendar"></i> {{ dateFormat "2006年01月02日" .post.PublishedAt }}</span>
                {{ end }}
                <span><i class="fa-regular fa-eye"></i> {{ .post.Views }} 次阅读</span>
                {{ if .post.Author }}
                <span><i class="fa-regular fa-user"></i> {{ .post.Author.Name }}</span>
                {{ end }}
            </div>
        </div>

        <!-- 文章正文 -->
        <article class="article-content">
            {{ markdownHTML .post.Content }}
        </article>

        <!-- 标签 -->
        {{ if .post.Tags }}
        <div class="article-tags">
            <span>标签：</span>
            {{ range .post.Tags }}
            <a href="/tag/{{ .Slug }}" class="tag-item">{{ .Name }}</a>
            {{ end }}
        </div>
        {{ end }}
    </div>

    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```

---

## 分类页（category.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .category.Name }} - {{ .site.title }}</title>
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "assets/css/list.css" }}">
    <style>:root{ --primary: {{ getConfig "primaryColor" "#2563eb" }}; }</style>
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                {{ range .headerMenu.Items }}
                <a href="{{ .URL }}" class="nav-link">{{ .Label }}</a>
                {{ end }}
            </nav>
        </div>
    </header>

    <div class="page-header">
        <h1>{{ .category.Name }}</h1>
        {{ if .category.Description }}
        <p>{{ .category.Description }}</p>
        {{ end }}
    </div>

    <main class="list-container">
        {{ range .posts }}
        <article class="article-card">
            <h2 class="article-card-title">
                <a href="/article/{{ .Slug }}">{{ .Title }}</a>
            </h2>
            <div class="article-card-meta">
                {{ if .PublishedAt }}
                <span>{{ dateFormat "2006-01-02" .PublishedAt }}</span>
                {{ end }}
                <span>{{ .Views }} 阅读</span>
            </div>
            {{ if .Excerpt }}
            <p class="article-card-excerpt">{{ substr .Excerpt 0 150 }}...</p>
            {{ end }}
            <a href="/article/{{ .Slug }}" class="read-more">阅读全文 →</a>
        </article>
        {{ else }}
        <div class="empty-state">
            <i class="fa-regular fa-folder-open"></i>
            <p>该分类下暂无文章</p>
        </div>
        {{ end }}

        <!-- 分页 -->
        {{ if gt .paginator.TotalPages 1 }}
        <nav class="pagination">
            {{ range .paginator.Pages }}
            <a href="?page={{ . }}"
               class="page-btn {{ if eq . $.paginator.Page }}active{{ end }}">{{ . }}</a>
            {{ end }}
        </nav>
        {{ end }}
    </main>

    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```

---

## 归档页（archives.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>归档 - {{ .site.title }}</title>
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "assets/css/list.css" }}">
    <style>:root{ --primary: {{ getConfig "primaryColor" "#2563eb" }}; }</style>
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                {{ range .headerMenu.Items }}
                <a href="{{ .URL }}" class="nav-link">{{ .Label }}</a>
                {{ end }}
            </nav>
        </div>
    </header>

    <div class="page-header">
        <h1>文章归档</h1>
    </div>

    <main class="list-container">
        {{ range $year, $posts := .archives }}
        <div class="archive-year">
            <h2 class="archive-year-title">{{ $year }}</h2>
            <ul class="archive-list">
                {{ range $posts }}
                <li class="archive-item">
                    <span class="archive-date">
                        {{ if .PublishedAt }}{{ dateFormat "01-02" .PublishedAt }}{{ end }}
                    </span>
                    <a href="/article/{{ .Slug }}" class="archive-title">{{ .Title }}</a>
                </li>
                {{ end }}
            </ul>
        </div>
        {{ end }}
    </main>

    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```

---

## 标签云页（tags.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>标签 - {{ .site.title }}</title>
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "assets/css/list.css" }}">
    <style>:root{ --primary: {{ getConfig "primaryColor" "#2563eb" }}; }</style>
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                {{ range .headerMenu.Items }}
                <a href="{{ .URL }}" class="nav-link">{{ .Label }}</a>
                {{ end }}
            </nav>
        </div>
    </header>

    <div class="page-header">
        <h1>标签云</h1>
    </div>

    <main class="list-container">
        {{ if gt (len .tags) 0 }}
        <div class="tags-cloud-container">
            {{ range .tags }}
            <a href="/tag/{{ .Slug }}" class="tag-item">
                {{ .Name }}
                <span class="tag-count">{{ .PostCount }}</span>
            </a>
            {{ end }}
        </div>
        {{ else }}
        <div class="empty-state">
            <i class="fa-solid fa-tags"></i>
            <p>暂无标签</p>
        </div>
        {{ end }}
    </main>

    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```

---

## 搜索页（search.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>搜索: {{ .query }} - {{ .site.title }}</title>
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "assets/css/list.css" }}">
    <style>:root{ --primary: {{ getConfig "primaryColor" "#2563eb" }}; }</style>
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                {{ range .headerMenu.Items }}
                <a href="{{ .URL }}" class="nav-link">{{ .Label }}</a>
                {{ end }}
            </nav>
        </div>
    </header>

    <div class="page-header">
        <h1>搜索结果</h1>
        <form class="search-form" action="/search" method="get">
            <input type="text" name="q" value="{{ .query }}" class="search-input" placeholder="搜索文章...">
            <button type="submit" class="search-btn">
                <i class="fa-solid fa-magnifying-glass"></i>
            </button>
        </form>
    </div>

    <main class="list-container">
        {{ if .query }}
        <p class="results-info">
            "{{ .query }}" 的搜索结果：共 {{ len .posts }} 篇
        </p>
        {{ end }}

        {{ range .posts }}
        <article class="article-card">
            <h2 class="article-card-title">
                <a href="/article/{{ .Slug }}">{{ .Title }}</a>
            </h2>
            <div class="article-card-meta">
                {{ if .PublishedAt }}
                <span>{{ dateFormat "2006-01-02" .PublishedAt }}</span>
                {{ end }}
            </div>
            {{ if .Excerpt }}
            <p class="article-card-excerpt">{{ substr .Excerpt 0 150 }}...</p>
            {{ end }}
        </article>
        {{ else }}
        <div class="empty-state">
            <i class="fa-solid fa-magnifying-glass"></i>
            <p>未找到相关文章</p>
            {{ if not .query }}
            <p class="search-hint">请在上方搜索框输入关键词</p>
            {{ end }}
        </div>
        {{ end }}
    </main>

    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```

---

## 独立页面（page.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .page.Title }} - {{ .site.title }}</title>
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <link rel="stylesheet" href="{{ asset_path "assets/css/article.css" }}">
    <style>:root{ --primary: {{ getConfig "primaryColor" "#2563eb" }}; }</style>
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
            <nav class="nav-menu">
                {{ range .headerMenu.Items }}
                <a href="{{ .URL }}" class="nav-link">{{ .Label }}</a>
                {{ end }}
            </nav>
        </div>
    </header>

    <div class="page-header">
        <h1>{{ .page.Title }}</h1>
    </div>

    <div class="article-wrapper">
        <article class="article-content">
            {{ markdownHTML .page.Content }}
        </article>
    </div>

    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```

---

## 404 页面（404.html）

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 - 页面不存在 | {{ .site.title }}</title>
    <link rel="stylesheet" href="/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
    <style>:root{ --primary: {{ getConfig "primaryColor" "#2563eb" }}; }</style>
</head>
<body>
    <header class="header">
        <div class="header-container">
            <a href="/" class="logo">{{ .site.title }}</a>
        </div>
    </header>

    <div class="not-found">
        <div class="not-found-inner">
            <h1 class="not-found-code">404</h1>
            <h2 class="not-found-title">页面不存在</h2>
            <p class="not-found-desc">您访问的页面已被删除或不存在</p>
            <a href="/" class="btn-primary">返回首页</a>
        </div>
    </div>

    <footer class="footer">
        <div class="footer-container">
            <p>&copy; {{ .site.year }} {{ .site.title }}</p>
        </div>
    </footer>
</body>
</html>
```
