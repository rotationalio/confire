package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidSpecification = errors.New("specification must be a struct pointer")
	ErrNotAStruct           = errors.New("cannot wrap a non-struct type")
	ErrNotExported          = errors.New("field is not exported")
	ErrNotSettable          = errors.New("field is not settable")
	ErrMissingRequired      = errors.New("required field is zero valued")
)

type ValidationErrors []*InvalidConfig

func (e ValidationErrors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d validation errors occurred:", len(e)))
	for _, err := range e {
		sb.WriteString(fmt.Sprintf("\n    - %s", err.Error()))
	}
	return sb.String()
}

func (e ValidationErrors) Is(target error) bool {
	for _, err := range e {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e ValidationErrors) Contains(target error) bool {
	return e.Is(target)
}

func Join(err error, errs ...error) error {
	var (
		verrs  ValidationErrors
		isverr bool
	)

	// If the first error is not nil and it is not ValidationErrors then use the
	// regular errors.Join function. If it is a ValidationError then continue. If the
	// first error is nil, then create new ValidationErrors.
	if err != nil {
		// If the error is an InvalidConfig error then create a new ValidationErrors
		if cerr, ok := err.(*InvalidConfig); ok {
			verrs = ValidationErrors{cerr}
		} else {
			// If the error is ValidationErrors then flatten it into the current ValidationErrors.
			verrs, isverr = err.(ValidationErrors)
			if !isverr {
				errs = append([]error{err}, errs...)
				return errors.Join(errs...)
			}
		}
	} else {
		verrs = make(ValidationErrors, 0, len(errs))
	}

	// Loop through the remaining errors and append them to the ValidationErrors if
	// they are not nil and they are InvalidConfig errors. If they are not InvalidConfig
	// errors then use the regular errors.Join function.
	for _, err := range errs {
		if err != nil {
			// If the error is ValidationErrors then flatten it into the current ValidationErrors.
			if errs, ok := err.(ValidationErrors); ok {
				verrs = append(verrs, errs...)
				continue
			}

			cerr, iscerr := err.(*InvalidConfig)
			if !iscerr {
				errs = append([]error{err}, errs...)
				return errors.Join(errs...)
			}
			verrs = append(verrs, cerr)
		}
	}

	// If the ValidationErrors is empty, then return nil.
	switch len(verrs) {
	case 0:
		return nil
	case 1:
		return verrs[0]
	default:
		return verrs
	}
}
