package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source/helper"
	"github.com/gin-gonic/gin"
)

// GetPeerInfoHandler 查询指定 ASN 的 WireGuard 隧道信息和 BGP Peer 信息
func GetPeerInfoHandler(c *gin.Context) {
	asn := c.Query("asn")
	if asn == "" {
		SendResponse(c, http.StatusBadRequest, "error", "asn is required")
		return
	}
	wginfo, err := helper.RawWGShow("dn42_" + asn)
	if err != nil {
		SendResponse(c, http.StatusNotFound, "error", "asn not found")
		return
	}

	// 获取所有 BGP peer 名字
	peers, err := helper.GetBGPPeerNames()
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, "error", "failed to get BGP peers")
		return
	}
	// 模糊查找对应 peer 名字
	peerName, found := helper.FindPeerNameByFuzzy(peers, asn)
	var output any = nil
	if found {
		// 查询详细 BGP 信息
		var err error
		output, err = helper.GetAndParseBGPPeerDetail(peerName)
		if err != nil {
			SendResponse(c, http.StatusInternalServerError, "error", "failed to get BGP peer detail")
			return
		}
	}

	info := map[string]any{
		"asn":     asn,
		"wginfo":  wginfo,
		"bgpinfo": output,
	}
	SendResponse(c, http.StatusOK, "success", info)
}
