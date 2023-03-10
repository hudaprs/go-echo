package structs

import "go-echo/models"

type RoleCreateEditForm struct {
	Name string `json:"name" validate:"required"`
}

type RolePermissionForm struct {
	ID      uint           `json:"id" validate:"required,numeric"`
	Actions models.Actions `json:"actions"`
}

type RoleAssignPermissionsForm struct {
	Permissions []RolePermissionForm `json:"permissions" validate:"required,dive,required"`
}

type RoleAssignUsersForm struct {
	Roles []uint `json:"roles" validate:"required"`
}
