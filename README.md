# Allowing custom fields with `go-playground/validator`

Better explained example to register `Valuer` types for validation with `go-playground/validator`

## Internal of CustomTypeFunc registration

Registered CustomTypeFunc maintained as simple map against registered custom types

```go
func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{}) {
	if v.customFuncs == nil {
		v.customFuncs = make(map[reflect.Type]CustomTypeFunc)
	}

	for _, t := range types {
		v.customFuncs[reflect.TypeOf(t)] = fn
	}

	v.hasCustomFuncs = true
}
```
