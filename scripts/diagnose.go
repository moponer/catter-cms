package main

import (
	"fmt"
	"log"

	"catter-cms/actions"
	"catter-cms/models"
	"github.com/gobuffalo/pop/v6"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("=== Diagnostic Test ===")
	
	// 测试数据库连接
	fmt.Println("Testing database connection...")
	err := pop.LoadConfigFile("database.yml")
	if err != nil {
		log.Fatal("Failed to load database config:", err)
	}
	
	tx, err := pop.Connect("development")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	models.DB = tx
	fmt.Println("Database connected successfully!")
	
	// 测试查询
	fmt.Println("Testing post query...")
	posts := models.Posts{}
	query := models.DB.Where("status = ?", "published").Order("published_at desc").Limit(10)
	if err := query.All(&posts); err != nil {
		log.Fatal("Failed to query posts:", err)
	}
	fmt.Printf("Found %d posts\n", len(posts))
	
	// 检查 App 是否初始化
	fmt.Println("Checking Buffalo App...")
	if actions.App == nil {
		log.Fatal("App is nil!")
	}
	fmt.Println("App initialized successfully!")
	
	fmt.Println("\n=== All tests passed! ===")
}
