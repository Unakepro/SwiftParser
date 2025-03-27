package main

import (
	"main_pack/db"
	"main_pack/routes"
	"net/http"
)

func main() {
	db.ConnectDatabase()
	//db.SeedDatabase(db.DB)
	router := routes.SetupRoutes()
	http.ListenAndServe(":8080", router)
}
