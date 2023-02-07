package structs

type RefreshTokenForm struct {
	UserID       uint   `json:"userId"`
	RefreshToken string `json:"refreshToken"`
}
