package controllers

import (
	"echo-rest/helpers"
	"echo-rest/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService services.UserService
}

// @description Get data list
// @param 		echo.Context
// @return		error
func (uc UserController) Index(c echo.Context) error {
	pagination := helpers.SetPagination(c, helpers.Pagination{})
	users, err := uc.UserService.Index(pagination)

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get data list success", users)
}

// @description Store data
// @param 		echo.Context
// @return		error
func (uc UserController) Store(c echo.Context) error {
	return helpers.Ok(http.StatusCreated, "Data created successfully", nil)
}

// @description Get single data
// @param 		echo.Context
// @return		error
func (uc UserController) Show(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Get single data success", nil)
}

// @description Update data
// @param 		echo.Context
// @return		error
func (uc UserController) Update(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Update data success", nil)
}

// @description Delete data
// @param 		echo.Context
// @return		error
func (uc UserController) Delete(c echo.Context) error {
	return helpers.Ok(http.StatusOK, "Delete data success", nil)
}
