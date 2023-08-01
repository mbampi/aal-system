package aal

import (
	"aalsystem/pkg/fuseki"
	"aalsystem/pkg/homeassistant"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Manager is the manager of the AAL system.
type Manager struct {
	hass   *homeassistant.Client
	sparql *fuseki.Client
	logger *logrus.Logger

	sensors map[string]string // home assistant entity -> ontology sensor name

	observationID int
	obsLock       sync.Mutex
}

// NewManager creates a new AAL system manager.
func NewManager(hass *homeassistant.Client, sparql *fuseki.Client, logger *logrus.Logger) *Manager {
	return &Manager{
		hass:          hass,
		sparql:        sparql,
		logger:        logger,
		sensors:       map[string]string{},
		observationID: 0,
	}
}

// ObsID returns a new observation ID.
// It is thread-safe.
func (m *Manager) ObsID() int {
	m.obsLock.Lock()
	id := m.observationID
	m.observationID += 1
	m.obsLock.Unlock()

	return id
}

// AddSensor adds a new sensor to the AAL system.
func (m *Manager) AddSensor(entityID, sensorID string) {
	m.logger.Debugf("Adding sensor: %s (%s)", sensorID, entityID)
	m.sensors[entityID] = sensorID
}

// Run starts the AAL system.
// It connects to Home Assistant and SPARQL, and starts listening to Home Assistant events.
func (m *Manager) Run() error {
	m.logger.Info("Starting AAL System")

	// Initialize Home Assistant websocket connection
	err := m.hass.InitWebsocket()
	if err != nil {
		return fmt.Errorf("failed to inititalize Home Assistant websocket connection: %w", err)
	}

	// Load initial state of all sensors via REST
	err = m.loadInitialStates()
	if err != nil {
		return fmt.Errorf("failed to get initial state of sensors: %w", err)
	}
	m.logger.Info("Got initial state of all sensors")

	// Listen to Home Assistant events via websocket
	events, err := m.hass.ListenEvents()
	if err != nil {
		return fmt.Errorf("failed to listen to Home Assistant events: %w", err)
	}

	// Handle Home Assistant events
	for {
		event, ok := <-events
		if !ok {
			return fmt.Errorf("home Assistant events channel closed")
		}
		m.logger.Debugf("- Got event: %s", event.ShortString())
		m.handleStateChangeEvent(&event)
	}
}

// loadInitialStates loads the initial state of all sensors.
func (m *Manager) loadInitialStates() error {
	m.logger.Debug("Getting initial state of all sensors")
	states, err := m.hass.GetStates()
	if err != nil {
		return fmt.Errorf("failed to get initial state of all sensors: %w", err)
	}
	for _, state := range states {
		event := homeassistant.EventFromState(state)
		m.logger.Debugf("- Initial state: %s", event.ShortString())
		err = m.handleStateChangeEvent(event)
		if err != nil {
			return fmt.Errorf("failed to handle initial state: %w", err)
		}
	}
	return nil
}

// handleStateChangeEvent handles a state change event from Home Assistant.
func (m *Manager) handleStateChangeEvent(event *homeassistant.Event) error {
	sensorID, ok := m.sensors[event.EntityID]
	if !ok {
		m.logger.Debugf("Sensor %s not found", event.EntityID)
		return nil
	}
	obs := Observation{
		ID:        fmt.Sprintf("%d", m.ObsID()),
		SensorID:  sensorID,
		Value:     event.State,
		Timestamp: time.Now(),
	}

	// Insert observation into ontology
	m.logger.Debugf("Inserting observation: %s (%s)", obs.SensorID, obs.Value)
	err := m.insertObservation(&obs)
	if err != nil {
		return fmt.Errorf("failed to insert observation: %w", err)
	}
	m.logger.Debugf("Inserted observation: %s (%s)", obs.SensorID, obs.Value)

	// Check trigger
	m.logger.Trace("Checking triggers activated by rules")
	err = m.checkTriggers()
	if err != nil {
		return fmt.Errorf("failed to check triggers: %w", err)
	}
	m.logger.Trace("Checked triggers activated by rules")
	return nil
}

// insertObservation inserts an observation into the ontology.
func (m *Manager) insertObservation(obs *Observation) error {
	query := obs.InsertQuery()

	err := m.sparql.Update(string(query))
	if err != nil {
		return err
	}
	return nil
}

// checkTriggers checks if any trigger was activated by the SWRL rules.
func (m *Manager) checkTriggers() error {
	query := `
PREFIX sosa: <http://www.w3.org/ns/sosa/>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?s
WHERE {
  ?s rdf:type :Trigger .
}`
	res, err := m.sparql.Query(query)
	if err != nil {
		return err
	}

	bindings := res.Results.Bindings
	if len(bindings) == 0 {
		m.logger.Debugf("No triggers activated")
		return nil
	}

	m.logger.Debugf("Got triggers: %v", bindings)
	// for _, binding := range res.Results.Bindings {
	// 	m.logger.Debugf("+ Action triggered: %v", binding)
	// 	// TODO: run action/alarm
	// }

	err = m.DeleteTriggers()
	if err != nil {
		return fmt.Errorf("failed to delete triggers: %w", err)
	}
	m.logger.Debugf("Deleted triggers")

	return nil
}

// DeleteTriggers deletes all triggers from the ontology.
func (m *Manager) DeleteTriggers() error {
	query := `
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

DELETE WHERE {
	  ?s rdf:type :Trigger .
}`

	err := m.sparql.Update(query)
	if err != nil {
		return err
	}
	return nil
}
