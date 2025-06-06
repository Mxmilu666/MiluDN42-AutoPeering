package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/helper"
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
	// 模糊查找所有匹配的 peer 名字
	matchedPeers := helper.FindPeerNamesByFuzzy(peers, asn)
	var mergedDetail *helper.BirdPeerDetail = nil
	if len(matchedPeers) > 0 {
		mergedDetail = &helper.BirdPeerDetail{
			PeerName: asn,
			Channels: []helper.BirdPeerChannelInfo{},
		}
		for idx, peerName := range matchedPeers {
			detail, err := helper.GetAndParseBGPPeerDetail(peerName)
			if err != nil {
				SendResponse(c, http.StatusInternalServerError, "error", "failed to get BGP peer detail")
				return
			}
			if idx == 0 {
				mergedDetail.BGPState = detail.BGPState
			}
			mergedDetail.Channels = append(mergedDetail.Channels, detail.Channels...)
		}
	}

	info := map[string]any{
		"asn":     asn,
		"wginfo":  wginfo,
		"bgpinfo": mergedDetail,
	}
	SendResponse(c, http.StatusOK, "success", info)
}
