package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gondolia/gondolia/services/identity/internal/config"
)

func main() {
	var command string
	var steps int

	flag.StringVar(&command, "command", "up", "Migration command: up, down, version, force, seed")
	flag.IntVar(&steps, "steps", 0, "Number of migrations to run (0 = all)")
	flag.Parse()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	switch command {
	case "up":
		if err := runMigrations(cfg.DatabaseURL(), steps, true); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		if steps == 0 {
			steps = 1 // Default to 1 step down for safety
		}
		if err := runMigrations(cfg.DatabaseURL(), steps, false); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")

	case "version":
		if err := showVersion(cfg.DatabaseURL()); err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}

	case "force":
		version := steps
		if version == 0 {
			log.Fatal("Version required for force command: -steps=<version>")
		}
		if err := forceVersion(cfg.DatabaseURL(), version); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}
		fmt.Printf("Forced version to %d\n", version)

	case "seed":
		if err := runSeed(cfg); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		fmt.Println("Seed completed successfully")

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

func loadConfig() (*config.Config, error) {
	// For migration, we don't require JWT secrets
	os.Setenv("JWT_ACCESS_SECRET", "not-needed-for-migration")
	os.Setenv("JWT_REFRESH_SECRET", "not-needed-for-migration")
	return config.Load()
}

func runMigrations(databaseURL string, steps int, up bool) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}
	defer m.Close()

	if up {
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	} else {
		err = m.Steps(-steps)
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func showVersion(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations applied yet")
			return nil
		}
		return err
	}

	fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
	return nil
}

func forceVersion(databaseURL string, version int) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}
	defer m.Close()

	return m.Force(version)
}

func runSeed(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Import here to avoid circular dependencies
	return seedDatabase(ctx, cfg.DatabaseURL())
}
