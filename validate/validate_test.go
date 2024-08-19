package validate_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rotationalio/confire/assert"
	confireErrors "github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/validate"
)

func TestValidator(t *testing.T) {
	type Specification struct {
		PropA Age
		PropB Age `ignored:"true"`
		PropC Age `validate:"ignored"`
		PropD Age `required:"true"`
		PropE Age `validate:"required"`
	}

	infos, err := validate.Gather(&Specification{})
	assert.Ok(t, err)
	validator := validate.ValidatorFrom(infos[0].Field)
	assert.Assert(t, validator != nil, "expected validator to implement the Validate interface")

	invalid := &Specification{PropA: Age(200)}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{PropA: Age(200), PropD: Age(300)}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	var target confireErrors.ValidationErrors
	ok := errors.As(err, &target)
	assert.True(t, ok)
	assert.Assert(t, len(target) == 3, "expected 3 validation errors got %s", target)

	valid := &Specification{PropD: Age(42), PropE: Age(63)}
	err = validate.Validate(valid)
	assert.Ok(t, err)
}

func TestRequired(t *testing.T) {
	type Specification struct {
		Skipped          string
		RequiredTag      string `required:"true"`
		NotRequiredTag   string `required:"false"`
		RequiredValidate string `validate:"required"`
		IgnoreRequired   string `required:"true" ignored:"true"`
		IgnoreValidate   string `required:"true" validate:"ignore"`
		Nested           Nested `required:"true"`
	}

	invalid := &Specification{}
	err := validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	var target confireErrors.ValidationErrors
	ok := errors.As(err, &target)
	assert.True(t, ok)
	assert.Assert(t, len(target) == 3, "expected 3 validation errors")

	partial := &Specification{RequiredTag: "foo", RequiredValidate: "bar"}
	err = validate.Validate(partial)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	var single *confireErrors.ValidationError
	ok = errors.As(err, &single)
	assert.True(t, ok)

	valid := &Specification{
		RequiredTag:      "foo",
		RequiredValidate: "bar",
		Nested:           Nested{PropB: 42},
	}
	err = validate.Validate(valid)
	assert.Ok(t, err)
}

func TestNestedRequired(t *testing.T) {
	type Specification struct {
		Nested NestedRequired `required:"true"`
	}

	invalid := &Specification{}
	err := validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{PropC: "baz"}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{PropA: "foo"}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{PropB: 32}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	valid := &Specification{Nested: NestedRequired{PropA: "foo", PropB: 32}}
	err = validate.Validate(valid)
	assert.Ok(t, err)
}

func TestNestedFieldsRequired(t *testing.T) {
	type Specification struct {
		Nested NestedRequired
	}

	invalid := &Specification{}
	err := validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{PropC: "baz"}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{PropA: "foo"}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: NestedRequired{PropB: 32}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	valid := &Specification{Nested: NestedRequired{PropA: "foo", PropB: 32}}
	err = validate.Validate(valid)
	assert.Ok(t, err)
}

func TestNestedPointerRequired(t *testing.T) {
	type Specification struct {
		Nested *NestedRequired
	}

	invalid := &Specification{}
	err := validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: &NestedRequired{}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: &NestedRequired{PropC: "baz"}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: &NestedRequired{PropA: "foo"}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	invalid = &Specification{Nested: &NestedRequired{PropB: 32}}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	valid := &Specification{Nested: &NestedRequired{PropA: "foo", PropB: 32}}
	err = validate.Validate(valid)
	assert.Ok(t, err)
}

func TestRequiredTypes(t *testing.T) {
	type Specification struct {
		Ptr       *string        `required:"true"`
		String    string         `required:"true"`
		Int       int            `required:"true"`
		Neg       int16          `required:"true"`
		Bool      bool           `required:"true"`
		Float     float64        `required:"true"`
		Uint      uint64         `required:"true"`
		Strings   []string       `required:"true"`
		Map       map[string]int `required:"true"`
		Dur       time.Duration  `required:"true"`
		TS        time.Time      `required:"true"`
		Nested    Nested         `required:"true"`
		NestedPtr *Nested        `required:"true"`
	}

	infos, err := validate.Gather(&Specification{})
	assert.Ok(t, err)

	invalid := &Specification{}
	err = validate.Validate(invalid)
	assert.Assert(t, err != nil, "expected a validation error to have occurred")

	var target confireErrors.ValidationErrors
	ok := errors.As(err, &target)
	assert.True(t, ok)
	assert.Assert(t, len(target) == len(infos), "expected %d validation errors", len(infos))

	foo := "foo"

	valid := &Specification{
		Ptr:       &foo,
		String:    foo,
		Int:       42,
		Neg:       -42,
		Bool:      true,
		Float:     0.5,
		Uint:      21,
		Strings:   []string{"a", "b", "c"},
		Map:       map[string]int{"a": 1, "b": 2},
		Dur:       5 * time.Minute,
		TS:        time.Now(),
		Nested:    Nested{PropA: foo},
		NestedPtr: &Nested{PropA: foo},
	}
	err = validate.Validate(valid)
	assert.Ok(t, err)

}

func TestUnknownValidator(t *testing.T) {
	type Specification struct {
		Whoopsie string `validate:"notthenameofanactualvalidatorbecausethisshouldnotbeone"`
	}

	err := validate.Validate(&Specification{})
	assert.Assert(t, err != nil, "expected an error to have occurred")
	assert.Equals(t, "unknown validator \"notthenameofanactualvalidatorbecausethisshouldnotbeone\"", err.Error())
}

type Nested struct {
	PropA string
	PropB int64
}

type NestedRequired struct {
	PropA string `required:"true"`
	PropB int64  `required:"true"`
	PropC string
}

type Age int

func (a Age) Validate() error {
	if a > 120 {
		return errors.New("humans don't live that long")
	}
	return nil
}

var _ validate.Validator = Age(1)
