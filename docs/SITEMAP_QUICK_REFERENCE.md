# Sitemap功能快速参考

## 概述
Catter CMS已实现完整的Sitemap生成功能，支持自动生成站点地图，符合搜索引擎优化标准。

## 访问地址

### 公开访问（无需认证）
- **Sitemap XML**: http://127.0.0.1:3001/sitemap.xml
- **浏览器访问**: 直接打开上述URL即可查看XML内容

### 管理后台API（需要认证）
- **获取Sitemap信息**: GET `/api/admin/sitemap/info`
- **重建Sitemap**: POST `/api/admin/sitemap/rebuild`

## 使用方法

### 方法1: 自动生成（推荐）
应用启动后会自动生成sitemap，无需手动操作。

### 方法2: 手动重建
1. 登录管理后台: http://127.0.0.1:3001/admin
2. 导航到: 设置 → Sitemap管理
3. 点击"重建Sitemap"按钮

### 方法3: API调用
```bash
# 重建Sitemap
curl -X POST http://127.0.0.1:3001/api/admin/sitemap/rebuild \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取统计信息
curl -X GET http://127.0.0.1:3001/api/admin/sitemap/info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Sitemap内容

### 包含的URL类型
1. **首页** - 优先级 1.0，更新频率 daily
2. **文章页** - 优先级 0.8，更新频率 weekly
3. **分类页** - 优先级 0.6，更新频率 weekly
4. **页面** - 优先级 0.5，更新频率 monthly

### XML格式示例
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://127.0.0.1:3001/</loc>
    <lastmod>2026-03-16</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>http://127.0.0.1:3001/blog/sample-post</loc>
    <lastmod>2026-03-16</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
</urlset>
```

## 技术细节

### 路由配置
**文件**: `actions/app.go`
```go
// Sitemap API（必须在静态文件服务之前）
App.GET("/sitemap.xml", admin.SitemapHandler)
protectedAPI.GET("/admin/sitemap/info", admin.SitemapInfoHandler)
protectedAPI.POST("/admin/sitemap/rebuild", admin.RebuildSitemapHandler)

// 静态文件服务
App.ServeFiles("/", http.Dir("public"))
```

### 服务实现
**文件**: `services/sitemap.go`
- `Generate()` - 生成sitemap XML
- `Rebuild()` - 重建sitemap
- `SetBaseURL()` - 设置基础URL
- `GetURLCount()` - 获取URL数量

### API处理器
**文件**: `actions/admin/sitemap.go`
- `SitemapHandler` - 返回sitemap.xml
- `SitemapInfoHandler` - 返回统计信息
- `RebuildSitemapHandler` - 重建sitemap

## 管理界面

### 路径
**文件**: `admin/src/views/sitemap/Index.vue`

### 功能
1. 查看Sitemap统计信息
2. 手动重建Sitemap
3. 预览Sitemap内容
4. 配置Sitemap参数

### 使用步骤
1. 登录管理后台
2. 点击左侧菜单 "系统设置"
3. 点击 "Sitemap管理"
4. 查看统计信息或点击"重建Sitemap"

## 验证和测试

### 快速验证脚本
```powershell
# 运行验证脚本
.\test_sitemap.ps1
```

### 手动验证
1. 启动应用: `.\catter-cms-test.exe`
2. 访问: http://127.0.0.1:3001/sitemap.xml
3. 检查返回的XML内容

### 在线验证工具
- Google Search Console
- Bing Webmaster Tools
- XML Sitemap Validator

## 配置选项

### 基础URL设置
Sitemap会自动使用请求的Host和Protocol，或从session中获取`base_url`。

### 优先级配置
```go
const (
    PriorityHomepage = "1.0"  // 首页
    PriorityPost     = "0.8"  // 文章
    PriorityPage     = "0.5"  // 页面
    PriorityCategory = "0.6"  // 分类
)
```

### 更新频率
- 首页: daily
- 文章: weekly
- 页面: monthly
- 分类: weekly

## 故障排除

### 问题1: 访问sitemap.xml返回404
**解决方案**: 确保路由在静态文件服务之前（已修复）

### 问题2: Sitemap内容不完整
**解决方案**: 点击"重建Sitemap"按钮重新生成

### 问题3: 缓存问题
**解决方案**: 清除浏览器缓存或等待1小时（缓存时间）

## 相关文档

- [Sitemap实现文档](SITEMAP_IMPLEMENTATION.md) - 详细技术实现
- [Sitemap代码审查报告](SITEMAP_CODE_REVIEW.md) - 代码质量审查
- [Sitemap修复总结](SITEMAP_FIX_SUMMARY.md) - 修复历程
- [Sitemap最终修复](SITEMAP_FINAL_FIX.md) - 最终修复报告
- [Sitemap 404修复](SITEMAP_404_FIX.md) - 404问题修复
- [Sitemap使用演示](SITEMAP_DEMO.md) - 使用演示

## 最佳实践

1. **定期重建**: 发布新内容后手动重建sitemap
2. **提交搜索引擎**: 将sitemap.xml提交给Google和Bing
3. **监控错误**: 使用搜索引擎工具监控sitemap错误
4. **保持更新**: 确保sitemap与站点内容同步

## 搜索引擎提交

### Google
1. 登录 Google Search Console
2. 选择站点
3. 导航到 Sitemaps
4. 输入: `sitemap.xml`
5. 点击提交

### Bing
1. 登录 Bing Webmaster Tools
2. 选择站点
3. 导航到 Sitemaps
4. 输入: `sitemap.xml`
5. 点击提交

## 性能优化

### 缓存策略
- HTTP缓存: `Cache-Control: public, max-age=3600` (1小时)
- 建议: 每次发布内容后手动重建

### 文件大小
- Sitemap大小通常 < 1MB
- 支持大量URL（最多50,000个）
- 超过限制可使用Sitemap索引

## 总结

✅ **功能完整**: 支持自动生成和手动重建
✅ **符合标准**: 遵循sitemaps.org规范
✅ **易于使用**: 提供管理界面和API
✅ **生产就绪**: 已通过编译和功能验证
✅ **SEO友好**: 支持搜索引擎优化

**状态**: 生产就绪 ✅
