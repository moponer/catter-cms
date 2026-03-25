# 长期记忆

## 项目架构
- 后端入口：`gin-app/cmd/main.go`，不是 `bootstrap/app.go`
- 服务端口：config.yaml 中 `port: "3002"`
- 主题路径：`./templates/themes/`
- 构建命令：后端 `go build -o catter-cms.exe ./cmd`，前端 `npm run build`

## 技术栈
- 后端：Gin (Go 1.25)，GORM (SQLite)，主题系统支持嵌套配置存储
- 前端：Vue 3 + Vite + Element Plus，管理后台在 `/admin/` 路径下
- SSR：Go 标准模板引擎（html/template），所有前台 handler 调用 `services.RenderTemplate()` 使用 `GetActiveTheme()`

## 系统设置架构
- **数据表**：`settings` 表，key-value 结构（JSON 编码）
- **后端读取**：`models.GetSettingAsString(key)` 实时从数据库读取
- **前端保存**：`SiteSettings.vue` 发送扁平结构到 `/api/settings`
- **前台渲染**：`getSiteData()` 返回 `{title, description, url, logo, favicon, year}`
- **模板访问**：`{{ .site.title }}`, `{{ .site.description }}` 等
- **关键修复（2026-03-25）**：前端字段名 `site_name` → `site_title`，与后端一致

## 静态资源路径规范
- **FontAwesome**：`/fontawesome/css/all.min.css`（不是 `/themes/fontawesome/`）
- **主题CSS/JS**：`/themes/{theme-name}/assets/...`
- **Google Fonts**：已移除，改用系统字体栈（国内访问优化）

## 主题系统核心逻辑
- `services.GetActiveTheme()` = 当前激活主题（内存），SSR 渲染只用它
- `services.GetAllThemes()` = 所有已加载主题（管理用）
- 主题配置存储：`models.SaveThemeConfig` → `theme_configs` 表，JSON 格式，**纯扁平格式**（dot-notation key）
- 关键修复（2026-03-22）：`LoadTheme` 中添加 `FlattenMap()` 将 YAML 嵌套 config 转为扁平格式，确保 `theme.Config` 始终为纯扁平格式
- API `/api/themes/:slug/config/full` 返回 `{code:0, data: {name, slug, version, configForm[], ...其他扁平key配置}}`
- 前端 ThemeSettings.vue 直接使用 `fullData[field.name]` 访问扁平值，无需 `getNestedValue/flattenToNested`

## 主题优化完成状态（2026-03-23）
- **Catter Blog**: CSS变量化
- **Catter CMS Official**: CSS变量化
- **Enterprise Website**: theme.yml精简

## Phase 1 & 2 完成（2026-03-23）
**核心架构**：
- themeConfig/config 两种访问方式兼容
- siteName 作为 site.title 别名
- 新增模板函数：json_decode、substr、mod、div

**Blog 主题完整配置项**：
- 侧边栏位置（left/right）
- 文章列表样式（list/card/magazine）
- 文章列表模式（all/category）
- 标签云样式（flat/colorful/heat）
- 热门文章开关
- 阅读量/阅读时间开关

**Enterprise 主题完整配置项**：
- Hero幻灯片配置（多图轮播/单图横幅/视频背景）
- 模块开关（Stats/Services/Features/News/Partners/CTA）
- 页面模板选择（default/services/contact）

## Phase 6 SSR性能与安全优化（2026-03-24）
**安全优化**：
- 新增 SecurityHeaders 中间件：X-Content-Type-Options、X-Frame-Options、X-XSS-Protection 等
- 启用 Gzip 压缩中间件（之前存在但未启用）
- Referrer-Policy 安全头
- API 健康检查添加 no-cache 头

**性能优化**：
- 模板预加载机制：主题激活时异步预加载 home/article_show/category/tag 等常用模板
- 重构 getTemplateFuncsForTheme 支持传入主题参数
- 模板缓存自动管理

**前端配置界面**：
- ThemeSettings.vue: 支持 show_if 条件显示（支持 && 和 || 运算符）
- FieldRenderer.vue: 支持14种字段类型
- RepeaterField.vue: 支持 repeater 嵌套和拖拽排序

## 管理后台静态文件服务（2026-03-22 修复）
- `r.Static("/admin", adminDistPath)` 会拦截所有 `/admin/*` 路径，包括不存在的目录（如 `/admin/themes`），导致 NoRoute 永远不会被触发
- 解决方案：使用 `r.Group("/admin")` + 中间件，只允许磁盘上存在的文件通过，否则让 NoRoute fallback 到 `index.html`
- 中间件逻辑：`filePath := adminDistPath + c.Request.URL.Path[len("/admin"):]` → 如果文件不存在则 `c.Abort()` 让 NoRoute 处理

## 后端日志输出
- Go `log` 包输出重定向到 `os.Stderr`，避免与 HTTP 响应混在一起
- Gin 强制使用 `gin.ReleaseMode`，避免 debug 日志干扰

## 管理后台路由
- `/themes` → ThemeList（主题列表）
- `/themes/:slug` → ThemeSettings（主题配置）
- `/themes/:slug/modules` → ModuleManager（模块管理）

## 主题配置界面已知 Bug（2026-03-23 修复）
- ThemeSettings.vue Tab 标签显示为空：`:label="section.label"` → `:label="section.label || section.section"`，因为 theme.yml 的 configDefinition 每个 section 只有 `section` 字段作为分组名，没有独立的 `label`
- FieldRenderer.vue 不支持 `code` 类型：新增 `v-else-if="type === 'code'"` 分支，用 el-input textarea + 暗色等宽字体样式渲染

## API 基础路径
- 后端端口 3002
- 前端 dev server: 5173
- 主题 API 前缀：`/api/themes/`

## 主题系统架构（2026-03-24 重构）
**配置文件分离**：
- `theme.yml` = 元数据（name, slug, version, author, description）
- `settings.yaml` = Halo风格表单定义（使用 FormKit Schema 格式）
- `assets/css/base.css` = 共享基础样式

**模块组件化**（catter-blog）：
- `modules/widget-categories.html` - 分类小部件
- `modules/widget-tags.html` - 标签云小部件
- `modules/widget-recent-posts.html` - 最新文章小部件
- `modules/widget-profile.html` - 站点资料小部件
- `modules/post-card.html` - 文章卡片组件

**注意**：Go 模板引擎不支持跨文件 template 调用，模块组件代码需内联到使用模板中。

**当前主题状态**：
- catter-blog: 首页和分类页已重构，使用共享CSS
- catter-cms-official: 企业主题，样式完善无需大改

**数据库问题修复**：
- site_title 和 site_description 值被错误设置为 `" catter-blog\"`（带引号和反斜杠）
- 修复：`UPDATE settings SET value='Catter Blog' WHERE key='site_title';`

## 模板预加载架构重构（2026-03-24 完成）
**问题**：之前模板渲染失败，原因是 `template home not found in theme catter-blog`

**解决方案**：
- 添加 `TemplateSet` 类型存储组合模板：`type TemplateSet struct { Root *template.Template; Template map[string]*template.Template }`
- 修改 `preloadAllTemplates` 函数：正确从文件路径提取模板名（`pages/home.html` → `home`）
- 添加 `dateFormat` 模板函数支持日期格式化

**修复记录**：
- 移除模板中不存在的字段引用：Cover、ReadTime、PostCount（PostCount是方法不是字段）
- 简化 footer links 移除不存在的 dict 函数

**当前主题状态**：
- catter-blog: 全新首页模板 + CSS，支持文章列表、标签云、最新文章、页脚
- 模板渲染正常工作，预加载5个模板成功
- CSS 风格已更新为参考项目的白底简洁风格（蓝色主色调 #2563eb）

## 模板函数列表（2026-03-24）
getConfig, len, asset_path, current_theme, safe, default, lower, upper, trim, gt, lt, eq, add, sub, iterate, markdownHTML, json_decode, substr, mod, first, int, div, dateFormat

## 主题架构规范化（2026-03-25）✅
**问题**：主题存在大量内联样式、不规范的 CSS 变量使用、settings.yaml 冗余

**修复内容**：
1. **dianyan-blog 主题**：
   - 移除 10+ 处内联样式，改用 CSS 类
   - 新增 CSS 类：`category-link`, `main-full`, `main-padded`, `full-width`, `span-full`, `categories-grid`, `category-card`, `tags-cloud-container`, `tag-item`, `archive-year`, `archive-list`, `archive-item`, `archive-date`, `archive-title`, `search-hint`
   - 删除冗余 `settings.yaml`（系统使用 `theme.yml` 的 `configDefinition`）

2. **catter-cms-official 主题**：
   - 移除 8 处硬编码颜色的内联样式
   - 新增 CSS 类：`empty-state`, `empty-state-full`, `search-form`, `search-input`, `search-btn`, `results-info`, `categories-grid`, `category-card` 系列, `archive-*` 系列
   - CSS 变量默认值统一放在 `base.css`
   - 颜色配置 key 规范化：`primaryDark`, `primaryDarker`, `accent`, `bgPrimary`, `textPrimary` 等

**⚠️ Go 模板限制（重要）**：
- 不支持 `{{ template "_head" . }}` 跨文件模板包含，会导致 `no such template "_head"` 错误
- 每个页面模板必须包含完整的 HTML 结构（head/body/footer）
- CSS/JS 通过 `<link>` 和 `<script>` 引用，不要用 template include

**CSS 样式规范**：
- ✅ 统一封装到 CSS 文件（base.css/home.css/list.css/article.css）
- ❌ 禁止在 HTML 中使用 `style="..."` 内联样式属性
- ❌ 禁止在 HTML 中嵌入 `<style>` 块（CSS 变量除外）

**CSS 变量命名规范**：
- CSS 变量名（如 `--primary-dark`）与 theme.yml 中 configDefinition 的 name（如 `primaryDark`）必须一致
- CSS 变量默认值放在 `base.css` 的 `:root` 中

**主题配置加载机制**：
- 后端 API `/api/themes/:slug/config/full` 返回 `configForm`（来自 theme.yml 的 configDefinition）和已保存配置值
- 前端 ThemeSettings.vue 通过 `getThemeConfigFull()` 加载并渲染配置表单
- `settings.yaml` 已废弃，所有配置定义必须在 `theme.yml` 的 `configDefinition` 中

**文档更新**：
- 创建 `docs/THEME_DEVELOPER_GUIDE_V3.md`，整合今日教训：CSS 样式规范、Go 模板限制、settings.yaml 废弃说明

## 参考项目样式要点（2026-03-24）
**配色方案**：
- primary/accent: #2563eb (蓝色)
- bg-primary: #f8fafc (浅灰白)
- bg-secondary: #ffffff (纯白)
- text-primary: #1e293b
- text-secondary: #64748b

**布局结构**：
- Header: 白色背景 sticky 导航栏
- Page Header: 页面大标题区域
- Posts Section: 文章列表
- Footer: 蓝色背景页脚

**CSS 变量**：
- --primary, --primary-hover
- --bg-primary, --bg-secondary
- --text-primary, --text-secondary, --text-tertiary
- --border-color, --accent

## 模板配置 Key 修复（2026-03-25）
**问题：** 主题配置不生效
- `postListStyle`：模板中错误使用 `listStyle`，已修复
- `sidebarPosition`：模板中未使用，已添加动态 class `sidebar-left`
- 站点名称：统一使用系统设置 `{{ .site.title }}`，主题配置中的 `siteName` 已移除（方案B：站点名称全局统一）

**修复位置：** `templates/themes/dianyan-blog/pages/home.html`, `theme.yml`

## 模板语法修复（2026-03-25）
**问题1：** dianyan-blog 主题 `{{ . Excerpt }}` 中 dot 和字段名之间有空格
- 原因：Go 模板将 `{{ . Excerpt }}` 解析为调用函数 `Excerpt`，而非访问字段
- 修复：改为 `{{ .Excerpt }}`（无空格）
- 日志错误：`function "Excerpt" not defined`

**问题2：** 模板引用 Post 不存在的字段 `.Cover`
- 原因：Post 模型没有 Cover 字段（封面图）
- 修复：从模板中移除该引用
- 日志错误：`can't evaluate field Cover in type models.Post`

**问题3：** ThemeSettings.vue 重复导入 `onUnmounted`
- 修复：删除第177行的重复 import

**问题4：** `.Category.URL` 等不存在的 URL 字段
- 原因：Category/Tag/Post 模型没有 URL 字段，需要在模板中拼接
- 修复：使用 `/category/{{ .Category.Slug }}` 形式拼接

**问题5：** `in` 函数参数顺序错误
- 原因：Go 模板中 `in` 函数是 `in item slice`，不是 `in slice item`
- 修复：使用 `in "profile" $widgets` 而非 `in $widgets "profile"`

**问题6：** 新增 `ensureSlice` 模板函数
- 位置：`gin-app/internal/services/theme.go`
- 用途：安全处理配置值可能是 `[]interface{}` 或 JSON 字符串的情况
- 使用：`{{ $widgets := ensureSlice (getConfig "sidebarWidgets" `["profile"]`) }}`

## dianyan-blog 主题修复（2026-03-25）
- 修复 footer 三栏布局 CSS
- 实现侧边栏模块配置（profile/categories/tags）
- 完善标签云样式
- 移除内联 CSS 变量（违反规范）
- 移除 recent-posts 模块（后端未提供数据）

## 文档更新
- 创建 `docs/BUG_FIXES.md`（问题修复记录手册）
- 更新 `docs/THEME_DEVELOPER_GUIDE_V3.md` 至 v3.1

## 主题配置界面优化（2026-03-25）
**颜色字段改进**：
- FieldRenderer.vue: 颜色输入框移至左侧，颜色选择器移至右侧作为辅助
- 输入框显示当前颜色值（带颜色预览色块）
- 颜色选择器用于可视化选择，选择后自动填入输入框
- 右侧显示默认值提示（如"默认: #2d5a4a"）
- 最终生效颜色值以输入框内的为准

**文章图片遮罩优化**：
- base.css: 给文章封面图添加渐变遮罩层
- 遮罩从左上到右下的微透明渐变，改善低质量图片的模糊感
- 占位图（picsum.photos）添加饱和度降低滤镜，与遮罩协调

**修改文件**：
- `admin/src/components/FieldRenderer.vue`: 颜色字段模板重构、添加 defaultValue prop
- `admin/src/views/Themes/ThemeSettings.vue`: 传递 defaultValue 给 FieldRenderer
- `templates/themes/dianyan-blog/assets/css/base.css`: 文章封面图遮罩样式

**博主信息模块（2026-03-25）**：
- 侧边栏"站点资料"改为"博主信息"卡片
- 新增配置项：authorName（博主名称）、authorBio（博主简介）、authorAvatar（博主头像）
- 新增社交链接：socialGithub、socialWeibo、socialWechat、socialEmail
- 博主头像默认使用 picsum.photos 随机头像
- 社交链接显示为图标（SVG 格式），悬停有动效
- 移除文章详情页侧边栏的"最新文章"小部件

**修改文件**：
- `templates/themes/dianyan-blog/pages/home.html`: 侧边栏博主信息模块
- `templates/themes/dianyan-blog/pages/article_show.html`: 同步更新博主信息模块
- `templates/themes/dianyan-blog/theme.yml`: 添加博主信息配置项
- `templates/themes/dianyan-blog/assets/css/base.css`: 博主信息卡片样式

## 系统设置数据流修复（2026-03-25）
**问题1**：前端字段名与后端 key 不匹配
- 前端错误使用 `site_name`，后端 key 是 `site_title`
- 修复：SiteSettings.vue 改用 `site_title`

**问题2**：前端发送嵌套结构
- 错误：`{ category: 'basic', settings: { ... } }`
- 正确：扁平对象 `{ site_title, site_description, ... }`

**问题3**：数据库值损坏
- settings 表部分值包含多余引号和反斜杠
- 修复：删除损坏记录，重新保存

**数据流验证**：前端字段名 → API key → settings 表 → getSiteData() → 模板变量 `{{ .site.title }}`

**修改文件**：
- `admin/src/views/settings/SiteSettings.vue`

**文档更新**：
- `docs/BUG_FIXES.md` 更新至 v1.1，新增问题 12-16
