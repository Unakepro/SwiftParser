package handlers

import (
	"encoding/json"
	"main_pack/db"
	"main_pack/services"
	"net/http"
)

type SwiftCodeResponse struct {
	Address       string         `json:"address"`
	BankName      string         `json:"bankName"`
	ISO2Code      string         `json:"ISO2Code"`
	CountryName   string         `json:"countryName"`
	IsHeadquarter bool           `json:"isHeadquarter"`
	SwiftCode     string         `json:"swiftCode"`
	Branches      []db.SwiftCode `json:"branches,omitempty"`
}

type CountryCodeResponse struct {
	ISO2Code    string         `json:"ISO2Code"`
	CountryName string         `json:"countryName"`
	SwiftCodes  []db.SwiftCode `json:"swiftCodes,omitempty"`
}

func GetSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	swiftCode := r.URL.Path[len("/v1/swift-codes/"):]

	swiftEntry, err := services.GetSwiftCodeByCode(swiftCode)
	if err != nil {
		http.Error(w, "SWIFT code not found", http.StatusNotFound)
		return
	}

	response := SwiftCodeResponse{
		Address:  swiftEntry.Address,
		BankName: swiftEntry.BankName,
		ISO2Code: swiftEntry.ISO2Code,
		// CountryName:   swiftEntry.,
		IsHeadquarter: swiftEntry.IsHeadquarter,
		SwiftCode:     swiftEntry.SwiftCode,
	}

	if swiftEntry.IsHeadquarter {
		var branches []db.SwiftCode
		db.DB.Where("bank_name = ? AND is_headquarter = ?", swiftEntry.BankName, false).Find(&branches)
		response.Branches = branches
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
