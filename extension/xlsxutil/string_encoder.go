package xlsxutil

import (
	"fmt"
	"unsafe"
)

type StringCodec struct{}

func (r *StringCodec) Encode(ptr unsafe.Pointer, writer interface{}) {
	fmt.Printf("%+v\n", ptr)
	s := *((*string)(ptr))
	_ = writer.(Writer).WriteString(s)
}
