/*
Package structs helps with reflection and is based on the archived
github.com/fatih/structs package. Reflection is used a tremendous amount in the confire
package, so it makes sense to have helpers for reflection to simplify things.
*/
package structs

import (
	"fmt"
	"reflect"

	"github.com/rotationalio/confire/errors"
)

// Struct encapsulates a struct type to provide reflection around the struct.
type Struct struct {
	raw   interface{}
	value reflect.Value
}

// New returns a wrapped struct ready for reflection. It returns an error if the spec
// is not a struct or a pointer to a struct.
func New(spec interface{}) (*Struct, error) {
	val := reflect.ValueOf(spec)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, errors.ErrNotAStruct
	}

	return &Struct{
		raw:   spec,
		value: val,
	}, nil
}

// Returns true if the wrapped struct is a pointer.
func (s *Struct) IsPointer() bool {
	val := reflect.ValueOf(s.raw)
	return val.Kind() == reflect.Pointer
}

// Name returns the struct type name within its package.
func (s *Struct) Name() string {
	return s.value.Type().Name()
}

// Names returns a slice with all of the field names on the struct.
func (s *Struct) Names() []string {
	fields := s.Fields()
	names := make([]string, 0, len(fields))

	for _, field := range fields {
		names = append(names, field.Name())
	}

	return names
}

func (s *Struct) NumField() int {
	return s.value.Type().NumField()
}

// Fields returns a slice of all the fields on the struct.
func (s *Struct) Fields() []*Field {
	return getFields(s.value)
}

// Field returns a field by the specified name, returning an error if not found.
func (s *Struct) Field(name string) (*Field, error) {
	t := s.value.Type()
	field, ok := t.FieldByName(name)
	if !ok {
		return nil, fmt.Errorf("no field named %q on %s", name, s.Name())
	}

	return &Field{field: field, value: s.value.FieldByName(name)}, nil
}

// IsZero returns true if all fields on the struct are a zero-value for their type.
// Primitive types such as string, bool, int have zero-values "", false, 0, etc.
// Collection types are generally zero-valued if they are empty or nil. Nested structs
// are zero-valued if they are nil (for pointers) or if they are also zero-valued, e.g.
// all the fields on the struct are the zero-value for their type. If the field
// implements the IsZero() method (e.g. such as time.Time) then it is used.
func (s *Struct) IsZero() bool {
	fields := s.structFields()
	for _, field := range fields {
		val := s.value.FieldByName(field.Name)

		// Check if the value implements the Zero interface
		if zero := zeroFrom(val); zero != nil {
			if !zero.IsZero() {
				return false
			}
			continue
		}

		// If it is a nested struct, use struct IsZero
		if IsStruct(val.Interface()) {
			nested, _ := New(val.Interface())
			if !nested.IsZero() {
				return false
			}
			continue
		}

		// Get the current and zero-value for the given field and check if they're equal
		zero := reflect.Zero(val.Type()).Interface()
		current := val.Interface()

		if !reflect.DeepEqual(current, zero) {
			return false
		}
	}
	return true
}

// HasZero returns true if any field in a struct is zero-valued for their type.
// Primitive types such as string, bool, int have zero-values "", false, 0, etc.
// Collection types are generally zero-valued if they are empty or nil. Nested structs
// are zero-valued if they are nil (for pointers) or if they are also zero-valued, e.g.
// all the fields on the struct are the zero-value for their type. If the field
// implements the IsZero() method (e.g. such as time.Time) then it is used.
func (s *Struct) HasZero() bool {
	fields := s.structFields()
	for _, field := range fields {
		val := s.value.FieldByName(field.Name)

		// Check if the value implements the Zero interface
		if zero := zeroFrom(val); zero != nil {
			if zero.IsZero() {
				return true
			}
			continue
		}

		// If it is a nested struct, use struct IsZero
		if IsStruct(val.Interface()) {
			nested, _ := New(val.Interface())
			if nested.IsZero() {
				return true
			}
			continue
		}

		// Get the current and zero-value for the given field and check if they're equal
		zero := reflect.Zero(val.Type()).Interface()
		current := val.Interface()

		if reflect.DeepEqual(current, zero) {
			return true
		}
	}
	return false
}

func (s *Struct) structFields() (fields []reflect.StructField) {
	t := s.value.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Ignore any unexported fields. Note that this is the idiomatic way to check
		// for unexported fields rather than using CanSet or checking the casing of the
		// field name because PkgPath already performs field qualification.
		if field.PkgPath != "" {
			continue
		}

		fields = append(fields, field)
	}

	return fields
}

// IsStruct returns true if the given variable is a struct or a pointer to a struct.
func IsStruct(spec interface{}) bool {
	val := reflect.ValueOf(spec)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// beware uninitialized zero-valued structs
	if val.Kind() == reflect.Invalid {
		return false
	}

	return val.Kind() == reflect.Struct
}
