package api

import (
	"fmt"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/helper"
)

type PeerInfoResponse struct {
	Status int      `json:"status"`
	Msg    string   `json:"msg"`
	Data   PeerInfo `json:"data"`
}

type PeerInfo struct {
	ASN     string  `json:"asn"`
	BGPInfo BGPInfo `json:"bgpinfo"`
	WGInfo  string  `json:"wginfo"`
}

type BGPInfo struct {
	PeerName string           `json:"PeerName"`
	BGPState string           `json:"BGPState"`
	Channels []BGPChannelInfo `json:"Channels"`
}

type BGPChannelInfo struct {
	ChannelType string `json:"ChannelType"`
	State       string `json:"State"`
	Imported    int    `json:"Imported"`
	Exported    int    `json:"Exported"`
	Preferred   int    `json:"Preferred"`
}

// GetPeerInfoByASN 请求指定 ASN 的 peer info
func GetPeerInfoByASN(nodeName, asn string) (*PeerInfo, error) {
	apiPath := fmt.Sprintf("/api/peer/info?asn=%s", asn)
	resp, err := helper.RequestNodeByName[PeerInfoResponse](nodeName, apiPath, "GET", nil)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
