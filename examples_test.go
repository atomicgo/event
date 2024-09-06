package event_test

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

func Example_demo() {
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

	// Output:
	// Player "Marvin" joined the game
	// Player "Bob" joined the game
	// Player "Alice" joined the game
}

func ExampleEvent_Close() {
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

	// Output:
	// 1
	// 2
	// 3
}
