package utils

func SliceAppend[M ~[]V, V any](m, n M) []V {
	if len(m) == 0 {
		m = make([]V, len(n))
	}
	m = append(m, n...)
	return m
}

func SliceDistinct[M ~[]V, V any](m M, fn func(V) any) []V {
	if len(m) == 0 {
		return m
	}
	ms := make(map[any]V, len(m))
	keys := make([]any, 0, len(m))
	for _, v := range m {
		key := fn(v)
		if _, ok := ms[key]; ok {
			continue
		} else {
			ms[key] = v
			keys = append(keys, key)
		}
	}
	mv := make([]V, 0, len(ms))
	for i := range keys {
		key := keys[i]
		mv = append(mv, ms[key])
	}
	return mv
}

func SliceFind[M ~[]V, V any](m M, f func(V) bool) int {
	for k, v := range m {
		if f(v) {
			return k
		}
	}
	return -1
}

func SliceMap[M ~[]V, V any](m M, f func(V) bool) {
	for _, v := range m {
		if !f(v) {
			return
		}
	}
	return
}
