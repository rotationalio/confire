package patch

import (
	"fmt"
	"reflect"

	"github.com/rotationalio/confire/structs"
)

// Patch updates the destination struct in-place with non-zero values from the source.
// Only fields with the same name and type get updated (the structs technically do not
// have to be of the same type, though this is the tested methodology).
//
// Returns true if any value has been changed on the destination struct.
func Patch(dst, src interface{}) (changed bool, err error) {
	var (
		target, patch *structs.Struct
		fields        []*structs.Field
	)

	if target, err = structs.New(dst); err != nil {
		return changed, err
	}

	if patch, err = structs.New(src); err != nil {
		return changed, err
	}

	// Continue to add nested fields to the work stack
	fields = patch.Fields()
	for n := len(fields); n > 0; n = len(fields) {
		var srcf, dstf *structs.Field

		// pop the top off of the work stack
		srcf = fields[n-1]
		fields = fields[:n-1]

		// skip unexported fields
		if !srcf.IsExported() {
			continue
		}

		// add embedded fields into the work stack
		if srcf.IsEmbedded() {
			fields = append(fields, srcf.Fields()...)
		}

		// skip zero-valued fields
		if srcf.IsZero() {
			continue
		}

		// attempt to get the target field by name
		name := srcf.Name()
		if dstf, err = target.Field(name); err != nil {
			// skip fields that do not exist
			continue
		}

		// if the src value is a struct then add it to work stack
		// if structs.IsStruct(srcf.Value()) && structs.IsStruct(dstf.Value()) {
		// 	// recursively patch nested structs
		// 	var nestedChanged bool
		// 	if nestedChanged, err = Patch(dstf.Value(), srcf.Value()); err != nil {
		// 		return changed, err
		// 	}

		// 	// only update if nested has changed to prevent overriding changed bool
		// 	if nestedChanged {
		// 		changed = true
		// 	}
		// 	continue
		// }

		// fetch the src value
		srcv := reflect.Indirect(reflect.ValueOf(srcf.Value()))

		// check the types to make sure we can set the value
		if skind, dkind := srcv.Kind(), dstf.Kind(); skind != dkind {
			return changed, fmt.Errorf(`field %s type mismatch while merging %s vs %s`, name, dkind, skind)
		}

		// test to see if the struct has changed
		if !changed && !reflect.DeepEqual(srcv.Interface(), dstf.Value()) {
			changed = true
		}

		// attempt to set the dstf with the value
		if err = dstf.Set(srcv.Interface()); err != nil {
			return changed, err
		}
	}

	return changed, nil
}
