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

		api.POST("/verify", handles.RequestVerify)
		api.POST("/verify/confirm", handles.ConfirmVerify)
		api.POST("/peer", handles.PeerHandler)
	}

	// 动态生成的/verify/:dir路由
	r.GET("/verify/:dir", handles.VerifyHandler)

	return r
}
