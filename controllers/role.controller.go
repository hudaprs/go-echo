package controllers

import (
	"go-echo/helpers"
	"go-echo/services"
	"go-echo/structs"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoleController struct {
	RoleService services.RoleService
}

// @description Get data list
// @param 		echo.Context
// @return		error
func (rc RoleController) Index(c echo.Context) error {
	pagination := helpers.SetPagination(c, helpers.Pagination{})
	roleList, err := rc.RoleService.Index(pagination)

	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get role success", roleList)
}

// @description Store data
// @param 		echo.Context
// @return		error
func (rc RoleController) Store(c echo.Context) error {
	form := new(structs.RoleCreateEditForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	newRole, err := rc.RoleService.Store(structs.RoleCreateEditForm{
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
func (rc RoleController) Show(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	roleDetail, statusCode, err := rc.RoleService.Show(uint(id))
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Get role detail success", roleDetail)
}

// @description Update data
// @param 		echo.Context
// @return		error
func (rc RoleController) Update(c echo.Context) error {
	form := new(structs.RoleCreateEditForm)
	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}
	if err := c.Validate(form); err != nil {
		return err
	}

	id, _ := strconv.Atoi(c.Param("id"))

	roleDetail, statusCode, err := rc.RoleService.Show(uint(id))
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	roleDetail.Name = form.Name

	updatedRole, err := rc.RoleService.Update(roleDetail)
	if err != nil {
		return helpers.ErrorServer(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Update role success", updatedRole)
}

// @description Delete data
// @param 		echo.Context
// @return		error
func (rc RoleController) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	deletedRole, statusCode, err := rc.RoleService.Delete(uint(id))
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Role deleted successfully", deletedRole)
}
