package confire_test

import (
	"testing"

	"github.com/rotationalio/confire"
	"github.com/rotationalio/confire/assert"
	"github.com/rotationalio/confire/errors"
)

func TestParserError(t *testing.T) {
	testCases := []struct {
		err error
		ok  bool
	}{
		{errors.ErrInvalidSpecification, false},
		{&errors.ParseError{Source: "a", Field: "b", Type: "c", Value: "d", Err: errors.ErrNotAStruct}, true},
		{&errors.ValidationError{Source: "foo", Err: errors.ErrNotExported}, false},
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
		{&errors.ValidationError{}, assert.False},
		{&errors.ValidationError{Source: "foo", Err: errors.ErrNotExported}, assert.False},
	}

	for _, tc := range testCases {
		tc.assert(t, confire.IsParseError(tc.err))
	}
}

func TestValidationError(t *testing.T) {
	testCases := []struct {
		err error
		ok  bool
	}{
		{errors.ErrInvalidSpecification, false},
		{&errors.ParseError{Source: "a", Field: "b", Type: "c", Value: "d", Err: errors.ErrNotAStruct}, false},
		{&errors.ValidationError{Source: "foo", Err: errors.ErrNotExported}, true},
	}

	for _, tc := range testCases {
		target, ok := confire.ValidationError(tc.err)
		if tc.ok {
			assert.True(t, ok)
			assert.Assert(t, target != nil, "expected target to be not nil")
		} else {
			assert.False(t, ok)
			assert.Assert(t, target == nil, "expected target to be nil")
		}

	}
}

func TestIsValidationError(t *testing.T) {
	testCases := []struct {
		err    error
		assert assert.BoolAssertion
	}{
		{errors.ErrInvalidSpecification, assert.False},
		{errors.ErrNotAStruct, assert.False},
		{errors.ErrNotExported, assert.False},
		{errors.ErrNotSettable, assert.False},
		{&errors.ParseError{}, assert.False},
		{&errors.ParseError{Source: "a", Field: "b", Type: "c", Value: "d", Err: errors.ErrNotAStruct}, assert.False},
		{&errors.ValidationError{}, assert.True},
		{&errors.ValidationError{Source: "foo", Err: errors.ErrNotExported}, assert.True},
	}

	for _, tc := range testCases {
		tc.assert(t, confire.IsValidationError(tc.err))
	}
}
