package structs

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotExported = errors.New("field is not exported")
	ErrNotSettable = errors.New("field is not settable")
)

func getFields(v reflect.Value) (fields []*Field) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		field := &Field{
			field: f,
			value: v.FieldByName(f.Name),
		}

		fields = append(fields, field)
	}

	return fields
}

// Field represents a single struct field and encapsulates high-level functionality for
// managing and working with the field using the reflect package.
type Field struct {
	value reflect.Value
	field reflect.StructField
}

// Tag returns the value assocaited with the key in the tag string. If there is no such
// key in the tag, an empty string is returned.
func (f *Field) Tag(key string) string {
	return f.field.Tag.Get(key)
}

// Value returns the underlying value of the field, panics if the field is not exported.
func (f *Field) Value() interface{} {
	return f.value.Interface()
}

// IsEmbedded returns true if the given field is an anonymous field.
func (f *Field) IsEmbedded() bool {
	return f.field.Anonymous
}

// IsExported returns true if the given field is exported.
func (f *Field) IsExported() bool {
	return f.field.PkgPath == ""
}

// IsZero returns true if the given field has a zero-value.
// Panics if the field is not exported
func (f *Field) IsZero() bool {
	if zero := zeroFrom(f.value); zero != nil {
		return zero.IsZero()
	}

	zero := reflect.Zero(f.value.Type()).Interface()
	current := f.Value()
	return reflect.DeepEqual(current, zero)
}

// Name returns the name of the given field
func (f *Field) Name() string {
	return f.field.Name
}

// Kind returns the fields kind, such as "string", "map", "bool", etc ..
func (f *Field) Kind() reflect.Kind {
	return f.value.Kind()
}

// Set sets the field to given value v. It returns an error if the field is not
// settable (not addressable or not exported) or if the given value's type
// doesn't match the fields type.
func (f *Field) Set(val interface{}) error {
	// Cannot set unexported fields (prevent panic)
	if !f.IsExported() {
		return ErrNotExported
	}

	// Ensure the value can be set
	if !f.value.CanSet() {
		return ErrNotSettable
	}

	given := reflect.ValueOf(val)
	if f.value.Kind() != given.Kind() {
		return fmt.Errorf("cannot set type %q on field type %q", given.Kind(), f.value.Kind())
	}

	f.value.Set(given)
	return nil
}

// Zero sets the field to its zero value.
func (f *Field) Zero() error {
	zero := reflect.Zero(f.value.Type()).Interface()
	return f.Set(zero)
}
