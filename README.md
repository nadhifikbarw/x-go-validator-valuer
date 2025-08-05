# Allowing custom fields with `go-playground/validator`

`go-playground/validator` is great, but documentation on some of its features can be too implicit and hard to parse through without reading its internals

Here's more explicitly documented example on how to register custom types to allow types such as `database/sql/driver.Valuer` for validation with `go-playground/validator`

## How CustomTypeFunc being maintained and used

> This is v10 behavior explanation, in case future internals get changed.

For each type registered using a particular `validator.CustomTypeFunc`, it's maintained as simple map between type and those function. It explains why you need to list all your initialized struct when calling this function

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

Whenever validator perform validation, validation accept `interface{}` as its input, so it uses reflection to guide how they would extract underlying value that need to be validated based on its validation tag.

Validator keeps track whether any custom type functions has been registered. It will then check whether CustomTypeFunc for a particular type exists to resolve underlying value of particular type.

Since validator maintain these CustomTypeFunc on type-level reflection, you're then required to register all custom type that need to be validated using validator.
