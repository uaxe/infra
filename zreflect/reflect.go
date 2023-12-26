package zreflect

import (
	"reflect"
)

func IsStructPtr(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr &&
		reflect.ValueOf(value).Elem().Kind() == reflect.Struct
}

func IsFunction(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Func
}

func IsStruct(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Struct
}

func HasElements(typ reflect.Type) bool {
	kind := typ.Kind()
	return kind == reflect.Ptr || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}
