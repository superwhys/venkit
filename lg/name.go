package lg

import (
	"reflect"
	"runtime"
)

func FuncName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func StructName(i any) string {
	tPtr := reflect.TypeOf(i)

	if tPtr.Kind() == reflect.Ptr {
		tPtr = tPtr.Elem()
	}

	return tPtr.Name()
}
