package api

import (
	"net/http"

	db "blog-go-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

type GetAllRoleRequest struct {
	Name     string `form:"name"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=10"`
}

// GetAllRole godoc
//
//	@Summary		Get All Role
//	@Description	Get All Role
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			name		query		string	false	"Role Name"
//	@Param			page_id		query		int		true	"Page ID"
//	@Param			page_size	query		int		true	"Page Size"
//	@Success		200			{object}	jsonResponseWithPaginate
//	@Router			/api/role [get]
//	@Security		BearerAuth
func (server *Server) GetAllRole(ctx *gin.Context) {

	var req GetAllRoleRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAllRoleParams{
		Lower:  req.Name,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := server.store.GetAllRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	count, err := server.store.CountAllRole(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := jsonResponseWithPaginate{
		jsonResponse: jsonResponse{
			Error:   false,
			Message: "successfully",
			Data:    roles,
		},
		Total: count,
	}

	ctx.JSON(http.StatusOK, payload)
}

type GetRoleByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetRoleById godoc
//
//	@Summary		Get Role By ID
//	@Description	Get Role By ID
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Role ID"
//	@Success		200	{object}	jsonResponse
//	@Router			/api/role/{id} [get]
//	@Security		BearerAuth
func (server *Server) GetRoleById(ctx *gin.Context) {

	var req GetRoleByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	role, err := server.store.GetRoleById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    role,
	})
}

type CreateRoleRequest struct {
	Name             string                      `json:"name" binding:"required"`
	PermissionGroups []getAllPermissionGroupData `json:"permission_groups"`
}

// CreateRole godoc
//
//	@Summary		Create Role
//	@Description	Create Role
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			input	body		CreateRoleRequest	true	"Role Information"
//	@Success		200		{object}	jsonResponse
//	@Router			/api/role [post]
//	@Security		BearerAuth
func (server *Server) CreateRole(ctx *gin.Context) {

	var req CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	role, err := server.store.CreateRole(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert Role Permission
	for _, pg := range req.PermissionGroups {
		for _, p := range pg.Permissions {

			// Check exists role_permission
			// if is_assigned == true, insert
			if p.IsAssigned {
				argCreateRolePermission := db.CreateRolePermissionParams{
					RoleID:       role.ID,
					PermissionID: p.ID,
				}

				_, err := server.store.CreateRolePermission(ctx, argCreateRolePermission)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    role,
	})
}

type UpdateRoleRequest struct {
	ID               int64                       `json:"id" binding:"required,min=1"`
	Name             string                      `json:"name" binding:"required"`
	PermissionGroups []getAllPermissionGroupData `json:"permission_groups"`
}

// UpdateRole godoc
//
//	@Summary		Update Role
//	@Description	Update Role
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			input	body		UpdateRoleRequest	true	"Role Information"
//	@Success		200		{object}	jsonResponse
//	@Router			/api/role [put]
//	@Security		BearerAuth
func (server *Server) UpdateRole(ctx *gin.Context) {

	var req UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateRoleParams{
		ID:   req.ID,
		Name: req.Name,
	}

	// Update Role
	err := server.store.UpdateRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Insert and Delete Role Permission
	for _, pg := range req.PermissionGroups {
		for _, p := range pg.Permissions {

			argGetRolePermissionByRoleIdAndPermissionId := db.GetRolePermissionByRoleIdAndPermissionIdParams{
				RoleID:       req.ID,
				PermissionID: p.ID,
			}

			// Check exists role_permission
			// if exists and is_assigned == false, delete
			// if not exists and is_assigned == true, insert
			rolePermissionExists, _ := server.store.GetRolePermissionByRoleIdAndPermissionId(ctx, argGetRolePermissionByRoleIdAndPermissionId)

			if rolePermissionExists.ID == 0 && p.IsAssigned {
				argCreateRolePermission := db.CreateRolePermissionParams{
					RoleID:       req.ID,
					PermissionID: p.ID,
				}

				_, err := server.store.CreateRolePermission(ctx, argCreateRolePermission)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
			} else if rolePermissionExists.ID != 0 && !p.IsAssigned {
				argDeleteRolePermission := db.DeleteRolePermissionParams{
					RoleID:       req.ID,
					PermissionID: p.ID,
				}

				err := server.store.DeleteRolePermission(ctx, argDeleteRolePermission)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    nil,
	})
}

type DeleteRoleRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// DeleteRole godoc
//
//	@Summary		Delete Role
//	@Description	Delete Role
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Role ID"
//	@Success		200	{object}	jsonResponse
//	@Router			/api/role/{id} [delete]
//	@Security		BearerAuth
func (server *Server) DeleteRole(ctx *gin.Context) {

	var req DeleteRoleRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Delete role by id
	err := server.store.DeleteRole(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete all role_permission by role_id
	err = server.store.DeleteRolePermissionByRoleId(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    nil,
	})
}

// GetRoleForDropDownList godoc
//
//	@Summary		Get Role For Drop Down List
//	@Description	Get Role For Drop Down List
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	jsonResponse
//	@Router			/api/role/dropdownlist [get]
//	@Security		BearerAuth
func (Server *Server) GetRoleForDropDownList(ctx *gin.Context) {

	roles, err := Server.store.GetRoleForDropDownList(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    roles,
	})
}
