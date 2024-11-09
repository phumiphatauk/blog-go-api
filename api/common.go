package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) getVersion(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, successResponse("%%CI_COMMIT_SHORT_SHA%%"))
}
