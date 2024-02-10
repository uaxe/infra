package telnet

import "net"

func Dial(addr string, options ...Option) (conn *Connection, err error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	conn = NewConnection(c, options)
	return
}
