package main

import (
	"fmt"

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
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	server.Setupserver()
}
