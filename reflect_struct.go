package reflectutil

import (
	"bytes"
	"github.com/modern-go/reflect2"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unsafe"
)

// Encoder is an internal type registered to cache as needed.
type Encoder interface {
	Encode(ptr unsafe.Pointer, buf *bytes.Buffer)
}

// FieldBinding describe how should we encode/decode the struct field
type FieldBinding struct {
	levels  []int
	Field   reflect2.StructField
	Name    string
	Encoder Encoder
}

type sortableFieldBinding []*FieldBinding

func (r sortableFieldBinding) Len() int {
	return len(r)
}

func (r sortableFieldBinding) Less(i, j int) bool {
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
	Type   reflect2.Type
	Fields []*FieldBinding
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

func DescribeStruct(cfg *Config, typ reflect2.Type) *StructDescriptor {
	var embeddedBindings []*FieldBinding
	var bindings []*FieldBinding
	concreteTyp := typ.(*reflect2.UnsafeStructType)
	for i := 0; i < concreteTyp.NumField(); i++ {
		field := concreteTyp.Field(i)
		tag, found := field.Tag().Lookup(cfg.getTagKey())
		if cfg.OnlyTaggedField && !found && !field.Anonymous() {
			continue
		}
		if tag == "-" {
			continue
		}
		tagParts := strings.Split(tag, ",")
		if field.Anonymous() && (tag == "" || tagParts[0] == "") {
			typ := field.Type()
			kind := typ.Kind()
			isPtr := kind == reflect.Ptr
			if isPtr {
				typ = typ.(*reflect2.UnsafePtrType).Elem()
				kind = typ.Kind()
			}
			if kind == reflect.Struct {
				structDescriptor := DescribeStruct(cfg, typ)
				if isPtr {
					for _, binding := range structDescriptor.Fields {
						binding.levels = append([]int{i}, binding.levels...)
						binding.Encoder = &StructFieldEncoder{field, &dereferenceEncoder{binding.Encoder}}
						embeddedBindings = append(embeddedBindings, binding)
					}
				} else {
					for _, binding := range structDescriptor.Fields {
						binding.levels = append([]int{i}, binding.levels...)
						binding.Encoder = &StructFieldEncoder{field, binding.Encoder}
						embeddedBindings = append(embeddedBindings, binding)
					}
				}
				continue
			}
		}
		fieldNames := convertFieldName(field.Name(), tagParts[0], tag)
		binding := &FieldBinding{
			levels:  []int{i},
			Field:   field,
			Name:    fieldNames,
			Encoder: nil, // TODO
		}
		bindings = append(bindings, binding)
	}
	// merge normal & embedded bindings & sort with original order
	allBindings := sortableFieldBinding(append(embeddedBindings, bindings...))
	sort.Sort(allBindings)
	return &StructDescriptor{Type: typ, Fields: allBindings}
}

func convertFieldName(fieldName string, tagName string, tag string) string {
	if tag == "-" {
		return ""
	}
	unexported := unicode.IsLower(rune(fieldName[0]))
	if unexported {
		return ""
	}
	if tagName != "" {
		return tagName
	}
	return fieldName
}
