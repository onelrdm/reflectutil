package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/conv"
	"time"
	"unsafe"
)

type AnyCodec struct {
	valType reflect2.Type
}

type Writer interface {
	WriteString(string) error
}

func (r AnyCodec) Encode(ptr unsafe.Pointer, writer interface{}) {
	v := r.valType.UnsafeIndirect(ptr)
	var s string
	switch v := v.(type) {
	case time.Time:
		s = v.Format("2006-01-02 15:04:05")
	case *time.Time:
		if v != nil {
			s = v.Format("2006-01-02 15:04:05")
		}
	default:
		s = conv.MustString(v)
	}
	_ = writer.(Writer).WriteString(s)
}
