package reflectutil

import (
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/modern-go/reflect2"
	"github.com/onelrdm/conv"
)

type GetEncoder func(reflect2.Type, reflect2.StructField) Encoder

// FieldBinding describe how should we encode/decode the struct field
type FieldBinding struct {
	order   int
	levels  []int
	Field   reflect2.StructField
	Name    string
	Encoder Encoder
}

func ToBindingName(fieldName string, tagName string) string {
	if unexported := unicode.IsLower(rune(fieldName[0])); unexported {
		return ""
	}
	if tagName != "" {
		return tagName
	}
	return fieldName
}

type sortableFieldBinding []*FieldBinding

func (r sortableFieldBinding) Len() int {
	return len(r)
}

func (r sortableFieldBinding) Less(i, j int) bool {
	if r[i].order < r[j].order {
		return true
	} else if r[i].order > r[j].order {
		return false
	}

	left := r[i].levels
	right := r[j].levels
	k := 0
	for {
		if left[k] < right[k] {
			return true
		} else if left[k] > right[k] {
			return false
		}
		k++
	}
}

func (r sortableFieldBinding) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// StructDescriptor describe how should we encode/decode the struct
type StructDescriptor struct {
	Type   *reflect2.UnsafeStructType
	Fields []*FieldBinding
}

func DescribeStruct(typ *reflect2.UnsafeStructType, getEncoder GetEncoder, opt *Option) *StructDescriptor {
	var embeddedBindings []*FieldBinding
	var bindings []*FieldBinding
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag, found := field.Tag().Lookup(opt.Tag())
		if opt.TaggedFieldOnly && !found && !field.Anonymous() {
			continue
		}
		if tag == "-" {
			continue
		}
		tagParts := strings.Split(tag, ",")
		if field.Anonymous() && (tag == "" || tagParts[0] == "") {
			typ := field.Type()
			kind := typ.Kind()
			isEmbeddedPtr := kind == reflect.Ptr
			if isEmbeddedPtr {
				typ = typ.(*reflect2.UnsafePtrType).Elem()
				kind = typ.Kind()
			}
			if kind == reflect.Struct {
				sd := DescribeStruct(typ.(*reflect2.UnsafeStructType), getEncoder, opt)
				if isEmbeddedPtr {
					for _, binding := range sd.Fields {
						binding.levels = append([]int{i}, binding.levels...)
						binding.Encoder = &StructFieldEncoder{field: field, fieldEncoder: &DereferenceEncoder{binding.Encoder}}
						embeddedBindings = append(embeddedBindings, binding)
					}
				} else {
					for _, binding := range sd.Fields {
						binding.levels = append([]int{i}, binding.levels...)
						binding.Encoder = &StructFieldEncoder{field, binding.Encoder}
						embeddedBindings = append(embeddedBindings, binding)
					}
				}
				continue
			}
		}
		order := 0
		if len(tagParts) > 1 {
			order = int(conv.MustInt64(tagParts[1]))
		}
		binding := &FieldBinding{
			order:   order,
			levels:  []int{i},
			Field:   field,
			Name:    ToBindingName(field.Name(), tagParts[0]),
			Encoder: &StructFieldEncoder{field: field, fieldEncoder: getEncoder(field.Type(), field)},
		}
		bindings = append(bindings, binding)
	}
	// merge normal & embedded bindings & sort with original order or order in tag
	allBindings := sortableFieldBinding(append(embeddedBindings, bindings...))
	sort.Sort(allBindings)
	return &StructDescriptor{Type: typ, Fields: allBindings}
}

// GetFieldBinding get one field from the descriptor by its name.
// Can not use map here to keep field orders.
func (r *StructDescriptor) GetFieldBinding(fieldName string) *FieldBinding {
	for _, binding := range r.Fields {
		if binding.Field.Name() == fieldName {
			return binding
		}
	}
	return nil
}
