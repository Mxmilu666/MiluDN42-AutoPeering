package main

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/logger"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/server"
)

func main() {
	l := logger.InitLogger()
	if l == nil {
		panic("Failed to initialize logger")
	}

	logger.Info("Nya!,MiluDN42-AutoPeering-Node")

	configPath := "config.yaml"
	_, err := source.LoadConfig(configPath)
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}

	server.Setupserver()
}
