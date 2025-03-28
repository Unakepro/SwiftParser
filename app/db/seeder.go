package db

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"main_pack/models"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	db.AutoMigrate(&models.Country{}, &models.SwiftCode{})

	countries, swift_codes, err := parseCSV("data.csv")
	if err != nil {
		log.Fatal("Seeding database failed:", err)
	}

	for _, country := range countries {
		db.FirstOrCreate(&country, models.Country{ISO2Code: country.ISO2Code})
	}

	for _, swift_code := range swift_codes {
		db.Create(&swift_code)
	}

	log.Println("Database seeding completed!")
}

func parseCSV(filepath string) ([]models.Country, []models.SwiftCode, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var swift_codes []models.SwiftCode
	var countries []models.Country

	for _, row := range rows[1:] {
		iso2Code := strings.ToUpper(row[0])
		timeZone := strings.TrimSpace(row[7])
		countryName := strings.ToUpper(row[6])

		countries = append(countries, models.Country{
			ISO2Code: iso2Code,
			Name:     countryName,
			TimeZone: timeZone,
		})

		swiftCode := strings.TrimSpace(row[1])
		isHQ := strings.HasSuffix(swiftCode, "XXX")

		swift_codes = append(swift_codes, models.SwiftCode{
			SwiftCode:     swiftCode,
			Name:          row[3],
			Address:       row[4],
			Town:          row[5],
			CountryCode:   iso2Code,
			CodeType:      row[2],
			IsHeadquarter: isHQ,
		})
	}

	return countries, swift_codes, nil
}
