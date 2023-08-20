package homeassistant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Event represents a Home Assistant event.
type Event struct {
	ID          int            `json:"id"`
	EntityID    string         `json:"entity_id"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastChanged time.Time      `json:"time_fired"`
}

func (e *Event) ShortString() string {
	return fmt.Sprintf("%s: %s", e.EntityID, e.State)
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
				c.logger.Errorf("error reading event: %s", err)
				// if the connection is closed, try to reconnect
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					c.logger.Warn("websocket connection closed, trying to reconnect")
					err = c.InitWebsocket()
					if err != nil {
						c.logger.Errorf("failed to reconnect to websocket: %s", err)
						return
					}
					c.logger.Info("reconnected to websocket")
					continue
				} else {
					return
				}
			}
			events <- wsEventToEvent(event)
		}
	}()
	return events, nil
}

// wsEventToEvent converts a websocket event to a Event.
func wsEventToEvent(event wsEvent) Event {
	return Event{
		ID:          event.ID,
		EntityID:    event.Event.Data.EntityID,
		State:       event.Event.Data.NewState.State,
		Attributes:  event.Event.Data.NewState.Attributes,
		LastChanged: event.Event.Data.NewState.LastChanged,
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

// EventFromState creates an event from a state.
func EventFromState(state *EntityState) *Event {
	return &Event{
		ID:          0,
		EntityID:    state.EntityID,
		Attributes:  state.Attributes,
		State:       state.State,
		LastChanged: state.LastChanged,
	}
}
