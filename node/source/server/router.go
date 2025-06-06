package server

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/server/handles"
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
		api.GET("/peer/restart", handles.RestartHandler)
		api.GET("/peer/remove", handles.RemoveHandler)
	}
	return r
}
