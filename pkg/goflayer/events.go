package goflayer

import (
	"context"
	"sync"
	"time"
)

// EventBus manages event subscriptions and emissions.
//
// EventBus is the core event system used throughout goflayer. It provides a thread-safe
// way to subscribe to events and emit them to all registered handlers.
type EventBus struct {
	handlers map[string][]*eventHandlerWrapper
	mu       sync.RWMutex
	// FIXED: Add counter for unique handler IDs
	nextID   uint64
}

// Event represents an event with data and timestamp.
type Event struct {
	Name string
	Data []interface{}
	Time time.Time
}

// EventHandler is a function that handles events.
// It receives variable arguments that contain event-specific data.
type EventHandler func(...interface{})

// Subscription allows unsubscribing from events.
type Subscription interface {
	Unsubscribe()
}

// NewEventBus creates a new event bus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]*eventHandlerWrapper),
		nextID:   1,
	}
}

// On subscribes to an event. Returns a subscription that can be used to unsubscribe.
//
// Multiple handlers can be subscribed to the same event. When an event is emitted,
// all handlers will be called concurrently with the same event data.
//
// The handler function should be safe for concurrent use.
func (eb *EventBus) On(event string, handler EventHandler) Subscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.handlers[event] == nil {
		eb.handlers[event] = make([]*eventHandlerWrapper, 0)
	}

	// FIXED: Use wrapper with unique ID for proper removal
	wrapper := &eventHandlerWrapper{
		id:      eb.nextID,
		handler: handler,
	}
	eb.nextID++

	eb.handlers[event] = append(eb.handlers[event], wrapper)

	return &subscription{
		eventBus: eb,
		event:    event,
		id:       wrapper.id,
	}
}

// Emit emits an event to all registered handlers.
//
// If no handlers are registered for the event, this is a no-op.
// All handlers are called concurrently in separate goroutines.
// Panics in handlers are recovered to prevent crashing the entire event system.
//
// FIXED: Wait for all handlers to complete before returning
func (eb *EventBus) Emit(event string, data ...interface{}) {
	eb.mu.RLock()
	handlers := eb.handlers[event]
	// FIXED: Copy handlers to avoid holding lock during execution
	handlersCopy := make([]*eventHandlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	eb.mu.RUnlock()

	if len(handlersCopy) == 0 {
		return
	}

	// FIXED: Wait for all handlers to complete
	var wg sync.WaitGroup
	for _, wrapper := range handlersCopy {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			// Recover from panics in handlers
			defer func() {
				if r := recover(); r != nil {
					// Could log panic here if logger available
				}
			}()
			h(data...)
		}(wrapper.handler)
	}
	wg.Wait()
}

// Once waits for an event once with optional condition checking and context cancellation.
//
// If checkCondition is nil, returns on first occurrence of the event.
// If checkCondition is provided, returns when the condition is true for event data.
// The context can be used to cancel the wait or add a timeout.
//
// Returns the event data that satisfied the condition (or first event if no condition).
func (eb *EventBus) Once(ctx context.Context, event string, checkCondition func(...interface{}) bool) ([]interface{}, error) {
	resultChan := make(chan []interface{}, 1)

	// FIXED: Ensure subscription is cleaned up
	sub := eb.On(event, func(data ...interface{}) {
		if checkCondition == nil || checkCondition(data...) {
			select {
			case resultChan <- data:
			case <-ctx.Done():
			}
		}
	})

	defer sub.Unsubscribe()

	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// RemoveAll removes all handlers for a specific event.
// Returns the number of handlers removed.
func (eb *EventBus) RemoveAll(event string) int {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	count := 0
	if handlers, ok := eb.handlers[event]; ok {
		count = len(handlers)
		delete(eb.handlers, event)
	}
	return count
}

// HasHandlers checks if there are any handlers registered for an event.
func (eb *EventBus) HasHandlers(event string) bool {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, ok := eb.handlers[event]
	return ok && len(handlers) > 0
}

// HandlerCount returns the number of handlers registered for an event.
func (eb *EventBus) HandlerCount(event string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if handlers, ok := eb.handlers[event]; ok {
		return len(handlers)
	}
	return 0
}

// eventHandlerWrapper wraps a handler with a unique ID for removal.
// FIXED: This solves the function comparison problem in Go
type eventHandlerWrapper struct {
	id      uint64
	handler EventHandler
}

// subscription implements Subscription interface.
type subscription struct {
	eventBus *EventBus
	event    string
	id       uint64
}

// Unsubscribe removes the handler from the event bus.
// After calling Unsubscribe(), the handler will no longer receive events.
// Multiple calls to Unsubscribe() are safe (idempotent).
// FIXED: Now properly identifies the handler by ID
func (s *subscription) Unsubscribe() {
	s.eventBus.mu.Lock()
	defer s.eventBus.mu.Unlock()

	handlers := s.eventBus.handlers[s.event]
	if handlers == nil {
		return
	}

	// Find and remove the handler by ID
	for i, wrapper := range handlers {
		if wrapper.id == s.id {
			// Remove by swapping with last and shrinking
			last := len(handlers) - 1
			handlers[i] = handlers[last]
			handlers[last] = nil
			s.eventBus.handlers[s.event] = handlers[:last]
			return
		}
	}
}

// EmitAsync emits an event asynchronously without waiting for handlers to complete.
// This is faster than Emit() but handlers will be executed after this function returns.
//
// FIXED: Use a single goroutine to manage async emission
func (eb *EventBus) EmitAsync(event string, data ...interface{}) {
	go eb.Emit(event, data...)
}

// EmitWithTimeout emits an event and waits for all handlers to complete or timeout.
// Returns true if all handlers completed before timeout, false if timeout occurred.
func (eb *EventBus) EmitWithTimeout(event string, timeout time.Duration, data ...interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return eb.EmitWithContext(ctx, event, data...)
}

// EmitWithContext emits an event and waits for all handlers to complete or context cancellation.
// Returns true if all handlers completed, false if context was cancelled.
func (eb *EventBus) EmitWithContext(ctx context.Context, event string, data ...interface{}) bool {
	eb.mu.RLock()
	handlers := eb.handlers[event]
	// Copy handlers to avoid holding lock
	handlersCopy := make([]*eventHandlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	eb.mu.RUnlock()

	if len(handlersCopy) == 0 {
		return true
	}

	// Use a channel to track completion
	done := make(chan struct{})
	var wg sync.WaitGroup

	for _, wrapper := range handlersCopy {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					// Recover from panic
				}
			}()

			// Check context before executing
			select {
			case <-ctx.Done():
				return
			default:
				h(data...)
			}
		}(wrapper.handler)
	}

	// FIXED: Properly wait for completion
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	}
}
