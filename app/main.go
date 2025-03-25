package main

import (
	"main_pack/db"
)

func main() {
	db.ConnectDatabase()
	db.SeedDatabase(db.DB)

}
