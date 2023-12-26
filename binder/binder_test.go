package binder_test

import (
	"fmt"
	"net/http"

	"github.com/uaxe/infra/binder"
)

type DstHeader struct {
	ContentType   string `header:"Content-Type"`
	ContentLength string `json:"Content-Length"`
}

func ExampleHeader() {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Content-Length", "1024")
	var dst DstHeader
	err := binder.Header.Binding(h, &dst)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", dst.ContentType)
	fmt.Printf("%#v\n", dst.ContentLength)
	// Output:
	// <nil>
	// "application/json"
	// ""
}
