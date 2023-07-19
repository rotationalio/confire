/*
Package env is a port of the github.com/kelseyhightower/envconfig package that is
modified for use in confire's multi-config system. This package decodes environment
variables based on a user defined specification in order to configure a struct from the
environment with required or default values.
*/
package env

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/rotationalio/confire/errors"
	"github.com/rotationalio/confire/structs"
)

const (
	tagIgnored    = "ignored"
	tagSplitWords = "split_words"
	tagEnv        = "envconfig"
	tagDefault    = "default"
	tagRequired   = "required"
)

type Info struct {
	Name  string         // Name of the field to compute the envvar from
	Alt   string         // String specified by the env tag
	Key   string         // The final environment variable key determined by the algorithm
	Field *structs.Field // The actual field to set the envvar from (along with tags)
}

var gatherRegexp = regexp.MustCompile("([^A-Z]+|[A-Z]+[^A-Z]+|[A-Z]+)")
var acronymRegexp = regexp.MustCompile("([A-Z]+)([A-Z][^A-Z]+)")

func Gather(prefix string, spec interface{}) (infos []Info, err error) {
	var s *structs.Struct
	if s, err = structs.New(spec); err != nil {
		return nil, errors.ErrInvalidSpecification
	}

	if !s.IsPointer() {
		return nil, errors.ErrInvalidSpecification
	}

	infos = make([]Info, 0, s.NumField())
	for _, field := range s.Fields() {
		// Skip any ignored fields or fields that cannot be set.
		if !field.CanSet() || isTrue(field.Tag(tagIgnored)) {
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
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		// Capture information about the config variable
		info := Info{
			Name:  field.Name(),
			Alt:   strings.ToUpper(field.Tag(tagEnv)),
			Field: field,
		}

		// Default to the field name as the envvar name (will be upcased)
		info.Key = info.Name

		// Best effort to un-pick camel casing as separate words
		if isTrue(field.Tag(tagSplitWords)) {
			words := gatherRegexp.FindAllStringSubmatch(field.Name(), -1)
			if len(words) > 0 {
				var name []string
				for _, words := range words {
					if m := acronymRegexp.FindStringSubmatch(words[0]); len(m) == 3 {
						name = append(name, m[1], m[2])
					} else {
						name = append(name, words[0])
					}
				}
				info.Key = strings.Join(name, "_")
			}
		}

		if info.Alt != "" {
			info.Key = info.Alt
		}

		if prefix != "" {
			info.Key = fmt.Sprintf("%s_%s", prefix, info.Key)
		}

		info.Key = strings.ToUpper(info.Key)
		infos = append(infos, info)

		if field.Kind() == reflect.Struct {
			// TODO: honor Decode if present
			innerPrefix := prefix
			if !field.IsEmbedded() {
				innerPrefix = info.Key
			}

			embeddedPtr := field.Pointer()
			embeddedInfos, err := Gather(innerPrefix, embeddedPtr)
			if err != nil {
				return nil, err
			}

			infos = append(infos[:len(infos)-1], embeddedInfos...)
			continue
		}

	}

	return infos, nil
}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
