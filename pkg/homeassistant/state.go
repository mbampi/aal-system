package homeassistant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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
	req := &http.Request{
		Method: http.MethodGet,
		URL:    c.baseURL.ResolveReference(&url.URL{Path: "/api/states"}),
		Header: c.getHeaders(),
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get states: %s", res.Status)
	}
	defer res.Body.Close()
	var states []*EntityState
	err = json.NewDecoder(res.Body).Decode(&states)
	if err != nil {
		return nil, err
	}
	return states, nil
}
