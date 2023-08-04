package utils

type Number interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64
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
