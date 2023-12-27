package assert

func Setup(assert func(err error), funcs ...any) {
	if assert == nil {
		return
	}
	if e := FindE(funcs...); e != nil {
		assert(e)
	}
}

func PanicE(errs ...any) {
	Setup(Panic, errs...)
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
