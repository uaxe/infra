package assert

func Assert(guard bool, fn func()) {
	if !guard {
		fn()
	}
}

func AssertE(err error, fn func(e error)) {
	Assert(err == nil, func() { fn(err) })
}
