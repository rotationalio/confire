/*
Package parse implements a reflection-based mechanism of converting a string (e.g. from
an environment variable or from a struct field) into a Go type for the struct.
*/
package parse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/structs"
)

func Parse(value string, field reflect.Value) error {
	// Attempt to use the decoder, setter, and unmarshalers to parse the field.
	if decoder := DecoderFromValue(field); decoder != nil {
		if err := decoder.Decode(value); err != nil {
			return &errors.ParseError{
				Source: "Decoder",
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if setter := SetterFromValue(field); setter != nil {
		if err := setter.Set(value); err != nil {
			return &errors.ParseError{
				Source: "Setter",
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if txt := TextUnmarshalerFromValue(field); txt != nil {
		if err := txt.UnmarshalText([]byte(value)); err != nil {
			return &errors.ParseError{
				Source: "TextUnmarshaler",
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if bin := BinaryUnmarshalerFromValue(field); bin != nil {
		// TODO: should we decode base64 or hex data here?
		if err := bin.UnmarshalBinary([]byte(value)); err != nil {
			return &errors.ParseError{
				Source: "BinaryUnmarshaler",
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if err := parse(value, field); err != nil {
		return &errors.ParseError{
			Type:  field.Type().Name(),
			Value: value,
			Err:   err,
		}
	}
	return nil
}

// ParseField parses the given type from the field and sets it.
func ParseField(value string, field *structs.Field) error {
	// Attempt to use the decoder, setter, and unmarshalers to parse the field.
	if decoder := DecoderFrom(field); decoder != nil {
		if err := decoder.Decode(value); err != nil {
			return &errors.ParseError{
				Source: "Decoder",
				Field:  field.Name(),
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if setter := SetterFrom(field); setter != nil {
		if err := setter.Set(value); err != nil {
			return &errors.ParseError{
				Source: "Setter",
				Field:  field.Name(),
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if txt := TextUnmarshalerFrom(field); txt != nil {
		if err := txt.UnmarshalText([]byte(value)); err != nil {
			return &errors.ParseError{
				Source: "TextUnmarshaler",
				Field:  field.Name(),
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if bin := BinaryUnmarshalerFrom(field); bin != nil {
		// TODO: should we decode base64 or hex data here?
		if err := bin.UnmarshalBinary([]byte(value)); err != nil {
			return &errors.ParseError{
				Source: "BinaryUnmarshaler",
				Field:  field.Name(),
				Type:   field.Type().Name(),
				Value:  value,
				Err:    err,
			}
		}
		return nil
	}

	if err := parse(value, field.Reflect()); err != nil {
		return &errors.ParseError{
			Field: field.Name(),
			Type:  field.Type().Name(),
			Value: value,
			Err:   err,
		}
	}
	return nil
}

func parse(value string, field reflect.Value) (err error) {
	typ := field.Type()

	// If this is a pointer to another value, make sure that value is allocated.
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}
		field = field.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var val int64
		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, typ.Bits())
		}

		if err != nil {
			return err
		}

		field.SetInt(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var val uint64
		if val, err = strconv.ParseUint(value, 0, typ.Bits()); err != nil {
			return err
		}
		field.SetUint(val)

	case reflect.Bool:
		var val bool
		if val, err = strconv.ParseBool(value); err != nil {
			return err
		}
		field.SetBool(val)

	case reflect.Float32, reflect.Float64:
		var val float64
		if val, err = strconv.ParseFloat(value, typ.Bits()); err != nil {
			return err
		}
		field.SetFloat(val)

	case reflect.Slice:
		sl := reflect.MakeSlice(typ, 0, 0)
		if typ.Elem().Kind() == reflect.Uint8 {
			sl = reflect.ValueOf([]byte(value))
		} else if strings.TrimSpace(value) != "" {
			vals := strings.Split(value, ",")
			sl = reflect.MakeSlice(typ, len(vals), len(vals))
			for i, val := range vals {
				if err = Parse(val, sl.Index(i)); err != nil {
					return err
				}
			}
		}
		field.Set(sl)

	case reflect.Map:
		mp := reflect.MakeMap(typ)
		if strings.TrimSpace(value) != "" {
			pairs := strings.Split(value, ",")
			for _, pair := range pairs {
				kvpair := strings.Split(pair, ":")
				if len(kvpair) != 2 {
					return fmt.Errorf("invalid map item: %q", pair)
				}

				k := reflect.New(typ.Key()).Elem()
				if err = Parse(kvpair[0], k); err != nil {
					return nil
				}

				v := reflect.New(typ.Elem()).Elem()
				if err = Parse(kvpair[1], v); err != nil {
					return nil
				}

				mp.SetMapIndex(k, v)
			}
		}
		field.Set(mp)
	}

	return nil
}
