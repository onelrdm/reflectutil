package xlsxutil

import (
	"github.com/onelrdm/conv"
	"unsafe"
)

type IntCodec struct{}

func (r IntCodec) Encode(ptr unsafe.Pointer, writer interface{}) {
	v := *((*int)(ptr))
	s := conv.MustString(v)
	_ = writer.(Writer).WriteString(s)
}
