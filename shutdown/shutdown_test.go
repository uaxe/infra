package shutdown

import (
	"log"
	"syscall"
	"testing"
	"time"
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
	pid := syscall.Getpid()
	go func() {
		time.Sleep(time.Second)
		_ = syscall.Kill(pid, syscall.SIGINT)
	}()
	h.WatchSignal()
}
