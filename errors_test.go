package confire_test

import (
	"testing"

	"go.rtnl.ai/confire"
	"go.rtnl.ai/confire/assert"
	"go.rtnl.ai/confire/errors"
)

func TestParserError(t *testing.T) {
	testCases := []struct {
		err error
		ok  bool
	}{
		{errors.ErrInvalidSpecification, false},
		{&errors.ParseError{Source: "a", Field: "b", Type: "c", Value: "d", Err: errors.ErrNotAStruct}, true},
		{errors.Required("", "foo"), false},
		{errors.Invalid("", "foo", "bar"), false},
		{errors.Parse("", "foo", errors.ErrNotExported), false},
		{errors.Wrap("", "foo", "bar", errors.ErrNotSettable, "qux"), false},
	}

	for _, tc := range testCases {
		target, ok := confire.ParseError(tc.err)
		if tc.ok {
			assert.True(t, ok)
			assert.Assert(t, target != nil, "expected target to be not nil")
		} else {
			assert.False(t, ok)
			assert.Assert(t, target == nil, "expected target to be nil")
		}

	}
}

func TestIsParserError(t *testing.T) {
	testCases := []struct {
		err    error
		assert assert.BoolAssertion
	}{
		{errors.ErrInvalidSpecification, assert.False},
		{errors.ErrNotAStruct, assert.False},
		{errors.ErrNotExported, assert.False},
		{errors.ErrNotSettable, assert.False},
		{&errors.ParseError{}, assert.True},
		{&errors.ParseError{Source: "a", Field: "b", Type: "c", Value: "d", Err: errors.ErrNotAStruct}, assert.True},
		{errors.Required("", "foo"), assert.False},
		{errors.Invalid("", "foo", "bar"), assert.False},
		{errors.Parse("", "foo", errors.ErrNotExported), assert.False},
		{errors.Wrap("", "foo", "bar", errors.ErrNotSettable, "qux"), assert.False},
	}

	for _, tc := range testCases {
		tc.assert(t, confire.IsParseError(tc.err))
	}
}
