package controllers

import (
	"go-echo/helpers"
	"go-echo/services"
	"go-echo/structs"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PermissionController struct {
	PermissionService services.PermissionService
}

// @description Get data list
// @param 		echo.Context
// @return		error
func (rc PermissionController) Index(c echo.Context) error {
	permissionList, err := rc.PermissionService.Index()

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get permission list success", permissionList)
}

// @description Assign permissions to role
// @param 		echo.Context
// @return		error
func (rc PermissionController) AssignPermissions(c echo.Context) error {
	form := new(structs.RoleAssignPermissionsForm)
	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	roleId, _ := strconv.Atoi(c.Param("roleId"))

	rolePermissionList, err := rc.PermissionService.AssignPermissions(uint(roleId), *form)
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Permission assigned successfully", rolePermissionList)
}
