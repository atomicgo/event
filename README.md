<!--



┌───────────────────────────────────────────────────────────────────┐
│                                                                   │
│                          IMPORTANT NOTE                           │
│                                                                   │
│               This file is automatically generated                │
│           All manual modifications will be overwritten            │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘



-->

<h1 align="center">AtomicGo | event</h1>

<p align="center">
<img src="https://img.shields.io/endpoint?url=https%3A%2F%2Fatomicgo.dev%2Fapi%2Fshields%2Fevent&style=flat-square" alt="Downloads">

<a href="https://github.com/atomicgo/event/releases">
<img src="https://img.shields.io/github/v/release/atomicgo/event?style=flat-square" alt="Latest Release">
</a>

<a href="https://codecov.io/gh/atomicgo/event" target="_blank">
<img src="https://img.shields.io/github/actions/workflow/status/atomicgo/event/go.yml?style=flat-square" alt="Tests">
</a>

<a href="https://codecov.io/gh/atomicgo/event" target="_blank">
<img src="https://img.shields.io/codecov/c/gh/atomicgo/event?color=magenta&logo=codecov&style=flat-square" alt="Coverage">
</a>

<a href="https://codecov.io/gh/atomicgo/event">
<!-- unittestcount:start --><img src="https://img.shields.io/badge/Unit_Tests-2-magenta?style=flat-square" alt="Unit test count"><!-- unittestcount:end -->
</a>

<a href="https://opensource.org/licenses/MIT" target="_blank">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>
  
<a href="https://goreportcard.com/report/github.com/atomicgo/event" target="_blank">
<img src="https://goreportcard.com/badge/github.com/atomicgo/event?style=flat-square" alt="Go report">
</a>   

</p>

---

<p align="center">
<strong><a href="https://pkg.go.dev/atomicgo.dev/event#section-documentation" target="_blank">Documentation</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CONTRIBUTING.md" target="_blank">Contributing</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CODE_OF_CONDUCT.md" target="_blank">Code of Conduct</a></strong>
</p>

---

<p align="center">
  <img src="https://raw.githubusercontent.com/atomicgo/atomicgo/main/assets/header.png" alt="AtomicGo">
</p>

<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get atomicgo.dev/event</pre></h3>
<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# event

```go
import "atomicgo.dev/event"
```

Package event provides a generic and thread\-safe event system for Go. It allows multiple listeners to subscribe to events carrying data of any type. Listeners can be added and notified when events are triggered, and the event can be closed to prevent further operations.





```go
package main

import (
	"fmt"
	"time"

	"atomicgo.dev/event"
)

func delay() {
	time.Sleep(time.Millisecond * 10)
}

type Player struct {
	Name string
}

// Create a new event
var PlayerJoinEvent = event.New[Player]()

func main() {
	// Listen to the event as many times as you want
	PlayerJoinEvent.Listen(func(p Player) {
		fmt.Printf("Player %q joined the game\n", p.Name)
	})

	PlayerJoinEvent.Listen(func(_ Player) {
		// Do something else
	})

	// ...

	// Trigger the event somewhere - can be in a different function or package
	PlayerJoinEvent.Trigger(Player{Name: "Marvin"})
	delay() // delay for deterministic output
	PlayerJoinEvent.Trigger(Player{Name: "Bob"})
	delay() // delay for deterministic output
	PlayerJoinEvent.Trigger(Player{Name: "Alice"})

	// Keep the program alive
	time.Sleep(time.Second)

}
```

#### Output

```
Player "Marvin" joined the game
Player "Bob" joined the game
Player "Alice" joined the game
```



## Index

- [Variables](<#variables>)
- [type Event](<#Event>)
  - [func New\[T any\]\(\) \*Event\[T\]](<#New>)
  - [func \(e \*Event\[T\]\) Close\(\)](<#Event[T].Close>)
  - [func \(e \*Event\[T\]\) Listen\(f func\(T\)\) error](<#Event[T].Listen>)
  - [func \(e \*Event\[T\]\) Trigger\(value T\) error](<#Event[T].Trigger>)


## Variables

<a name="ErrEventClosed"></a>ErrEventClosed is returned when an operation is attempted on a closed event.

```go
var ErrEventClosed = errors.New("event is closed")
```

<a name="Event"></a>
## type [Event](<https://github.com/atomicgo/event/blob/main/event.go#L13-L17>)

Event represents a generic, thread\-safe event system that can handle multiple listeners. The type parameter T specifies the type of data that the event carries when triggered.

```go
type Event[T any] struct {
    // contains filtered or unexported fields
}
```

<a name="New"></a>
### func [New](<https://github.com/atomicgo/event/blob/main/event.go#L20>)

```go
func New[T any]() *Event[T]
```

New creates and returns a new Event instance for the specified type T.

<a name="Event[T].Close"></a>
### func \(\*Event\[T\]\) [Close](<https://github.com/atomicgo/event/blob/main/event.go#L74>)

```go
func (e *Event[T]) Close()
```

Close closes the event system, preventing any new listeners from being added or events from being triggered. After calling Close, any subsequent calls to Trigger or Listen will return ErrEventClosed. Existing listeners are removed, and resources are cleaned up.





```go
package main

import (
	"fmt"
	"time"

	"atomicgo.dev/event"
)

func delay() {
	time.Sleep(time.Millisecond * 10)
}

func main() {
	// Create a new event
	exampleEvent := event.New[int]()

	// Listen to the event
	exampleEvent.Listen(func(v int) {
		fmt.Println(v)
	})

	// Trigger the event
	exampleEvent.Trigger(1)
	delay() // delay for deterministic output
	exampleEvent.Trigger(2)
	delay() // delay for deterministic output
	exampleEvent.Trigger(3)

	// Time for listeners to process the event
	delay()

	// Close the event
	exampleEvent.Close()

	// Trigger the event again
	exampleEvent.Trigger(4)
	delay() // delay for deterministic output
	exampleEvent.Trigger(5)
	delay() // delay for deterministic output
	exampleEvent.Trigger(6)

	// Keep the program alive
	time.Sleep(time.Second)

}
```

#### Output

```
1
2
3
```



<a name="Event[T].Listen"></a>
### func \(\*Event\[T\]\) [Listen](<https://github.com/atomicgo/event/blob/main/event.go#L58>)

```go
func (e *Event[T]) Listen(f func(T)) error
```

Listen registers a new listener callback function for the event. The listener will be invoked with the event's data whenever Trigger is called. Returns ErrEventClosed if the event has been closed.

<a name="Event[T].Trigger"></a>
### func \(\*Event\[T\]\) [Trigger](<https://github.com/atomicgo/event/blob/main/event.go#L27>)

```go
func (e *Event[T]) Trigger(value T) error
```

Trigger notifies all registered listeners by invoking their callback functions with the provided value. It runs each listener in a separate goroutine and waits for all listeners to complete. Returns ErrEventClosed if the event has been closed.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)
