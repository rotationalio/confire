package structs_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rotationalio/confire/assert"
	"github.com/rotationalio/confire/structs"
)

func TestFieldMethods(t *testing.T) {
	spec := &WebColor{Color: Color{Red: 42, Green: 31, Blue: 0}, HexCode: "#42ac32"}
	s, err := structs.New(spec)
	assert.Ok(t, err)

	hexField, err := s.Field("HexCode")
	assert.Ok(t, err)

	// Assert that we can fetch a tag
	assert.Equals(t, "hex", hexField.Tag("json"))
	assert.Equals(t, "", hexField.Tag("yaml"))

	// Assert that we can fetch a value
	assert.Equals(t, "#42ac32", hexField.Value())

	// Check embedded fields
	embeddedField, err := s.Field("Color")
	assert.Ok(t, err)
	assert.Assert(t, embeddedField.IsEmbedded(), "expected color to be embedded")
	assert.Assert(t, !hexField.IsEmbedded(), "expected hexField to not be embedded")

	// Check exported fields
	assert.Assert(t, hexField.IsExported(), "expected hexField to be exported")

	unexported, err := s.Field("unexported")
	assert.Ok(t, err)
	assert.Assert(t, !unexported.IsExported(), "expected unexported field to not be exported")

	// Check is zero
	assert.Assert(t, !hexField.IsZero(), "expected hexfield not not be zero-valued")

	createdField, err := s.Field("Created")
	assert.Ok(t, err)
	assert.Assert(t, createdField.IsZero(), "expected created field to be zero-valued")

	// Check names and types
	assert.Equals(t, reflect.String, hexField.Kind())
	assert.Equals(t, reflect.Struct, createdField.Kind())

	// Test setting
	err = hexField.Set("hello world")
	assert.Ok(t, err)

	// Should not be able to set the wrong type
	err = createdField.Set("hello world")
	assert.Equals(t, "field is not settable: cannot set type \"string\" on field type \"struct\"", err.Error())

	err = unexported.Set("hello world")
	assert.Equals(t, "field is not exported", err.Error())

	// Test zeroing out the field
	err = hexField.Zero()
	assert.Ok(t, err)
	assert.Equals(t, "", spec.HexCode)
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
