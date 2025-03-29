package models

type Country struct {
	ISO2Code string `gorm:"type:varchar(2);primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	TimeZone string `gorm:"type:varchar(100)not null"`
}

type SwiftCode struct {
	SwiftCode     string `gorm:"type:varchar(11);primaryKey;not null"`
	Name          string `gorm:"type:varchar(255);not null"`
	Address       string `gorm:"type:text"`
	Town          string `gorm:"type:varchar(100)"`
	CountryCode   string `gorm:"type:varchar(2);not null;index"`
	CodeType      string `gorm:"type:varchar(10)"`
	IsHeadquarter bool   `gorm:"type:boolean"`

	Country Country `gorm:"foreignKey:CountryCode;references:ISO2Code"`
}

type SwiftCodeResponse struct {
	Address       string             `json:"address"`
	BankName      string             `json:"bankName"`
	ISO2Code      string             `json:"ISO2Code"`
	CountryName   string             `json:"countryName"`
	IsHeadquarter bool               `json:"isHeadquarter"`
	SwiftCode     string             `json:"swiftCode"`
	Branches      []BankDataResponse `json:"branches,omitempty"`
}

type CountryCodeResponse struct {
	ISO2Code    string             `json:"ISO2Code"`
	CountryName string             `json:"countryName"`
	BankCodes   []BankDataResponse `json:"swiftCodes,omitempty"`
}

type BankDataResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	ISO2Code      string `json:"ISO2Code"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type AddSwiftCodeRequest struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	ISO2Code      string `json:"ISO2Code"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
