package aal

import (
	"aalsystem/pkg/fuseki"
	"aalsystem/pkg/homeassistant"

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
