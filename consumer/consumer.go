package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type consumer struct {
	name, url, topic string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run consumer.go <topic>")
		os.Exit(1)
	}
	topic := os.Args[1]
	uid := uuid.New()
	sub := &consumer{
		name:  uid.String(),
		url:   "ws://localhost:9092",
		topic: topic,
	}

	url := fmt.Sprintf("%s/ws?name=%s&topic=%s", sub.url, sub.name, sub.topic)
	log.Printf("connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		log.Printf("received: %s", message)
	}
}
