package events

import (
	"goapp/internal/app"
)

// Common event types
const (
	// User events
	UserCreated   EventType = "user.created"
	UserUpdated   EventType = "user.updated"
	UserDeleted   EventType = "user.deleted"
	UserLoggedIn  EventType = "user.logged_in"
	UserLoggedOut EventType = "user.logged_out"

	// Product events
	ProductCreated EventType = "product.created"
	ProductUpdated EventType = "product.updated"
	ProductDeleted EventType = "product.deleted"
	StockUpdated   EventType = "product.stock_updated"

	// System events
	SystemStarted   EventType = "system.started"
	SystemShutdown  EventType = "system.shutdown"
	ConfigReloaded  EventType = "system.config_reloaded"
	DatabaseError   EventType = "system.database_error"
	CacheError      EventType = "system.cache_error"
	SecurityAlert   EventType = "system.security_alert"
	BackupComplete  EventType = "system.backup_complete"
	MaintenanceMode EventType = "system.maintenance_mode"
)

var (
	// DefaultBus is the default event bus instance
	DefaultBus *EventBus
)

// InitEventBus initializes the event bus and subscribes to system events
func InitEventBus() {
	DefaultBus = NewEventBus()
	app.Info("Event bus initialized")

	// Subscribe to system events for logging
	DefaultBus.Subscribe(SystemStarted, func(e Event) {
		app.Info("System started", "details", e.Payload)
	})

	DefaultBus.Subscribe(SystemShutdown, func(e Event) {
		app.Info("System shutdown", "details", e.Payload)
	})

	DefaultBus.Subscribe(DatabaseError, func(e Event) {
		app.Error("Database error", "details", e.Payload)
	})

	DefaultBus.Subscribe(CacheError, func(e Event) {
		app.Error("Cache error", "details", e.Payload)
	})

	DefaultBus.Subscribe(SecurityAlert, func(e Event) {
		app.Warn("Security alert", "details", e.Payload)
	})
}

func init() {
	// Initialize with an empty event bus, which will be properly
	// initialized after the logger in main.go
	DefaultBus = NewEventBus()
}

// Publish is a convenience function that publishes an event to the default bus
func Publish(eventType EventType, payload interface{}) {
	DefaultBus.Publish(Event{
		Type:    eventType,
		Payload: payload,
	})
}

// Subscribe is a convenience function that subscribes to an event on the default bus
func Subscribe(eventType EventType, handler Handler) {
	DefaultBus.Subscribe(eventType, handler)
}

// SubscribeMany is a convenience function that subscribes to multiple events on the default bus
func SubscribeMany(eventTypes []EventType, handler Handler) {
	DefaultBus.SubscribeMany(eventTypes, handler)
}

// Unsubscribe is a convenience function that unsubscribes from an event on the default bus
func Unsubscribe(eventType EventType, handler Handler) {
	DefaultBus.Unsubscribe(eventType, handler)
}

// Clear is a convenience function that clears all subscribers for an event type on the default bus
func Clear(eventType EventType) {
	DefaultBus.Clear(eventType)
}

// ClearAll is a convenience function that clears all subscribers on the default bus
func ClearAll() {
	DefaultBus.ClearAll()
}
