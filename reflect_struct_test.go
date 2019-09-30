package reflectutil

import (
	"github.com/modern-go/reflect2"
	"github.com/onelrdm/conv"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

type TestContext struct {
	cfg *Config
}

func NewContext(cfg *Config) *TestContext {
	return &TestContext{cfg: cfg}
}

func (r TestContext) Config() *Config {
	return r.cfg
}

func (r TestContext) NewEncoder(typ reflect2.Type) Encoder {
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

type TestWriter struct{}

func (TestWriter) WriteString(string) error{
	return nil
}

func (r AnyCodec) Encode(ptr unsafe.Pointer, writer interface{}) {
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
	_ = writer.(Writer).WriteString(s)
}

func TestDescribeStruct(t *testing.T) {
	now := time.Now()
	type Embed2 struct {
		ID2           int         `reflect:"id2"`
		Name2         string      `reflect:"name2"`
		Valid2        bool        `reflect:"valid2"`
		Time2         time.Time   `reflect:"time2"`
		Ignored       interface{} `reflect:"-"`
		unexported    interface{} `reflect:"unexported"`
		OtherTagField interface{} `other_tag_key:"unexported"`
		FieldName     interface{}
	}
	type Embed struct {
		*Embed2
		ID   int        `reflect:"id"`
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
		assert.NotNil(t, DescribeStruct(NewContext(&Config{}), reflect2.TypeOf(val)))
		assert.NotNil(t, DescribeStruct(NewContext(&Config{TaggedFieldOnly: true}), reflect2.TypeOf(val)))
		assert.NotNil(t, DescribeStruct(NewContext(&Config{TaggedFieldOnly: true, TagKey: "other_tag_key"}), reflect2.TypeOf(val)))
	}
	{
		type Embed struct {
			ID   int    `reflect:"id,3"`
			Name string `reflect:"name,1"`
			Data string `reflect:"data,2"`
		}
		type Val struct {
			Embed
		}
		val := Val{Embed: Embed{ID: 1, Name: "name", Data: "data"}}
		sd := DescribeStruct(NewContext(&Config{}), reflect2.TypeOf(val))
		assert.NotNil(t, sd)
		assert.Equal(t, sd.Fields[0].Name, "name")
		assert.Equal(t, sd.Fields[1].Name, "data")
		assert.Equal(t, sd.Fields[2].Name, "id")
	}
	{
		type Val struct {
			ID   int    `reflect:"id,3"`
			Name string `reflect:"name,1"`
		}
		val := Val{ID: 1, Name: "name"}
		DescribeStruct(NewContext(&Config{}), reflect2.TypeOf(val), func(field *reflect2.StructField) {
		})
	}
}

func TestDescribeStruct2(t *testing.T) {
	{
		type Embed struct {
			ID   int    `reflect:"id,3"`
			Name string `reflect:"name,1"`
			Data string `reflect:"data,2"`
		}
		type Val struct {
			*Embed
		}
		val := Val{Embed: &Embed{ID: 1, Name: "name", Data: "data"}}
		sd := DescribeStruct(NewContext(&Config{}), reflect2.TypeOf(val))
		assert.NotNil(t, sd)
		assert.Equal(t, sd.Fields[0].Name, "name")
		assert.Equal(t, sd.Fields[1].Name, "data")
		assert.Equal(t, sd.Fields[2].Name, "id")
		sd.Fields[0].Encoder.Encode(reflect2.PtrOf(val.Embed.Name), &TestWriter{})
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
