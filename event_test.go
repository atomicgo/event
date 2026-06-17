package event

import "testing"

func TestZeroValue(t *testing.T) {
	var ev Event[string]

	// Can listen to zero value without panicing
	ev.Listen(func(v string) { println(v) })
}
