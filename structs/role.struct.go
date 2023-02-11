package structs

import "echo-rest/models"

type RoleCreateEditForm struct {
	Name string `json:"name" validate:"required"`
}

type RolePermissionForm struct {
	ID     uint          `json:"id" validate:"required,numeric"`
	Action models.Action `json:"action"`
}

type RoleAssignPermissionsForm struct {
	Permissions []RolePermissionForm `json:"permissions" validate:"required,dive,required"`
}

type RoleAssignUsersForm struct {
	Roles []uint `json:"roles" validate:"required"`
}
