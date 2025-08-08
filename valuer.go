package main

import (
	"database/sql/driver"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var _ validator.CustomTypeFunc = ValidateValuer
var _ validator.CustomTypeFunc = NullValidateValuer

// ValidateValuer returns underlying value from [driver.Valuer] type for validation.
// Type registered using this must its nil/zero value and always return nil error value
//
// It panics if non-`Valuer` field is passed for validation.
// It panics when [driver.Valuer] returns non-nil error value.
//
// See [NullValidateValuer] for Valuer type that is allowed to return nil as fallback.
//
// ValidateValuer implements [validator.CustomTypeFunc] interface.
func ValidateValuer(field reflect.Value) interface{} {
	valuer, ok := field.Interface().(driver.Valuer)
	if !ok {
		panic("not a Valuer field")
	}

	// Panic when `Valuer` return non-nill error
	// `Valuer` registered against this function
	// typically attach semantic meaning to nil value
	// hence it panics to avoid incorrect validation behavior
	val, err := valuer.Value()
	if err != nil {
		panic("Valuer field returns non-nil error")
	}

	return val
}

// NullValidateValuer returns underlying value from [driver.Valuer] type for validation.
// It returns nil as fallback when field returns non-nil error value.
//
// It panics if non-`Valuer` field is passed for validation.
//
// See [ValidateValuer] for stricter behavior that panics when Valuer returns error
//
// NullValidateValuer implements [validator.CustomTypeFunc] interface.
func NullValidateValuer(field reflect.Value) interface{} {
	valuer, ok := field.Interface().(driver.Valuer)
	if !ok {
		panic("not a Valuer field")
	}

	// Return nil when `Valuer` return non-nill error
	// `Valuer` registered against this function typically doesn't
	// attach nil value semantically hence using nil value as fallback
	// value is permissible
	val, err := valuer.Value()
	if err != nil {
		return nil
	}

	return val
}
