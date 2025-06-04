package helper

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source/models"
)

type PeerType = models.PeerType

// peerType 到模板文件名的映射
var birdTemplateMap = map[PeerType]string{
	models.PeerTypeIPv4Only:                                "bird_ipv4_only.conf",
	models.PeerTypeIPv6Only:                                "bird_ipv6_only.conf",
	models.PeerTypeDualStack:                               "bird_ipv4_ipv6.conf",
	models.PeerTypeMultiProtocol:                           "bird_ipv6_multi_protocol.conf",
	models.PeerTypeIPv6OnlyLocalLinkv6:                     "bird_ipv6_only.conf",
	models.PeerTypeDualStackLocalLinkv6:                    "bird_ipv4_ipv6.conf",
	models.PeerTypeMultiProtocolLocalLinkv6:                "bird_ipv6_multi_protocol.conf",
	models.PeerTypeMultiProtocolExtendedNextHop:            "bird_ipv6_multi_protocol_extended_next_hop.conf",
	models.PeerTypeMultiProtocolExtendedNextHopLocalLinkv6: "bird_ipv6_multi_protocol_extended_next_hop.conf",
}

// 构造 Bird 模板数据
type PeerRequestForTemplate struct {
	ASN             string
	IPv4            string
	IPv6            string
	PublicIP        string
	ExtendedNextHop bool
}

// RenderBirdConf 根据 PeerType 渲染 bird 配置模板
type BirdTemplateData map[string]any

func RenderBirdConf(peerType PeerType, data BirdTemplateData) (string, error) {
	tplName, ok := birdTemplateMap[peerType]
	if !ok {
		return "", fmt.Errorf("unknown peer type: %s", peerType)
	}
	tplPath := filepath.Join("template", tplName)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return buf.String(), nil
}

func BuildBirdTemplateData(req PeerRequestForTemplate) BirdTemplateData {
	peering_dn42_v6 := req.IPv6
	if len(req.IPv6) >= 4 && req.IPv6[:4] == "fe80" {
		peering_dn42_v6 = fmt.Sprintf("ip %% 'dn42_%s'", req.ASN)
	}
	return BirdTemplateData{
		"peering_asn":       req.ASN,
		"peering_dn42_v4":   req.IPv4,
		"peering_dn42_v6":   peering_dn42_v6,
		"public_ip":         req.PublicIP,
		"extended_next_hop": req.ExtendedNextHop,
	}
}
