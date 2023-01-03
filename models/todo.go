package models

import (
	"echo-rest/database"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type ITodoForm struct {
	Title     string
	Completed bool
}

type Todo struct {
	ID        uint      `gorm:"primaryKey" column:"id" json:"id"`
	Title     string    `gorm:"title" json:"title"`
	Completed bool      `gorm:"completed" json:"completed"`
	CreatedAt time.Time `gorm:"column:createdAt;" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func GetTodoList() ([]Todo, error) {
	db := database.DatabaseConnection()
	var todoList []Todo

	query := db.Find(&todoList)

	return todoList, query.Error
}

func CreateTodo(payload ITodoForm) (Todo, error) {
	db := database.DatabaseConnection()
	todo := Todo{
		Title:     payload.Title,
		Completed: payload.Completed,
	}

	query := db.Create(&todo)

	return todo, query.Error
}

func GetTodo(id int) (Todo, error, int) {
	db := database.DatabaseConnection()

	var todo Todo

	query := db.First(&todo, id)

	isNotFound := errors.Is(query.Error, gorm.ErrRecordNotFound)

	var statusCode int

	if isNotFound {
		statusCode = http.StatusNotFound
	} else if query.Error != nil {
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = http.StatusOK
	}

	return todo, query.Error, statusCode
}

func UpdateTodo(todo Todo) error {
	db := database.DatabaseConnection()

	query := db.Save(&todo)

	return query.Error
}

func DeleteTodo(todo Todo) error {
	db := database.DatabaseConnection()

	query := db.Delete(&todo)

	return query.Error
}
