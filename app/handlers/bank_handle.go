package handlers

import (
	"encoding/json"
	"main_pack/db"
	"main_pack/services"
	"net/http"

	"github.com/gorilla/mux"
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

func GetSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

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
	vars := mux.Vars(r)
	countryISO2 := vars["countryISO2"]

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

func PostSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	var request AddSwiftCodeRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var country db.Country
	if err := db.DB.Where("iso2_code = ?", request.ISO2Code).First(&country).Error; err != nil {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	swiftCode := db.SwiftCode{
		Address:       request.Address,
		BankName:      request.BankName,
		ISO2Code:      request.ISO2Code,
		IsHeadquarter: request.IsHeadquarter,
		SwiftCode:     request.SwiftCode,
	}

	if err := db.DB.Create(&swiftCode).Error; err != nil {
		http.Error(w, "Failed to create SWIFT code", http.StatusInternalServerError)
		return
	}

	response := MessageResponse{
		Message: "SWIFT code added successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	if err := db.DB.Where("swift_code = ?", swiftCode).Delete(&db.SwiftCode{}).Error; err != nil {
		http.Error(w, "Failed to delete SWIFT code", http.StatusInternalServerError)
		return
	}

	response := MessageResponse{
		Message: "SWIFT code deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
