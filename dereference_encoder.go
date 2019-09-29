package reflectutil

import (
	"unsafe"
)

type dereferenceEncoder struct {
	ValueEncoder Encoder
}

func (encoder *dereferenceEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	encoder.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), writer)
}

func (encoder *dereferenceEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	deReferenced := *((*unsafe.Pointer)(ptr))
	if deReferenced == nil {
		return true
	}
	isEmbeddedPtrNil, converted := encoder.ValueEncoder.(IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := unsafe.Pointer(deReferenced)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}
