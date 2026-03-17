# Catter CMS 主题配置示例

这个文档展示了如何在 Catter CMS 中创建类似 Halo 的主题配置。

## 配置文件结构

主题配置文件位于主题目录下的 `theme.yml`。

### 示例：类似 Halo Earth 主题的配置

```yaml
# ============================================
# 主题基本信息
# ============================================
name: "Earth"
version: "1.0.0"
description: "一个简洁优雅的博客主题"
author: "Your Name"
authorEmail: "your.email@example.com"
authorUrl: "https://yourwebsite.com"
license: "MIT"
tags:
  - "博客"
  - "极简"
  - "响应式"
thumbnail: "screenshot.png"

# ============================================
# 主题配置定义（在管理后台显示的配置项）
# ============================================
configDefinition:
  # ============================================
  # 第一组：布局设置
  # ============================================
  - section: "布局设置"
    icon: "grid"  # 使用 Element Plus 图标名称
    description: "配置网站的布局方式"
    fields:
      - name: "layout_style"
        type: "select"
        label: "布局样式"
        description: "选择网站的布局方式"
        options:
          - label: "两栏布局（侧边栏在右侧）"
            value: "two-columns-right"
          - label: "两栏布局（侧边栏在左侧）"
            value: "two-columns-left"
          - label: "单栏布局（无侧边栏）"
            value: "single-column"
          - label: "三栏布局（双侧边栏）"
            value: "three-columns"
        default: "two-columns-right"

      - name: "content_width"
        type: "number"
        label: "内容宽度"
        description: "文章内容区域的最大宽度（像素）"
        min: 600
        max: 1200
        step: 50
        default: 800

      - name: "enable_breadcrumb"
        type: "switch"
        label: "显示面包屑导航"
        description: "在文章页面顶部显示面包屑导航"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

  # ============================================
  # 第二组：样式设置
  # ============================================
  - section: "样式设置"
    icon: "brush"
    description: "自定义主题的外观样式"
    fields:
      - name: "color_scheme"
        type: "select"
        label: "配色方案"
        description: "选择主题的配色方案"
        options:
          - label: "默认蓝"
            value: "default-blue"
          - label: "清新绿"
            value: "fresh-green"
          - label: "温暖橙"
            value: "warm-orange"
          - label: "优雅紫"
            value: "elegant-purple"
        default: "default-blue"

      - name: "border_radius"
        type: "slider"
        label: "圆角大小"
        description: "调整卡片和按钮的圆角大小"
        min: 0
        max: 20
        step: 1
        default: 8

      - name: "font_size_base"
        type: "select"
        label: "基础字号"
        description: "选择全局基础字体大小"
        options:
          - label: "14px（小）"
            value: "14"
          - label: "16px（标准）"
            value: "16"
          - label: "18px（大）"
            value: "18"
        default: "16"

      - name: "article_image_style"
        type: "radio"
        label: "文章封面图样式"
        description: "选择文章封面的显示样式"
        options:
          - label: "圆角"
            value: "rounded"
          - label: "方角"
            value: "square"
          - label: "圆形"
            value: "circle"
        default: "rounded"

      - name: "enable_card_shadow"
        type: "switch"
        label: "卡片阴影"
        description: "为文章卡片添加阴影效果"
        activeText: "启用"
        inactiveText: "禁用"
        default: true

      - name: "enable_gradient_background"
        type: "switch"
        label: "渐变背景"
        description: "使用渐变色作为背景"
        activeText: "启用"
        inactiveText: "禁用"
        default: false

  # ============================================
  # 第三组：侧边栏设置
  # ============================================
  - section: "侧边栏设置"
    icon: "menu"
    description: "配置侧边栏的内容和样式"
    fields:
      - name: "sidebar_position"
        type: "radio"
        label: "侧边栏位置"
        description: "设置侧边栏在页面中的位置"
        options:
          - label: "右侧"
            value: "right"
          - label: "左侧"
            value: "left"
        default: "right"

      - name: "enable_recent_posts"
        type: "switch"
        label: "显示最新文章"
        description: "在侧边栏显示最新的文章列表"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "recent_posts_count"
        type: "number"
        label: "最新文章数量"
        description: "显示的最新文章数量"
        min: 3
        max: 20
        step: 1
        default: 5

      - name: "enable_categories"
        type: "switch"
        label: "显示分类列表"
        description: "在侧边栏显示文章分类"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "enable_tags"
        type: "switch"
        label: "显示标签云"
        description: "在侧边栏显示文章标签"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "enable_search_box"
        type: "switch"
        label: "显示搜索框"
        description: "在侧边栏显示搜索输入框"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "search_placeholder"
        type: "text"
        label: "搜索框提示文字"
        placeholder: "请输入搜索关键词"
        default: "搜索文章..."

  # ============================================
  # 第四组：页脚设置
  # ============================================
  - section: "页脚设置"
    icon: "bottom"
    description: "自定义页脚信息"
    fields:
      - name: "footer_style"
        type: "select"
        label: "页脚样式"
        description: "选择页脚的显示方式"
        options:
          - label: "简洁型"
            value: "simple"
          - label: "完整型（包含所有链接）"
            value: "full"
          - label: "自定义HTML"
            value: "custom"
        default: "full"

      - name: "footer_text"
        type: "textarea"
        label: "页脚文字"
        description: "在页脚显示的自定义文字"
        rows: 3
        placeholder: "输入页脚文字，支持HTML"
        default: "&copy; 2024 我的博客. 保留所有权利."

      - name: "show_powered_by"
        type: "switch"
        label: "显示Powered By"
        description: "在页脚显示Powered by Catter CMS"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "footer_columns"
        type: "number"
        label: "页脚列数"
        description: "页脚内容的列数（仅在完整型样式下有效）"
        min: 1
        max: 4
        step: 1
        default: 3

      - name: "show_social_links"
        type: "switch"
        label: "显示社交链接"
        description: "在页脚显示社交媒体链接"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

  # ============================================
  # 第五组：备案设置
  # ============================================
  - section: "备案设置"
    icon: "document"
    description: "配置网站备案信息"
    fields:
      - name: "show_icp"
        type: "switch"
        label: "显示ICP备案号"
        description: "在页脚显示ICP备案信息"
        activeText: "显示"
        inactiveText: "隐藏"
        default: false

      - name: "icp_number"
        type: "text"
        label: "ICP备案号"
        placeholder: "请输入ICP备案号，例如：京ICP备12345678号"
        description: "您的网站ICP备案编号"

      - name: "icp_link"
        type: "text"
        label: "ICP链接"
        placeholder: "工信部备案系统链接"
        default: "https://beian.miit.gov.cn/"

      - name: "show_police_icp"
        type: "switch"
        label: "显示公安备案号"
        description: "在页脚显示公安备案信息"
        activeText: "显示"
        inactiveText: "隐藏"
        default: false

      - name: "police_icp_number"
        type: "text"
        label: "公安备案号"
        placeholder: "请输入公安备案号，例如：京公网安备11010502012345号"

      - name: "police_icp_icon"
        type: "image"
        label: "公安备案图标"
        description: "上传公安备案的图标文件"

  # ============================================
  # 第六组：文章设置
  # ============================================
  - section: "文章设置"
    icon: "document"
    description: "配置文章的显示方式"
    fields:
      - name: "article_list_style"
        type: "select"
        label: "文章列表样式"
        options:
          - label: "列表视图"
            value: "list"
          - label: "卡片视图"
            value: "card"
          - label: "缩略图视图"
            value: "thumbnail"
        default: "card"

      - name: "show_article_summary"
        type: "switch"
        label: "显示摘要"
        description: "在文章列表中显示文章摘要"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "summary_length"
        type: "number"
        label: "摘要长度"
        description: "文章摘要的最大字符数"
        min: 50
        max: 500
        step: 10
        default: 200

      - name: "show_read_time"
        type: "switch"
        label: "显示阅读时间"
        description: "在文章列表中显示预计阅读时间"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "show_views_count"
        type: "switch"
        label: "显示浏览次数"
        description: "在文章列表中显示文章浏览次数"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

  # ============================================
  # 第七组：高级设置
  # ============================================
  - section: "高级设置"
    icon: "setting"
    description: "高级配置选项"
    fields:
      - name: "custom_css"
        type: "code"
        label: "自定义CSS"
        description: "添加自定义CSS代码来覆盖主题样式"
        language: "css"
        rows: 10
        placeholder: "/* 在这里输入自定义CSS */"

      - name: "custom_js"
        type: "code"
        label: "自定义JavaScript"
        description: "添加自定义JavaScript代码"
        language: "javascript"
        rows: 10
        placeholder: "// 在这里输入自定义JavaScript"

      - name: "enable_pjax"
        type: "switch"
        label: "启用PJAX"
        description: "使用PJAX实现页面无刷新切换"
        activeText: "启用"
        inactiveText: "禁用"
        default: false

      - name: "enable_lazy_load"
        type: "switch"
        label: "图片懒加载"
        description: "延迟加载页面图片以提高性能"
        activeText: "启用"
        inactiveText: "禁用"
        default: true

# ============================================
# 主题默认配置值（用户修改后的值保存在这里）
# ============================================
settings:
  layout_style: "two-columns-right"
  content_width: 800
  enable_breadcrumb: true
  color_scheme: "default-blue"
  border_radius: 8
  font_size_base: "16"
  article_image_style: "rounded"
  enable_card_shadow: true
  enable_gradient_background: false
  sidebar_position: "right"
  enable_recent_posts: true
  recent_posts_count: 5
  enable_categories: true
  enable_tags: true
  enable_search_box: true
  search_placeholder: "搜索文章..."
  footer_style: "full"
  footer_text: "&copy; 2024 我的博客. 保留所有权利."
  show_powered_by: true
  footer_columns: 3
  show_social_links: true
  show_icp: false
  icp_number: ""
  icp_link: "https://beian.miit.gov.cn/"
  show_police_icp: false
  police_icp_number: ""
  police_icp_icon: ""
  article_list_style: "card"
  show_article_summary: true
  summary_length: 200
  show_read_time: true
  show_views_count: true
  custom_css: ""
  custom_js: ""
  enable_pjax: false
  enable_lazy_load: true
```

## 配置字段类型说明

### 1. text - 文本输入
```yaml
- name: "site_name"
  type: "text"
  label: "网站名称"
  placeholder: "请输入网站名称"
  default: "我的网站"
```

### 2. textarea - 文本域
```yaml
- name: "description"
  type: "textarea"
  label: "网站描述"
  rows: 3
  placeholder: "请输入网站描述"
```

### 3. number - 数字输入
```yaml
- name: "posts_per_page"
  type: "number"
  label: "每页文章数"
  min: 5
  max: 50
  step: 1
  default: 10
```

### 4. color - 颜色选择器
```yaml
- name: "primary_color"
  type: "color"
  label: "主色调"
  showAlpha: false
  default: "#3B82F6"
```

### 5. switch - 开关
```yaml
- name: "enable_dark_mode"
  type: "switch"
  label: "启用深色模式"
  activeText: "启用"
  inactiveText: "禁用"
  default: false
```

### 6. select - 下拉选择
```yaml
- name: "layout"
  type: "select"
  label: "布局方式"
  options:
    - label: "两栏布局"
      value: "two-columns"
    - label: "单栏布局"
      value: "single-column"
  default: "two-columns"
```

### 7. multi-select - 多选下拉
```yaml
- name: "tags"
  type: "multi-select"
  label: "标签选择"
  options:
    - label: "技术"
      value: "tech"
    - label: "生活"
      value: "life"
```

### 8. radio - 单选按钮
```yaml
- name: "sidebar_position"
  type: "radio"
  label: "侧边栏位置"
  options:
    - label: "左侧"
      value: "left"
    - label: "右侧"
      value: "right"
  default: "right"
```

### 9. checkbox - 复选框
```yaml
- name: "features"
  type: "checkbox"
  label: "启用功能"
  options:
    - label: "评论"
      value: "comments"
    - label: "分享"
      value: "share"
```

### 10. image - 图片上传
```yaml
- name: "logo"
  type: "image"
  label: "网站Logo"
  description: "上传网站Logo图片"
```

### 11. datetime - 日期时间
```yaml
- name: "event_date"
  type: "datetime"
  label: "活动时间"
  placeholder: "选择日期时间"
```

### 12. daterange - 日期范围
```yaml
- name: "date_range"
  type: "daterange"
  label: "日期范围"
  placeholder: "选择日期范围"
```

### 13. slider - 滑块
```yaml
- name: "border_radius"
  type: "slider"
  label: "圆角大小"
  min: 0
  max: 20
  step: 1
  default: 8
```

### 14. code - 代码编辑器
```yaml
- name: "custom_css"
  type: "code"
  label: "自定义CSS"
  language: "css"
  rows: 10
  placeholder: "/* 输入CSS代码 */"
```

## 在模板中使用配置

```plush-html
<!-- 获取主题配置 -->
<%= theme.Settings.LayoutStyle %>
<%= theme.Settings.ContentWidth %>
<%= theme.Settings.EnableBreadcrumb %>

<!-- 在CSS中使用变量 -->
<style>
  :root {
    --content-width: <%= theme.Settings.ContentWidth %>px;
    --border-radius: <%= theme.Settings.BorderRadius %>px;
  }

  .container {
    max-width: var(--content-width);
  }

  .card {
    border-radius: var(--border-radius);
  }
</style>

<!-- 在HTML中使用开关 -->
<% if (theme.Settings.EnableBreadcrumb) { %>
  <nav class="breadcrumb">
    <!-- 面包屑导航 -->
  </nav>
<% } %>
```

## 注意事项

1. **configDefinition** 是主题开发者定义的配置表单
2. **settings** 是用户修改后的配置值
3. 字段的 `name` 必须唯一
4. 所有配置项都应该有默认值
5. 分组的 `section` 名称要清晰易懂
6. 使用描述文字帮助用户理解每个配置项

## API 说明

### 获取主题配置
```
GET /api/themes/{slug}/config/full
```

响应示例：
```json
{
  "config": {
    "configForm": [...],  // configDefinition（配置表单定义）
    "settings": {...}     // 用户配置的值
  }
}
```

### 更新主题配置
```
PUT /api/themes/{slug}/config/full
```

请求体：
```json
{
  "settings": {
    "layout_style": "single-column",
    "content_width": 900
  }
}
```

## 对比 Halo 主题配置

Catter CMS 的主题配置系统与 Halo 非常相似：

| 功能 | Halo | Catter CMS |
|------|------|------------|
| 配置分组 | ✅ | ✅ |
| 多种字段类型 | ✅ | ✅ (14种) |
| 默认值 | ✅ | ✅ |
| 必填验证 | ✅ | ✅ |
| 条件显示 | ✅ | ✅ (showIf) |
| 图标支持 | ✅ | ✅ |
| 自定义CSS/JS | ✅ | ✅ |

Catter CMS 完全支持像 Halo Earth 主题那样的配置方式！
