# Catter CMS 原子级开发计划

> 版本：1.0 | 日期：2026-03-20 | 状态：待执行

---

## 总览

本计划基于《主题系统升级方案》，分为 4 个阶段，共 15 个原子任务。

---

## Phase 1: API 密钥系统（4天）

### 任务 1.1: API 密钥模型
**文件**: `gin-app/internal/models/api_key.go`

**创建内容**:
```go
type APIKey struct {
    ID          uint
    Name        string      // 密钥名称
    Key         string      // sk-xxxxx 格式
    Permissions string      // JSON 权限列表
    Enabled     bool        // 是否启用
    UsageCount  int         // 使用次数
    LastUsedAt  *time.Time  // 最后使用时间
    ExpiresAt   *time.Time  // 过期时间
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 方法
- GenerateAPIKey() 生成随机 Key
- GetAPIByKey(key) 通过 Key 查询
- HasPermission(resource, action) 检查权限
- IncrementUsage() 增加使用次数
```

**验收标准**:
- [ ] 文件创建成功
- [ ] 结构体定义完整
- [ ] 方法实现正确

---

### 任务 1.2: API 密钥 CRUD API
**文件**: `gin-app/internal/handlers/api/api_keys.go`

**创建内容**:
```
GET    /api/api-keys          # 获取所有密钥
POST   /api/api-keys          # 创建新密钥
PUT    /api/api-keys/:id      # 更新密钥
DELETE /api/api-keys/:id      # 删除密钥
POST   /api/api-keys/:id/reset # 重置密钥
```

**响应格式**:
```json
{
  "id": 1,
  "name": "AI助手专用",
  "key": "sk-a1b2c3...",
  "permissions": ["posts", "pages"],
  "enabled": true,
  "usage_count": 42,
  "created_at": "2026-03-20T12:00:00Z"
}
```

**验收标准**:
- [ ] 5 个端点实现
- [ ] 返回格式正确
- [ ] 错误处理完整

---

### 任务 1.3: API Key 认证中间件
**文件**: `gin-app/internal/middleware/apikey.go`

**创建内容**:
```go
func APIKeyAuth() gin.HandlerFunc {
    // 1. 从 Header 获取 Authorization: Bearer sk-xxxxx
    // 2. 验证 Key 是否存在且启用
    // 3. 检查是否过期
    // 4. 检查权限
    // 5. 增加使用统计
    // 6. 设置到 context
}
```

**权限检查逻辑**:
```go
func CheckPermission(apiKey *APIKey, method string, path string) bool {
    // GET    /posts → posts.read
    // POST   /posts → posts.create
    // PUT    /posts/:id → posts.update
    // DELETE /posts/:id → posts.delete
}
```

**验收标准**:
- [ ] 中间件实现
- [ ] 权限检查正确
- [ ] 错误响应正确

---

### 任务 1.4: 管理后台 API 密钥界面
**文件**: `admin/src/views/api-keys/Index.vue`

**UI 包含**:
- 密钥列表显示
- 创建新密钥对话框
- 复制密钥功能
- 禁用/启用切换
- 重置密钥功能
- 删除确认

**验收标准**:
- [ ] 列表显示正确
- [ ] 创建功能正常
- [ ] 复制功能正常
- [ ] 禁用/启用正常
- [ ] 重置功能正常
- [ ] 删除功能正常

---

## Phase 2: AI 文档接口（3天）

### 任务 2.1: 能力发现接口
**文件**: `gin-app/internal/handlers/api/ai_docs.go`

**端点**: `GET /api/v1/ai/capabilities`

**返回内容**:
```json
{
  "name": "Catter CMS",
  "version": "1.0.0",
  "description": "内容管理系统",
  "auth": {
    "type": "bearer",
    "header": "Authorization",
    "format": "Bearer {api_key}"
  },
  "capabilities": [
    {
      "id": "posts",
      "name": "文章管理",
      "description": "创建、编辑、删除、查询文章",
      "endpoints": [
        "GET /posts",
        "GET /posts/:id",
        "POST /posts",
        "PUT /posts/:id",
        "DELETE /posts/:id"
      ]
    }
  ]
}
```

**验收标准**:
- [ ] 端点实现
- [ ] 返回格式正确
- [ ] 包含所有资源

---

### 任务 2.2: 操作指南接口
**文件**: `gin-app/internal/handlers/api/ai_docs.go`

**端点**: `GET /api/v1/ai/guide`

**返回内容**:
```json
{
  "quick_start": {
    "step_1": "获取 API Key",
    "step_2": "设置环境变量",
    "step_3": "使用 Bearer token 认证"
  },
  "common_tasks": [
    {
      "task": "创建文章",
      "method": "POST",
      "path": "/api/v1/posts",
      "required": ["title", "content"],
      "optional": ["status", "category_id", "tag_ids"],
      "example": {
        "title": "文章标题",
        "content": "# 内容...",
        "status": "published"
      }
    }
  ],
  "restrictions": [
    "不能删除管理员用户",
    "批量操作限制 100 条/次"
  ]
}
```

**验收标准**:
- [ ] 端点实现
- [ ] 包含所有常见任务
- [ ] 示例代码正确

---

### 任务 2.3: 详细文档接口
**文件**: `gin-app/internal/handlers/api/ai_docs.go`

**端点**:
```
GET /api/v1/ai/docs              # 所有端点文档
GET /api/v1/ai/docs/:operation   # 特定操作文档
```

**验收标准**:
- [ ] 端点实现
- [ ] 文档完整准确
- [ ] 示例可执行

---

### 任务 2.4: OpenAPI 规范接口
**文件**: `gin-app/internal/handlers/api/ai_docs.go`

**端点**: `GET /api/v1/ai/openapi.json`

**验收标准**:
- [ ] 端点实现
- [ ] OpenAPI 3.0 格式正确
- [ ] 可被 AI 工具解析

---

## Phase 3: 页面模板系统（3天）

### 任务 3.1: 模板配置定义
**文件**: `templates/themes/*/theme.yaml`

**添加内容**:
```yaml
customTemplates:
  page:
    - name: "默认页面"
      file: "page_default.html"
      description: "标准内容页布局"
    - name: "服务列表"
      file: "page_services.html"
      description: "带图片卡片的服务展示页"
```

**验收标准**:
- [ ] 配置格式正确
- [ ] 模板文件存在

---

### 任务 3.2: 页面模板选择器
**文件**: `gin-app/cmd/main.go` (pageShowHandler)

**验收标准**:
- [ ] 路由支持模板选择
- [ ] 模板不存在时回退默认

---

### 任务 3.3: 管理后台模板选择
**文件**: `admin/src/views/pages/PageEdit.vue`

**验收标准**:
- [ ] 下拉框显示模板列表
- [ ] 选择后保存到数据库

---

## Phase 4: 模块化首页（3天）

### 任务 4.1: 模块配置定义
**文件**: `templates/themes/*/settings.yaml`

**验收标准**:
- [ ] 配置格式正确
- [ ] repeater 类型支持

---

### 任务 4.2: 首页模块化渲染
**文件**: `templates/themes/*/pages/home.html`

**验收标准**:
- [ ] 模块可启用/禁用
- [ ] 模块内容可配置

---

### 任务 4.3: Repeater 组件
**文件**: `admin/src/components/Repeater.vue`

**验收标准**:
- [ ] 组件实现
- [ ] 与配置系统集成
- [ ] 数据正确保存

---

## 任务依赖关系

```
Phase 1: API 密钥系统
├── 1.1 模型 → 1.2 CRUD API → 1.3 中间件 → 1.4 界面
│
Phase 2: AI 文档接口（依赖 Phase 1）
├── 2.1 能力发现 → 2.2 操作指南 → 2.3 详细文档 → 2.4 OpenAPI
│
Phase 3: 页面模板系统（独立）
├── 3.1 配置定义 → 3.2 路由修改 → 3.3 管理界面
│
Phase 4: 模块化首页（独立）
├── 4.1 配置定义 → 4.2 模板修改 → 4.3 Repeater 组件
```

---

## 执行顺序

**并行执行**:
- Phase 3 和 Phase 4 可以并行执行

**串行执行**:
- Phase 2 必须在 Phase 1 之后

**建议顺序**:
1. Phase 1: API 密钥系统
2. Phase 3 + Phase 4: 并行执行
3. Phase 2: AI 文档接口

---

## 验收标准

**每个任务完成后**:
- [ ] 代码编译通过
- [ ] 功能测试通过
- [ ] API 响应正确

**全部完成后**:
- [ ] 全面编译测试
- [ ] 端到端测试
- [ ] 文档更新

---

*计划待执行*
