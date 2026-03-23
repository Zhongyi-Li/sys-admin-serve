package database

import (
	"context"
	"fmt"
	"time"

	"sys-admin-serve/internal/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const pingTimeout = 3 * time.Second
const maxPingAttempts = 10
const pingRetryInterval = 2 * time.Second

func NewMySQL(cfg config.MySQLConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get mysql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	var pingErr error
	for attempt := 1; attempt <= maxPingAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		pingErr = sqlDB.PingContext(ctx)
		cancel()
		if pingErr == nil {
			return db, nil
		}

		if attempt < maxPingAttempts {
			time.Sleep(pingRetryInterval)
		}
	}

	return nil, fmt.Errorf("ping mysql after retries: %w", pingErr)
}
