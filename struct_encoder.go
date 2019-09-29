package reflectutil

import (
	"github.com/modern-go/reflect2"
	"unsafe"
)

type StructFieldEncoder struct {
	field        reflect2.StructField
	fieldEncoder Encoder
}

func (r *StructFieldEncoder) Encoder() Encoder {
	return r.fieldEncoder
}

func (r *StructFieldEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	fieldPtr := r.field.UnsafeGet(ptr)
	r.fieldEncoder.Encode(fieldPtr, writer)
}
