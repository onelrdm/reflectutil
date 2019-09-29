package xlsxutil

import (
	"bytes"
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/conv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntCodec_Encode(t *testing.T) {
	v := 123
	codec := IntCodec{}
	var buf bytes.Buffer
	codec.Encode(reflect2.PtrOf(v), &buf)
	assert.Equal(t, conv.MustString(v), buf.String())
}
