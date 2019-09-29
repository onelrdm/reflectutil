package reflectutil

import (
	"github.com/modern-go/reflect2"
	"reflect"
	"testing"
	"unsafe"
)

func TestStructFieldEncoder_Encoder(t *testing.T) {
	(&StructFieldEncoder{}).Encoder()
}

type MockFieldEncoder struct{}

func (MockFieldEncoder) Encode(ptr unsafe.Pointer, writer interface{}) {}

func TestStructFieldEncoder_Encode(t *testing.T) {
	encoder := &StructFieldEncoder{
		field: &reflect2.UnsafeStructField{
			StructField: reflect.StructField{},
		},
		fieldEncoder: &MockFieldEncoder{},
	}
	encoder.Encode(reflect2.PtrOf(1), nil)
}
