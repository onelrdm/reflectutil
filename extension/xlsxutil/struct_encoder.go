package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/reflectutil"
	"github.com/tealeg/xlsx"
	"reflect"
	"unsafe"
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
	default:
		return &AnyCodec{valType: typ}
	}
}

type structEncoder struct {
	typ reflect2.Type
}

func (r *structEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	ctx := NewStructContext(&reflectutil.Config{TaggedFieldOnly: true})
	sd := reflectutil.DescribeStruct(ctx, r.typ)
	row := writer.(*xlsx.Row)
	for _, binding := range sd.Fields {
		cell := row.AddCell()
		w := CellWriter{Cell: cell}
		fieldPtr := binding.Field.UnsafeGet(ptr)
		binding.Encoder.Encode(fieldPtr, &w)
	}
}
