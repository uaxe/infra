package shutdown

import (
	"log"
	"testing"
)

type Server struct {
	id int
}

func (s *Server) Close() {
	log.Print("shutting down: ", s.id)
}

func TestShutdown(t *testing.T) {
	h := NewHook()
	server1 := &Server{id: 1}
	h.Add(func() { server1.Close() })
	server2 := &Server{id: 2}
	h.Add(func() { server2.Close() })

	h.WatchSignal()
}
