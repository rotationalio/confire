package parse_test

import (
	"encoding"
	"testing"

	"github.com/rotationalio/confire/assert"
	. "github.com/rotationalio/confire/parse"
	"github.com/rotationalio/confire/structs"
)

type Does struct{}

type DoesNot struct{}

func (d *Does) Decode(string) error {
	return nil
}

func (d *Does) Set(string) error {
	return nil
}

func (d *Does) UnmarshalText([]byte) error {
	return nil
}

func (d *Does) UnmarshalBinary([]byte) error {
	return nil
}

type Spec struct {
	Does       *Does
	DoesNot    *DoesNot
	DoesVal    Does
	DoesNotVal DoesNot
}

func TestDecoderFrom(t *testing.T) {
	does := &Does{}
	doesnot := &DoesNot{}

	var val interface{}
	val = does

	_, ok := val.(Decoder)
	assert.Assert(t, ok, "expected that does implements Decoder")

	val = doesnot
	_, ok = val.(Decoder)
	assert.Assert(t, !ok, "expected that doesnot does not implements Decoder")

	spec := &Spec{Does: &Does{}, DoesNot: &DoesNot{}}
	s, err := structs.New(spec)
	assert.Ok(t, err)

	for _, field := range s.Fields() {
		decoder := DecoderFrom(field)
		decoderValue := DecoderFromValue(field.Reflect())

		if field.Name() == "Does" || field.Name() == "DoesVal" {
			assert.Assert(t, decoder != nil, "expected to extract a decoder")
			assert.Assert(t, decoderValue != nil, "expected to extract a decoder value")
		} else {
			assert.Assert(t, decoder == nil, "expected to not extract a decoder")
			assert.Assert(t, decoderValue == nil, "expected to not extract a decoder value")
		}
	}

}

func TestSetterFrom(t *testing.T) {
	does := &Does{}
	doesnot := &DoesNot{}

	var val interface{}
	val = does

	_, ok := val.(Setter)
	assert.Assert(t, ok, "expected that does implements Setter")

	val = doesnot
	_, ok = val.(Setter)
	assert.Assert(t, !ok, "expected that doesnot does not implements Setter")

	spec := &Spec{Does: &Does{}, DoesNot: &DoesNot{}}
	s, err := structs.New(spec)
	assert.Ok(t, err)

	for _, field := range s.Fields() {
		setter := SetterFrom(field)
		setterValue := SetterFromValue(field.Reflect())

		if field.Name() == "Does" || field.Name() == "DoesVal" {
			assert.Assert(t, setter != nil, "expected to extract a Setter")
			assert.Assert(t, setterValue != nil, "expected to extract a Setter value")
		} else {
			assert.Assert(t, setter == nil, "expected to not extract a Setter")
			assert.Assert(t, setterValue == nil, "expected to not extract a Setter value")
		}
	}
}

func TestTextUnmarshalerFrom(t *testing.T) {
	does := &Does{}
	doesnot := &DoesNot{}

	var val interface{}
	val = does

	_, ok := val.(encoding.TextUnmarshaler)
	assert.Assert(t, ok, "expected that does implements TextUnmarshaler")

	val = doesnot
	_, ok = val.(encoding.TextUnmarshaler)
	assert.Assert(t, !ok, "expected that doesnot does not implements TextUnmarshaler")

	spec := &Spec{Does: &Does{}, DoesNot: &DoesNot{}}
	s, err := structs.New(spec)
	assert.Ok(t, err)

	for _, field := range s.Fields() {
		text := TextUnmarshalerFrom(field)
		textValue := TextUnmarshalerFromValue(field.Reflect())

		if field.Name() == "Does" || field.Name() == "DoesVal" {
			assert.Assert(t, text != nil, "expected to extract a TextUnmarshaler")
			assert.Assert(t, textValue != nil, "expected to extract a TextUnmarshaler value")
		} else {
			assert.Assert(t, text == nil, "expected to not extract a TextUnmarshaler")
			assert.Assert(t, textValue == nil, "expected to not extract a TextUnmarshaler value")
		}
	}
}

func TestBinaryUnmarshalerFrom(t *testing.T) {
	does := &Does{}
	doesnot := &DoesNot{}

	var val interface{}
	val = does

	_, ok := val.(encoding.BinaryUnmarshaler)
	assert.Assert(t, ok, "expected that does implements BinaryUnmarshaler")

	val = doesnot
	_, ok = val.(encoding.BinaryUnmarshaler)
	assert.Assert(t, !ok, "expected that doesnot does not implements BinaryUnmarshaler")

	spec := &Spec{Does: &Does{}, DoesNot: &DoesNot{}}
	s, err := structs.New(spec)
	assert.Ok(t, err)

	for _, field := range s.Fields() {
		bin := BinaryUnmarshalerFrom(field)
		binValue := BinaryUnmarshalerFromValue(field.Reflect())

		if field.Name() == "Does" || field.Name() == "DoesVal" {
			assert.Assert(t, bin != nil, "expected to extract a BinaryUnmarshaler")
			assert.Assert(t, binValue != nil, "expected to extract a BinaryUnmarshaler value")
		} else {
			assert.Assert(t, bin == nil, "expected to not extract a BinaryUnmarshaler")
			assert.Assert(t, binValue == nil, "expected to not extract a BinaryUnmarshaler value")
		}
	}
}
