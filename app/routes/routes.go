package routes

import (
	"swiftapi/app/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/swift-codes/{swiftCode}", handlers.GetSwiftCodeHandler).Methods("GET")

	r.HandleFunc("/v1/swift-codes/country/{countryISO2}", handlers.GetSwiftCodesByCountryHandler).Methods("GET")

	r.HandleFunc("/v1/swift-codes", handlers.PostSwiftCodeHandler).Methods("POST")

	r.HandleFunc("/v1/swift-codes/{swiftCode}", handlers.DeleteSwiftCodeHandler).Methods("DELETE")

	return r
}
