package xlsxutil

import (
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
	"github.com/tealeg/xlsx"

	"github.com/onelrdm/reflectutil"
)

type StructContext struct {
	cfg *reflectutil.Config
}

func NewStructContext(cfg *reflectutil.Config) *StructContext {
	return &StructContext{cfg: cfg}
}

func (r StructContext) Config() *reflectutil.Config {
	return r.cfg
}

func (r StructContext) NewEncoder(typ reflect2.Type) reflectutil.Encoder {
	kind := typ.Kind()
	switch kind {
	case reflect.Int:
		return intCodec
	case reflect.String:
		return stringCodec
	case reflect.Ptr:
		typ := typ.(*reflect2.UnsafePtrType)
		encoder := r.NewEncoder(typ.Elem())
		return &reflectutil.DereferenceEncoder{Encoder: encoder}
	default:
		return &AnyCodec{typ: typ}
	}
}

type structEncoder struct {
	*reflectutil.StructDescriptor
}

func (r *structEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	row := writer.(*xlsx.Row)
	for _, binding := range r.Fields {
		cell := row.AddCell()
		w := CellWriter{Cell: cell}
		fieldPtr := binding.Field.UnsafeGet(ptr)
		binding.Encoder.Encode(fieldPtr, &w)
	}
}
