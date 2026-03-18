package models

import (
	"log"

	"catter-cms-gin/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&Post{},
		&Category{},
		&Tag{},
		&Comment{},
		&Page{},
		&Media{},
		&Setting{},
		&ThemeConfig{},
		&Menu{},
		&MenuItem{},
	)
}
