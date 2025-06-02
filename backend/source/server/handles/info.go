package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source"
	"github.com/gin-gonic/gin"
)

type NodeInfo struct {
	IPv4 string `json:"ipv4,omitempty"`
	IPv6 string `json:"ipv6,omitempty"`
	ASN  string `json:"asn"`
}

// GetInfo 获取节点的基本信息
func GetInfo(c *gin.Context) {
	info := NodeInfo{
		IPv4: source.AppConfig.DN42.IPv4,
		IPv6: source.AppConfig.DN42.IPv6,
		ASN:  source.AppConfig.DN42.ASN,
	}

	SendResponse(c, http.StatusOK, "success", info)
}
