package main

import (
	"github.com/dominikpalatynski/EventService/connection"
)

func main() {
	server := connection.NewAPIServer()
	server.Run()
}