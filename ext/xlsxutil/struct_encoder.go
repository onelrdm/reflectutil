package xlsxutil

import (
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
	"github.com/tealeg/xlsx"

	"github.com/onelrdm/reflectutil"
)

func GetStructFieldEncoder(typ reflect2.Type, field reflect2.StructField) reflectutil.Encoder {
	kind := typ.Kind()
	switch kind {
	case reflect.Int:
		return intCodec
	case reflect.String:
		return stringCodec
	case reflect.Ptr:
		typ := typ.(*reflect2.UnsafePtrType)
		encoder := GetStructFieldEncoder(typ.Elem(), field)
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
		binding.Encoder.Encode(ptr, &w)
	}
}
