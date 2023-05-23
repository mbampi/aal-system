package main

import (
	"aalsystem/pkg/homeassistant"
	"aalsystem/pkg/sparql"
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
		logger.Fatal("Failed to connect to Home Assistant")
	}
	logger.Info("Connected to Home Assistant")

	logger.Debug("Connecting to SPARQL")
	ds := "sno"
	sparqlServer := sparql.NewClient(ds)
	if !sparqlServer.IsConnected() {
		logger.Fatal("Failed to connect to SPARQL")
	}
	logger.Info("Connected to SPARQL")

	logger.Debug("Getting Home Assistant states")
	states, err := hass.GetStates()
	if err != nil {
		logger.Fatal("Failed to get Home Assistant states:", err)
	}
	logger.Info("Got Home Assistant states")
	for _, state := range states {
		logger.Debugf(" - %s : %s : %s", state.EntityID, state.State, state.Attributes["friendly_name"])
	}

	logger.Debug("Getting SPARQL stats")
	stats, err := sparqlServer.GetStats()
	if err != nil {
		logger.Fatal("Failed to get SPARQL stats:", err)
	}
	logger.Info("Got SPARQL stats")
	logger.Debugf(" - %v", stats)

	logger.Debug("Doing SPARQL Query")
	query := `PREFIX owl: <http://www.w3.org/2002/07/owl%23>
	PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema%23>
	SELECT DISTINCT ?class ?label ?description
	WHERE {
	  ?class a owl:Class.
	  OPTIONAL { ?class rdfs:label ?label}
	  OPTIONAL { ?class rdfs:comment ?description}
	}
	LIMIT 10`

	results, err := sparqlServer.Query(query)
	if err != nil {
		logger.Fatal("Failed to do SPARQL query:", err)
	}
	logger.Info("Got SPARQL results")
	logger.Debugf(" - %v", results)
}
