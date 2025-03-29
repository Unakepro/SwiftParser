package services_test

import (
	"testing"

	"swiftapi/app/db"
	"swiftapi/app/models"
	"swiftapi/app/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := testDB.AutoMigrate(&models.Country{}, &models.SwiftCode{}); err != nil {
		t.Fatalf("Failed to migrate models: %v", err)
	}
	return testDB
}

func TestGetSwiftCodeByCode(t *testing.T) {
	dbInstance := setupTestDB(t)
	db.SetDatabase(dbInstance)

	country := models.Country{
		ISO2Code: "US",
		Name:     "United States",
		TimeZone: "America/New_York",
	}
	if err := dbInstance.Create(&country).Error; err != nil {
		t.Fatalf("Failed to create country: %v", err)
	}

	swift := models.SwiftCode{
		SwiftCode:     "TESTUSXX",
		Name:          "Test Bank",
		Address:       "Test Address",
		Town:          "Test Town",
		CountryCode:   "US",
		CodeType:      "BIC",
		IsHeadquarter: true,
	}
	if err := dbInstance.Create(&swift).Error; err != nil {
		t.Fatalf("Failed to create swift code: %v", err)
	}

	result, err := services.GetSwiftCodeByCode("TESTUSXX")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Name != "Test Bank" {
		t.Errorf("Expected Bank Name 'Test Bank', got %s", result.Name)
	}
	if result.Country.Name != "United States" {
		t.Errorf("Expected Country 'United States', got %s", result.Country.Name)
	}
}

func TestGetSwiftCodesByCountryCode(t *testing.T) {
	dbInstance := setupTestDB(t)
	db.SetDatabase(dbInstance)

	country := models.Country{
		ISO2Code: "GB",
		Name:     "United Kingdom",
		TimeZone: "Europe/London",
	}
	if err := dbInstance.Create(&country).Error; err != nil {
		t.Fatalf("Failed to create country: %v", err)
	}

	swift1 := models.SwiftCode{
		SwiftCode:     "GBTESTXXX",
		Name:          "Test Bank GB",
		Address:       "Address 1",
		Town:          "London",
		CountryCode:   "GB",
		CodeType:      "BIC",
		IsHeadquarter: false,
	}
	swift2 := models.SwiftCode{
		SwiftCode:     "GBHEADXXX",
		Name:          "Test Bank GB",
		Address:       "Address 2",
		Town:          "London",
		CountryCode:   "GB",
		CodeType:      "BIC",
		IsHeadquarter: true,
	}
	if err := dbInstance.Create(&swift1).Error; err != nil {
		t.Fatalf("Failed to create swift1: %v", err)
	}
	if err := dbInstance.Create(&swift2).Error; err != nil {
		t.Fatalf("Failed to create swift2: %v", err)
	}

	swiftCodes, countryName, err := services.GetSwiftCodesByCountryCode("GB")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if countryName != "United Kingdom" {
		t.Errorf("Expected country name 'United Kingdom', got %s", countryName)
	}
	if len(swiftCodes) != 2 {
		t.Errorf("Expected 2 swift codes, got %d", len(swiftCodes))
	}
}
