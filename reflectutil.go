// Package reflectutil implements reflect2 utility functions.
package reflectutil

import (
	"github.com/modern-go/reflect2"
	"reflect"
)

func TypeOf(obj interface{}) reflect2.Type {
	typ := reflect2.TypeOf(obj)
	kind := typ.Kind()
	if kind == reflect.Ptr {
		ptrType := typ.(*reflect2.UnsafePtrType)
		typ = ptrType.Elem()
	}
	return typ
}
