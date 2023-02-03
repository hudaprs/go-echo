package structs

import "echo-rest/models"

type RoleCreateEditForm struct {
	Name string `json:"name" validate:"required"`
}

type RolePermissionForm struct {
	PermissionID uint          `json:"permissionId" validate:"required,numeric"`
	Action       models.Action `json:"action"`
}

type RoleAssignPermissionsForm struct {
	Permissions []RolePermissionForm `json:"permissions" validate:"required,dive,required"`
}
