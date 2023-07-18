package utils

type Number interface {
	uint8 | int | int64 | float32 | float64
}

func Max[V Number](x, y V) V {
	if x > y {
		return x
	}
	return y
}

func Min[V Number](x, y V) V {
	if x < y {
		return x
	}
	return y
}
