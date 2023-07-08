package homeassistant

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Event represents a Home Assistant event.
type Event struct {
	ID         int            `json:"id"`
	EntityID   string         `json:"entity_id"`
	Attributes map[string]any `json:"attributes"`
	TimeFired  time.Time      `json:"time_fired"`
}

type wsEvent struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Event struct {
		Context struct {
			ID       string      `json:"id"`
			ParentID interface{} `json:"parent_id"`
			UserID   interface{} `json:"user_id"`
		} `json:"context"`
		Data struct {
			EntityID string `json:"entity_id"`
			NewState struct {
				Attributes map[string]any `json:"attributes"`
				Context    struct {
					ID       string      `json:"id"`
					ParentID interface{} `json:"parent_id"`
					UserID   interface{} `json:"user_id"`
				} `json:"context"`
				EntityID    string    `json:"entity_id"`
				LastChanged time.Time `json:"last_changed"`
				LastUpdated time.Time `json:"last_updated"`
				State       string    `json:"state"`
			} `json:"new_state"`
			OldState struct {
				Attributes map[string]any `json:"attributes"`
				Context    struct {
					ID       string      `json:"id"`
					ParentID interface{} `json:"parent_id"`
					UserID   interface{} `json:"user_id"`
				} `json:"context"`
				EntityID    string    `json:"entity_id"`
				LastChanged time.Time `json:"last_changed"`
				LastUpdated time.Time `json:"last_updated"`
				State       string    `json:"state"`
			} `json:"old_state"`
		} `json:"data"`
		EventType string    `json:"event_type"`
		Origin    string    `json:"origin"`
		TimeFired time.Time `json:"time_fired"`
	} `json:"event"`
}

// ListenEvents listens for events.
func (c *Client) ListenEvents() (<-chan Event, error) {
	if c.wsConn == nil {
		return nil, fmt.Errorf("websocket connection not established")
	}
	err := c.wsConn.WriteJSON(map[string]any{
		"id":         1,
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
		return nil, fmt.Errorf("failed to subscribe to events: %v", result)
	}

	// read events and send them to the events channel until the connection is closed.
	events := make(chan Event)
	go func() {
		defer close(events)
		for {
			var event wsEvent
			err := c.wsConn.ReadJSON(&event)
			if err != nil {
				log.Printf("error reading event: %s", err)
				return
			}
			events <- wsEventToEvent(event)
		}
	}()
	return events, nil
}

// wsEventToEvent converts a websocket event to a Event.
func wsEventToEvent(event wsEvent) Event {
	return Event{
		ID:         event.ID,
		EntityID:   event.Event.Data.EntityID,
		Attributes: event.Event.Data.NewState.Attributes,
		TimeFired:  event.Event.TimeFired,
	}
}

type restEvent struct {
	Event         string `json:"event"`
	ListenerCount int    `json:"listener_count"`
}

// GetEvents returns the events.
func (c *Client) GetEvents() ([]*restEvent, error) {
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
	var events []*restEvent
	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
