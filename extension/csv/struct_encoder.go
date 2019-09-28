package csv

import (
	"bytes"
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/reflectutil"
	"unsafe"
)

type StructEncoder struct {
	typ    reflect2.Type
	fields []reflectutil.StructField
}

func (encoder *StructEncoder) Encode(ptr unsafe.Pointer, buf *bytes.Buffer) {
	buf.WriteObjectStart()
	for _, field := range encoder.fields {
		if field.Encoder.IsEmbeddedPtrNil(ptr) {
			continue
		}
		buf.WriteObjectField(field.Name)
		buf.WriteMore()
		field.Encoder.Encode(ptr, buf)
	}
	buf.WriteObjectEnd()
}