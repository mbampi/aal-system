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

// ListenEvents listens for events.
func (c *Client) ListenEvents() (<-chan Event, error) {
	if c.wsConn == nil {
		return nil, fmt.Errorf("websocket connection not established")
	}
	err := c.wsConn.WriteJSON(map[string]string{
		"id":         "1",
		"type":       "subscribe_events",
		"event_type": "state_changed",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to events: %s", err)
	}

	var result websocketResult
	err = c.wsConn.ReadJSON(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to events: %s", err)
	}
	if !result.Success {
		return nil, fmt.Errorf("failed to subscribe to events: %s", result.Result)
	}

	events := make(chan Event)
	go func() {
		defer close(events)
		for {
			var event Event
			err := c.wsConn.ReadJSON(&event)
			if err != nil {
				return
			}
			events <- event
		}
	}()
	return events, nil
}
