package reflectutil

import (
	"unsafe"
)

type DereferenceEncoder struct {
	Encoder Encoder
}

func (r *DereferenceEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	r.Encoder.Encode(*((*unsafe.Pointer)(ptr)), writer)
}
