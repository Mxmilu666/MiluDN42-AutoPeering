package main

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source/logger"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source/server"
)

func main() {
	l := logger.InitLogger()
	if l == nil {
		panic("Failed to initialize logger")
	}

	logger.Info("Nya!,MiluDN42-AutoPeering-Backend")

	configPath := "config.yaml"
	_, err := source.LoadConfig(configPath)
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return
	}

	server.Setupserver()
}
