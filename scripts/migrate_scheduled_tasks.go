package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库连接
	dbPath := filepath.Join(".", "catter.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	defer db.Close()

	// 读取SQL文件
	sqlFile := "migrations/20260316000001_create_scheduled_tasks.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("读取SQL文件失败: %v", err)
	}

	// 执行SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("执行SQL失败: %v", err)
	}

	fmt.Println("✅ 定时任务表创建成功!")
}
