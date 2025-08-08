# Validating Custom Type with `go-playground/validator`

[`go-playground/validator`](https://github.com/go-playground/validator) is typically the de-facto recommendation if you want ready-made validation in Go. I found some documentation  of its features can be too implicit for may taste and hard to parse through without reading its internals.

Here I expanded the example on how to configure validator to properly handle custom type, specifically allowing custom type that implements `database/sql/driver.Valuer` interface for validation.

## How `CustomTypeFunc` being maintained and used

> This explanation is v10 behavior, future internals may get changed.

Any custom type that need want to be validated need to be registered alongside its `validator.CustomTypeFunc`, this `CustomTypeFunc` is maintained as simple map.

Since Validator uses reflection internally, it's necessary to provide `CustomTypeFunc` for each custom type to provide the ability to get underlying field value from such custom type.

```go
func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{}) {
	if v.customFuncs == nil {
        // Allocate simple map
		v.customFuncs = make(map[reflect.Type]CustomTypeFunc)
	}

    // Map each registered type with specified [CustomTypeFunc]
	for _, t := range types {
		v.customFuncs[reflect.TypeOf(t)] = fn
	}

	v.hasCustomFuncs = true
}
```

Usage

```go
v.RegisterCustomTypeFunc(
    NullValidateValuer,
    pgtype.Bool{},
    pgtype.Float4{},
    pgtype.Int4{},
    pgtype.Text{},
    pgtype.Timestamp{},
    pgtype.Timestamptz{},
)
```

When validating a struct, validator accepts `interface{}` as input, it uses reflection to guide what each field type they're dealing with, if it encounters field of custom type, it then consults this internal `CustomTypeFunc` map and call the registered function to resolve the underlying value of the custom field.

To be more explicit, this value resolution is performed recursively. So your `CustomTypeFunc` may actually return  another custom type, the validator will keep trying to resolve the underlying value of the subsequent custom type until it encounters primitives that they recognize or no longer recognize a custom type that's unregistered.
