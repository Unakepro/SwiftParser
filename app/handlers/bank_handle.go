package handlers

import (
	"encoding/json"
	"main_pack/db"
	"main_pack/services"
	"net/http"
)

type SwiftCodeResponse struct {
	Address       string              `json:"address"`
	BankName      string              `json:"bankName"`
	ISO2Code      string              `json:"ISO2Code"`
	CountryName   string              `json:"countryName"`
	IsHeadquarter bool                `json:"isHeadquarter"`
	SwiftCode     string              `json:"swiftCode"`
	Branches      []SwiftCodeResponse `json:"branches,omitempty"`
}

type BankDataResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	ISO2Code      string `json:"ISO2Code"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type CountryCodeResponse struct {
	ISO2Code    string             `json:"ISO2Code"`
	CountryName string             `json:"countryName"`
	BankCodes   []BankDataResponse `json:"swiftCodes,omitempty"`
}

func GetSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.URL.Path[len("/v1/swift-codes/"):]

	swiftEntry, err := services.GetSwiftCodeByCode(swiftCode)
	if err != nil {
		http.Error(w, "SWIFT code not found", http.StatusNotFound)
		return
	}

	response := SwiftCodeResponse{
		Address:       swiftEntry.Address,
		BankName:      swiftEntry.BankName,
		ISO2Code:      swiftEntry.ISO2Code,
		CountryName:   swiftEntry.Country.CountryName,
		IsHeadquarter: swiftEntry.IsHeadquarter,
		SwiftCode:     swiftEntry.SwiftCode,
	}

	if swiftEntry.IsHeadquarter {
		var branches []db.SwiftCode
		db.DB.Preload("Country").Where("bank_name = ? AND is_headquarter = ?", swiftEntry.BankName, false).Find(&branches)

		var branchResponses []SwiftCodeResponse
		for _, branch := range branches {
			branchResponses = append(branchResponses, SwiftCodeResponse{
				Address:       branch.Address,
				BankName:      branch.BankName,
				ISO2Code:      branch.ISO2Code,
				SwiftCode:     branch.SwiftCode,
				IsHeadquarter: branch.IsHeadquarter,
				CountryName:   branch.Country.CountryName,
			})
		}
		response.Branches = branchResponses
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetSwiftCodesByCountryHandler(w http.ResponseWriter, r *http.Request) {
	countryISO2 := r.URL.Path[len("/v1/swift-codes/country/"):]

	swiftCodes, countryName, err := services.GetSwiftCodesByCountryCode(countryISO2)
	if err != nil {
		http.Error(w, "Country code not found", http.StatusNotFound)
		return
	}

	var bankCodeResponses []BankDataResponse
	for _, swiftCode := range swiftCodes {
		bankCodeResponses = append(bankCodeResponses, BankDataResponse{
			Address:       swiftCode.Address,
			BankName:      swiftCode.BankName,
			ISO2Code:      swiftCode.ISO2Code,
			IsHeadquarter: swiftCode.IsHeadquarter,
			SwiftCode:     swiftCode.SwiftCode,
		})
	}

	response := CountryCodeResponse{
		ISO2Code:    countryISO2,
		CountryName: countryName,
		BankCodes:   bankCodeResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
