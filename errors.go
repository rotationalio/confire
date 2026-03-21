package confire

import (
	"errors"

	confireErrors "go.rtnl.ai/confire/errors"
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

// Extract a configuration error from an error if it is one.
func InvalidConfig(err error) (*confireErrors.InvalidConfig, bool) {
	target := &confireErrors.InvalidConfig{}
	if ok := errors.As(err, &target); ok {
		return target, true
	}
	return nil, false
}

// Returns true if the underlying error is a configuration error.
func IsInvalidConfig(err error) bool {
	target := &confireErrors.InvalidConfig{}
	return errors.As(err, &target)
}

// Bringing in the errors from the errors package for convenience.
var (
	Required = confireErrors.Required
	Invalid  = confireErrors.Invalid
	Parse    = confireErrors.Parse
	Wrap     = confireErrors.Wrap
	Join     = confireErrors.Join
)
