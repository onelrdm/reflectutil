package reflectutil

import (
	"bytes"
	"github.com/modern-go/reflect2"
	"testing"
	"time"
	"unsafe"
)

func TestDescribeStruct(t *testing.T) {
	cfg := Config{
		TagKey: "",
	}

	type Embed2 struct {
		ID2   int    `reflect:"id2"`
		Name2 string `reflect:"name2"`
	}
	type Embed struct {
		Embed2
		ID   int    `reflect:"id"`
		Name string `reflect:"name"`
	}
	type Val struct {
		Embed
		Data *time.Time `reflect:"data"`
	}
	val := Val{
		Embed: Embed{
			Embed2: Embed2{
				ID2:   2,
				Name2: "name2",
			},
			ID:   1,
			Name: "name",
		},
		Data: nil,
	}
	desc := DescribeStruct(&cfg, reflect2.TypeOf(val))
	t.Logf("%+v\n", desc.Type)
	for _, binding := range desc.Fields {
		t.Logf("%+v\n", binding)
	}

	field := desc.GetFieldBinding("Name2").Field
	fieldPtr := field.UnsafeGet(reflect2.PtrOf(val))
	//fieldVal := field.Get(reflect2.PtrOf(val))
	t.Logf("%+v %+v\n", field, field.Type())
	obj := *((*string)(fieldPtr))
	t.Logf("%+v\n", obj)

	id2 := desc.GetFieldBinding("ID2").Field
	id2Ptr := id2.UnsafeGet(reflect2.PtrOf(val))
	//fieldVal := field.Get(reflect2.PtrOf(val))
	t.Logf("%+v %+v\n", id2, id2.Type())
	id2Obj := *((*int)(id2Ptr))
	t.Logf("%+v\n", id2Obj)

	// Slice
	slice := []interface{}{val, val, val}
	typ := reflect2.TypeOf(slice)
	sliceType := typ.(*reflect2.UnsafeSliceType)
	t.Logf("%+v\n", sliceType.Elem())
	slicePtr := reflect2.PtrOf(slice)
	length := sliceType.UnsafeLengthOf(slicePtr)
	t.Logf("slice length: %+v\n", length)
	slicePtrObj := sliceType.UnsafeGetIndex(slicePtr, 0)
	t.Logf("%+v\n", slicePtrObj)

	sliceObj := sliceType.Elem().UnsafeIndirect(slicePtrObj)
	desc = DescribeStruct(&cfg, reflect2.TypeOf(sliceObj))
	t.Logf("%+v\n", desc.Type)
	for _, binding := range desc.Fields {
		t.Logf("%+v\n", binding)
	}

}

type stringCodec struct{}

func (codec *stringCodec) Encode(ptr unsafe.Pointer, buf *bytes.Buffer) {
	str := *((*string)(ptr))
	buf.WriteString(str)
}

type sliceEncoder struct {
	sliceType   *reflect2.UnsafeSliceType
	elemEncoder Encoder
}

//func encoderOfSlice( typ reflect2.Type) Encoder {
//	sliceType := typ.(*reflect2.UnsafeSliceType)
//	encoder := encoderOfType(sliceType.Elem())
//	return &sliceEncoder{sliceType, encoder}
//}
