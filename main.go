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
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting AAL System")

	logger.Debug("Getting Home Assistant token")
	accessToken := os.Getenv("HASSIO_TOKEN")
	if accessToken == "" {
		logger.Fatal("Failed to get Home Assistant token")
	}

	logger.Debug("Connecting to Home Assistant")
	hass := homeassistant.NewClient()
	hass.SetLongLivedAccessToken(accessToken)
	if !hass.IsConnected() {
		logger.Warnf("Is your HomeAssistant also connected to network \"%s\" ?", utils.GetWIFIName())
		logger.Fatal("Failed to connect to Home Assistant")
	}
	logger.Info("Connected to Home Assistant")

	logger.Debug("Connecting to SPARQL")
	ds := "med"
	sparqlServer := fuseki.NewClient(ds)
	if !sparqlServer.IsConnected() {
		logger.Fatal("Failed to connect to SPARQL")
	}
	logger.Info("Connected to SPARQL")

	aalManager := aal.NewManager(hass, sparqlServer, logger)
	err := aalManager.Run()
	if err != nil {
		logger.Fatal("Failed to run AAL system:", err)
	}

	// logger.Debug("Getting Home Assistant states")
	// states, err := hass.GetStates()
	// if err != nil {
	// 	logger.Fatal("Failed to get Home Assistant states:", err)
	// }
	// logger.Info("Got Home Assistant states")
	// for _, state := range states {
	// 	logger.Debugf(" - %s : %s : %s", state.EntityID, state.State, state.Attributes["friendly_name"])
	// }

	// logger.Debug("Getting SPARQL status")
	// status, err := sparqlServer.GetStatus()
	// if err != nil {
	// 	logger.Fatal("Failed to get SPARQL status:", err)
	// }
	// logger.Info("Got SPARQL status")
	// logger.Debugf("%s", prettyfy(status))

	// logger.Debug("Doing SPARQL Query")
	// query := snomed.SubclassesOf(snomed.BodyTemperature)
	// results, err := sparqlServer.Query(query)
	// if err != nil {
	// 	logger.Fatal("Failed to do SPARQL query:", err)
	// }
	// logger.Info("Got SPARQL results")
	// logger.Debugf("%s", prettyfy(results))
}
