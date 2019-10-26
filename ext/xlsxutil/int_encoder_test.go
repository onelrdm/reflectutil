package xlsxutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/conv"
	"github.com/stretchr/testify/assert"
	"github.com/tealeg/xlsx"
	"testing"
)

func TestIntCodec_Encode(t *testing.T) {
	v := 123
	codec := IntCodec{}
	w := CellWriter{Cell: &xlsx.Cell{}}
	codec.Encode(reflect2.PtrOf(v), &w)
	assert.Equal(t, conv.MustString(v), w.String())
}
