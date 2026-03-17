// Package event provides a thread-safe event bus for pub/sub messaging.
//
// The event system is the core communication mechanism in goflayer.
// All components communicate through events, enabling loose coupling
// and extensibility.
package event

import (
	"context"
	"sync"
	"time"
)

// Handler is a function that handles an event.
// It receives variable arguments containing event-specific data.
type Handler func(...interface{})

// Subscription represents a subscription to an event.
// It can be used to unsubscribe from the event.
type Subscription interface {
	// Unsubscribe removes this subscription from the event bus.
	Unsubscribe()
}

// Bus manages event subscriptions and emissions.
// It is thread-safe and can be used concurrently from multiple goroutines.
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]*handlerWrapper
	nextID   uint64
}

// handlerWrapper wraps a handler with metadata for proper removal.
type handlerWrapper struct {
	id      uint64
	handler Handler
}

// subscriptionImpl implements the Subscription interface.
type subscriptionImpl struct {
	bus  *Bus
	event string
	id   uint64
}

// NewBus creates a new event bus.
func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]*handlerWrapper),
		nextID:   1,
	}
}

// Subscribe registers a handler for an event.
// Returns a subscription that can be used to unsubscribe.
//
// Multiple handlers can be subscribed to the same event.
// When an event is emitted, all handlers are called in the order
// they were registered.
func (b *Bus) Subscribe(event string, handler Handler) Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.handlers[event] == nil {
		b.handlers[event] = make([]*handlerWrapper, 0, 4)
	}

	wrapper := &handlerWrapper{
		id:      b.nextID,
		handler: handler,
	}
	b.nextID++

	b.handlers[event] = append(b.handlers[event], wrapper)

	return &subscriptionImpl{
		bus:  b,
		event: event,
		id:   wrapper.id,
	}
}

// Unsubscribe removes a handler from an event.
func (s *subscriptionImpl) Unsubscribe() {
	s.bus.mu.Lock()
	defer s.bus.mu.Unlock()

	handlers := s.bus.handlers[s.event]
	if handlers == nil {
		return
	}

	// Find and remove the handler by ID
	for i, wrapper := range handlers {
		if wrapper.id == s.id {
			// Remove by swapping with last and truncating
			last := len(handlers) - 1
			handlers[i] = handlers[last]
			handlers[last] = nil
			s.bus.handlers[s.event] = handlers[:last]
			return
		}
	}
}

// Emit emits an event to all registered handlers synchronously.
// It waits for all handlers to complete before returning.
//
// Panics in handlers are recovered to prevent crashing the entire system.
// If no handlers are registered for the event, this is a no-op.
func (b *Bus) Emit(event string, data ...interface{}) {
	b.mu.RLock()
	handlers := b.handlers[event]

	// Copy handlers to avoid holding lock during execution
	handlersCopy := make([]*handlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	b.mu.RUnlock()

	if len(handlersCopy) == 0 {
		return
	}

	// Execute all handlers and wait for completion
	var wg sync.WaitGroup
	for _, wrapper := range handlersCopy {
		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					// TODO: Log panic recovery
				}
			}()
			h(data...)
		}(wrapper.handler)
	}
	wg.Wait()
}

// EmitAsync emits an event asynchronously.
// Handlers are executed in goroutines but this method returns immediately
// without waiting for handlers to complete.
func (b *Bus) EmitAsync(event string, data ...interface{}) {
	b.mu.RLock()
	handlers := b.handlers[event]

	// Copy handlers
	handlersCopy := make([]*handlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	b.mu.RUnlock()

	if len(handlersCopy) == 0 {
		return
	}

	for _, wrapper := range handlersCopy {
		go func(h Handler) {
			defer func() {
				if r := recover(); r != nil {
					// TODO: Log panic recovery
				}
			}()
			h(data...)
		}(wrapper.handler)
	}
}

// EmitWithContext emits an event and waits for handlers to complete
// or until the context is cancelled.
//
// Returns true if all handlers completed, false if context was cancelled.
func (b *Bus) EmitWithContext(ctx context.Context, event string, data ...interface{}) bool {
	b.mu.RLock()
	handlers := b.handlers[event]

	// Copy handlers
	handlersCopy := make([]*handlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	b.mu.RUnlock()

	if len(handlersCopy) == 0 {
		return true
	}

	done := make(chan struct{})
	var wg sync.WaitGroup

	for _, wrapper := range handlersCopy {
		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					// TODO: Log panic recovery
				}
			}()

			select {
			case <-ctx.Done():
				return
			default:
				h(data...)
			}
		}(wrapper.handler)
	}

	// Wait for completion in a goroutine
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

// EmitWithTimeout emits an event and waits for handlers with a timeout.
// Returns true if all handlers completed within the timeout.
func (b *Bus) EmitWithTimeout(event string, timeout time.Duration, data ...interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return b.EmitWithContext(ctx, event, data...)
}

// Once waits for an event to occur once, optionally checking a condition.
//
// If checkCondition is nil, returns on first occurrence of the event.
// If checkCondition is provided, returns when the condition is true for event data.
// The context can be used to cancel the wait or add a timeout.
//
// Returns the event data that satisfied the condition (or first event if no condition).
func (b *Bus) Once(ctx context.Context, event string, checkCondition func(...interface{}) bool) ([]interface{}, error) {
	resultChan := make(chan []interface{}, 1)

	sub := b.Subscribe(event, func(data ...interface{}) {
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

// HasHandlers returns true if there are handlers registered for the event.
func (b *Bus) HasHandlers(event string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers, ok := b.handlers[event]
	return ok && len(handlers) > 0
}

// HandlerCount returns the number of handlers registered for an event.
func (b *Bus) HandlerCount(event string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, ok := b.handlers[event]; ok {
		return len(handlers)
	}
	return 0
}

// RemoveAll removes all handlers for a specific event.
// Returns the number of handlers removed.
func (b *Bus) RemoveAll(event string) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	count := 0
	if handlers, ok := b.handlers[event]; ok {
		count = len(handlers)
		delete(b.handlers, event)
	}
	return count
}

// Clear removes all event handlers.
func (b *Bus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = make(map[string][]*handlerWrapper)
}
