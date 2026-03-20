# Catter CMS 主题系统升级方案

> 版本：1.0 | 日期：2026-03-20 | 状态：待确认

---

## 一、当前问题分析

### 1.1 核心问题：主题与内容未分离

当前主题模板中存在**大量硬编码的业务内容**：

```html
<!-- enterprise-website/pages/home.html -->
<section class="section">
    <h2 class="section-title">我们的服务</h2>  <!-- ❌ 硬编码标题 -->
    <div class="services-grid">
        <div class="service-card">
            <h3>软件开发</h3>  <!-- ❌ 硬编码服务名称 -->
            <p>定制化软件开发服务...</p>  <!-- ❌ 硬编码描述 -->
        </div>
    </div>
</section>
```

**后果**：物流公司使用此主题时，必须修改源码才能显示"仓储物流"、"冷链物流"等内容。

### 1.2 硬编码内容清单

| 主题 | 硬编码内容 | 影响 |
|------|------------|------|
| enterprise-website | 导航链接（服务、方案、产品） | 必须改源码 |
| enterprise-website | 6个服务项目（软件开发、云服务等） | 必须改源码 |
| enterprise-website | 统计数据（500+、98%） | 必须改源码 |
| catter-blog | 回退导航链接 | 影响较小 |
| catter-blog | "Powered by Catter CMS" | 应可配置 |

---

## 二、优秀实践参考

### 2.1 Halo CMS

**核心理念**：主题只提供布局和样式，内容由数据驱动

```
my-theme/
├── templates/
│   ├── index.html
│   └── fragment/
├── settings.yaml    # 用户配置定义
└── theme.yaml       # 主题元数据
```

### 2.2 WordPress

**核心理念**：模板标签 + 自定义字段 + Widgets

- Customizer：实时预览配置
- Widgets：可拖拽内容区块
- Custom Fields：扩展数据结构

---

## 三、升级方案

### 3.1 设计原则

```
┌─────────────────────────────────────────────────────┐
│         主题（Theme）                                 │
│  职责：布局、样式、配色、动画、交互                      │
│  不包含：业务内容、文案、数据                            │
└─────────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────────┐
│       内容数据（Content）                              │
│  来源：数据库、主题配置、用户输入                        │
└─────────────────────────────────────────────────────┘
```

### 3.2 新的主题结构

```
enterprise-website/
├── templates/
│   ├── home.html
│   ├── page.html
│   └── fragments/
│       ├── header.html
│       └── footer.html
├── static/
├── theme.yaml      # 主题元数据
└── settings.yaml   # 主题设置定义
```

### 3.3 新的配置定义结构

```yaml
sections:
  - name: "company"
    label: "公司信息"
    fields:
      - name: "company.name"
        label: "公司名称"
        type: "text"
      - name: "company.logo"
        label: "公司Logo"
        type: "image"

  - name: "hero"
    label: "首屏大图"
    fields:
      - name: "hero.title"
        label: "主标题"
        type: "text"
      - name: "hero.backgroundImage"
        label: "背景图片"
        type: "image"

  - name: "services"
    label: "服务项目"
    fields:
      - name: "services.items"
        label: "服务列表"
        type: "repeater"  # 可重复字段
        fields:
          - name: "icon"
            label: "图标"
          - name: "title"
            label: "服务名称"
          - name: "description"
            label: "服务描述"

  - name: "stats"
    label: "统计数据"
    fields:
      - name: "stats.items"
        label: "统计项"
        type: "repeater"
        fields:
          - name: "number"
            label: "数字"
          - name: "label"
            label: "描述"
```

### 3.4 模板改造示例

```html
<!-- home.html - 改造后 -->
<section class="section">
    <h2>{{ .settings.services.title }}</h2>
    <div class="services-grid">
        {{ range .settings.services.items }}
        <div class="service-card">
            <div class="icon">{{ .icon }}</div>
            <h3>{{ .title }}</h3>
            <p>{{ .description }}</p>
        </div>
        {{ end }}
    </div>
</section>
```

---

## 四、实施计划

| 阶段 | 任务 | 时间 |
|------|------|------|
| 1 | 重构配置数据结构，支持嵌套和 repeater | 2-3天 |
| 2 | 改造 enterprise-website 主题 | 3-5天 |
| 3 | 改造 catter-blog 主题 | 2-3天 |
| 4 | 管理后台 UI 组件（repeater、嵌套） | 2-3天 |
| 5 | 文档和示例 | 1-2天 |

---

## 五、预期效果

**改造前**：
```
用户：我想用企业官网主题做物流公司网站
开发者：需要改源码，把"软件开发"改成"仓储物流"...
```

**改造后**：
```
用户：我想用企业官网主题做物流公司网站
系统：请在主题配置中填写服务项目
用户：自己改！"软件开发"→"仓储物流"
```

---

*待确认后实施*
