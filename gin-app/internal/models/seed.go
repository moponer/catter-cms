package models

import (
	"fmt"
	"log"
	"time"
)

// SeedData 初始化示例数据
func SeedData() error {
	// 1. 确保有管理员用户
	if err := CreateDefaultAdmin(); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	// 2. 确保创建基础页面（不管是否有文章数据）
	ensurePages()

	// 3. 确保创建基础菜单
	ensureMenus()

	// 检查是否已有分类数据
	var postCount int64
	DB.Model(&Post{}).Count(&postCount)
	if postCount > 0 {
		log.Println("Posts already exist, skipping post seed")
		return nil
	}

	log.Println("Seeding initial data...")

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
			Status:      "published",
			Views:       100,
			CategoryID:  &categories[0].ID,
			AuthorID:    admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "Go 语言快速入门指南",
			Slug:        "go-quick-start-guide",
			Content:     "# Go 语言快速入门\n\nGo 语言是一种简洁、高效、并发友好的编程语言。\n\n## 环境安装\n\n```bash\n# 下载 Go\nwget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz\n\n# 解压\nsudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz\n\n# 配置环境变量\nexport PATH=$PATH:/usr/local/go/bin\n```\n\n## 第一个程序\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```\n\n## 基本语法\n\n- 变量声明\n- 函数定义\n- 流程控制\n- 并发编程\n\nGo 语言的并发模型非常强大，使用 goroutine 和 channel 可以轻松实现并发编程。",
			Excerpt:     "Go 语言是一种简洁、高效、并发友好的编程语言，本文将带你快速入门。",
			Status:      "published",
			Views:       85,
			CategoryID:  &categories[4].ID,
			AuthorID:    admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "Vue 3 组合式 API 完全指南",
			Slug:        "vue3-composition-api-guide",
			Content:     "# Vue 3 组合式 API 完全指南\n\nVue 3 引入了组合式 API，这是 Vue 的一次重大升级。\n\n## 为什么选择组合式 API\n\n1. 更好的代码组织\n2. 更好的逻辑复用\n3. 更好的 TypeScript 支持\n\n## 基本用法\n\n```vue\n<script setup>\nimport { ref, computed, onMounted } from 'vue'\n\nconst count = ref(0)\nconst doubled = computed(() => count.value * 2)\n\nfunction increment() {\n  count.value++\n}\n\nonMounted(() => {\n  console.log('Component mounted!')\n})\n</script>\n\n<template>\n  <button @click=\"increment\">\n    Count: {{ count }}, Doubled: {{ doubled }}\n  </button>\n</template>\n```\n\n组合式 API 让组件逻辑更加清晰，是现代 Vue 开发的推荐方式。",
			Excerpt:     "Vue 3 引入了组合式 API，这是 Vue 的一次重大升级本文将详细介绍。",
			Status:      "published",
			Views:       120,
			CategoryID:  &categories[3].ID,
			AuthorID:    admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "Docker 入门实战教程",
			Slug:        "docker-getting-started",
			Content:     "# Docker 入门实战\n\nDocker 是现代软件开发和部署的重要工具。\n\n## 什么是 Docker\n\nDocker 是一个开源的容器化平台，可以将应用及其依赖打包成轻量级的容器。\n\n## 基本命令\n\n```bash\n# 拉取镜像\ndocker pull nginx\n\n# 运行容器\ndocker run -d -p 80:80 nginx\n\n# 查看容器\ndocker ps\n\n# 停止容器\ndocker stop <container_id>\n```\n\n## Dockerfile 示例\n\n```dockerfile\nFROM node:18-alpine\nWORKDIR /app\nCOPY package*.json ./\nRUN npm install\nCOPY . .\nEXPOSE 3000\nCMD [\"npm\", \"start\"]\n```\n\nDocker 让开发、测试、部署变得更加简单和一致。",
			Excerpt:     "Docker 是现代软件开发和部署的重要工具，本文将带你快速上手。",
			Status:      "published",
			Views:       90,
			CategoryID:  &categories[0].ID,
			AuthorID:    admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "我的第一次创业经历",
			Slug:        "my-first-startup-experience",
			Content:     "# 我的第一次创业经历\n\n创业是一条充满挑战但也很有趣的道路。\n\n## 为什么会创业\n\n毕业后，我在一家互联网公司工作了三年，逐渐意识到自己想要做更多的事情。\n\n## 遇到的困难\n\n1. 资金不足\n2. 团队组建\n3. 产品定位\n4. 市场推广\n\n## 收获与成长\n\n虽然项目最终没有成功，但这段经历让我学到了很多：\n\n- 技术能力提升\n- 团队协作经验\n- 解决问题的能力\n- 对商业的理解\n\n创业不仅仅是为了成功，更是一种成长和历练。",
			Excerpt:     "创业是一条充满挑战但也很有趣的道路，记录我的第一次创业经历。",
			Status:      "published",
			Views:       200,
			CategoryID:  &categories[1].ID,
			AuthorID:    admin.ID,
			PublishedAt: publishedAt,
		},
		{
			Title:       "UI 设计基本原则",
			Slug:        "ui-design-principles",
			Content:     "# UI 设计基本原则\n\n好的 UI 设计让产品更加易用和美观。\n\n## 一致性\n\n保持界面元素的一致性是 UI 设计的基础。\n\n## 层次结构\n\n通过大小、颜色、对比度等手段建立清晰的视觉层次。\n\n## 反馈\n\n为用户的每个操作提供及时的反馈。\n\n## 可用性\n\n- 清晰的导航\n- 易于理解的图标\n- 适当的留白\n- 合理的配色\n\n## 总结\n\n好的设计是隐藏的设计，让用户专注于内容而不是界面本身。",
			Excerpt:     "好的 UI 设计让产品更加易用和美观，本文介绍基本原则。",
			Status:      "published",
			Views:       75,
			CategoryID:  &categories[2].ID,
			AuthorID:    admin.ID,
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

	log.Println("Seed data created successfully!")
	log.Printf("  - Categories: %d", len(categories))
	log.Printf("  - Tags: %d", len(tags))
	log.Printf("  - Posts: %d", len(posts))

	return nil
}

// ensurePages 确保基础页面存在
func ensurePages() {
	// 需要确保存在的基础页面
	pages := []struct {
		Slug      string
		Title     string
		Content   string
		Excerpt   string
		SortOrder int
	}{
		{
			Slug:      "about",
			Title:     "关于我们",
			Content:   "# 关于我们\n\n欢迎来到我们的网站！\n\n我们致力于提供最好的产品和服务。\n\n## 我们的使命\n\n帮助企业实现数字化转型，提升业务效率。\n\n## 联系方式\n\n- Email: contact@example.com\n- 电话: 400-123-4567",
			Excerpt:   "了解我们的团队和使命",
			SortOrder: 1,
		},
		{
			Slug:      "contact",
			Title:     "联系我们",
			Content:   "# 联系我们\n\n有任何问题或建议？请随时联系我们！\n\n## 联系方式\n\n- **地址**: 北京市朝阳区科技园A座\n- **电话**: 400-123-4567\n- **邮箱**: contact@example.com\n- **工作时间**: 周一至周五 9:00-18:00\n\n我们会在24小时内回复您的留言。",
			Excerpt:   "获取联系方式",
			SortOrder: 2,
		},
		{
			Slug:      "privacy",
			Title:     "隐私政策",
			Content:   "# 隐私政策\n\n我们重视您的隐私保护。\n\n## 信息收集\n\n我们收集您在使用服务时主动提供的信息。\n\n## 信息使用\n\n您的信息仅用于改善服务质量，不会与第三方共享。\n\n## 联系隐私团队\n\n如有任何隐私相关问题，请联系 privacy@example.com",
			Excerpt:   "了解我们的隐私政策",
			SortOrder: 3,
		},
		{
			Slug:      "features",
			Title:     "产品功能",
			Content:   "# 产品功能\n\n## 核心特性\n\n### 1. 简洁易用\n\n直观的界面设计，让内容管理变得轻松简单。\n\n### 2. 主题系统\n\n支持多主题切换，随时更换网站外观。\n\n### 3. Markdown支持\n\n内置 Markdown 编辑器，写作更加高效。\n\n### 4. SEO优化\n\n内置 SEO 功能，让搜索引擎更喜欢你的网站。\n\n### 5. 响应式设计\n\n完美适配各种设备，提供最佳用户体验。\n\n### 6. 插件扩展\n\n支持插件扩展，满足个性化需求。",
			Excerpt:   "了解产品的核心功能",
			SortOrder: 4,
		},
		{
			Slug:      "pricing",
			Title:     "价格方案",
			Content:   "# 价格方案\n\n选择适合您的方案\n\n## 免费版\n\n- 最多 5 篇文章\n- 1 个主题\n- 基础统计\n- 社区支持\n\n## 专业版\n\n- 无限文章\n- 所有主题\n- 高级统计\n- 优先支持\n- 自定义域名\n\n## 企业版\n\n- 私有部署\n- 专属客服\n- 定制开发\n- SLA 保障",
			Excerpt:   "查看价格方案",
			SortOrder: 5,
		},
		{
			Slug:      "services",
			Title:     "服务项目",
			Content:   "# 服务项目\n\n## 我们的服务\n\n### 企业数字化转型\n\n帮助企业实现业务流程数字化，提升运营效率。\n\n### 技术咨询\n\n提供专业的技术咨询服务，助力技术决策。\n\n### 定制开发\n\n根据企业需求，提供定制化的软件开发服务。\n\n### 运维支持\n\n7x24 小时运维支持，保障系统稳定运行。\n\n### 培训服务\n\n提供产品使用培训和技能提升课程。",
			Excerpt:   "了解我们的专业服务",
			SortOrder: 6,
		},
		{
			Slug:      "solutions",
			Title:     "解决方案",
			Content:   "# 行业解决方案\n\n## 针对不同行业的解决方案\n\n### 制造业\n\n- 生产管理数字化\n- 供应链协同\n- 质量追溯\n\n### 零售业\n\n- 全渠道营销\n- 会员管理\n- 库存优化\n\n### 服务业\n\n- 客户服务管理\n- 预约调度\n- 绩效分析\n\n### 教育行业\n\n- 在线学习平台\n- 学情分析\n- 资源管理",
			Excerpt:   "查看行业解决方案",
			SortOrder: 7,
		},
	}

	for _, page := range pages {
		var existing Page
		if err := DB.Where("slug = ?", page.Slug).First(&existing).Error; err != nil {
			// 页面不存在，创建它
			newPage := Page{
				Title:     page.Title,
				Slug:      page.Slug,
				Content:   page.Content,
				Excerpt:   page.Excerpt,
				Status:    "published",
				Template:  "default",
				SortOrder: page.SortOrder,
			}
			if err := DB.Create(&newPage).Error; err != nil {
				log.Printf("Warning: Failed to create page %s: %v", page.Slug, err)
			} else {
				log.Printf("Created page: %s", page.Slug)
			}
		}
	}
}

// ensureMenus 确保基础菜单存在
func ensureMenus() {
	// 需要确保存在的基础菜单
	menus := []struct {
		Location    string
		Name        string
		Description string
		Items       []struct {
			Label string
			URL   string
			Sort  int
		}
	}{
		{
			Location:    "header",
			Name:        "顶部导航",
			Description: "网站顶部主导航菜单",
			Items: []struct {
				Label string
				URL   string
				Sort  int
			}{
				{Label: "首页", URL: "/", Sort: 1},
				{Label: "文章", URL: "/article", Sort: 2},
				{Label: "分类", URL: "/categories", Sort: 3},
				{Label: "标签", URL: "/tags", Sort: 4},
				{Label: "归档", URL: "/archives", Sort: 5},
				{Label: "关于", URL: "/page/about", Sort: 6},
			},
		},
		{
			Location:    "footer",
			Name:        "底部导航",
			Description: "网站底部导航菜单",
			Items: []struct {
				Label string
				URL   string
				Sort  int
			}{
				{Label: "首页", URL: "/", Sort: 1},
				{Label: "文章", URL: "/article", Sort: 2},
				{Label: "关于", URL: "/page/about", Sort: 3},
				{Label: "联系方式", URL: "/page/contact", Sort: 4},
				{Label: "隐私政策", URL: "/page/privacy", Sort: 5},
			},
		},
	}

	for _, menu := range menus {
		var existing Menu
		if err := DB.Where("location = ?", menu.Location).First(&existing).Error; err != nil {
			// 菜单不存在，创建它
			newMenu := Menu{
				ID:          menu.Location,
				Name:        menu.Name,
				Location:    menu.Location,
				Description: menu.Description,
			}
			if err := DB.Create(&newMenu).Error; err != nil {
				log.Printf("Warning: Failed to create menu %s: %v", menu.Location, err)
				continue
			}

			// 创建菜单项
			for i, item := range menu.Items {
				menuItem := MenuItem{
					ID:      fmt.Sprintf("%s-item-%d", menu.Location, i+1),
					MenuID:  menu.Location,
					Label:   item.Label,
					URL:     item.URL,
					Sort:    item.Sort,
					Enabled: true,
				}
				if err := DB.Create(&menuItem).Error; err != nil {
					log.Printf("Warning: Failed to create menu item %s: %v", item.Label, err)
				}
			}

			log.Printf("Created menu: %s", menu.Location)
		}
	}
}
