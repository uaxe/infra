package assert_test

import (
	"fmt"
	"github.com/uaxe/infra/assert"
)

func ExampleAssert() {
	assert.Assert(false, func() {
		fmt.Printf("%#v\n", false)
	})
	assert.Assert(true, func() {
		fmt.Printf("%#v\n", true)
	})
	// Output:
	// false
}

func ExampleAssertE() {
	assert.AssertE(fmt.Errorf("assert error"), func(e error) {
		fmt.Printf("%s\n", e)
	})
	// Output:
	// assert error
}
