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
	form := new(structs.UserCreateEditForm)
	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	// Check user email
	_, userDetailStatusCode, _ := uc.UserService.Show(structs.UserAttrsFind{Email: form.Email})
	if userDetailStatusCode == 200 {
		return helpers.ErrorBadRequest("Email already used")
	}

	// Create new user
	createdUser, err := uc.UserService.StoreOrUpdate(structs.UserCreateEditForm{
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
	id, _ := (strconv.Atoi(c.Param("id")))
	uintId := uint(id)

	form := new(structs.UserCreateEditForm)
	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	// Get user detail by id
	userDetailById, userDetailByIdStatusCode, err := uc.UserService.Show(structs.UserAttrsFind{ID: uintId})

	// Check email
	if userDetailById.Email != form.Email {
		// Check again if user update the email it-self
		userDetailByEmail, _, _ := uc.UserService.Show(structs.UserAttrsFind{Email: form.Email})
		if userDetailByEmail.Email == form.Email {
			return helpers.ErrorBadRequest("Email already used")

		}
	}

	// Check if user that will be updated not found
	if userDetailByIdStatusCode >= 400 && err != nil {
		return helpers.ErrorDynamic(userDetailByIdStatusCode, err.Error())
	}

	// Update selected user
	updatedUser, err := uc.UserService.StoreOrUpdate(structs.UserCreateEditForm{
		ID:    &userDetailById.ID,
		Name:  form.Name,
		Email: form.Email,
		Roles: form.Roles,
	})
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusCreated, "User updated successfully", updatedUser)
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
