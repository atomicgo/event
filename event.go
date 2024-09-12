package event

import (
	"errors"
	"sync"
)

// ErrEventClosed is returned when an operation is attempted on a closed event.
var ErrEventClosed = errors.New("event is closed")

// Event represents a generic, thread-safe event system that can handle multiple listeners.
// The type parameter T specifies the type of data that the event carries when triggered.
type Event[T any] struct {
	listeners []func(T)
	mu        sync.RWMutex
	closed    bool
}

// New creates and returns a new Event instance for the specified type T.
func New[T any]() *Event[T] {
	return &Event[T]{}
}

// Trigger notifies all registered listeners by invoking their callback functions with the provided value.
// It runs each listener in a separate goroutine and waits for all listeners to complete.
// Returns ErrEventClosed if the event has been closed.
func (e *Event[T]) Trigger(value T) error {
	e.mu.RLock()
	if e.closed {
		e.mu.RUnlock()
		return ErrEventClosed
	}

	// Copy the listeners to avoid holding the lock during execution.
	// This ensures that triggering the event is thread-safe even if listeners are added or removed concurrently.
	listeners := make([]func(T), len(e.listeners))
	copy(listeners, e.listeners)
	e.mu.RUnlock()

	var wg sync.WaitGroup
	for _, listener := range listeners {
		wg.Add(1)

		go func(f func(T)) {
			defer wg.Done()
			f(value)
		}(listener)
	}

	wg.Wait()

	return nil
}

// Listen registers a new listener callback function for the event.
// The listener will be invoked with the event's data whenever Trigger is called.
// Returns ErrEventClosed if the event has been closed.
func (e *Event[T]) Listen(f func(T)) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return ErrEventClosed
	}

	e.listeners = append(e.listeners, f)

	return nil
}

// Close closes the event system, preventing any new listeners from being added or events from being triggered.
// After calling Close, any subsequent calls to Trigger or Listen will return ErrEventClosed.
// Existing listeners are removed, and resources are cleaned up.
func (e *Event[T]) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return
	}

	e.closed = true
	e.listeners = nil // Release references to listener functions
}
