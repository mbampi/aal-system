package main

import (
	"aalsystem/pkg/homeassistant"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting AAL System")

	logger.Debug("Getting Home Assistant token")
	accessToken := os.Getenv("HASSIO_TOKEN")
	if accessToken == "" {
		logger.Fatal("Failed to get Home Assistant token")
	}

	hass := homeassistant.NewClient()
	hass.SetLongLivedAccessToken(accessToken)
	if hass.IsConnected() {
		logger.Info("Connected to Home Assistant")
	} else {
		logger.Error("Failed to connect to Home Assistant")
	}
}
