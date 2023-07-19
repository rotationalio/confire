package structs

import (
	"fmt"
	"reflect"

	"github.com/rotationalio/confire/errors"
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

// Reflect returns the underlying reflect value for the field
func (f *Field) Reflect() reflect.Value {
	return f.value
}

// Pointer returns a pointer interface representing the address of f. It panics if
// CanAddr() returns false. This is typically used to used to obtain a pointer to a
// struct field or slice element in order to call a method that requires a pointer
// receiver.
func (f *Field) Pointer() interface{} {
	return f.value.Addr().Interface()
}

// IsEmbedded returns true if the given field is an anonymous field.
func (f *Field) IsEmbedded() bool {
	return f.field.Anonymous
}

// IsExported returns true if the given field is exported.
func (f *Field) IsExported() bool {
	return f.field.PkgPath == ""
}

// CanSet returns true if the given field is settable.
func (f *Field) CanSet() bool {
	return f.value.CanSet()
}

// CanAddr reports whether the value's pointer can be obtained with Pointer(). Such
// values are called addressable. A value is addressable if it is an element of a slice,
// an element of an addressable array, a field of an addressable struct, or the result
// of dereferencing a pointer. If CanAddr returns false, calling Addr will panic.
func (f *Field) CanAddr() bool {
	return f.value.CanAddr()
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

// IsNil reports if the field value is nil. The value must be a chan, func, interface,
// map, pointer, or slice value, otherwise this method panics.
func (f *Field) IsNil() bool {
	return f.value.IsNil()
}

// Name returns the name of the given field
func (f *Field) Name() string {
	return f.field.Name
}

// Kind returns the fields kind, such as "string", "map", "bool", etc ..
func (f *Field) Kind() reflect.Kind {
	return f.value.Kind()
}

// Kind returns the type of the field's kind, e.g. Array, Chan, Map, Pointer, or Slice.
func (f *Field) TypeKind() reflect.Kind {
	return f.value.Type().Elem().Kind()
}

// Elem returns a Field with the value that the interface f contains or that the pointer
// f points to. It panics if f's Kind is not Interface or Pointer. It returns the zero
// Value if f is nil.
func (f *Field) Elem() *Field {
	return &Field{
		value: f.value.Elem(),
		field: f.field,
	}
}

// Type returns the field value's type.
func (f *Field) Type() reflect.Type {
	return f.value.Type()
}

// Set sets the field to given value v. It returns an error if the field is not
// settable (not addressable or not exported) or if the given value's type
// doesn't match the fields type.
func (f *Field) Set(val interface{}) error {
	// Cannot set unexported fields (prevent panic)
	if !f.IsExported() {
		return errors.ErrNotExported
	}

	// Ensure the value can be set
	if !f.value.CanSet() {
		return errors.ErrNotSettable
	}

	given := reflect.ValueOf(val)
	if f.value.Kind() != given.Kind() {
		return fmt.Errorf("%w: cannot set type %q on field type %q", errors.ErrNotSettable, given.Kind(), f.value.Kind())
	}

	f.value.Set(given)
	return nil
}

// Zero sets the field to its zero value.
func (f *Field) Zero() error {
	zero := reflect.Zero(f.value.Type()).Interface()
	return f.Set(zero)
}

// InterfaceFrom is a complex type assertion that allows you to pass in a function that
// performs a type assertion to the underlying field value. If the underlying value can
// interface and implements the specified assertion, then the value and assertion bool
// are set on the arguments, allowing you to "extract" an interface type from the value.
//
// That was a bit complicated, so here is an example. Say you wanted to use the field as
// a encoding.TextUnmarshaler (e.g. you want to call the UnmarshalText() method of the
// value in the field). You would pass in a function as follows:
//
// field.InterfaceFrom(func(v interface{}, ok *bool) { t, *ok := v.(encoding.TextUnmarshaler )})
//
// The variable t is now a TextUnmarshaler, and you can call t.UnmarshalText() so long
// as the ok is true or t is not nil.
func (f *Field) InterfaceFrom(fn func(interface{}, *bool)) {
	if !f.value.CanInterface() {
		return
	}

	var ok bool
	fn(f.value.Interface(), &ok)
	if !ok && f.value.CanAddr() {
		fn(f.value.Addr().Interface(), &ok)
	}
}

// Fields returns a slice of Fields, usually used to get the fields from a nested struct.
func (f *Field) Fields() []*Field {
	return getFields(f.value)
}
