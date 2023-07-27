package aal

import (
	"aalsystem/pkg/fuseki"
	"aalsystem/pkg/homeassistant"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Manager is the manager of the AAL system.
type Manager struct {
	hass   *homeassistant.Client
	sparql *fuseki.Client
	logger *logrus.Logger

	sensors map[string]string // home assistant entity id -> ontology sensor id
}

// NewManager creates a new AAL system manager.
func NewManager(hass *homeassistant.Client, sparql *fuseki.Client, logger *logrus.Logger) *Manager {
	return &Manager{
		hass:   hass,
		sparql: sparql,
		logger: logger,
	}
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

	// Load initial state of all sensors
	err := m.loadInitialStates()
	if err != nil {
		return fmt.Errorf("failed to get initial state of sensors: %w", err)
	}
	m.logger.Info("Got initial state of all sensors")

	// Listen to Home Assistant events
	err = m.hass.InitWebsocket()
	if err != nil {
		return fmt.Errorf("failed to inititalize Home Assistant websocket connection: %w", err)
	}
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
		// m.handleStateChangeEvent(&event)
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
		m.handleStateChangeEvent(event)
	}
	return nil
}

// handleStateChangeEvent handles a state change event from Home Assistant.
func (m *Manager) handleStateChangeEvent(event *homeassistant.Event) {
	sensorID, ok := m.sensors[event.EntityID]
	if !ok {
		m.logger.Debugf("Sensor %s not found", event.EntityID)
		return
	}

	obs := Observation{
		ID:        fmt.Sprint(event.ID),
		SensorID:  sensorID,
		Value:     event.State,
		Timestamp: time.Now(),
	}

	query := obs.InsertQuery()
	m.logger.Debugf("Inserting observation: %s", query)

	err := m.sparql.Update(string(query))
	if err != nil {
		m.logger.Errorf("Failed to insert observation: %s", err)
		return
	}
	m.logger.Debugf("Inserted observation")
}
