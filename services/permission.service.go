package services

import (
	"echo-rest/models"
	"echo-rest/structs"

	"gorm.io/datatypes"
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

func (rs *PermissionService) AssignPermissions(roleId uint, payload structs.RoleAssignPermissionsForm) ([]models.RolePermissionResponse, error) {
	var existedPermissionList []models.RolePermissionResponse
	var mergedPermissionList []models.RolePermissionResponse

	// Check if theres any permissions existed before
	// Delete all data, and create an new one
	query := rs.DB.Where(models.RolePermissionResponse{RoleID: roleId}).Find(&existedPermissionList).Delete(&models.RolePermission{})
	if query.Error != nil {
		return []models.RolePermissionResponse{}, query.Error
	}

	// Create New Permissions
	for _, rolePermissionPayload := range payload.Permissions {
		var permissionDetail models.Permission
		if query := rs.DB.Find(&permissionDetail, models.Permission{ID: rolePermissionPayload.ID}); query.Error != nil {
			return []models.RolePermissionResponse{}, query.Error
		}

		mergedPermissionList = append(mergedPermissionList, models.RolePermissionResponse{
			RoleID:       roleId,
			PermissionID: rolePermissionPayload.ID,
			Actions: datatypes.JSONType[models.Actions]{
				Data: rolePermissionPayload.Actions,
			},

			// Append permission code to API response
			Code: permissionDetail.Code,
		})
	}
	query = rs.DB.Create(&mergedPermissionList)
	if query.Error != nil {
		return []models.RolePermissionResponse{}, query.Error
	}

	if len(mergedPermissionList) > 0 {
		return mergedPermissionList, query.Error
	} else {
		return []models.RolePermissionResponse{}, query.Error
	}
}
