package consumer

import (
	"fmt"
	"mini-kafka/broker"
)

type Consumer struct {
	topic string
	ch    <-chan string
}

func New(b *broker.Broker, topic string) *Consumer {
	// NOTE: Consumer does not need to keep a reference to the broker
	ch := b.Subscribe(topic) // channel is returned when subscribing
	return &Consumer{
		topic: topic,
		ch:    ch, // keep the channel
	}
}

func (c *Consumer) Listen(name string) {
	// Starting a goroutine to listen for messages
	go func() {
		for msg := range c.ch {
			fmt.Printf("[%s] received on topic '%s': %s\n", name, c.topic, msg)
		}
	}()
}
