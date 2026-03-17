package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/plush/v4"
)

func main() {
	// 读取模板文件
	content, err := os.ReadFile("templates/themes/catter-blog/pages/home.plush.html")
	if err != nil {
		log.Fatal("Error reading template:", err)
	}

	// 创建测试数据
	type Post struct {
		ID          int
		Title       string
		Slug        string
		Excerpt     string
		Views       int
		PublishedAt string
	}
	
	type Category struct {
		ID        int
		Name      string
		Slug      string
		PostCount int
	}
	
	type Tag struct {
		ID   int
		Name string
		Slug string
	}

	posts := []Post{
		{ID: 1, Title: "Test Post 1", Slug: "test-1", Excerpt: "Excerpt 1", Views: 10, PublishedAt: "2026-03-15"},
		{ID: 2, Title: "Test Post 2", Slug: "test-2", Excerpt: "Excerpt 2", Views: 20, PublishedAt: "2026-03-14"},
	}

	categories := []Category{
		{ID: 1, Name: "Tech", Slug: "tech", PostCount: 5},
		{ID: 2, Name: "Life", Slug: "life", PostCount: 3},
	}

	tags := []Tag{
		{ID: 1, Name: "Go", Slug: "go"},
		{ID: 2, Name: "Web", Slug: "web"},
	}

	site := map[string]string{
		"title":       "Test Site",
		"description": "Test Description",
	}

	// 创建上下文
	ctx := plush.NewContext()
	ctx.Set("posts", posts)
	ctx.Set("categories", categories)
	ctx.Set("tags", tags)
	ctx.Set("site", site)
	ctx.Set("theme", "catter-blog")

	// 添加 len 函数
	ctx.Set("len", func(s interface{}) int {
		switch v := s.(type) {
		case []Post:
			return len(v)
		case []Category:
			return len(v)
		case []Tag:
			return len(v)
		case string:
			return len(v)
		default:
			return 0
		}
	})

	// 渲染模板
	result, err := plush.Render(string(content), ctx)
	if err != nil {
		log.Fatal("Template error:", err)
	}

	fmt.Println("Template rendered successfully!")
	fmt.Println("Result length:", len(result))
	
	// 只打印前500个字符
	if len(result) > 500 {
		fmt.Println("First 500 chars:")
		fmt.Println(result[:500])
	} else {
		fmt.Println("Full result:")
		fmt.Println(result)
	}
}
