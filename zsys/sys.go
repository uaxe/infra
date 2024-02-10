package zsys

import (
	"reflect"
	"runtime"
)

func NameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func MethodName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown caller"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}
