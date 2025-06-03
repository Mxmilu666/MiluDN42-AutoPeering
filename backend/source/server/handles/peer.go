package handles

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PeerRequest struct {
	ASN             string `json:"asn"`
	IPv4            string `json:"ipv4"`
	IPv6            string `json:"ipv6"`
	PublicIP        string `json:"public_ip"`
	ExtendedNextHop bool   `json:"extended_next_hop"`
	Routes          string `json:"routes"` // "ipv4", "ipv6", "both"
	MultiProtocol   bool   `json:"multi_protocol"`
}

type PeerType string

const (
	PeerTypeIPv4Only                                PeerType = "ipv4_only"
	PeerTypeIPv6Only                                PeerType = "ipv6_only"
	PeerTypeDualStack                               PeerType = "ipv4_ipv6"
	PeerTypeMultiProtocol                           PeerType = "ipv6_multi_protocol"
	PeerTypeIPv6OnlyLocalLinkv6                     PeerType = "ipv6_only_local_linkv6"
	PeerTypeDualStackLocalLinkv6                    PeerType = "ipv4_ipv6_local_linkv6"
	PeerTypeMultiProtocolLocalLinkv6                PeerType = "ipv6_multi_protocol_local_linkv6"
	PeerTypeMultiProtocolExtendedNextHop            PeerType = "ipv6_multi_protocol_extended_next_hop"
	PeerTypeMultiProtocolExtendedNextHopLocalLinkv6 PeerType = "ipv6_multi_protocol_extended_next_hop_local_linkv6"
)

func isLocalLink(ipv6 string) bool {
	return len(ipv6) >= 4 && ipv6[:4] == "fe80"
}

func JudgePeerType(req PeerRequest) PeerType {
	isLocal := isLocalLink(req.IPv6)

	switch {
	case req.MultiProtocol && req.ExtendedNextHop && isLocal:
		return PeerTypeMultiProtocolExtendedNextHopLocalLinkv6
	case req.MultiProtocol && req.ExtendedNextHop:
		return PeerTypeMultiProtocolExtendedNextHop
	case req.MultiProtocol && isLocal:
		return PeerTypeMultiProtocolLocalLinkv6
	case req.MultiProtocol:
		return PeerTypeMultiProtocol
	case isLocal && req.Routes == "both":
		return PeerTypeDualStackLocalLinkv6
	case isLocal && req.Routes == "ipv6":
		return PeerTypeIPv6OnlyLocalLinkv6
	case req.Routes == "both":
		return PeerTypeDualStack
	case req.Routes == "ipv4":
		return PeerTypeIPv4Only
	case req.Routes == "ipv6":
		return PeerTypeIPv6Only
	default:
		return "unknown"
	}
}

func PeerHandler(c *gin.Context) {
	var req PeerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	typeStr := JudgePeerType(req)
	c.JSON(http.StatusOK, gin.H{"type": typeStr})
}
