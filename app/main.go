package main

import (
	"net/http"
	"swiftapi/app/db"
	"swiftapi/app/routes"
)

func main() {
	db.ConnectDatabase()
	db.SeedDatabase(db.DB)
	router := routes.SetupRoutes()
	http.ListenAndServe(":8080", router)
}
