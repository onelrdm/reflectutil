package xlsxutil

import (
	"reflect"

	"github.com/modern-go/reflect2"

	"github.com/onelrdm/reflectutil"
)

var stringCodec = &StringCodec{}
var intCodec = &IntCodec{}

func GetEncoder(typ reflect2.Type, field reflect2.StructField) reflectutil.Encoder {
	kind := typ.Kind()
	switch kind {
	case reflect.Int:
		return intCodec
	case reflect.String:
		return stringCodec
	case reflect.Struct:
		sd := reflectutil.DescribeStruct(typ.(*reflect2.UnsafeStructType), GetStructFieldEncoder, DefaultOption)
		return &structEncoder{sd}
	case reflect.Ptr:
		typ := typ.(*reflect2.UnsafePtrType)
		encoder := GetEncoder(typ.Elem(), field)
		return &reflectutil.DereferenceEncoder{Encoder: encoder}
	case reflect.Slice:
		typ := typ.(*reflect2.UnsafeSliceType)
		encoder := GetEncoder(typ.Elem(), field)
		return &sliceEncoder{typ: typ, encoder: encoder}
	default:
		return &AnyCodec{typ: typ}
	}
}
