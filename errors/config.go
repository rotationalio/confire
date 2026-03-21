package errors

import (
	"errors"
	"fmt"
)

func Required(conf, field string) *InvalidConfig {
	return &InvalidConfig{
		conf:  conf,
		field: field,
		issue: "is required but not set",
		err:   ErrMissingRequired,
	}
}

func Invalid(conf, field, issue string, args ...any) *InvalidConfig {
	return &InvalidConfig{
		conf:  conf,
		field: field,
		issue: fmt.Sprintf(issue, args...),
	}
}

func Parse(conf, field string, err error) *InvalidConfig {
	return &InvalidConfig{
		conf:  conf,
		field: field,
		issue: fmt.Sprintf("could not parse value: %s", err.Error()),
		err:   err,
	}
}

func Wrap(conf, field, issue string, err error, args ...any) *InvalidConfig {
	return &InvalidConfig{
		conf:  conf,
		field: field,
		issue: fmt.Sprintf(issue, args...),
		err:   err,
	}
}

// Invalid is a field-specific configuration validation error and is returned either by
// field specification validation or by the user in a custom Validate() method.
type InvalidConfig struct {
	conf  string
	field string
	issue string
	err   error
}

func (e *InvalidConfig) Error() string {
	field := e.field
	if e.conf != "" {
		field = e.conf + "." + e.field
	}
	return fmt.Sprintf("invalid configuration: %s %s", field, e.issue)
}

func (e *InvalidConfig) Field() string {
	if e.conf != "" {
		return e.conf + "." + e.field
	}
	return e.field
}

func (e *InvalidConfig) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e *InvalidConfig) Unwrap() error {
	return e.err
}
