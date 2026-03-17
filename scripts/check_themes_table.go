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

	// 检查 themes 表结构
	fmt.Println("=== Themes Table Schema ===")
	rows, err := db.Query("PRAGMA table_info(themes)")
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
		fmt.Printf("  Column: %s, Type: %s\n", name, dtype)
	}

	// 检查 themes 数据
	fmt.Println("\n=== Themes Data ===")
	rows2, err := db.Query("SELECT * FROM themes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	cols, _ := rows2.Columns()
	fmt.Println("  Columns:", cols)

	for rows2.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows2.Scan(valuePtrs...)
		fmt.Println("  Row:", values)
	}
}
