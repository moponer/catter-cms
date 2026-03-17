package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库连接
	db, err := sql.Open("sqlite3", "./catter.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查表是否已存在
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='comments'").Scan(&tableName)
	if err == nil {
		fmt.Println("表 'comments' 已存在,跳过创建")
		return
	}

	// 创建comments表
	createTableSQL := `
	CREATE TABLE comments (
		id TEXT PRIMARY KEY,
		post_id TEXT NOT NULL,
		user_id TEXT,
		parent_id TEXT,
		author_name VARCHAR(255),
		author_email VARCHAR(255),
		author_url VARCHAR(500),
		content TEXT NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		ip_address VARCHAR(50),
		user_agent TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// 创建索引
	indexes := []string{
		"CREATE INDEX idx_comments_post_id ON comments(post_id);",
		"CREATE INDEX idx_comments_user_id ON comments(user_id);",
		"CREATE INDEX idx_comments_parent_id ON comments(parent_id);",
		"CREATE INDEX idx_comments_status ON comments(status);",
		"CREATE INDEX idx_comments_created_at ON comments(created_at);",
	}

	for _, idxSQL := range indexes {
		_, err = db.Exec(idxSQL)
		if err != nil {
			log.Printf("创建索引失败: %v", err)
		}
	}

	fmt.Println("✓ 表 'comments' 创建成功!")
}
