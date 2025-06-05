package server

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source/server/handles"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) *gin.Engine {
	// API路由
	api := r.Group("/api", AuthMiddleware())
	{
		node := api.Group("/node")
		{
			node.GET("/info", handles.GetInfo)
		}
		api.POST("/peer", handles.PeerHandler)
		api.GET("/peer/info", handles.GetPeerInfoHandler)
	}
	return r
}
