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
	roleDetail.Permissions = form.Permissions

	updatedRole, err := role.Update(roleDetail)
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Update role success", updatedRole)
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
