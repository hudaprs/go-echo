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
	todoList, err := todo.GetList()

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(c, http.StatusOK, "Get todo list success", todoList)
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
	createdTodo, err := todo.Store(models.TodoForm{
		Title:     form.Title,
		Completed: false,
	})

	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(c, http.StatusCreated, "Todo created successfully", createdTodo)
}

// @description Get todo detail
// @param 		echo.Context
// @return		error
func (TodoController) Show(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todo := models.Todo{}
	todoDetail, statusCode, err := todo.GetDetail(id)

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(c, http.StatusOK, "Get todo success", todoDetail)
}

// @description Update todo
// @param 		echo.Context
// @return		error
func (TodoController) Update(c echo.Context) error {
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

	return helpers.Ok(c, http.StatusOK, "Todo updated successfully", todoDetail)

}

// @description Delete todo
// @param 		echo.Context
// @return		error
func (TodoController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todo := models.Todo{}
	todoDetail, statusCode, err := todo.GetDetail(id)

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	deleteErr := todo.Delete(todoDetail)

	if deleteErr != nil {
		return helpers.ErrorServer(deleteErr.Error())
	}

	return helpers.Ok(c, http.StatusOK, "Todo deleted successfully", todoDetail)
}
