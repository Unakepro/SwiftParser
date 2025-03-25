package db

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

type Country struct {
	ISO2Code    string `gorm:"primaryKey;size:2"`
	CountryName string `gorm:"size:100;not null"`
}

type SwiftCode struct {
	SwiftCode string `gorm:"primaryKey;size:11"`
	BankName  string `gorm:"size:255;not null"`
	ISO2Code  string `gorm:"size:2;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;references:ISO2Code"`

	Address  string `gorm:"type:text"`
	TownName string `gorm:"size:100"`
	TimeZone string `gorm:"size:50"`

	IsHeadquarter bool    `gorm:"not null"`
	HeadSwiftCode *string `gorm:"size:11"`
}

func SeedDatabase(db *gorm.DB) {
	db.AutoMigrate(&Country{}, &SwiftCode{})

	countries, swiftRecords, err := parseCSV("data.csv")
	if err != nil {
		log.Fatal("Seeding database failed:", err)
	}

	for _, country := range countries {
		db.FirstOrCreate(&country, Country{ISO2Code: country.ISO2Code})
	}

	for _, record := range swiftRecords {
		db.Create(&record)
	}

	log.Println("Database seeding completed!")
}

func parseCSV(filepath string) ([]Country, []SwiftCode, error) {
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

	countryMap := make(map[string]string)
	var swiftRecords []SwiftCode

	for _, row := range rows[1:] {
		// fmt.Printf("first row %s\n", row[6])
		iso2Code := strings.ToUpper(row[0])
		countryMap[iso2Code] = strings.ToUpper(row[6])

		swiftCode := row[1]
		isHQ := strings.HasSuffix(swiftCode, "XXX")
		//var headSwiftCode *string
		// if !isHQ {
		// 	hqCode := swiftCode[:8] + "XXX"
		// 	headSwiftCode = &hqCode
		// }

		swiftRecord := SwiftCode{
			SwiftCode:     swiftCode,
			BankName:      row[3],
			ISO2Code:      iso2Code,
			Address:       row[4],
			TownName:      row[5],
			TimeZone:      row[7],
			IsHeadquarter: isHQ,
			// HeadSwiftCode: headSwiftCode,
		}

		swiftRecords = append(swiftRecords, swiftRecord)
	}

	var countries []Country
	for code, name := range countryMap {
		countries = append(countries, Country{ISO2Code: code, CountryName: name})
	}
	//}

	return countries, swiftRecords, nil
}
