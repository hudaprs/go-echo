package services

import (
	"go-echo/models"
	"go-echo/structs"

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
	rolePermissionResponse := []models.RolePermissionResponse{}
	var permissionIds []uint

	if len(payload.Permissions) > 0 {
		for _, permission := range payload.Permissions {
			permissionIds = append(permissionIds, permission.ID)
		}
	}

	// Check if role didn't include permissions
	if len(permissionIds) == 0 {
		if query := rs.DB.Where("role_id = ?", roleId).Delete(&models.RolePermissionResponse{}); query.Error != nil {
			return rolePermissionResponse, query.Error
		}
	}

	// Check if role already have permissions before
	// If role didn't sent the previous permissions, remove
	if query := rs.DB.Where("role_id = ?", roleId).Where("permission_id NOT IN ?", permissionIds).Delete(&models.RolePermissionResponse{}); query.Error != nil {
		return rolePermissionResponse, query.Error
	}

	// Check if role have existing permissions
	var existedRolePermissions []models.RolePermission
	if query := rs.DB.Where(models.RolePermission{RoleID: roleId}).Where("permission_id IN ?", permissionIds).Find(&existedRolePermissions); query.Error != nil {
		return rolePermissionResponse, query.Error
	}

	// Find unique role that never be assigned before to user
	var assignedPermissions []structs.RolePermissionForm
	for _, permission := range payload.Permissions {
		// Make identifier to skip if not exists
		skip := false
		for _, existedRolePermission := range existedRolePermissions {
			// If data found
			if permission.ID == existedRolePermission.PermissionID {
				// Just update the actions column
				if query := rs.DB.Model(&existedRolePermission).Update("actions", datatypes.JSONType[models.Actions]{
					Data: permission.Actions,
				}); query.Error != nil {
					return rolePermissionResponse, query.Error
				}

				skip = true
				break
			}
		}

		// If role not found, just make a new one
		if !skip {
			assignedPermissions = append(assignedPermissions, permission)
		}
	}

	// Check if theres any unique role that user didn't assign before
	// Create that unique permission and assign to specific roles
	if len(assignedPermissions) > 0 {
		for _, permission := range assignedPermissions {
			// Check if permission exists
			if query := rs.DB.First(&models.Permission{}, permission.ID); query.Error != nil {
				return rolePermissionResponse, query.Error
			}

			newUserPermission := &models.RolePermissionResponse{
				RoleID:       roleId,
				PermissionID: permission.ID,
				Actions: datatypes.JSONType[models.Actions]{
					Data: permission.Actions,
				},
			}

			if query := rs.DB.Create(newUserPermission); query.Error != nil {
				return rolePermissionResponse, query.Error
			}
		}
	}

	// Look up new assigned role through database
	if query := rs.DB.Select("role_permissions.*, permissions.code, permissions.id").Where("role_id = ?", roleId).Joins("left join permissions ON permissions.id = role_permissions.permission_id").Order("created_at desc").Find(&rolePermissionResponse); query.Error != nil {
		return rolePermissionResponse, query.Error
	}

	return rolePermissionResponse, nil
}
