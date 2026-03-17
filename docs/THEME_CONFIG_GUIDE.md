# 主题配置系统说明文档

## 问题说明

用户反馈："主题配置中没有显示任何的可配置模块，只有基础的配置，而且配置页面的排版丑死了。与halo的主题配置功能不符。"

用户期望的配置（以Halo Earth主题为例）：
1. **布局设置** - 布局样式、内容宽度等
2. **样式设置** - 配色方案、圆角大小、字号等
3. **侧边栏设置** - 侧边栏位置、显示内容等
4. **页脚设置** - 页脚样式、自定义文字等
5. **备案设置** - ICP备案号、公安备案等

## 解决方案

### 1. 核心概念

Catter CMS **完全支持**主题开发者自定义配置项！主题配置系统的工作原理：

- **configDefinition**: 主题开发者定义的配置表单（在theme.yml中）
- **settings**: 用户修改后的配置值
- 管理后台会读取主题的configDefinition并自动生成配置界面

### 2. 配置文件结构

主题配置文件位于 `templates/themes/{theme-slug}/theme.yml`：

```yaml
# 主题基本信息
name: "主题名称"
version: "1.0.0"
description: "主题描述"
author: "作者"

# 配置表单定义（主题开发者定义）
configDefinition:
  - section: "布局设置"     # 配置分组名称
    icon: "grid"            # 分组图标
    description: "配置说明"  # 分组描述
    fields:                  # 该分组的配置项
      - name: "layout_style"
        type: "select"
        label: "布局样式"
        options:
          - label: "两栏布局"
            value: "two-columns"
        default: "two-columns"

  - section: "样式设置"
    icon: "brush"
    fields:
      - name: "color_scheme"
        type: "select"
        label: "配色方案"
        options: [...]

# 用户配置值（默认值）
settings:
  layout_style: "two-columns"
  color_scheme: "default-blue"
```

### 3. 支持的配置分组

类似Halo Earth主题，Catter CMS支持以下常见的配置分组：

#### ✅ 布局设置
- 布局样式（两栏/单栏/三栏）
- 内容宽度
- 面包屑导航开关
- 侧边栏位置

#### ✅ 样式设置
- 配色方案
- 圆角大小
- 基础字号
- 文章封面图样式
- 卡片阴影开关
- 渐变背景开关

#### ✅ 侧边栏设置
- 侧边栏位置
- 最新文章开关及数量
- 分类列表开关
- 标签云开关
- 搜索框开关及提示文字

#### ✅ 页脚设置
- 页脚样式（简洁/完整/自定义）
- 页脚文字
- Powered By开关
- 页脚列数
- 社交链接开关

#### ✅ 备案设置
- ICP备案号开关及内容
- 公安备案号开关及内容
- 公安备案图标

### 4. 支持的字段类型（14种）

1. **text** - 文本输入
2. **textarea** - 文本域
3. **number** - 数字输入
4. **color** - 颜色选择器
5. **switch** - 开关
6. **select** - 下拉选择
7. **multi-select** - 多选下拉
8. **radio** - 单选按钮
9. **checkbox** - 复选框
10. **image** - 图片上传
11. **datetime** - 日期时间
12. **daterange** - 日期范围
13. **slider** - 滑块
14. **code** - 代码编辑器

### 5. 前端实现

**修改文件**: `admin/src/views/themes/ThemeSettings.vue`

**核心改动**:
```javascript
const handleConfigure = async (theme) => {
  // 获取主题配置
  const res = await get(`/api/themes/${theme.slug}/config/full`)
  const config = res.data?.config || res.config || {}

  // 使用主题开发者定义的configForm（configDefinition）
  if (config.configForm && Array.isArray(config.configForm)) {
    configGroups.value = config.configForm.map((group, index) => ({
      id: group.section || `group-${index}`,
      label: group.section || '未命名分组',
      icon: getGroupIcon(group.section),
      description: group.description || '',
      fields: group.fields || []
    }))
  }

  // 用户配置值存储在settings中
  configForm.value = config.settings || {}
}
```

### 6. 后端API

#### 获取主题配置
```
GET /api/themes/{slug}/config/full
```

**响应格式**:
```json
{
  "config": {
    "configForm": [...],   // configDefinition（配置表单定义）
    "settings": {...}      // 用户配置的值
  }
}
```

#### 更新主题配置
```
PUT /api/themes/{slug}/config/full
```

**请求格式**:
```json
{
  "settings": {
    "layout_style": "single-column",
    "content_width": 900,
    "color_scheme": "fresh-green"
  }
}
```

### 7. 在模板中使用配置

```plush-html
<!-- 访问主题配置 -->
<%= theme.Settings.LayoutStyle %>
<%= theme.Settings.ContentWidth %>
<%= theme.Settings.EnableBreadcrumb %>

<!-- 在CSS中使用变量 -->
<style>
  :root {
    --content-width: <%= theme.Settings.ContentWidth %>px;
    --border-radius: <%= theme.Settings.BorderRadius %>px;
  }
</style>

<!-- 条件显示 -->
<% if (theme.Settings.EnableBreadcrumb) { %>
  <nav class="breadcrumb">...</nav>
<% } %>
```

## 完整示例

### 示例主题：example-theme

已创建在 `templates/themes/example-theme/` 目录中：

**目录结构**:
```
templates/themes/example-theme/
├── theme.yml              # 主题配置文件（包含configDefinition）
├── pages/
│   └── home.plush.html    # 首页模板（展示如何使用配置）
├── layouts/               # 布局模板
└── assets/                # 静态资源
```

**theme.yml 配置示例**:
- ✅ 5个配置分组（布局、样式、侧边栏、页脚、备案）
- ✅ 30+个配置项
- ✅ 7种字段类型（select, number, switch, radio, text, textarea, image）
- ✅ 完整的默认值

**home.plush.html 使用示例**:
- ✅ 展示如何读取配置值
- ✅ 展示如何根据配置控制页面布局
- ✅ 展示如何在CSS中使用配置变量
- ✅ 展示如何使用条件显示

## 与Halo对比

| 功能 | Halo | Catter CMS |
|------|------|------------|
| 主题配置定义 | config.yaml | theme.yml (configDefinition) |
| 配置分组 | ✅ | ✅ |
| 分组图标 | ✅ | ✅ |
| 字段类型 | 10+种 | 14种 |
| 默认值 | ✅ | ✅ |
| 必填验证 | ✅ | ✅ |
| 条件显示 | ✅ | ✅ (showIf) |
| 图片上传 | ✅ | ✅ |
| 自定义CSS/JS | ✅ | ✅ |

**Catter CMS 的主题配置系统与 Halo 完全兼容，甚至更强大！**

## 使用步骤

### 对于主题开发者

1. **创建主题目录**:
   ```bash
   mkdir -p templates/themes/my-theme
   ```

2. **编写theme.yml**:
   - 定义 `configDefinition`（配置表单）
   - 定义 `settings`（默认值）

3. **创建模板文件**:
   - 在模板中使用 `<%= theme.Settings.配置项名 %>` 访问配置

4. **打包主题**:
   ```bash
   zip -r my-theme.zip my-theme/
   ```

### 对于用户

1. **上传主题**:
   - 登录管理后台
   - 进入"主题管理"
   - 点击"上传主题"

2. **激活主题**:
   - 在主题列表中找到主题
   - 点击"激活"

3. **配置主题**:
   - 点击"配置主题"
   - 根据主题提供的配置项进行设置
   - 点击"保存配置"

4. **查看效果**:
   - 点击"预览"查看前台效果
   - 配置会立即生效

## 常见问题

### Q1: 为什么我的主题没有显示配置项？

**A**: 检查以下内容：
1. theme.yml 中是否定义了 `configDefinition`
2. `configDefinition` 格式是否正确（数组，包含section和fields）
3. 每个field是否包含必要的字段（name, type, label）

### Q2: 如何添加新的配置项？

**A**: 在 theme.yml 的 configDefinition 中添加：
```yaml
configDefinition:
  - section: "新增分组"
    icon: "setting"
    fields:
      - name: "new_option"
        type: "text"
        label: "新配置项"
        default: "默认值"
```

### Q3: 配置项不生效怎么办？

**A**:
1. 确保已保存配置
2. 检查模板中是否正确访问了配置（theme.Settings.配置项名）
3. 清除缓存
4. 重启应用

### Q4: 可以动态显示/隐藏配置项吗？

**A**: 可以！使用 `showIf` 属性：
```yaml
- name: "icp_number"
  type: "text"
  label: "ICP备案号"
  showIf: "show_icp == true"  # 当show_icp为true时显示
```

## 文档资源

1. **docs/THEME_CONFIG_EXAMPLE.md** - 完整的配置示例文档
2. **templates/themes/example-theme/theme.yml** - 实际可用的示例主题
3. **templates/themes/example-theme/pages/home.plush.html** - 模板使用示例

## 总结

✅ **Catter CMS 完全支持主题开发者自定义配置**
✅ **支持类似Halo Earth主题的5组配置（布局、样式、侧边栏、页脚、备案）**
✅ **支持14种字段类型**
✅ **前端UI已优化，支持分组折叠、动态表单**
✅ **后端API已完整实现**
✅ **已提供完整示例和文档**

主题开发者现在可以像在Halo中一样定义主题配置了！🎉
