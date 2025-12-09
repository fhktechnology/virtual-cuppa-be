package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"virtual-cuppa-be/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	migrateUp := migrateCmd.Bool("up", false, "Run all pending migrations")
	migrateDown := migrateCmd.Int("down", 0, "Rollback N migrations")
	migrateVersion := migrateCmd.Bool("version", false, "Show current migration version")

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [command] [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  migrate -up           Run all pending migrations")
		fmt.Println("  migrate -down N       Rollback N migrations")
		fmt.Println("  migrate -version      Show current migration version")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate":
		migrateCmd.Parse(os.Args[2:])

		config.ConnectDatabase()
		sqlDB, err := config.GetSQLDB()
		if err != nil {
			log.Fatal("Failed to get database:", err)
		}

		migrationsPath := "./migrations"

		if *migrateUp {
			if err := config.RunMigrations(sqlDB, migrationsPath); err != nil {
				log.Fatal("Migration failed:", err)
			}
			fmt.Println("Migrations completed successfully")
		} else if *migrateDown > 0 {
			if err := config.RollbackMigration(sqlDB, migrationsPath, *migrateDown); err != nil {
				log.Fatal("Rollback failed:", err)
			}
			fmt.Printf("Rolled back %d migration(s)\n", *migrateDown)
		} else if *migrateVersion {
			version, dirty, err := config.GetMigrationVersion(sqlDB, migrationsPath)
			if err != nil {
				log.Fatal("Failed to get version:", err)
			}
			fmt.Printf("Current version: %d (dirty: %v)\n", version, dirty)
		} else {
			migrateCmd.Usage()
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
