package helpers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type ValidationResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type DatabaseDynamicMessage struct {
	NotFound string
}

func ErrorBadRequest(message string) error {
	return echo.NewHTTPError(http.StatusBadRequest, Response{
		Message: message,
		Status:  http.StatusBadRequest,
	})
}

func ErrorValidation(validations []ValidationResponse) error {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, Response{
		Message: "Validation error",
		Status:  http.StatusUnprocessableEntity,
		Result:  validations,
	})
}

func ErrorServer(message string) error {
	return echo.NewHTTPError(http.StatusInternalServerError, Response{
		Message: message,
		Status:  http.StatusInternalServerError,
	})
}

func ErrorDynamic(code int, message string) error {
	return echo.NewHTTPError(code, Response{
		Message: message,
		Status:  code,
	})
}

func ErrorUnauthorized() error {
	return echo.NewHTTPError(http.StatusUnauthorized, Response{
		Message: "Unauthorized",
		Status:  http.StatusUnauthorized,
	})
}

func ErrorForbidden() error {
	return echo.NewHTTPError(http.StatusForbidden, Response{
		Message: "You don't have access to this content!",
		Status:  http.StatusForbidden,
	})
}

func Ok(code int, message string, data interface{}) error {
	return echo.NewHTTPError(code, Response{
		Message: message,
		Status:  code,
		Result:  data,
	})
}

func ErrorDatabaseDynamic(queryError error, message DatabaseDynamicMessage) (int, error) {
	isNotFound := errors.Is(queryError, gorm.ErrRecordNotFound)

	var statusCode int

	if isNotFound {
		statusCode = http.StatusNotFound
		if message.NotFound != "" {
			queryError = errors.New(message.NotFound)
		}
	} else if queryError != nil {
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = http.StatusOK
	}

	return statusCode, queryError
}
