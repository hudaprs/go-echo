package helpers

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func HandleResponse(c echo.Context, response Response) error {
	return c.JSON(response.Status, Response{
		Status:  response.Status,
		Message: response.Message,
		Result:  response.Result,
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
