package event

import (
	"errors"
	"sync"
)

// ErrEventClosed is returned when an operation is attempted on a closed event.
var ErrEventClosed = errors.New("event is closed")

// ErrUnknownListener is returned when attempting to unregister an unknown id.
var ErrUnknownListener = errors.New("listener id is unknown")

// Event represents a generic, thread-safe event system that can handle multiple listeners.
// The type parameter T specifies the type of data that the event carries when triggered.
type Event[T any] struct {
	listeners map[int]func(T)
	mu        sync.RWMutex
	closed    bool
	nextID    int
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
	listeners := make([]func(T), 0, len(e.listeners))
	for _, l := range e.listeners {
		listeners = append(listeners, l)
	}
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
// Returns an ID which can be used with StopListening to deregister the listener.
// Returns ErrEventClosed if the event has been closed.
func (e *Event[T]) Listen(f func(T)) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Lazy init to allow use of zero value
	if e.listeners == nil {
		e.listeners = make(map[int]func(T))
	}

	if e.closed {
		return -1, ErrEventClosed
	}

	id := e.getID()

	e.listeners[id] = f

	return id, nil
}

// StopListening unregisters a listener, using the ID returned from Listen.
// The callback which was registered with that ID will no longer be called
// and any associated resources will be released.
// If the ID passed in is not registered, ErrUnknownListener will be returned.
func (e *Event[T]) StopListening(id int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, ok := e.listeners[id]
	if !ok {
		return ErrUnknownListener
	}
	delete(e.listeners, id)
	return nil
}

func (e *Event[T]) getID() int {
	id := e.nextID
	e.nextID++
	return id
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
