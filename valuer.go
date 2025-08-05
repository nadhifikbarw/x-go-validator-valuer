package main

import (
	"database/sql/driver"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var _ validator.CustomTypeFunc = ValidateValuer
var _ validator.CustomTypeFunc = NullValidateValuer

// ValidateValuer returns underlying value from [driver.Valuer] type for validation.
//
// It panics if non-`Valuer` field is passed for validation.
// It panics when [driver.Valuer] returns non-nil error value.
//
// See [NullValidateValuer] if you want to fallback returning nil when non-nil error returned.
//
// ValidateValuer implements [validator.CustomTypeFunc] interface.
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		if val, err := valuer.Value(); err == nil {
			return val
		}
		// Panic when `Valuer` return non-nill error, because `Valuer` registered against this function
		// attach semantic meaning to nil value, panic to avoid implicit false validation behavior
		panic("Valuer field returns non-nil error")
	}
	panic("not a Valuer field")
}

// NullValidateValuer returns underlying value from [driver.Valuer] type for validation.
// It returns nil as fallback when field returns non-nil error value.
//
// It panics if non-`Valuer` field is passed for validation.
//
// See [ValidateValuer] for strict behavior to handle non-nil error from
//
// NullValidateValuer implements [validator.CustomTypeFunc] interface.
func NullValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		if val, err := valuer.Value(); err == nil {
			return val
		}
		// Returns nil as fallback when receiving non-nill error
		return nil
	}
	panic("not a Valuer field")
}
