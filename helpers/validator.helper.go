package helpers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func renderCustomValidationMessage(actualMessage string, param string) string {
	switch actualMessage {
	case "email":
		return "valid email"
	case "gte":
		return "greater or equal than " + param
	default:
		return actualMessage
	}
}

func ValidateNotFoundData(err error) int {
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)

	var statusCode int

	if isNotFound {
		statusCode = http.StatusNotFound
	} else if err != nil {
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = http.StatusOK
	}

	return statusCode
}

func (customValidator *CustomValidator) Validate(i interface{}) error {

	if err := customValidator.Validator.Struct(i); err != nil {
		var validations []ValidationResponse

		for _, err := range err.(validator.ValidationErrors) {

			validations = append(validations, ValidationResponse{
				Field:   err.Field(),
				Message: err.Field() + " should be " + renderCustomValidationMessage(err.ActualTag(), err.Param()),
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
