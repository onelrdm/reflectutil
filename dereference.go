package reflectutil

import (
	"bytes"
	"unsafe"
)

type dereferenceEncoder struct {
	ValueEncoder ValEncoder
}

func (encoder *dereferenceEncoder) Encode(ptr unsafe.Pointer, buf *bytes.Buffer) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		buf.Write([]byte{'n', 'u', 'l', 'l'})
	} else {
		encoder.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), buf)
	}
}

func (encoder *dereferenceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	dePtr := *((*unsafe.Pointer)(ptr))
	if dePtr == nil {
		return true
	}
	return encoder.ValueEncoder.IsEmpty(dePtr)
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
