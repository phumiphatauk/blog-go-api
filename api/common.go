package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@Summary		Get API version
//	@Description	Get the current version of the API
//	@Tags			Version
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/api/version [get]
func (server *Server) getVersion(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, successResponse("%%CI_COMMIT_SHORT_SHA%%"))
}
