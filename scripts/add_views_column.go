package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "catter_cms_development.sqlite"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查 views 列是否存在
	var columnName string
	err = db.QueryRow("SELECT name FROM pragma_table_info('posts') WHERE name='views'").Scan(&columnName)
	if err == sql.ErrNoRows {
		// 列不存在，添加它
		fmt.Println("Adding 'views' column to posts table...")
		_, err = db.Exec("ALTER TABLE posts ADD COLUMN views INTEGER DEFAULT 0")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Column 'views' added successfully!")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Column 'views' already exists")
	}

	// 更新现有记录
	result, err := db.Exec("UPDATE posts SET views = 0 WHERE views IS NULL")
	if err != nil {
		log.Fatal(err)
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("Updated %d rows\n", rows)

	// 验证
	rows2, err := db.Query("SELECT id, title, views FROM posts LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	fmt.Println("\nSample posts:")
	for rows2.Next() {
		var id, views int
		var title string
		rows2.Scan(&id, &title, &views)
		fmt.Printf("  ID: %d, Title: %s, Views: %d\n", id, title, views)
	}
}
