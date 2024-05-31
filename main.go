package main

import (
	"fmt"

	"github.com/dominikpalatynski/EventService/connection"
	"github.com/dominikpalatynski/EventService/storage"
)

func main() {

	dataBase, err := storage.NewMongoDbStorage()
	if err != nil {
		fmt.Printf("dataBase init failed %v", err)
	}

	server := connection.NewAPIServer(dataBase)
	server.Run()
}