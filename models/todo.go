package models

import (
	"echo-rest/database"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type TodoForm struct {
	Title     string `json:"title" validate:"required"`
	Completed bool   `json:"completed"`
}

type Todo struct {
	ID        uint      `gorm:"primaryKey" column:"id" json:"id"`
	Title     string    `gorm:"title" json:"title"`
	Completed bool      `gorm:"completed" json:"completed"`
	CreatedAt time.Time `gorm:"column:createdAt;" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (Todo) TableName() string {
	return "todos"
}

func (Todo) GetList() ([]Todo, error) {
	db := database.DatabaseConnection()
	var todoList []Todo

	query := db.Find(&todoList)

	return todoList, query.Error
}

func (Todo) Store(payload TodoForm) (Todo, error) {
	db := database.DatabaseConnection()
	todo := Todo{
		Title:     payload.Title,
		Completed: payload.Completed,
	}

	query := db.Create(&todo)

	return todo, query.Error
}

func (Todo) GetDetail(id int) (Todo, int, error) {
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

	return todo, statusCode, query.Error
}

func (Todo) Update(todo Todo) error {
	db := database.DatabaseConnection()

	query := db.Save(&todo)

	return query.Error
}

func (Todo) Delete(todo Todo) error {
	db := database.DatabaseConnection()

	query := db.Delete(&todo)

	return query.Error
}
