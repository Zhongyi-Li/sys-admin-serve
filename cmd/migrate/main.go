package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"sys-admin-serve/internal/pkg/config"
	"sys-admin-serve/internal/pkg/database"

	"github.com/golang-migrate/migrate/v4"
)

const defaultMigrationsDir = "migrations"

func main() {
	configPath := os.Getenv("APP_CONFIG")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	migrator, cleanup, err := database.NewMigrator(cfg.MySQL, defaultMigrationsDir)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer func() {
		if closeErr := cleanup(); closeErr != nil {
			log.Printf("close migrator: %v", closeErr)
		}
	}()

	if err := runCommand(migrator, command, os.Args[2:]); err != nil {
		log.Fatalf("run migration command %s: %v", command, err)
	}
}

func runCommand(migrator *migrate.Migrate, command string, args []string) error {
	switch command {
	case "up":
		return ignoreNoChange(migrator.Up())
	case "down":
		steps := 1
		if len(args) > 0 {
			value, err := strconv.Atoi(args[0])
			if err != nil || value <= 0 {
				return fmt.Errorf("invalid down steps %q", args[0])
			}
			steps = value
		}
		return ignoreNoChange(migrator.Steps(-steps))
	case "version":
		version, dirty, err := migrator.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("version: none")
			return nil
		}
		if err != nil {
			return fmt.Errorf("get migration version: %w", err)
		}
		fmt.Printf("version: %d dirty: %t\n", version, dirty)
		return nil
	case "force":
		if len(args) == 0 {
			return errors.New("force requires a version")
		}
		version, err := strconv.Atoi(args[0])
		if err != nil || version < 0 {
			return fmt.Errorf("invalid force version %q", args[0])
		}
		return migrator.Force(version)
	default:
		return fmt.Errorf("unsupported command %q", command)
	}
}

func ignoreNoChange(err error) error {
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("no change")
		return nil
	}

	return err
}
