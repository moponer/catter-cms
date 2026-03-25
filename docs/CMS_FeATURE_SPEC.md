# Catter CMS 功能规格说明书

**版本**: 1.0
**日期**: 2026-03-23
**状态**: 规划中

---

## 一、设计理念

### 1.1 核心目标
- **主题与数据解耦**: 主题只负责渲染，数据全部来自配置或数据库
- **模块化设计**: 每个页面版块都是独立模块，可配置、可排序、可开关
- **所见即所得的配置体验**: 后台配置实时预览，无需反复刷新
- **Halo 兼容性**: 参考 Halo 2.0 的主题架构设计

### 1.2 关键词
- **主题 (Theme)**: 决定网站视觉呈现的模板+样式+配置
- **模块 (Module)**: 页面中的可配置版块（Hero、Features、Sidebar 等）
- **模板 (Template)**: 决定页面结构的文件（home.html、article.html 等）
- **配置项 (Setting)**: 控制主题行为的可调参数

---

## 二、主题系统架构

### 2.1 目录结构

```
templates/themes/{theme-slug}/
├── theme.yaml           # 主题元信息
├── settings.yaml         # 配置项定义（参考 Halo）
├── pages/                # 页面模板
│   ├── home.html         # 首页
│   ├── article.html      # 文章详情
│   ├── article_list.html # 文章列表
│   ├── category.html     # 分类页
│   ├── tag.html          # 标签页
│   ├── page.html         # 单页面
│   ├── search.html       # 搜索页
│   └── archives.html     # 归档页
├── modules/              # 可复用模块（可选）
│   ├── hero.html
│   ├── footer.html
│   └── sidebar.html
├── styles/               # 样式文件
│   ├── base.css
│   ├── components.css
│   └── variables.css      # CSS 变量（配色、字体等）
├── assets/               # 静态资源
│   ├── images/
│   └── js/
└── screenshots/          # 主题截图
    └── screenshot.png
```

### 2.2 theme.yaml 元信息

```yaml
name: 主题名称
slug: theme-slug
version: 1.0.0
author: 作者
description: 主题描述
thumbnail: screenshot.png  # 缩略图

# 主题提供的模板列表
templates:
  - name: 首页
    file: home.html
    route: /
  - name: 文章详情
    file: article.html
    route: /article/:slug
  - name: 文章列表
    file: article_list.html
    route: /article
  - name: 分类页
    file: category.html
    route: /category/:slug
  - name: 标签页
    file: tag.html
    route: /tag/:slug
  - name: 单页面
    file: page.html
    route: /:slug
  - name: 搜索页
    file: search.html
    route: /search
  - name: 归档页
    file: archives.html
    route: /archives

# 页面模板列表（简化写法，等同于上面的 templates）
pageTemplates:
  - default        # 默认模板
  - full-width     # 全宽模板
  - blank          # 空白模板（仅内容）
  - landing        # 着陆页模板
```

### 2.3 页面类型与路由

| 页面类型 | 路由规则 | 说明 |
|---------|---------|------|
| 首页 | `/` | 网站主页 |
| 文章列表 | `/article` | 所有文章列表 |
| 文章详情 | `/article/:slug` | 单篇文章 |
| 分类页 | `/category/:slug` | 某分类下的文章 |
| 标签页 | `/tag/:slug` | 某标签下的文章 |
| 单页面 | `/:slug` | 独立页面（如关于、联系） |
| 搜索页 | `/search` | 搜索结果 |
| 归档页 | `/archives` | 按时间归档 |
| 作者页 | `/author/:name` | 作者文章列表 |

---

## 三、主题配置系统（settings.yaml）

### 3.1 配置结构（参考 Halo）

```yaml
apiVersion: v1alpha1
kind: Setting
metadata:
  name: {theme-slug}-setting
spec:
  forms:
    - group: basic
      label: 基本设置
      formSchema:
        - $formkit: text
          name: site_name
          label: 网站名称
          value: "我的网站"

    - group: style
      label: 样式设置
      formSchema:
        - $formkit: select
          name: color_scheme
          label: 配色方案
          value: "light"
          options:
            - label: 浅色模式
              value: "light"
            - label: 深色模式
              value: "dark"

    - group: layout
      label: 布局设置
      formSchema:
        - $formkit: radio
          name: sidebar_position
          label: 侧边栏位置
          value: "right"
          options:
            - label: 左侧
              value: "left"
            - label: 右侧
              value: "right"
            - label: 不显示
              value: "none"

    - group: seo
      label: SEO 设置
      formSchema:
        - $formkit: text
          name: meta_keywords
          label: 关键词
        - $formkit: textarea
          name: meta_description
          label: 描述
```

### 3.2 字段类型

| 类型 | 说明 | 适用场景 |
|-----|------|---------|
| text | 单行文本 | 标题、名称 |
| textarea | 多行文本 | 描述、备注 |
| number | 数字 | 数量、长度 |
| color | 颜色选择 | 主题色、背景色 |
| switch | 开关 | 启用/禁用功能 |
| select | 下拉选择 | 布局、样式选择 |
| radio | 单选 | 排列方式、显示位置 |
| checkbox | 多选 | 功能开关组合 |
| image | 图片上传 | Logo、Banner |
| date | 日期 | 活动时间 |
| datetime | 日期时间 | 定时发布 |
| daterange | 日期范围 | 筛选条件 |
| slider | 滑块 | 透明度、大小比例 |
| code | 代码编辑 | 自定义 CSS/JS |
| repeater | 重复器 | 列表类配置 |

### 3.3 条件显示（if）

```yaml
- $formkit: radio
  name: layout_type
  label: 布局类型
  value: "sidebar"

- $formkit: select
  name: sidebar_width
  label: 侧边栏宽度
  if: "$get(layout_type).value === 'sidebar'"
  options:
    - label: 窄版 (200px)
      value: "200px"
    - label: 标准 (250px)
      value: "250px"
    - label: 宽版 (300px)
      value: "300px"
```

---

## 四、Blog 主题配置项设计

### 4.1 分组设计

```
1. 基本设置
   - 网站名称/Logo/Favicon
   - 博主信息（昵称、头像、简介、社交链接）

2. 首页布局
   - Hero 区域开关/内容配置
   - 侧边栏位置（左侧/右侧/不显示）
   - 文章列表样式（纯文字/图文卡片/杂志风）
   - 文章列表模式（全部/按分类）
   - 每页显示数量
   - 是否显示摘要/字数/阅读时间

3. 侧边栏小工具
   - 博主信息模块开关
   - 分类列表模块开关
   - 标签云模块开关（纯色/彩色/标签数量）
   - 最新文章模块开关（数量）
   - 热门文章模块开关（数量/热度计算方式）
   - 广告位模块开关

4. 文章详情页
   - 相关文章模块开关（数量）
   - 评论区开关
   - 上一篇/下一篇导航开关
   - 文章目录（TOC）开关
   - 打赏模块开关

5. SEO 设置
   - 首页 Meta Title/Keywords/Description
   - 文章页 Meta Title 格式
   - Open Graph 图片

6. 样式定制
   - 主题色（CSS 变量）
   - 字体配置
   - 自定义 CSS
   - 自定义 JavaScript

7. 高级设置
   - 首页横幅图片/链接/文字
   - 评论区设置（需要审核/匿名评论）
   - 订阅设置（RSS/邮件订阅链接）
```

### 4.2 关键配置项详解

#### 4.2.1 文章列表样式

```yaml
# 三种模式可选
post_list_style:
  type: radio
  label: 文章列表样式
  value: "card"
  options:
    - label: 纯文字列表
      value: "list"
    - label: 图文卡片
      value: "card"
    - label: 杂志风（交错布局）
      value: "magazine"

# 杂志风特定配置
magazine_layout:
  type: select
  label: 交错方式
  if: "$get(post_list_style).value === 'magazine'"
  options:
    - label: 左右交错
      value: "alternate"
    - label: 左对齐瀑布流
      value: "masonry"
```

#### 4.2.2 侧边栏小工具

```yaml
sidebar_widgets:
  type: repeater
  label: 侧边栏模块
  value: []
  items:
    - name: widget_type
      type: select
      options:
        - label: 博主信息
          value: "author"
        - label: 分类列表
          value: "categories"
        - label: 标签云
          value: "tags"
        - label: 最新文章
          value: "recent_posts"
        - label: 热门文章
          value: "hot_posts"
        - label: 广告位
          value: "ads"

    - name: title
      type: text
      label: 标题

    - name: enabled
      type: switch
      label: 启用

    # 博主信息特定
    - name: bio
      type: textarea
      if: "$get(widget_type).value === 'author'"

    # 标签云特定
    - name: tag_style
      type: select
      if: "$get(widget_type).value === 'tags'"
      options:
        - label: 纯色
          value: "flat"
        - label: 彩色
          value: "colorful"
        - label: 按热度着色
          value: "heat"

    # 最新/热门文章特定
    - name: post_count
      type: number
      if: "$get(widget_type).value === 'recent_posts' || $get(widget_type).value === 'hot_posts'"
```

#### 4.2.3 首页 Hero 配置

```yaml
hero_enabled:
  type: switch
  label: 显示 Hero 区域

hero_style:
  type: select
  label: Hero 样式
  if: "hero_enabled === true"
  options:
    - label: 全屏背景图
      value: "fullscreen"
    - label: 全宽横幅
      value: "banner"
    - label: 简洁文字
      value: "minimal"

hero_fullscreen:
  type: group
  if: "hero_style === 'fullscreen'"
  fields:
    - name: background_image
      type: image
      label: 背景图片
    - name: overlay_color
      type: color
      label: 遮罩颜色
    - name: overlay_opacity
      type: slider
      label: 遮罩透明度
      min: 0
      max: 100
    - name: title
      type: text
      label: 主标题
    - name: subtitle
      type: textarea
      label: 副标题
    - name: button_text
      type: text
      label: 按钮文字
    - name: button_link
      type: text
      label: 按钮链接
```

---

## 五、Enterprise 主题配置项设计

### 5.1 分组设计

```
1. 基本设置
   - 公司名称/Logo/Favicon
   - 公司简介
   - 联系方式（电话/邮箱/地址）
   - 社交媒体链接

2. 导航设置
   - 顶部导航菜单
   - 移动端导航样式
   - Logo 位置（居左/居中/居右）

3. 首页模块（模块化版块，可拖拽排序）
   - Hero/Banner 区域
   - 服务项目区域
   - 产品展示区域
   - 统计数据区域
   - 客户案例/合作伙伴
   - 新闻动态区域
   - CTA 行动召唤区域
   - 底部 Footer 区域

4. Banner/Hero 配置
   - 幻灯片数量（单图/多图轮播）
   - 每张幻灯片配置：
     * 图片
     * 标题
     * 副标题
     * 按钮文字/链接
     * 是否启用

5. 服务项目配置
   - 每项服务：
     * 图标
     * 标题
     * 描述
     * 链接

6. 产品展示配置
   - 展示样式（网格/列表/轮播）
   - 每项产品：
     * 图片
     * 名称
     * 描述
     * 价格（可选）
     * 链接

7. 统计数据配置
   - 每项数据：
     * 数值
     * 标签
     * 前缀/后缀

8. 新闻动态配置
   - 显示数量
   - 每条新闻：
     * 标题
     * 日期
     * 描述
     * 链接
     * 缩略图

9. CTA 区域配置
   - 标题
   - 副标题
   - 按钮配置

10. Footer 配置
    - 公司信息
    - 快速链接
    - 分类链接
    - 社交媒体
    - 备案号
    - Copyright

11. 页面模板配置
    - 服务详情页模板
    - 产品详情页模板
    - 案例详情页模板
    - 联系我们页模板

12. SEO 设置
13. 样式定制
14. 高级设置
```

### 5.2 首页模块化配置

```yaml
homepage_modules:
  type: repeater
  label: 首页模块
  items:
    - name: module_type
      type: select
      options:
        - label: Hero 幻灯片
          value: "hero_slides"
        - label: 服务项目
          value: "services"
        - label: 产品展示
          value: "products"
        - label: 统计数据
          value: "stats"
        - label: 客户案例
          value: "case_studies"
        - label: 新闻动态
          value: "news"
        - label: CTA 区域
          value: "cta"
        - label: 合作伙伴
          value: "partners"
        - label: 自定义 HTML
          value: "custom_html"

    - name: enabled
      type: switch

    - name: title
      type: text
      label: 区域标题

    - name: subtitle
      type: text
      label: 区域副标题

    - name: order
      type: number
      label: 排序顺序
```

---

## 六、前端配置界面设计

### 6.1 ThemeSettings.vue 重构

参考 Halo 的 Console 设计：

```
┌─────────────────────────────────────────────────────────────┐
│  主题名称 - 主题配置                              [返回列表] │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────┬─────────┬─────────┬─────────┬─────────┐       │
│  │ 基本设置 │ 样式设置 │ 布局设置 │ SEO设置 │ 高级设置 │       │
│  └─────────┴─────────┴─────────┴─────────┴─────────┘       │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                                                      │    │
│  │  [表单字段...]                                       │    │
│  │                                                      │    │
│  │  条件显示示例：                                      │    │
│  │  ┌──────────────────┐                              │    │
│  │  │ 布局类型: ○左 ○右 │  ← 用户选择                  │    │
│  │  └──────────────────┘                              │    │
│  │                                                      │    │
│  │  ┌──────────────────┐                              │    │
│  │  │ 侧边栏宽度: [下拉] │  ← 只有选"右"时才显示       │    │
│  │  └──────────────────┘                              │    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    [重置]    [保存配置]              │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 6.2 核心功能

1. **Tab 分组**: 配置项按 group 分组为 Tab
2. **条件显示**: 支持 `show_if` 条件，动态显示/隐藏字段
3. **Repeater 字段**: 支持数组类型配置（features、services 等）
4. **实时预览**: 字段修改后可点击预览查看效果
5. **配置验证**: 表单提交前验证必填项和数据格式
6. **批量重置**: 一键恢复所有配置到默认值

### 6.3 Repeater 字段 UI

对于数组类型配置（services、features 等），使用可拖拽的卡片列表：

```
┌─────────────────────────────────────────────┐
│ 服务项目                              [+ 添加] │
├─────────────────────────────────────────────┤
│ ┌─────────────────────────────────────────┐ │
│ │ [拖拽图标] 服务 1                    [删除] │ │
│ │  ┌──────────────────────────────────┐   │ │
│ │  │ 图标: [cloud-upload        ]     │   │ │
│ │  │ 标题: [云服务器]                  │   │ │
│ │  │ 描述: [高性能云服务器...]         │   │ │
│ │  │ 链接: [/services/cloud]           │   │ │
│ │  └──────────────────────────────────┘   │ │
│ └─────────────────────────────────────────┘ │
│ ┌─────────────────────────────────────────┐ │
│ │ [拖拽图标] 服务 2                    [删除] │ │
│ │  ...                                     │ │
│ └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
```

---

## 七、SSR 性能优化

### 7.1 模板渲染优化

```go
// 当前问题：每次请求都解析模板文件
func RenderTemplate(name string, data map[string]interface{}) (string, error) {
    // 应该：模板文件解析后缓存到内存
}

// 优化方案：使用 template.ParseFiles + template.Execute 的缓存版本
```

### 7.2 缓存策略

| 资源类型 | 缓存策略 | 时长 |
|---------|---------|------|
| 主题模板文件 | 内存缓存 | 进程周期 |
| 主题配置 | 内存缓存 + Redis | 5 分钟 |
| 静态资源 | CDN/浏览器缓存 | 1 天 |
| 数据库查询 | Query Cache | 1 分钟 |

### 7.3 性能指标目标

| 指标 | 目标值 |
|-----|-------|
| TTFB (Time To First Byte) | < 200ms |
| First Contentful Paint | < 1s |
| Largest Contentful Paint | < 2.5s |
| Server Render Time | < 50ms |

---

## 八、安全考虑

### 8.1 模板安全

```go
// 1. 禁止执行任意代码
// Go html/template 默认自动转义，比 Jinja2 更安全

// 2. 防止 XSS
// 用户输入必须经过 html/template 的自动转义
// 富文本内容使用 markdownHTML 函数处理

// 3. 避免模板注入
// 不要使用 text/template，应该使用 html/template
```

### 8.2 配置安全

```yaml
# 配置项应该支持类型验证
# 例如：color 类型应该验证是否为有效颜色值
# number 类型应该验证范围
# url 类型应该验证格式
```

### 8.3 API 安全

| 风险 | 防护措施 |
|-----|---------|
| CSRF | Token 验证 |
| XSS | 输出转义 |
| SQL 注入 | ORM 参数化查询 |
| 文件上传 | 类型检查 + 路径限制 |
| 配置注入 | Schema 验证 |

---

## 九、实现优先级

### Phase 1: 核心架构 (P0)
1. 重构 theme.yaml + settings.yaml 结构
2. 添加模板函数：json_decode、substr、条件渲染支持
3. 实现主题配置的正确传递（themeConfig 到模板）
4. 修复前端配置界面样式问题

### Phase 2: Blog 主题完整配置 (P1)
1. 侧边栏位置配置（左右/不显示）
2. 文章列表样式（纯文字/图文/杂志风）
3. 侧边栏小工具模块化
4. SEO 配置

### Phase 3: Enterprise 主题完整配置 (P1)
1. 首页模块化版块系统
2. Banner/幻灯片配置
3. 服务/产品/新闻配置
4. CTA 区域配置

### Phase 4: 前端配置界面 (P2)
1. Repeater 字段支持（拖拽排序）
2. 条件显示逻辑完善
3. 配置预览功能
4. 表单验证

### Phase 5: 性能优化 (P2)
1. 模板缓存
2. 主题配置缓存
3. 静态资源优化

### Phase 6: 安全加固 (P2)
1. 配置 Schema 验证
2. 输入过滤
3. 安全头设置

---

## 十、参考实现

### 10.1 Halo settings.yaml 示例

```yaml
apiVersion: halo.run/v1alpha1
kind: ThemeSetting
metadata:
  name: theme-starter-setting
spec:
  forms:
    - group: basic
      label: 基本设置
      formSchema:
        - $formkit: text
          name: site_name
          label: 网站名称
          value: "My Site"
        - $formkit: image
          name: site_logo
          label: 网站 Logo

    - group: style
      label: 样式设置
      formSchema:
        - $formkit: color
          name: primary_color
          label: 主题色
          value: "#007bff"
        - $formkit: select
          name: font_family
          label: 字体
          value: "system-ui"
          options:
            - label: 系统字体
              value: "system-ui"
            - label: 思源黑体
              value: "Noto Sans SC"
```

### 10.2 主题配置传递

```go
// 后端：在 data 中注入 themeConfig
func homeHandler(c *gin.Context) {
    theme := services.GetActiveTheme()

    data := map[string]interface{}{
        "site":        siteData,
        "pageTitle":   "首页",
        "posts":       posts,
        // 关键：将主题配置注入模板
        "themeConfig": theme.Config,
        // 页面级配置（可选）
        "pageConfig":  getPageConfig(c),
    }

    services.RenderTemplate("home", data)
}
```

### 10.3 模板中使用配置

```html
<!-- 侧边栏位置 -->
<div class="layout {{ if eq .themeConfig.sidebarPosition "left" }}layout-sidebar-left{{ end }}">

<!-- 文章列表样式 -->
{{ if eq .themeConfig.postListStyle "card" }}
  {{ range .posts }}
  <article class="post-card">
    <img src="{{ .Thumbnail }}" />
    <h3>{{ .Title }}</h3>
  </article>
  {{ end }}
{{ else if eq .themeConfig.postListStyle "list" }}
  {{ range .posts }}
  <article class="post-item">
    <h3>{{ .Title }}</h3>
    <p>{{ .Excerpt }}</p>
  </article>
  {{ end }}
{{ end }}

<!-- Repeater 数组渲染 -->
{{ range (json_decode .themeConfig.features) }}
<div class="feature">
  <span class="icon">{{ .icon }}</span>
  <h4>{{ .title }}</h4>
  <p>{{ .desc }}</p>
</div>
{{ end }}
```

---

## 附录 A: 配置项检查清单

### Blog 主题必须实现

- [ ] 侧边栏位置（left/right/none）
- [ ] 文章列表样式（list/card/magazine）
- [ ] 文章列表模式（all/category）
- [ ] 每页数量
- [ ] 摘要/字数/阅读时间开关
- [ ] 博主信息模块
- [ ] 分类列表模块
- [ ] 标签云模块（flat/colorful/heat）
- [ ] 最新文章模块
- [ ] 热门文章模块
- [ ] Hero 区域配置
- [ ] SEO 配置

### Enterprise 主题必须实现

- [ ] Hero 幻灯片配置
- [ ] 服务项目配置
- [ ] 产品展示配置
- [ ] 统计数据配置
- [ ] 新闻动态配置
- [ ] CTA 区域配置
- [ ] Footer 配置
- [ ] 页面模板选择
- [ ] 联系方式配置

---

**文档版本历史**:
- v1.0 (2026-03-23): 初始版本
