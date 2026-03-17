package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "catter_cms_development.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查激活的主题
	fmt.Println("=== Checking Active Theme ===")
	var id, name, slug string
	var isActive bool
	err = db.QueryRow("SELECT id, name, slug, is_active FROM themes WHERE is_active = 1").Scan(&id, &name, &slug, &isActive)
	if err != nil {
		log.Fatal("Error querying active theme:", err)
	}
	fmt.Printf("Active Theme: ID=%s, Name=%s, Slug=%s\n", id, name, slug)

	// 检查主题目录是否存在
	themePath := fmt.Sprintf("templates/themes/%s/pages/home.plush.html", slug)
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		fmt.Printf("ERROR: Template file not found: %s\n", themePath)
	} else {
		fmt.Printf("OK: Template file exists: %s\n", themePath)
	}

	// 检查 posts 数据
	fmt.Println("\n=== Checking Posts ===")
	rows, err := db.Query("SELECT id, title, views FROM posts WHERE status = 'published' ORDER BY published_at DESC LIMIT 5")
	if err != nil {
		log.Fatal("Error querying posts:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, views int
		var title string
		rows.Scan(&id, &title, &views)
		fmt.Printf("  Post: ID=%d, Title=%s, Views=%d\n", id, title, views)
	}

	// 检查 categories 数据
	fmt.Println("\n=== Checking Categories ===")
	rows2, err := db.Query("SELECT id, name, slug FROM categories")
	if err != nil {
		log.Fatal("Error querying categories:", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var id int
		var name, slug string
		rows2.Scan(&id, &name, &slug)
		fmt.Printf("  Category: ID=%d, Name=%s, Slug=%s\n", id, name, slug)
	}
}
