package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/reflectutil"
	"github.com/tealeg/xlsx"
	"io"
	"reflect"
)

func Convert(obj interface{}, writer io.Writer) error {
	if obj == nil {
		return ErrIsNil
	}
	ctx := NewContext(&reflectutil.Config{TaggedFieldOnly: true})
	typ := reflectutil.TypeOf(obj)
	if typ.Kind() != reflect.Struct {
		return ErrMustBeStruct
	}

	file := xlsx.NewFile()
	sd := reflectutil.DescribeStruct(ctx, typ)
	var hasSlice bool
	for _, binding := range sd.Fields {
		field := binding.Field
		typ := field.Type()
		kind := typ.Kind()
		if kind != reflect.Slice {
			continue
		}
		sheet, err := file.AddSheet(binding.Name)
		if err != nil {
			return err
		}
		ptr := field.UnsafeGet(reflect2.PtrOf(obj))
		binding.Encoder.Encode(ptr, sheet)
		hasSlice = true
	}
	if !hasSlice {
		return ErrMustHaveSlice
	}
	if err := file.Write(writer); err != nil {
		return err
	}
	return nil
}
