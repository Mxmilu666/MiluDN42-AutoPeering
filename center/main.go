package main

import (
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/logger"
)

func main() {
	l := logger.InitLogger()
	if l == nil {
		panic("Failed to initialize logger")
	}

	logger.Info("Nya!,MiluDN42-AutoPeering-Center")

}
