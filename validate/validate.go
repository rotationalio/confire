package validate

import (
	goerrors "errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/structs"
)

const (
	tagRequired  = "required"
	tagValidator = "validate"
	tagIgnored   = "ignored"
)

type Validator interface {
	Validate() error
}

// Validate runs the given struct through the validation workflow as follows: if the
// required tag is set to true and the field is zero-valued then an error is returned.
// Otherwise, if the field is a Validator its validate method is called. Finally if a
// built-in validator is specified by the validate tag, then it is applied to the field
// value. The invalid configuration is returned as a multi-error. If the validate tag
// is set to ignored or the field has no required/validate tag and is not a Validator,
// then no validation is applied to the field.
func Validate(spec interface{}) (err error) {
	var infos []Info
	if infos, err = Gather(spec); err != nil {
		return err
	}

	errs := make(errors.ValidationErrors, 0, len(infos))
	for _, info := range infos {
		if verr := info.Validate.Validate(); verr != nil {
			errs = append(errs, asValidationError(verr, info.Field.Name()))
		}
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errs
	}
}

func Gather(spec interface{}) (infos []Info, err error) {
	var s *structs.Struct
	if s, err = structs.New(spec); err != nil {
		return nil, errors.ErrInvalidSpecification
	}

	infos = make([]Info, 0, s.NumField())
	for _, field := range s.Fields() {
		// If the ignored tag is set or the validator is set to ignored, skip the field
		if isTrue(field.Tag(tagIgnored)) || ignoreValidation(field.Tag(tagValidator)) {
			continue
		}

		// If this is a struct then gather validators for the nested fields
		if field.Kind() == reflect.Struct {
			var subinfos []Info
			if subinfos, err = Gather(field.Pointer()); err != nil {
				return nil, err
			}
			infos = append(infos, subinfos...)
		}

		// Chain validators together if necessary
		validators := make([]Validator, 0, 3)

		// Check if the required tag is set and add required validator if it is
		if isTrue(field.Tag(tagRequired)) {
			validators = append(validators, Required(field))
		}

		// Check if the field is a validator
		if validator := ValidatorFrom(field); validator != nil {
			validators = append(validators, validator)
		}

		// Check if there is a validator tag
		if validator := field.Tag(tagValidator); validator != "" {
			// Lookup the specified validator in the validation library.
			switch validator {
			case "required":
				validators = append(validators, Required(field))
			default:
				return nil, fmt.Errorf("unknown validator %q", validator)
			}
		}

		// If no validators were specified by the user, ignore this field
		if len(validators) == 0 {
			continue
		}

		info := Info{Field: field}
		if len(validators) == 1 {
			info.Validate = validators[0]
		} else {
			info.Validate = Many(validators...)
		}

		infos = append(infos, info)
	}

	return infos, nil
}

type Info struct {
	Field    *structs.Field // The actual field to get the validation info from (along with tags)
	Validate Validator      // The validator to apply to the field
}

// Attempts to get a Validator variable from the specified field.
func ValidatorFrom(field *structs.Field) (v Validator) {
	field.InterfaceFrom(func(i interface{}, ok *bool) { v, *ok = i.(Validator) })
	return v
}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func ignoreValidation(s string) bool {
	s = strings.ToLower(s)
	if s == "ignore" || s == "ignored" {
		return true
	}
	return false
}

func asValidationError(err error, source string) *errors.ValidationError {
	target := &errors.ValidationError{}
	if goerrors.As(err, &target) {
		return target
	}
	return &errors.ValidationError{
		Source: source,
		Err:    err,
	}
}
