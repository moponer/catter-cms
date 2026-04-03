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
  "content": "文章内容（见下方格式说明）",
  "content_type": "markdown",
  "excerpt": "摘要",
  "cover": "封面图URL",
  "status": "draft",
  "category_id": 1,
  "tag_names": ["标签1", "标签2"]
}
```

**`content_type` 参数：**
- `"html"`（默认）：内容已经是 HTML 格式，直接存储
- `"markdown"`：**推荐**。AI 生成 Markdown，后端自动转换为 HTML 存储

**`status` 可选值：** `draft`（草稿）、`published`（已发布）、`archived`（已归档）

### 文章内容格式规范

**推荐：使用 Markdown 格式（设置 `content_type: "markdown"`）**

AI 生成文章时应使用标准 Markdown 语法，确保前台正确渲染：

```markdown
# 文章标题（一级标题，自动从 title 字段获取，可选）

这是文章开头的介绍文字...

## 二级标题

正文内容，支持**粗体**、*斜体*、~~删除线~~

### 列表

- 列表项一
- 列表项二
- 列表项三

### 有序列表

1. 第一步
2. 第二步
3. 第三步

### 引用

> 这是一段引用文字

### 代码块

```javascript
const greeting = "Hello, World!";
console.log(greeting);
```

### 图片

![图片描述](https://example.com/image.jpg)

### 链接

[链接文字](https://example.com)

### 表格

| 列1 | 列2 | 列3 |
|-----|-----|-----|
| A   | B   | C   |
| D   | E   | F   |
```

**图片使用规范：**
- 使用图床或绝对 URL（如 `https://example.com/image.jpg`）
- 不要使用相对路径
- 建议先通过 `POST /api/v1/media` 上传图片获取 URL
- 每个 `![]()` 图片语法不要超过 200 字符

**排版要点：**
- 标题层级清晰，不要跳级（一级 → 二级 → 三级）
- 段落之间留空行
- 列表项不要过长（每项建议 100 字以内）
- 重要内容可用加粗强调，但不要整段加粗

**加粗语法注意事项：**
- 中文环境下，`**文字：**`（标点在加粗内部）可能导致渲染失败
- **正确写法**：`**标题文字**：后续内容`（标点在加粗外部）
- **错误写法**：`**标题文字：**后续内容`（标点在加粗内部）
- 将标点符号放在 `**` 外部，确保正确渲染

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
| GET | /api/v1/media | 媒体列表（支持 page, page_size, file_type, search） | read |
| GET | /api/v1/media/:id | 媒体详情 | read |
| POST | /api/v1/media/upload | 上传文件（multipart/form-data，字段名 file） | write |
| POST | /api/v1/media/upload-cover | 上传封面图（自动压缩到 1200px 宽，85% 质量） | write |
| DELETE | /api/v1/media/:id | 删除媒体（同时删除物理文件和数据库记录） | admin |

**上传响应格式：**
```json
{
  "code": 0,
  "message": "文件上传成功",
  "data": {
    "id": 1,
    "filename": "原始文件名.jpg",
    "path": "服务器文件路径",
    "url": "/uploads/1234567890_123456.jpg",
    "file_size": 102400,
    "uploaded_by": 1,
    "created_at": "2026-04-02T15:00:00Z"
  }
}
```

**封面图上传额外返回：**
```json
{
  "data": {
    "id": 1,
    "url": "/uploads/cover_xxx.jpg",
    "width": 1200,
    "height": 675,
    "original_size": 500000,
    "compressed_size": 120000
  }
}
```

### 图片上传与使用流程

当用户提供图片（本地文件、截图等）需要在文章中使用：

**Step 1: 上传图片获取 URL**

使用 PowerShell Socket 上传图片到 CMS 媒体库（Windows 推荐）：

```powershell
$cmsUrl = [Environment]::GetEnvironmentVariable("CATTER_CMS_URL", "User")
$apiKey = [Environment]::GetEnvironmentVariable("CATTER_CMS_KEY", "User")
$filePath = "C:\path\to\image.jpg"

# 读取文件
$fileBytes = [System.IO.File]::ReadAllBytes($filePath)
$fileName = [System.IO.Path]::GetFileName($filePath)

# 构造 multipart/form-data 请求体
$boundary = [System.Guid]::NewGuid().ToString()
$LF = "`r`n"
$bodyLines = @(
    "--$boundary",
    "Content-Disposition: form-data; name=`"file`"; filename=`"$fileName`"",
    "Content-Type: application/octet-stream$LF",
    ([System.Text.Encoding]::GetEncoding("iso-8859-1").GetString($fileBytes)),
    "--$boundary--$LF"
) -join $LF

$bodyBytes = [System.Text.Encoding]::GetEncoding("iso-8859-1").GetBytes($bodyLines)

# 解析 URL
$uri = [System.Uri]$cmsUrl
$hostName = $uri.Host
$port = if ($uri.Port -gt 0) { $uri.Port } else { 80 }

# Socket 连接
$tcp = New-Object System.Net.Sockets.TcpClient
$tcp.Connect($hostName, $port)
$stream = $tcp.GetStream()

$request = "POST /api/v1/media/upload HTTP/1.1$LF" +
    "Host: $hostName`:$port$LF" +
    "Authorization: Bearer $apiKey$LF" +
    "Content-Type: multipart/form-data; boundary=$boundary$LF" +
    "Content-Length: $($bodyBytes.Length)$LF" +
    "Connection: close$LF$LF"

$requestBytes = [System.Text.Encoding]::ASCII.GetBytes($request)
$stream.Write($requestBytes, 0, $requestBytes.Length)
$stream.Write($bodyBytes, 0, $bodyBytes.Length)

$reader = New-Object System.IO.StreamReader($stream)
$response = $reader.ReadToEnd()
$tcp.Close()

# 解析响应（跳过 HTTP 头，提取 JSON body）
$jsonStart = $response.IndexOf("{")
$jsonBody = $response.Substring($jsonStart)
$result = $jsonBody | ConvertFrom-Json

# 获取图片 URL
$imageUrl = "$cmsUrl$($result.data.url)"
Write-Host "图片 URL: $imageUrl"
```

**Step 2: 在文章中使用图片 URL**

将获取到的 URL 插入到 Markdown 内容中：

```markdown
![图片描述](https://your-cms.com/uploads/1234567890_123456.jpg)
```

**Step 3: 创建/更新文章时引用图片**

```json
{
  "title": "文章标题",
  "content": "正文内容...![配图](https://your-cms.com/uploads/xxx.jpg)...更多内容",
  "content_type": "markdown",
  "cover": "https://your-cms.com/uploads/cover_xxx.jpg",
  "status": "draft"
}
```

**注意事项：**
- 图片必须先上传到 CMS 媒体库，不要使用外部图床
- 上传后建议通过 `GET /api/v1/media` 确认上传成功
- 封面图使用 `upload-cover` 接口会自动压缩优化
- 普通图片使用 `upload` 接口保持原图
- 每个 `![]()` 图片语法不要超过 200 字符
- 不要使用相对路径，必须使用完整 URL

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
