package structs

type TodoForm struct {
	Title     string `json:"title" validate:"required"`
	Completed bool   `json:"completed"`
	UserID    uint   `json:"userId"`
}
