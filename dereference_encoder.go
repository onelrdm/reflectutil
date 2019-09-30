package reflectutil

import (
	"unsafe"
)

type DereferenceEncoder struct {
	Encoder Encoder
}

func (r *DereferenceEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	deReferenced := *((*unsafe.Pointer)(ptr))
	if deReferenced == nil {
		return
	}
	r.Encoder.Encode(deReferenced, writer)
}
