package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
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

func Ok(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, Response{
		Message: message,
		Status:  code,
		Result:  data,
	})
}

func HandleJSON(c echo.Context) (map[string]interface{}, bool) {
	_json := make(map[string]interface{})
	json.NewDecoder(c.Request().Body).Decode(&_json)

	if len(_json) == 0 {
		return nil, true
	}

	return _json, false
}
