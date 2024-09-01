package event_test

import (
	"fmt"
	"time"

	"atomicgo.dev/event"
)

type Player struct {
	Name string
}

var PlayerJoinEvent = event.New[Player]()

func Example_demo() {
	// Listen to the event as many times as you want
	PlayerJoinEvent.Listen(func(p Player) {
		fmt.Printf("Player %q joined the game\n", p.Name)
	})

	PlayerJoinEvent.Listen(func(p Player) {
		// Do something else
	})

	// ...

	// Trigger the event somewhere - can be in a different function or package
	PlayerJoinEvent.Trigger(Player{Name: "Marvin"})
	PlayerJoinEvent.Trigger(Player{Name: "Bob"})
	PlayerJoinEvent.Trigger(Player{Name: "Alice"})

	// Keep the program alive
	time.Sleep(time.Second)

	// Output:
	// Player "Marvin" joined the game
	// Player "Bob" joined the game
	// Player "Alice" joined the game
}
