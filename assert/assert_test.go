package assert_test

import (
	"errors"
	"fmt"

	"github.com/uaxe/infra/assert"
)

func ExampleSetup() {

	var targetErr error

	fn := func(i int) error {
		targetErr = fmt.Errorf("err:%d", i)
		return targetErr
	}
	assert.Setup(func(err error) {
		fmt.Printf("%#v\n", errors.Is(err, targetErr))
	},
		nil,
		func() (string, error) {
			return "1", nil
		},
		fn(1),
		func() error {
			targetErr = errors.New("string error")
			return targetErr
		},
	)

	// Output:
	// true
}

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
