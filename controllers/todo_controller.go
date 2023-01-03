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

/**
* @description Get todo list
*
* @param {echo.Context} c
*
* @return error
 */
func (TodoController) Index(c echo.Context) error {
	todoList, err := models.GetTodoList()

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

	todo, err := models.CreateTodo(models.TodoForm{
		Title:     form.Title,
		Completed: false,
	})

	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(c, http.StatusCreated, "Todo created successfully", todo)
}

/*
* @description Get single todo
*
* @param {echo.Context} c
*
* @return error
 */
func (TodoController) Show(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todo, statusCode, err := models.GetTodo(id)

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(c, http.StatusOK, "Get todo success", todo)
}

/*
* @description Update todo
*
* @param {echo.Context} c
*
* @return error
 */
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

	todo, statusCode, err := models.GetTodo(id)

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	// Update record
	todo.Title = form.Title
	todo.Completed = form.Completed

	// Save record
	updateErr := models.UpdateTodo(todo)

	if updateErr != nil {
		return helpers.ErrorServer(updateErr.Error())
	}

	return helpers.Ok(c, http.StatusOK, "Todo updated successfully", todo)

}

/*
* @description Delete todo
*
* @param {echo.Context} c
*
* @return error
 */
func (TodoController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	todo, statusCode, err := models.GetTodo(id)

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	deleteErr := models.DeleteTodo(todo)

	if deleteErr != nil {
		return helpers.ErrorServer(deleteErr.Error())
	}

	return helpers.Ok(c, http.StatusOK, "Todo deleted successfully", todo)

}
