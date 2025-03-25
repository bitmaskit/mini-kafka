package broker

import "sync"

type Broker struct {
	rwmu        sync.RWMutex             // nil value is fine
	subscribers map[string][]chan string // map topic -> subscribers
}

func New() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan string),
	}
}

// Subscribe returns a new channel for a given topic
func (b *Broker) Subscribe(topic string) <-chan string {
	ch := make(chan string, 10) // buffered to avoid blocking
	b.rwmu.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	b.rwmu.Unlock()
	return ch
}

// Publish sends a message to all subscribers of a topic
func (b *Broker) Publish(topic, message string) {
	b.rwmu.RLock()               // RLock because we are only reading the map
	subs := b.subscribers[topic] // subs is []chan string
	b.rwmu.RUnlock()             // RUnlock because we are done reading the map

	for _, ch := range subs { // iterate over all subscribers
		select {
		case ch <- message: // fast consumer, success
		default: // slow consumer, drop message
		}
	}
}

// Unsubscribe removes a specific channel from the topic
func (b *Broker) Unsubscribe(topic string, ch <-chan string) {
	b.rwmu.Lock()
	defer b.rwmu.Unlock()

	subs := b.subscribers[topic]
	for i, sub := range subs {
		if sub == ch {
			b.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// Close all channels
func (b *Broker) Close() {
	b.rwmu.Lock()
	defer b.rwmu.Unlock()

	for _, subs := range b.subscribers { // subs []chan string
		for _, ch := range subs { // ch is chan string
			close(ch) // close the channel
		}
	}
	b.subscribers = nil // set the map to nil
}
