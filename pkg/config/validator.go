package config

import validation "github.com/go-ozzo/ozzo-validation/v4"

func Validate(v any) error {
	oldTag := validation.ErrorTag
	validation.ErrorTag = "yaml"
	defer func() {
		validation.ErrorTag = oldTag
	}()
	return validation.Validate(v)
}

func ValidateStruct(v any, fields ...*validation.FieldRules) error {
	oldTag := validation.ErrorTag
	validation.ErrorTag = "yaml"
	defer func() {
		validation.ErrorTag = oldTag
	}()
	return validation.ValidateStruct(v, fields...)
}
