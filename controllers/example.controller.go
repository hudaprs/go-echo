package controllers

import (
	"go-echo/helpers"
	"go-echo/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ExampleController struct {
	ExampleService services.ExampleService
}

// @description Get data list
// @param 		echo.Context
// @return		error
func (ec ExampleController) Index(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Get data list success", nil)
}

// @description Store data
// @param 		echo.Context
// @return		error
func (ec ExampleController) Store(c echo.Context) error {
	return helpers.Ok(http.StatusCreated, "Data created successfully", nil)
}

// @description Get single data
// @param 		echo.Context
// @return		error
func (ec ExampleController) Show(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Get single data success", nil)
}

// @description Update data
// @param 		echo.Context
// @return		error
func (ec ExampleController) Update(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Update data success", nil)
}

// @description Delete data
// @param 		echo.Context
// @return		error
func (ec ExampleController) Delete(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Delete data success", nil)
}
