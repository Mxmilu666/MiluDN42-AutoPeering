package helper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source"
)

// SetupConfFiles 会在 Bird 和 WireGuard 配置都渲染完毕后，将内容写入目录
func SetupConfFiles(asn string, birdConf string, wgConf string) error {
	birdDir := source.AppConfig.Bird.PeerConfPath
	wgDir := source.AppConfig.Wireguard.ConfigPath
	filename := fmt.Sprintf("dn42_%s.conf", asn)

	birdPath := filepath.Join(birdDir, filename)
	wgPath := filepath.Join(wgDir, filename)

	// 如果文件已存在则不覆盖
	if _, err := os.Stat(birdPath); err == nil {
		return fmt.Errorf("bird conf already exists")
	}
	if _, err := os.Stat(wgPath); err == nil {
		return fmt.Errorf("wireguard conf already exists")
	}

	if err := os.WriteFile(birdPath, []byte(birdConf), 0644); err != nil {
		return fmt.Errorf("write bird conf: %w", err)
	}
	if err := os.WriteFile(wgPath, []byte(wgConf), 0644); err != nil {
		return fmt.Errorf("write wireguard conf: %w", err)
	}
	return nil
}

// RunWgQuickAndBirdc 执行 wg-quick up 和 birdc c 命令，并返回错误信息（如有）
func RunWgQuickAndBirdc(asn string) error {
	wgConf := fmt.Sprintf("dn42_%s.conf", asn)
	wgConfPath := filepath.Join(source.AppConfig.Wireguard.ConfigPath, wgConf)

	// 执行 wg-quick up
	wgCmd := exec.Command("wg-quick", "up", wgConfPath)
	wgOut, wgErr := wgCmd.CombinedOutput()
	if wgErr != nil {
		return fmt.Errorf("wg-quick up failed: %v, output: %s", wgErr, string(wgOut))
	}

	// 执行 birdc c
	birdcCmd := exec.Command("birdc", "c")
	birdcOut, birdcErr := birdcCmd.CombinedOutput()
	if birdcErr != nil {
		return fmt.Errorf("birdc c failed: %v, output: %s", birdcErr, string(birdcOut))
	}

	return nil
}
