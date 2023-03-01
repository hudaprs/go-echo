package controllers

import (
	"go-echo/helpers"
	"go-echo/services"
	"go-echo/structs"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TodoController struct {
	TodoService services.TodoService
}

// @description Get data list
// @param 		echo.Context
// @return		error
func (tc TodoController) Index(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)
	pagination := helpers.SetPagination(c, helpers.Pagination{})
	todoList, err := tc.TodoService.Index(authenticatedUser.ID, pagination)

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get todo list success", todoList)
}

// @description Store data
// @param 		echo.Context
// @return		error
func (tc TodoController) Store(c echo.Context) error {
	form := new(structs.TodoForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	authenticatedUser := helpers.JwtGetClaims(c)
	createdTodo, err := tc.TodoService.Store(structs.TodoForm{
		Title:     form.Title,
		Completed: false,
		UserID:    authenticatedUser.ID,
	})

	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(http.StatusCreated, "Todo created successfully", createdTodo)
}

// @description Get single data
// @param 		echo.Context
// @return		error
func (tc TodoController) Show(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todoDetail, statusCode, err := tc.TodoService.Show(id)

	// Check if not correct user
	if !tc.TodoService.IsCorrectUser(authenticatedUser.ID, todoDetail.UserID) {
		return helpers.ErrorForbidden()
	}

	// Check if something went wrong when trying to lookup todo detail
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get todo success", todoDetail)
}

// @description Update data
// @param 		echo.Context
// @return		error
func (tc TodoController) Update(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)
	form := new(structs.TodoForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todoDetail, statusCode, err := tc.TodoService.Show(id)

	// Check if not correct user
	if !tc.TodoService.IsCorrectUser(authenticatedUser.ID, todoDetail.UserID) {
		return helpers.ErrorForbidden()
	}

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	// Update record
	todoDetail.Title = form.Title
	todoDetail.Completed = form.Completed

	// Save record
	updateErr := tc.TodoService.Update(todoDetail)

	if updateErr != nil {
		return helpers.ErrorServer(updateErr.Error())
	}

	return helpers.Ok(http.StatusOK, "Todo updated successfully", todoDetail)
}

// @description Delete data
// @param 		echo.Context
// @return		error
func (tc TodoController) Delete(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todoDetail, statusCode, err := tc.TodoService.Show(id)

	// Check if not correct user
	if !tc.TodoService.IsCorrectUser(authenticatedUser.ID, todoDetail.UserID) {
		return helpers.ErrorForbidden()
	}

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	deleteErr := tc.TodoService.Delete(todoDetail)

	if deleteErr != nil {
		return helpers.ErrorServer(deleteErr.Error())
	}

	return helpers.Ok(http.StatusOK, "Todo deleted successfully", todoDetail)
}
