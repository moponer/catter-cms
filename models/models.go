package models

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		return
	}
	pop.Debug = env == "development"

	log.Println("数据库连接成功")
}
