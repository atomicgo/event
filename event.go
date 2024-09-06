package event

import "sync"

// Event represents an event system that can handle multiple listeners.
type Event[T any] struct {
	listeners []chan T
	mu        sync.Mutex
	closed    bool
}

// New creates a new event.
func New[T any]() *Event[T] {
	// Create a new event
	return &Event[T]{
		listeners: []chan T{},
	}
}

// Trigger triggers the event and notifies all listeners.
func (e *Event[T]) Trigger(value T) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, listener := range e.listeners {
		go func(l chan T) {
			if !e.closed {
				l <- value
			}
		}(listener)
	}
}

// Listen gets called when the event is triggered.
func (e *Event[T]) Listen(f func(T)) {
	// Check if the event is closed
	if e.closed {
		return
	}

	// Create listener slice if it doesn't exist
	if e.listeners == nil {
		e.listeners = []chan T{}
	}

	// Create a new channel
	ch := make(chan T)

	e.mu.Lock()
	e.listeners = append(e.listeners, ch)
	e.mu.Unlock()

	go func() {
		for v := range ch {
			if !e.closed {
				f(v)
			}
		}
	}()
}

// Close closes the event and all its listeners.
// After calling this method, the event can't be used anymore and new listeners can't be added.
func (e *Event[T]) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, listener := range e.listeners {
		close(listener)
	}

	e.listeners = nil
	e.closed = true
}
