package publisher

import "mini-kafka/broker"

type Producer struct {
	Broker *broker.Broker
}

func (p *Producer) Publish(topic, message string) {
	p.Broker.Publish(topic, message)
}
