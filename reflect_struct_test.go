package reflectutil

import (
	"reflect"
	"strconv"
	"testing"
	"time"
	"unsafe"

	"github.com/modern-go/reflect2"
	"github.com/onelrdm/conv"
	"github.com/stretchr/testify/assert"
)

func newEncoder(typ reflect2.Type, _ reflect2.StructField) Encoder {
	kind := typ.Kind()
	switch kind {
	default:
		return &AnyCodec{valType: typ}
	}
}

type AnyCodec struct {
	valType reflect2.Type
}

type Writer interface {
	WriteString(string) error
}

type TestWriter struct {
	b []byte
}

func (r *TestWriter) WriteString(s string) error {
	r.b = append(r.b, s...)
	return nil
}

func (r AnyCodec) Encode(ptr unsafe.Pointer, w interface{}) {
	v := r.valType.UnsafeIndirect(ptr)
	var s string
	switch v := v.(type) {
	case time.Time:
		s = v.Format("2006-01-02 15:04:05")
	case *time.Time:
		if v != nil {
			s = v.Format("2006-01-02 15:04:05")
		}
	default:
		s = conv.MustString(v)
	}
	_ = w.(Writer).WriteString(s)
}

func TestDescribeStruct(t *testing.T) {
	now := time.Now()
	type Embed2 struct {
		ID2           int         `reflect:"id2"`
		Name2         string      `reflect:"name2"`
		Valid2        bool        `reflect:"valid2"`
		Time2         time.Time   `reflect:"time2"`
		Ignored       interface{} `reflect:"-"`
		unexported    interface{} ``
		OtherTagField interface{} `other_tag_key:"unexported"`
		FieldName     interface{}
	}
	type Embed struct {
		ID int `reflect:"id"`
		*Embed2
		Name string     `reflect:"name"`
		Time *time.Time `reflect:"time"`
	}
	type Val struct {
		Embed
		Data *time.Time `reflect:"data"`
	}
	{
		val := Val{
			Embed: Embed{
				Embed2: &Embed2{
					ID2:    2,
					Name2:  "name2",
					Valid2: true,
					Time2:  now,
				},
				ID:   1,
				Name: "name",
				Time: &now,
			},
			Data: &now,
		}
		typ := reflect2.TypeOf(val).(*reflect2.UnsafeStructType)
		assert.NotNil(t, DescribeStruct(typ, newEncoder, &Option{}))

		sd := DescribeStruct(typ, newEncoder, &Option{TaggedFieldOnly: true})
		assert.NotNil(t, sd)
		assert.Equal(t, "id", sd.Fields[0].Name)
		assert.Equal(t, "id2", sd.Fields[1].Name)
		assert.Equal(t, "name2", sd.Fields[2].Name)
		assert.Equal(t, "valid2", sd.Fields[3].Name)
		assert.Equal(t, "time2", sd.Fields[4].Name)
		assert.Equal(t, "name", sd.Fields[5].Name)
		assert.Equal(t, "time", sd.Fields[6].Name)
		assert.Equal(t, "data", sd.Fields[7].Name)

		assert.NotNil(t, DescribeStruct(typ, newEncoder, &Option{TaggedFieldOnly: true, TagKey: "other_tag_key"}))
	}
	{
		type Embed struct {
			ID   int    `reflect:"id,4"`
			Name string `reflect:"name,2"`
			Data string `reflect:"data,3"`
		}
		type Val struct {
			Embed
			Count int `reflect:"count,1"`
		}
		val := Val{Embed: Embed{ID: 1, Name: "name", Data: "data"}}
		typ := reflect2.TypeOf(val).(*reflect2.UnsafeStructType)
		sd := DescribeStruct(typ, newEncoder, &Option{})
		assert.NotNil(t, sd)
		assert.Equal(t, sd.Fields[0].Name, "count")
		assert.Equal(t, sd.Fields[1].Name, "name")
		assert.Equal(t, sd.Fields[2].Name, "data")
		assert.Equal(t, sd.Fields[3].Name, "id")

		var w TestWriter
		for _, binding := range sd.Fields {
			binding.Encoder.Encode(reflect2.PtrOf(val), &w)
		}
		assert.Equal(t, strconv.Itoa(val.Count)+val.Embed.Name+val.Embed.Data+strconv.Itoa(val.Embed.ID), string(w.b))
	}
}

func BenchmarkDescribeStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func TestStructDescriptor_GetFieldBinding(t *testing.T) {
	r := &StructDescriptor{
		Type: nil,
		Fields: []*FieldBinding{
			{
				levels: nil,
				Field: &reflect2.UnsafeStructField{
					StructField: reflect.StructField{
						Name:      "Test",
						PkgPath:   "",
						Type:      nil,
						Tag:       "",
						Offset:    0,
						Index:     nil,
						Anonymous: false,
					},
				},
				Name:    "Test",
				Encoder: nil,
			},
		},
	}
	assert.NotNil(t, r.GetFieldBinding("Test"))
	assert.Nil(t, r.GetFieldBinding(""))
}
