package helper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/logger"
)

// RemoveTunnelAndBird 删除指定 ASN 的 WireGuard 隧道和 BIRD 会话
func RemoveTunnelAndBird(asn string) error {
	iface := fmt.Sprintf("dn42_%s", asn)

	// 删除 Bird 配置文件
	birdConfPath := filepath.Join(source.AppConfig.Bird.PeerConfPath, iface+".conf")
	if err := os.Remove(birdConfPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove bird conf: %w", err)
	}

	// 关闭 WireGuard 隧道
	downCmd := exec.Command("wg-quick", "down", iface)
	downOut, downErr := downCmd.CombinedOutput()
	if downErr != nil {
		logger.Warn("wg-quick down failed", "error", downErr, "output", string(downOut))
	}

	// 删除 WireGuard 配置文件
	wgConfPath := filepath.Join(source.AppConfig.Wireguard.ConfigPath, iface+".conf")
	if err := os.Remove(wgConfPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove wg conf: %w", err)
	}

	// birdc c 重新加载配置
	birdcCmd := exec.Command("birdc", "c")
	birdcOut, birdcErr := birdcCmd.CombinedOutput()
	if birdcErr != nil {
		return fmt.Errorf("birdc c failed: %v, output: %s", birdcErr, string(birdcOut))
	}

	return nil
}
