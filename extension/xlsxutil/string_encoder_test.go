package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"github.com/tealeg/xlsx"
	"testing"
)

func TestStringCodec_Encode(t *testing.T) {
	codec := StringCodec{}
	v := "TestStringCodec_Encode"
	w := CellWriter{Cell: &xlsx.Cell{}}
	codec.Encode(reflect2.PtrOf(v), &w)
	assert.Equal(t, v, w.String())
}
