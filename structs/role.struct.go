package structs

type RoleCreateEditForm struct {
	Name string `json:"name" validate:"required"`
}

type RolePermissionForm struct {
	PermissionID uint `json:"permissionId" validate:"required,numeric"`
	RoleID       uint `json:"roleId" validate:"required,numeric"`
}

type RoleAssignPermissionsForm struct {
	Permissions []RolePermissionForm `json:"permissions" validate:"required,dive,required"`
}
