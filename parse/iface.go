package parse

import (
	"encoding"
	"reflect"

	"github.com/rotationalio/confire/structs"
)

// Decoder has the same semantics as Setter, but takes higher precedence.
type Decoder interface {
	Decode(value string) error
}

// Setter is implemented by types can self-deserialize values.
// Any type that implements flag.Value also implements Setter.
type Setter interface {
	Set(value string) error
}

// Attempts to get a Decoder variable from the specified field.
func DecoderFrom(field *structs.Field) (d Decoder) {
	field.InterfaceFrom(func(v interface{}, ok *bool) { d, *ok = v.(Decoder) })
	return d
}

// Attempts to get a Setter variable from the specified field.
func SetterFrom(field *structs.Field) (s Setter) {
	field.InterfaceFrom(func(v interface{}, ok *bool) { s, *ok = v.(Setter) })
	return s
}

// Attempts to get an encoding.TextUnmarshaler variable from the specified field.
func TextUnmarshalerFrom(field *structs.Field) (t encoding.TextUnmarshaler) {
	field.InterfaceFrom(func(v interface{}, ok *bool) { t, *ok = v.(encoding.TextUnmarshaler) })
	return t
}

// Attempts to get an encoding.BinaryUnmarshaler variable from the specified field.
func BinaryUnmarshalerFrom(field *structs.Field) (b encoding.BinaryUnmarshaler) {
	field.InterfaceFrom(func(v interface{}, ok *bool) { b, *ok = v.(encoding.BinaryUnmarshaler) })
	return b
}

func interfaceFromValue(field reflect.Value, fn func(interface{}, *bool)) {
	// it may be impossible for a struct field to fail this check
	if !field.CanInterface() {
		return
	}
	var ok bool
	fn(field.Interface(), &ok)
	if !ok && field.CanAddr() {
		fn(field.Addr().Interface(), &ok)
	}
}

// Attempts to get a Decoder variable from the specified value.
func DecoderFromValue(field reflect.Value) (d Decoder) {
	interfaceFromValue(field, func(v interface{}, ok *bool) { d, *ok = v.(Decoder) })
	return d
}

// Attempts to get a Setter variable from the specified value.
func SetterFromValue(field reflect.Value) (s Setter) {
	interfaceFromValue(field, func(v interface{}, ok *bool) { s, *ok = v.(Setter) })
	return s
}

// Attempts to get an encoding.TextUnmarshaler variable from the specified value.
func TextUnmarshalerFromValue(field reflect.Value) (t encoding.TextUnmarshaler) {
	interfaceFromValue(field, func(v interface{}, ok *bool) { t, *ok = v.(encoding.TextUnmarshaler) })
	return t
}

// Attempts to get an encoding.BinaryUnmarshaler variable from the specified value.
func BinaryUnmarshalerFromValue(field reflect.Value) (b encoding.BinaryUnmarshaler) {
	interfaceFromValue(field, func(v interface{}, ok *bool) { b, *ok = v.(encoding.BinaryUnmarshaler) })
	return b
}
