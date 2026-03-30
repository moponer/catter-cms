# 长期记忆 - Catter CMS 项目

## 项目架构
- **后端**：Gin (Go)，GORM (SQLite)，入口 `gin-app/cmd/main.go`
- **前端**：Vue 3 + Vite + Element Plus，管理后台在 `/admin/`
- **SSR**：Go 标准模板引擎（html/template），`services.RenderTemplate()` 渲染前台
- **服务端口**：`config.yaml` 中 `port: "3000"`
- **主题路径**：`./templates/themes/`
- **构建命令**：后端 `go build -o catter-cms.exe ./cmd`，前端 `npm run build`

## 主题系统核心
- 主题配置存储：扁平格式，`theme_configs` 表（JSON）
- 配置加载：`LoadTheme` → `FlattenMap()` → `GetThemeConfig()` → `DeepMerge()`
- API：`/api/themes/:slug/config/full` 返回 `{configForm[], ...扁平key配置}`
- **`settings.yaml` 已废弃**，配置定义必须在 `theme.yml` 的 `configDefinition` 中
- 模板函数：`getConfig, asset_path, markdownHTML, json_decode, ensureSlice, dateFormat` 等

## 各页面真实数据字段（已验证源码）

```
home.html         → { site, config, pageTitle, posts[], recentPosts[], categories[], tags[] }
article_show.html → { site, config, pageTitle, post, categories[], tags[], recentPosts[] }
article_list.html → { site, config, pageTitle, posts[], categories[], tags[], pagination{} }
category.html     → { site, config, pageTitle, category, posts[], categories[], tags[] }
tag.html          → { site, config, pageTitle, tag, posts[] }
categories.html   → { site, config, pageTitle, categories[] }  含 PostCount
tags.html         → { site, config, pageTitle, tags[] }        含 PostCount
archives.html     → { site, config, pageTitle, archives{年月:posts[]}, posts[] }
search.html       → { site, config, pageTitle, query, posts[], count }
page.html         → { site, config, pageTitle, page, recentPosts[] }
```
- **分页字段名是 `pagination`**（不是 `paginator`），结构：`{page, pageSize, total, totalPages}`
- **仅 article_list.html 有分页**，category/tag 返回全量
- `page.Template` 不为空时渲染 `page-{template}.html`，找不到则回退 `page.html`
- `archives` 是 Go map，遍历无序（`range $yearMonth, $posts := .archives`）

## 关键 Go 模板约束
1. **不支持跨文件 template include**（`{{ template "_head" . }}` 会报错）
2. **每页必须包含完整 HTML 结构**（DOCTYPE、head、body）
3. `.字段名` 中间不能有空格（`{{ . Excerpt }}` 错误）
4. 比较用函数 `eq/gt/lt/ge/le/ne`
5. Post 模型有 Cover 字段，无 ReadTime、URL 字段
6. `in` 函数参数顺序：`in item slice`
7. range 循环取索引：`{{ range $index, $item := .posts }}`
8. **`markdownHTML` 已返回 `template.HTML`，禁止再接 `| safe`**
9. **Page 模型无 Description 字段**，用 Excerpt 替代
10. `.container` 必须设置 `width: 100%` 否则 flex `space-between` 不生效

## CSS 规范
- 禁止内联 `style="..."` 属性
- 禁止 HTML 中嵌入 `<style>` 块（CSS 变量初始化除外）
- 所有颜色必须使用 CSS 变量
- CSS 变量默认值在 `base.css` 的 `:root` 中定义
- CSS 变量名（`--primary-dark`）与 configDefinition 的 name（`primaryDark`）必须一致

## 静态资源路径规范
- FontAwesome：`/fontawesome/css/all.min.css`
- 主题 CSS/JS：`{{ asset_path "assets/css/base.css" }}`（不要带 `/themes/xxx/`）
- `asset_path` 函数会自动添加 `/themes/{themeName}/` 前缀，传入相对路径即可
- 错误示例：`{{ asset_path "/themes/catter-cms-official/assets/css/base.css" }}` → 会变成 `/themes/catter-cms-official//themes/catter-cms-official/...`
- 禁止 Google Fonts（国内问题），改用系统字体栈
- **禁止使用 Emoji 作为 UI 图标**，必须用 FontAwesome `<i class="fa-solid fa-xxx"></i>`

## ⛔ 文件写入安全规范（高优先级）
- **绝对禁止用 PowerShell 修改 `theme.yml` / CSS / 任何配置文件**（双 BOM 问题）
- 修改现有文件 → 用 `replace_in_file` 工具
- 必须全量覆写 → 用 Python：`open(..., encoding='utf-8-sig')` 读 + `open(..., encoding='utf-8', newline='\n')` 写

## 系统设置 & API
- 字段名：`site_title`（非 `site_name`）、`site_description`
- API 响应：`{code: 0, message, data}` / `{code: 1, message, data: null}`

## 页面模板数据架构（2026-03-27）
- `Page` 模型有 `TemplateData` JSON 字段
- `theme.yml` 的 `pageTemplates` 定义模板类型和字段
- `PageEdit.vue` 根据模板类型动态渲染编辑组件
- 前台 `pageShowHandler` 解析 `template_data` 为 `.td` 对象传入模板

## 已开发主题列表
- **catter-blog**：蓝色简洁博客，侧边栏、标签云
- **catter-cms-official**：企业官网风格，深蓝色
- **enterprise-website**：商务专业风格
- **dianyan-blog**：淡雅博客，博主信息侧边栏
- **duocai-blog**：多彩博客，Hero 背景轮播，image-list 字段
- **shengyi-logistics**：盛意物流，危化品运输企业站

## 前端组件
- `FieldRenderer.vue`：14 种字段类型
- `ImageUploadField.vue` / `ImageListUploadField.vue`：图片上传
- `RepeaterField.vue`：可重复字段，拖拽排序
- `ThemeSettings.vue`：主题配置，支持 show_if 条件显示

## 安全与架构
- SPA 中间件：`middleware/spa.go`
- XSS 防护：`html.EscapeString()`
- 安全头：SecurityHeaders 中间件
- Gzip 压缩已启用

## catter-cms-themes Skill
- **位置**：`f:\Catter-CMS\.workbuddy\skills\catter-cms-themes\`（项目级）
- **zip 包**：`f:\Catter-CMS\.workbuddy\skills\catter-cms-themes.zip`
- **文件**：SKILL.md + references/（template-examples、css-base-template、theme-yml-examples、bug-patterns）
- **Bug 模式**：B001-B021，最常命中 #1 getConfig 默认值、#8 markdownHTML|safe、#8 Emoji 禁止
- **触发时机**：开发、修改、调试任何 Catter CMS 主题时加载

## 2026-03-30 修复记录

### shengyi-logistics 主题问题修复
1. **Hero 模块不显示**：根因是 `homeModuleOrder` 数据库值缺少 hero，前端保存配置时会覆盖。修复：更新数据库 + ThemeSettings.vue 加载时自动补充缺失的 hero 模块
2. **服务项目按钮不对齐**：修改 home.css，`.service-card` 改 flex 纵向布局，按钮 `margin-top: auto` 固定底部
3. **首页标题多出"首页"**：home.html 标题改为只显示 `.site.title`
4. **模板 unexpected EOF**：删除 home.html 第 520 行重复的 `{{ if getConfig "statsEnabled" true }}`
5. **空模块列表时 Hero 不显示**：home.html `{{ if $modules }}` 改为 `{{ if gt (len $modules) 0 }}`

### theme.yml 规范
- `homeModuleOrder` 默认值必须包含 hero 和 news 模块，否则新激活主题时对应模块不会显示
- ThemeSettings.vue 加载时自动补充缺失的 hero 和 news 模块（向后兼容旧数据）

### shengyi-logistics 新闻动态模块不显示
6. **新闻动态模块不显示**：根因是 `homeModuleOrder` 默认值和数据库存储都缺少 `news` 模块。修复：
   - theme.yml `homeModuleOrder` default 添加 `{id: "news", name: "新闻动态", dataKey: "newsItems"}`
   - theme.yml 添加 `newsEnabled`（switch）、`newsTitle`（text）、`newsShowDate`（switch）配置定义
   - ThemeSettings.vue 加载时自动补充缺失的 news 模块（与 hero 同理）
7. **新闻模块文章不显示**：模板在 `{{ range $modules }}` 循环内使用 `.posts`/`.recentPosts`，但 `.` 指向模块项。修复：改为 `$.posts`/`$.recentPosts` 访问根数据
8. **文章列表链接 404**：模板使用 `/article-list` 但路由是 `/article`。修复：所有链接改为 `/article`，同时更新数据库 footerNav 配置
