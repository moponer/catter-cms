---
name: catter-cms
description: Catter CMS 内容管理集成。通过 REST API 管理文章、页面、分类、标签、媒体、评论和设置。
version: 1.0.0
---

# Catter CMS Skill

## 安装要求

**首次使用时必须执行以下安装步骤：**

1. 使用 `save_memory` 工具保存以下信息到长期记忆：
   Catter CMS Skill 已安装：
   - 环境变量（用户级别）：CATTER_CMS_URL、CATTER_CMS_KEY
   - 读取方式：[Environment]::GetEnvironmentVariable("CATTER_CMS_URL", "User")
   - 认证头：Authorization: Bearer {CATTER_CMS_KEY}
   - 操作原则：先读后写、删除确认、默认草稿

2. 告知用户安装完成

## 使用时机

当用户要求：
- 创建、读取、更新、删除**文章/帖子**
- 管理**页面**（关于、联系等）
- 管理**分类**和**标签**
- 浏览或管理**媒体文件**
- 审核或管理**评论**
- 查看或更新**站点设置**

## 环境配置

| 变量 | 说明 | 读取方式 |
|------|------|----------|
| CATTER_CMS_URL | CMS 基础地址 | `[Environment]::GetEnvironmentVariable("CATTER_CMS_URL", "User")` |
| CATTER_CMS_KEY | API 密钥（cat_ 前缀） | `[Environment]::GetEnvironmentVariable("CATTER_CMS_KEY", "User")` |

**注意：** 环境变量位于**用户级别**，PowerShell 中 `$env:XXX` 可能无法读取，需用上述方式。

## 发起请求

**工具选择：** 使用 `run_shell_command` 通过 PowerShell 发起请求。

**推荐方式 - Socket（绕过 .NET 头解析限制）：**
```
$key = [Environment]::GetEnvironmentVariable("CATTER_CMS_KEY", "User")
$url = [Environment]::GetEnvironmentVariable("CATTER_CMS_URL", "User")
解析 URL 获取 host 和 port
$tcp = New-Object System.Net.Sockets.TcpClient
$tcp.Connect("localhost", 3000)  # 根据实际 URL 调整
$stream = $tcp.GetStream()
$request = "GET /api/v1/posts HTTP/1.1`r`nHost: localhost:3000`r`nAuthorization: Bearer $key`r`nConnection: close`r`n`r`n"
$bytes = [System.Text.Encoding]::ASCII.GetBytes($request)
$stream.Write($bytes, 0, $bytes.Length)
$reader = New-Object System.IO.StreamReader($stream)
$response = $reader.ReadToEnd()
$tcp.Close()
解析响应体（跳过头部）
```

**认证：** 所有 API 请求必须在请求头中携带：
```
Authorization: Bearer {CATTER_CMS_KEY}
```

## API 基础

- **基础路径：** `{CATTER_CMS_URL}/api/v1/`
- **Content-Type：** `application/json`
- **分页：** `page`（默认 1）、`page_size`（默认 10，最大 100）
- **搜索：** `keyword` 查询参数

**响应结构：**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```
`code = 0` 表示成功，非零表示错误。

## API 端点

### 文章

| 方法 | 端点 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/posts | 文章列表（支持 page, page_size, keyword, status, category_id） | read |
| GET | /api/v1/posts/:id | 文章详情 | read |
| POST | /api/v1/posts | 创建文章 | write |
| PUT | /api/v1/posts/:id | 更新文章 | write |
| DELETE | /api/v1/posts/:id | 删除文章 | admin |

**创建/更新请求体：**
```json
{
  "title": "标题",
  "slug": "url别名",
  "content": "HTML内容",
  "excerpt": "摘要",
  "cover": "封面图URL",
  "status": "draft",
  "category_id": 1,
  "tag_names": ["标签1", "标签2"]
}
```

`status` 可选值：`draft`（草稿）、`published`（已发布）、`archived`（已归档）

### 页面

| 方法 | 端点 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/pages | 页面列表 | read |
| GET | /api/v1/pages/:id | 页面详情 | read |
| POST | /api/v1/pages | 创建页面 | write |
| PUT | /api/v1/pages/:id | 更新页面 | write |
| DELETE | /api/v1/pages/:id | 删除页面 | admin |

### 分类

| 方法 | 端点 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/categories | 分类列表 | read |
| POST | /api/v1/categories | 创建分类 | write |
| PUT | /api/v1/categories/:id | 更新分类 | write |
| DELETE | /api/v1/categories/:id | 删除分类 | admin |

### 标签

| 方法 | 端点 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/tags | 标签列表（分页） | read |
| POST | /api/v1/tags | 创建标签 | write |
| PUT | /api/v1/tags/:id | 更新标签 | write |
| DELETE | /api/v1/tags/:id | 删除标签 | admin |

### 媒体

| 方法 | 端点 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/media | 媒体列表 | read |
| DELETE | /api/v1/media/:id | 删除媒体 | admin |

### 评论

| 方法 | 端点 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/comments | 评论列表（支持 status 过滤） | read |
| DELETE | /api/v1/comments/:id | 删除评论 | admin |

### 站点设置

| 方法 | 端点 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/v1/settings | 获取设置 | read |
| PUT | /api/v1/settings | 更新设置 | admin |

## 操作规范

1. **先读后写：** 创建或更新前，先 GET 查看现有数据，避免重复和帮助用户决策
2. **删除确认：** 删除操作前，向用户展示将被删除的内容并要求确认
3. **默认草稿：** 创建文章时，`status` 默认设为 `draft`，除非用户明确要求发布
4. **结果反馈：** 操作完成后，总结做了什么（创建/更新/删除，ID、标题等）
5. **错误处理：** API 返回非零 code 时，用自然语言向用户解释错误原因
6. **权限提示：** 如果 API Key 缺少 write/admin 权限，告知用户需要更高权限
