package routes

import (
	"main_pack/handlers"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/v1/swift-codes/", handlers.GetSwiftCodeHandler)
	// http.HandleFunc("/v1/swift-codes/country", handlers.GetSwiftCodesByCountryHandler)
}
