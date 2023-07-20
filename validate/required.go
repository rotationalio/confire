package validate

import (
	"github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/structs"
)

// Required returns a validation error if the specified field is zero valued.
func Required(field *structs.Field) Validator {
	return &required{field: field}
}

type required struct {
	field *structs.Field
}

func (r required) Validate() error {
	if r.field.IsZero() {
		return &errors.ValidationError{
			Source: r.field.Name(),
			Err:    errors.ErrMissingRequiredField,
		}
	}
	return nil
}

// Many returns a validator that is comprised of many sub-validators, that are each
// applied in turn. If one of the validators fails, that error is immediately returned.
func Many(validators ...Validator) Validator {
	return &many{validators: validators}
}

type many struct {
	validators []Validator
}

func (m many) Validate() (err error) {
	for _, validator := range m.validators {
		if err = validator.Validate(); err != nil {
			return err
		}
	}
	return nil
}
