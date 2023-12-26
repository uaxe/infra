package binder

type (
	Binder interface {
		Name() string
		Binding(src any, dst any) error
	}
)

func Bindings(bs ...Binder) {

}
