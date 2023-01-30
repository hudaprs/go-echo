package models

type Action struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type Permission struct {
	ID          string `json:"id" validate:"required"`
	Description string `json:"description" validate:"required"`
	Action      Action `json:"action" validate:"required"`
}

func (Permission) GeneratePermissions() []Permission {
	permissions := []Permission{
		{ID: "Todo", Description: "Menu permissions for todo feature", Action: Action{
			Create: false,
			Read:   false,
			Update: false,
			Delete: false,
		},
		},
	}

	return permissions
}
