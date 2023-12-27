package zreflect

import (
	"reflect"
)

func IsStructPtr(x any) bool {
	return reflect.ValueOf(x).Kind() == reflect.Ptr &&
		reflect.ValueOf(x).Elem().Kind() == reflect.Struct
}

func IsFunction(x any) bool {
	return reflect.ValueOf(x).Kind() == reflect.Func
}

func IsStruct(x any) bool {
	return reflect.ValueOf(x).Kind() == reflect.Struct
}

func HasElements(typ reflect.Type) bool {
	kind := typ.Kind()
	return kind == reflect.Ptr || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}

func TypeAndValue(x any) (reflect.Type, reflect.Value) {
	t, v := reflect.TypeOf(x), reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	return t, v
}

func TypeAndKind(x any) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(x)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}
