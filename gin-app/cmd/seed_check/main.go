package main

import (
	"fmt"
	"log"

	"catter-cms-gin/internal/models"
	"catter-cms-gin/config"
)

func main() {
	// 加载配置
	cfg := config.Get()

	// 初始化数据库
	dbConfig := &config.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Database: cfg.Database.Database,
	}
	if err := models.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 查询数据
	var users []models.User
	models.DB.Find(&users)
	fmt.Printf("Users: %d\n", len(users))
	for _, u := range users {
		fmt.Printf("  - %s (%s) - Role: %s\n", u.Name, u.Email, u.Role)
	}

	var categories []models.Category
	models.DB.Find(&categories)
	fmt.Printf("\nCategories: %d\n", len(categories))
	for _, c := range categories {
		fmt.Printf("  - %s (%s)\n", c.Name, c.Slug)
	}

	var tags []models.Tag
	models.DB.Find(&tags)
	fmt.Printf("\nTags: %d\n", len(tags))
	for _, t := range tags {
		fmt.Printf("  - %s (%s)\n", t.Name, t.Slug)
	}

	var posts []models.Post
	models.DB.Find(&posts)
	fmt.Printf("\nPosts: %d\n", len(posts))
	for _, p := range posts {
		fmt.Printf("  - %s (%s) - Status: %s\n", p.Title, p.Slug, p.Status)
	}

	var pages []models.Page
	models.DB.Find(&pages)
	fmt.Printf("\nPages: %d\n", len(pages))
	for _, p := range pages {
		fmt.Printf("  - %s (%s) - Status: %s\n", p.Title, p.Slug, p.Status)
	}
}
