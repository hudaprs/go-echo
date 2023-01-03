package helpers

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (customValidator *CustomValidator) Validate(i interface{}) error {

	if err := customValidator.Validator.Struct(i); err != nil {
		var validations []ValidationResponse

		for _, err := range err.(validator.ValidationErrors) {
			validations = append(validations, ValidationResponse{
				Field:   err.Field(),
				Message: err.Field() + " should be " + err.ActualTag(),
			})

			// fmt.Println(err.Namespace()) // can differ when a custom TagNameFunc is registered or
			// fmt.Println(err.Field())     // by passing alt name to ReportError like below
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
		}

		// Optionally, you could return the error to give each route more control over the status code
		return ErrorValidation(validations)

	}
	return nil
}
