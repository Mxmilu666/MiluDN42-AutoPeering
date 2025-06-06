package helper

import (
	"fmt"
	"os/exec"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/logger"
)

// RestartTunnelAndBird 重启指定 ASN 的 WireGuard 隧道和 BIRD 会话
func RestartTunnelAndBird(asn string) (string, any, error) {
	iface := fmt.Sprintf("dn42_%s", asn)

	// 先 down
	downCmd := exec.Command("wg-quick", "down", iface)
	downOut, downErr := downCmd.CombinedOutput()
	if downErr != nil {
		logger.Warn("wg-quick down failed", "error", downErr, "output", string(downOut))
	}

	// 再 up
	upCmd := exec.Command("wg-quick", "up", iface)
	upOut, upErr := upCmd.CombinedOutput()
	if upErr != nil {
		return "", nil, fmt.Errorf("wg-quick up failed: %v, output: %s", upErr, string(upOut))
	}

	// 获取所有 BGP peer 名字
	peers, err := GetBGPPeerNames()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get BGP peers: %w", err)
	}
	// 模糊查找对应 peer 名字
	peerName, found := FindPeerNameByFuzzy(peers, asn)
	var output any = nil
	if found {
		birdcCmd := exec.Command("birdc", "restart", peerName)
		birdcOut, birdcErr := birdcCmd.CombinedOutput()
		if birdcErr != nil {
			return peerName, output, fmt.Errorf("birdc c failed: %v, output: %s", birdcErr, string(birdcOut))
		}
	}

	return peerName, output, nil
}
