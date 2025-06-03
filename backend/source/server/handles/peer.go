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

	// 规则校验
	if req.MultiProtocol && req.Routes != "both" {
		return "unknown"
	}
	if req.ExtendedNextHop && (!req.MultiProtocol || req.IPv6 == "") {
		return "unknown"
	}

	// 只允许单IP类型
	if req.Routes == "ipv4" {
		if req.IPv4 != "" && req.IPv6 == "" {
			return PeerTypeIPv4Only
		}
		return "unknown"
	}
	if req.Routes == "ipv6" {
		if req.IPv6 != "" && req.IPv4 == "" {
			if isLocal {
				return PeerTypeIPv6OnlyLocalLinkv6
			}
			return PeerTypeIPv6Only
		}
		return "unknown"
	}

	// 需要双IP类型
	if req.IPv4 == "" || req.IPv6 == "" {
		return "unknown"
	}

	// 多协议类型
	if req.MultiProtocol {
		if req.ExtendedNextHop {
			if isLocal {
				return PeerTypeMultiProtocolExtendedNextHopLocalLinkv6
			}
			return PeerTypeMultiProtocolExtendedNextHop
		}
		if isLocal {
			return PeerTypeMultiProtocolLocalLinkv6
		}
		return PeerTypeMultiProtocol
	}

	// 非多协议类型
	if isLocal {
		return PeerTypeDualStackLocalLinkv6
	}
	return PeerTypeDualStack
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
