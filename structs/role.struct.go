package structs

import "echo-rest/models"

type RoleCreateEditForm struct {
	Name string `json:"name" validate:"required"`
}

type RolePermissionForm struct {
	Code   string        `json:"code" validate:"required"`
	Action models.Action `json:"action"`
}

type RoleAssignPermissionsForm struct {
	Permissions []RolePermissionForm `json:"permissions" validate:"required,dive,required"`
}

type RoleAssignUsersForm struct {
	Roles []uint `json:"roles" validate:"required"`
}
