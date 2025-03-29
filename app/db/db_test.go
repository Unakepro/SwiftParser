package db_test

import (
	"log"
	"testing"

	"main_pack/db"
	"main_pack/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	testDB.AutoMigrate(&models.Country{}, &models.SwiftCode{})

	return testDB
}

func TestDatabaseConnection(t *testing.T) {
	testDB := SetupTestDB()
	db.SetDatabase(testDB)

	if db.DB == nil {
		t.Fatal("Database should be initialized, but it's nil")
	}
}
