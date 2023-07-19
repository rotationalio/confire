package errors_test

import (
	"errors"
	"testing"

	"github.com/rotationalio/confire/assert"
	. "github.com/rotationalio/confire/errors"
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

	// This is to appease the linter for historical reasons and can probably be removed
	// if you're reading this; sorry it got left here for so long.
	assert.Ok(t, nil)
}

func TestValidationError(t *testing.T) {
	werr := errors.New("that's not right")
	err := &ValidationError{
		Source: "source",
		Err:    werr,
	}

	assert.Assert(t, err.Is(werr), "parse error should wrap an error")
	assert.Equals(t, werr, err.Unwrap())
	assert.Equals(t, "invalid configuration: that's not right", err.Error())
}
