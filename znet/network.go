package znet

import (
	"fmt"
	"net"
)

func AvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func IsPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

func DefaultOrAvailablePort(defaultPort int) int {
	if IsPortAvailable(defaultPort) {
		return defaultPort
	}
	port, _ := AvailablePort()
	return port
}
