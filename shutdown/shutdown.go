package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var _ Hook = (*hook)(nil)

type Hook interface {
	WithSignals(signals ...syscall.Signal)

	Add(func())

	WatchSignal()
}

type hook struct {
	signal   chan os.Signal
	handlers []func()
}

func NewHook() Hook {
	h := &hook{signal: make(chan os.Signal, 1)}
	h.Add(func() { signal.Stop(h.signal) })
	h.WithSignals(syscall.SIGINT, syscall.SIGTERM)
	return h
}

func (h *hook) WithSignals(signals ...syscall.Signal) {
	for _, s := range signals {
		signal.Notify(h.signal, s)
	}
}

func (h *hook) Add(f func()) {
	h.handlers = append([]func(){f}, h.handlers...)
}

func (h *hook) WatchSignal() {
	<-h.signal
	for _, handler := range h.handlers {
		handler()
	}
}
