package hook

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrNilProvide       = errors.New("nil provide")
	ErrNotMatchProvider = errors.New("not match provider")
)

type Hook[T any] interface {
	Register(p T, index ...int) error

	Get(match func(p T) bool) (T, error)

	Update(match func(T) bool, p T) (T, error)
}

type providerWrap[T any] struct {
	provide T
	index   int
}

type IHook[T any] struct {
	lock      sync.Mutex
	providers []*providerWrap[T]
}

func (h *IHook[T]) Register(provide T, indexes ...int) error {
	if reflect.TypeOf(provide) == nil {
		return ErrNilProvide
	}
	index := 0
	if len(indexes) > 0 {
		index = indexes[0]
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	insertIndex := len(h.providers)
	for i, wrap := range h.providers {
		if wrap.index <= index {
			insertIndex = i
			break
		}
	}
	wrap := &providerWrap[T]{provide: provide, index: index}
	curr := append(append([]*providerWrap[T]{}, h.providers[0:insertIndex]...), wrap)
	h.providers = append(curr, h.providers[insertIndex:]...)
	return nil
}

func (h *IHook[T]) Get(match func(T) bool) (m T, e error) {
	for _, wrap := range h.providers {
		if match(wrap.provide) {
			return wrap.provide, nil
		}
	}
	e = ErrNotMatchProvider
	return
}

func (h *IHook[T]) Update(match func(T) bool, provide T) (m T, e error) {
	if reflect.TypeOf(provide) == nil {
		e = ErrNilProvide
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	index := 0
	for i, wrap := range h.providers {
		if match(wrap.provide) {
			index = i
			m = wrap.provide
			break
		}
	}

	if reflect.TypeOf(m) == nil {
		e = ErrNotMatchProvider
		return
	}

	wrap := &providerWrap[T]{provide: provide, index: index}
	h.providers[index] = wrap

	return
}
