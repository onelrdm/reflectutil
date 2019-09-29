package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/reflectutil"
	"github.com/tealeg/xlsx"
	"io"
	"reflect"
)

func Convert(obj interface{}, writer io.Writer) error {
	ctx := NewContext(&reflectutil.Config{TaggedFieldOnly: true})
	typ := reflect2.TypeOf(obj)
	kind := typ.Kind()

	if kind == reflect.Ptr {
		ptrType := typ.(*reflect2.UnsafePtrType)
		typ = ptrType.Elem()
		kind = typ.Kind()
	}

	if kind != reflect.Struct {
		return ErrMustBeStruct
	}

	file := xlsx.NewFile()
	sd := reflectutil.DescribeStruct(ctx, typ)
	var hasSlice bool
	for _, binding := range sd.Fields {
		kind := binding.Field.Type().Kind()
		if kind != reflect.Slice {
			continue
		}
		hasSlice = true
		field := binding.Field
		ptr := field.UnsafeGet(reflect2.PtrOf(obj))
		sheet, err := file.AddSheet(binding.Name)
		if err != nil {
			return err
		}
		binding.Encoder.Encode(ptr, sheet)
	}
	if !hasSlice {
		return ErrMustHaveSlice
	}
	if err := file.Write(writer); err != nil {
		return err
	}
	return nil
}
