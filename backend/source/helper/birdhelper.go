package helper

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

// 获取所有 BGP peer 名字
func GetBGPPeerNames() ([]string, error) {
	cmd := exec.Command("birdc", "s", "p")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	var peers []string
	re := regexp.MustCompile(`^([\w\-\.]+)\s+BGP\s+`)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			peers = append(peers, matches[1])
		}
	}
	return peers, nil
}

// 获取并解析指定 BGP peer 的详细信息
func GetAndParseBGPPeerDetail(peerName string) (*BirdPeerDetail, error) {
	cmd := exec.Command("birdc", "s", "p", "a", peerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return ParseBirdPeerDetail(peerName, string(output))
}

// BGP Peer详细信息结构体
type BirdPeerChannelInfo struct {
	ChannelType string // "ipv4" 或 "ipv6"
	State       string
	Imported    int
	Exported    int
	Preferred   int
}

type BirdPeerDetail struct {
	PeerName string
	BGPState string
	Channels []BirdPeerChannelInfo
}

// 解析 birdc s p a <peer_name> 输出，提取关键信息
func ParseBirdPeerDetail(peerName, output string) (*BirdPeerDetail, error) {
	lines := strings.Split(output, "\n")
	var detail BirdPeerDetail
	detail.PeerName = peerName
	var currentChannel *BirdPeerChannelInfo
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "BGP state:") {
			detail.BGPState = strings.TrimSpace(strings.TrimPrefix(line, "BGP state:"))
		}
		if strings.HasPrefix(line, "Channel ") {
			// 遇到新channel前，先保存上一个channel
			if currentChannel != nil {
				detail.Channels = append(detail.Channels, *currentChannel)
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentChannel = &BirdPeerChannelInfo{ChannelType: parts[1]}
			} else {
				currentChannel = nil
			}
			continue
		}
		if currentChannel != nil && strings.HasPrefix(line, "State:") {
			currentChannel.State = strings.TrimSpace(strings.TrimPrefix(line, "State:"))
		}
		if currentChannel != nil && strings.HasPrefix(line, "Routes:") {
			re := regexp.MustCompile(`(\d+) imported, (\d+) exported, (\d+) preferred`)
			matches := re.FindStringSubmatch(line)
			if len(matches) == 4 {
				currentChannel.Imported = parseInt(matches[1])
				currentChannel.Exported = parseInt(matches[2])
				currentChannel.Preferred = parseInt(matches[3])
			}
		}
	}
	// 循环结束后，别忘了保存最后一个channel
	if currentChannel != nil {
		detail.Channels = append(detail.Channels, *currentChannel)
	}
	return &detail, nil
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// 根据关键字模糊查找最匹配的peer名字
func FindPeerNameByFuzzy(peers []string, keyword string) (string, bool) {
	keyword = strings.ToLower(keyword)
	for _, peer := range peers {
		if strings.Contains(strings.ToLower(peer), keyword) {
			return peer, true
		}
	}
	return "", false
}
