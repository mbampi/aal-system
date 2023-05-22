package homeassistant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Event struct {
	Event         string `json:"event"`
	ListenerCount int    `json:"listener_count"`
}

// GetEvents returns the events.
func (c *Client) GetEvents() ([]*Event, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    c.baseURL.ResolveReference(&url.URL{Path: "/api/events"}),
		Header: c.getHeaders(),
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get events: %s", res.Status)
	}
	defer res.Body.Close()
	var events []*Event
	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
