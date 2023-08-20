package aal

import (
	"aalsystem/pkg/fuseki"
	"aalsystem/pkg/homeassistant"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/sirupsen/logrus"
)

// Manager is the manager of the AAL system.
type Manager struct {
	hass   *homeassistant.Client
	sparql *fuseki.Client
	logger *logrus.Logger

	sensors map[string]Sensor // home assistant entity -> ontology sensor name

	currentFindings   []Finding
	findingsChan      chan *Finding
	observationChan   chan *Observation
	receivedEventTime time.Time
}

// NewManager creates a new AAL system manager.
func NewManager(hass *homeassistant.Client, sparql *fuseki.Client, logger *logrus.Logger) *Manager {
	return &Manager{
		hass:            hass,
		sparql:          sparql,
		logger:          logger,
		sensors:         map[string]Sensor{},
		findingsChan:    make(chan *Finding),
		observationChan: make(chan *Observation),
	}
}

// AddSensor adds a new sensor to the AAL system.
func (m *Manager) AddSensor(entityID, sensorID string) {
	m.logger.Debugf("Adding sensor: %s (%s)", sensorID, entityID)

	query := GetSensorQuery(sensorID)
	res, err := m.sparql.Query(string(query))
	if err != nil {
		m.logger.Errorf("error getting sensor: query=%v", query)
		return
	}
	if len(res.Results.Bindings) == 0 {
		m.logger.Errorf("sensor not found in ontology: %s", sensorID)
		return
	}

	m.logger.Debugf("Sensor from ontology: %v", res.Results.Bindings[0])

	sensor := Sensor{
		ID:                 sensorID,
		ObservableProperty: res.Results.Bindings[0]["observableProperty"].Value,
		InstalledAt:        URI(res.Results.Bindings[0]["installedAt"].Value),
	}

	m.sensors[entityID] = sensor

	m.logger.Debugf("Added sensor: %v", sensor)
}

// Run starts the AAL system.
// It connects to Home Assistant and SPARQL, and starts listening to Home Assistant events.
func (m *Manager) Run() error {
	m.logger.Info("Starting AAL System")

	initMetrics()

	server := NewServer(m.logger, m.findingsChan, m.observationChan)
	go server.Run()

	// Initialize Home Assistant websocket connection
	err := m.hass.InitWebsocket()
	if err != nil {
		return fmt.Errorf("failed to inititalize Home Assistant websocket connection: %w", err)
	}

	err = retry.Do(func() error {
		if !m.sparql.IsConnected() {
			return fmt.Errorf("failed to connect to SPARQL server")
		}
		return nil
	},
		retry.Attempts(20),
		retry.Delay(1),
		retry.OnRetry(func(n uint, err error) {
			m.logger.Warnf("Failed to connect to SPARQL server. Retrying (%d/%d)", n+1, 20)
		}),
	)
	if err != nil {
		m.logger.Fatal("Failed to connect to SPARQL server")
	}
	m.logger.Info("Connected to SPARQL")

	// real sensors
	// m.AddSensor("sensor.emfitqs_000ebc_heart_rate", "emfit_heartrate")
	// m.AddSensor("sensor.emfitqs_000ebc_respiratory_rate", "emfit_breathrate")
	// m.AddSensor("binary_sensor.emfitqs_000ebc_bed_presence", "emfit_bedpresence")

	// simulated sensors
	m.AddSensor("sensor.relative_humidity", "bedroom_humidity")
	m.AddSensor("sensor.heart_rate", "emfit_heartrate")
	m.AddSensor("sensor.breath_rate", "emfit_breathrate")

	// Listen to Home Assistant events via websocket
	events, err := m.hass.ListenEvents()
	if err != nil {
		return fmt.Errorf("failed to listen to Home Assistant events: %w", err)
	}
	m.logger.Info("Listening to Home Assistant events")

	// Load initial state of all sensors via REST
	err = m.loadInitialStates()
	if err != nil {
		return fmt.Errorf("failed to get initial state of sensors: %w", err)
	}
	m.logger.Info("Got initial state of all sensors")

	// Handle Home Assistant events
	m.logger.Debug("Waiting new Home Assistant events")
	for {
		event, ok := <-events
		if !ok {
			return fmt.Errorf("home Assistant events channel closed")
		}
		m.logger.Debugf("- Got event: %s", event.ShortString())
		err := m.handleStateChangeEvent(&event)
		if err != nil {
			m.logger.Errorf("failed to handle state change event: %s", err.Error())
		}
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
		err = m.handleStateChangeEvent(event)
		if err != nil {
			return fmt.Errorf("failed to handle initial state: %w", err)
		}
	}
	return nil
}

// handleStateChangeEvent handles a state change event from Home Assistant.
func (m *Manager) handleStateChangeEvent(event *homeassistant.Event) error {
	sensor, ok := m.sensors[event.EntityID]
	if !ok {
		m.logger.Debugf("Sensor %s not found", event.EntityID)
		return nil
	}
	if event.EntityID == "sensor.heart_rate" {
		m.receivedEventTime = time.Now()
	}
	obs := Observation{
		Name:      sensor.ObservableProperty,
		Sensor:    sensor.ID,
		Value:     event.State,
		Timestamp: time.Now(),
	}
	if event.EntityID == "sensor.heart_rate" {
		duration := time.Since(m.receivedEventTime)
		eventToObservationDuration.WithLabelValues(obs.Name).Observe(duration.Seconds())
	}
	observationsCount.WithLabelValues(obs.Name).Inc()
	m.observationChan <- &obs

	// Insert observation into ontology
	m.logger.Debugf("Inserting observation: %s (%s)", obs.Sensor, obs.Value)
	err := m.insertObservation(&obs)
	if err != nil {
		return fmt.Errorf("failed to insert observation: %w", err)
	}
	m.logger.Infof("Inserted observation: sensor=%s value=%s", obs.Sensor, obs.Value)

	// Check finding
	m.logger.Trace("Checking findings activated by rules")
	go func() {
		err = m.checkFindings()
		if err != nil {
			m.logger.Errorf("failed to check findings: %s", err.Error())
		}
	}()
	return nil
}

// insertObservation inserts an observation into the ontology.
func (m *Manager) insertObservation(obs *Observation) error {
	query := obs.InsertQuery()

	err := m.sparql.Update(string(query))
	if err != nil {
		m.logger.Errorf("error inserting observation: query=%v", query)
		return err
	}
	return nil
}

// checkFindings checks if any finding was inferred by the SWRL rules.
func (m *Manager) checkFindings() error {
	query := findingsQuery()
	res, err := m.sparql.Query(query)
	if err != nil {
		return err
	}

	m.logger.Tracef("Got %d findings", len(res.Results.Bindings))
	findings := resultToFindings(*res)
	if len(findings) == 0 {
		m.logger.Debugf("No findings found")
		return nil
	}

	for _, finding := range findings {
		if finding.IsDuplicate(m.currentFindings) {
			continue
		}
		m.logger.Infof("++ New finding: %s has %s (%s)", finding.Patient, finding.Name, finding.Value)
		if finding.Name == "Tachycardia" {
			duration := time.Since(m.receivedEventTime)
			eventToFindingDuration.WithLabelValues(finding.Name).Observe(duration.Seconds())
		}
		findingsCount.WithLabelValues(finding.Name).Inc()
		m.findingsChan <- &finding
	}
	m.currentFindings = findings

	return nil
}
