package main

import (
	"aalsystem/pkg/aal"
	"aalsystem/pkg/fuseki"
	"aalsystem/pkg/homeassistant"
	"aalsystem/pkg/utils"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.Info("Starting AAL System")

	// Environment variables
	logger.Debug("Getting Home Assistant token")
	accessToken := os.Getenv("HASSIO_TOKEN")
	if accessToken == "" {
		logger.Fatal("Failed to get Home Assistant token")
	}

	// Home Assistant
	logger.Debug("Connecting to Home Assistant")
	hass := homeassistant.NewClient()
	hass.SetLogger(logger)
	hass.SetLongLivedAccessToken(accessToken)
	if !hass.IsConnected() {
		logger.Warnf("Is your HomeAssistant also connected to network \"%s\" ?", utils.GetWIFIName())
		logger.Fatal("Failed to connect to Home Assistant")
	}
	logger.Info("Connected to Home Assistant")

	// SPARQL
	logger.Debug("Connecting to SPARQL")
	ds := "med"
	sparqlServer := fuseki.NewClient(ds)

	// AAL System
	aalManager := aal.NewManager(hass, sparqlServer, logger)

	err := aalManager.Run()
	if err != nil {
		logger.Fatal("Failed to run AAL system:", err)
	}
}
