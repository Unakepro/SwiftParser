package handlers

import (
	"encoding/json"
	"net/http"
	"swiftapi/app/db"
	"swiftapi/app/models"
	"swiftapi/app/services"

	"github.com/gorilla/mux"
)

func GetSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	swiftEntry, err := services.GetSwiftCodeByCode(swiftCode)
	if err != nil {
		http.Error(w, "SWIFT code not found", http.StatusNotFound)
		return
	}

	response := models.SwiftCodeResponse{
		Address:       swiftEntry.Address,
		BankName:      swiftEntry.Name,
		ISO2Code:      swiftEntry.Country.ISO2Code,
		CountryName:   swiftEntry.Country.Name,
		IsHeadquarter: swiftEntry.IsHeadquarter,
		SwiftCode:     swiftEntry.SwiftCode,
	}

	if swiftEntry.IsHeadquarter {
		var branches []models.SwiftCode
		db.DB.Preload("Country").Where("name = ? AND is_headquarter = ?", swiftEntry.Name, false).Find(&branches)

		var branchResponses []models.BankDataResponse
		for _, branch := range branches {
			branchResponses = append(branchResponses, models.BankDataResponse{
				Address:       branch.Address,
				BankName:      branch.Name,
				ISO2Code:      branch.Country.ISO2Code,
				SwiftCode:     branch.SwiftCode,
				IsHeadquarter: branch.IsHeadquarter,
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

	var bankCodeResponses []models.BankDataResponse
	for _, swiftCode := range swiftCodes {
		bankCodeResponses = append(bankCodeResponses, models.BankDataResponse{
			Address:       swiftCode.Address,
			BankName:      swiftCode.Name,
			ISO2Code:      swiftCode.Country.ISO2Code,
			IsHeadquarter: swiftCode.IsHeadquarter,
			SwiftCode:     swiftCode.SwiftCode,
		})
	}

	response := models.CountryCodeResponse{
		ISO2Code:    countryISO2,
		CountryName: countryName,
		BankCodes:   bankCodeResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PostSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	var request models.AddSwiftCodeRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var country models.Country
	if err := db.DB.Where("iso2_code = ?", request.ISO2Code).First(&country).Error; err != nil {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	swiftCode := models.SwiftCode{
		Address:       request.Address,
		Name:          request.BankName,
		CountryCode:   request.ISO2Code,
		IsHeadquarter: request.IsHeadquarter,
		SwiftCode:     request.SwiftCode,
	}

	if err := db.DB.Create(&swiftCode).Error; err != nil {
		http.Error(w, "Failed to create SWIFT code", http.StatusInternalServerError)
		return
	}

	response := models.MessageResponse{
		Message: "SWIFT code added successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	result := db.DB.Where("swift_code = ?", swiftCode).Delete(&models.SwiftCode{})

	if result.RowsAffected == 0 {
		http.Error(w, "SWIFT code not found", http.StatusNotFound)
		return
	}

	if result.Error != nil {
		http.Error(w, "Failed to delete SWIFT code", http.StatusInternalServerError)
		return
	}

	response := models.MessageResponse{
		Message: "SWIFT code deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
