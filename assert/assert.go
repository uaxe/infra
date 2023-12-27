package assert

import (
	"reflect"
)

func Assert(guard bool, fn func()) {
	if !guard {
		fn()
	}
}

func FindE(errs ...any) error {
	for _, e := range errs {
		if e == nil {
			continue
		}
		t := reflect.TypeOf(e)
		switch t.Kind() {
		case reflect.Ptr:
			if e != nil {
				return e.(error)
			}
		case reflect.Func:
			if t.NumOut() < 1 {
				panic("Setup  requires a function with the signature 'func() error'")
			}
			if t.Out(t.NumOut()-1) != reflect.TypeOf((*error)(nil)).Elem() {
				panic("Setup  requires a function with the signature 'func() error'")
			}
			in := make([]reflect.Value, 0, t.NumIn())
			for i := 0; i < t.NumIn(); i++ {
				in = append(in, reflect.New(t.In(i).Elem()))
			}
			result := reflect.ValueOf(e).Call(in)[t.NumOut()-1].Interface()
			if result != nil {
				return result.(error)
			}
		default:
		}
	}
	return nil
}
