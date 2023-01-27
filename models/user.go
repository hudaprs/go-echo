package models

import (
	"echo-rest/database"
	"echo-rest/helpers"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Email     string    `gorm:"column:email;index:index_email,unique" json:"email"`
	Password  string    `gorm:"column:password" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
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

type UserRefreshForm struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type UserLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type UserMeResponse struct {
	ID    uint `json:"id"`
	Name  uint `json:"name"`
	Email uint `json:"email"`
}

func (User) TableName() string {
	return "users"
}

func (User) CheckEmail(email string) (User, int, error) {
	var user User
	db := database.Connect()

	query := db.Where("email = ?", email).First(&user)

	findUserStatusCode := helpers.ValidateNotFoundData(query.Error)

	return user, findUserStatusCode, query.Error
}

func (User) Show(id uint) (User, int, error) {
	db := database.Connect()

	var user User

	query := db.First(&user, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return user, statusCode, err
}

func (User) Store(payload UserStoreForm) (User, error) {
	db := database.Connect()

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
