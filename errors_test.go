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
}
