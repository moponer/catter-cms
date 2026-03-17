# Catter CMS

Catter CMS 是一个基于 Go + Buffalo + Vue 3 + Element Plus 开发的内容管理系统，提供了完整的内容管理功能，包括文章管理、分类管理、用户管理、主题管理等。

## 技术栈

- **后端**：Go 1.25.0 + Buffalo 1.1.3 + SQLite
- **前端**：Vue 3 + Element Plus + Vite
- **认证**：JWT

## 功能特性

- ✅ 用户认证系统（登录/登出）
- ✅ 文章管理（创建/编辑/删除/发布）
- ✅ 分类管理
- ✅ 标签管理
- ✅ 媒体管理
- ✅ 主题管理
- ✅ 系统设置
- ✅ API 接口

## 安装指南

### 1. 环境要求

- Go 1.25.0 或更高版本
- Node.js 18.0 或更高版本
- npm 或 yarn

### 2. 后端安装

```bash
# 克隆仓库
git clone <repository-url>
cd catter-cms

# 安装依赖
go mod tidy

# 初始化数据库
buffalo pop create -a

# 运行开发服务器
go run cmd/app/main.go
```

### 3. 前端安装

```bash
# 进入前端目录
cd admin

# 安装依赖
npm install

# 构建前端资源
npm run build
```

## 运行方式

### 开发模式

```bash
# 后端
go run cmd/app/main.go

# 前端（开发模式）
cd admin
npm run dev
```

### 生产模式

```bash
# 构建前端资源
cd admin
npm run build

# 运行后端服务器
go run cmd/app/main.go
```

## API 接口

### 认证接口
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息

### 文章接口
- `GET /api/posts` - 获取文章列表
- `GET /api/posts/{id}` - 获取文章详情
- `POST /api/posts` - 创建文章
- `PUT /api/posts/{id}` - 更新文章
- `DELETE /api/posts/{id}` - 删除文章

### 分类接口
- `GET /api/categories` - 获取分类列表
- `GET /api/categories/{id}` - 获取分类详情
- `POST /api/categories` - 创建分类
- `PUT /api/categories/{id}` - 更新分类
- `DELETE /api/categories/{id}` - 删除分类

## 项目结构

```
catter-cms/
├── actions/           # 控制器
├── admin/             # 前端代码
├── cmd/               # 命令行入口
├── config/            # 配置文件
├── middleware/        # 中间件
├── models/            # 数据模型
├── repositories/      # 数据访问层
├── services/          # 业务逻辑层
├── templates/         # 服务器端模板
├── public/            # 静态文件
└── docs/              # 文档
```

## 配置说明

- 数据库配置：`database.yml`
- 应用配置：`config/config.go`

## 注意事项

- JWT 密钥目前硬编码在代码中，生产环境应使用环境变量或配置文件
- 前端构建时暂时跳过了类型检查，生产环境应修复类型问题
- 数据库默认使用 SQLite，生产环境可配置为 MySQL 或 PostgreSQL

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 许可证

MIT License
