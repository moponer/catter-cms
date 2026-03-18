package main

import (
	"fmt"
	"log"

	"catter-cms-gin/internal/models"
	"catter-cms-gin/config"
)

func main() {
	cfg := config.Get()
	dbConfig := &config.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Database: cfg.Database.Database,
	}
	if err := models.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	var user models.User
	if err := models.DB.Where("email = ?", "admin@example.com").First(&user).Error; err != nil {
		log.Fatalf("User not found: %v", err)
	}

	fmt.Printf("User found:\n")
	fmt.Printf("  ID: %d\n", user.ID)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Role: %s\n", user.Role)
	fmt.Printf("  PasswordHash: %s\n", user.PasswordHash)

	// Test password check
	valid := user.CheckPassword("admin123")
	fmt.Printf("\nPassword 'admin123' valid: %v\n", valid)
}
