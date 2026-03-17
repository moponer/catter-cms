# Catter CMS 主题配置最佳实践指南

## 概述

本文档基于Halo CMS主题配置系统的最佳实践，为Catter CMS提供完整的主题配置指南。确保后台配置可以实时改变前台显示。

## 核心概念

### 1. 配置文件结构

主题配置文件位于 `templates/themes/{theme-slug}/theme.yml`，包含以下主要部分：

```yaml
# 主题基本信息
name: "主题名称"
version: "2.0.0"
apiVersion: "2.0"
description: "主题描述"
author: "作者"
license: "MIT"

# 主题能力声明
capabilities:
  layouts: []
  modules: []
  widgets: []
  postFormats: []
  customFields: true
  colorSchemes: true
  darkMode: false

# 默认配置（包含所有模块、页面、SEO等配置）
config:
  appearance:
    colors: {}
    typography: {}
    layout: {}
  modules:
    hero:
      enabled: true
      config: {}
  pages:
    home:
      title: "首页"
      enabled: true
  seo:
    title: "网站标题"
    description: "网站描述"
  custom:
    css: ""
    js: ""
  advanced:
    performance: {}
    social: {}

# 配置定义（用于生成管理后台配置表单）
configDefinition:
  - section: "布局设置"
    icon: "grid"
    fields:
      - name: "config.path"
        type: "text"
        label: "配置项名称"

# 菜单位置定义
menuLocations:
  header:
    name: "顶部导航"
    maxDepth: 3

# 小工具区域定义
widgetAreas:
  - id: "sidebar"
    name: "侧边栏"

# 自定义字段定义
customFields:
  post:
    - name: "subtitle"
      type: "text"
      label: "副标题"

# 主题特定设置
settings:
  colorSchemes: []
```

### 2. 关键配置说明

#### configDefinition - 配置表单定义

这是主题开发者定义的配置表单，用于在管理后台生成配置界面。

**重要原则**：
1. **字段命名必须与config路径对应**
   - `name: "appearance.colors.primary"` 对应 `config.appearance.colors.primary`
2. **支持条件显示**：使用 `show_if` 属性
3. **提供默认值**：每个字段都应该有 `default` 值
4. **清晰的描述**：使用 `description` 帮助用户理解

**示例**：
```yaml
configDefinition:
  - section: "外观设置"
    icon: "brush"
    description: "配置主题颜色、字体等"
    fields:
      - name: "appearance.colors.primary"
        type: "color"
        label: "主色调"
        default: "#3B82F6"
        description: "网站的主要品牌颜色"

      - name: "appearance.layout.sidebar"
        type: "select"
        label: "侧边栏位置"
        options:
          - label: "右侧"
            value: "right"
          - label: "左侧"
            value: "left"
        default: "right"

      - name: "modules.hero.enabled"
        type: "switch"
        label: "启用Hero模块"
        activeText: "显示"
        inactiveText: "隐藏"
        default: true

      - name: "modules.hero.config.title"
        type: "text"
        label: "主标题"
        default: "欢迎访问"
        show_if: "modules.hero.enabled"
```

#### config - 默认配置

这是主题的默认配置值，包含所有可配置的选项。

**结构建议**：
```yaml
config:
  # 外观配置
  appearance:
    colors:
      primary: "#3B82F6"
      secondary: "#764ba2"
      accent: "#06B6D4"
      background: "#ffffff"
      text: "#1f2937"
    typography:
      baseFont: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
      headingFont: "inherit"
      fontSize:
        base: "14px"
        scale: 1.125
    layout:
      width: "1200px"
      sidebar: "right"
      borderRadius: "8px"

  # 模块配置
  modules:
    hero:
      enabled: true
      config:
        title: "欢迎访问"
        subtitle: "这是一个示例"
        backgroundType: "gradient"

    features:
      enabled: true
      config:
        title: "核心特性"
        columns: 3

  # 页面配置
  pages:
    home:
      title: "首页"
      template: "home.plush.html"
      enabled: true
      config:
        showHero: true
        showFeatures: true

  # SEO配置
  seo:
    title: "网站标题"
    description: "网站描述"
    keywords: ["关键词1", "关键词2"]
    opengraph:
      enabled: true
      image: "/og-image.jpg"

  # 自定义代码
  custom:
    css: |
      /* 自定义CSS */
    js: |
      // 自定义JS

  # 高级配置
  advanced:
    performance:
      lazyLoadImages: true
    social:
      email: "contact@example.com"
    features:
      backToTop: true
      toc: true
```

### 3. 在模板中使用配置

在 `.plush.html` 模板中，可以通过 `theme.Settings` 访问用户配置的值。

**基础用法**：
```plush-html
<!-- 访问颜色配置 -->
<style>
  :root {
    --primary-color: <%= theme.Settings.Appearance.Colors.Primary %>;
    --secondary-color: <%= theme.Settings.Appearance.Colors.Secondary %>;
  }

  body {
    background-color: <%= theme.Settings.Appearance.Colors.Background %>;
    color: <%= theme.Settings.Appearance.Colors.Text %>;
  }
</style>

<!-- 访问布局配置 -->
<div class="container" style="max-width: <%= theme.Settings.Appearance.Layout.Width %>;">
  <!-- 内容 -->
</div>

<!-- 条件显示模块 -->
<% if (theme.Settings.Modules.Hero.Enabled) { %>
  <section class="hero">
    <h1><%= theme.Settings.Modules.Hero.Config.Title %></h1>
    <p><%= theme.Settings.Modules.Hero.Config.Subtitle %></p>
  </section>
<% } %>

<!-- 循环显示特性 -->
<% if (theme.Settings.Modules.Features.Enabled) { %>
  <section class="features">
    <h2><%= theme.Settings.Modules.Features.Config.Title %></h2>
    <div class="grid-<%= theme.Settings.Modules.Features.Config.Columns %>">
      <!-- 特性内容 -->
    </div>
  </section>
<% } %>

<!-- SEO Meta标签 -->
<title><%= theme.Settings.Seo.Title %></title>
<meta name="description" content="<%= theme.Settings.Seo.Description %>">
<meta name="keywords" content="<%= theme.Settings.Seo.Keywords %>">

<!-- 自定义CSS -->
<style>
  <%= raw(theme.Settings.Custom.CSS) %>
</style>

<!-- 自定义JS -->
<script>
  <%= raw(theme.Settings.Custom.JS) %>
</script>
```

## 配置类型详解

### 1. 布局设置

**适用场景**：博客、官网、企业站

**常见配置项**：
- 内容宽度
- 侧边栏位置
- 容器模式（盒状/全宽）
- 区块间距
- 圆角大小

**示例**：
```yaml
- section: "布局设置"
  icon: "grid"
  fields:
    - name: "appearance.layout.width"
      type: "text"
      label: "内容宽度"
      placeholder: "例如: 1200px, 90%"
      default: "1200px"

    - name: "appearance.layout.sidebar"
      type: "select"
      label: "侧边栏位置"
      options:
        - label: "右侧"
          value: "right"
        - label: "左侧"
          value: "left"
        - label: "隐藏"
          value: "none"
      default: "right"

    - name: "appearance.layout.borderRadius"
      type: "slider"
      label: "圆角大小"
      min: 0
      max: 20
      step: 1
      default: 8
```

### 2. 样式设置

**适用场景**：所有类型主题

**常见配置项**：
- 主色调
- 次要颜色
- 强调色
- 背景色
- 文字颜色
- 圆角大小
- 动画开关

**示例**：
```yaml
- section: "样式设置"
  icon: "brush"
  fields:
    - name: "appearance.colors.primary"
      type: "color"
      label: "主色调"
      default: "#3B82F6"
      description: "网站的主要品牌颜色"

    - name: "appearance.colors.secondary"
      type: "color"
      label: "次要颜色"
      default: "#764ba2"

    - name: "appearance.animations.enabled"
      type: "switch"
      label: "启用动画"
      activeText: "启用"
      inactiveText: "禁用"
      default: true
```

### 3. 模块配置

**适用场景**：官网、落地页

**常见模块**：
- Hero模块
- 特性展示
- 价格方案
- 团队介绍
- 客户评价
- 联系表单

**示例**：
```yaml
- section: "Hero设置"
  icon: "view"
  fields:
    - name: "modules.hero.enabled"
      type: "switch"
      label: "启用Hero模块"
      default: true

    - name: "modules.hero.config.title"
      type: "text"
      label: "主标题"
      default: "欢迎访问"
      show_if: "modules.hero.enabled"

    - name: "modules.hero.config.backgroundType"
      type: "select"
      label: "背景类型"
      options:
        - label: "渐变色"
          value: "gradient"
        - label: "纯色"
          value: "color"
        - label: "图片"
          value: "image"
      default: "gradient"
      show_if: "modules.hero.enabled"

    - name: "modules.hero.config.backgroundImage"
      type: "image"
      label: "背景图片"
      show_if: "modules.hero.enabled && modules.hero.config.backgroundType == 'image'"

    - name: "modules.hero.config.showCta"
      type: "switch"
      label: "显示按钮"
      default: true
      show_if: "modules.hero.enabled"

    - name: "modules.hero.config.ctaText"
      type: "text"
      label: "按钮文本"
      default: "立即开始"
      show_if: "modules.hero.enabled && modules.hero.config.showCta"
```

### 4. 侧边栏设置

**适用场景**：博客、文档站

**常见配置项**：
- 侧边栏位置
- 最新文章开关及数量
- 分类列表开关
- 标签云开关
- 搜索框开关

**示例**：
```yaml
- section: "侧边栏设置"
  icon: "menu"
  fields:
    - name: "sidebar.position"
      type: "select"
      label: "侧边栏位置"
      options:
        - label: "右侧"
          value: "right"
        - label: "左侧"
          value: "left"
      default: "right"

    - name: "widgets.search.enabled"
      type: "switch"
      label: "显示搜索框"
      default: true

    - name: "widgets.recentPosts.enabled"
      type: "switch"
      label: "显示最新文章"
      default: true

    - name: "widgets.recentPosts.config.count"
      type: "number"
      label: "文章数量"
      min: 3
      max: 20
      default: 5
      show_if: "widgets.recentPosts.enabled"
```

### 5. 页脚设置

**适用场景**：所有类型主题

**常见配置项**：
- 页脚样式（简洁/完整/自定义）
- 页脚文字
- Powered By开关
- 社交链接开关

**示例**：
```yaml
- section: "页脚设置"
  icon: "bottom"
  fields:
    - name: "footer.style"
      type: "select"
      label: "页脚样式"
      options:
        - label: "简洁型"
          value: "simple"
        - label: "完整型"
          value: "full"
      default: "full"

    - name: "footer.text"
      type: "textarea"
      label: "页脚文字"
      rows: 3
      default: "&copy; 2024 我的网站. 保留所有权利."

    - name: "footer.showPoweredBy"
      type: "switch"
      label: "显示Powered By"
      default: true
```

### 6. 备案设置

**适用场景**：中国境内的网站

**常见配置项**：
- ICP备案号
- 公安备案号
- 公安备案图标

**示例**：
```yaml
- section: "备案设置"
  icon: "document"
  fields:
    - name: "beian.showIcp"
      type: "switch"
      label: "显示ICP备案号"
      default: false

    - name: "beian.icpNumber"
      type: "text"
      label: "ICP备案号"
      placeholder: "京ICP备12345678号"
      show_if: "beian.showIcp"

    - name: "beian.showPoliceIcp"
      type: "switch"
      label: "显示公安备案号"
      default: false

    - name: "beian.policeIcpNumber"
      type: "text"
      label: "公安备案号"
      placeholder: "京公网安备11010502012345号"
      show_if: "beian.showPoliceIcp"

    - name: "beian.policeIcpIcon"
      type: "image"
      label: "公安备案图标"
      show_if: "beian.showPoliceIcp"
```

### 7. SEO设置

**适用场景**：所有类型主题

**常见配置项**：
- 网站标题
- 网站描述
- 关键词
- Open Graph图片

**示例**：
```yaml
- section: "SEO设置"
  icon: "search"
  fields:
    - name: "seo.title"
      type: "text"
      label: "网站标题"
      required: true
      default: "我的网站"

    - name: "seo.description"
      type: "textarea"
      label: "网站描述"
      rows: 3
      required: true
      default: "网站描述"

    - name: "seo.keywords"
      type: "text"
      label: "关键词"
      placeholder: "用逗号分隔多个关键词"

    - name: "seo.opengraph.image"
      type: "image"
      label: "Open Graph图片"
      description: "在社交媒体上分享时显示的图片"
```

### 8. 自定义代码

**适用场景**：所有类型主题

**常见配置项**：
- 自定义CSS
- 自定义JavaScript

**示例**：
```yaml
- section: "自定义代码"
  icon: "code"
  fields:
    - name: "custom.css"
      type: "code"
      label: "自定义CSS"
      language: "css"
      rows: 10
      description: "添加自定义CSS样式，会覆盖主题的默认样式"

    - name: "custom.js"
      type: "code"
      label: "自定义JavaScript"
      language: "javascript"
      rows: 10
      description: "添加自定义JavaScript代码"
```

## 字段类型完整列表

### 1. text - 文本输入
```yaml
- name: "site.name"
  type: "text"
  label: "网站名称"
  placeholder: "请输入网站名称"
  default: "我的网站"
  required: true
```

### 2. textarea - 文本域
```yaml
- name: "site.description"
  type: "textarea"
  label: "网站描述"
  rows: 3
  placeholder: "请输入网站描述"
  default: "这是我的网站"
```

### 3. number - 数字输入
```yaml
- name: "posts.perPage"
  type: "number"
  label: "每页文章数"
  min: 5
  max: 50
  step: 1
  default: 10
```

### 4. color - 颜色选择器
```yaml
- name: "colors.primary"
  type: "color"
  label: "主色调"
  showAlpha: false
  default: "#3B82F6"
```

### 5. switch - 开关
```yaml
- name: "features.darkMode"
  type: "switch"
  label: "启用深色模式"
  activeText: "启用"
  inactiveText: "禁用"
  default: false
```

### 6. select - 下拉选择
```yaml
- name: "layout.style"
  type: "select"
  label: "布局样式"
  options:
    - label: "两栏布局"
      value: "two-columns"
    - label: "单栏布局"
      value: "single-column"
  default: "two-columns"
```

### 7. multi-select - 多选下拉
```yaml
- name: "tags.selected"
  type: "multi-select"
  label: "标签选择"
  options:
    - label: "技术"
      value: "tech"
    - label: "生活"
      value: "life"
  default: ["tech"]
```

### 8. radio - 单选按钮
```yaml
- name: "sidebar.position"
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
- name: "features.enabled"
  type: "checkbox"
  label: "启用功能"
  options:
    - label: "评论"
      value: "comments"
    - label: "分享"
      value: "share"
  default: ["comments", "share"]
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
- name: "event.date"
  type: "datetime"
  label: "活动时间"
  placeholder: "选择日期时间"
```

### 12. daterange - 日期范围
```yaml
- name: "date.range"
  type: "daterange"
  label: "日期范围"
  placeholder: "选择日期范围"
```

### 13. slider - 滑块
```yaml
- name: "border.radius"
  type: "slider"
  label: "圆角大小"
  min: 0
  max: 20
  step: 1
  default: 8
```

### 14. code - 代码编辑器
```yaml
- name: "custom.css"
  type: "code"
  label: "自定义CSS"
  language: "css"
  rows: 10
  placeholder: "/* 输入CSS代码 */"
```

## 条件显示（show_if）

`show_if` 属性允许你根据其他字段的值来动态显示或隐藏字段。

**语法**：
- 简单条件：`show_if: "field_name == 'value'"`
- 布尔值：`show_if: "field_name"`
- 多条件：`show_if: "field1 == 'value' && field2"`
- 或条件：`show_if: "field1 == 'value' || field2 == 'value'"`

**示例**：
```yaml
- name: "hero.enabled"
  type: "switch"
  label: "启用Hero模块"
  default: true

- name: "hero.title"
  type: "text"
  label: "主标题"
  show_if: "hero.enabled"

- name: "hero.backgroundType"
  type: "select"
  label: "背景类型"
  options:
    - label: "渐变色"
      value: "gradient"
    - label: "纯色"
      value: "color"
    - label: "图片"
      value: "image"
  show_if: "hero.enabled"

- name: "hero.backgroundColor"
  type: "color"
  label: "背景颜色"
  show_if: "hero.enabled && hero.backgroundType == 'color'"

- name: "hero.backgroundImage"
  type: "image"
  label: "背景图片"
  show_if: "hero.enabled && hero.backgroundType == 'image'"
```

## 配置更新机制

### 1. 前端配置保存

用户在管理后台修改配置后，前端会调用：

```javascript
const saveConfig = async () => {
  await put(`/api/themes/${theme.slug}/config/full`, {
    settings: {
      // 更新的配置值
      "appearance.colors.primary": "#FF0000",
      "modules.hero.config.title": "新标题"
    }
  })
}
```

### 2. 后端配置更新

后端API会：
1. 接收配置更新
2. 更新 `theme.yml` 中的 `settings` 字段
3. 保存配置到文件系统
4. 返回成功响应

### 3. 前台实时生效

配置保存后，前台页面会：
1. 读取更新后的配置
2. 应用到模板渲染
3. 实时显示效果

## 最佳实践总结

### 1. 配置设计原则

✅ **清晰分组**
- 使用有意义的分组名称
- 每个分组包含相关的配置项
- 使用图标增强识别度

✅ **提供默认值**
- 所有字段都应该有默认值
- 默认值应该是最常用的配置
- 确保默认值让主题可用

✅ **详细描述**
- 使用 `description` 说明配置项的作用
- 提供使用示例
- 说明配置的影响范围

✅ **条件显示**
- 使用 `show_if` 减少不相关配置的显示
- 提供清晰的逻辑
- 避免过度复杂的条件

### 2. 命名规范

✅ **配置路径命名**
- 使用点号分隔的路径：`appearance.colors.primary`
- 保持一致性：`modules.hero.enabled`, `modules.features.enabled`
- 使用英文命名：`primary`, `secondary`, `background`

✅ **分组命名**
- 使用中文分组名称：`布局设置`, `样式设置`, `Hero设置`
- 使用有意义的名称：避免`分组1`, `配置A`
- 保持简洁：`SEO设置`, `自定义代码`

### 3. 性能优化

✅ **避免过度配置**
- 只配置必要的选项
- 避免创建太多可选配置
- 考虑默认覆盖机制

✅ **缓存配置**
- 配置应该被缓存
- 避免每次请求都读取文件
- 使用内存缓存提高性能

### 4. 用户体验

✅ **实时预览**
- 提供配置预览功能
- 配置改变时立即显示效果
- 使用预览窗口

✅ **配置重置**
- 提供重置为默认的功能
- 确认重置操作
- 保留备份

✅ **导入导出**
- 支持配置导入
- 支持配置导出
- 提供预设模板

## 常见问题

### Q1: 配置不生效怎么办？

**A**: 检查以下几点：
1. 配置字段名称是否与 `configDefinition` 中的 `name` 一致
2. 模板中是否正确访问了配置：`theme.Settings.配置项名`
3. 注意大小写：`Settings` 是大写S
4. 清除浏览器缓存
5. 重启应用

### Q2: 如何添加新的配置项？

**A**:
1. 在 `configDefinition` 中添加配置项定义
2. 在 `config` 中设置默认值
3. 在模板中使用配置
4. 保存并测试

### Q3: show_if 条件不工作？

**A**:
1. 检查条件语法是否正确
2. 确保引用的字段存在
3. 检查值的格式是否匹配
4. 使用调试工具查看实际值

### Q4: 如何支持多语言？

**A**:
1. 配置字段支持多语言选项
2. 使用语言标识符区分
3. 在模板中根据语言显示对应值
4. 提供语言切换功能

## 示例主题

### 1. Catter Blog（博客主题）

**特点**：
- 完整的博客功能
- Hero模块
- 文章列表
- 侧边栏
- 评论系统

**配置分组**：
- 外观设置
- Hero模块
- 文章列表
- 侧边栏设置
- SEO设置
- 自定义代码

### 2. Catter CMS Official（官网主题）

**特点**：
- 落地页设计
- 功能展示
- 价格方案
- 文档中心

**配置分组**：
- 布局设置
- 样式设置
- Hero设置
- 特性展示
- 价格方案
- SEO设置
- 自定义代码

### 3. Enterprise Website（企业官网）

**特点**：
- 企业展示
- 服务介绍
- 团队展示
- 客户评价
- 联系表单

**配置分组**：
- 布局设置
- 样式设置
- Hero设置
- 关于我们
- 服务设置
- 团队设置
- 联系设置
- SEO设置
- 自定义代码

## 总结

Catter CMS 的主题配置系统基于Halo CMS的最佳实践，提供了：

✅ **完整的配置支持**
- 14种字段类型
- 条件显示
- 默认值
- 详细描述

✅ **灵活的配置能力**
- 支持任意数量的配置分组
- 支持嵌套配置结构
- 支持自定义字段

✅ **优秀的用户体验**
- 分组折叠面板
- 实时预览
- 配置重置
- 导入导出

✅ **强大的功能**
- 后台配置实时改变前台
- 配置缓存
- 性能优化
- 多语言支持

遵循本文档的最佳实践，您可以创建功能强大、易于使用的主题配置系统！

---

**文档版本**: 1.0.0
**最后更新**: 2026-03-16
**维护者**: WorkBuddy Team
