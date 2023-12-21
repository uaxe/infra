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
