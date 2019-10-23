package xlsxutil

import (
	"unsafe"
)

type StringCodec struct{}

func (r *StringCodec) Encode(ptr unsafe.Pointer, writer interface{}) {
	s := *((*string)(ptr))
	_ = writer.(Writer).WriteString(s)
}
