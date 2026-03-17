package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库连接
	db, err := sql.Open("sqlite3", "catter_cms_development.sqlite")
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	// 读取迁移文件
	sqlBytes, err := ioutil.ReadFile("migrations/20260312_create_users_table.sql")
	if err != nil {
		log.Fatalf("无法读取迁移文件: %v", err)
	}
	sqlScript := string(sqlBytes)

	// 执行SQL脚本
	_, err = db.Exec(sqlScript)
	if err != nil {
		log.Fatalf("执行SQL脚本失败: %v", err)
	}

	fmt.Println("数据库迁移成功完成！")
	fmt.Println("默认管理员用户已创建:")
	fmt.Println("邮箱: admin@catter-cms.com")
	fmt.Println("密码: admin123")
}
