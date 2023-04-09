package structs

import "reflect"

type Zero interface {
	IsZero() bool
}

func interfaceFrom(field reflect.Value, fn func(interface{}, *bool)) {
	if !field.CanInterface() {
		return
	}

	var ok bool
	fn(field.Interface(), &ok)
	if !ok && field.CanAddr() {
		fn(field.Addr().Interface(), &ok)
	}
}

func zeroFrom(field reflect.Value) (z Zero) {
	interfaceFrom(field, func(v interface{}, ok *bool) { z, *ok = v.(Zero) })
	return z
}
