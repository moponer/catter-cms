package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "catter_cms_development.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查 posts 表结构
	fmt.Println("=== Posts Table Schema ===")
	rows, err := db.Query("PRAGMA table_info(posts)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, dtype string
		var notNull, pk int
		var dfltValue sql.NullString
		rows.Scan(&cid, &name, &dtype, &notNull, &pk, &dfltValue)
		fmt.Printf("  Column: %s, Type: %s, Default: %v\n", name, dtype, dfltValue.String)
	}

	// 检查 posts 数据
	fmt.Println("\n=== Posts Data ===")
	rows2, err := db.Query("SELECT id, title, slug, status, views FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var id, views int
		var title, slug, status string
		rows2.Scan(&id, &title, &slug, &status, &views)
		fmt.Printf("  ID: %d, Title: %s, Slug: %s, Status: %s, Views: %d\n", id, title, slug, status, views)
	}

	// 检查 themes 表和激活的主题
	fmt.Println("\n=== Active Theme ===")
	var themeName string
	err = db.QueryRow("SELECT name FROM themes WHERE is_active = 1").Scan(&themeName)
	if err != nil {
		fmt.Println("  No active theme found or error:", err)
	} else {
		fmt.Println("  Active theme:", themeName)
	}
}
