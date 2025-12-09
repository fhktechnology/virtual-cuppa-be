package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "virtual_cuppa")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	migrationsPath := getMigrationsPath()
	if err := RunMigrations(sqlDB, migrationsPath); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	DB = database
}

func GetSQLDB() (*sql.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return DB.DB()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
