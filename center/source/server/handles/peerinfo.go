package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/api"
	"github.com/gin-gonic/gin"
)

// GetPeerInfo 获取指定 ASN 的 Peer 信息
func GetPeerInfo(c *gin.Context) {
	asn := c.Query("asn")
	node := c.Query("node") // 支持通过 node 参数指定节点
	if asn == "" {
		SendResponse(c, http.StatusBadRequest, "asn is required", nil)
		return
	}
	if node == "" {
		SendResponse(c, http.StatusBadRequest, "node is required", nil)
		return
	}
	info, err := api.GetPeerInfoByASN(node, asn)
	if err != nil {
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	SendResponse(c, http.StatusOK, "success", info)
}
