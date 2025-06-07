package main

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/logger"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/server"
)

func main() {
	l := logger.InitLogger()
	if l == nil {
		panic("Failed to initialize logger")
	}

	logger.Info("Nya!,MiluDN42-AutoPeering-Center")

	configPath := "config.yaml"
	_, err := source.LoadConfig(configPath)
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}

	server.Setupserver()
}
