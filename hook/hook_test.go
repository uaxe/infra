package hook_test

import (
	"errors"
	"testing"

	"github.com/uaxe/infra/hook"
)

func Test_HookAny(t *testing.T) {
	var (
		defaultHook = &hook.IHook[any]{}
		Register    = defaultHook.Register
		Get         = defaultHook.Get
	)
	provide := 1
	if err := Register(provide); err != nil {
		t.FailNow()
	}
	curr, err := Get(func(p any) bool {
		return provide == p
	})
	if err != nil {
		t.Logf("err: %+v", err)
		t.FailNow()
	}
	if curr != provide {
		t.Logf("curr: %+v", curr)
		t.FailNow()
	}
}

func Test_Hook(t *testing.T) {
	type Server struct {
		Addr string
		Port int
	}
	var (
		defaultHook = &hook.IHook[*Server]{}
		Register    = defaultHook.Register
		Get         = defaultHook.Get
		Update      = defaultHook.Update
	)

	provide := &Server{Port: 80}

	if err := Register(provide); err != nil {
		t.FailNow()
	}

	if err := Register(&Server{Port: 81}); err != nil {
		t.FailNow()
	}

	curr, err := Get(func(s *Server) bool {
		return provide.Port == s.Port
	})

	if err != nil {
		t.Logf("err: %+v", err)
		t.FailNow()
	}

	if curr.Port != provide.Port {
		t.Logf("curr: %+v", curr)
		t.FailNow()
	}

	_, err = Get(func(s *Server) bool { return provide.Addr == "127.0.0.1" })
	if !errors.Is(err, hook.ErrNotMatchProvider) {
		t.Logf("err: %+v", err)
		t.FailNow()
	}

	newProvide := &Server{Port: 8000}
	oldProvide, err := Update(func(s *Server) bool {
		return provide.Port == s.Port
	}, newProvide)
	if err != nil {
		t.Logf("err: %+v", err)
		t.FailNow()
	}

	if oldProvide.Port != provide.Port {
		t.Logf("oldProvide: %+v", curr)
		t.FailNow()
	}

}
