package bus

import (
	"context"
	"sync"
)

type Handler func(ctx context.Context, event interface{}) error

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]Handler),
	}
}

func (b *EventBus) Subscribe(topic string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], handler)
}

func (b *EventBus) Publish(ctx context.Context, topic string, event interface{}) error {
	b.mu.RLock()
	handlers := b.subscribers[topic]
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
