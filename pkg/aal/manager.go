package aal

import (
	"aalsystem/pkg/fuseki"
	"aalsystem/pkg/homeassistant"
	"aalsystem/pkg/utils"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Manager is the manager of the AAL system.
type Manager struct {
	hass   *homeassistant.Client
	sparql *fuseki.Client
	logger *logrus.Logger

	sensors map[string]*Sensor // map of sensors, indexed by sensor ID
}

// NewManager creates a new AAL system manager.
func NewManager(hass *homeassistant.Client, sparql *fuseki.Client, logger *logrus.Logger) *Manager {
	return &Manager{
		hass:   hass,
		sparql: sparql,
		logger: logger,
	}
}

// AddSensor adds a sensor to the AAL system.
func (m *Manager) AddSensor(sensor *Sensor) {
	m.sensors[sensor.ID] = sensor
}

// Run starts the AAL system.
// It connects to Home Assistant and SPARQL, and starts listening to Home Assistant events.
func (m *Manager) Run() error {
	m.logger.Info("Starting AAL System")

	err := m.hass.InitWebsocket()
	if err != nil {
		return fmt.Errorf("failed to inititalize Home Assistant websocket connection: %w", err)
	}
	events, err := m.hass.ListenEvents()
	if err != nil {
		return fmt.Errorf("failed to listen to Home Assistant events: %w", err)
	}
	for {
		event, ok := <-events
		if !ok {
			return fmt.Errorf("home Assistant events channel closed")
		}
		m.logger.Debugf("- Got event: %s", utils.Prettyfy(event))
		m.handleStateChangeEvent(&event)
	}
}

// handleStateChangeEvent handles a state change event from Home Assistant.
func (m *Manager) handleStateChangeEvent(event *homeassistant.Event) {
	sensor := m.sensors[event.EntityID]
	obs := Observation{
		ID:       fmt.Sprint(event.ID),
		SensorID: event.EntityID,
		Value:    fmt.Sprint(event.Attributes[sensor.ValueField]),
		Unit:     sensor.Unit,
	}

	query := obs.InsertQuery()

	res, err := m.sparql.Query(string(query))
	if err != nil {
		m.logger.Errorf("Failed to insert observation: %s", err)
	}
	m.logger.Debugf("Inserted observation: %s", utils.Prettyfy(res))
}
