package api

import (
	"encoding/json"
	"fmt"

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
	respBytes, err := helper.RequestNodeByName(nodeName, "/api/node/info")
	if err != nil {
		return nil, err
	}
	var resp NodeInfoResponse
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("Node returned error: %s", resp.Msg)
	}
	return &resp.Data, nil
}
