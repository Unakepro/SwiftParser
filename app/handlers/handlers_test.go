package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"main_pack/db"
	"main_pack/models"
	"main_pack/routes"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIntegrationDB(t *testing.T) *gorm.DB {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := testDB.AutoMigrate(&models.Country{}, &models.SwiftCode{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return testDB
}

func seedData(dbInstance *gorm.DB) {
	country := models.Country{
		ISO2Code: "US",
		Name:     "United States",
		TimeZone: "America/New_York",
	}
	dbInstance.Create(&country)

	swift := models.SwiftCode{
		SwiftCode:     "TESTUSXX",
		Name:          "Test Bank",
		Address:       "Test Address",
		Town:          "Test Town",
		CountryCode:   "US",
		CodeType:      "BIC",
		IsHeadquarter: true,
	}
	dbInstance.Create(&swift)
}

func TestGetSwiftCodeHandler(t *testing.T) {
	dbInstance := setupIntegrationDB(t)
	db.SetDatabase(dbInstance)
	seedData(dbInstance)

	router := routes.SetupRoutes()

	req, _ := http.NewRequest("GET", "/v1/swift-codes/TESTUSXX", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var response models.SwiftCodeResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if response.BankName != "Test Bank" {
		t.Errorf("Expected BankName 'Test Bank', got %s", response.BankName)
	}
}

func TestGetSwiftCodesByCountryHandler(t *testing.T) {
	dbInstance := setupIntegrationDB(t)
	db.SetDatabase(dbInstance)

	country := models.Country{
		ISO2Code: "CA",
		Name:     "Canada",
		TimeZone: "America/Toronto",
	}
	dbInstance.Create(&country)
	swift1 := models.SwiftCode{
		SwiftCode:     "CATEST001",
		Name:          "Test Bank CA",
		Address:       "Address 1",
		Town:          "Toronto",
		CountryCode:   "CA",
		CodeType:      "BIC",
		IsHeadquarter: false,
	}
	swift2 := models.SwiftCode{
		SwiftCode:     "CAHEAD001",
		Name:          "Test Bank CA",
		Address:       "Address 2",
		Town:          "Toronto",
		CountryCode:   "CA",
		CodeType:      "BIC",
		IsHeadquarter: true,
	}
	dbInstance.Create(&swift1)
	dbInstance.Create(&swift2)

	router := routes.SetupRoutes()

	req, _ := http.NewRequest("GET", "/v1/swift-codes/country/CA", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var response models.CountryCodeResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if response.CountryName != "Canada" {
		t.Errorf("Expected CountryName 'Canada', got %s", response.CountryName)
	}
	if len(response.BankCodes) != 2 {
		t.Errorf("Expected 2 bank codes, got %d", len(response.BankCodes))
	}
}

func TestPostSwiftCodeHandler(t *testing.T) {
	dbInstance := setupIntegrationDB(t)
	db.SetDatabase(dbInstance)

	country := models.Country{
		ISO2Code: "AU",
		Name:     "Australia",
		TimeZone: "Australia/Sydney",
	}
	dbInstance.Create(&country)

	router := routes.SetupRoutes()

	newSwift := models.AddSwiftCodeRequest{
		Address:       "New Address",
		BankName:      "New Bank",
		ISO2Code:      "AU",
		IsHeadquarter: false,
		SwiftCode:     "AUNEW001",
	}

	payload, err := json.Marshal(newSwift)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	req, _ := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 on POST, got %d", rr.Code)
	}

	var response models.MessageResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode POST response: %v", err)
	}
	if response.Message != "SWIFT code added successfully" {
		t.Errorf("Expected success message, got %s", response.Message)
	}
}

func TestDeleteSwiftCodeHandler(t *testing.T) {
	dbInstance := setupIntegrationDB(t)
	db.SetDatabase(dbInstance)

	country := models.Country{
		ISO2Code: "DE",
		Name:     "Germany",
		TimeZone: "Europe/Berlin",
	}
	dbInstance.Create(&country)
	swift := models.SwiftCode{
		SwiftCode:     "DETEST001",
		Name:          "Test Bank DE",
		Address:       "Address DE",
		Town:          "Berlin",
		CountryCode:   "DE",
		CodeType:      "BIC",
		IsHeadquarter: false,
	}
	dbInstance.Create(&swift)

	router := routes.SetupRoutes()

	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/DETEST001", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 on DELETE, got %d", rr.Code)
	}

	var response models.MessageResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode DELETE response: %v", err)
	}
	if response.Message != "SWIFT code deleted successfully" {
		t.Errorf("Expected deletion message, got %s", response.Message)
	}
}
