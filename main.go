package main

import (
	"mini-kafka/broker"
	"mini-kafka/server"
)

func main() {
	brk := broker.New()
	srv := server.Server{Broker: brk}

	srv.Start(":9092")
}
