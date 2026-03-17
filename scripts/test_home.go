package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID          int
	Title       string
	Slug        string
	Content     string
	Excerpt     string
	Status      string
	UserID      int
	CategoryID  *int
	Views       int
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func main() {
	db, err := sql.Open("sqlite3", "catter_cms_development.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 查询文章（模拟 home.go 中的查询）
	fmt.Println("=== Querying Posts ===")
	rows, err := db.Query("SELECT id, title, slug, content, excerpt, status, user_id, category_id, views, published_at, created_at, updated_at FROM posts WHERE status = ? ORDER BY published_at DESC LIMIT 10", "published")
	if err != nil {
		log.Fatal("Error querying posts:", err)
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		var publishedAt sql.NullTime
		err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Excerpt, &p.Status, &p.UserID, &p.CategoryID, &p.Views, &publishedAt, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			log.Fatal("Error scanning post:", err)
		}
		if publishedAt.Valid {
			p.PublishedAt = &publishedAt.Time
		}
		posts = append(posts, p)
		fmt.Printf("  Post: ID=%d, Title=%s, Views=%d\n", p.ID, p.Title, p.Views)
	}

	fmt.Printf("\nTotal posts: %d\n", len(posts))

	// 测试模板渲染
	fmt.Println("\n=== Testing Template Rendering ===")
	
	// 使用 plush 渲染
	/*
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
	*/
}
