package api

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/helper"
)

// CreatePeerRequest 定义 peer 创建请求体
type CreatePeerRequest struct {
	ASN             string `json:"asn"`
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

// CreatePeerResponse 定义 peer 创建响应体
type CreatePeerResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   string `json:"data"`
}

// CreatePeer 向指定 node 发送创建 peer 请求
func CreatePeer(nodeName string, req *CreatePeerRequest) (*CreatePeerResponse, error) {
	resp, err := helper.RequestNodeByName[CreatePeerResponse](nodeName, "/api/peer", "POST", req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
