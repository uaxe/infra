package utils_test

import (
	"fmt"

	"github.com/uaxe/infra/utils"
)

func ExampleValue() {

	fmt.Printf("%#v\n", utils.Value(false, true))
	fmt.Printf("%#v\n", utils.Value(1, 0))
	fmt.Printf("%#v\n", utils.Value(0, 0))

	fmt.Printf("%#v\n", utils.AssertV(false, "zkep", "infra"))
	fmt.Printf("%#v\n", utils.AssertV(true, "zkep", "infra"))
	fmt.Printf("%#v\n", utils.AssertV(true, func() error { return nil },
		func() error {
			return fmt.Errorf("not nil")
		})())
	// Output:
	// true
	// 1
	// 0
	// "infra"
	// "zkep"
	// <nil>
}
