package services

import (
	"main_pack/db"
)

func GetSwiftCodeByCode(code string) (*db.SwiftCode, error) {
	var swiftCode db.SwiftCode
	result := db.DB.Where("swift_code = ?", code).First(&swiftCode)
	if result.Error != nil {
		return nil, result.Error
	}
	return &swiftCode, nil
}
