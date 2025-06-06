package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/helper"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/models"
	"github.com/gin-gonic/gin"
)

type WireGuardInfo struct {
	PubKey   string `json:"pubkey"`
	Endpoint string `json:"endpoint"`
}

type PeerRequest struct {
	ASN             string        `json:"asn"`
	IPv4            string        `json:"ipv4"`
	IPv6            string        `json:"ipv6"`
	PublicIP        string        `json:"public_ip"`
	ExtendedNextHop bool          `json:"extended_next_hop"`
	Routes          string        `json:"routes"` // "ipv4", "ipv6", "both"
	MultiProtocol   bool          `json:"multi_protocol"`
	WireGuard       WireGuardInfo `json:"wireguard"`
}

type PeerType = models.PeerType

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
		if req.IPv4 != "" {
			return models.PeerTypeIPv4Only
		}
		return "unknown"
	}
	if req.Routes == "ipv6" {
		if req.IPv6 != "" {
			if isLocal {
				return models.PeerTypeIPv6OnlyLocalLinkv6
			}
			return models.PeerTypeIPv6Only
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
				return models.PeerTypeMultiProtocolExtendedNextHopLocalLinkv6
			}
			return models.PeerTypeMultiProtocolExtendedNextHop
		}
		if isLocal {
			return models.PeerTypeMultiProtocolLocalLinkv6
		}
		return models.PeerTypeMultiProtocol
	}

	// 非多协议类型
	if isLocal {
		return models.PeerTypeDualStackLocalLinkv6
	}
	return models.PeerTypeDualStack
}

func PeerHandler(c *gin.Context) {
	var req PeerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, http.StatusBadRequest, "error", "invalid request")
		return
	}

	typeStr := JudgePeerType(req)
	if typeStr == "unknown" {
		SendResponse(c, http.StatusBadRequest, "error", "invalid peer type")
		return
	}

	// 构造 Bird 模板数据
	birdData := helper.PeerRequestForTemplate{
		ASN:             req.ASN,
		IPv4:            req.IPv4,
		IPv6:            req.IPv6,
		PublicIP:        req.PublicIP,
		ExtendedNextHop: req.ExtendedNextHop,
	}
	data := helper.BuildBirdTemplateData(birdData)
	conf, err := helper.RenderBirdConf(helper.PeerType(typeStr), data)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, "error", "failed to render bird config")
		return
	}

	// WireGuard 配置校验
	if req.WireGuard.PubKey == "" || req.WireGuard.Endpoint == "" {
		SendResponse(c, http.StatusBadRequest, "error", "WireGuard pubkey and endpoint are required")
	}

	// 构造 WireGuard 模板数据
	wgData := helper.PeerRequestWGForTemplate{
		ASN:  req.ASN,
		IPv4: req.IPv4,
		WireGuard: helper.WireGuardInfoForTemplate{
			PubKey:   req.WireGuard.PubKey,
			Endpoint: req.WireGuard.Endpoint,
		},
	}

	wgConf := helper.BuildWireGuardTemplateData(wgData)
	wgConfStr, wgErr := helper.RenderWireGuardConf(wgConf)
	if wgErr != nil {
		SendResponse(c, http.StatusInternalServerError, "error", "failed to render wireguard config")
		return
	}

	err = helper.SetupConfFiles(req.ASN, conf, wgConfStr)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, "error", "failed to setup config files")
		return
	}

	// 新增：调用wg-quick up和birdc c
	err = helper.RunWgQuickAndBirdc(req.ASN)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, "error", "failed to setup")
		return
	}

	SendResponse(c, http.StatusOK, "success", "peer setup successfully type: "+typeStr)
}
