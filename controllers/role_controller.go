package controllers

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoleController struct {
}

// @description Get data list
// @param 		echo.Context
// @return		error
func (RoleController) Index(c echo.Context) error {
	role := models.Role{}
	pagination := helpers.SetPagination(c, helpers.Pagination{})
	roleList, err := role.Index(pagination)

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get role success", roleList)
}

// @description Store data
// @param 		echo.Context
// @return		error
func (RoleController) Store(c echo.Context) error {
	form := new(models.RoleCreateForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	role := models.Role{}
	newRole, err := role.Store(models.RoleCreateForm{
		Name: form.Name,
	})

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusCreated, "Role created successfully", newRole)
}

// @description Get single data
// @param 		echo.Context
// @return		error
func (RoleController) Show(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	role := models.Role{}
	roleDetail, statusCode, err := role.Show(uint(id))
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get role detail success", roleDetail)
}

// @description Update data
// @param 		echo.Context
// @return		error
func (RoleController) Update(c echo.Context) error {
	form := new(models.RoleUpdateForm)
	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	id, _ := strconv.Atoi(c.Param("id"))

	role := models.Role{}
	roleDetail, statusCode, err := role.Show(uint(id))
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	roleDetail.Name = form.Name

	updatedRole, err := role.Update(roleDetail)
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Update role success", updatedRole)
}

// @description Update permission
// @param 		echo.Context
// @return		error
func (RoleController) UpdatePermission(c echo.Context) error {
	form := new(models.RoleUpdatePermissionForm)
	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	roleId, _ := strconv.Atoi(c.Param("roleId"))
	role := models.Role{}

	// Get role detail
	roleDetail, statusCode, err := role.Show(uint(roleId))
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	// Update permissions
	var mapRolePermissions []models.RolePermission
	for _, permission := range form.Permissions {
		mapRolePermissions = append(mapRolePermissions, models.RolePermission{
			ID:     permission.ID,
			RoleID: uint(roleId),
			Menu:   permission.Menu,
			Action: []byte(permission.Action),
		})
	}
	roleDetail.RolePermission = mapRolePermissions
	updatedRole, err := role.UpdatePermission(roleDetail)
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Permissions updated successfully", updatedRole)
}

// @description Delete data
// @param 		echo.Context
// @return		error
func (RoleController) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	role := models.Role{}
	deletedRole, statusCode, err := role.Delete(uint(id))
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Role deleted successfully", deletedRole)
}

// @description Get permissions list
// @param 		echo.Context
// @return		error
func (RoleController) PermissionList(c echo.Context) error {
	permission := models.Permission{}

	return helpers.Ok(http.StatusOK, "Get permission list success", permission.GeneratePermissions())
}
