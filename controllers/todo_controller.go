package controllers

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TodoController struct {
}

// @description Get todo list
// @param 		echo.Context
// @return		error
func (TodoController) Index(c echo.Context) error {
	todo := models.Todo{}
	authenticatedUser := helpers.JwtGetClaims(c)
	pagination := helpers.SetPagination(c, helpers.Pagination{})
	todoList, err := todo.GetList(authenticatedUser.ID, pagination)

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get todo list success", todoList)
}

// @description Create todo
// @param 		is echo.Context
// @return 		error
func (TodoController) Store(c echo.Context) error {
	form := new(models.TodoForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	todo := models.Todo{}
	authenticatedUser := helpers.JwtGetClaims(c)
	createdTodo, err := todo.Store(models.TodoForm{
		Title:     form.Title,
		Completed: false,
		UserID:    authenticatedUser.ID,
	})

	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(http.StatusCreated, "Todo created successfully", createdTodo)
}

// @description Get todo detail
// @param 		echo.Context
// @return		error
func (TodoController) Show(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todo := models.Todo{}
	todoDetail, statusCode, err := todo.GetDetail(id)

	// Check if not correct user
	if !todo.IsCorrectUser(authenticatedUser.ID, todoDetail) {
		return helpers.ErrorForbidden()
	}

	// Check if something went wrong when trying to lookup todo detail
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get todo success", todoDetail)
}

// @description Update todo
// @param 		echo.Context
// @return		error
func (TodoController) Update(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)
	form := new(models.TodoForm)

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

	todo := models.Todo{}
	todoDetail, statusCode, err := todo.GetDetail(id)

	// Check if not correct user
	if !todo.IsCorrectUser(authenticatedUser.ID, todoDetail) {
		return helpers.ErrorForbidden()
	}

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	// Update record
	todoDetail.Title = form.Title
	todoDetail.Completed = form.Completed

	// Save record
	updateErr := todo.Update(todoDetail)

	if updateErr != nil {
		return helpers.ErrorServer(updateErr.Error())
	}

	return helpers.Ok(http.StatusOK, "Todo updated successfully", todoDetail)
}

// @description Delete todo
// @param 		echo.Context
// @return		error
func (TodoController) Delete(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todo := models.Todo{}
	todoDetail, statusCode, err := todo.GetDetail(id)

	// Check if not correct user
	if !todo.IsCorrectUser(authenticatedUser.ID, todoDetail) {
		return helpers.ErrorForbidden()
	}

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	deleteErr := todo.Delete(todoDetail)

	if deleteErr != nil {
		return helpers.ErrorServer(deleteErr.Error())
	}

	return helpers.Ok(http.StatusOK, "Todo deleted successfully", todoDetail)
}
