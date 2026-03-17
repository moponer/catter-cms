package grifts

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/grift/grift"
	"catter-cms/models"
)

var _ = grift.Namespace("db", func() {
	// 运行主题配置迁移
	grift.Desc("migrate_themes", "运行主题配置数据库迁移")
	grift.Add("migrate_themes", func(c *grift.Context) error {
		return migrateThemes()
	})

	// 回滚主题配置迁移
	grift.Desc("rollback_themes", "回滚主题配置数据库迁移")
	grift.Add("rollback_themes", func(c *grift.Context) error {
		return rollbackThemes()
	})

	// 初始化主题数据
	grift.Desc("seed_themes", "初始化主题数据")
	grift.Add("seed_themes", func(c *grift.Context) error {
		return seedThemes()
	})
})

// migrateThemes 运行主题配置迁移
func migrateThemes() error {
	log.Println("开始运行主题配置迁移...")

	// 读取并执行迁移SQL文件
	migrationFile := filepath.Join("migrations", "20260314120000_create_theme_configs.up.sql")
	
	if err := executeSQLFile(migrationFile); err != nil {
		return fmt.Errorf("执行迁移失败: %w", err)
	}

	log.Println("主题配置迁移完成")
	return nil
}

// rollbackThemes 回滚主题配置迁移
func rollbackThemes() error {
	log.Println("开始回滚主题配置迁移...")

	// 读取并执行回滚SQL文件
	rollbackFile := filepath.Join("migrations", "20260314120000_create_theme_configs.down.sql")
	
	if err := executeSQLFile(rollbackFile); err != nil {
		return fmt.Errorf("执行回滚失败: %w", err)
	}

	log.Println("主题配置回滚完成")
	return nil
}

// seedThemes 初始化主题数据
func seedThemes() error {
	log.Println("开始初始化主题数据...")

	// 检查是否已有主题数据
	var count int
	err := models.DB.RawQuery("SELECT COUNT(*) FROM themes").First(&count)
	if err != nil {
		return fmt.Errorf("查询主题数据失败: %w", err)
	}

	if count == 0 {
		log.Println("没有主题数据,开始初始化...")
		
		// 插入默认主题
		themes := []map[string]interface{}{
			{
				"name":        "Catter Blog",
				"slug":        "catter-blog",
				"version":     "2.0.0",
				"description": "现代化的博客主题，支持模块化配置",
				"author":      "Catter Team",
				"thumbnail":   "/themes/catter-blog/assets/images/thumbnail.jpg",
				"is_active":   true,
			},
			{
				"name":        "Catter CMS Official",
				"slug":        "catter-cms-official",
				"version":     "2.0.0",
				"description": "Catter CMS官方官网主题",
				"author":      "Catter Team",
				"thumbnail":   "/themes/catter-cms-official/assets/images/thumbnail.jpg",
				"is_active":   false,
			},
		}

		for _, theme := range themes {
			query := `INSERT INTO themes (name, slug, version, description, author, thumbnail, is_active, created_at, updated_at)
					 VALUES (:name, :slug, :version, :description, :author, :thumbnail, :is_active, datetime('now'), datetime('now'))`
			
			if err := models.DB.RawQuery(query, theme).Exec(); err != nil {
				log.Printf("插入主题 %s 失败: %v", theme["name"], err)
			} else {
				log.Printf("已插入主题: %s", theme["name"])
			}
		}
	} else {
		log.Printf("已有 %d 个主题,跳过初始化", count)
	}

	log.Println("主题数据初始化完成")
	return nil
}

// executeSQLFile 执行SQL文件
func executeSQLFile(sqlFile string) error {
	// 检查文件是否存在
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return fmt.Errorf("SQL文件不存在: %s", sqlFile)
	}

	// 读取SQL文件内容
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %w", err)
	}

	sql := string(content)

	// 分割SQL语句(按分号分隔)
	statements := strings.Split(sql, ";")

	// 执行每条SQL语句
	for _, statement := range statements {
		stmt := strings.TrimSpace(statement)
		if stmt == "" {
			continue
		}

		// 跳过注释
		if strings.HasPrefix(stmt, "--") || strings.HasPrefix(stmt, "/*") {
			continue
		}

		// 执行SQL
		if err := models.DB.RawQuery(stmt).Exec(); err != nil {
			// 忽略"表已存在"的错误
			if !strings.Contains(err.Error(), "already exists") {
				log.Printf("执行SQL语句失败: %v\n语句: %s", err, stmt[:min(100, len(stmt))])
				// 不返回错误,继续执行下一条
			}
		}
	}

	return nil
}

// min 返回最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
