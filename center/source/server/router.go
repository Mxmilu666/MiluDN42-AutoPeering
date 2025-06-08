package server

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/server/handles"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	{
		api.GET("/nodes/info", handles.GetAllNodesInfo)
		api.GET("/node/info/:name", handles.GetNodeInfo)

		api.GET("/asn/verify", handles.SendASNVerifyCode)
		api.POST("/asn/verify", handles.VerifyASNCodeAndIssueJWT)

		peer := api.Group("/peer", handles.JWTMiddleware("asn-verify"))
		{
			peer.GET("/info", handles.GetPeerInfo)
			peer.POST("/create", handles.CreatePeerHandler)
		}
	}
	return r
}
