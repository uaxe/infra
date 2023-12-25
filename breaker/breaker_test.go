package breaker

import (
	"fmt"
)

func ExampleBreaker() {
	b := NewBreaker()
	fmt.Printf("%#v\n", len(b.Name()) > 0)
	_, err := b.Allow()
	fmt.Printf("%#v\n", err)
	// Output:
	// true
	// <nil>
}
