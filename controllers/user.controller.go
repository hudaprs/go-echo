package controllers

import (
	"echo-rest/helpers"
	"echo-rest/services"
	"echo-rest/structs"
	"net/http"
	"strconv"

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

	return helpers.Ok(http.StatusOK, "Get user list success", users)
}

// @description Store data
// @param 		echo.Context
// @return		error
func (uc UserController) Store(c echo.Context) error {
	form := new(structs.UserCreateForm)
	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	// Check user email
	isEmailExists, _, _, _ := uc.UserService.CheckEmail(form.Email)
	if isEmailExists {
		return helpers.ErrorBadRequest("Email already used")
	}

	// Create new user
	createdUser, err := uc.UserService.Store(structs.UserCreateForm{
		Name:  form.Name,
		Email: form.Email,
		Roles: form.Roles,
	})
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusCreated, "User created successfully", createdUser)
}

// @description Get single data
// @param 		echo.Context
// @return		error
func (uc UserController) Show(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Get single user
	user, statusCode, err := uc.UserService.Show(structs.UserAttrsFind{ID: uint(id)})
	if statusCode >= 400 && err != nil {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get user detail success", user)
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
	id, _ := strconv.Atoi(c.Param("id"))

	// Get single user
	user, statusCode, err := uc.UserService.Delete(structs.UserAttrsFind{ID: uint(id)})
	if statusCode >= 400 && err != nil {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "User deleted successfully", user)
}
