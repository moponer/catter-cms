# Catter CMS（Buffalo + Go）架构设计

# Catter CMS 架构设计（Buffalo + Go 版）

---

name: catter-cms-design-buffalo

description: "Catter CMS - Buffalo (Go) 轻量级内容管理系统架构设计"

---

## 1. 项目概述

### 1.1 项目背景

Catter CMS（Go 版）是基于 **Buffalo** 框架构建的轻量级内容管理系统，专为中小企业设计。Buffalo 作为 Go 生态中开箱即用的全栈 Web 框架（对标 Rails/Django），提供了路由、模板、ORM、迁移、认证等一站式能力，本架构充分利用其特性，实现高性能、易部署、强类型的 CMS 系统。

### 1.2 核心特性

|特性|描述|
|---|---|
|全栈 SSR|前后端一体化渲染，无需分离部署，兼顾性能与开发效率|
|多主题系统|主题隔离，配置驱动，支持热切换|
|多语言|内置 i18n 支持，内容与界面多语言适配|
|多数据库支持|开发用 SQLite，生产可无缝切换 PostgreSQL|
|权限管理|管理员 + 编辑者双角色，细粒度权限控制|
|增量静态生成|核心页面预生成，更新时增量重建|
### 1.3 目标用户

- 内容编辑者：创建、编辑、发布内容

- 网站管理员：系统配置、用户管理、主题管理

- 普通客户：访问前台内容

## 2. 技术架构

### 2.1 技术栈

|层级|技术|说明|
|---|---|---|
|核心框架|Buffalo v0.18+|Go 全栈 Web 框架，提供路由、中间件、模板等核心能力|
|语言|Go 1.21+|强类型、高性能，编译后无依赖部署|
|ORM|Pop (Buffalo 内置)|支持 SQLite/PostgreSQL，提供迁移、模型映射能力|
|模板引擎|Plush (Buffalo 内置)|兼容 Go 模板语法，支持组件化、布局嵌套|
|前端样式|Tailwind CSS 4|原子化 CSS，主题样式隔离|
|认证|JWT + bcrypt|基于 Cookie/Session 的认证，支持 JWT 接口认证|
|i18n|Buffalo i18n|内置多语言支持，适配界面与内容翻译|
|静态生成|自定义增量生成工具|基于 Go 脚本实现页面预生成，支持增量更新|
|缓存|Redis (可选)|前台页面缓存、API 结果缓存|
### 2.2 渲染模式

```Plain Text

┌─────────────────────────────────────────────────────────┐
│                      Catter CMS (Go)                    │
├────────────────────────┬────────────────────────────────┤
│      管理后台          │           前台                  │
├────────────────────────┼────────────────────────────────┤
│  渲染模式: SSR + SPA   │  渲染模式: SSR + 增量静态生成   │
│  路径: /admin/*       │  路径: /* (除 /admin/*)         │
│  交互: 实时表单/操作   │  交互: 静态页面 + 动态数据      │
│  无需 SEO             │  完整 SEO 支持                 │
│  权限: 登录验证        │  权限: 公开访问                 │
└────────────────────────┴────────────────────────────────┘
```

### 2.3 核心配置（Buffalo）

```Go

// config/app.go
package config

import (
 "github.com/gobuffalo/buffalo"
 "github.com/gobuffalo/envy"
 "github.com/rs/cors"
 "log"
)

// App 初始化核心配置
func App() *buffalo.App {
 app := buffalo.New(buffalo.Options{
 Env:         envy.Get("GO_ENV", "development"),
 SessionName: "_catter_cms_session",
 })

 // 全局中间件
 app.Use(buffalo.RequestLogger)
 app.Use(buffalo.Recoverer)
 app.Use(buffalo.SessionSaver)

 // CORS 配置（API 跨域）
 app.Use(cors.New(cors.Options{
 AllowedOrigins:   []string{"*"},
 AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
 AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
 AllowCredentials: true,
 }).Handler)

 // 路由规则
 // 前台路由 - 静态优先，增量生成
 app.GET("/", HomeHandler)
 app.GET("/blog/{slug}", BlogShowHandler)
 app.GET("/{slug}", PageShowHandler)

 // 管理后台路由 - SSR + 权限验证
 admin := app.Group("/admin")
 admin.Use(AuthMiddleware) // 登录验证
 admin.GET("/", AdminDashboardHandler)
 admin.GET("/login", AdminLoginHandler)
 admin.POST("/login", AdminLoginPostHandler)
 admin.RESTful("/posts", PostsResource{}) // RESTful 资源路由

 // API 路由
 api := app.Group("/api")
 api.Use(APIAuthMiddleware)
 api.GET("/posts", APIPostsIndexHandler)
 api.POST("/posts", APIPostsCreateHandler)

 return app
}

// 数据库配置
func DB() *pop.Connection {
 env := envy.Get("GO_ENV", "development")
 db, err := pop.Connect(env)
 if err != nil {
 log.Fatal(err)
 }
 return db
}
```

## 3. 项目结构

### 3.1 完整目录结构

```Plain Text

catter-cms/
├── actions/                     # 处理请求的核心逻辑（控制器）
│   ├── admin/                   # 管理后台 Actions
│   │   ├── dashboard.go         # 仪表盘
│   │   ├── login.go             # 登录
│   │   ├── posts.go             # 文章管理
│   │   ├── media.go             # 媒体库
│   │   ├── users.go             # 用户管理
│   │   ├── themes.go            # 主题管理
│   │   └── settings.go          # 系统设置
│   ├── front/                   # 前台 Actions
│   │   ├── home.go              # 首页
│   │   ├── blog.go              # 博客列表/详情
│   │   ├── page.go              # 静态页面
│   │   └── sitemap.go           # 站点地图
│   ├── api/                     # API Actions
│   │   ├── posts.go             # 文章 API
│   │   ├── media.go             # 媒体 API
│   │   └── themes.go            # 主题 API
│   └── init.go                  # Actions 初始化
│
├── models/                      # 数据模型
│   ├── user.go                  # 用户模型
│   ├── post.go                  # 文章模型
│   ├── category.go              # 分类模型
│   ├── tag.go                   # 标签模型
│   ├── media.go                 # 媒体模型
│   ├── menu.go                  # 菜单模型
│   ├── theme.go                 # 主题模型
│   ├── setting.go               # 系统设置模型
│   └── db.go                    # 数据库初始化
│
├── templates/                   # 模板文件
│   ├── admin/                   # 管理后台模板
│   │   ├── layouts/             # 后台布局
│   │   │   └── application.plush.html
│   │   ├── partials/            # 后台公共组件
│   │   │   ├── header.plush.html
│   │   │   ├── sidebar.plush.html
│   │   │   └── footer.plush.html
│   │   ├── dashboard.plush.html # 仪表盘页面
│   │   ├── login.plush.html     # 登录页
│   │   └── posts/               # 文章管理模板
│   │
│   ├── themes/                  # 主题模板（可插拔）
│   │   ├── blog/                # Blog 主题
│   │   │   ├── layouts/         # 主题布局
│   │   │   │   └── default.plush.html
│   │   │   ├── partials/        # 主题组件
│   │   │   │   ├── header.plush.html
│   │   │   │   ├── footer.plush.html
│   │   │   │   └── post_card.plush.html
│   │   │   ├── home.plush.html  # 首页模板
│   │   │   ├── blog.plush.html  # 博客列表模板
│   │   │   └── blog_show.plush.html # 博客详情模板
│   │   └── corporate/           # 企业官网主题
│   │       ├── layouts/
│   │       ├── partials/
│   │       └── pages/
│   │
│   └── shared/                  # 共享模板组件
│       ├── form_helpers.plush.html
│       └── pagination.plush.html
│
├── assets/                      # 静态资源
│   ├── css/                     # 样式文件
│   │   ├── admin/               # 后台样式
│   │   │   └── application.css
│   │   └── themes/              # 主题样式
│   │       ├── blog/
│   │       │   └── theme.css
│   │       └── corporate/
│   │           └── theme.css
│   ├── js/                      # 脚本文件
│   │   ├── admin/               # 后台脚本
│   │   └── themes/              # 主题脚本
│   └── images/                  # 图片资源
│
├── public/                      # 编译后静态文件（自动生成）
│   ├── assets/                  # 编译后的 CSS/JS
│   ├── uploads/                 # 用户上传文件
│   └── static/                  # 静态生成的页面
│
├── migrations/                  # 数据库迁移文件
│   ├── 20260310120000_create_users.up.sql
│   ├── 20260310120100_create_posts.up.sql
│   └── ...
│
├── seeds/                       # 种子数据
│   ├── users.go
│   └── posts.go
│
├── locales/                     # 多语言文件
│   ├── zh-CN.yaml
│   └── en.yaml
│
├── tasks/                       # Buffalo 任务（增量生成、备份等）
│   ├── generate_static.go       # 静态页面生成任务
│   ├── revalidate.go            # 增量更新任务
│   └── backup_db.go             # 数据库备份任务
│
├── config/                      # 配置文件
│   ├── app.go                   # 应用配置
│   ├── database.yml             # 数据库配置
│   └── secrets.yml              # 密钥配置
│
├── plugins/                     # 扩展插件
│   ├── seo/                     # SEO 插件
│   └── analytics/               # 统计插件
│
├── .env                         # 环境变量
├── .env.example                 # 环境变量示例
├── go.mod                       # Go 模块
├── go.sum                       # Go 依赖
├── main.go                      # 入口文件
└── catter-cms                   # 编译后二进制文件
```

### 3.2 目录职责

|目录|职责|
|---|---|
|`actions/`|处理 HTTP 请求的核心逻辑（对应 MVC 的 Controller）|
|`models/`|数据模型与数据库交互（对应 MVC 的 Model）|
|`templates/`|模板文件（对应 MVC 的 View），按后台/主题隔离|
|`assets/`|未编译的静态资源（CSS/JS/图片）|
|`public/`|编译后静态文件，可直接部署|
|`migrations/`|数据库迁移文件，管理表结构变更|
|`locales/`|多语言翻译文件|
|`tasks/`|自定义 CLI 任务（静态生成、备份等）|
|`config/`|框架与应用配置|
|`templates/themes/`|可插拔主题模板，支持热切换|
## 4. 数据模型

### 4.1 核心模型定义（Go）

```Go

// models/user.go
package models

import (
 "time"

 "github.com/gobuffalo/pop/v6"
 "github.com/gobuffalo/validate/v3"
 "github.com/gobuffalo/validate/v3/validators"
 "github.com/gofrs/uuid"
 "golang.org/x/crypto/bcrypt"
)

// User 定义用户模型
type User struct {
 ID        uuid.UUID `json:"id" db:"id"`
 Email     string    `json:"email" db:"email"`
 Name      string    `json:"name" db:"name"`
 Password  string    `json:"-" db:"password"` // 不返回给前端
 Role      string    `json:"role" db:"role"` // admin / editor
 Avatar    string    `json:"avatar,omitempty" db:"avatar"`
 CreatedAt time.Time `json:"created_at" db:"created_at"`
 UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate 验证用户数据
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
 return validate.Validate(
 &validators.StringIsPresent{Field: u.Email, Name: "Email"},
 &validators.StringIsPresent{Field: u.Name, Name: "Name"},
 &validators.StringIsPresent{Field: u.Password, Name: "Password"},
 &validators.StringInclusion{Field: u.Role, Name: "Role", List: []string{"admin", "editor"}},
 ), nil
}

// BeforeSave 保存前加密密码
func (u *User) BeforeSave(tx *pop.Connection) error {
 if len(u.Password) > 0 {
 hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
 if err != nil {
 return err
 }
 u.Password = string(hash)
 }
 return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
 err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
 return err == nil
}

// models/post.go
package models

import (
 "time"

 "github.com/gobuffalo/pop/v6"
 "github.com/gobuffalo/validate/v3"
 "github.com/gobuffalo/validate/v3/validators"
 "github.com/gofrs/uuid"
)

// Post 定义文章模型
type Post struct {
 ID            uuid.UUID `json:"id" db:"id"`
 Title         string    `json:"title" db:"title"`
 Slug          string    `json:"slug" db:"slug"`
 Content       string    `json:"content" db:"content"`
 Excerpt       string    `json:"excerpt,omitempty" db:"excerpt"`
 FeaturedImage string    `json:"featured_image,omitempty" db:"featured_image"`
 Status        string    `json:"status" db:"status"` // draft / published / archived
 AuthorID      uuid.UUID `json:"author_id" db:"author_id"`
 Author        User      `json:"author,omitempty" belongs_to:"users"`
 CategoryID    uuid.UUID `json:"category_id,omitempty" db:"category_id"`
 Category      Category  `json:"category,omitempty" belongs_to:"categories"`
 Tags          []Tag     `json:"tags,omitempty" many_to_many:"post_tags"`
 Locale        string    `json:"locale" db:"locale"` // zh-CN / en
 SEOTitle      string    `json:"seo_title,omitempty" db:"seo_title"`
 SEODescription string   `json:"seo_description,omitempty" db:"seo_description"`
 PublishedAt   *time.Time `json:"published_at,omitempty" db:"published_at"`
 CreatedAt     time.Time `json:"created_at" db:"created_at"`
 UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Validate 验证文章数据
func (p *Post) Validate(tx *pop.Connection) (*validate.Errors, error) {
 return validate.Validate(
 &validators.StringIsPresent{Field: p.Title, Name: "Title"},
 &validators.StringIsPresent{Field: p.Slug, Name: "Slug"},
 &validators.StringIsPresent{Field: p.Content, Name: "Content"},
 &validators.StringInclusion{Field: p.Status, Name: "Status", List: []string{"draft", "published", "archived"}},
 &validators.StringIsPresent{Field: p.Locale, Name: "Locale"},
 ), nil
}

// 其他核心模型（Category/Tag/Media/Menu/Theme/Setting）类似
```

### 4.2 数据库配置

```YAML

# config/database.yml
development:
  dialect: sqlite3
  database: catter_development.db
  pool: 5
  host: localhost
  migrations_table: schema_migrations

test:
  dialect: sqlite3
  database: catter_test.db
  pool: 5
  host: localhost
  migrations_table: schema_migrations

production:
  dialect: postgres
  database: <%= env("DATABASE_NAME") %>
  user: <%= env("DATABASE_USER") %>
  password: <%= env("DATABASE_PASSWORD") %>
  host: <%= env("DATABASE_HOST") %>
  port: <%= env("DATABASE_PORT") %>
  pool: 10
  ssl_mode: require
  migrations_table: schema_migrations
```

## 5. API 设计

### 5.1 核心 API 路由（RESTful）

|方法|路径|处理函数|描述|权限|
|---|---|---|---|---|
|GET|/api/posts|APIPostsIndex|文章列表|公开|
|GET|/api/posts/{id}|APIPostsShow|文章详情|公开|
|POST|/api/posts|APIPostsCreate|创建文章|编辑+|
|PUT|/api/posts/{id}|APIPostsUpdate|更新文章|编辑+|
|DELETE|/api/posts/{id}|APIPostsDestroy|删除文章|编辑+|
|GET|/api/categories|APICategoriesIndex|分类列表|公开|
|POST|/api/categories|APICategoriesCreate|创建分类|编辑+|
|GET|/api/themes|APIThemesIndex|主题列表|公开|
|POST|/api/themes/{id}/activate|APIThemesActivate|激活主题|管理员|
|GET|/api/auth/me|APIAuthMe|当前用户信息|登录|
|POST|/api/auth/login|APIAuthLogin|用户登录|公开|
|POST|/api/auth/logout|APIAuthLogout|用户登出|登录|
### 5.2 API 处理函数示例

```Go

// actions/api/posts.go
package api

import (
 "net/http"

 "github.com/gobuffalo/buffalo"
 "github.com/gobuffalo/pop/v6"
 "github.com/gofrs/uuid"

 "catter-cms/models"
 "catter-cms/tasks"
)

// APIPostsIndex 获取文章列表
func APIPostsIndex(c buffalo.Context) error {
 tx := c.Value("tx").(*pop.Connection)
 posts := &models.Posts{}

 // 查询参数
 limit := c.Param("limit")
 if limit == "" {
 limit = "10"
 }
 page := c.Param("page")
 if page == "" {
 page = "1"
 }
 status := c.Param("status")
 locale := c.Param("locale")

 // 构建查询
 q := tx.Paginate(page, limit)
 q = q.Where("status = ?", "published")
 if locale != "" {
 q = q.Where("locale = ?", locale)
 }

 // 执行查询
 if err := q.All(posts); err != nil {
 return c.Error(http.StatusInternalServerError, err)
 }

 // 返回 JSON
 return c.JSON(http.StatusOK, map[string]interface{}{
 "posts": posts,
 "meta": map[string]interface{}{
 "total_pages": q.TotalPages,
 "total_count": q.TotalEntries,
 "page":        q.Page,
 "per_page":    q.PerPage,
 },
 })
}

// APIPostsCreate 创建文章
func APIPostsCreate(c buffalo.Context) error {
 tx := c.Value("tx").(*pop.Connection)
 post := &models.Post{}

 // 解析请求体
 if err := c.Bind(post); err != nil {
 return c.Error(http.StatusBadRequest, err)
 }

 // 设置默认值
 user := c.Value("current_user").(*models.User)
 post.AuthorID = user.ID
 post.Status = "draft"
 if post.Locale == "" {
 post.Locale = "zh-CN"
 }

 // 验证并保存
 verrs, err := post.Validate(tx)
 if err != nil {
 return c.Error(http.StatusInternalServerError, err)
 }
 if verrs.HasAny() {
 return c.Error(http.StatusUnprocessableEntity, verrs)
 }

 if err := tx.Create(post); err != nil {
 return c.Error(http.StatusInternalServerError, err)
 }

 // 触发增量更新
 go tasks.RevalidatePost(post.Slug)

 return c.JSON(http.StatusCreated, post)
}
```

## 6. 主题系统

### 6.1 主题结构

```Plain Text

templates/themes/{theme-slug}/
├── layouts/             # 布局模板
│   ├── default.plush.html
│   ├── blog.plush.html
│   └── fullwidth.plush.html
├── partials/            # 组件模板
│   ├── header.plush.html
│   ├── footer.plush.html
│   ├── post_card.plush.html
│   └── hero.plush.html
├── pages/               # 页面模板
│   ├── home.plush.html
│   ├── blog.plush.html
│   └── blog_show.plush.html
├── config/              # 主题配置
│   ├── menus.yml        # 菜单配置
│   └── settings.yml     # 主题设置
├── theme.yml            # 主题元数据
└── assets/              # 主题静态资源
    ├── css/
    │   └── theme.css
    ├── js/
    │   └── theme.js
    └── images/
        └── thumbnail.jpg
```

### 6.2 主题元数据

```YAML

# templates/themes/blog/theme.yml
name: "Blog Theme"
slug: "blog"
version: "1.0.0"
description: "现代博客主题"
author: "Catter Team"
thumbnail: "/themes/blog/assets/images/thumbnail.jpg"
requires:
  buffalo: ">=0.18.0"
config:
  default_locale: "zh-CN"
  primary_color: "#3B82F6"
  layouts: ["default", "blog", "fullwidth"]
  modules: ["post-list", "hero", "categories", "tags"]
```

### 6.3 主题切换逻辑

```Go

// actions/admin/themes.go
package admin

import (
 "net/http"
 "os"
 "path/filepath"

 "github.com/gobuffalo/buffalo"
 "gopkg.in/yaml.v3"

 "catter-cms/models"
 "catter-cms/tasks"
)

// ThemesActivate 激活主题
func ThemesActivate(c buffalo.Context) error {
 tx := c.Value("tx").(*pop.Connection)
 themeSlug := c.Param("slug")

 // 1. 查询主题
 theme := &models.Theme{}
 if err := tx.Where("slug = ?", themeSlug).First(theme); err != nil {
 return c.Error(http.StatusNotFound, err)
 }

 // 2. 更新系统设置（激活主题）
 setting := &models.Setting{}
 if err := tx.Where("key = ?", "active_theme").First(setting); err != nil {
 setting = &models.Setting{
 Key:      "active_theme",
 Value:    themeSlug,
 Type:     "string",
 Category: "general",
 }
 } else {
 setting.Value = themeSlug
 }
 if err := tx.Save(setting); err != nil {
 return c.Error(http.StatusInternalServerError, err)
 }

 // 3. 加载主题菜单
 menuConfigPath := filepath.Join("templates", "themes", themeSlug, "config", "menus.yml")
 menuConfig, err := loadYAMLConfig(menuConfigPath)
 if err != nil {
 return c.Error(http.StatusInternalServerError, err)
 }
 // 保存菜单到数据库
 if err := saveThemeMenus(tx, themeSlug, menuConfig); err != nil {
 return c.Error(http.StatusInternalServerError, err)
 }

 // 4. 触发全量静态生成
 go tasks.GenerateStaticSite()

 // 5. 返回结果
 return c.JSON(http.StatusOK, map[string]string{
 "message": "主题已激活",
 "theme":   themeSlug,
 })
}

// loadYAMLConfig 加载 YAML 配置文件
func loadYAMLConfig(path string) (map[string]interface{}, error) {
 data, err := os.ReadFile(path)
 if err != nil {
 return nil, err
 }
 var config map[string]interface{}
 if err := yaml.Unmarshal(data, &config); err != nil {
 return nil, err
 }
 return config, nil
}

// saveThemeMenus 保存主题菜单到数据库
func saveThemeMenus(tx *pop.Connection, themeSlug string, config map[string]interface{}) error {
 // 实现菜单保存逻辑（根据配置创建/更新 Menu 模型）
 return nil
}
```

## 7. 多语言支持

### 7.1 语言配置

```YAML

# locales/zh-CN.yaml
common:
  home: 首页
  blog: 博客
  about: 关于
  contact: 联系我们
admin:
  dashboard: 仪表盘
  posts: 文章管理
  media: 媒体库
  settings: 系统设置
post:
  title: 标题
  content: 内容
  publish: 发布
  draft: 草稿

# locales/en.yaml
common:
  home: Home
  blog: Blog
  about: About
  contact: Contact Us
admin:
  dashboard: Dashboard
  posts: Posts
  media: Media
  settings: Settings
post:
  title: Title
  content: Content
  publish: Publish
  draft: Draft
```

### 7.2 多语言模板使用

```HTML

<!-- templates/shared/partials/header.plush.html -->
<nav>
  <a href="/">{{ t("common.home") }}</a>
  <a href="/blog">{{ t("common.blog") }}</a>
  <a href="/about">{{ t("common.about") }}</a>
  
  <!-- 语言切换 -->
  <div class="locale-switcher">
    <a href="/?locale=zh-CN">简体中文</a>
    <a href="/?locale=en">English</a>
  </div>
</nav>
```

### 7.3 多语言处理中间件

```Go

// actions/middleware/i18n.go
package middleware

import (
 "net/http"

 "github.com/gobuffalo/buffalo"
 "github.com/gobuffalo/mw-i18n/v2"
 "github.com/gobuffalo/packr/v2"
)

// I18n 中间件初始化
func I18n() buffalo.MiddlewareFunc {
 // 加载语言文件
 box := packr.New("locales", "../locales")
 translator, err := i18n.New(box, "zh-CN")
 if err != nil {
 panic(err)
 }

 // 返回中间件
 return translator.Middleware()
}

// SetLocale 设置当前语言
func SetLocale(next buffalo.Handler) buffalo.Handler {
 return func(c buffalo.Context) error {
 // 从 URL 参数获取语言
 locale := c.Param("locale")
 if locale == "" {
 // 从 Cookie 获取
 cookie, err := c.Request().Cookie("locale")
 if err == nil {
 locale = cookie.Value
 }
 }
 if locale == "" {
 // 默认语言
 locale = "zh-CN"
 }

 // 设置上下文语言
 c.Set("locale", locale)
 i18n.SetLocale(c, locale)

 // 保存到 Cookie
 c.SetCookie(&http.Cookie{
 Name:  "locale",
 Value: locale,
 Path:  "/",
 })

 return next(c)
 }
}
```

## 8. 权限系统

### 8.1 角色定义

|角色|权限|
|---|---|
|管理员|所有操作（用户管理、主题管理、系统设置、内容管理）|
|编辑者|内容管理（文章/分类/标签/媒体）、个人信息修改|
### 8.2 权限中间件

```Go

// actions/middleware/auth.go
package middleware

import (
 "fmt"
 "net/http"

 "github.com/gobuffalo/buffalo"
 "github.com/gobuffalo/pop/v6"
 "github.com/gofrs/uuid"

 "catter-cms/models"
)

// Auth 登录验证中间件
func Auth(next buffalo.Handler) buffalo.Handler {
 return func(c buffalo.Context) error {
 // 从 Session 获取用户 ID
 userIDStr := c.Session().Get("user_id")
 if userIDStr == nil {
 // 重定向到登录页
 return c.Redirect(http.StatusSeeOther, "/admin/login")
 }

 // 查询用户
 userID, err := uuid.FromString(userIDStr.(string))
 if err != nil {
 return c.Redirect(http.StatusSeeOther, "/admin/login")
 }

 tx := c.Value("tx").(*pop.Connection)
 user := &models.User{}
 if err := tx.Find(user, userID); err != nil {
 return c.Redirect(http.StatusSeeOther, "/admin/login")
 }

 // 设置当前用户到上下文
 c.Set("current_user", user)
 return next(c)
 }
}

// AdminOnly 管理员权限中间件
func AdminOnly(next buffalo.Handler) buffalo.Handler {
 return func(c buffalo.Context) error {
 user := c.Value("current_user").(*models.User)
 if user.Role != "admin" {
 return c.Error(http.StatusForbidden, fmt.Errorf("权限不足"))
 }
 return next(c)
 }
}
```

## 9. 增量静态生成

### 9.1 静态生成任务

```Go

// tasks/generate_static.go
package tasks

import (
 "html/template"
 "os"
 "path/filepath"

 "github.com/gobuffalo/buffalo"
 "github.com/gobuffalo/pop/v6"

 "catter-cms/models"
)

// GenerateStaticSite 全量生成静态站点
func GenerateStaticSite() error {
 // 初始化应用上下文
 app := buffalo.New(buffalo.Options{Env: "production"})
 tx, err := pop.Connect("production")
 if err != nil {
 return err
 }
 defer tx.Close()

 // 1. 获取激活的主题
 activeTheme := getActiveTheme(tx)

 // 2. 生成首页
 if err := generateHomePage(app, tx, activeTheme); err != nil {
 return err
 }

 // 3. 生成所有已发布文章
 posts := &models.Posts{}
 if err := tx.Where("status = ?", "published").All(posts); err != nil {
 return err
 }
 for _, post := range *posts {
 if err := generatePostPage(app, tx, activeTheme, post); err != nil {
 return err
 }
 }

 // 4. 生成分类/标签页面
 if err := generateCategoryPages(app, tx, activeTheme); err != nil {
 return err
 }

 return nil
}

// RevalidatePost 增量更新单篇文章
func RevalidatePost(slug string) error {
 // 仅重新生成该文章页面，不影响其他页面
 app := buffalo.New(buffalo.Options{Env: "production"})
 tx, err := pop.Connect("production")
 if err != nil {
 return err
 }
 defer tx.Close()

 activeTheme := getActiveTheme(tx)
 post := &models.Post{}
 if err := tx.Where("slug = ?", slug).First(post); err != nil {
 return err
 }

 return generatePostPage(app, tx, activeTheme, *post)
}

// getActiveTheme 获取激活的主题
func getActiveTheme(tx *pop.Connection) string {
 setting := &models.Setting{}
 if err := tx.Where("key = ?", "active_theme").First(setting); err != nil {
 return "blog" // 默认主题
 }
 return setting.Value
}

// generatePostPage 生成单篇文章静态页面
func generatePostPage(app *buffalo.App, tx *pop.Connection, theme string, post models.Post) error {
 // 1. 准备模板数据
 data := map[string]interface{}{
 "post":  post,
 "theme": theme,
 "site":  getSiteSettings(tx),
 }

 // 2. 渲染模板
 tmplPath := filepath.Join("templates", "themes", theme, "pages", "blog_show.plush.html")
 tmpl, err := template.ParseFiles(tmplPath)
 if err != nil {
 return err
 }

 // 3. 创建输出目录
 outputDir := filepath.Join("public", "static", "blog")
 if err := os.MkdirAll(outputDir, 0755); err != nil {
 return err
 }

 // 4. 写入静态文件
 outputPath := filepath.Join(outputDir, post.Slug+".html")
 file, err := os.Create(outputPath)
 if err != nil {
 return err
 }
 defer file.Close()

 return tmpl.Execute(file, data)
}

// generateHomePage 生成首页静态页面
func generateHomePage(app *buffalo.App, tx *pop.Connection, theme string) error {
 // 实现首页生成逻辑
 return nil
}

// generateCategoryPages 生成分类页面
func generateCategoryPages(app *buffalo.App, tx *pop.Connection, theme string) error {
 // 实现分类页面生成逻辑
 return nil
}

// getSiteSettings 获取站点设置
func getSiteSettings(tx *pop.Connection) map[string]interface{} {
 // 从数据库查询站点设置并返回
 return map[string]interface{}{
 "title":       "Catter CMS",
 "description": "轻量级内容管理系统",
 }
}
```

## 10. 部署

### 10.1 开发环境

```Bash

# 安装依赖
go mod download

# 初始化数据库
buffalo db create
buffalo db migrate
buffalo db seed

# 启动开发服务器
buffalo dev
```

### 10.2 生产环境

```Bash

# 编译二进制文件
buffalo build --production

# 部署步骤
1. 上传编译后的二进制文件（catter-cms）到服务器
2. 配置环境变量（数据库、密钥等）
3. 执行数据库迁移：./catter-cms db migrate
4. 启动服务：./catter-cms serve -port 8080

# 可选：使用 systemd 管理服务
# /etc/systemd/system/catter-cms.service
[Unit]
Description=Catter CMS
After=network.target

[Service]
User=www-data
WorkingDirectory=/opt/catter-cms
Environment="GO_ENV=production"
Environment="DATABASE_URL=postgresql://user:pass@localhost:5432/catter"
ExecStart=/opt/catter-cms/catter-cms serve -port 8080
Restart=always

[Install]
WantedBy=multi-user.target
```

### 10.3 数据库切换

```Bash

# 开发环境（SQLite）
export GO_ENV=development
buffalo db create

# 生产环境（PostgreSQL）
export GO_ENV=production
export DATABASE_NAME=catter
export DATABASE_USER=postgres
export DATABASE_PASSWORD=secret
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
buffalo db create
buffalo db migrate
```

## 11. 开发规范与约束

### 11.1 代码规范

- 严格遵循 Go 官方编码规范（`gofmt` 格式化）

- 所有模型必须实现 Validate 方法，确保数据合法性

- 所有数据库操作必须通过 Pop ORM，禁止直接执行 SQL

- 权限验证必须通过中间件实现，禁止在 Handler 中硬编码

- 主题模板必须遵循统一的命名规范，确保可切换性

### 11.2 安全规范

- 用户密码必须通过 bcrypt 加密，禁止明文存储

- 所有用户输入必须验证，防止 SQL 注入、XSS 等攻击

- API 接口必须验证权限，公开接口需做频率限制

- 上传文件必须验证类型和大小，存储路径需安全处理

- 敏感配置（密钥、数据库密码）必须通过环境变量注入

### 11.3 架构约束

- 禁止随意修改核心目录结构

- 禁止使用 Buffalo 未支持的第三方 ORM

- 主题开发必须遵循隔离原则，禁止修改核心代码

- 多语言实现必须基于内置 i18n 模块，禁止自定义方案

- 静态生成逻辑必须支持增量更新，禁止每次全量重建

---

*文档版本: 1.0*

*最后更新: 2026-03-10*

### 总结

1. Catter CMS 基于 Buffalo（Go）构建全栈 SSR 架构，核心采用 Go 语言开发，前端通过 Plush 模板引擎 + 轻量 JS 实现，无需 Node.js 环境；

2. 项目采用 MVC 架构，核心分为 actions（控制器）、models（数据模型）、templates（模板）三层，支持可插拔主题、多语言、增量静态生成等核心 CMS 能力；

3. 部署形态为单一 Go 二进制文件，开发环境使用 SQLite，生产环境可无缝切换 PostgreSQL，内存占用极低且性能高效。
> （注：文档部分内容可能由 AI 生成）