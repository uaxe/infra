package zreflect_test

import (
	"fmt"

	"github.com/uaxe/infra/zreflect"
)

type Src struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

type Dst struct {
	Name   string `json:"name,omitempty"`
	Age    int    `json:"age,omitempty"`
	Gender string `json:"gender,omitempty"`
	Job    string `json:"job,omitempty"`
}

func ExampleMergeStruct() {
	src := &Src{Name: "uaxe", Age: 1}
	dst := &Dst{}
	err := zreflect.MergeStruct(src, dst)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", dst.Name == "uaxe")
	fmt.Printf("%#v\n", dst.Age == 1)
	// Output:
	// <nil>
}
