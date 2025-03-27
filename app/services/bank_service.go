package services

import (
	"main_pack/db"
)

func GetSwiftCodeByCode(code string) (*db.SwiftCode, error) {
	var swiftCode db.SwiftCode
	result := db.DB.Preload("Country").Where("swift_code = ?", code).First(&swiftCode)
	if result.Error != nil {
		return nil, result.Error
	}
	return &swiftCode, nil
}

func GetSwiftCodesByCountryCode(countryISO2 string) ([]db.SwiftCode, string, error) {
	var swiftCodes []db.SwiftCode
	var country db.Country

	if err := db.DB.Where("iso2_code = ?", countryISO2).First(&country).Error; err != nil {
		return nil, "", err
	}

	if err := db.DB.Where("iso2_code = ?", countryISO2).Find(&swiftCodes).Error; err != nil {
		return nil, "", err
	}

	return swiftCodes, country.CountryName, nil
}
