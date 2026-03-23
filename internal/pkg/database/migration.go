package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"sys-admin-serve/internal/pkg/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationTimeout = 5 * time.Second

func EnsureDatabaseExists(cfg config.MySQLConfig) error {
	adminDB, err := sql.Open("mysql", cfg.AdminDSN())
	if err != nil {
		return fmt.Errorf("open mysql admin connection: %w", err)
	}
	defer func() {
		_ = adminDB.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), migrationTimeout)
	defer cancel()

	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET %s COLLATE utf8mb4_unicode_ci",
		cfg.DBName,
		cfg.Charset,
	)

	if _, err := adminDB.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("create database %s: %w", cfg.DBName, err)
	}

	return nil
}

func NewMigrator(cfg config.MySQLConfig, migrationsDir string) (*migrate.Migrate, func() error, error) {
	if err := EnsureDatabaseExists(cfg); err != nil {
		return nil, nil, err
	}

	absoluteDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve migrations dir: %w", err)
	}

	sqlDB, err := sql.Open("mysql", migrationDSN(cfg))
	if err != nil {
		return nil, nil, fmt.Errorf("open mysql for migration: %w", err)
	}

	driver, err := mysqlmigrate.WithInstance(sqlDB, &mysqlmigrate.Config{DatabaseName: cfg.DBName})
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("create mysql migration driver: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", filepath.ToSlash(absoluteDir))
	migrator, err := migrate.NewWithDatabaseInstance(sourceURL, cfg.DBName, driver)
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("create migrator: %w", err)
	}

	cleanup := func() error {
		sourceErr, databaseErr := migrator.Close()
		if sourceErr != nil {
			return fmt.Errorf("close migration source: %w", sourceErr)
		}
		if databaseErr != nil {
			return fmt.Errorf("close migration database: %w", databaseErr)
		}
		return nil
	}

	return migrator, cleanup, nil
}

func migrationDSN(cfg config.MySQLConfig) string {
	return fmt.Sprintf("%s&multiStatements=true", cfg.DSN())
}
