package homeassistant

import "time"

type EntityState struct {
	EntityID    string         `json:"entity_id"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastChanged time.Time      `json:"last_changed"`
	LastUpdated time.Time      `json:"last_updated"`
	Context     struct {
		ID       string `json:"id"`
		ParentID string `json:"parent_id"`
		UserID   string `json:"user_id"`
	} `json:"context"`
}

// GetStates returns the states of all entities.
func (c *Client) GetStates() ([]*EntityState, error) {
	return nil, nil
}

// GetState returns the state of an entity.
func (c *Client) GetState(entityID string) (*EntityState, error) {
	return nil, nil
}

// SetState sets the state of an entity.
func (c *Client) SetState(entityID string, state string) error {
	return nil
}
