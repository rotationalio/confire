package confire

import (
	"errors"

	confireErrors "github.com/rotationalio/confire/errors"
)

// Extract a parse error from an error if it is one.
func ParseError(err error) (*confireErrors.ParseError, bool) {
	target := &confireErrors.ParseError{}
	if ok := errors.As(err, &target); ok {
		return target, true
	}
	return nil, false
}

// Returns true if the underlying error is a parser error.
func IsParseError(err error) bool {
	target := &confireErrors.ParseError{}
	return errors.As(err, &target)
}

// Extract a validation error from an error if it is one.
func ValidationError(err error) (*confireErrors.ValidationError, bool) {
	target := &confireErrors.ValidationError{}
	if ok := errors.As(err, &target); ok {
		return target, true
	}
	return nil, false
}

// Returns true if the underlying error is a validation error.
func IsValidationError(err error) bool {
	target := &confireErrors.ValidationError{}
	return errors.As(err, &target)
}
