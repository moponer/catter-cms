# Catter CMS theme.yml 完整示例集

> 包含三种典型主题类型的完整 theme.yml 示例。

---

## 示例一：个人博客主题

```yaml
name: My Clean Blog
slug: my-clean-blog
version: 1.0.0
author: 开发者名称
description: 简洁现代的个人博客主题，专注内容阅读体验

config:
  # 颜色
  primaryColor: "#2563eb"
  primaryDark: "#1e40af"
  accentColor: "#60a5fa"
  bgPrimary: "#f8fafc"
  bgSecondary: "#ffffff"
  textPrimary: "#1e293b"
  textSecondary: "#64748b"
  borderColor: "#e5e7eb"

  # 功能开关
  heroEnabled: true
  sidebarEnabled: true
  showViews: true
  showReadMore: true

  # Hero 配置
  heroTitle: "技术与生活"
  heroSubtitle: "记录所思所想，分享技术干货"
  heroButtonText: "开始阅读"
  heroButtonLink: "/article"

  # 文章列表
  postListStyle: "card"   # card / list / magazine

  # 侧边栏
  sidebarWidgets: '["categories","tags"]'

  # 博主信息
  authorName: ""
  authorBio: ""
  authorAvatar: ""

  # 页脚
  footerText: ""
  footerLinks: '[{"name":"关于","url":"/page/about"},{"name":"归档","url":"/archives"}]'

configDefinition:
  - section: style
    label: 外观风格
    icon: fa-solid fa-palette
    fields:
      - name: primaryColor
        type: color
        label: 主色调
        default: "#2563eb"
        description: 导航栏、按钮、链接等主要颜色
      - name: primaryDark
        type: color
        label: 深主色
        default: "#1e40af"
      - name: accentColor
        type: color
        label: 强调色
        default: "#60a5fa"
      - name: bgPrimary
        type: color
        label: 背景色
        default: "#f8fafc"
      - name: textPrimary
        type: color
        label: 主文字色
        default: "#1e293b"

  - section: hero
    label: Hero 区域
    icon: fa-solid fa-star
    fields:
      - name: heroEnabled
        type: switch
        label: 显示 Hero 区域
        default: true
      - name: heroTitle
        type: text
        label: 主标题
        show_if: "heroEnabled == true"
        default: 技术与生活
      - name: heroSubtitle
        type: textarea
        label: 副标题
        show_if: "heroEnabled == true"
        default: 记录所思所想
      - name: heroButtonText
        type: text
        label: 按钮文字
        show_if: "heroEnabled == true"
        default: 开始阅读
      - name: heroButtonLink
        type: text
        label: 按钮链接
        show_if: "heroEnabled == true"
        default: /article

  - section: layout
    label: 布局设置
    icon: fa-solid fa-table-columns
    fields:
      - name: sidebarEnabled
        type: switch
        label: 显示侧边栏
        default: true
      - name: sidebarWidgets
        type: checkbox
        label: 侧边栏模块
        show_if: "sidebarEnabled == true"
        options:
          - value: categories
            label: 分类列表
          - value: tags
            label: 标签云
          - value: profile
            label: 博主信息
      - name: postListStyle
        type: radio
        label: 文章列表样式
        default: card
        options:
          - value: card
            label: 卡片式
          - value: list
            label: 列表式
          - value: magazine
            label: 杂志式

  - section: author
    label: 博主信息
    icon: fa-solid fa-user
    fields:
      - name: authorName
        type: text
        label: 博主昵称
        default: ""
      - name: authorBio
        type: textarea
        label: 博主简介
        default: ""
      - name: authorAvatar
        type: image
        label: 博主头像
        default: ""

  - section: display
    label: 显示设置
    icon: fa-solid fa-eye
    fields:
      - name: showViews
        type: switch
        label: 显示阅读量
        default: true

  - section: footer
    label: 页脚设置
    icon: fa-solid fa-bars
    fields:
      - name: footerText
        type: text
        label: 版权附加文字
        default: ""
        placeholder: 如：All rights reserved.
      - name: footerLinks
        type: repeater
        label: 页脚链接
        default: '[]'
        fields:
          - name: name
            type: text
            label: 链接名称
          - name: url
            type: text
            label: 链接地址

  - section: seo
    label: SEO 设置
    icon: fa-solid fa-magnifying-glass
    fields:
      - name: metaKeywords
        type: text
        label: 全局关键词
        default: ""
      - name: gaTrackingId
        type: text
        label: Google Analytics ID
        default: ""

  - section: custom
    label: 自定义代码
    icon: fa-solid fa-code
    fields:
      - name: customCss
        type: code
        label: 自定义 CSS
        default: ""
      - name: customJs
        type: code
        label: 自定义 JS
        default: ""
```

---

## 示例二：企业官网主题

```yaml
name: Enterprise Pro
slug: enterprise-pro
version: 1.0.0
author: 开发者名称
description: 专业商务企业官网主题，模块化设计，适合中大型企业展示

config:
  # 颜色
  primaryColor: "#0f3460"
  primaryDark: "#0a2444"
  accentColor: "#e94560"
  bgPrimary: "#f8f9fa"
  bgSecondary: "#ffffff"
  textPrimary: "#1a1a2e"
  textSecondary: "#6b7280"
  borderColor: "#e5e7eb"

  # Hero 区域
  heroEnabled: true
  heroTitle: "专业赋能，助力企业腾飞"
  heroSubtitle: "我们提供全方位的企业服务解决方案"
  heroBgType: "gradient"  # gradient / image / video
  heroBgImage: ""

  # 模块开关
  moduleStatsEnabled: true
  moduleFeaturesEnabled: true
  moduleServicesEnabled: true
  moduleNewsEnabled: true
  moduleCTAEnabled: true

  # Stats 模块
  statsItems: '[{"value":"500+","label":"服务客户"},{"value":"10年","label":"行业经验"},{"value":"98%","label":"客户满意度"}]'

  # Features 模块
  featuresTitle: "我们的优势"
  featuresSubtitle: "专注于为客户提供最优质的服务"
  featuresItems: '[{"icon":"fa-solid fa-bolt","title":"高效快速","desc":"快速响应，高效交付"},{"icon":"fa-solid fa-shield-halved","title":"安全可靠","desc":"多重安全保障"},{"icon":"fa-solid fa-headset","title":"专业支持","desc":"7×24 专业服务"}]'

  # CTA 模块
  ctaTitle: "准备好开始了吗？"
  ctaSubtitle: "联系我们，获取专属解决方案"
  ctaButtonText: "立即咨询"
  ctaButtonLink: "/page/contact"

configDefinition:
  - section: branding
    label: 品牌设置
    icon: fa-solid fa-building
    fields:
      - name: primaryColor
        type: color
        label: 品牌主色
        default: "#0f3460"
      - name: accentColor
        type: color
        label: 强调色
        default: "#e94560"
      - name: bgPrimary
        type: color
        label: 背景色
        default: "#f8f9fa"

  - section: hero
    label: Hero 区域
    icon: fa-solid fa-image
    fields:
      - name: heroEnabled
        type: switch
        label: 显示 Hero
        default: true
      - name: heroTitle
        type: text
        label: 主标题
        show_if: "heroEnabled == true"
      - name: heroSubtitle
        type: textarea
        label: 副标题
        show_if: "heroEnabled == true"
      - name: heroBgType
        type: radio
        label: 背景类型
        show_if: "heroEnabled == true"
        default: gradient
        options:
          - value: gradient
            label: 渐变色
          - value: image
            label: 背景图
      - name: heroBgImage
        type: image
        label: 背景图片
        show_if: "heroEnabled == true && heroBgType == image"

  - section: modules
    label: 模块管理
    icon: fa-solid fa-puzzle-piece
    fields:
      - name: moduleStatsEnabled
        type: switch
        label: 数据统计模块
        default: true
      - name: moduleFeaturesEnabled
        type: switch
        label: 特性优势模块
        default: true
      - name: moduleServicesEnabled
        type: switch
        label: 服务项目模块
        default: true
      - name: moduleNewsEnabled
        type: switch
        label: 新闻动态模块
        default: true
      - name: moduleCTAEnabled
        type: switch
        label: 行动号召模块
        default: true

  - section: stats
    label: 数据统计
    icon: fa-solid fa-chart-bar
    fields:
      - name: statsItems
        type: repeater
        label: 统计项目
        show_if: "moduleStatsEnabled == true"
        fields:
          - name: value
            type: text
            label: 数值（如 500+）
          - name: label
            type: text
            label: 描述文字

  - section: features
    label: 特性优势
    icon: fa-solid fa-star
    fields:
      - name: featuresTitle
        type: text
        label: 模块标题
        show_if: "moduleFeaturesEnabled == true"
        default: 我们的优势
      - name: featuresItems
        type: repeater
        label: 特性列表
        show_if: "moduleFeaturesEnabled == true"
        fields:
          - name: icon
            type: text
            label: 图标（Font Awesome 类名）
            placeholder: fa-solid fa-bolt
          - name: title
            type: text
            label: 标题
          - name: desc
            type: textarea
            label: 描述

  - section: cta
    label: 行动号召
    icon: fa-solid fa-bullhorn
    fields:
      - name: ctaTitle
        type: text
        label: 标题
        show_if: "moduleCTAEnabled == true"
        default: 准备好开始了吗？
      - name: ctaSubtitle
        type: text
        label: 副标题
        show_if: "moduleCTAEnabled == true"
      - name: ctaButtonText
        type: text
        label: 按钮文字
        show_if: "moduleCTAEnabled == true"
        default: 立即咨询
      - name: ctaButtonLink
        type: text
        label: 按钮链接
        show_if: "moduleCTAEnabled == true"
        default: /page/contact

  - section: custom
    label: 自定义代码
    icon: fa-solid fa-code
    fields:
      - name: customCss
        type: code
        label: 自定义 CSS
        default: ""
```

---

## 示例三：新中式博客主题

> 使用深色调、宋体字体栈、传统美学的博客主题

```yaml
name: 墨韵
slug: moyun-blog
version: 1.0.0
author: 开发者名称
description: 新中式美学博客主题，墨色书写，诗意生活

config:
  primaryColor: "#2D4A3E"
  primaryDark: "#1a2e26"
  accentColor: "#C9A84C"
  bgPrimary: "#FAF7F0"
  bgSecondary: "#FFFFFF"
  textPrimary: "#2C2C2C"
  textSecondary: "#666666"
  borderColor: "#E8E0D0"

  heroEnabled: true
  heroTitle: "墨韵"
  heroSubtitle: "以文会友，以字传情"

  sidebarEnabled: true

  # 特色配置：作者信息
  authorName: ""
  authorPoem: "日日行，不怕千万里；常常做，不怕千万事。"
  authorAvatar: ""

  footerIcp: ""
  footerText: ""

configDefinition:
  - section: colors
    label: 颜色设置
    icon: fa-solid fa-palette
    fields:
      - name: primaryColor
        type: color
        label: 主色（墨绿）
        default: "#2D4A3E"
      - name: accentColor
        type: color
        label: 强调色（金色）
        default: "#C9A84C"
      - name: bgPrimary
        type: color
        label: 背景色（宣纸）
        default: "#FAF7F0"

  - section: hero
    label: 首页 Banner
    icon: fa-solid fa-image
    fields:
      - name: heroEnabled
        type: switch
        label: 显示 Banner
        default: true
      - name: heroTitle
        type: text
        label: 站点名称（大字）
        show_if: "heroEnabled == true"
        default: 墨韵
      - name: heroSubtitle
        type: textarea
        label: 站点格言
        show_if: "heroEnabled == true"
        default: 以文会友，以字传情

  - section: author
    label: 作者简介
    icon: fa-solid fa-user-pen
    fields:
      - name: authorName
        type: text
        label: 笔名
        default: ""
      - name: authorPoem
        type: textarea
        label: 座右铭
        default: ""
      - name: authorAvatar
        type: image
        label: 头像
        default: ""

  - section: footer
    label: 页脚
    icon: fa-solid fa-bars
    fields:
      - name: footerIcp
        type: text
        label: ICP 备案号
        default: ""
        placeholder: 如：京ICP备XXXXXXXX号
      - name: footerText
        type: text
        label: 附加信息
        default: ""

  - section: custom
    label: 自定义代码
    icon: fa-solid fa-code
    fields:
      - name: customCss
        type: code
        label: 自定义 CSS
        default: ""
```
