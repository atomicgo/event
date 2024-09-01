package event

import "sync"

// Event represents an event system that can handle multiple listeners.
type Event[T any] struct {
	listeners []chan T
	mu        sync.Mutex
}

func New[T any]() *Event[T] {
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
			l <- value
		}(listener)
	}
}

// Listen gets called when the event is triggered.
func (e *Event[T]) Listen(f func(T)) {
	ch := make(chan T)

	e.mu.Lock()
	e.listeners = append(e.listeners, ch)
	e.mu.Unlock()

	go func() {
		for v := range ch {
			f(v)
		}
	}()
}

func (e *Event[T]) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, listener := range e.listeners {
		close(listener)
	}

	e.listeners = nil
}
