package reflectutil

import (
	"bytes"
	"github.com/modern-go/reflect2"
	"unsafe"
)

type StructFieldEncoder struct {
	field        reflect2.StructField
	fieldEncoder Encoder
}

func (encoder *StructFieldEncoder) Encode(ptr unsafe.Pointer, buf *bytes.Buffer) {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	encoder.fieldEncoder.Encode(fieldPtr, buf)
}

func (encoder *StructFieldEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	isEmbeddedPtrNil, converted := encoder.fieldEncoder.(IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}

type IsEmbeddedPtrNil interface {
	IsEmbeddedPtrNil(ptr unsafe.Pointer) bool
}
