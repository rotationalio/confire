package merge

import (
	"fmt"
	"reflect"

	"github.com/rotationalio/confire/structs"
)

// Merge updates the destination struct in-place with non-zero values from the source.
// Only fields with the same name and type get updated (the structs technically do not
// have to be of the same type, though this is untested as the goal is to merge structs
// of the same type). Note that because non-zero values are skipped to merge a non-zero
// value (such as false) into the target, you either need to use a pointer such that
// nil is the non-zero value or use patch.
//
// Returns true if any value has been changed on the destination struct.
func Merge(dst, src interface{}) (changed bool, err error) {
	// Find all of the correlated fields between src and dst, ensuring that embedded
	// structs are correctly added and that nested structs are recursively merged.
	// Gathering pairs also skips unexported fields on the src and unsettable fields on
	// the dst since these will not be able to be merged.
	var pairs []*fieldPair
	if pairs, err = getMergePairs(dst, src); err != nil {
		return changed, err
	}

	for _, pair := range pairs {
		dstf, srcf := pair.dstf, pair.srcf

		// Do not merge zero-valued fields
		if srcf.IsZero() {
			continue
		}

		// Check the types to ensure that the fields match
		if skind, dkind := srcf.Kind(), dstf.Kind(); skind != dkind {
			return changed, fmt.Errorf(`field %s type mismatch while merging %s vs %s`, pair.name, dkind, skind)
		}

		// Test to see if the struct has changed; this has to be done before Set but if
		// Set errors there is an edge case where changed is incorrect.
		if !changed && !reflect.DeepEqual(srcf.Value(), dstf.Value()) {
			changed = true
		}

		// Set the dstf with the value of the src
		if err = dstf.Set(srcf.Value()); err != nil {
			return changed, fmt.Errorf("could not set field %s: %w", dstf.Name(), err)
		}
	}
	return changed, nil
}

type fieldPair struct {
	name string
	srcf *structs.Field
	dstf *structs.Field
}

func getMergePairs(dst, src interface{}) (mergePairs []*fieldPair, err error) {
	var target, patch *structs.Struct
	if target, err = structs.New(dst); err != nil {
		return nil, err
	}

	if patch, err = structs.New(src); err != nil {
		return nil, err
	}

	// Search the src struct for fields to patch into the dst struct
	fields := patch.Fields()
	for _, srcf := range fields {
		// Skip unexported fields
		if !srcf.IsExported() {
			continue
		}

		// Identify the target field
		var dstf *structs.Field
		if dstf, err = target.Field(srcf.Name()); err != nil {
			continue
		}

		// Add embedded fields to the pairs
		if srcf.IsEmbedded() {
			var embeddedPairs []*fieldPair
			if embeddedPairs, err = getMergePairs(dstf.Value(), srcf.Value()); err != nil {
				continue
			}
			mergePairs = append(mergePairs, embeddedPairs...)
			continue
		}

		// Otherwise create pair from src and dst fields
		mergePairs = append(mergePairs, &fieldPair{srcf.Name(), srcf, dstf})
	}

	return mergePairs, nil
}
