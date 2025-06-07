package server

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/server/handles"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	{
		api.GET("/nodes/info", handles.GetAllNodesInfo)
	}
	return r
}
