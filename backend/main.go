package main

import (
	"aalsystem/pkg/aal"
	"aalsystem/pkg/fuseki"
	"aalsystem/pkg/homeassistant"
	"aalsystem/pkg/utils"
	"fmt"
	"os"

	"github.com/avast/retry-go"
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
	err := retry.Do(func() error {
		if !sparqlServer.IsConnected() {
			return fmt.Errorf("failed to connect to SPARQL server")
		}
		return nil
	},
		retry.Attempts(20),
		retry.Delay(1),
		retry.OnRetry(func(n uint, err error) {
			logger.Warnf("Failed to connect to SPARQL server. Retrying (%d/%d)", n+1, 20)
		}),
	)
	if err != nil {
		logger.Fatal("Failed to connect to SPARQL server")
	}

	logger.Info("Connected to SPARQL")

	// AAL System
	aalManager := aal.NewManager(hass, sparqlServer, logger)

	// real sensors
	// aalManager.AddSensor("sensor.emfitqs_000ebc_heart_rate", "emfit_heartrate")
	// aalManager.AddSensor("sensor.emfitqs_000ebc_respiratory_rate", "emfit_breathrate")
	// aalManager.AddSensor("binary_sensor.emfitqs_000ebc_bed_presence", "emfit_bedpresence")

	// simulated sensors
	aalManager.AddSensor("sensor.relative_humidity", "bedroom_humidity")
	aalManager.AddSensor("sensor.heart_rate", "emfit_heartrate")
	aalManager.AddSensor("sensor.breath_rate", "emfit_breathrate")

	err = aalManager.Run()
	if err != nil {
		logger.Fatal("Failed to run AAL system:", err)
	}
}
