package models

import (
	"echo-rest/database"
	"echo-rest/helpers"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"name" json:"name"`
	Email     string    `gorm:"email;index:index_email,unique" json:"email"`
	Password  string    `gorm:"password" json:"-"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

type UserStoreForm struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type UserLoginForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserAuthResponse struct {
	Token string `json:"token"`
}

func (User) TableName() string {
	return "users"
}

func (User) CheckEmail(email string) (User, int, error) {
	var user User
	db := database.DatabaseConnection()

	query := db.Where("email = ?", email).First(&user)

	findUserStatusCode := helpers.ValidateNotFoundData(query.Error)

	return user, findUserStatusCode, query.Error
}

func (User) Store(payload UserStoreForm) (User, error) {
	db := database.DatabaseConnection()

	hashedPassword, err := helpers.PasswordHash(payload.Password)

	if err != nil {
		panic("User Store: failed when start to hash password")
	}

	user := User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	query := db.Create(&user)

	return user, query.Error
}
