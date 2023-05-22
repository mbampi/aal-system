package homeassistant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Service struct {
	Domain   string           `json:"domain"`
	Services []map[string]any `json:"services"`
}

// GetServices returns the services.
func (c *Client) GetServices() ([]*Service, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    c.baseURL.ResolveReference(&url.URL{Path: "/api/services"}),
		Header: c.getHeaders(),
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get services: %s", res.Status)
	}
	defer res.Body.Close()
	var services []*Service
	err = json.NewDecoder(res.Body).Decode(&services)
	if err != nil {
		return nil, err
	}
	return services, nil
}
