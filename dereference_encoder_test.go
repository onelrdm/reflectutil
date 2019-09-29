package reflectutil

import (
	"github.com/modern-go/reflect2"
	"testing"
)

func Test_dereferenceEncoder_Encode(t *testing.T) {
	encoder := &DereferenceEncoder{ValueEncoder: &MockFieldEncoder{}}
	encoder.Encode(reflect2.PtrOf(1), nil)
}