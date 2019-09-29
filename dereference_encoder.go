package reflectutil

import (
	"unsafe"
)

type DereferenceEncoder struct {
	ValueEncoder Encoder
}

func (r *DereferenceEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	r.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), writer)
}
