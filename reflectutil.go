// Package reflectutil implements reflect2 utility functions.
package reflectutil

import (
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

// Encoder is an internal type registered to cache as needed.
type Encoder interface {
	Encode(ptr unsafe.Pointer, writer interface{})
}

func TypeOf(obj interface{}) reflect2.Type {
	typ := reflect2.TypeOf(obj)
	kind := typ.Kind()
	if kind == reflect.Ptr {
		ptrType := typ.(*reflect2.UnsafePtrType)
		typ = ptrType.Elem()
	}
	return typ
}

func EncoderOf(encoder Encoder) interface{} {
	if encoder, ok := encoder.(*DereferenceEncoder); ok {
		return encoder.Encoder
	}
	return encoder
}
