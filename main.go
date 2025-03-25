package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mini-kafka/broker"
	"mini-kafka/consumer"
	"mini-kafka/producer"
)

var topics = []string{"foo", "bar", "baz"}

func main() {
	// Start the broker
	b := broker.New()
	defer b.Close()

	// Create a new publisher
	p := publisher.Producer{Broker: b}

	// Create 100 consumers for each topic
	for i, topic := range topics {
		for j := 0; j < 100; j++ {
			c := consumer.New(b, topic)
			// NOTE: Listen starts a goroutine to listen for messages
			c.Listen(fmt.Sprintf("%s-consumer-%d-%d", topic, i, j))
		}
	}

	fmt.Println("Enter messages in format: <topic> <message>")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid input. Format: <topic> <message>")
			continue
		}
		p.Publish(parts[0], parts[1])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
