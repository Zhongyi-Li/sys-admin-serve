package main

import (
	"fmt"
	"log"
	"os"

	"sys-admin-serve/internal/pkg/config"
	"sys-admin-serve/internal/pkg/database"
)

func main() {
	configPath := os.Getenv("APP_CONFIG")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get mysql db: %v", err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	if err := database.SeedInitialData(db); err != nil {
		log.Fatalf("seed initial data: %v", err)
	}

	adminUsername := os.Getenv("SEED_ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}

	fmt.Printf("seeded admin user %q\n", adminUsername)
}
