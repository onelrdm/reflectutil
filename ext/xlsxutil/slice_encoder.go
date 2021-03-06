package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/reflectutil"
	"github.com/tealeg/xlsx"
	"unsafe"
)

type sliceEncoder struct {
	typ     *reflect2.UnsafeSliceType
	encoder reflectutil.Encoder
}

func (r *sliceEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {
	if r.typ.UnsafeIsNil(ptr) {
		return
	}
	length := r.typ.UnsafeLengthOf(ptr)
	if length == 0 {
		return
	}
	sheet := writer.(*xlsx.Sheet)

	// Column names
	encoder := reflectutil.EncoderOf(r.encoder).(*structEncoder)
	row := sheet.AddRow()
	for _, binding := range encoder.Fields {
		cell := row.AddCell()
		cell.Value = binding.Name
	}

	for i := 0; i < length; i++ {
		row := sheet.AddRow()
		elemPtr := r.typ.UnsafeGetIndex(ptr, i)
		r.encoder.Encode(elemPtr, row)
	}
}
