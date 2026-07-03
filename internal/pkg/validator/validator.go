package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validate *validator.Validate
}

func New() *StructValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			return field.Name
		}
		return strings.Split(tag, ",")[0]
	})
	return &StructValidator{validate: v}
}

func (sv *StructValidator) Validate(out any) error {
	if err := sv.validate.Struct(out); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fmt.Errorf("validation failed: %s", formatValidationErrors(validationErrors))
		}
		return err
	}
	return nil
}

func formatValidationErrors(validationErrors validator.ValidationErrors) string {
	parts := make([]string, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		parts = append(parts, fmt.Sprintf("%s is %s", fieldErr.Field(), fieldErr.ActualTag()))
	}
	return strings.Join(parts, ", ")
}
