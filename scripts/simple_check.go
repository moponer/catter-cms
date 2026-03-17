package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

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
}

func main() {
	db, err := sql.Open("sqlite3", "catter_cms_development.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 查询文章
	rows, err := db.Query("SELECT id, title, slug, content, excerpt, status, user_id, category_id, views FROM posts WHERE status = 'published'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Excerpt, &p.Status, &p.UserID, &p.CategoryID, &p.Views)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, p)
	}

	fmt.Printf("Found %d posts\n", len(posts))
	for i, p := range posts {
		fmt.Printf("Post %d: ID=%d, Title=%s, Views=%d\n", i, p.ID, p.Title, p.Views)
		
		// 使用反射检查字段
		v := reflect.ValueOf(p)
		t := v.Type()
		
		// 尝试获取 Views 字段
		viewsField := v.FieldByName("Views")
		if viewsField.IsValid() {
			fmt.Printf("  Views field: value=%v, kind=%s, type=%s\n", viewsField.Interface(), viewsField.Kind(), viewsField.Type())
		} else {
			fmt.Println("  Views field not found!")
		}
		
		// 检查所有字段
		fmt.Printf("  All fields (%d):\n", t.NumField())
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			fieldValue := v.Field(j)
			fmt.Printf("    %s: %v (type: %s, kind: %s)\n", field.Name, fieldValue.Interface(), fieldValue.Type(), fieldValue.Kind())
		}
	}
}
