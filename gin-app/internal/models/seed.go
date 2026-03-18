package models

import (
	"log"
	"time"
)

// SeedData 初始化示例数据
func SeedData() error {
	// 检查是否已有数据
	var postCount int64
	DB.Model(&Post{}).Count(&postCount)
	if postCount > 0 {
		log.Println("Data already exists, skipping seed")
		return nil
	}

	log.Println("Seeding initial data...")

	// 1. 确保有管理员用户
	if err := CreateDefaultAdmin(); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	// 获取管理员用户ID
	var admin User
	if err := DB.Where("role = ?", "admin").First(&admin).Error; err != nil {
		log.Printf("Warning: No admin user found: %v", err)
		return nil
	}

	// 2. 创建分类
	categories := []Category{
		{Name: "技术教程", Slug: "tutorial", Description: "各类技术教程和指南", SortOrder: 1},
		{Name: "生活随笔", Slug: "lifestyle", Description: "日常生活中的点滴记录", SortOrder: 2},
		{Name: "产品设计", Slug: "design", Description: "产品设计和用户体验相关", SortOrder: 3},
		{Name: "前端开发", Slug: "frontend", Description: "前端技术文章和最佳实践", SortOrder: 4},
		{Name: "后端开发", Slug: "backend", Description: "后端技术与架构", SortOrder: 5},
	}

	if err := DB.Create(&categories).Error; err != nil {
		log.Printf("Warning: Failed to create categories: %v", err)
	}

	// 3. 创建标签
	tags := []Tag{
		{Name: "Go", Slug: "go"},
		{Name: "Vue", Slug: "vue"},
		{Name: "React", Slug: "react"},
		{Name: "Node.js", Slug: "nodejs"},
		{Name: "Python", Slug: "python"},
		{Name: "数据库", Slug: "database"},
		{Name: "Docker", Slug: "docker"},
		{Name: "Kubernetes", Slug: "kubernetes"},
		{Name: "AI", Slug: "ai"},
		{Name: "机器学习", Slug: "machine-learning"},
	}

	if err := DB.Create(&tags).Error; err != nil {
		log.Printf("Warning: Failed to create tags: %v", err)
	}

	// 4. 创建文章
	now := time.Now()
	publishedAt := &now

	posts := []Post{
		{
			Title:       "欢迎使用 Catter CMS",
			Slug:        "welcome-to-catter-cms",
			Content:     "# 欢迎使用 Catter CMS\n\n这是一个现代化的博客系统，基于 Go 语言 Gin 框架开发。\n\n## 主要特性\n\n- 简洁易用的管理后台\n- 支持 Markdown 编辑器\n- 响应式设计，适配各种设备\n- 强大的主题系统\n- SEO 友好\n\n## 开始使用\n\n1. 登录管理后台\n2. 创建分类和标签\n3. 开始撰写文章\n4. 选择喜欢的主题\n\n祝您使用愉快！",
			Excerpt:     "这是一个现代化的博客系统，基于 Go 语言 Gin 框架开发。",
			Status:       "published",
			Views:        100,
			CategoryID:  &categories[0].ID,
			AuthorID:     admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "Go 语言快速入门指南",
			Slug:        "go-quick-start-guide",
			Content:     "# Go 语言快速入门\n\nGo 语言是一种简洁、高效、并发友好的编程语言。\n\n## 环境安装\n\n```bash\n# 下载 Go\nwget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz\n\n# 解压\nsudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz\n\n# 配置环境变量\nexport PATH=$PATH:/usr/local/go/bin\n```\n\n## 第一个程序\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```\n\n## 基本语法\n\n- 变量声明\n- 函数定义\n- 流程控制\n- 并发编程\n\nGo 语言的并发模型非常强大，使用 goroutine 和 channel 可以轻松实现并发编程。",
			Excerpt:     "Go 语言是一种简洁、高效、并发友好的编程语言，本文将带你快速入门。",
			Status:       "published",
			Views:        85,
			CategoryID:  &categories[4].ID,
			AuthorID:     admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "Vue 3 组合式 API 完全指南",
			Slug:        "vue3-composition-api-guide",
			Content:     "# Vue 3 组合式 API 完全指南\n\nVue 3 引入了组合式 API，这是 Vue 的一次重大升级。\n\n## 为什么选择组合式 API\n\n1. 更好的代码组织\n2. 更好的逻辑复用\n3. 更好的 TypeScript 支持\n\n## 基本用法\n\n```vue\n<script setup>\nimport { ref, computed, onMounted } from 'vue'\n\nconst count = ref(0)\nconst doubled = computed(() => count.value * 2)\n\nfunction increment() {\n  count.value++\n}\n\nonMounted(() => {\n  console.log('Component mounted!')\n})\n</script>\n\n<template>\n  <button @click=\"increment\">\n    Count: {{ count }}, Doubled: {{ doubled }}\n  </button>\n</template>\n```\n\n组合式 API 让组件逻辑更加清晰，是现代 Vue 开发的推荐方式。",
			Excerpt:     "Vue 3 引入了组合式 API，这是 Vue 的一次重大升级本文将详细介绍。",
			Status:       "published",
			Views:        120,
			CategoryID:  &categories[3].ID,
			AuthorID:     admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "Docker 入门实战教程",
			Slug:        "docker-getting-started",
			Content:     "# Docker 入门实战\n\nDocker 是现代软件开发和部署的重要工具。\n\n## 什么是 Docker\n\nDocker 是一个开源的容器化平台，可以将应用及其依赖打包成轻量级的容器。\n\n## 基本命令\n\n```bash\n# 拉取镜像\ndocker pull nginx\n\n# 运行容器\ndocker run -d -p 80:80 nginx\n\n# 查看容器\ndocker ps\n\n# 停止容器\ndocker stop <container_id>\n```\n\n## Dockerfile 示例\n\n```dockerfile\nFROM node:18-alpine\nWORKDIR /app\nCOPY package*.json ./\nRUN npm install\nCOPY . .\nEXPOSE 3000\nCMD [\"npm\", \"start\"]\n```\n\nDocker 让开发、测试、部署变得更加简单和一致。",
			Excerpt:     "Docker 是现代软件开发和部署的重要工具，本文将带你快速上手。",
			Status:       "published",
			Views:        90,
			CategoryID:  &categories[0].ID,
			AuthorID:     admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "我的第一次创业经历",
			Slug:        "my-first-startup-experience",
			Content:     "# 我的第一次创业经历\n\n创业是一条充满挑战但也很有趣的道路。\n\n## 为什么会创业\n\n毕业后，我在一家互联网公司工作了三年，逐渐意识到自己想要做更多的事情。\n\n## 遇到的困难\n\n1. 资金不足\n2. 团队组建\n3. 产品定位\n4. 市场推广\n\n## 收获与成长\n\n虽然项目最终没有成功，但这段经历让我学到了很多：\n\n- 技术能力提升\n- 团队协作经验\n- 解决问题的能力\n- 对商业的理解\n\n创业不仅仅是为了成功，更是一种成长和历练。",
			Excerpt:     "创业是一条充满挑战但也很有趣的道路，记录我的第一次创业经历。",
			Status:       "published",
			Views:        200,
			CategoryID:  &categories[1].ID,
			AuthorID:     admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "UI 设计基本原则",
			Slug:        "ui-design-principles",
			Content:     "# UI 设计基本原则\n\n好的 UI 设计让产品更加易用和美观。\n\n## 一致性\n\n保持界面元素的一致性是 UI 设计的基础。\n\n## 层次结构\n\n通过大小、颜色、对比度等手段建立清晰的视觉层次。\n\n## 反馈\n\n为用户的每个操作提供及时的反馈。\n\n## 可用性\n\n- 清晰的导航\n- 易于理解的图标\n- 适当的留白\n- 合理的配色\n\n## 总结\n\n好的设计是隐藏的设计，让用户专注于内容而不是界面本身。",
			Excerpt:     "好的 UI 设计让产品更加易用和美观，本文介绍基本原则。",
			Status:       "published",
			Views:        75,
			CategoryID:  &categories[2].ID,
			AuthorID:     admin.ID,
			PublishedAt: publishedAt,
		},
	}

	if err := DB.Create(&posts).Error; err != nil {
		log.Printf("Warning: Failed to create posts: %v", err)
	}

	// 5. 为文章添加标签关联
	// 第一篇文章关联 Go, Vue 标签
	if len(tags) >= 2 {
		DB.Model(&posts[0]).Association("Tags").Append(&tags[0], &tags[1])
	}
	// 第二篇关联 Go, Docker 标签
	if len(tags) >= 7 {
		DB.Model(&posts[1]).Association("Tags").Append(&tags[0], &tags[6])
	}
	// 第三篇关联 Vue, React 标签
	if len(tags) >= 3 {
		DB.Model(&posts[2]).Association("Tags").Append(&tags[1], &tags[2])
	}
	// 第四篇关联 Docker, Kubernetes 标签
	if len(tags) >= 8 {
		DB.Model(&posts[3]).Association("Tags").Append(&tags[6], &tags[7])
	}
	// 第五篇关联 AI, 机器学习 标签
	if len(tags) >= 10 {
		DB.Model(&posts[4]).Association("Tags").Append(&tags[8], &tags[9])
	}

	// 6. 创建页面
	pages := []Page{
		{
			Title:       "关于",
			Slug:        "about",
			Content:     "# 关于我们\n\n欢迎来到 Catter CMS！\n\n我们致力于打造一个简洁、强大的博客系统。\n\n## 我们的使命\n\n让每个人都能轻松创建自己的博客。\n\n## 联系方式\n\n- Email: contact@example.com\n- GitHub: github.com/catter-cms",
			Excerpt:     "了解 Catter CMS 团队和我们的使命",
			Status:      "published",
			Template:    "default",
			SortOrder:   1,
		},
		{
			Title:       "联系方式",
			Slug:        "contact",
			Content:     "# 联系我们\n\n有问题或建议？请联系我们！\n\n## 联系方式\n\n- **Email**: contact@example.com\n- **GitHub**: github.com/catter-cms\n- **Twitter**: @cattercms\n\n我们会尽快回复您的问题。",
			Excerpt:     "有问题或建议？请联系我们！",
			Status:      "published",
			Template:    "default",
			SortOrder:   2,
		},
		{
			Title:       "隐私政策",
			Slug:        "privacy",
			Content:     "# 隐私政策\n\n我们重视您的隐私。\n\n## 数据收集\n\n我们仅收集必要的数据来提供博客服务。\n\n## 数据使用\n\n您的数据仅用于提供服务，不会出售给第三方。\n\n## 联系我们\n\n如有问题，请联系我们的隐私团队。",
			Excerpt:     "了解我们的隐私政策",
			Status:      "published",
			Template:    "default",
			SortOrder:   3,
		},
	}

	if err := DB.Create(&pages).Error; err != nil {
		log.Printf("Warning: Failed to create pages: %v", err)
	}

	log.Println("Seed data created successfully!")
	log.Printf("  - Categories: %d", len(categories))
	log.Printf("  - Tags: %d", len(tags))
	log.Printf("  - Posts: %d", len(posts))
	log.Printf("  - Pages: %d", len(pages))

	return nil
}
