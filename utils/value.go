package utils

import (
	"reflect"
)

func Value[V any](v, d V) V {
	val := reflect.ValueOf(v)
	if val.IsZero() {
		return d
	}
	return v
}

func AssertV[V any](pass bool, v, d V) V {
	if pass {
		return v
	}
	return d
}
