package structs

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
