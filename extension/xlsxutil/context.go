package xlsxutil

import (
	"reflect"

	"github.com/modern-go/reflect2"

	"github.com/onelrdm/reflectutil"
)

type Context struct {
	cfg *reflectutil.Config
}

func NewContext(cfg *reflectutil.Config) *Context {
	return &Context{cfg: cfg}
}

func (r Context) Config() *reflectutil.Config {
	return r.cfg
}

var stringCodec = &StringCodec{}
var intCodec = &IntCodec{}

func (r Context) NewEncoder(typ reflect2.Type) reflectutil.Encoder {
	kind := typ.Kind()
	switch kind {
	case reflect.Int:
		return intCodec
	case reflect.String:
		return stringCodec
	case reflect.Struct:
		ctx := NewStructContext(&reflectutil.Config{TaggedFieldOnly: true})
		sd := reflectutil.DescribeStruct(ctx, typ)
		return &structEncoder{sd}
	case reflect.Ptr:
		typ := typ.(*reflect2.UnsafePtrType)
		encoder := r.NewEncoder(typ.Elem())
		return &reflectutil.DereferenceEncoder{Encoder: encoder}
	case reflect.Slice:
		typ := typ.(*reflect2.UnsafeSliceType)
		encoder := r.NewEncoder(typ.Elem())
		return &sliceEncoder{typ: typ, encoder: encoder}
	default:
		return &AnyCodec{typ: typ}
	}
}
