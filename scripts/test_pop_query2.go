package main

import (
	"fmt"
	"log"

	"catter-cms/models"

	"github.com/gobuffalo/plush/v4"
	"github.com/gobuffalo/pop/v6"
)

func main() {
	// 初始化数据库连接
	err := pop.LoadConfigFile("database.yml")
	if err != nil {
		log.Fatal("Error loading database config:", err)
	}

	tx, err := pop.Connect("development")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	models.DB = tx

	// 查询文章
	posts := models.Posts{}
	query := models.DB.Where("status = ?", "published").Order("published_at desc").Limit(10)
	if err := query.All(&posts); err != nil {
		log.Fatal("Error querying posts:", err)
	}

	fmt.Printf("Found %d posts:\n", len(posts))
	for i, p := range posts {
		fmt.Printf("  Post %d: ID=%d, Title=%s, Views=%d\n", i, p.ID, p.Title, p.Views)
	}

	// 测试模板渲染
	fmt.Println("\n=== Testing Plush Template ===")
	ctx := plush.NewContext()
	ctx.Set("posts", posts)

	tmpl := `<%= for (post) in posts { %>
Title: <%= post.Title %>, Views: <%= post.Views %>
<% } %>`

	result, err := plush.Render(tmpl, ctx)
	if err != nil {
		log.Fatal("Template error:", err)
	}
	fmt.Println("Result:")
	fmt.Println(result)
}
