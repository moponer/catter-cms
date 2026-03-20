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
            <p>定制化软件开发服务...</p>
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
主题 = 布局 + 样式 + 配色 + 动画
       ↓
用户 = 公司信息 + 服务项目 + 团队成员 + 联系方式
       ↓
网站 = 独一无二的企业官网
```

### 3.2 新的主题结构

```
enterprise-website/
├── templates/
│   ├── home.html
│   ├── page_default.html
│   ├── page_services.html
│   ├── page_team.html
│   ├── page_contact.html
│   └── modules/
│       ├── hero.html
│       ├── services.html
│       ├── stats.html
│       └── news.html
├── static/
├── theme.yaml
└── settings.yaml
```

---

## 四、页面模板系统

### 4.1 多模板支持

```yaml
# theme.yaml
customTemplates:
  page:
    - name: "默认页面"
      file: "page_default.html"
      description: "标准内容页布局"
    - name: "服务列表"
      file: "page_services.html"
      description: "带图片卡片的服务展示页"
    - name: "团队介绍"
      file: "page_team.html"
      description: "团队成员卡片网格布局"
    - name: "联系我们"
      file: "page_contact.html"
      description: "带表单和地图的联系页"
    - name: "图片画廊"
      file: "page_gallery.html"
      description: "图片网格展示"
    - name: "FAQ问答"
      file: "page_faq.html"
      description: "常见问题折叠面板"
```

### 4.2 用户操作

创建页面时可选择模板：
- 默认页面 → 标准内容布局
- 服务列表 → 服务卡片展示
- 团队介绍 → 成员卡片网格

---

## 五、模块化首页系统

### 5.1 首页模块定义

```yaml
modules:
  - id: "hero"
    name: "首屏大图"
    description: "页面顶部的大图区域"
    enabled: true
    order: 1
  - id: "services"
    name: "服务展示"
    description: "公司服务项目卡片"
    enabled: true
    order: 2
  - id: "stats"
    name: "统计数据"
    description: "关键业务数据展示"
    enabled: true
    order: 3
  - id: "news"
    name: "新闻动态"
    description: "最新文章或手动配置"
    enabled: true
    order: 4
  - id: "contact"
    name: "联系我们"
    description: "联系方式和表单"
    enabled: true
    order: 5
```

### 5.2 模块内容配置

```yaml
# 服务模块配置（Repeater）
services:
  enabled: true
  title: "我们的服务"
  subtitle: "为企业提供全方位的解决方案"
  items:
    - icon: "🚛"
      title: "仓储物流"
      description: "专业的仓储管理服务"
      image: ""
    - icon: "❄️"
      title: "冷链物流"
      description: "全程温控运输"
      image: ""
    - icon: "📦"
      title: "快递配送"
      description: "高效配送网络"
      image: ""

# 统计模块配置（Repeater）
stats:
  enabled: true
  items:
    - number: "1000+"
      label: "服务客户"
    - number: "50万+"
      label: "年配送量"
    - number: "99%"
      label: "准时率"
```

### 5.3 模板渲染

```html
<!-- home.html -->
{{ if .settings.hero.enabled }}
    {{ include "modules/hero.html" }}
{{ end }}

{{ if .settings.services.enabled }}
    {{ include "modules/services.html" }}
{{ end }}

{{ if .settings.stats.enabled }}
    {{ include "modules/stats.html" }}
{{ end }}
```

---

## 六、侧边栏小组件

### 6.1 小组件定义

```yaml
sidebar:
  enabled: true
  position: "right"
  widgets:
    - type: "search"
      title: "搜索文章"
      enabled: true
    - type: "categories"
      title: "文章分类"
      enabled: true
    - type: "tags"
      title: "热门标签"
      enabled: true
    - type: "recent_posts"
      title: "最新文章"
      enabled: true
      config:
        count: 5
    - type: "author"
      title: "关于博主"
      enabled: true
      config:
        name: "博主"
        bio: "..."
        avatar: ""
    - type: "html"
      title: "自定义内容"
      enabled: false
      config:
        content: "<div>自定义HTML</div>"
```

---

## 七、管理后台 UI

### 7.1 主题配置页面

```
┌─────────────────────────────────────────────────────────────┐
│  企业官网主题配置                              [保存配置]      │
├─────────────────────────────────────────────────────────────┤
│  [外观设置] [首页模块] [页面模板] [侧边栏] [自定义代码]        │
│                                                             │
│  📄 页面模板                                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ [预览图]     │ │ [预览图]     │ │ [预览图]     │           │
│  │ 默认页面     │ │ 服务列表     │ │ 团队介绍     │           │
│  │ ● 已选择    │ │ ○           │ │ ○           │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
│                                                             │
│  🧱 首页模块（拖拽排序）                                      │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ☰ 首屏大图    [启用 ✓]  [编辑]                        │   │
│  │ ☰ 服务展示    [启用 ✓]  [编辑]                        │   │
│  │ ☰ 统计数据    [启用 ✓]  [编辑]                        │   │
│  │ ☰ 新闻动态    [启用 ✓]  [编辑]                        │   │
│  │ ☰ 联系我们    [启用 ✓]  [编辑]                        │   │
│  └─────────────────────────────────────────────────────┘   │
│  [ + 添加模块 ]                                              │
└─────────────────────────────────────────────────────────────┘
```

### 7.2 Repeater 编辑器

```
┌─────────────────────────────────────────────────────────────┐
│  编辑: 服务展示模块                    [保存] [取消]          │
├─────────────────────────────────────────────────────────────┤
│  区块标题: [我们的服务________________________________]       │
│                                                             │
│  服务项目列表:                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 项目 1                                    [↑] [↓] [×]│   │
│  │ 图标: [🚛]  标题: [仓储物流]  描述: [专业仓储服务]    │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 项目 2                                    [↑] [↓] [×]│   │
│  │ 图标: [❄️]  标题: [冷链物流]  描述: [全程温控运输]    │   │
│  └─────────────────────────────────────────────────────┘   │
│  [ + 添加项目 ]                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 八、实施路线图

| 阶段 | 任务 | 时间 |
|------|------|------|
| 1 | 基础架构：repeater、模块系统、条件显示 | 1周 |
| 2 | 主题改造：拆分模块、创建模板、settings.yaml | 2周 |
| 3 | 管理后台：repeater组件、拖拽排序、图片选择 | 1周 |
| 4 | 完善测试：改造blog主题、文档、端到端测试 | 1周 |

---

## 九、预期效果

**物流公司使用企业官网主题：**
1. 选择"默认页面"模板
2. 修改公司名称、Logo、联系方式
3. 编辑服务项目：仓储物流、冷链物流、快递配送
4. 修改统计数据：1000+客户、50万+配送量
5. 上传公司照片
6. **完成！无需修改任何代码**

---

*文档完*

---

## 十、AI 原生支持方案

### 10.1 设计理念

```
用户：帮我写一篇关于"物流行业趋势"的文章
AI：调用 CMS API 创建文章...
    ↓
CMS：创建文章成功，自动生成 SEO，推荐标签
```

**核心原则：**
- CMS 提供完整的 REST API
- 用户在管理后台生成 API Key
- AI 通过 API Key 认证调用 API
- API Key 只存储在用户本地，不传给 AI

### 10.2 架构

```
┌─────────────────────────────────────────────────────────────┐
│                    AI 客户端（本地）                          │
│  ChatGPT / Claude / Cursor / 其他 AI                        │
│  ↓                                                          │
│  配置: API_ENDPOINT + API_KEY（环境变量，不传给 AI）          │
├─────────────────────────────────────────────────────────────┤
│                    Catter CMS API                           │
│  /api/v1/...（API Key 认证）                                 │
│  /api/v1/ai/capabilities（AI 能力发现）                      │
│  /api/v1/ai/guide（AI 操作指南）                             │
│  /api/v1/ai/docs（AI 详细文档）                              │
└─────────────────────────────────────────────────────────────┘
```

### 10.3 密钥安全

**存储位置：**
- 系统环境变量（最安全）
- 本地配置文件（~/.catter/config）
- AI 客户端配置文件

**安全原则：**
- API Key 永远不传给 AI
- 不出现在对话中
- 不出现在 AI 上下文中
- 只在本地环境使用

---

## 十一、API 密钥管理

### 11.1 管理界面

```
┌─────────────────────────────────────────────────────────────┐
│  API 密钥管理                                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  API 端点: https://your-cms.com/api/v1                      │
│                                                             │
│  我的 API Keys:                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 🔑 AI助手专用 - 2026-03-20                           │   │
│  │    sk-a1b2c3d4e5f6... [复制]                         │   │
│  │    权限: 文章、页面、分类、标签                        │   │
│  │    状态: ✅ 启用  使用: 42 次                         │   │
│  │    [禁用] [删除] [重置]                               │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  [ + 创建新的 API Key ]                                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 11.2 权限控制

```yaml
permissions:
  - name: "文章管理"
    key: "posts"
    actions: ["read", "create", "update", "delete"]
  - name: "页面管理"
    key: "pages"
    actions: ["read", "create", "update", "delete"]
  - name: "分类管理"
    key: "categories"
    actions: ["read", "create", "update", "delete"]
  - name: "标签管理"
    key: "tags"
    actions: ["read", "create", "update", "delete"]
  - name: "菜单管理"
    key: "menus"
    actions: ["read", "update"]
  - name: "主题配置"
    key: "themes"
    actions: ["read", "update"]
  - name: "系统设置"
    key: "settings"
    actions: ["read"]  # 只读
```

---

## 十二、AI 文档接口

### 12.1 能力发现

```
GET /api/v1/ai/capabilities

返回：
{
  "name": "Catter CMS",
  "version": "1.0.0",
  "capabilities": [
    {
      "id": "posts",
      "name": "文章管理",
      "description": "创建、编辑、删除、查询文章",
      "actions": ["list", "get", "create", "update", "delete"]
    },
    {
      "id": "pages",
      "name": "页面管理",
      "description": "创建、编辑、删除、查询页面",
      "actions": ["list", "get", "create", "update", "delete"]
    }
  ]
}
```

### 12.2 操作指南

```
GET /api/v1/ai/guide

返回：
{
  "quick_start": {
    "step_1": "获取 API Key（管理后台 → API 密钥管理）",
    "step_2": "设置环境变量 CATTER_API_KEY",
    "step_3": "使用 Bearer token 认证调用 API"
  },
  "common_tasks": [
    {
      "task": "创建一篇新文章",
      "method": "POST",
      "path": "/api/v1/posts",
      "required_fields": ["title", "content"],
      "optional_fields": ["status", "category_id", "tag_ids"],
      "example": {
        "title": "文章标题",
        "content": "文章内容（Markdown）",
        "status": "published"
      }
    }
  ],
  "restrictions": [
    "不能删除管理员用户",
    "批量操作限制：单次最多 100 条"
  ]
}
```

### 12.3 详细文档

```
GET /api/v1/ai/docs

返回完整 API 文档，包括：
- 所有端点的详细说明
- 请求参数和响应格式
- 错误代码说明
- 使用示例

GET /api/v1/ai/docs/{operation}

返回特定操作的详细说明：
- 参数约束
- 响应示例
- 常见错误
- 相关操作
```

### 12.4 OpenAPI 兼容

```
GET /api/v1/ai/openapi.json

返回标准 OpenAPI 3.0 格式，兼容：
- ChatGPT Custom GPT Actions
- Claude Tool Use
- Cursor API 集成
- 其他支持 OpenAPI 的 AI 工具
```

---

## 十三、AI 使用流程

### 13.1 用户配置

```
1. 登录管理后台 → API 密钥管理
2. 创建 API Key，选择权限
3. 复制 API Key
4. 配置本地环境变量：
   # Windows
   $env:CATTER_API_KEY="sk-xxxxx"
   $env:CATTER_API_URL="https://your-cms.com/api/v1"
   
   # macOS/Linux
   export CATTER_API_KEY="sk-xxxxx"
   export CATTER_API_URL="https://your-cms.com/api/v1"
5. 告诉 AI："使用 CMS API 完成以下任务..."
```

### 13.2 AI 工作流程

```
AI 收到任务："帮我创建一篇文章"
    ↓
1. 调用 GET /api/v1/ai/capabilities（发现可用能力）
    ↓
2. 调用 GET /api/v1/ai/docs/posts/create（了解如何创建文章）
    ↓
3. 调用 GET /api/v1/categories（获取现有分类）
    ↓
4. 调用 POST /api/v1/posts（创建文章）
    ↓
5. 返回结果给用户
```

### 13.3 常见使用场景

**场景 1：AI 写文章**
```
用户：帮我写一篇关于"2026年物流趋势"的文章
AI：创建文章，自动填充标题、内容、分类、标签
```

**场景 2：AI 调整配置**
```
用户：把网站主色调改成蓝色
AI：调用 update_theme_config 更新颜色配置
```

**场景 3：AI 管理内容**
```
用户：把所有文章都加上"AI"标签
AI：批量更新文章标签
```

**场景 4：AI 开发主题**
```
用户：创建一个新的企业官网主题
AI：生成主题文件，调用 API 注册主题
```

---

## 十四、安全措施

| 措施 | 说明 |
|------|------|
| 环境变量存储 | 密钥不传给 AI，只在本地 |
| 细粒度权限 | 只授予必要权限 |
| 使用统计 | 记录每次 API 调用 |
| 禁用/删除 | 随时可以禁用或删除 Key |
| 重置功能 | 泄露时可立即重置 |
| IP 白名单（可选） | 限制访问来源 |
| 速率限制 | 防止滥用 |
| 操作日志 | 记录所有 AI 操作 |

---

## 十五、实施路线图（更新）

| 阶段 | 任务 | 时间 |
|------|------|------|
| 1 | 基础架构：repeater、模块系统、条件显示 | 1周 |
| 2 | 主题改造：拆分模块、创建模板、settings.yaml | 2周 |
| 3 | 管理后台：repeater组件、拖拽排序、图片选择 | 1周 |
| 4 | API 密钥管理、权限控制 | 3天 |
| 5 | AI 文档接口（capabilities/guide/docs） | 2天 |
| 6 | 完善测试：改造blog主题、文档、端到端测试 | 1周 |

---

## 十六、预期效果

### 主题系统

**物流公司使用企业官网主题：**
1. 选择"默认页面"模板
2. 修改公司名称、Logo、联系方式
3. 编辑服务项目：仓储物流、冷链物流、快递配送
4. 修改统计数据：1000+客户、50万+配送量
5. **完成！无需修改任何代码**

### AI 支持

**AI 自动创建文章：**
1. 用户配置 API Key
2. 告诉 AI："创建一篇关于物流的文章"
3. AI 自动调用 API 创建文章
4. 文章自动发布到前台
5. **完成！无需手动操作**

---

*文档完*
