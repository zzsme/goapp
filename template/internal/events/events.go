package events

import (
	"sync"
)

// EventType represents a type of event
type EventType string

// System events
const (
	SystemStarted  EventType = "system.started"
	SystemShutdown EventType = "system.shutdown"
)

// User events
const (
	UserCreated EventType = "user.created"
	UserUpdated EventType = "user.updated"
	UserDeleted EventType = "user.deleted"
)

// Product events
const (
	ProductCreated EventType = "product.created"
	ProductUpdated EventType = "product.updated"
	ProductDeleted EventType = "product.deleted"
)

// EventHandler represents a function that handles an event
type EventHandler func(data map[string]interface{})

// eventBus manages event subscriptions and publishing
type eventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

var (
	defaultBus *eventBus
	once       sync.Once
)

// InitEventBus initializes the event bus system
func InitEventBus() {
	once.Do(func() {
		defaultBus = &eventBus{
			handlers: make(map[EventType][]EventHandler),
		}
	})
}

// Subscribe registers a handler for a specific event type
func Subscribe(eventType EventType, handler EventHandler) {
	defaultBus.mu.Lock()
	defer defaultBus.mu.Unlock()

	if defaultBus.handlers == nil {
		defaultBus.handlers = make(map[EventType][]EventHandler)
	}

	defaultBus.handlers[eventType] = append(defaultBus.handlers[eventType], handler)
}

// Publish sends an event to all registered handlers
func Publish(eventType EventType, data map[string]interface{}) {
	defaultBus.mu.RLock()
	handlers := defaultBus.handlers[eventType]
	defaultBus.mu.RUnlock()

	// Execute handlers in goroutines
	for _, handler := range handlers {
		go func(h EventHandler) {
			h(data)
		}(handler)
	}
}

// Unsubscribe removes a handler for a specific event type
func Unsubscribe(eventType EventType, handler EventHandler) {
	defaultBus.mu.Lock()
	defer defaultBus.mu.Unlock()

	if handlers, ok := defaultBus.handlers[eventType]; ok {
		for i, h := range handlers {
			if &h == &handler {
				defaultBus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Clear removes all handlers for a specific event type
func Clear(eventType EventType) {
	defaultBus.mu.Lock()
	defer defaultBus.mu.Unlock()

	delete(defaultBus.handlers, eventType)
}

// ClearAll removes all event handlers
func ClearAll() {
	defaultBus.mu.Lock()
	defer defaultBus.mu.Unlock()

	defaultBus.handlers = make(map[EventType][]EventHandler)
}
