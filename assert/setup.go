package assert

func Setup(fn func(err error), errs ...error) {
	for _, e := range errs {
		AssertE(e, fn)
	}
}

func SetupPanic(errs ...error) {
	for _, e := range errs {
		AssertE(e, func(e error) {
			panic(e)
		})
	}
}
