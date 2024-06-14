package main

import (
	"log"

	"ecom/cmd/api"
	"ecom/config"
	"ecom/db"
)

func main() {

	db, q, err := db.NewPostgresStorage(config.Envs.DATABASE_URL)
	if err != nil {
		log.Fatal("Couldn't connect to the database")
	}

	server := api.NewAPIServer(config.Envs.PORT, db, q)

	log.Fatal(server.Run())
}
