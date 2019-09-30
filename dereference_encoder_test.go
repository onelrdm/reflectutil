package reflectutil

import (
	"github.com/modern-go/reflect2"
	"testing"
)

func Test_dereferenceEncoder_Encode(t *testing.T) {
	encoder := &DereferenceEncoder{Encoder: &MockFieldEncoder{}}
	encoder.Encode(reflect2.PtrOf(1), nil)
	type MockEmbeddedStruct struct {}
	type MockStruct struct {
		MockEmbeddedStruct
	}
	encoder.Encode(reflect2.PtrOf(MockStruct{}), nil)
}
