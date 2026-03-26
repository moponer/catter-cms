# Catter CMS 问题修复记录

> 记录所有已解决的问题、解决方案和经验教训，作为问题排查手册。

**创建日期：** 2026-03-25
**版本：** 1.1

---

## 目录

- [Go 模板语法问题](#go-模板语法问题)
- [主题配置问题](#主题配置问题)
- [前端编译问题](#前端编译问题)
- [主题样式问题](#主题样式问题)
- [模板函数问题](#模板函数问题)
- [静态资源加载问题](#静态资源加载问题)
- [系统设置问题](#系统设置问题)
- [经验教训总结](#经验教训总结)
- [附录：Post 模型字段参考](#附录post-模型字段参考)
- [附录：系统设置数据流架构](#附录系统设置数据流架构)

---

## Go 模板语法问题

### 问题 1: `function "Excerpt" not defined`

**错误信息：**
```
template execute error: template: home.html:81: function "Excerpt" not defined
```

**原因：**
Go 模板中 `. Excerpt`（dot 和字段名之间有空格）会被解析为**函数调用**，而不是字段访问。

**错误写法：**
```html
<p class="post-excerpt">{{ . Excerpt }}</p>
```

**正确写法：**
```html
<p class="post-excerpt">{{ .Excerpt }}</p>
```

**影响文件：**
- `templates/themes/dianyan-blog/pages/home.html`

**修复日期：** 2026-03-25

---

### 问题 2: `can't evaluate field Cover in type models.Post`

**错误信息：**
```
executing "home.html" at <.Cover>: can't evaluate field Cover in type models.Post
```

**原因：**
模板引用了 `Post` 模型中不存在的 `Cover` 字段。Post 模型没有封面图字段。

**Post 模型可用字段：**
```go
type Post struct {
    ID, Title, Slug, Content, Excerpt, Views, Status
    PublishedAt, CreatedAt, UpdatedAt
    CategoryID, Category, Tags
    AuthorID, Author
    // 注意：没有 Cover、ReadTime 等字段
}
```

**错误写法：**
```html
{{ if .Cover }}
<div class="post-thumbnail">
    <img src="{{ .Cover }}" alt="{{ .Title }}">
</div>
{{ end }}
```

**正确做法：**
移除缩略图判断块，或使用实际存在的字段。

**影响文件：**
- `templates/themes/dianyan-blog/pages/home.html`

**修复日期：** 2026-03-25

---

### 问题 3: `can't evaluate field URL in type *models.Category`

**错误信息：**
```
executing "home.html" at <.Category.URL>: can't evaluate field URL in type *models.Category
```

**原因：**
Category/Tag/Post 模型**没有 URL 字段**，URL 需要在模板中手动拼接。

**错误写法：**
```html
<a href="{{ .Category.URL }}">{{ .Category.Name }}</a>
```

**正确写法：**
```html
<a href="/category/{{ .Category.Slug }}">{{ .Category.Name }}</a>
```

**URL 拼接规则：**
| 类型 | 模板写法 |
|------|---------|
| 文章 | `/article/{{ .Slug }}` |
| 分类 | `/category/{{ .Slug }}` |
| 标签 | `/tag/{{ .Slug }}` |
| 页面 | `/page/{{ .Slug }}` |

**影响文件：**
- `templates/themes/dianyan-blog/pages/home.html`
- `templates/themes/dianyan-blog/pages/tag.html`
- `templates/themes/dianyan-blog/pages/category.html`
- `templates/themes/dianyan-blog/pages/article_show.html`
- `templates/themes/dianyan-blog/pages/article_list.html`

**修复日期：** 2026-03-25

---

## 主题配置问题

### 问题 4: 主题配置不生效

**症状：**
管理后台修改主题配置后，前台页面没有变化。

**可能原因：**

1. **配置 key 名不匹配**
   - 模板中使用的 key 与 `theme.yml` 中 `configDefinition` 的 `name` 不一致
   - 例如：模板用 `listStyle`，yml 配置用 `postListStyle`

2. **配置未正确保存到数据库**
   - 检查 `SaveThemeConfigToDB` 函数是否正确执行

3. **配置未正确加载到内存**
   - 检查 `LoadTheme` 函数是否正确读取数据库配置

4. **模板缓存未刷新**
   - 重启服务器或重新激活主题

5. **后端未重启**
   - 添加新函数（如 `ensureSlice`）后必须重启服务器

**排查步骤：**

1. 检查服务器日志中的配置加载信息：
   ```
   [DEBUG] Theme config keys: [primaryColor sidebarWidgets ...]
   ```

2. 在管理后台重新保存配置，观察日志输出

3. 重启服务器让配置重新加载

4. 确认 `getConfig` 的 key 与 `theme.yml` 中 `configDefinition` 的 `name` 完全一致

**配置保存流程（代码追踪）：**

```
前端 saveConfig() 
  → saveThemeConfig(configValues, slug)
    → PUT /api/themes/:slug/config
      → handlers.SaveThemeConfig()
        → services.UpdateThemeConfig(slug, config)
          → models.SaveThemeConfig() // 保存到数据库
          → models.GetThemeConfig()  // 重新加载
          → theme.Config = reloadedConfig  // 更新内存
```

**配置加载流程（代码追踪）：**

```
theme.LoadTheme()
  → yaml.Unmarshal()  // 解析 theme.yml 的 config
  → FlattenMap()  // 扁平化嵌套 config
  → models.GetThemeConfig()  // 加载数据库配置
  → DeepMerge(theme.Config, savedConfig)  // 合并
  → theme.Config  // 最终配置
```

**修复日期：** 2026-03-25

---

### 问题 4.1: dianyan-blog 配置 key 名不匹配

**具体问题：**
- 模板使用 `getConfig "listStyle" "card"`
- theme.yml 使用 `postListStyle: card`

**修复：**
统一使用 `postListStyle` 作为 key 名。

**影响文件：**
- `templates/themes/dianyan-blog/pages/home.html`

**修复日期：** 2026-03-25

---

### 问题 5: `in` 函数参数顺序错误

**错误信息：**
```
wrong type for value; expected string; got []interface {}
```

**原因：**
Go 模板中 `in` 函数的参数顺序是 `in item slice`，而不是 `in slice item`。

**错误写法：**
```html
{{ if in $widgets "profile" }}
```

**正确写法：**
```html
{{ if in "profile" $widgets }}
```

**函数签名：**
```go
"in": func(item interface{}, slice interface{}) bool
```

**影响文件：**
- `templates/themes/dianyan-blog/pages/home.html`

**修复日期：** 2026-03-25

---

### 问题 6: `ensureSlice` 函数 - 处理配置值类型转换

**背景：**
`getConfig` 函数返回的配置值可能是：
- 已解析的 `[]interface{}`（从数据库读取后自动解析）
- JSON 字符串（默认值）
- 其他类型

直接传给 `json_decode` 会报错。

**解决方案：**
添加 `ensureSlice` 函数统一处理：

```go
"ensureSlice": func(v interface{}) []interface{} {
    if v == nil {
        return []interface{}{}
    }
    rv := reflect.ValueOf(v)
    switch rv.Kind() {
    case reflect.Slice, reflect.Array:
        // 已经是切片，直接返回
        result := make([]interface{}, rv.Len())
        for i := 0; i < rv.Len(); i++ {
            result[i] = rv.Index(i).Interface()
        }
        return result
    case reflect.String:
        // 如果是字符串，尝试解析为JSON
        var decoded interface{}
        if err := json.Unmarshal([]byte(rv.String()), &decoded); err == nil {
            if s, ok := decoded.([]interface{}); ok {
                return s
            }
        }
        return []interface{}{}
    }
    return []interface{}{}
}
```

**使用方式：**
```html
{{ $widgets := ensureSlice (getConfig "sidebarWidgets" `["profile","categories","tags"]`) }}
```

**影响文件：**
- `gin-app/internal/services/theme.go`

**修复日期：** 2026-03-25

---

## 前端编译问题

### 问题 7: `onUnmounted` 重复导入

**错误信息：**
```
Identifier 'onUnmounted' has already been declared.
```

**原因：**
`ThemeSettings.vue` 中 `onUnmounted` 被导入了两次。

**错误代码：**
```javascript
// 第137行
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

// 第177行（重复）
import { onUnmounted } from 'vue'
onUnmounted(() => {
```

**修复方式：**
删除第177行的重复 import 语句。

**影响文件：**
- `admin/src/views/themes/ThemeSettings.vue`

**修复日期：** 2026-03-25

---

## 主题样式问题

### 问题 8: Footer 三栏布局不生效

**症状：**
Footer 三个区块垂直排列，而不是横向三栏布局。

**原因：**
`footer-content` 和 `footer-section` CSS 类没有定义样式。

**解决方案：**
添加 CSS 样式：

```css
.footer-content {
    display: flex;
    justify-content: space-between;
    gap: var(--space-xl);
    margin-bottom: var(--space-xl);
}

.footer-section {
    flex: 1;
    text-align: left;
}
```

**影响文件：**
- `templates/themes/dianyan-blog/assets/css/base.css`

**修复日期：** 2026-03-25

---

### 问题 9: 侧边栏模块未实现配置逻辑

**症状：**
- 管理后台启用了多个侧边栏模块，但前台只显示标签云
- 标签云样式丑，间距过大

**原因：**
- 模板中硬编码了侧边栏内容，没有使用 `sidebarWidgets` 配置
- CSS 样式不完善

**解决方案：**

1. 实现配置驱动的侧边栏：
```html
{{ $widgets := ensureSlice (getConfig "sidebarWidgets" `["profile","categories","tags"]`) }}

{{ if in "profile" $widgets }}
<div class="sidebar-widget">...站点资料...</div>
{{ end }}

{{ if in "categories" $widgets }}
<div class="sidebar-widget">...分类...</div>
{{ end }}

{{ if in "tags" $widgets }}
<div class="sidebar-widget">...标签云...</div>
{{ end }}
```

2. 完善 CSS 样式：
```css
.tag-item {
    display: inline-flex;
    align-items: center;
    padding: var(--space-xs) var(--space-md);
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    font-size: var(--text-xs);
    border-radius: var(--radius-sm);
    border: 1px solid var(--border-color);
    transition: all 0.2s ease;
}

.tag-item:hover {
    background: var(--primary);
    color: white;
    border-color: var(--primary);
}
```

**影响文件：**
- `templates/themes/dianyan-blog/pages/home.html`
- `templates/themes/dianyan-blog/assets/css/base.css`

**修复日期：** 2026-03-25

---

### 问题 10: 内联 CSS 变量（违反规范）

**问题：**
在 HTML 中使用 `<style>` 块定义 CSS 变量，违反主题开发规范。

**错误做法：**
```html
<head>
    <style>
        :root {
            --primary: {{ getConfig "primaryColor" "#2d5a4a" }};
        }
    </style>
</head>
```

**正确做法：**
CSS 变量默认值定义在 `base.css` 的 `:root` 中：
```css
:root {
    --primary: #2d5a4a;
    --primary-light: #4a7c6f;
    /* ... */
}
```

**影响文件：**
- `templates/themes/dianyan-blog/pages/home.html`

**修复日期：** 2026-03-25

---

## 模板函数问题

### 问题 11: Go 模板不支持跨文件包含

**错误信息：**
```
no such template "_head"
```

**原因：**
Go 标准库 `html/template` 不支持 `{{ template "_head" . }}` 跨文件模板包含。

**解决方案：**
每个页面模板必须包含完整的 HTML 结构：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>{{ .site.title }}</title>
    <link rel="stylesheet" href="{{ asset_path "/themes/my-theme/assets/css/base.css" }}">
</head>
<body>
    <!-- 完整内容 -->
</body>
</html>
```

**禁止做法：**
```html
<!-- ❌ 错误：会报 "no such template" 错误 -->
{{ template "_head" . }}

<!-- ❌ 错误：每个页面应包含完整的 head -->
<head>
    {{ template "_head" . }}  <!-- 不要拆分 partial -->
</head>
```

**影响文件：**
- 所有主题模板

**记录日期：** 2026-03-25

---

## 静态资源加载问题

### 问题 12: FontAwesome MIME 类型错误（资源被阻止）

**症状：**
```
来自 'http://localhost:3002/themes/fontawesome/css/all.min.css' 的资源已被阻止，
因为 MIME 类型不匹配（X-Content-Type-Options: nosniff）
```

**原因：**
模板中使用的静态资源路径与后端配置的静态文件服务路径不匹配。

**错误路径：**
```html
<link rel="stylesheet" href="/themes/fontawesome/css/all.min.css">
```

**正确路径：**
```html
<link rel="stylesheet" href="/fontawesome/css/all.min.css">
```

**路径映射说明：**
| 资源类型 | 后端配置 | 模板使用 |
|---------|---------|---------|
| FontAwesome | `r.Static("/fontawesome", "./static/fontawesome")` | `/fontawesome/...` |
| 主题资源 | `r.Static("/themes", "./templates/themes")` | `/themes/{theme-name}/assets/...` |

**批量修复命令（VS Code 正则）：**
- 查找：`/themes/fontawesome/`
- 替换：`/fontawesome/`

**影响文件：**
- `templates/themes/dianyan-blog/pages/*.html`（10个文件）
- `templates/themes/catter-cms-official/pages/*.html`（10个文件）

**修复日期：** 2026-03-25

---

### 问题 13: Google Fonts 加载慢（国内访问问题）

**症状：**
页面加载时 `fonts.googleapis.com` 和 `fonts.gstatic.com` 请求超时，严重影响页面加载速度。

**尝试的解决方案（失败）：**
1. 使用国内镜像 `fonts.googleapis.cn.com` → 仍无法访问
2. 使用其他镜像 `fonts.font.im` → 同样失败

**最终解决方案：**
完全移除 Google Fonts，使用系统字体栈：

```css
body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
                 "Helvetica Neue", Arial, "Noto Sans SC", sans-serif,
                 "Apple Color Emoji", "Segoe UI Emoji";
}
```

**修复步骤：**
1. 移除模板中的 Google Fonts 链接：
   ```html
   <!-- ❌ 删除 -->
   <link href="https://fonts.googleapis.com/css2?family=..." rel="stylesheet">
   ```

2. 更新 CSS font-family 使用系统字体栈

3. 对于需要特殊字体的场景（如新中式风格），使用本地字体：
   ```css
   @font-face {
       font-family: 'Noto Serif SC';
       src: url('/fonts/NotoSerifSC.woff2') format('woff2');
   }
   ```

**影响文件：**
- `templates/themes/dianyan-blog/pages/*.html`
- `templates/themes/dianyan-blog/assets/css/base.css`

**修复日期：** 2026-03-25

---

## 系统设置问题

### 问题 14: 系统设置保存无效

**症状：**
- 点击保存后提示成功，但刷新页面设置项全部恢复为空
- 数据库中 settings 表没有新记录

**原因：**
前端发送的数据结构格式不正确。

**错误写法（嵌套结构）：**
```javascript
const handleSave = async () => {
    await put('/api/settings', {
        category: 'basic',  // ❌ 后端不期望这个结构
        settings: { ...settings.basic }
    })
}
```

**正确写法（扁平结构）：**
```javascript
const handleSave = async () => {
    const allSettings = {
        ...settings.basic,
        ...settings.seo,
        ...settings.email
    }
    await put('/api/settings', allSettings)  // ✅ 扁平对象
}
```

**后端处理逻辑：**
```go
// PUT /api/settings
func UpdateSettings(c *gin.Context) {
    var data map[string]interface{}
    c.BindJSON(&data)
    
    for key, value := range data {
        models.SaveSetting(key, value)  // 逐个保存
    }
}
```

**影响文件：**
- `admin/src/views/settings/SiteSettings.vue`

**修复日期：** 2026-03-25

---

### 问题 15: 系统设置前台不显示

**症状：**
- 管理后台保存设置成功，再次进入也能看到保存的值
- 但前台页面没有任何变化，站点名称仍显示默认值

**原因：**
前端字段名与后端数据库 key 不匹配。

**字段名映射问题：**
| 前端字段名（错误） | 后端 key（正确） | 模板访问 |
|------------------|----------------|---------|
| `site_name` | `site_title` | `{{ .site.title }}` |
| `site_desc` | `site_description` | `{{ .site.description }}` |

**修复方案（选择方案 B）：**
统一使用后端定义的 key 名称：

```javascript
// SiteSettings.vue
const settings = reactive({
    basic: {
        site_title: '',      // ✅ 与后端一致
        site_description: '', // ✅ 与后端一致
        site_url: '',
        logo: '',
        favicon: ''
    }
})
```

**数据流追踪：**
```
前端保存（site_title）
  → API: PUT /api/settings { site_title: "..." }
    → 后端: settings 表 { key: "site_title", value: "..." }
      → 渲染: getSiteData() 返回 { title: "..." }
        → 模板: {{ .site.title }}
```

**关键代码验证：**
```go
// services/site.go
func getSiteData() map[string]interface{} {
    return map[string]interface{}{
        "title":       models.GetSettingAsString("site_title"),
        "description": models.GetSettingAsString("site_description"),
        // ...
    }
}
```

**影响文件：**
- `admin/src/views/settings/SiteSettings.vue`

**修复日期：** 2026-03-25

---

### 问题 16: 数据库设置值损坏

**症状：**
部分设置项值异常，包含多余的引号和反斜杠。

**损坏示例：**
```sql
-- site_title 的值
value: '" catter-blog\"'  -- ❌ 损坏

-- favicon 的值  
value: '"\"'  -- ❌ 损坏
```

**原因：**
可能是前端发送数据时 JSON 编码错误，或后端保存时重复编码。

**修复方案：**
直接修复数据库：

```sql
-- 删除损坏的记录
DELETE FROM settings WHERE key IN ('site_title', 'site_description', 'favicon', 'logo', 'posts_per_page');

-- 重新保存正确的值（通过管理后台）
```

**预防措施：**
后端保存设置时进行数据清洗：
```go
func SaveSetting(key string, value interface{}) error {
    // 处理字符串值，去除多余的引号
    if strVal, ok := value.(string); ok {
        strVal = strings.Trim(strVal, "\"")
        value = strVal
    }
    // ... 保存到数据库
}
```

**修复日期：** 2026-03-25

---

### 问题 17: 站点名称显示为 URL 或乱码

**症状：**
- 前台显示站点名称为 `http://localhost:3002\` 或其他异常值
- 激活主题失败，日志显示 `Theme " 老灯Blog\ not found`

**原因：**
数据库中多个设置值被污染：
1. `site_title` 被错误设置为 `" http://localhost:3002\"`（带引号和反斜杠）
2. `site_name` 字段不应存在（正确的字段是 `site_title`）
3. `active_theme` 被错误设置为 `"老灯Blog\"`（站点名称被误存为主题名）

**数据污染示例：**
```sql
-- 错误数据
site_title = " http://localhost:3002\"
site_name = "老灯Blog"      -- 多余字段
active_theme = "老灯Blog\"  -- 应该是主题名如 "dianyan-blog"

-- 正确数据
site_title = "老灯Blog"
active_theme = "dianyan-blog"
```

**修复方案：**
```sql
-- 删除多余字段
DELETE FROM settings WHERE key='site_name';

-- 修正站点名称（JSON 格式，带引号）
UPDATE settings SET value='"老灯Blog"' WHERE key='site_title';

-- 修正主题名称
UPDATE settings SET value='"dianyan-blog"' WHERE key='active_theme';
```

**⚠️ 注意：** settings 表中的字符串值必须是 JSON 格式（带引号），因为后端使用 `json.Marshal` 保存，使用 `json.Unmarshal` 读取。

**影响文件：**
- `catter.db` (settings 表)

**修复日期：** 2026-03-25

---

## 经验教训总结

### Go 模板语法要点

1. **dot 和字段名之间不能有空格**
   - ✅ `{{ .Excerpt }}`
   - ❌ `{{ . Excerpt }}`

2. **字段访问 vs 函数调用**
   - `{{ .Title }}` - 访问字段
   - `{{ Title . }}` - 调用函数，传入 dot

3. **函数参数顺序要正确**
   - `in` 函数：`in item slice`（先 item，后 slice）

### 主题开发规范

1. **CSS 样式必须封装在 CSS 文件中**
   - 禁止内联 `style="..."` 属性
   - 禁止在 HTML 中嵌入 `<style>` 块（CSS 变量除外）

2. **每个页面模板包含完整 HTML 结构**
   - 不要尝试拆分 partial 模板

3. **CSS 变量默认值在 base.css 的 :root 中定义**

4. **配置 key 名必须与模板中 getConfig 的 key 一致**

### 静态资源路径规范

1. **FontAwesome 路径**
   - ✅ 正确：`/fontawesome/css/all.min.css`
   - ❌ 错误：`/themes/fontawesome/css/all.min.css`

2. **主题资源路径**
   - ✅ 正确：`/themes/{theme-name}/assets/css/...`
   - ❌ 错误：`/themes/assets/css/...`

3. **Google Fonts 替代方案**
   - 国内环境：使用系统字体栈
   - 必须使用特殊字体：下载本地 woff2 文件

### 系统设置开发规范

1. **前端字段名必须与后端 key 一致**
   ```javascript
   // ✅ 正确
   { site_title: '', site_description: '' }
   
   // ❌ 错误
   { site_name: '', site_desc: '' }
   ```

2. **发送扁平结构数据**
   ```javascript
   // ✅ 正确：扁平对象
   { site_title: 'My Site', site_description: '...' }
   
   // ❌ 错误：嵌套结构
   { category: 'basic', settings: { ... } }
   ```

3. **数据流验证**
   - 前端字段名 → 后端 key → 数据库 → getSiteData() → 模板变量
   - 每个环节都要检查字段名是否一致

4. **⚠️ settings 表存储格式（重要）**
   
   后端使用 `json.Marshal` 保存值，所以数据库中存储的是 JSON 格式：
   
   ```sql
   -- ✅ 正确格式（字符串值带引号）
   INSERT INTO settings (key, value) VALUES ('site_title', '"老灯Blog"');
   INSERT INTO settings (key, value) VALUES ('active_theme', '"dianyan-blog"');
   INSERT INTO settings (key, value) VALUES ('posts_per_page', '10');
   
   -- ❌ 错误格式（缺少引号）
   INSERT INTO settings (key, value) VALUES ('site_title', '老灯Blog');
   ```
   
   **为什么需要引号？**
   ```go
   // 后端保存逻辑
   func SetSetting(key string, value interface{}) error {
       valueStr, _ := json.Marshal(value)  // "老灯Blog" → `"老灯Blog"`
       // 存储到数据库
   }
   
   // 后端读取逻辑
   func GetSettingAsString(key string) string {
       value := "..." // 从数据库读取
       var result string
       json.Unmarshal([]byte(value), &result)  // `"老灯Blog"` → "老灯Blog"
       return result
   }
   ```
   
   **手动修复数据库时注意**：
   ```sql
   -- ✅ 正确
   UPDATE settings SET value='"老灯Blog"' WHERE key='site_title';
   
   -- ❌ 错误（会导致前台显示带引号）
   UPDATE settings SET value='老灯Blog' WHERE key='site_title';
   ```

5. **禁止使用的字段名**
   - ❌ `site_name` → 应使用 `site_title`
   - ❌ `site_desc` → 应使用 `site_description`

---

## 附录：Post 模型字段参考

```go
type Post struct {
    ID              uint       `json:"id"`
    Title           string     `json:"title"`
    Slug            string     `json:"slug"`
    Content         string     `json:"content"`
    Excerpt         string     `json:"excerpt"`
    Status          string     `json:"status"`       // "published" | "draft"
    Views           int        `json:"views"`
    MetaTitle       string     `json:"meta_title"`
    MetaDescription string     `json:"meta_description"`
    CategoryID      *uint      `json:"category_id"`
    Category        *Category  `json:"category"`     // 关联分类
    AuthorID        uint       `json:"author_id"`
    Author          *User      `json:"author"`       // 关联作者
    Tags            []Tag      `json:"tags"`         // 关联标签
    PublishedAt     *time.Time `json:"published_at"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
}
```

---

## 附录：系统设置数据流架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         前端 (Vue 3)                             │
├─────────────────────────────────────────────────────────────────┤
│  SiteSettings.vue                                               │
│  ├─ settings.basic.site_title                                   │
│  ├─ settings.basic.site_description                             │
│  └─ handleSave() → PUT /api/settings                            │
│      发送扁平对象: { site_title, site_description, ... }         │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                         后端 (Gin)                               │
├─────────────────────────────────────────────────────────────────┤
│  handlers.UpdateSettings(c *gin.Context)                        │
│  ├─ 解析 JSON: map[string]interface{}                           │
│  └─ 循环保存: models.SaveSetting(key, value)                    │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                       数据库 (SQLite)                            │
├─────────────────────────────────────────────────────────────────┤
│  settings 表                                                    │
│  ┌─────────────────┬────────────────────────┐                   │
│  │ key             │ value                  │                   │
│  ├─────────────────┼────────────────────────┤                   │
│  │ site_title      │ "Catter Blog"          │                   │
│  │ site_description│ "A modern blog"        │                   │
│  │ site_url        │ "https://example.com"  │                   │
│  └─────────────────┴────────────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      模板渲染 (SSR)                              │
├─────────────────────────────────────────────────────────────────┤
│  services.getSiteData()                                         │
│  └─ 读取: models.GetSettingAsString("site_title")               │
│      返回: map[string]interface{}{                              │
│          "title": "Catter Blog",                                │
│          "description": "A modern blog",                        │
│          ...                                                    │
│      }                                                          │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      前台模板 (Go Template)                      │
├─────────────────────────────────────────────────────────────────┤
│  <title>{{ .site.title }}</title>                               │
│  <meta name="description" content="{{ .site.description }}">    │
└─────────────────────────────────────────────────────────────────┘
```

**关键字段映射表：**

| 前端字段名 | 后端 Key | 模板变量 |
|-----------|---------|---------|
| `site_title` | `site_title` | `.site.title` |
| `site_description` | `site_description` | `.site.description` |
| `site_url` | `site_url` | `.site.url` |
| `logo` | `logo` | `.site.logo` |
| `favicon` | `favicon` | `.site.favicon` |

---

**文档版本：** 1.3
**最后更新：** 2026-03-26

---

### 问题 19: `asset_path` 路径重复导致 MIME 类型错误

**错误信息：**
```
Refused to apply style from 'http://localhost:3000/themes/duocai-blog//themes/duocai-blog/assets/css/base.css'
because its MIME type ('text/html') is not a supported stylesheet MIME type
```

**原因：**
`asset_path` 函数会自动添加 `/themes/{theme-name}/` 前缀，但传入的路径已经包含了这部分，导致路径重复。

**错误写法：**
```html
{{ asset_path "/themes/duocai-blog/assets/css/base.css" }}
<!-- 结果: /themes/duocai-blog//themes/duocai-blog/assets/css/base.css -->
```

**正确写法：**
```html
{{ asset_path "assets/css/base.css" }}
<!-- 结果: /themes/duocai-blog/assets/css/base.css -->
```

**`asset_path` 函数实现：**
```go
"asset_path": func(path string) string {
    if theme == nil {
        return "/" + path
    }
    return "/themes/" + theme.Name + "/" + path
},
```

**影响文件：**
- `templates/themes/duocai-blog/pages/home.html`
- `templates/themes/duocai-blog/pages/article_show.html`
- `templates/themes/duocai-blog/pages/article_list.html`
- `templates/themes/duocai-blog/pages/category.html`
- `templates/themes/duocai-blog/pages/tag.html`
- `templates/themes/duocai-blog/pages/tags.html`
- `templates/themes/duocai-blog/pages/archives.html`
- `templates/themes/duocai-blog/pages/categories.html`
- `templates/themes/duocai-blog/pages/search.html`
- `templates/themes/duocai-blog/pages/page.html`
- `templates/themes/duocai-blog/pages/404.html`

**修复日期：** 2026-03-26

---

## 新增问题（2026-03-26）

### 问题 18: `cannot index slice/array with type models.Post`

**错误信息：**
```
template execute error: template: home.html:240:61: executing "home.html" at <index $.posts .>: error calling index: cannot index slice/array with type models.Post
```

**原因：**
在 `range` 循环中，`.` 是**当前迭代对象**，不是索引。错误地使用 `index $.posts .` 试图用 Post 对象作为索引。

**错误写法：**
```html
{{ range first 5 .posts }}
    <span class="popular-number">{{ add (index $.posts .) 1 }}</span>
    <!-- ❌ 错误：. 是 Post 对象，不能用作数组索引 -->
```

**正确写法：**
```html
{{ range $index, $post := first 5 .posts }}
    <span class="popular-number">{{ add $index 1 }}</span>
    <!-- ✅ 正确：使用 $index 获取位置 -->
```

**Go 模板 range 循环语法：**
```html
{{ range $index, $item := collection }}
    {{ $index }}  <!-- 当前索引（从0开始）-->
    {{ $item }}   <!-- 当前元素 -->
{{ end }}

{{ range $item := collection }}
    {{ . }}       <!-- 当前元素（隐式 dot）-->
{{ end }}
```

**影响文件：**
- `templates/themes/duocai-blog/pages/home.html`

**修复日期：** 2026-03-26
