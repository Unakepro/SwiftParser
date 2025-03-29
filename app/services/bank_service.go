package services

import (
	"swiftapi/app/db"
	"swiftapi/app/models"
)

func GetSwiftCodeByCode(code string) (*models.SwiftCode, error) {
	var swiftCode models.SwiftCode
	result := db.DB.Preload("Country").Where("swift_code = ?", code).First(&swiftCode)
	if result.Error != nil {
		return nil, result.Error
	}
	return &swiftCode, nil
}

func GetSwiftCodesByCountryCode(countryISO2 string) ([]models.SwiftCode, string, error) {
	var swiftCodes []models.SwiftCode
	var country models.Country

	if err := db.DB.Where("iso2_code = ?", countryISO2).First(&country).Error; err != nil {
		return nil, "", err
	}

	if err := db.DB.Where("country_code = ?", countryISO2).Find(&swiftCodes).Error; err != nil {
		return nil, "", err
	}

	return swiftCodes, country.Name, nil
}
