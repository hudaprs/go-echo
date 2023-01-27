package models

import (
	"gorm.io/datatypes"
)

type Action struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type Permission struct {
	ID          string                     `json:"id"`
	Description string                     `json:"description"`
	Action      datatypes.JSONType[Action] `json:"action"`
}

func (Permission) GeneratePermissions() datatypes.JSONType[[]Permission] {
	permission := []Permission{
		{ID: "Todo", Description: "Menu permissions for todo feature", Action: datatypes.JSONType[Action]{
			Data: Action{
				Create: false,
				Read:   false,
				Update: false,
				Delete: false,
			},
		},
		},
	}

	return datatypes.JSONType[[]Permission]{
		Data: permission,
	}
}
