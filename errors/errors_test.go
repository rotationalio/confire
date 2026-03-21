package errors_test

import (
	"errors"
	"testing"

	"go.rtnl.ai/confire/assert"
	. "go.rtnl.ai/confire/errors"
)

func TestValidationErrors(t *testing.T) {
	errs := make(ValidationErrors, 0, 3)
	assert.Equals(t, "0 validation errors occurred:", errs.Error())

	errs = append(errs, Required("", "foo"))
	assert.Equals(t, "invalid configuration: foo is required but not set", errs.Error())
	assert.True(t, errs.Contains(ErrMissingRequired))
	assert.False(t, errs.Contains(ErrNotAStruct))

	errs = append(errs, Wrap("", "colors", "at least one color must be specified", errors.New("at least one color should be specified")))
	assert.Equals(t, "2 validation errors occurred:\n    - invalid configuration: foo is required but not set\n    - invalid configuration: colors at least one color must be specified", errs.Error())

	invalid := Invalid("", "port", "port number out of range")
	errs = append(errs, invalid)
	assert.Equals(t, "3 validation errors occurred:\n    - invalid configuration: foo is required but not set\n    - invalid configuration: colors at least one color must be specified\n    - invalid configuration: port port number out of range", errs.Error())
	assert.True(t, errs.Contains(invalid))
	assert.False(t, errs.Contains(ErrNotAStruct))
}

func TestJoin(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Equals(t, nil, Join(nil))
		assert.Equals(t, nil, Join(nil, nil))
		assert.Equals(t, nil, Join(nil, nil, nil))
		assert.Equals(t, nil, Join(nil, nil, nil, nil))
		assert.Equals(t, nil, Join(ValidationErrors{}, nil))
		assert.Equals(t, nil, Join(ValidationErrors{}, nil, nil))
		assert.Equals(t, nil, Join(ValidationErrors{}, nil, nil, nil))
		assert.Equals(t, nil, Join(nil, ValidationErrors{}))
		assert.Equals(t, nil, Join(nil, ValidationErrors{}, nil))
		assert.Equals(t, nil, Join(nil, ValidationErrors{}, nil, ValidationErrors{}))
		assert.Equals(t, nil, Join(ValidationErrors{}, nil, ValidationErrors{}, nil))
		assert.Equals(t, nil, Join(ValidationErrors{}, ValidationErrors{}, ValidationErrors{}, ValidationErrors{}))
	})

	t.Run("Single", func(t *testing.T) {
		required := Required("", "foo")
		assert.Equals(t, required, Join(required))
		assert.Equals(t, required, Join(required, nil))
		assert.Equals(t, required, Join(nil, required))
		assert.Equals(t, required, Join(ValidationErrors{}, required))
		assert.Equals(t, required, Join(required, nil, nil))
		assert.Equals(t, required, Join(nil, required, nil))
		assert.Equals(t, required, Join(ValidationErrors{}, nil, required, nil, nil))
		assert.Equals(t, required, Join(ValidationErrors{}, nil, nil, nil, required))
	})

	t.Run("Multiple", func(t *testing.T) {
		required := Required("", "foo")
		badport := Invalid("", "port", "port number out of range")
		parse := Parse("", "bind_addr", errors.New("invalid bind address"))
		colors := Wrap("", "colors", "at least one color must be specified", errors.New("at least one color should be specified"))

		testCases := []struct {
			err              error
			errs             []error
			expectedLength   int
			expectedContains []error
		}{
			{required, []error{badport}, 2, []error{required, badport}},
			{required, []error{badport, parse}, 3, []error{required, badport, parse}},
			{required, []error{badport, parse, colors}, 4, []error{required, badport, parse, colors}},
			{nil, []error{required, badport, parse, colors}, 4, []error{required, badport, parse, colors}},
			{ValidationErrors{}, []error{required, badport, parse, colors}, 4, []error{required, badport, parse, colors}},
			{nil, []error{nil, required, nil, badport, nil, nil}, 2, []error{required, badport}},
			{ValidationErrors{}, []error{nil, nil, badport, nil, colors}, 2, []error{badport, colors}},
			{nil, []error{ValidationErrors{required, badport, parse, colors}}, 4, []error{required, badport, parse, colors}},
			{ValidationErrors{required, badport}, []error{parse, colors}, 4, []error{required, badport, parse, colors}},
			{ValidationErrors{required, badport}, []error{ValidationErrors{parse, colors}}, 4, []error{required, badport, parse, colors}},
			{ValidationErrors{required}, []error{ValidationErrors{badport}, ValidationErrors{parse}, ValidationErrors{colors}}, 4, []error{required, badport, parse, colors}},
		}

		for _, tc := range testCases {
			err := Join(tc.err, tc.errs...)
			verrs, isverr := err.(ValidationErrors)
			assert.True(t, isverr)
			assert.Equals(t, tc.expectedLength, len(verrs))
			for _, err := range tc.expectedContains {
				assert.True(t, verrs.Contains(err))
			}
		}
	})

	t.Run("UglyDuck", func(t *testing.T) {
		other := errors.New("not a validation error")
		required := Required("", "foo")
		badport := Invalid("", "port", "port number out of range")

		testCases := []struct {
			err        error
			errs       []error
			expectedIs []error
		}{
			{other, []error{required, badport}, []error{other, required, badport}},
			{nil, []error{other, required, badport}, []error{other, required, badport}},
			{ValidationErrors{}, []error{required, other, badport}, []error{other, required, badport}},
			{nil, []error{required, badport, other}, []error{other, required, badport}},
		}

		for _, tc := range testCases {
			errs := Join(tc.err, tc.errs...)
			_, isverr := errs.(ValidationErrors)
			assert.False(t, isverr)
			for _, err := range tc.expectedIs {
				assert.True(t, errors.Is(errs, err))
			}
		}
	})
}
