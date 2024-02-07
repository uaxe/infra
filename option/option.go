package option

type Option[T any] interface {
	Apply(T)
}

type applyOption[T any] struct {
	f func(T)
}

func (x *applyOption[T]) Apply(do T) {
	x.f(do)
}

func NewApplyOption[T any](f func(T)) Option[T] {
	return &applyOption[T]{f: f}
}
