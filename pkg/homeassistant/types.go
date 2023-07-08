package homeassistant

type websocketResult struct {
	ID      int            `json:"id"`
	Type    string         `json:"type"`
	Success bool           `json:"success"`
	Result  map[string]any `json:"result"`
}
