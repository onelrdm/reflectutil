// Package reflectutil implements reflect2 utility functions.
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

// ValEncoder is an internal type registered to cache as needed.
type ValEncoder interface {
	IsEmpty(ptr unsafe.Pointer) bool
	Encode(ptr unsafe.Pointer, buf *bytes.Buffer)
}

// FieldBinding describe how should we encode/decode the struct field
type FieldBinding struct {
	levels    []int
	Field     reflect2.StructField
	FromNames []string
	ToNames   []string
	Encoder   ValEncoder
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

// GetField get one field from the descriptor by its name.
// Can not use map here to keep field orders.
func (r *StructDescriptor) GetField(fieldName string) *FieldBinding {
	for _, binding := range r.Fields {
		if binding.Field.Name() == fieldName {
			return binding
		}
	}
	return nil
}

func createStructDescriptor(ctx *ctx, typ reflect2.Type, bindings []*FieldBinding, embeddedBindings []*FieldBinding) *StructDescriptor {
	// merge normal & embedded bindings & sort with original order
	allBindings := sortableFieldBinding(append(embeddedBindings, bindings...))
	sort.Sort(allBindings)
	return &StructDescriptor{Type: typ, Fields: allBindings}
}

func DescribeStruct(ctx *ctx, typ reflect2.Type) *StructDescriptor {
	var embeddedBindings []*FieldBinding
	var bindings []*FieldBinding
	concreteTyp := typ.(*reflect2.UnsafeStructType)
	for i := 0; i < concreteTyp.NumField(); i++ {
		field := concreteTyp.Field(i)
		tag, found := field.Tag().Lookup(ctx.getTagKey())
		if ctx.OnlyTaggedField && !found && !field.Anonymous() {
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
				structDescriptor := DescribeStruct(ctx, typ)
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

		fieldNames := convertFieldNames(field.Name(), tagParts[0], tag)
		encoder := ctx.encoder
		binding := &FieldBinding{
			levels:    []int{i},
			Field:     field,
			FromNames: fieldNames,
			ToNames:   fieldNames,
			Encoder:   encoder,
		}
		bindings = append(bindings, binding)
	}
	return createStructDescriptor(ctx, typ, bindings, embeddedBindings)
}

func convertFieldNames(fieldName string, tagName string, tag string) []string {
	// ignore
	if tag == "-" {
		return []string{}
	}
	// private
	unexported := unicode.IsLower(rune(fieldName[0]))
	if unexported {
		return []string{}
	}
	// rename
	if tagName != "" {
		return []string{tagName}
	}
	return []string{fieldName}
}
