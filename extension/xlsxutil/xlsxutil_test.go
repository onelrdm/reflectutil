package xlsxutil

import (
	"bytes"
	"errors"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
	"time"
)

type ErrorWriter struct{}

func (ErrorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("")
}

func TestConvert(t *testing.T) {
	type Embed2 struct {
		ID2    int       `reflect:"id2"`
		Name2  string    `reflect:"name2"`
		Valid2 bool      `reflect:"valid2"`
		Time2  time.Time `reflect:"time2"`
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
	type Slice struct {
		List         []Val   `reflect:"list"`
		ExtendedID   int     `reflect:"extended_id"`
		ExtendedName string  `reflect:"extended_name"`
		ExtendedData float64 `reflect:"extended_name"`
	}

	tests := []func(){
		func() {
			now := time.Now()
			val := Val{
				Embed: Embed{
					Embed2: &Embed2{
						ID2:    2,
						Name2:  "name2",
						Valid2: true,
						Time2:  now,
					},
					ID:   123,
					Name: "name",
					Time: &now,
				},
				Data: &now,
			}
			slice := Slice{
				List:         []Val{val, val, val},
				ExtendedID:   999999,
				ExtendedName: "ExtendedName",
				ExtendedData: 555,
			}
			var buf bytes.Buffer
			err := Convert(slice, &buf)
			assert.NoError(t, err)

			slice = Slice{List: nil}
			err = Convert(slice, &buf)
			assert.NoError(t, err)

			slice = Slice{List: []Val{}}
			err = Convert(slice, &buf)
			assert.NoError(t, err)

			err = Convert(slice, &ErrorWriter{})
			assert.Error(t, err)

			err = Convert(1, &buf)
			assert.Error(t, err)

			emptyStruct := struct{ ID int `reflect:"id"` }{}
			err = Convert(emptyStruct, &buf)
			assert.Error(t, err)

			type DuplicatedSlice struct {
				List  []Val `reflect:"list"`
				List2 []Val `reflect:"list"`
			}
			duplicatedSlice := DuplicatedSlice{
				List:  []Val{val, val, val},
				List2: []Val{val, val, val},
			}
			err = Convert(duplicatedSlice, nil)
			assert.Error(t, err)
		},
		func() {
			type SliceListPtr struct {
				List []*Val `reflect:"list"`
			}
			now := time.Now()
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
			slice := &SliceListPtr{
				List: []*Val{&val},
			}
			var buf bytes.Buffer
			err := Convert(slice, &buf)
			assert.NoError(t, err)
		},
		func() {
			type Slice struct {
				List []*Val `reflect:"list"`
			}
			type SliceInterface struct {
				Data interface{}
			}
			now := time.Now()
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
			slice := &Slice{
				List: []*Val{&val},
			}
			sliceInterface := SliceInterface{
				Data: slice,
			}
			//var buf bytes.Buffer
			//assert.NoError(t, Convert(sliceInterface.Data, &buf))
			//assert.Error(t, Convert(nil, &buf))

			file, err := os.OpenFile("test1.xlsx", os.O_RDWR|os.O_CREATE, 0755)
			assert.NoError(t, err)
			assert.NoError(t, Convert(sliceInterface.Data, file))
		},
	}
	tests[2]()
	t.Skip()
	for _, test := range tests {
		test()
	}
}

func TestConvert2(t *testing.T) {
	type Embed2 struct {
		ID2    int       `reflect:"id2"`
		Name2  string    `reflect:"name2"`
		Valid2 bool      `reflect:"valid2"`
		Time2  time.Time `reflect:"time2"`
	}
	type Embed struct {
		Embed2
		ID   int        `reflect:"id"`
		Name string     `reflect:"name"`
		Time *time.Time `reflect:"time"`
	}
	type Val struct {
		*Embed
		Data *time.Time `reflect:"data"`
	}

	tests := []func(){
		func() {
			type Slice struct {
				List []*Val `reflect:"list"`
			}
			type SliceInterface struct {
				Data interface{}
			}
			now := time.Now()
			val := Val{
				Embed: &Embed{
					Embed2: Embed2{
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
			slice := &Slice{
				List: []*Val{&val, &val},
			}
			//sliceInterface := SliceInterface{
			//	Data: slice,
			//}
			//var buf bytes.Buffer
			//assert.NoError(t, Convert(sliceInterface.Data, &buf))
			//assert.Error(t, Convert(nil, &buf))

			file, err := os.OpenFile("test1.xlsx", os.O_RDWR|os.O_CREATE, 0755)
			assert.NoError(t, err)
			assert.NoError(t, Convert(slice, file))
		},
	}
	for _, test := range tests {
		test()
	}
}

func TestInterfacePassing(t *testing.T) {
	type Val struct {
		Name string
	}
	type Data struct {
		ID  interface{}
		Val interface{}
	}
	var obj interface{}
	obj = Data{
		Val: Val{},
	}
	fn := func(field *reflect2.StructField) {
		typ := (*field).Type()
		kind := typ.Kind()
		if kind == reflect.Interface {
			ptr := (*field).UnsafeGet(reflect2.PtrOf(obj))
			v := typ.UnsafeIndirect(ptr)
			typ := reflect2.TypeOf(v)
			(*field).(*reflect2.UnsafeStructField).StructField.Type = typ.Type1()
		}
	}
	reflect.ValueOf(obj).Type()
	typ := reflect2.TypeOf(obj)
	concreteTyp := typ.(*reflect2.UnsafeStructType)
	field := concreteTyp.Field(1)
	t.Logf("field.Type %+v", field.Type())
	fn(&field)
	t.Logf("field %+v", field)
}
