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

	// 检查 themes 表
	fmt.Println("=== Themes Table ===")
	rows, err := db.Query("SELECT id, name, slug, is_active FROM themes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, slug string
		var isActive bool
		rows.Scan(&id, &name, &slug, &isActive)
		fmt.Printf("  ID: %d, Name: %s, Slug: %s, Active: %v\n", id, name, slug, isActive)
	}
}
