package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var _ Hook = (*hook)(nil)

// Hook a graceful shutdown hook, default with signals of SIGINT and SIGTERM
type Hook interface {
	// WithSignals add more signals into hook
	WithSignals(signals ...syscall.Signal)

	// ADD register shutdown handles
	Add(func())

	// WatchSignal with signal
	WatchSignal()
}

type hook struct {
	signal   chan os.Signal
	handlers []func()
}

// NewHook create a Hook instance
func NewHook() Hook {
	h := &hook{
		signal: make(chan os.Signal, 1),
	}
	h.Add(func() {
		signal.Stop(h.signal)
	})
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
