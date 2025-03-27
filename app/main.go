package main

import (
	"log"
	"main_pack/db"
	"main_pack/routes"
	"net/http"
)

func main() {
	db.ConnectDatabase()
	//db.SeedDatabase(db.DB)
	routes.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
