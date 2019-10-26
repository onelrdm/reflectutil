package xlsxutil

import (
	"io"
	"reflect"

	"github.com/modern-go/reflect2"
	"github.com/tealeg/xlsx"

	"github.com/onelrdm/reflectutil"
)

var DefaultOption = &reflectutil.Option{TaggedFieldOnly: true}

func Convert(obj interface{}, writer io.Writer) error {
	if obj == nil {
		return ErrIsNil
	}
	typ := reflectutil.TypeOf(obj)
	if typ.Kind() != reflect.Struct {
		return ErrMustBeStruct
	}

	file := xlsx.NewFile()
	sd := reflectutil.DescribeStruct(typ.(*reflect2.UnsafeStructType), GetEncoder, DefaultOption)
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
		binding.Encoder.Encode(reflect2.PtrOf(obj), sheet)
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
