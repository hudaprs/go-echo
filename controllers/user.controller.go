package controllers

import (
	"go-echo/helpers"
	"go-echo/locales"
	"go-echo/services"
	"go-echo/structs"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService services.UserService
	RoleService services.RoleService
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

	return helpers.Ok(http.StatusOK, locales.LocalesGet("user.rest.index"), users)
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
		return helpers.ErrorBadRequest(locales.LocalesGet("validation.emailAlreadyUsed"))
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

	return helpers.Ok(http.StatusCreated, locales.LocalesGet("user.rest.store"), createdUser)
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

	return helpers.Ok(http.StatusOK, locales.LocalesGet("user.rest.show"), user)
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
			return helpers.ErrorBadRequest(locales.LocalesGet("validation.emailAlreadyUsed"))

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

	return helpers.Ok(http.StatusCreated, locales.LocalesGet("user.rest.update"), updatedUser)
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

	return helpers.Ok(http.StatusOK, locales.LocalesGet("user.rest.destroy"), user)
}

// @description Dropdown data
// @param 		echo.Context
// @return		error
func (uc UserController) RoleDropdown(c echo.Context) error {
	roleList, err := uc.RoleService.Dropdown()
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, locales.LocalesGet("role.rest.index"), roleList)
}
