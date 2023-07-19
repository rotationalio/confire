package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSpecification = errors.New("specification must be a struct pointer")
	ErrNotAStruct           = errors.New("cannot wrap a non-struct type")
	ErrNotExported          = errors.New("field is not exported")
	ErrNotSettable          = errors.New("field is not settable")
)

type ParseError struct {
	Source string
	Field  string
	Type   string
	Value  string
	Err    error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("confire: could not parse %[2]s from %[1]s: converting %[4]q to type %[3]s: %[5]s", e.Source, e.Field, e.Type, e.Value, e.Err)
}

func (e *ParseError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Source string
	Err    error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid configuration: %s", e.Err)
}

func (e *ValidationError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}
