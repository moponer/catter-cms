package main

import (
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/gobuffalo/plush/v4"
)

type Post struct {
	ID          int
	Title       string
	Slug        string
	Content     string
	Excerpt     string
	Status      string
	Views       int
	PublishedAt string
}

func main() {
	// 创建测试数据
	posts := []Post{
		{ID: 1, Title: "Test Post 1", Slug: "test-1", Excerpt: "Excerpt 1", Views: 10, PublishedAt: "2026-03-15"},
		{ID: 2, Title: "Test Post 2", Slug: "test-2", Excerpt: "Excerpt 2", Views: 20, PublishedAt: "2026-03-14"},
	}

	// 创建上下文
	ctx := plush.NewContext()
	ctx.Set("posts", posts)

	// 测试模板
	tmpl := `<%= for (post) in posts { %>
Title: <%= post.Title %>, Views: <%= post.Views %>
<% } %>`

	// 渲染模板
	result, err := plush.Render(tmpl, ctx)
	if err != nil {
		log.Fatal("Template error:", err)
	}

	fmt.Println("Result:")
	fmt.Println(result)

	// 测试更复杂的模板
	fmt.Println("\n=== Testing with struct slice ===")
	ctx2 := plush.NewContext()
	ctx2.Set("posts", posts)

	tmpl2 := `<%= for (post) in posts { %>
<%= post.Title %> - <%= post.Views %> views
<% } %>`

	result2, err := plush.Render(tmpl2, ctx2)
	if err != nil {
		log.Fatal("Template error 2:", err)
	}

	fmt.Println(result2)

	// 测试 helper 函数
	fmt.Println("\n=== Testing len() helper ===")
	ctx3 := plush.NewContext()
	ctx3.Set("posts", posts)
	
	// 添加 len 函数到上下文
	ctx3.Set("len", func(s interface{}) int {
		switch v := s.(type) {
		case []Post:
			return len(v)
		case string:
			return len(v)
		default:
			return 0
		}
	})

	tmpl3 := `Total posts: <%= len(posts) %>`
	result3, err := plush.Render(tmpl3, ctx3)
	if err != nil {
		log.Fatal("Template error 3:", err)
	}

	fmt.Println(result3)

	// 测试访问单个结构体的字段
	fmt.Println("\n=== Testing single struct field access ===")
	ctx4 := plush.NewContext()
	post := posts[0]
	ctx4.Set("post", post)

	tmpl4 := `Title: <%= post.Title %>, Views: <%= post.Views %>`
	result4, err := plush.Render(tmpl4, ctx4)
	if err != nil {
		log.Fatal("Template error 4:", err)
	}

	fmt.Println(result4)

	// 打印结构体信息
	fmt.Println("\n=== Struct info ===")
	fmt.Printf("Post type: %T\n", post)
	fmt.Printf("Post.Views type: %T\n", post.Views)
	
	// 测试 template 包来检查是否有问题
	fmt.Println("\n=== Testing with standard template package ===")
	t := template.Must(template.New("test").Parse(`{{range .Posts}}Title: {{.Title}}, Views: {{.Views}}
{{end}}`))
	var buf strings.Builder
	err = t.Execute(&buf, map[string]interface{}{"Posts": posts})
	if err != nil {
		log.Fatal("Standard template error:", err)
	}
	fmt.Println(buf.String())
}
