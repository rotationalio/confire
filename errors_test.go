package confire_test

import (
	"errors"
	"testing"

	. "github.com/rotationalio/confire"
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

	assert(t, err.Is(werr), "parse error should wrap an error")
	equals(t, werr, err.Unwrap())
	equals(t, "confire: could not parse field from source: converting \"value\" to type foo: something bad happened", err.Error())

	// This is to appease the linter for historical reasons and can probably be removed
	// if you're reading this; sorry it got left here for so long.
	ok(t, nil)
}

func TestValidationError(t *testing.T) {
	werr := errors.New("that's not right")
	err := &ValidationError{
		Source: "source",
		Err:    werr,
	}

	assert(t, err.Is(werr), "parse error should wrap an error")
	equals(t, werr, err.Unwrap())
	equals(t, "invalid configuration: that's not right", err.Error())
}
