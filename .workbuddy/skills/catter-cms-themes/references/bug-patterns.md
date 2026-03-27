# Catter CMS 已知 Bug 模式与修复速查

> 每次遇到问题，先查此表；每次发现新 Bug，更新此表。

---

## 快速索引

| Bug 类型 | 错误信息 | 跳转 |
|---------|---------|------|
| dot 空格问题 | `function "xxx" not defined` | [#B001](#b001) |
| 不存在的字段 | `can't evaluate field xxx` | [#B002](#b002) |
| URL 字段不存在 | `can't evaluate field URL` | [#B003](#b003) |
| 跨文件模板 | `no such template` | [#B004](#b004) |
| asset_path 路径重复 | MIME type error（路径含`//`）| [#B005](#b005) |
| in 函数参数顺序 | `wrong type for value` | [#B006](#b006) |
| range 循环取索引 | `cannot index slice/array` | [#B007](#b007) |
| CSS 变量不生效 | 颜色配置无变化 | [#B008](#b008) |
| FontAwesome 路径 | MIME type error | [#B009](#b009) |
| Google Fonts 加载慢 | 页面加载超时 | [#B010](#b010) |
| ensureSlice 缺失 | `wrong type for value` | [#B011](#b011) |
| 配置 key 不匹配 | 配置不生效（无报错） | [#B012](#b012) |
| 配置分组无标签 | 后台显示"未命名分组" | [#B013](#b013) |
| markdownHTML 未用 | 文章内容显示源码 | [#B014](#b014) |
| 动画关键帧缺失 | 轮播不工作 | [#B015](#b015) |
| 分页字段名错误 | 分页不显示 / nil 错误 | [#B016](#b016) |
| **markdownHTML \| safe** | **模板静默失败回退到 page.html** | [#B017](#b017) |
| **page.Description 不存在** | **模板渲染失败，页面空白/回退** | [#B018](#b018) |
| **template_data 双重序列化** | **range 迭代字符串报错，相册/自定义页面不显示** | [#B019](#b019) |
| **container 宽度缺失** | **导航栏 Logo 和菜单无法左右分布** | [#B020](#b020) |

---

## B001

**错误信息：** `function "Excerpt" not defined` 或 `function "Title" not defined`

**原因：** Go 模板中 `.` 和字段名之间有空格，被解析为函数调用

```html
<!-- ❌ 错误 -->
{{ . Excerpt }}
{{ . Category.Name }}

<!-- ✅ 修复：紧接写，无空格 -->
{{ .Excerpt }}
{{ .Category.Name }}
```

---

## B002

**错误信息：** `can't evaluate field Cover in type models.Post`（旧版本）

**更新（2026-03-27）：Post 模型已有 Cover 字段！**

```
Post 实际可用字段（已验证源码）：
ID, Title, Slug, Content, Excerpt, Cover, Views, Status
PublishedAt, CreatedAt, UpdatedAt
CategoryID, Category{ID,Name,Slug,Description,ParentID}
AuthorID, Author{ID,Name,Email}
Tags[]{ID,Name,Slug}
MetaTitle, MetaDescription

✅ Cover 字段存在（封面图 URL 字符串），但可能为空字符串，渲染前需 if 判断
❌ 仍然不存在：ReadTime, URL（URL 需手动拼接 /article/{{ .Slug }}）
```

**正确用法：**
```html
{{ if .post.Cover }}
<img src="{{ .post.Cover }}" alt="{{ .post.Title }}">
{{ end }}
```

---

## B003

**错误信息：** `can't evaluate field URL in type *models.Category`

**原因：** Category/Tag/Post 等模型没有 URL 字段，URL 需要手动拼接

```html
<!-- ❌ 错误 -->
<a href="{{ .Category.URL }}">{{ .Category.Name }}</a>
<a href="{{ .URL }}">{{ .Label }}</a>

<!-- ✅ 修复：手动拼接 -->
<a href="/category/{{ .Category.Slug }}">{{ .Category.Name }}</a>
<a href="/article/{{ .Slug }}">{{ .Title }}</a>
<a href="/tag/{{ .Slug }}">{{ .Name }}</a>
<a href="/page/{{ .Slug }}">{{ .Title }}</a>
```

---

## B004

**错误信息：** `no such template "_head"` 或 `no such template "_footer"`

**原因：** Go html/template 不支持跨文件的 template include

```html
<!-- ❌ 错误：尝试引用外部 partial 模板 -->
{{ template "_head" . }}
{{ template "partials/nav" . }}

<!-- ✅ 修复：每个页面包含完整 HTML 结构（复制粘贴 header/footer） -->
<!DOCTYPE html>
<html lang="zh-CN">
<head>...</head>
<body>
    <header>...</header>  <!-- 每个页面都要写完整 -->
    <main>...</main>
    <footer>...</footer>
</body>
</html>
```

---

## B005

**错误信息：** MIME type error，路径出现 `//themes/xxx//themes/xxx/`

**原因：** `asset_path` 函数自动添加 `/themes/{slug}/` 前缀，传参不要再带此前缀

```html
<!-- ❌ 错误：路径重复 -->
{{ asset_path "/themes/my-theme/assets/css/base.css" }}
<!-- 结果：/themes/my-theme//themes/my-theme/assets/css/base.css -->

<!-- ✅ 修复：只传相对路径 -->
{{ asset_path "assets/css/base.css" }}
<!-- 结果：/themes/my-theme/assets/css/base.css -->
```

---

## B006

**错误信息：** `wrong type for value; expected string; got []interface {}`

**原因：** `in` 函数参数顺序错误（先 item，后 slice）

```html
<!-- ❌ 错误：顺序反了 -->
{{ if in $widgets "profile" }}

<!-- ✅ 修复：in item slice -->
{{ if in "profile" $widgets }}
```

---

## B007

**错误信息：** `cannot index slice/array with type models.Post`

**原因：** 在 range 循环中，`.` 是当前元素（Post 对象），不是索引数字

```html
<!-- ❌ 错误：用 Post 对象当索引 -->
{{ range .posts }}
<span>{{ add (index $.posts .) 1 }}</span>
{{ end }}

<!-- ✅ 修复：用 $index, $item 语法 -->
{{ range $index, $post := .posts }}
<span>{{ add $index 1 }}</span>
<h3>{{ $post.Title }}</h3>
{{ end }}
```

---

## B008

**症状：** 用户在管理后台修改颜色，前台没有变化

**可能原因及修复：**

1. **CSS 变量名与配置 key 不对应**
   ```yaml
   # theme.yml
   - name: primaryDark    # camelCase
   ```
   ```css
   /* base.css */
   --primary-dark: #xxx;  /* kebab-case，对应关系 primaryDark ↔ --primary-dark */
   ```
   ```html
   <!-- 模板中注入用户配置 -->
   <style>:root{ --primary-dark: {{ getConfig "primaryDark" "#0f3460" }}; }</style>
   ```

2. **模板中没有注入 CSS 变量覆盖**
   ```html
   <!-- ✅ 每个模板的 <head> 中必须有 CSS 变量注入块 -->
   <link rel="stylesheet" href="{{ asset_path "assets/css/base.css" }}">
   <style>
   :root {
       --primary: {{ getConfig "primaryColor" "#2563eb" }};
       --primary-dark: {{ getConfig "primaryDark" "#0f3460" }};
   }
   </style>
   ```

3. **CSS 中硬编码颜色（未用变量）**
   ```css
   /* ❌ 错误 */
   .hero { background: #0f3460; }

   /* ✅ 修复 */
   .hero { background: var(--primary-dark); }
   ```

---

## B009

**错误信息：** MIME type error，FontAwesome CSS 被拒绝加载

**原因：** FontAwesome 路径错误

```html
<!-- ❌ 错误路径 -->
<link href="/themes/fontawesome/css/all.min.css">

<!-- ✅ 正确路径（全局路由，不在 themes 目录下） -->
<link rel="stylesheet" href="/fontawesome/css/all.min.css">
```

---

## B010

**症状：** 页面加载极慢，控制台显示 `fonts.googleapis.com` 超时

**原因：** Google Fonts 在国内无法访问

```html
<!-- ❌ 错误：使用 Google Fonts -->
<link href="https://fonts.googleapis.com/css2?family=Noto+Sans+SC&display=swap">

<!-- ✅ 修复：使用系统字体栈 -->
<!-- 直接删除 Google Fonts link，在 CSS 中设置系统字体 -->
```

```css
/* base.css */
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
                 'Helvetica Neue', Arial, 'Noto Sans SC', sans-serif;
}
```

---

## B011

**错误信息：** `wrong type for value` 或 range 循环报错

**原因：** 从数据库读取的配置值可能已被解析为 `[]interface{}`，直接用 `json_decode` 会出错

```html
<!-- ❌ 错误：直接 json_decode 可能失败 -->
{{ range (json_decode (getConfig "sidebarWidgets" `["profile"]`)) }}

<!-- ✅ 修复：用 ensureSlice 包装 -->
{{ $widgets := ensureSlice (getConfig "sidebarWidgets" `["profile","categories","tags"]`) }}
{{ range $widgets }}
    {{ . }}
{{ end }}
```

---

## B012

**症状：** 管理后台保存配置成功，但前台内容没有变化（无报错）

**排查流程：**
1. 确认模板中使用的 key 与 theme.yml 中 `configDefinition` 的 `name` 完全一致（区分大小写）
2. 检查 `config` 部分是否也定义了该 key（两处必须都有）
3. 检查 `getConfig` 函数的 key 参数是否拼写正确

```yaml
# ✅ 两处必须一致
config:
  postListStyle: "card"           # key 在这里

configDefinition:
  - section: layout
    fields:
      - name: postListStyle       # ← 必须与上面 key 完全一致
```

```html
<!-- ✅ 模板中 key 必须与上面一致 -->
{{ getConfig "postListStyle" "card" }}
```

---

## B013

**症状：** 管理后台主题配置界面显示"未命名分组"

**原因：** `configDefinition` 中缺少 `label` 字段（只有 `section`）

```yaml
# ❌ 错误：没有 label
configDefinition:
  - section: colors
    fields: ...

# ✅ 修复：必须同时有 section 和 label
configDefinition:
  - section: colors
    label: 颜色设置     # ← 必须有 label
    fields: ...
```

---

## B014

**症状：** 文章内容页显示 Markdown 源码（`# 标题` `**粗体**` 等）

**原因：** 文章内容未经过 Markdown 渲染

```html
<!-- ❌ 错误：直接输出原始内容 -->
<div class="article-content">{{ .post.Content }}</div>

<!-- ✅ 修复：使用 markdownHTML（不需要 | safe，见 B017） -->
<article class="article-content">
    {{ markdownHTML .post.Content }}
</article>
```

---

## B015

**症状：** CSS 动画/轮播功能不正常（显示空白或只显示第一帧）

**原因：** CSS `animation` 属性和对应的 `@keyframes` 必须同时存在

```css
/* ❌ 错误：有 animation 但无 @keyframes */
.slide { animation: fadeInOut 8s infinite; }
/* 缺少 @keyframes fadeInOut */

/* ✅ 修复：两者必须配套 */
.carousel-slide:nth-child(1) { animation: slideShow 12s infinite 0s; }
.carousel-slide:nth-child(2) { animation: slideShow 12s infinite 4s; }
.carousel-slide:nth-child(3) { animation: slideShow 12s infinite 8s; }

@keyframes slideShow {
    0%, 100% { opacity: 0; }
    8%, 25% { opacity: 1; }
    33%, 99% { opacity: 0; }
}
```

**轮播配置公式（N 张图，总时长 T 秒）：**
- 每张显示时间 ≈ T/N 秒
- 每张延迟：index × (T/N) 秒
- @keyframes 显示区间：`0% 到 (100/N)%`

---

## B016

**症状：** 分页控件不显示，或出现 `nil pointer` 错误

**原因：** 分页字段名是 `pagination`，不是 `paginator`。且字段结构为 `{page, pageSize, total, totalPages}`，没有 `Pages[]` 数组。

```html
<!-- ❌ 错误：字段名错误 -->
{{ .paginator.TotalPages }}
{{ range .paginator.Pages }}

<!-- ✅ 正确：字段名是 pagination，字段是小写 -->
{{ with .pagination }}
    当前第 {{ .page }} 页 / 共 {{ .totalPages }} 页
{{ end }}
```

**⚠️ 注意：**
- `pagination` 只在 `article_list.html` 中存在，其他页面（category、tag、home 等）**没有分页**
- `category.html` 和 `tag.html` 目前返回全量文章，没有分页参数

---

## B017 ⚠️ 高危（2026-03-27 新增）

**症状：** 页面只显示标题和导航，无正文/内容区域。服务器日志出现模板执行错误，**静默回退到默认模板**。表现为使用 `page-gallery.html` 时实际渲染了 `page.html`，页面内容为空。

**错误信息：** `expected string; got template.HTML` 或服务器日志 `[ERROR] ExecuteTemplate failed`

**原因：** `markdownHTML` 函数返回的是 `template.HTML` 类型（已标记安全），不需要再接 `| safe`。`| safe` 期望接收 `string` 类型参数，类型冲突导致模板执行失败。

```html
<!-- ❌ 错误：markdownHTML 已返回 template.HTML，| safe 会类型冲突 -->
<div class="article-body">{{ markdownHTML .post.Content | safe }}</div>
<div class="article-body">{{ markdownHTML .page.Content | safe }}</div>

<!-- ✅ 修复：直接使用 markdownHTML，不加 | safe -->
<div class="article-body">{{ markdownHTML .post.Content }}</div>
<div class="article-body">{{ markdownHTML .page.Content }}</div>
```

**⚠️ 此 Bug 的危险性在于：** 模板执行失败时 Go 模板引擎不会返回 HTTP 500，而是由 `pageShowHandler` 的回退逻辑渲染默认 `page.html`。页面看起来"能打开"但内容不对，非常难排查。**排查方法：检查服务器 stderr 日志中的 `[ERROR]` 行。**

---

## B018（2026-03-27 新增）

**症状：** 自定义页面（page-gallery.html 等）渲染失败，静默回退到默认 page.html

**错误信息：** `can't evaluate field Description in type models.Page`

**原因：** Page 模型没有 `Description` 字段，只有 `Excerpt` 字段

```
Page 实际可用字段（已验证源码）：
ID, Title, Slug, Content, Excerpt, Status, Template,
TemplateData, MetaTitle, MetaDescription, Views, SortOrder,
CreatedAt, UpdatedAt

⚠️ Page 没有：Description, Category, Tags, Author
```

```html
<!-- ❌ 错误：Page 没有 Description 字段 -->
{{ if .page.Description }}
<p>{{ .page.Description }}</p>
{{ end }}

<!-- ✅ 修复：使用 Excerpt -->
{{ if .page.Excerpt }}
<p>{{ .page.Excerpt }}</p>
{{ end }}
```

---

## B019 ⚠️ 高危（2026-03-27 新增）

**症状：** 相册/产品展示等自定义页面只显示标题，无图片、无网格。服务器日志显示 `range can't iterate over [...]`

**错误信息：** `range can't iterate over [{"url":"/uploads/..."}]`

**原因：** `template_data` 中数组字段被存储为 JSON **字符串**（双重序列化），而非真正的 JSON 数组。`range` 无法迭代字符串。

```
数据库中的错误格式（images 是字符串）：
{"images":"[{\"url\":\"...\"}]", "columns":"4"}
                        ↑ 被多包了一层 JSON 字符串

正确格式（images 是数组）：
{"images":[{"url":"..."}], "columns":"4"}
```

**根本原因：** 前端 `PageEdit.vue` 保存 `image-list` 类型的 `template_data` 时，将图片数组序列化为 JSON 字符串后直接存入外层 JSON，导致双重序列化。

**模板层防护（推荐）：** 使用 `ensureSlice` + `json_decode` 处理不确定类型：

```html
<!-- 方案 A：修复数据库后直接 range（根本修复） -->
{{ range .td.images }}
<div class="gallery-item">
    <img src="{{ .url }}" alt="{{ .title }}">
</div>
{{ end }}

<!-- 方案 B：模板中添加防护（兼容可能存在的字符串值） -->
{{ $images := ensureSlice .td.images }}
{{ if eq (printf "%T" (index $images 0)) "string" }}
    {{ range $img := $images }}
        {{ $item := json_decode $img }}
        <img src="{{ $item.url }}">
    {{ end }}
{{ else }}
    {{ range $images }}
    <img src="{{ .url }}">
    {{ end }}
{{ end }}
```

**数据库修复脚本（Python）：**
```python
import sqlite3, json
conn = sqlite3.connect('catter.db')
cur = conn.cursor()
cur.execute('SELECT id, slug, template_data FROM pages WHERE template_data IS NOT NULL')
for id, slug, td_str in cur.fetchall():
    data = json.loads(td_str)
    changed = False
    for key, val in data.items():
        if isinstance(val, str):
            try:
                parsed = json.loads(val)
                if isinstance(parsed, (list, dict)):
                    data[key] = parsed
                    changed = True
            except: pass
    if changed:
        cur.execute('UPDATE pages SET template_data = ? WHERE id = ?',
                    (json.dumps(data, ensure_ascii=False), id))
conn.commit()
conn.close()
```

---

## B020（2026-03-27 新增）

**症状：** 导航栏中 Logo 和菜单没有左右分开，全部挤在中间或左对齐

**原因：** `.container` 没有设置 `width: 100%`，导致 flex 子元素无法撑满容器宽度，`justify-content: space-between` 无法生效

```css
/* ❌ 错误：container 没有设置 width */
.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}

/* ✅ 修复：添加 width: 100% */
.container {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}
```

同时确保 `.nav-inner` 也设置了 `width: 100%`：
```css
.nav-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;   /* ← 关键 */
    height: 100%;
}
```

**排查流程：** 导航栏布局异常时，依次检查：1) `.container` 是否有 `width: 100%`  2) `.nav-inner` 是否有 `width: 100%`  3) 是否使用了 `justify-content: space-between`

---

## B021（2026-03-27 新增）

**症状：** 图标在不同系统上显示不一致（Windows 上正常、macOS 上变色、Linux 上显示方块），或图标颜色无法通过 CSS 变量控制

**原因：** 使用 Emoji 字符（🚛📞📍✉️💰📅📁👁 等）替代 FontAwesome 图标

```html
<!-- ❌ 错误：使用 Emoji 作为 UI 图标 -->
<div class="nav-logo-icon">🚛</div>
<span class="phone-icon">📞</span>
<div class="feature-icon">💰</div>
<div class="post-meta">📅 2026-01-01</div>
<span>👁 1234 阅读</span>

<!-- ✅ 修复：使用 FontAwesome 图标 -->
<div class="nav-logo-icon"><i class="fa-solid fa-truck"></i></div>
<span class="phone-icon"><i class="fa-solid fa-phone"></i></span>
<div class="feature-icon"><i class="fa-solid fa-coins"></i></div>
<div class="post-meta"><i class="fa-solid fa-calendar"></i> 2026-01-01</div>
<span><i class="fa-solid fa-eye"></i> 1234 阅读</span>
```

**常见 Emoji → FontAwesome 速查表：**

| Emoji | FontAwesome | 用途 |
|-------|------------|------|
| 🚛 | `fa-solid fa-truck` | 卡车/运输 |
| 📞 | `fa-solid fa-phone` | 电话 |
| 📍 | `fa-solid fa-location-dot` | 定位/地址 |
| ✉️ | `fa-solid fa-envelope` | 邮箱 |
| 📰 | `fa-solid fa-newspaper` | 新闻 |
| 🛡️ | `fa-solid fa-shield-halved` | 安全/防护 |
| 📡 | `fa-solid fa-satellite-dish` | 卫星/GPS |
| 💰 | `fa-solid fa-coins` | 资金/价格 |
| 🔰 | `fa-solid fa-certificate` | 证书/资质 |
| 🏭 | `fa-solid fa-industry` | 工厂/工业 |
| 📦 | `fa-solid fa-box` | 包裹/运输 |
| 📋 | `fa-solid fa-list` | 列表 |
| ⏰ | `fa-solid fa-clock` | 时钟/时间 |
| ⚡ | `fa-solid fa-bolt` | 快捷/效率 |
| 🌐 | `fa-solid fa-globe` | 全球/国际化 |
| ✅ | `fa-solid fa-check` | 确认/通过 |
| ☁️ | `fa-solid fa-cloud` | 云/在线 |
| 📠 | `fa-solid fa-fax` | 传真 |
| 💬 | `fa-solid fa-comments` | 评论/咨询 |
| 📷 | `fa-solid fa-camera` | 相机/图片 |
| 🔍 | `fa-solid fa-magnifying-glass` | 搜索 |
| 📁 | `fa-solid fa-folder` | 文件夹/分类 |
| 📅 | `fa-solid fa-calendar` | 日期/日历 |
| 👁 | `fa-solid fa-eye` | 浏览/阅读量 |
| ❮ | `fa-solid fa-chevron-left` | 左箭头 |
| ❯ | `fa-solid fa-chevron-right` | 右箭头 |

**批量检测命令（Python）：**
```python
import os, sys
sys.stdout.reconfigure(encoding='utf-8')
for root, dirs, files in os.walk('templates/themes/THEME_NAME'):
    for fname in files:
        if not fname.endswith('.html'): continue
        with open(os.path.join(root, fname), 'r', encoding='utf-8-sig') as f:
            for i, line in enumerate(f, 1):
                for c in line:
                    cp = ord(c)
                    if (0x2600<=cp<=0x27BF) or (0x1F300<=cp<=0x1FAFF) or cp in (0xFE0F, 0x200D):
                        print(f'{fname}:{i} emoji detected')
                        break
```

**排查流程：** 1) 用检测命令扫描所有 HTML 文件  2) 根据 Emoji → FontAwesome 速查表替换  3) 确保 `<head>` 中引入了 `/fontawesome/css/all.min.css`
