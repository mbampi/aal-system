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
}

// NewManager creates a new AAL system manager.
func NewManager(hass *homeassistant.Client, sparql *fuseki.Client, logger *logrus.Logger) *Manager {
	return &Manager{
		hass:   hass,
		sparql: sparql,
		logger: logger,
	}
}

// Run starts the AAL system.
func (m *Manager) Run() error {
	m.logger.Info("Starting AAL System")
	err := m.hass.InitWebsocket()
	if err != nil {
		return fmt.Errorf("failed to init Home Assistant websocket: %w", err)
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
	}
	return nil
}
