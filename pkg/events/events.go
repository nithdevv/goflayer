// Package events предоставляет типизированный event bus.
package events

import (
	"context"
	"sync"
	"time"
)

// Handler is a function that handles an event.
type Handler func(...interface{})

// Subscription represents an event subscription.
type Subscription struct {
	bus      *Bus
	event    string
	id       uint64
	cancel   context.CancelFunc
	done     chan struct{}
}

// Unsubscribe removes the subscription.
func (s *Subscription) Unsubscribe() {
	if s.cancel != nil {
		s.cancel()
	}
	<-s.done
}

// Bus is a thread-safe event bus.
type Bus struct {
	mu        sync.RWMutex
	handlers  map[string]map[uint64]Handler
	nextID    uint64
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewBus creates a new event bus.
func NewBus() *Bus {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bus{
		handlers: make(map[string]map[uint64]Handler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Subscribe subscribes to an event.
func (b *Bus) Subscribe(event string, handler Handler) *Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.handlers[event] == nil {
		b.handlers[event] = make(map[uint64]Handler)
	}

	id := b.nextID
	b.nextID++

	b.handlers[event][id] = handler

	subCtx, subCancel := context.WithCancel(b.ctx)
	done := make(chan struct{})

	sub := &Subscription{
		bus:    b,
		event:  event,
		id:     id,
		cancel: subCancel,
		done:   done,
	}

	// Wait for unsubscribe
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		<-subCtx.Done()
		b.removeHandler(event, id)
		close(done)
	}()

	return sub
}

// removeHandler removes a handler (must be called with unlocked mutex).
func (b *Bus) removeHandler(event string, id uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if handlers, ok := b.handlers[event]; ok {
		delete(handlers, id)
		if len(handlers) == 0 {
			delete(b.handlers, event)
		}
	}
}

// Emit emits an event synchronously.
func (b *Bus) Emit(event string, data ...interface{}) {
	b.mu.RLock()
	handlers, ok := b.handlers[event]
	if !ok || len(handlers) == 0 {
		b.mu.RUnlock()
		return
	}

	// Copy handlers to avoid holding lock
	handlersCopy := make([]Handler, 0, len(handlers))
	for _, h := range handlers {
		handlersCopy = append(handlersCopy, h)
	}
	b.mu.RUnlock()

	for _, h := range handlersCopy {
		h(data...)
	}
}

// EmitAsync emits an event asynchronously.
func (b *Bus) EmitAsync(event string, data ...interface{}) {
	b.mu.RLock()
	handlers, ok := b.handlers[event]
	if !ok || len(handlers) == 0 {
		b.mu.RUnlock()
		return
	}

	handlersCopy := make([]Handler, 0, len(handlers))
	for _, h := range handlers {
		handlersCopy = append(handlersCopy, h)
	}
	b.mu.RUnlock()

	for _, h := range handlersCopy {
		go h(data...)
	}
}

// EmitWithTimeout emits an event with a timeout.
func (b *Bus) EmitWithTimeout(event string, timeout time.Duration, data ...interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return b.EmitWithContext(ctx, event, data...)
}

// EmitWithContext emits an event with context.
func (b *Bus) EmitWithContext(ctx context.Context, event string, data ...interface{}) bool {
	b.mu.RLock()
	handlers, ok := b.handlers[event]
	if !ok || len(handlers) == 0 {
		b.mu.RUnlock()
		return true
	}

	handlersCopy := make([]Handler, 0, len(handlers))
	for _, h := range handlers {
		handlersCopy = append(handlersCopy, h)
	}
	b.mu.RUnlock()

	// Run handlers with context cancellation
	for _, h := range handlersCopy {
		go func(handler Handler) {
			defer func() { recover() }()

			select {
			case <-ctx.Done():
				return
			default:
				handler(data...)
			}
		}(h)
	}

	return true
}

// Close closes the event bus.
func (b *Bus) Close() error {
	b.cancel()
	b.wg.Wait()
	return nil
}

// HasHandlers returns true if event has handlers.
func (b *Bus) HasHandlers(event string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers, ok := b.handlers[event]
	return ok && len(handlers) > 0
}

// HandlerCount returns the number of handlers for an event.
func (b *Bus) HandlerCount(event string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, ok := b.handlers[event]; ok {
		return len(handlers)
	}
	return 0
}

// Clear removes all handlers.
func (b *Bus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = make(map[string]map[uint64]Handler)
}

// Once waits for an event once.
func (b *Bus) Once(ctx context.Context, event string, predicate func(...interface{}) bool) ([]interface{}, error) {
	resultCh := make(chan []interface{}, 1)

	sub := b.Subscribe(event, func(data ...interface{}) {
		if predicate == nil || predicate(data...) {
			select {
			case resultCh <- data:
			case <-ctx.Done():
			}
		}
	})
	defer sub.Unsubscribe()

	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
