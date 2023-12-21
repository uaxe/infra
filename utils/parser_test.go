package utils_test

import (
	"fmt"
	"github.com/uaxe/infra/utils"
	"net/http"
)

type DstHeader struct {
	ContentType   string `header:"Content-Type"`
	ContentLength string `json:"Content-Length"`
}

func ExampleHeaderParser() {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Content-Length", "1024")
	var dst DstHeader
	err := utils.HeaderParser.Parse(h, &dst)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", dst.ContentType)
	fmt.Printf("%#v\n", dst.ContentLength)
	// Output:
	// <nil>
	// "application/json"
	// ""
}
