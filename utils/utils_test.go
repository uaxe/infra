package utils_test

import (
	"fmt"
	"github.com/uaxe/infra/utils"
)

func ExampleAssert() {
	utils.Assert(false, func() {
		fmt.Printf("%#v\n", false)
	})
	utils.Assert(true, func() {
		fmt.Printf("%#v\n", true)
	})
	// Output:
	// false
}

func ExampleAssertE() {
	utils.AssertE(fmt.Errorf("assert error"), func(e error) {
		fmt.Printf("%s\n", e)
	})
	// Output:
	// assert error
}

func ExampleRainbow() {
	rainbow := utils.Rainbow("infra")
	fmt.Printf("%#v\n", rainbow)
	// Output:
	// "\x1b[0;0;31mi\x1b[0m\x1b[0;0;33mn\x1b[0m\x1b[0;0;32mf\x1b[0m\x1b[0;0;36mr\x1b[0m\x1b[0;0;34ma\x1b[0m"
}
