package errors

import (
	"errors"
	"testing"

	"go.rtnl.ai/confire/assert"
)

func TestRequired(t *testing.T) {
	err := Required("", "bind_addr")
	assert.Equals(t, "invalid configuration: bind_addr is required but not set", err.Error())
	assert.Assert(t, err.Is(ErrMissingRequired), "required should wrap a missing required field error")
	assert.Equals(t, ErrMissingRequired, err.Unwrap())
}

func TestInvalid(t *testing.T) {
	err := Invalid("", "mode", "invalid mode %q", "foo")
	assert.Equals(t, "invalid configuration: mode invalid mode \"foo\"", err.Error())
	assert.Equals(t, nil, err.Unwrap())
}

func TestParse(t *testing.T) {
	err := Parse("", "bind_addr", errors.New("invalid bind address"))
	assert.Equals(t, "invalid configuration: bind_addr could not parse value: invalid bind address", err.Error())
	assert.Equals(t, errors.New("invalid bind address"), err.Unwrap())
}

func TestWrap(t *testing.T) {
	err := Wrap("", "bind_addr", "invalid bind address %q", errors.New("invalid bind address"), "foo")
	assert.Equals(t, "invalid configuration: bind_addr invalid bind address \"foo\"", err.Error())
	assert.Equals(t, errors.New("invalid bind address"), err.Unwrap())
}

func TestInvalidConfig(t *testing.T) {
	testCases := []struct {
		conf     string
		field    string
		issue    string
		err      error
		errstr   string
		fieldstr string
	}{
		{
			conf:     "",
			field:    "bind_addr",
			issue:    "is required but not set",
			errstr:   "invalid configuration: bind_addr is required but not set",
			fieldstr: "bind_addr",
		},
		{
			conf:     "",
			field:    "bind_addr",
			issue:    "is required but not set",
			err:      errors.New("required field is zero valued"),
			errstr:   "invalid configuration: bind_addr is required but not set",
			fieldstr: "bind_addr",
		},
		{
			conf:     "telemetry",
			field:    "service_name",
			issue:    "cannot have spaces or start with a number",
			errstr:   "invalid configuration: telemetry.service_name cannot have spaces or start with a number",
			fieldstr: "telemetry.service_name",
		},
		{
			conf:     "telemetry",
			field:    "service_name",
			issue:    "cannot have spaces or start with a number",
			err:      errors.New("invalid service name"),
			errstr:   "invalid configuration: telemetry.service_name cannot have spaces or start with a number",
			fieldstr: "telemetry.service_name",
		},
	}

	for _, tc := range testCases {
		err := Wrap(tc.conf, tc.field, tc.issue, tc.err)
		assert.Equals(t, tc.errstr, err.Error())
		assert.Equals(t, tc.fieldstr, err.Field())
		if tc.err != nil {
			assert.Assert(t, err.Is(tc.err), "wrap should wrap the error")
			assert.Equals(t, tc.err, err.Unwrap())
		} else {
			assert.Equals(t, nil, err.Unwrap())
		}
	}
}
