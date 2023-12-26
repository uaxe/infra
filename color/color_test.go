package color_test

import (
	"fmt"

	"github.com/uaxe/infra/color"
)

func ExampleRainbow() {
	rainbow := color.Rainbow("infra")
	fmt.Printf("%#v\n", rainbow)
	// Output:
	// "\x1b[0;0;31mi\x1b[0m\x1b[0;0;33mn\x1b[0m\x1b[0;0;32mf\x1b[0m\x1b[0;0;36mr\x1b[0m\x1b[0;0;34ma\x1b[0m"
}
