package xlsxutil

import (
	"testing"
	"time"

	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"github.com/tealeg/xlsx"
)

func TestAnyCodec_Encode(t *testing.T) {
	{
		v := time.Now()
		codec := AnyCodec{typ: reflect2.TypeOf(v)}
		w := CellWriter{Cell: &xlsx.Cell{}}
		codec.Encode(reflect2.PtrOf(v), &w)
		assert.Equal(t, v.Format("2006-01-02 15:04:05"), w.String())
	}
	{
		v := 1
		codec := AnyCodec{typ: reflect2.TypeOf(v)}
		w := CellWriter{Cell: &xlsx.Cell{}}
		codec.Encode(reflect2.PtrOf(v), &w)
		assert.Equal(t, "1", w.String())
	}
}
