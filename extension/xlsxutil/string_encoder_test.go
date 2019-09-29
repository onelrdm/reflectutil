package xlsxutil

import (
	"bytes"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringCodec_Encode(t *testing.T) {
	var buf bytes.Buffer
	codec := StringCodec{}
	v := "TestStringCodec_Encode"
	codec.Encode(reflect2.PtrOf(v), &buf)
	assert.Equal(t, v, buf.String())
}
