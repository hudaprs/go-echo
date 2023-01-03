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
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	return helpers.HandleResponse(c, helpers.Response{
		Message: "Get todo list success",
		Status:  http.StatusOK,
		Result:  todoList,
	})
}

/*
* @description Create todo
*
* @param {echo.Context} c
*
* @return error
 */
func (TodoController) Store(c echo.Context) error {
	// Simple validation for payload
	form, isNotValid := helpers.HandleJSON(c)
	title, ok := form["title"].(string)

	if (isNotValid && !ok) || title == "" {
		return helpers.HandleResponse(c, helpers.Response{
			Message: "Please fill all forms",
			Status:  http.StatusUnprocessableEntity,
			Result: map[string]string{
				"title": "Title is required",
			},
		})
	}

	todo, err := models.CreateTodo(models.ITodoForm{
		Title:     title,
		Completed: false,
	})

	if err != nil {
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
	}

	return helpers.HandleResponse(c, helpers.Response{
		Message: "User successfully created",
		Status:  http.StatusCreated,
		Result:  todo,
	})
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
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
	}

	todo, err, statusCode := models.GetTodo(id)

	if err != nil && statusCode >= 400 {
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  statusCode,
		})
	}

	return helpers.HandleResponse(c, helpers.Response{
		Message: "Get todo success",
		Status:  http.StatusOK,
		Result:  todo,
	})
}

/*
* @description Update todo
*
* @param {echo.Context} c
*
* @return error
 */
func (TodoController) Update(c echo.Context) error {
	// Simple validation for payload
	form, isNotValid := helpers.HandleJSON(c)
	title, titleOk := form["title"].(string)
	completed, completedOk := form["completed"].(bool)

	if (isNotValid && (!titleOk || !completedOk)) || title == "" {
		return helpers.HandleResponse(c, helpers.Response{
			Message: "Please fill all forms",
			Status:  http.StatusUnprocessableEntity,
			Result: map[string]string{
				"title":     "Title is required",
				"completed": "Completed is required",
			},
		})
	}

	id, err := strconv.Atoi(c.Param("id"))

	// Check for query params
	if err != nil {
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
	}

	todo, err, statusCode := models.GetTodo(id)

	if err != nil && statusCode >= 400 {
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  statusCode,
		})
	}

	// Update record
	todo.Title = title
	todo.Completed = completed

	// Save record
	updateErr := models.UpdateTodo(todo)

	if updateErr != nil {
		return helpers.HandleResponse(c, helpers.Response{
			Message: updateErr.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	return helpers.HandleResponse(c, helpers.Response{
		Message: "Todo successfully updated",
		Status:  http.StatusOK,
		Result:  todo,
	})
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
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
	}

	todo, err, statusCode := models.GetTodo(id)

	if err != nil && statusCode >= 400 {
		return helpers.HandleResponse(c, helpers.Response{
			Message: err.Error(),
			Status:  statusCode,
		})
	}

	deleteErr := models.DeleteTodo(todo)

	if deleteErr != nil {
		return helpers.HandleResponse(c, helpers.Response{
			Message: deleteErr.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	return helpers.HandleResponse(c, helpers.Response{
		Message: "Todo successfully deleted",
		Status:  http.StatusOK,
		Result:  todo,
	})
}
