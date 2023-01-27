package controllers

import (
	"echo-rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ExampleController struct {
}

// @description Get data list
// @param 		echo.Context
// @return		error
func (ExampleController) Index(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Get data list success", nil)
}

// @description Store data
// @param 		echo.Context
// @return		error
func (ExampleController) Store(c echo.Context) error {
	return helpers.Ok(http.StatusCreated, "Data created successfully", nil)

}

// @description Get single data
// @param 		echo.Context
// @return		error
func (ExampleController) Show(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Get single data success", nil)
}

// @description Update data
// @param 		echo.Context
// @return		error
func (ExampleController) Update(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Update data success", nil)
}

// @description Delete data
// @param 		echo.Context
// @return		error
func (ExampleController) Delete(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Delete data success", nil)
}
