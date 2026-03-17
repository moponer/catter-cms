package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("=== Simple Diagnostic Test ===")

	db, err := sql.Open("sqlite3", "catter_cms_development.sqlite")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// 测试查询
	fmt.Println("Testing post query...")
	rows, err := db.Query("SELECT id, title, views FROM posts WHERE status = 'published'")
	if err != nil {
		log.Fatal("Failed to query posts:", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, views int
		var title string
		rows.Scan(&id, &title, &views)
		fmt.Printf("  Post %d: ID=%d, Title=%s, Views=%d\n", count, id, title, views)
		count++
	}

	fmt.Printf("\nFound %d posts\n", count)
	fmt.Println("\n=== Test completed! ===")
}
