package znet_test

import (
	"fmt"
	"net"

	"github.com/uaxe/infra/znet"
)

func ExampleDefaultOrAvailablePort() {
	defaultPort := 8007
	port, err := znet.DefaultOrAvailablePort(defaultPort)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}
	defer func() { _ = listener.Close() }()

	port = znet.DefaultOrAvailablePortWithFunc(defaultPort, func(err error) {
		fmt.Printf("%#v\n", err)
	})
	fmt.Printf("%#v\n", port != defaultPort)
	// Output:
	// <nil>
	// 8007
	// true
}
