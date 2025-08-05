package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgtype"
)

// `validator.Validate` singleton caches struct info
var v *validator.Validate

func main() {
	// Opt-in to v11+ behavior
	v = validator.New(validator.WithRequiredStructEnabled())

	// Internally [validator.Validate] uses reflection to support data validation.
	// The implementation need to support beyond primitive data types, it needs to be able to handle complex struct meaningfully.
	// `time.Time` is a good example of complex struct that might be used to represent a piece of information.
	//
	// In order to aid [validator.Validate] to extract custom type value that will be run through validation rules
	// we need to provide CustomTypeFunc that will be used to provide underlying value against custom types registered with it.

	// As mention in docs, RegisterCustomTypeFunc is not thread-safe it is intended that these all be registered prior to any validation

	// Why these types: https://github.com/nadhifikbarw/x-go-painless-null/
	// Register `sql.NullXxX` types
	v.RegisterCustomTypeFunc(
		// ValidateValuer provides stricter behavior that nil value
		// always come from underlying `Valuer` field, since it panics
		// when Valuer return non-nill error value leaving us with no
		// good value as fallback that doesn't semantically being used
		//
		// See [NullValidateValuer] for alternative behavior where it
		// returns nil as fallback instead of panic.
		ValidateValuer,
		sql.NullBool{},
		sql.NullByte{},
		sql.NullFloat64{},
		sql.NullInt16{},
		sql.NullInt32{},
		sql.NullInt64{},
		sql.NullString{},
		sql.NullTime{},
	)

	// Register `guregu/null.XxX` types
	v.RegisterCustomTypeFunc(
		ValidateValuer,
		null.Bool{},
		null.Byte{},
		null.Float{},
		null.Int{},
		null.Int16{},
		null.Int32{},
		null.Int64{},
		null.String{},
		null.Time{},
	)

	// Register `pgtype.XxX` types (Non-exhaustive)
	v.RegisterCustomTypeFunc(
		NullValidateValuer,
		pgtype.Bool{},
		pgtype.Float4{},
		pgtype.Int4{},
		pgtype.Text{},
		pgtype.Timestamp{},
		pgtype.Timestamptz{},
	)

	// Example struct using Valuer types
	example := struct {
		Name         pgtype.Text `validate:"required"`
		Description  sql.NullString
		FailingField null.String `validate:"omitempty,max=1"`
		Category     null.String `validate:"omitempty,max=1"`
	}{
		Name:        pgtype.Text{String: "Example Name", Valid: true},
		Description: sql.NullString{}, // Null value
		// Deliberately provide failing input
		FailingField: null.NewString("Failing fail", true),
		Category:     null.NewString("", false),
	}

	// This should return meaningful validation errors
	if err := v.Struct(example); err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range e {
				slog.Info(
					fmt.Sprintf("Error: %s", fieldError.Error()),
				)
			}
		}
	}
}
