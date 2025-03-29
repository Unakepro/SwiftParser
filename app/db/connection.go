package db

import (
	"fmt"
	"log"
	"swiftapi/app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Database does not exist. Trying to create it...")
		createDatabase(cfg)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Database connection failed:", err)
		}
	}

	log.Println("Database connected successfully!")
	DB = db
}

func createDatabase(cfg *config.Config) {
	tmp_dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)

	tempDB, err := gorm.Open(postgres.Open(tmp_dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL server:", err)
	}

	tempDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
}

func SetDatabase(database *gorm.DB) {
	DB = database
}
