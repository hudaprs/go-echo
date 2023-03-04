package services

import (
	"go-echo/helpers"
	"go-echo/locales"
	"go-echo/models"
	"go-echo/structs"

	"gorm.io/gorm"
)

type TodoService struct {
	DB *gorm.DB
}

func (t *TodoService) Index(userId uint, pagination helpers.Pagination) (*helpers.Pagination, error) {
	var todoList []models.TodoUserResponse

	query := t.DB.Scopes(helpers.Paginate(todoList, &pagination, t.DB)).Preload("User").Where(&models.Todo{UserID: userId}).Find(&todoList)
	pagination.Rows = todoList

	return &pagination, query.Error
}

func (t *TodoService) Store(payload structs.TodoForm) (models.TodoResponse, error) {
	todo := models.TodoResponse{
		Title:     payload.Title,
		Completed: payload.Completed,
		UserID:    payload.UserID,
	}

	query := t.DB.Create(&todo)

	return todo, query.Error
}

func (t *TodoService) Show(id int) (models.TodoUserResponse, int, error) {
	var todo models.TodoUserResponse

	query := t.DB.Preload("User").First(&todo, id)
	statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
		NotFound: locales.LocalesGet("todo.validation.notFound"),
	})

	return todo, statusCode, err
}

func (t *TodoService) Update(todo models.TodoUserResponse) error {
	query := t.DB.Save(&todo)

	return query.Error
}

func (t *TodoService) Delete(todo models.TodoUserResponse) error {
	query := t.DB.Delete(&todo)

	return query.Error
}

func (t *TodoService) IsCorrectUser(userId uint, userTodoId uint) bool {
	return userId == userTodoId
}
