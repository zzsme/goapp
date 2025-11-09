package events

import (
	"goapp/internal/app"
	"sync"
)

// EventType defines the type of event
type EventType string

// Event represents an event with a type and payload
type Event struct {
	Type    EventType   // Type of the event
	Payload interface{} // Data associated with the event
}

// Handler is a function that processes an event
type Handler func(event Event)

// EventBus manages event subscriptions and distributions
type EventBus struct {
	subscribers map[EventType][]Handler
	mu          sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]Handler),
	}
}

// Subscribe registers a handler function for specific event types
func (eb *EventBus) Subscribe(eventType EventType, handler Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if _, exists := eb.subscribers[eventType]; !exists {
		eb.subscribers[eventType] = []Handler{}
	}
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
	app.Info("Subscribed to event type: %s", "event_type", eventType)
}

// SubscribeMany registers a handler function for multiple event types
func (eb *EventBus) SubscribeMany(eventTypes []EventType, handler Handler) {
	for _, eventType := range eventTypes {
		eb.Subscribe(eventType, handler)
	}
}

// Publish sends an event to all subscribers of its type
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers, exists := eb.subscribers[event.Type]
	eb.mu.RUnlock()

	if !exists {
		app.Debug("No subscribers for event type: %s", "event_type", event.Type)
		return
	}

	app.Debug("Publishing event", "event_type", event.Type, "subscribers", len(handlers))

	// Async event handling
	for _, handler := range handlers {
		go func(h Handler) {
			defer func() {
				if r := recover(); r != nil {
					app.Error("Panic in event handler", "error", r)
				}
			}()
			h(event)
		}(handler)
	}
}

// AsyncPublish is an alias for Publish since we already handle events asynchronously
func (eb *EventBus) AsyncPublish(event Event) {
	eb.Publish(event)
}

// Unsubscribe removes a handler for a specific event type
func (eb *EventBus) Unsubscribe(eventType EventType, handler Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handlers, exists := eb.subscribers[eventType]; exists {
		for i, h := range handlers {
			if &h == &handler {
				eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
				app.Info("Unsubscribed from event type", "event_type", eventType)
				return
			}
		}
	}
}

// Clear removes all subscribers for a specific event type
func (eb *EventBus) Clear(eventType EventType) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	delete(eb.subscribers, eventType)
	app.Info("Cleared all subscribers for event type", "event_type", eventType)
}

// ClearAll removes all subscribers for all event types
func (eb *EventBus) ClearAll() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers = make(map[EventType][]Handler)
	app.Info("Cleared all event subscribers")
}
