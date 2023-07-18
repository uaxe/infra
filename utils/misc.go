package utils

import "fmt"

func Assert(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func Setup(errs ...error) {
	for _, e := range errs {
		Assert(e == nil, fmt.Sprintf("setup panic: %s", e.Error()))
	}
}

func SetupF(f func(error), errs ...error) {
	for _, e := range errs {
		f(e)
	}
}

func SetupBreak(f func(error) bool, errs ...error) error {
	for _, e := range errs {
		if f(e) {
			return e
		}
	}
	return nil
}

const (
	_ = 1.0 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)
