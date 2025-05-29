package main

import (
	"log"

	"github.com/zgiai/ginext/config"
	"github.com/zgiai/ginext/internal/modules/user"
	"github.com/zgiai/ginext/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	migrator := db.Migrator()

	// 如果表存在，先删除
	if migrator.HasTable(&user.User{}) {
		if err := migrator.DropTable(&user.User{}); err != nil {
			log.Fatalf("Failed to drop users table: %v", err)
		}
		log.Println("Dropped existing users table")
	}

	// 创建新表
	if err := migrator.CreateTable(&user.User{}); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	log.Println("Database migration completed successfully")
}
