package api

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/helper"
)

type NodeInfoResponse struct {
	Status int      `json:"status"`
	Msg    string   `json:"msg"`
	Data   NodeInfo `json:"data"`
}

type NodeInfo struct {
	IPv4               string `json:"ipv4"`
	IPv6               string `json:"ipv6"`
	ASN                string `json:"asn"`
	WireguardPublicKey string `json:"wireguard_public_key"`
}

// GetNodeInfoByName 请求指定 node 的 info
func GetNodeInfoByName(nodeName string) (*NodeInfo, error) {
	resp, err := helper.RequestNodeByName[NodeInfoResponse](nodeName, "/api/node/info", "GET", nil)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
