package zreflect

import (
	"fmt"
	"reflect"
	"strings"
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

func GetPointer(v any) any {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		return v
	}
	return reflect.New(vv.Type()).Interface()
}

func Indirect(v reflect.Value) reflect.Value {
	return reflect.Indirect(v)
}

func TypeOf(t any) (reflect.Type, bool) {
	v := reflect.TypeOf(t)
	isPtr := false
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		isPtr = true
	}
	return v, isPtr
}

func ValueOf(t any) (reflect.Value, bool) {
	v := reflect.ValueOf(t)
	isPtr := false
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		isPtr = true
	}
	return v, isPtr
}

func SetValue(rcvr any, fieldName string, value any) bool {
	v := reflect.ValueOf(rcvr)
	t := reflect.TypeOf(rcvr)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return false
	}

	if f, ok := t.Elem().FieldByName(fieldName); ok {
		fvalue := v.Elem().FieldByIndex(f.Index)

		if !fvalue.CanSet() {
			return false
		}

		switch fvalue.Kind() {
		case reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int64,
			reflect.Int, reflect.Uint8, reflect.Uint, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32,
			reflect.Float64, reflect.Bool, reflect.Array, reflect.Slice, reflect.Map:
			fvalue.Set(reflect.ValueOf(value))
		case reflect.String:
			fvalue.SetString(fmt.Sprintf("%s", value))
		case reflect.Chan:
			return fvalue.TrySend(reflect.ValueOf(value))
		default:
		}
		return true
	}
	return false
}

func GetValue[S any](rcvr any, fieldName string) (S, bool) {
	v := reflect.ValueOf(rcvr)
	var ret S

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return ret, false
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ret, false
	}

	if f, ok := v.Type().FieldByName(fieldName); ok {
		fvalue := v.FieldByIndex(f.Index)
		if fvalue.CanInterface() {
			ret, ok = fvalue.Interface().(S)
			return ret, ok
		}
	}
	return ret, false
}

func GetMethods(typ reflect.Type) map[string]*reflect.Method {
	methods := make(map[string]*reflect.Method)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		methods[mname] = &method
	}
	return methods
}

func CallMethodByName(rcvr any, methodName string, params ...any) ([]reflect.Value, error) {
	typ := reflect.TypeOf(rcvr)
	kind := typ.Kind()
	if kind == reflect.Pointer {
		kind = typ.Elem().Kind()
	}

	if kind != reflect.Struct && kind != reflect.Interface {
		return nil, fmt.Errorf("param 'rcvr' should be struct or interface type.")
	}
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mname := method.Name

		if strings.EqualFold(mname, methodName) {
			return callMethod(rcvr, &method, params...)
		}
	}
	return nil, fmt.Errorf("method name '%s' not found", methodName)
}

func callMethod(rcvr any, method *reflect.Method, params ...any) ([]reflect.Value, error) {
	paramSize := len(params) + 1
	paramValues := make([]reflect.Value, paramSize)
	paramValues[0] = reflect.ValueOf(rcvr)
	i := 1
	for _, v := range params {
		paramValues[i] = reflect.ValueOf(v)
		i++
	}
	returnValues := method.Func.Call(paramValues)
	return returnValues, nil
}
