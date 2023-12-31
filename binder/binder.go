package binder

type (
	Binder interface {
		Name() string
		Binding(src any, dst any) error
	}
)

func Bindings(src, dst any, bs ...Binder) error {
	for _, bind := range bs {
		if err := bind.Binding(src, dst); err != nil {
			return err
		}
	}
	return nil
}
