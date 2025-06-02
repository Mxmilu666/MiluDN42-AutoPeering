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
	Routes          string `json:"routes"`
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
	mode := req.Routes
	isLocalLink := isLocalLink(req.IPv6)

	if req.ExtendedNextHop {
		if isLocalLink && mode == "both" {
			return PeerTypeMultiProtocolExtendedNextHopLocalLinkv6
		} else if mode == "both" {
			return PeerTypeMultiProtocolExtendedNextHop
		}
	}
	if isLocalLink {
		if mode == "ipv4" {
			return PeerTypeDualStackLocalLinkv6
		}
		if mode == "ipv6" {
			return PeerTypeIPv6OnlyLocalLinkv6
		}
		if mode == "both" {
			return PeerTypeMultiProtocolLocalLinkv6
		}
	}
	if mode == "both" && !isLocalLink && !req.ExtendedNextHop {
		return PeerTypeMultiProtocol
	}
	if mode == "ipv4" {
		return PeerTypeIPv4Only
	}
	if mode == "ipv6" {
		return PeerTypeIPv6Only
	}
	if mode == "both" {
		return PeerTypeDualStack
	}
	return "unknown"
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
