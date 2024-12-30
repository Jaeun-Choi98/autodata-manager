package broker

import "sync"

type Message struct {
	Topic string
	Body  string
}

type Broker struct {
	subscribers map[string][]chan Message
	lock        sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan Message),
	}
}

func (b *Broker) Publish(topic string, msg Message) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if subs, ok := b.subscribers[topic]; ok {
		for _, ch := range subs {
			go func(c chan Message) {
				c <- msg
			}(ch)
		}
	}
}

func (b *Broker) Subscribe(topic string) <-chan Message {
	b.lock.Lock()
	defer b.lock.Unlock()

	ch := make(chan Message, 10)
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	return ch
}

func (b *Broker) Unsubscribe(topic string, ch <-chan Message) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if subs, ok := b.subscribers[topic]; ok {
		for i, sub := range subs {
			if sub == ch {
				b.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(sub)
				break
			}
		}
	}
}
