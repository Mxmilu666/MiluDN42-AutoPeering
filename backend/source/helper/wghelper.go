package helper

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source"
)

// 构造 WireGuard 模板数据
type WireGuardInfoForTemplate struct {
	PubKey   string
	Endpoint string
}

type PeerRequestWGForTemplate struct {
	ASN       string
	IPv4      string
	WireGuard WireGuardInfoForTemplate
}

// WireGuardTemplateData 用于渲染 WireGuard 配置模板的数据结构
func BuildWireGuardTemplateData(req PeerRequestWGForTemplate) WireGuardTemplateData {
	asnLen := len(req.ASN)
	peeringPort := req.ASN
	if asnLen > 5 {
		peeringPort = req.ASN[asnLen-5:]
	}
	return WireGuardTemplateData{
		"peering_port":       peeringPort,
		"my_link_local":      source.AppConfig.DN42.LocalLink,
		"my_dn42_v6":         source.AppConfig.DN42.IPv6,
		"my_dn42_v4":         source.AppConfig.DN42.IPv4,
		"peering_dn42_v4":    req.IPv4,
		"peering_public_key": req.WireGuard.PubKey,
		"peering_endpoint":   req.WireGuard.Endpoint,
	}
}

// RenderWireGuardConf 渲染 WireGuard 配置模板
func RenderWireGuardConf(data WireGuardTemplateData) (string, error) {
	tplPath := filepath.Join("template", "wireguard.conf")
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

type WireGuardTemplateData map[string]any

// RawWGShow 执行 wg show 并返回原始输出
func RawWGShow(iface string) (string, error) {
	cmd := exec.Command("wg", "show", iface)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run wg show %s: %w", iface, err)
	}
	return string(out), nil
}
