# 主题升级计划 3.0

> 基于 UI/UX Pro Max 技能标准，全面优化 Catter CMS 主题系统
> 创建时间：2026-03-24
> 状态：进行中

---

## 一、优化清单与进度

### P0 紧急修复（影响核心功能）

| 序号 | 优化点 | 状态 | 负责人 | 完成时间 | 备注 |
|------|--------|------|--------|----------|------|
| P0-1 | 修复主题预览功能 | ✅ 已完成 | AI | 2026-03-24 | ThemeList.vue:180 |
| P0-2 | 修复重置功能逻辑 | ✅ 已完成 | AI | 2026-03-24 | 恢复默认值非上次保存 |
| P0-3 | 实现 repeater 可视化编辑器 | ✅ 已完成 | AI | 2026-03-24 | 替换 JSON textarea |

### P1 重要优化（提升体验）

| 序号 | 优化点 | 状态 | 负责人 | 完成时间 | 备注 |
|------|--------|------|--------|----------|------|
| P1-1 | 实现 DeleteTheme 逻辑 | ✅ 已完成 | AI | 2026-03-24 | 删除文件系统+数据库 |
| P1-2 | 实现配置导入/导出功能 | ✅ 已完成 | AI | 2026-03-24 | 前端+后端 |
| P1-3 | 主题截图支持 | ✅ 已完成 | AI | 2026-03-24 | 自动检测 screenshot.png |

### P2 体验优化（持续改进）

| 序号 | 优化点 | 状态 | 负责人 | 完成时间 | 备注 |
|------|--------|------|--------|----------|------|
| P2-1 | 统一 CSS 管理方式 | ✅ 已完成 | AI | 2026-03-24 | 创建 assets/css/ 目录，提取内联 CSS |
| P2-2 | 添加快捷键支持 | ✅ 已完成 | AI | 2026-03-24 | Ctrl+S/Ctrl+R |
| P2-3 | 完善可访问性 | ✅ 已完成 | AI | 2026-03-24 | ARIA/焦点管理 |

---

## 二、现状分析

### 2.1 主题数量与状态

| 主题 | 模板文件 | 配置文件 | 问题数 |
|------|----------|----------|--------|
| `catter-blog` | 6个页面 | theme.yml (113行) | 0 |
| `catter-cms-official` | 11个页面 | theme.yml (264行) | 6个 JSON textarea |

### 2.2 核心架构

```
配置层：theme.yml (configDefinition) ──→ 后端 services/theme.go ──→ 前端 ThemeSettings.vue
                                    ↑
                          models/theme_config.go (数据库)
```

### 2.3 发现的问题详情

#### P0 问题

**P0-1: 主题预览功能错误**
- 位置：`admin/src/views/themes/ThemeList.vue:180`
- 问题：`window.open('/', '_blank')` 始终打开首页
- 期望：打开该主题的预览页面

**P0-2: 重置功能逻辑错误**
- 位置：`admin/src/views/themes/ThemeSettings.vue:295`
- 问题：重置为"上次保存"而非"默认值"
- 代码：
```javascript
// 当前实现
configValues.value = JSON.parse(JSON.stringify(originalConfig.value))
// originalConfig 是上次保存的值，不是 YAML 默认值
```

**P0-3: JSON textarea 反人类设计**
- 位置：`templates/themes/catter-cms-official/theme.yml`
- 问题：features、stats、pricing、testimonials、footerLinks、footerSocial 使用 textarea 存储 JSON
- 影响：用户体验极差，容易出错

#### P1 问题

**P1-1: DeleteTheme 未实现**
- 位置：`gin-app/internal/handlers/api/themes.go:261`
- 问题：TODO 注释，只是返回成功

**P1-2: 配置导入/导出未实现**
- 位置：`admin/src/api/theme.js` 声明但未实现

**P1-3: 主题截图支持缺失**
- 位置：ThemeList.vue 始终显示 placeholder

---

## 三、UI/UX Pro Max 设计标准

### 3.1 视觉设计

| 检查项 | 标准 | 当前状态 |
|--------|------|----------|
| 配色系统 | 统一的 design token | 各主题独立 |
| 字体规范 | 主题级字体设置 | 无统一规范 |
| 间距系统 | 8px 网格系统 | 硬编码 px |
| 圆角/阴影 | 统一设计系统 | 各主题自定义 |

### 3.2 交互设计

| 检查项 | 标准 | 当前状态 |
|--------|------|----------|
| 表单验证 | 实时校验反馈 | 基础 |
| 加载状态 | 骨架屏+渐进式 | 基础 |
| 错误处理 | 上下文错误提示 | 基础 message |
| 快捷键 | 键盘导航支持 | 无 |

### 3.3 可访问性

| 检查项 | 标准 | 当前状态 |
|--------|------|----------|
| ARIA 标签 | 表单字段关联 | 部分 |
| 焦点管理 | 明确的焦点样式 | 基础 |
| 颜色对比 | WCAG AA 合规 | 部分合规 |
| 屏幕阅读器 | aria-live 区域 | 未优化 |

---

## 四、具体实施方案

### P0-1: 修复主题预览功能

**修改文件**：`admin/src/views/themes/ThemeList.vue`

**修改内容**：
```javascript
// 第179-181行
function handlePreview(theme) {
  window.open('/?theme=' + theme.slug, '_blank')
}
```

**验收标准**：
- [ ] 点击预览按钮，新窗口打开 `/?theme={slug}`
- [ ] URL 参数正确传递

---

### P0-2: 修复重置功能逻辑

**修改文件**：`admin/src/views/themes/ThemeSettings.vue`

**修改内容**：
1. 添加 `defaultValues` 状态存储 YAML 默认值
2. 修改 `resetConfig` 函数使用默认值
3. 修改确认对话框文案

**验收标准**：
- [ ] 重置按钮点击后，配置恢复到 YAML 中定义的 `default` 值
- [ ] 对话框文案更新为"这将恢复到默认值"

---

### P0-3: 实现 repeater 可视化编辑器

**修改文件**：
1. `templates/themes/catter-cms-official/theme.yml`
2. `admin/src/components/RepeaterField.vue` (如需增强)

**修改内容**：将以下字段从 `type: textarea` 改为 `type: repeater`

| 字段名 | children 定义 |
|--------|---------------|
| features | icon (text), title (text), desc (textarea) |
| stats | value (text), label (text) |
| products | name (text), desc (textarea), link (text), image (image) |
| pricing | name (text), price (text), period (text), features (multi-select), button (text), link (text), recommended (switch) |
| testimonials | name (text), title (text), content (textarea), avatar (image) |
| footerLinks | text (text), link (text) |
| footerSocial | icon (select), link (text) |

**验收标准**：
- [ ] 以上字段在 ThemeSettings 中显示为可视化 repeater 编辑器
- [ ] 支持拖拽排序
- [ ] 支持添加/删除
- [ ] 数据正确序列化为 JSON 存储

---

### P1-1: 实现 DeleteTheme 逻辑

**修改文件**：`gin-app/internal/handlers/api/themes.go`

**修改内容**：
```go
// 第261行附近
// TODO: 实现从文件系统删除主题的逻辑
// 改为：
// 1. 获取主题目录路径
// 2. 删除主题目录
// 3. 删除数据库中的主题配置
// 4. 重新加载主题服务
```

**验收标准**：
- [ ] 删除非激活主题时，文件系统目录被删除
- [ ] 数据库中的主题配置被清除
- [ ] 无法删除当前激活的主题

---

### P1-2: 实现配置导入/导出功能

**修改文件**：
1. `gin-app/internal/handlers/api/themes.go`
2. `admin/src/api/theme.js`
3. `admin/src/views/themes/ThemeSettings.vue`

**修改内容**：
- 后端：添加 `ExportThemeConfig` 和 `ImportThemeConfig` Handler
- 前端：在 ThemeSettings.vue 添加导入/导出按钮

**验收标准**：
- [ ] 可以导出当前主题配置为 JSON 文件
- [ ] 可以导入 JSON 文件恢复配置
- [ ] 导入前有确认提示

---

### P1-3: 主题截图支持

**修改文件**：
1. `admin/src/views/themes/ThemeList.vue`
2. `templates/themes/catter-blog/theme.yml`
3. `templates/themes/catter-cms-official/theme.yml`

**修改内容**：
- 在 theme.yml 添加 `thumbnail` 字段指向截图路径
- 或在主题目录添加 `screenshot.png`

**验收标准**：
- [ ] 主题列表显示真实缩略图
- [ ] 无截图时显示优雅的 placeholder

---

### P2-1: 统一 CSS 管理方式

**目标**：所有主题使用主题目录内的 CSS 文件，通过 `asset_path` 引用

**验收标准**：
- [ ] 统一使用 `assets/css/` 目录
- [ ] 不使用内联样式

---

### P2-2: 添加快捷键支持

**修改文件**：`admin/src/views/themes/ThemeSettings.vue`

**快捷键定义**：
| 快捷键 | 功能 |
|--------|------|
| Ctrl+S | 保存配置 |
| Ctrl+R | 重置配置 |
| Esc | 关闭对话框 |

**验收标准**：
- [ ] 快捷键正常工作
- [ ] 有键盘焦点时快捷键可用

---

### P2-3: 完善可访问性

**修改文件**：多个 Vue 组件

**验收标准**：
- [ ] 表单字段有正确的 `aria-describedby` 关联
- [ ] 焦点状态清晰可见
- [ ] 颜色对比度符合 WCAG AA

---

## 五、审查清单

### 功能审查

- [ ] 主题列表显示正确
- [ ] 主题激活功能正常
- [ ] 主题配置保存/加载正常
- [ ] 主题预览功能正常
- [ ] 主题删除功能正常
- [ ] Repeater 字段可视化编辑正常
- [ ] 配置导入/导出正常

### UI/UX 审查

- [ ] 视觉风格一致
- [ ] 交互反馈及时
- [ ] 加载状态优雅
- [ ] 错误处理友好
- [ ] 快捷键工作正常

### 可访问性审查

- [ ] 键盘导航正常
- [ ] ARIA 标签完整
- [ ] 颜色对比合规
- [ ] 焦点管理正确

---

## 六、版本历史

| 版本 | 日期 | 修改内容 |
|------|------|----------|
| 1.0.0 | 2026-03-24 | 初始文档 |
| 1.0.1 | 2026-03-24 | P0-1、P0-2、P0-3、P1-1、P1-2、P1-3、P2-2、P2-3 完成 |
| 1.0.2 | 2026-03-24 | P2-1 CSS 统一重构完成 |
| 1.0.3 | 2026-03-24 | 全部优化完成 (9/9) |

---

## 七、完成总结

### 已完成优化 (9/9)

| 优先级 | 序号 | 优化点 | 修改文件 |
|--------|------|--------|----------|
| P0 | P0-1 | 修复主题预览功能 | admin/src/views/themes/ThemeList.vue |
| P0 | P0-2 | 修复重置功能逻辑 | admin/src/views/themes/ThemeSettings.vue |
| P0 | P0-3 | 实现 repeater 可视化编辑器 | templates/themes/catter-cms-official/theme.yml |
| P1 | P1-1 | 实现 DeleteTheme 逻辑 | gin-app/internal/handlers/api/themes.go, services/theme.go |
| P1 | P1-2 | 实现配置导入/导出功能 | gin-app/internal/handlers/api/themes.go, admin/src/views/themes/ThemeSettings.vue, admin/src/api/theme.js |
| P1 | P1-3 | 主题截图支持 | gin-app/internal/services/theme.go |
| P2 | P2-1 | 统一 CSS 管理方式 | templates/themes/catter-cms-official/assets/css/ |
| P2 | P2-2 | 添加快捷键支持 | admin/src/views/themes/ThemeSettings.vue |
| P2 | P2-3 | 完善可访问性 | admin/src/views/themes/ThemeSettings.vue |

### 待完成优化 (0/9)

| 优先级 | 序号 | 优化点 | 说明 |
|--------|------|--------|------|
| - | - | 全部完成 | - |

### 关键改进点

1. **repeater 可视化编辑器**：将 catter-cms-official 主题的 6 个 JSON textarea 字段改为可视化 repeater 编辑器
2. **DeleteTheme 功能**：完整实现从文件系统和数据库删除主题
3. **配置导入/导出**：支持配置数据的导出为 JSON 和从 JSON 导入
4. **主题截图**：自动检测主题目录下的 screenshot.png 等图片文件
5. **快捷键**：Ctrl+S 保存、Ctrl+R 重置、Esc 关闭对话框
6. **可访问性**：添加 ARIA 属性、焦点管理、键盘导航支持
7. **统一 CSS 管理**：创建 `assets/css/` 目录，提取内联 CSS 为独立文件（base.css、home.css、article.css、list.css）

---

*最后更新：2026-03-24 21:00*
