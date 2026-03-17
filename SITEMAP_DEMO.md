# Sitemap 功能演示

## 功能概述

Catter CMS 的 Sitemap 生成功能自动为您的网站创建符合搜索引擎优化（SEO）标准的 XML sitemap，帮助搜索引擎更好地索引您的内容。

## 主要特性

### 🎯 自动生成
- 自动收录所有已发布的文章
- 自动收录所有已发布的页面
- 自动收录所有分类
- 自动包含首页

### ⚙️ 智能配置
- 可自定义站点基础 URL
- 可调整更新频率（从总是到从不）
- 可设置内容优先级（0.0 - 1.0）
- 可选择包含的内容类型

### 📊 实时统计
- 显示当前 sitemap 中的 URL 总数
- 显示最后更新时间
- 一键重新生成 sitemap

### 🔍 SEO 优化
- 自动设置 `lastmod`（最后修改时间）
- 合理的 `changefreq`（更新频率）
- 科学的 `priority`（优先级）分配
- 符合 sitemaps.org 标准

## 快速开始

### 1. 访问 Sitemap 管理页面

登录管理后台后，点击左侧菜单的 **Sitemap** 选项。

### 2. 配置基础信息

在配置表单中填写以下信息：

```
站点URL: https://yourdomain.com
更新频率: weekly
包含内容: ✅ 文章 ✅ 页面 ✅ 分类
文章优先级: 0.8
页面优先级: 0.6
```

### 3. 生成 Sitemap

点击 **重新生成 Sitemap** 按钮，系统会立即生成新的 sitemap.xml 文件。

### 4. 访问 Sitemap

生成后，您可以：
- 点击 **查看 Sitemap** 在新窗口打开
- 点击 **复制 Sitemap URL** 复制链接

## 使用示例

### 示例 1: 博客网站

对于典型的博客网站，推荐配置：

```yaml
站点URL: https://myblog.com
更新频率: daily
包含内容:
  - 文章: ✓
  - 页面: ✓
  - 分类: ✓
优先级设置:
  - 文章: 0.8
  - 页面: 0.6
```

**生成的 sitemap 包含：**
- 首页（优先级 1.0，更新频率 daily）
- 所有已发布文章（优先级 0.8，更新频率 weekly）
- 所有页面（优先级 0.6，更新频率 monthly）
- 所有分类（优先级 0.5，更新频率 weekly）

### 示例 2: 企业官网

对于内容更新较少的企业官网：

```yaml
站点URL: https://company.com
更新频率: weekly
包含内容:
  - 文章: ✓
  - 页面: ✓
  - 分类: ✗
优先级设置:
  - 文章: 0.7
  - 页面: 0.9
```

**生成的 sitemap 包含：**
- 首页（优先级 1.0，更新频率 weekly）
- 所有新闻文章（优先级 0.7，更新频率 monthly）
- 所有页面（优先级 0.9，更新频率 monthly）

## Sitemap URL 访问

### 公开访问
您的 sitemap 可以通过以下 URL 公开访问：

```
https://yourdomain.com/sitemap.xml
```

**注意：** 此 URL 无需登录即可访问，方便搜索引擎抓取。

### 查看示例

访问 `https://yourdomain.com/sitemap.xml` 您将看到类似以下内容：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://yourdomain.com/</loc>
    <lastmod>2026-03-16T10:30:00Z</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://yourdomain.com/posts/welcome-to-catter-cms</loc>
    <lastmod>2026-03-16T09:00:00Z</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>https://yourdomain.com/posts/how-to-use-catter-cms</loc>
    <lastmod>2026-03-15T14:30:00Z</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>https://yourdomain.com/pages/about-us</loc>
    <lastmod>2026-03-10T11:00:00Z</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.6</priority>
  </url>
  <url>
    <loc>https://yourdomain.com/categories/technology</loc>
    <lastmod>2026-03-16T08:00:00Z</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.5</priority>
  </url>
</urlset>
```

## 提交到搜索引擎

### Google Search Console

1. 访问 [Google Search Console](https://search.google.com/search-console)
2. 选择您的网站
3. 在左侧菜单选择 "站点地图"
4. 点击 "添加新的站点地图"
5. 输入 `sitemap.xml`
6. 点击 "提交"

### Bing Webmaster Tools

1. 访问 [Bing Webmaster Tools](https://www.bing.com/webmasters)
2. 登录并选择您的网站
3. 在左侧菜单选择 "Sitemaps"
4. 在 "提交站点地图" 框中输入 `sitemap.xml`
5. 点击 "提交"

### Yandex.Webmaster

1. 访问 [Yandex.Webmaster](https://webmaster.yandex.com)
2. 选择您的网站
3. 在 "索引" 部分选择 "站点地图"
4. 添加 `sitemap.xml`
5. 点击 "添加"

### 添加到 robots.txt

在您的网站根目录创建或编辑 `robots.txt` 文件，添加以下内容：

```txt
User-agent: *
Allow: /

Sitemap: https://yourdomain.com/sitemap.xml
```

## 定时更新建议

为了保持 sitemap 的最新状态，建议：

### 选项 1: 手动更新
- 每次发布新内容后，在管理后台点击 "重新生成 Sitemap"
- 适合内容更新不频繁的网站

### 选项 2: 定时任务
使用 Catter CMS 的定时任务功能，设置每天自动重建 sitemap：

```
任务名称: 每日重建 Sitemap
任务类型: 自定义命令
Cron 表达式: 0 2 * * *
命令: curl -X POST http://yourdomain.com/api/admin/sitemap/rebuild
```

## 故障排除

### 问题 1: Sitemap 访问 404

**原因：** 路由未正确配置

**解决方案：**
1. 检查 `actions/app.go` 中是否包含 sitemap 路由
2. 确保应用已重新编译和部署
3. 重启应用服务

### 问题 2: Sitemap 内容为空

**原因：** 数据库中没有已发布的内容

**解决方案：**
1. 确保至少有一篇已发布的文章或页面
2. 检查内容的 `status` 字段是否为 `published`
3. 在管理后台创建测试内容

### 问题 3: URL 不正确

**原因：** 基础 URL 配置错误

**解决方案：**
1. 在 Sitemap 管理页面检查 "站点URL" 配置
2. 确保以 `http://` 或 `https://` 开头
3. 确保不以 `/` 结尾
4. 重新生成 sitemap

### 问题 4: 搜索引擎无法抓取

**原因：** 服务器阻止或防火墙限制

**解决方案：**
1. 检查服务器日志
2. 确保搜索引擎 IP 允许访问
3. 检查 robots.txt 配置
4. 验证 sitemap.xml 可以正常访问

## 性能优化

### 缓存策略

Sitemap 默认缓存 1 小时，这意味着：
- 搜索引擎在 1 小时内多次访问会返回缓存的 sitemap
- 减少数据库查询压力
- 提高响应速度

### 大型网站优化

如果您的网站有大量内容（超过 50,000 个 URL），建议：

1. **使用 Sitemap Index**
   - 将 sitemap 分割成多个文件
   - 创建 sitemap index 文件
   - 每个子 sitemap 最多 50,000 个 URL

2. **异步生成**
   - 在后台定时生成 sitemap
   - 避免请求时实时生成
   - 提高用户体验

3. **CDN 加速**
   - 将生成的 sitemap 部署到 CDN
   - 减轻源服务器压力
   - 加快全球访问速度

## 最佳实践

### 1. 定期更新
- 每次发布新内容后重建 sitemap
- 设置定时任务自动更新
- 保持 sitemap 始终是最新的

### 2. 合理设置优先级
- 首页: 1.0
- 重要页面: 0.9-1.0
- 常规内容: 0.6-0.8
- 分类/标签: 0.5
- 不重要页面: 0.1-0.4

### 3. 设置更新频率
- 频繁更新: daily, hourly
- 定期更新: weekly, monthly
- 很少更新: yearly, never

### 4. 监控索引状态
- 在 Google Search Console 监控索引状态
- 检查是否有错误
- 及时修复问题

### 5. 保持简洁
- 只包含真正需要索引的 URL
- 排除重定向页面
- 排除已删除内容

## 技术细节

### 支持的内容类型

| 类型 | URL 模式 | 默认优先级 | 默认频率 |
|------|----------|-----------|----------|
| 首页 | `/` | 1.0 | daily |
| 文章 | `/posts/{slug}` | 0.8 | weekly |
| 页面 | `/pages/{slug}` | 0.6 | monthly |
| 分类 | `/categories/{slug}` | 0.5 | weekly |

### API 端点

```
GET  /sitemap.xml              # 获取 sitemap XML
GET  /api/admin/sitemap/info   # 获取统计信息
POST /api/admin/sitemap/rebuild # 重建 sitemap
```

### 数据来源

- **文章**: `posts` 表中 `status = 'published'` 的记录
- **页面**: `pages` 表中 `status = 'published'` 的记录
- **分类**: `categories` 表中的所有记录

## 总结

Catter CMS 的 Sitemap 功能提供了：
- ✅ 自动化的 sitemap 生成
- ✅ 灵活的配置选项
- ✅ 完善的管理界面
- ✅ 符合 SEO 标准
- ✅ 易于使用的操作

通过合理配置和使用 sitemap，可以显著提升您网站在搜索引擎中的表现，提高内容索引效率，最终带来更多的有机流量。

立即开始使用 Sitemap 功能，让您的网站获得更好的搜索引擎优化效果！
