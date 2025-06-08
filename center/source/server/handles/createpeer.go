package handles

import (
	"net"
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/api"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/helper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomCreatePeerRequest struct {
	IPv4            string `json:"ipv4"`
	IPv6            string `json:"ipv6"`
	PublicIP        string `json:"public_ip"`
	ExtendedNextHop bool   `json:"extended_next_hop"`
	Routes          string `json:"routes"`
	MultiProtocol   bool   `json:"multi_protocol"`
	Wireguard       struct {
		Pubkey   string `json:"pubkey"`
		Endpoint string `json:"endpoint"`
	} `json:"wireguard"`
}

// CreatePeerHandler 处理创建 peer 的 API
func CreatePeerHandler(c *gin.Context) {
	var reqBody CustomCreatePeerRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		SendResponse(c, http.StatusBadRequest, "invalid request", nil)
		return
	}

	node := c.Query("node")
	if node == "" {
		SendResponse(c, http.StatusBadRequest, "node is required", nil)
		return
	}

	claims, exists := c.Get("jwtClaims")
	if !exists {
		SendResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	asn := claims.(jwt.MapClaims)["data"].(map[string]interface{})["asn"].(string)

	if asn == "" {
		SendResponse(c, http.StatusUnauthorized, "asn not found in token", nil)
		return
	}

	// 判断 Endpoint 是否为中国IP
	endpointHost := reqBody.Wireguard.Endpoint
	// 去除端口
	if host, _, err := net.SplitHostPort(endpointHost); err == nil {
		endpointHost = host
	}
	// 如果是域名，解析为IP
	ips, err := net.LookupIP(endpointHost)
	if err != nil || len(ips) == 0 {
		SendResponse(c, http.StatusBadRequest, "Invalid endpoint: cannot resolve host", nil)
		return
	}
	for _, ip := range ips {
		isCN, geoErr := helper.IsIPCNAny(ip.String())
		if geoErr != nil {
			SendResponse(c, http.StatusInternalServerError, "geoip check failed: "+geoErr.Error(), nil)
			return
		}
		if isCN {
			SendResponse(c, http.StatusBadRequest, "Endpoint IP is in China Mainland, not allowed", nil)
			return
		}
	}

	// 构造 api.CreatePeerRequest
	apiReq := api.CreatePeerRequest{
		ASN:             asn,
		IPv4:            reqBody.IPv4,
		IPv6:            reqBody.IPv6,
		PublicIP:        reqBody.PublicIP,
		ExtendedNextHop: reqBody.ExtendedNextHop,
		Routes:          reqBody.Routes,
		MultiProtocol:   reqBody.MultiProtocol,
		Wireguard: struct {
			Pubkey   string `json:"pubkey"`
			Endpoint string `json:"endpoint"`
		}{
			Pubkey:   reqBody.Wireguard.Pubkey,
			Endpoint: reqBody.Wireguard.Endpoint,
		},
	}

	resp, err := api.CreatePeer(node, &apiReq)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	SendResponse(c, http.StatusOK, resp.Msg, resp.Data)
}
