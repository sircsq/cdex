package main

import (
	"cdex/api"
	"cdex/db"
	"cdex/utils"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	nartDB := db.NewNartDB(config.DBSource)
	server := api.NewServer(nartDB)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
