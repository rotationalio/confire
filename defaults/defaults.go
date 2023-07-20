/*
Package default allows users to initialize structs with default values as defined in
struct tags. The default values are parsed from the struct tag string in the same way
that environment variables from the confire env package are parsed.
*/
package defaults

import (
	"reflect"

	"github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/parse"
	"github.com/rotationalio/confire/structs"
)

const tagDefault = "default"

// Process a struct to add the defaults from the default struct tag, parsing the struct
// tag string into the correct type using the same methodology as processing environment
// variables. Most types are parsed using the strconv package (time.Duration is handled
// specially). If the type implements Decoder, Setter, TextUnmarshaler, or
// BinaryUnmarshaler, then those decoders are used to parse the default value.
func Process(spec interface{}) (err error) {
	var s *structs.Struct
	if s, err = structs.New(spec); err != nil {
		return errors.ErrInvalidSpecification
	}

	if !s.IsPointer() {
		return errors.ErrInvalidSpecification
	}

	for _, field := range s.Fields() {
		// Skip any fields that cannot be set.
		if !field.CanSet() {
			continue
		}

		// Handle pointers if necessary
		for field.Kind() == reflect.Ptr {
			if field.IsNil() {
				if field.TypeKind() != reflect.Struct {
					// nil pointer to a non-struct: leave it alone.
					break
				}

				// nil pointer to a struct: create a zero-instance
				if err = field.Init(); err != nil {
					panic(err)
				}
			}
			field = field.Elem()
		}

		// Check if this field has a default value
		if value := field.Tag(tagDefault); value != "" {
			if err = parse.ParseField(value, field); err != nil {
				return err
			}
		} else if field.Kind() == reflect.Struct {
			if err = Process(field.Pointer()); err != nil {
				return err
			}
		}

	}

	return nil
}

// MustProcess is the same as Process but panics if an error occurs
func MustProcess(spec interface{}) {
	if err := Process(spec); err != nil {
		panic(err)
	}
}

// SetDefaults is an alias of Process
var SetDefaults func(spec interface{}) error = Process
