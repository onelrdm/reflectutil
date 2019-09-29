package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/reflectutil"
	"github.com/tealeg/xlsx"
	"unsafe"
)

type sliceEncoder struct {
	sliceType   *reflect2.UnsafeSliceType
	elemEncoder reflectutil.Encoder
}

func (r *sliceEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	if r.sliceType.UnsafeIsNil(ptr) {
		return
	}
	length := r.sliceType.UnsafeLengthOf(ptr)
	if length == 0 {
		return
	}

	sheet := writer.(*xlsx.Sheet)
	for i := 0; i < length; i++ {
		row := sheet.AddRow()
		elemPtr := r.sliceType.UnsafeGetIndex(ptr, i)
		r.elemEncoder.Encode(elemPtr, row)
	}
}
