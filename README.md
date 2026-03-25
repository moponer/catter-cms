# Catter CMS

一款基于 Go + Gin + Vue 3 + Element Plus 的现代化内容管理系统，支持主题定制、定时任务、Sitemap 生成等功能，开箱即用。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.25 + Gin Framework + GORM + SQLite |
| 前端 | Vue 3 + Element Plus + Vite |
| 认证 | JWT |
| 部署 | 单二进制文件，无额外依赖 |

## 功能特性

### 核心功能
- ✅ **用户认证** - JWT 登录/登出，角色权限控制
- ✅ **内容管理** - 文章、页面、分类、标签、评论
- ✅ **媒体库** - 图片/文件上传与管理
- ✅ **菜单管理** - 可视化菜单配置
- ✅ **主题系统** - 多主题支持，配置可视化编辑
- ✅ **定时任务** - 备份、清理、定时发布
- ✅ **Sitemap** - 自动生成 XML Sitemap

### 安全特性
- ✅ XSS 防护（用户输入自动转义）
- ✅ CSRF 防护
- ✅ 安全响应头（CSP、X-Frame-Options 等）
- ✅ JWT 密钥强度验证（生产环境 >=32 字符）
- ✅ 随机密码生成

### 系统特性
- ✅ 配置热加载（环境变量覆盖）
- ✅ 开发/生产环境切换
- ✅ Gzip 压缩
- ✅ SPA 路由支持（管理后台刷新不空白）

## 快速开始

### 环境要求
- Go 1.25+
- Node.js 18+
- npm / yarn

### 1. 克隆项目
```bash
git clone <repository-url>
cd Catter-CMS
```

### 2. 配置
编辑 `config.yaml` 或使用环境变量：

```yaml
app:
  port: "3000"
  env: "development"  # development / production
  jwt_secret: "your-secret-key-min-32-chars"
  admin_dist: "./admin/dist"
```

**环境变量覆盖**：
```bash
export APP_ENV=production
export JWT_SECRET=your-production-secret-key-min-32-chars
export PORT=3000
```

### 3. 构建
```bash
# 构建后端
cd gin-app
go build -o catter-cms.exe ./cmd

# 构建前端
cd ../admin
npm install
npm run build
```

### 4. 运行
```bash
# Windows
.\run.bat

# 或直接运行
./gin-app/catter-cms.exe
```

访问 `http://localhost:3000/admin` 进入管理后台。

### 5. 开发模式
```bash
# 后端热重载
cd gin-app
go run ./cmd

# 前端开发服务器（可选）
cd admin
npm run dev
```

## 项目结构

```
Catter-CMS/
├── admin/                      # Vue 3 前端管理后台
│   ├── src/
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 公共组件
│   │   ├── router/            # 路由配置
│   │   ├── utils/             # 工具函数
│   │   └── App.vue            # 根组件
│   ├── dist/                  # 构建输出（部署用）
│   └── vite.config.*          # Vite 配置
│
├── gin-app/                    # Go 后端
│   ├── cmd/
│   │   └── main.go            # 入口文件
│   ├── config/
│   │   └── config.go          # 配置管理
│   ├── internal/
│   │   ├── handlers/          # HTTP 处理器
│   │   ├── middleware/        # 中间件（认证、安全、SPA）
│   │   ├── models/           # 数据模型
│   │   └── services/          # 业务逻辑
│   └── templates/             # SSR 模板
│
├── templates/                  # 主题模板
│   └── themes/
│       ├── catter-blog/       # 博客主题
│       └── catter-cms-official/ # 企业主题
│
├── config.yaml                 # 应用配置
├── run.bat                     # Windows 启动脚本
└── catter.db                   # SQLite 数据库
```

## API 接口

### 认证
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/auth/login` | 用户登录 |
| GET | `/api/auth/me` | 获取当前用户 |

### 内容
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/posts` | 文章列表 |
| GET | `/api/posts/:id` | 文章详情 |
| POST | `/api/posts` | 创建文章 |
| PUT | `/api/posts/:id` | 更新文章 |
| DELETE | `/api/posts/:id` | 删除文章 |

### 分类/标签/媒体/用户
类似 RESTful 设计，完整 API 列表见 [API_STATUS.md](API_STATUS.md)

### 系统
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/settings` | 获取设置 |
| PUT | `/api/settings` | 保存设置 |
| GET | `/api/themes` | 主题列表 |
| POST | `/api/themes/upload` | 上传主题 |
| POST | `/api/themes/:id/activate` | 激活主题 |
| GET | `/sitemap.xml` | Sitemap |

## 主题开发

主题放在 `templates/themes/` 目录下，每个主题包含：

```
theme-name/
├── theme.yml           # 主题元数据
├── settings.yaml       # 配置表单定义（可选）
├── assets/             # 静态资源
│   ├── css/base.css    # 基础样式
│   └── js/main.js
└── pages/             # 页面模板
    ├── home.html       # 首页
    ├── article_list.html  # 文章列表
    ├── article_show.html  # 文章详情
    └── ...
```

详细开发指南见 `docs/THEME_DEVELOPER_GUIDE_V3.md`

## 配置说明

### config.yaml
```yaml
app:
  port: "3000"           # 服务端口
  env: "development"    # 环境模式
  jwt_secret: "dev-secret-key-change-in-production"
  admin_dist: "./admin/dist"  # 前端构建目录
  log_level: "debug"    # 日志级别

database:
  type: "sqlite"
  path: "./catter.db"

scheduler:
  enabled: true
  max_retries: 3
```

### 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `APP_ENV` | 环境模式 | `development` |
| `PORT` | 服务端口 | `3000` |
| `JWT_SECRET` | JWT 密钥（生产环境必填） | - |
| `DB_PATH` | 数据库路径 | `./catter.db` |

## 部署

### 生产环境
```bash
# 1. 设置环境变量
export APP_ENV=production
export JWT_SECRET=your-production-secret-key-min-32-chars

# 2. 构建
cd gin-app && go build -o catter-cms ./cmd
cd admin && npm install && npm run build

# 3. 运行
./catter-cms
```

### Docker（可选）
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o catter-cms ./gin-app/cmd

FROM alpine:latest
COPY --from=builder /app/catter-cms .
COPY --from=builder /app/admin/dist ./admin/dist
COPY --from=builder /app/templates ./templates
EXPOSE 3000
CMD ["./catter-cms"]
```

## 安全建议

1. **修改默认密码** - 首次部署后立即修改管理员密码
2. **设置强 JWT 密钥** - 生产环境使用 >=32 字符的随机字符串
3. **启用生产模式** - 设置 `APP_ENV=production`
4. **配置 HTTPS** - 使用反向代理（Nginx/Caddy）终止 TLS
5. **定期备份** - 使用定时任务自动备份数据库

## 许可证

MIT License
