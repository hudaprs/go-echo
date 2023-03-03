package middlewares

import (
	"go-echo/database"
	"go-echo/helpers"
	"go-echo/locales"
	"go-echo/models"

	"github.com/labstack/echo/v4"
)

type PermissionAction string

const (
	Create PermissionAction = "Create"
	Read   PermissionAction = "Read"
	Update PermissionAction = "Update"
	Delete PermissionAction = "Delete"
)

func RoleCheck(permissionCode string, permissionAction PermissionAction) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var permissionDetail *models.Permission
			var userActiveRole *models.RoleUser
			var rolePermission *models.RolePermission
			authenticatedUser := helpers.JwtGetClaims(c)

			conn := database.Connect()

			// Find permission
			if query := conn.Where(&models.Permission{Code: permissionCode}).First(&permissionDetail); query.Error != nil {
				statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
					NotFound:           locales.LocalesGet("validation.permissionAccess") + ": " + locales.LocalesGet("validation.notFound"),
					NeedAuthentication: true,
				})
				return helpers.ErrorDynamic(statusCode, err.Error())
			}

			// Find active role
			if query := conn.Where(&models.RoleUser{UserID: authenticatedUser.ID, IsActive: true}).First(&userActiveRole); query.Error != nil {
				statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
					NotFound:           locales.LocalesGet("validation.permissionAccess") + ": " + locales.LocalesGet("roleUser.validation.notFound"),
					NeedAuthentication: true,
				})
				return helpers.ErrorDynamic(statusCode, err.Error())
			}

			// Find permission of role
			if userActiveRole != nil && permissionDetail != nil {
				if query := conn.Where(&models.RolePermission{RoleID: userActiveRole.RoleID, PermissionID: permissionDetail.ID}).First(&rolePermission); query.Error != nil {
					statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
						NotFound:           locales.LocalesGet("validation.permissionAccess") + ": " + locales.LocalesGet("rolePermission.validation.notFound"),
						NeedAuthentication: true,
					})
					return helpers.ErrorDynamic(statusCode, err.Error())
				}
			}

			// Check role permission
			// Check the permission action
			// If match, pass the user in
			// If not, throw forbidden error response
			if rolePermission != nil {
				// Create Action
				if permissionAction == Create && rolePermission.Actions.Data.Create {
					return next(c)
				}

				// Read Action
				if permissionAction == Read && rolePermission.Actions.Data.Read {
					return next(c)
				}

				// Update Action
				if permissionAction == Update && rolePermission.Actions.Data.Update {
					return next(c)
				}

				// Delete Action
				if permissionAction == Delete && rolePermission.Actions.Data.Delete {
					return next(c)
				}
			}

			return helpers.ErrorForbidden()
		}
	}
}
