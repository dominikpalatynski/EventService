package main

import (
	"MonitorService/queue"
	"MonitorService/storage"
	"fmt"
)

func main() {
	dataBase, err := storage.NewMongoDbStorage()
	if err != nil {
		fmt.Printf("dataBase init failed %v", err)
	}

	queue := queue.NewQueueHandler(dataBase)

	go queue.StartMonitor()

	select {}
}