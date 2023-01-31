package services

import (
	"echo-rest/models"
	"echo-rest/structs"

	"gorm.io/gorm"
)

type PermissionService struct {
	DB *gorm.DB
}

func (ps *PermissionService) Index() ([]models.PermissionResponse, error) {
	var permissions []models.PermissionResponse

	query := ps.DB.Find(&permissions)

	return permissions, query.Error
}

func (rs *PermissionService) AssignPermissions(payload structs.RoleAssignPermissionsForm) error {
	var rolePermissionList []models.RolePermission

	for _, rolePermissionPayload := range payload.Permissions {
		rolePermissionList = append(rolePermissionList, models.RolePermission{
			RoleID:       rolePermissionPayload.RoleID,
			PermissionID: rolePermissionPayload.PermissionID,
		})
	}

	query := rs.DB.Create(&rolePermissionList)

	return query.Error
}
