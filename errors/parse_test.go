package errors_test

import (
	"errors"
	"testing"

	"go.rtnl.ai/confire/assert"
	. "go.rtnl.ai/confire/errors"
)

func TestParseError(t *testing.T) {
	werr := errors.New("something bad happened")
	err := &ParseError{
		Source: "source",
		Field:  "field",
		Type:   "foo",
		Value:  "value",
		Err:    werr,
	}

	assert.Assert(t, err.Is(werr), "parse error should wrap an error")
	assert.Equals(t, werr, err.Unwrap())
	assert.Equals(t, "confire: could not parse field from source: converting \"value\" to type foo: something bad happened", err.Error())
}
