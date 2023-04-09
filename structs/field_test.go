package structs_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rotationalio/confire/structs"
)

func TestFieldMethods(t *testing.T) {
	spec := &WebColor{Color: Color{Red: 42, Green: 31, Blue: 0}, HexCode: "#42ac32"}
	s, err := structs.New(spec)
	ok(t, err)

	hexField, err := s.Field("HexCode")
	ok(t, err)

	// Assert that we can fetch a tag
	equals(t, "hex", hexField.Tag("json"))
	equals(t, "", hexField.Tag("yaml"))

	// Assert that we can fetch a value
	equals(t, "#42ac32", hexField.Value())

	// Check embedded fields
	embeddedField, err := s.Field("Color")
	ok(t, err)
	assert(t, embeddedField.IsEmbedded(), "expected color to be embedded")
	assert(t, !hexField.IsEmbedded(), "expected hexField to not be embedded")

	// Check exported fields
	assert(t, hexField.IsExported(), "expected hexField to be exported")

	unexported, err := s.Field("unexported")
	ok(t, err)
	assert(t, !unexported.IsExported(), "expected unexported field to not be exported")

	// Check is zero
	assert(t, !hexField.IsZero(), "expected hexfield not not be zero-valued")

	createdField, err := s.Field("Created")
	ok(t, err)
	assert(t, createdField.IsZero(), "expected created field to be zero-valued")

	// Check names and types
	equals(t, reflect.String, hexField.Kind())
	equals(t, reflect.Struct, createdField.Kind())

	// Test setting
	err = hexField.Set("hello world")
	ok(t, err)

	// Should not be able to set the wrong type
	err = createdField.Set("hello world")
	equals(t, "cannot set type \"string\" on field type \"struct\"", err.Error())

	err = unexported.Set("hello world")
	equals(t, "field is not exported", err.Error())

	// Test zeroing out the field
	err = hexField.Zero()
	ok(t, err)
	equals(t, "", spec.HexCode)
}

type WebColor struct {
	Color
	HexCode    string `json:"hex"`
	Created    time.Time
	unexported string
}

func (w *WebColor) IVal() string {
	// appease linter
	return w.unexported
}
