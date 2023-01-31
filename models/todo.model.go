package models

import (
	"echo-rest/database"
	"echo-rest/helpers"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type TodoForm struct {
	Title     string `json:"title" validate:"required"`
	Completed bool   `json:"completed"`
	UserID    uint   `json:"userId"`
}

type Todo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"userId"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"user,omitempty"`
	Title     string    `gorm:"column:title" json:"title"`
	Completed bool      `gorm:"column:completed" json:"completed"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Todo) TableName() string {
	return "todos"
}

func (Todo) GetList(userId uint, pagination helpers.Pagination) (*helpers.Pagination, error) {
	db := database.Connect()
	var todoList []Todo

	query := db.Scopes(helpers.Paginate(todoList, &pagination, db)).Preload("User").Where(&Todo{UserID: userId}).Find(&todoList)
	pagination.Rows = todoList

	return &pagination, query.Error
}

func (Todo) Store(payload TodoForm) (Todo, error) {
	db := database.Connect()
	todo := Todo{
		Title:     payload.Title,
		Completed: payload.Completed,
		UserID:    payload.UserID,
	}

	query := db.Preload("User").Create(&todo)

	return todo, query.Error
}

func (Todo) Show(id int) (Todo, int, error) {
	db := database.Connect()

	var todo Todo

	query := db.Preload("User").First(&todo, id)

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
	db := database.Connect()

	query := db.Save(&todo)

	return query.Error
}

func (Todo) Delete(todo Todo) error {
	db := database.Connect()

	query := db.Delete(&todo)

	return query.Error
}

func (Todo) IsCorrectUser(userId uint, todo Todo) bool {
	return userId == todo.UserID
}
