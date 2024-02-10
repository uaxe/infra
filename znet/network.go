package znet

import (
	"fmt"
	"net"
)

func AvailablePort() (int, error) {
	address := fmt.Sprintf(":%d", 0)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return 0, err
	}
	defer func() { _ = listener.Close() }()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func IsPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer func() { _ = listener.Close() }()
	return true
}

func DefaultOrAvailablePort(port int) (int, error) {
	if !IsPortAvailable(port) {
		return AvailablePort()
	}
	return port, nil
}

func DefaultOrAvailablePortWithFunc(port int, fn func(err error)) int {
	if IsPortAvailable(port) {
		return port
	}
	port, err := AvailablePort()
	if err != nil {
		fn(err)
	}
	return port
}
