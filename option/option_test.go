package option_test

import (
	"testing"

	"github.com/uaxe/infra/option"
)

type options struct {
	cache bool
	size  int
}

var defaultOptions = options{}

func Size(size int) option.Option[*options] {
	return option.NewApplyOption[*options](func(o *options) {
		o.size = size
	})
}

func Cache() option.Option[*options] {
	return option.NewApplyOption[*options](func(o *options) {
		o.cache = true
	})
}

type Service struct {
	opts options
}

func New(opt ...option.Option[*options]) *Service {
	opts := defaultOptions
	for _, o := range opt {
		o.Apply(&opts)
	}

	return &Service{
		opts: opts,
	}
}

func TestService(t *testing.T) {
	s := New(Cache(), Size(10))
	if !s.opts.cache {
		t.Logf("%+v", s.opts)
		t.FailNow()
	}
	if s.opts.size != 10 {
		t.Logf("%+v", s.opts)
		t.FailNow()
	}
}
