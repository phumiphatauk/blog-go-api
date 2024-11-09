package api

import (
	db "blog-go-api/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getAllPermissionGroupRequest struct {
	RoleId int64 `form:"role_id"`
}

type getAllPermissionGroupData struct {
	db.PermissionGroup
	Permissions []db.GetPermissionByPermissionGroupIdAndRoleIdRow `json:"permissions"`
}

func GetAllPermissionGroupResponse(permissionGroup db.PermissionGroup, permissions []db.GetPermissionByPermissionGroupIdAndRoleIdRow) getAllPermissionGroupData {

	return getAllPermissionGroupData{
		PermissionGroup: permissionGroup,
		Permissions:     permissions,
	}
}

func (server *Server) GetAllPermissionGroup(ctx *gin.Context) {
	// Get User from Token
	var req getAllPermissionGroupRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	permissionGroups, err := server.store.GetAllPermissionGroup(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response []getAllPermissionGroupData

	for _, permissionGroup := range permissionGroups {

		arg := db.GetPermissionByPermissionGroupIdAndRoleIdParams{
			PermissionGroupID: permissionGroup.ID,
			RoleID:            req.RoleId,
		}

		permissions, err := server.store.GetPermissionByPermissionGroupIdAndRoleId(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		permissionGroup := GetAllPermissionGroupResponse(permissionGroup, permissions)

		response = append(response, permissionGroup)
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Error:   false,
		Message: "successfully",
		Data:    response,
	})
}
