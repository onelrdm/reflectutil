package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/reflectutil"
	"reflect"
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
		return &structEncoder{typ: typ}
	case reflect.Slice:
		sliceType := typ.(*reflect2.UnsafeSliceType)
		encoder := r.NewEncoder(sliceType.Elem())
		return &sliceEncoder{sliceType: sliceType, elemEncoder: encoder}
	default:
		return &AnyCodec{valType: typ}
	}
}
